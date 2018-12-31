package gowhere

import (
	"fmt"
	"regexp"
)

// Rule represents one redirect rule
type Rule struct {
	// The line of the input file where the rule was found
	LineNum int
	// The Apache directive ("redirect" or "redirectmatch")
	Directive string
	// The HTTP response code ("301", etc.)
	Code string
	// The pattern to match (a literal for "redirect" and a regexp for
	// "redirectmatch")
	Pattern string
	// The destination of the redirection. May include regexp group
	// substitutions for "redirectmatch" (e.g., "$1")
	Target string
	re     *regexp.Regexp
}

// Return a nicely formatted version of the Rule
func (r *Rule) String() string {
	return fmt.Sprintf("[line %d] %s %s %s %s",
		r.LineNum, r.Directive, r.Pattern, r.Code, r.Target)
}

// Check represents a test for one Rule
type Check struct {
	// The line of the input file where the check was found
	LineNum int
	// The input to give to the RuleSet
	Input string
	// The expected HTTP response code
	Code string
	// The expected destination of the redirection
	Expected string
}

// Match holds the values for a Rule that has matched and the
// destination of the redirect
type Match struct {
	Rule
	// The matched destination for the redirection
	Match string
}

// RuleSet holds a group of Rules to be applied together
type RuleSet struct {
	rules []Rule
}

// NewRule creates a Rule from the strings on the input line
func NewRule(lineNum int, params []string) (*Rule, error) {
	if len(params) < 3 {
		return nil, fmt.Errorf("Not enough parameters on line %d: %v",
			lineNum, params)
	}
	if len(params) > 4 {
		return nil, fmt.Errorf("Too many parameters on line %d: %v",
			lineNum, params)
	}

	r := Rule{LineNum: lineNum, Directive: params[0]}

	if len(params) == 4 {
		// redirect code pattern target
		r.Code = params[1]
		r.Pattern = params[2]
		r.Target = params[3]
	} else if params[1] == "410" {
		// The page has been deleted and is not coming
		// back (nil target).
		r.Code = params[1]
		r.Pattern = params[2]
	} else {
		// redirect pattern target
		// (code is implied)
		r.Code = "301"
		r.Pattern = params[1]
		r.Target = params[2]
	}

	// Verify that we understand the directive and compile the
	// regexp if there is one.
	switch r.Directive {
	case "redirect":
	case "redirectmatch":
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			return nil, fmt.Errorf("Could not understand regexp '%s' in rule on line %d: %v",
				r.Pattern, lineNum, params)
		}
		r.re = re
	default:
		return nil, fmt.Errorf("Could not understand dirctive '%s' in rule on line %d: %v",
			r.Directive, lineNum, params)
	}

	return &r, nil
}

// NewCheck creates a Check from the strings on the input line
func NewCheck(lineNum int, params []string) (*Check, error) {
	var t Check

	t.LineNum = lineNum

	if len(params) == 3 {
		// input code expected
		t.Input = params[0]
		t.Code = params[1]
		t.Expected = params[2]
		return &t, nil
	}

	if len(params) == 2 {
		// input code
		// (no expected redirect)
		t.Input = params[0]
		t.Code = params[1]
		return &t, nil
	}

	return nil, fmt.Errorf("Could not understand check on line %d: %v",
		lineNum, params)
}

// Match tests whether the rule matches the target string.
//
// Returns the matching string, so when the rule pattern is a regexp
// and the target includes substitutions the return value is the
// actual path to which the redirect would send the browser.
func (r *Rule) Match(target string) string {
	switch r.Directive {

	case "redirect":
		if r.Pattern == target {
			return r.Target
		}

	case "redirectmatch":
		// if the pattern matches, expand the references in the target
		// to what was matched in the input so we can return a real
		// path rather than a regexp
		result := []byte{}
		for _, submatches := range r.re.FindAllStringSubmatchIndex(target, -1) {
			result = r.re.ExpandString(result, r.Target, target, submatches)
		}
		return string(result)
	}

	return ""
}

func (rs *RuleSet) firstMatch(target string, verbose bool) *Match {
	if verbose {
		fmt.Printf("\nfirstMatch '%s'\n", target)
	}

	for _, r := range rs.rules {
		if verbose {
			fmt.Printf("checking: '%s' against %s '%s'\n", target,
				r.Directive, r.Pattern)
		}

		s := r.Match(target)
		if s != "" {
			m := Match{r, s}
			return &m
		}
	}

	return nil
}

// FindMatches locates all of the Rules that match the Check
func (rs *RuleSet) FindMatches(check *Check, settings Settings) []Match {
	var r []Match

	seen := make(map[string]bool)
	match := rs.firstMatch(check.Input, settings.Verbose)
	for {
		if match == nil {
			if settings.Verbose {
				fmt.Printf("no more matches\n")
			}
			break
		}

		if settings.Verbose {
			fmt.Printf("matched: %v\n", *match)
		}

		if seen[match.Match] {
			// cycle detected
			if settings.Verbose {
				fmt.Printf("cycle\n")
			}
			break
		}
		r = append(r, *match)
		seen[match.Match] = true

		if settings.MaxHops > 0 && len(r) > settings.MaxHops {
			if settings.Verbose {
				fmt.Printf("max hops\n")
			}
			break
		}

		if match.Match == "" {
			// a redirect that doesn't point to a path,
			// like code 410
			if settings.Verbose {
				fmt.Printf("no-target redirect\n")
			}
			break
		}

		// look for another item in a redirect chain
		match = rs.firstMatch(match.Match, settings.Verbose)
	}

	return r
}
