package gowhere

import (
	"fmt"
)

type Mismatched struct {
	rule Rule
	input string
	expected string
}

type Cycle struct {
	test RuleTest
	matches []Rule
}

type Results struct {
	// inputs that did not match the expected value
	mismatched []Mismatched
	// rules that result in too many hops
	exceeded_hops []Mismatched
	// inputs that result in redirect cycles
	cycles []Cycle
	// rules that never matched
	unmatched []Rule
}

func findMatches(rules *RuleSet, test *RuleTest, max_hops int) ([]Rule, error) {
	var r []Rule

	seen := make(map[int]bool)
	for match := rules.Match(test.input); match != nil; match = rules.Match(match.target) {
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

func ProcessTests(rules *RuleSet, tests *RuleTestSet, max_hops int) (*Results, error) {
	r := Results{}

	for _, test := range tests.tests {
		fmt.Printf("test: %v\n", test)
		matches, err := findMatches(rules, &test, max_hops)
		if err != nil {
			return &r, err
		}
		for _, m := range matches {
			fmt.Printf("match: %v\n", m)
		}
	}

	return &r, nil
}
