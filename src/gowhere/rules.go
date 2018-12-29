package gowhere

import (
	"fmt"
	"regexp"
)

type Rule struct {
	line_num int
	directive string
	code string
	pattern string
	target string
	re *regexp.Regexp
}

type RuleTest struct {
	line_num int
	input string
	code string
	expected string
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

	r.line_num = line_num
	r.directive = params[0]

	if len(params) == 4 {
		// redirect code pattern target
		r.code = params[1]
		r.pattern = params[2]
		r.target = params[3]
	} else if len(params) == 3 {
		if params[1] == "410" {
			// The page has been deleted and is not coming
			// back (nil target).
			r.code = params[1]
			r.pattern = params[2]
		} else {
			// redirect pattern target
			// (code is implied)
			r.code = "301"
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

	t.line_num = line_num

	if len(params) == 3 {
		// input code expected
		t.input = params[0]
		t.code = params[1]
		t.expected = params[2]
		return &t, nil
	}

	if len(params) == 2 {
		// input code
		// (no expected redirect)
		t.input = params[0]
		t.code = params[1]
		return &t, nil
	}

	return nil, fmt.Errorf("Could not understand test on line %d: %v",
		line_num, params)
}

func (r *Rule) Match (target string) (*Rule) {
	fmt.Printf("checking: '%s' against %s '%s'\n", target,
		r.directive, r.pattern)

	switch r.directive {

	case "redirect":
		if r.pattern == target {
			fmt.Printf("matched: %v\n", *r)
			return r
		}

	case "redirectmatch":
		match := r.re.FindStringSubmatch(target)
		if len(match) > 0 {
			fmt.Printf("matched: %v\n", *r)
			return r
		}
	}

	return nil
}

func (rs *RuleSet) firstMatch (target string) (*Rule) {

	for _, r := range rs.rules {
		if r.Match(target) != nil {
			return &r
		}
	}

	return nil
}

func (rs *RuleSet) FindMatches(test *RuleTest, max_hops int) ([]Rule, error) {
	var r []Rule

	seen := make(map[int]bool)
	for match := rs.firstMatch(test.input); match != nil; match = rs.firstMatch(match.target) {
		if len(r) > max_hops {
			break
		}
		r = append(r, *match)
		if seen[match.line_num] {
			// cycle detected
			break
		}
		seen[match.line_num] = true
		if match.target == "" {
			// a redirect that doesn't point to a path,
			// like code 410
			break
		}
	}

	return r, nil
}
