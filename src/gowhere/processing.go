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

func ProcessTests(rules *RuleSet, tests *RuleTestSet, max_hops int) (*Results, error) {
	r := Results{}

	for _, test := range tests.tests {
		fmt.Printf("test: %v\n", test)
		matches, err := rules.FindMatches(&test, max_hops)
		if err != nil {
			return &r, err
		}
		for _, m := range matches {
			fmt.Printf("match: %v\n", m)
		}
	}

	return &r, nil
}
