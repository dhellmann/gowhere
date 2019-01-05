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
