package gowhere

import (
	"fmt"
)

type Rule struct {
	line_num int
	code string
	pattern string
	target string
}

func NewRule(line_num int, params []string) (*Rule, error) {
	var r Rule

	r.line_num = line_num

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

type RuleSet struct {
	rules []Rule
}

func (rules *RuleSet) Add(line_num int, params []string) error {
	r, err := NewRule(line_num, params)
	if err != nil {
		return err
	}
	rules.rules = append(rules.rules, *r)
	return nil
}
