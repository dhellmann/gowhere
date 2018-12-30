package gowhere

import (
	"fmt"
	"regexp"
)

type Rule struct {
	LineNum   int
	directive string
	Code      string
	pattern   string
	target    string
	re        *regexp.Regexp
}

func (r *Rule) String() string {
	return fmt.Sprintf("[line %d] %s %s %s %s",
		r.LineNum, r.directive, r.pattern, r.Code, r.target)
}

type RuleTest struct {
	LineNum  int
	Input    string
	Code     string
	Expected string
}

type Match struct {
	Rule
	Match string
}

type RuleSet struct {
	rules []Rule
}

type RuleTestSet struct {
	tests []RuleTest
}

func NewRule(line_num int, params []string) (*Rule, error) {
	var r Rule

	if len(params) < 3 {
		return nil, fmt.Errorf("Not enough parameters on line %d: %v",
			line_num, params)
	}

	r.LineNum = line_num
	r.directive = params[0]

	if len(params) == 4 {
		// redirect code pattern target
		r.Code = params[1]
		r.pattern = params[2]
		r.target = params[3]
	} else if len(params) == 3 {
		if params[1] == "410" {
			// The page has been deleted and is not coming
			// back (nil target).
			r.Code = params[1]
			r.pattern = params[2]
		} else {
			// redirect pattern target
			// (code is implied)
			r.Code = "301"
			r.pattern = params[1]
		}
	} else {
		return nil, fmt.Errorf("Could not understand rule on line %d: %v",
			line_num, params)
	}

	// Verify that we understand the directive and compile the
	// regexp if there is one.
	switch r.directive {
	case "redirect":
	case "redirectmatch":
		re, err := regexp.Compile(r.pattern)
		if err != nil {
			return nil, fmt.Errorf("Could not understand regexp '%s' in rule on line %d: %v",
				r.pattern, line_num, params)
		}
		r.re = re
	default:
		return nil, fmt.Errorf("Could not understand dirctive '%s' in rule on line %d: %v",
			r.directive, line_num, params)
	}

	return &r, nil
}

func NewRuleTest(line_num int, params []string) (*RuleTest, error) {
	var t RuleTest

	t.LineNum = line_num

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

	return nil, fmt.Errorf("Could not understand test on line %d: %v",
		line_num, params)
}

func (r *Rule) Match(target string) string {
	switch r.directive {

	case "redirect":
		if r.pattern == target {
			return r.target
		}

	case "redirectmatch":
		// if the pattern matches, expand the references in the target
		// to what was matched in the input so we can return a real
		// path rather than a regexp
		result := []byte{}
		for _, submatches := range r.re.FindAllStringSubmatchIndex(target, -1) {
			result = r.re.ExpandString(result, r.target, target, submatches)
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
				r.directive, r.pattern)
		}

		s := r.Match(target)
		if s != "" {
			m := Match{r, s}
			return &m
		}
	}

	return nil
}

func (rs *RuleSet) FindMatches(test *RuleTest, settings Settings) ([]Match, error) {
	var r []Match

	seen := make(map[string]bool)
	match := rs.firstMatch(test.Input, settings.Verbose)
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

	return r, nil
}
