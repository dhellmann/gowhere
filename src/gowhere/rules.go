package gowhere

import (
	"fmt"
)

type Rule struct {
	line_num int
	directive string
	code string
	pattern string
	target string
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
		return &r, nil
	}

	if len(params) == 3 {
		if params[1] == "410" {
			// The page has been deleted and is not coming
			// back (nil target).
			r.code = params[1]
			r.pattern = params[2]
			return &r, nil
		} else {
			// redirect pattern target
			// (code is implied)
			r.code = "301"
			r.pattern = params[1]
			r.target = params[2]
			return &r, nil
		}
	}

	return nil, fmt.Errorf("Could not understand rule on line %d: %v",
		line_num, params)
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
