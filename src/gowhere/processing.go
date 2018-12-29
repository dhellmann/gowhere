package gowhere

import (
	"fmt"
)

type Mismatched struct {
	test RuleTest
	matches []Match
}

type Cycle struct {
	test RuleTest
	matches []Match
}

type Results struct {
	// inputs that did not match the expected value
	Mismatched []Mismatched
	// rules that result in too many hops
	ExceededHops []Mismatched
	// inputs that result in redirect cycles
	Cycles []Cycle
	// rules that never matched
	Unmatched []Rule
	// rules that were matched properly
	Matched []Rule
}

func ProcessTests(rules *RuleSet, tests *RuleTestSet, max_hops int) (*Results, error) {
	r := Results{}
	used := make(map[int]bool)

	for _, test := range tests.tests {
		fmt.Printf("\ntest: %v\n", test)
		matches, err := rules.FindMatches(&test, max_hops)
		if err != nil {
			return &r, err
		}
		fmt.Printf("found %d matches: %v\n", len(matches), matches)
		if len(matches) == 0 {
			if test.code == "200" {
				// The test is ensuring that a URL
				// does *not* redirect, so the test is
				// passing.
			} else {
				// The test did not match any rules,
				// so record the mismatch as having
				// not redirected.
				r.Mismatched = append(
					r.Mismatched,
					Mismatched{test, matches})
			}
		} else {
			// Record only the first match as used,
			// encouraging individual tests for each rule.
			used[matches[0].line_num] = true

			// Look for cycles, mismatches, etc.
			finalMatch := matches[len(matches)-1]
			if test.input == finalMatch.match {
				// The matches resulted in going back to
				// the starting point, so we have a cycle
				r.Cycles = append(r.Cycles,
					Cycle{test, matches})
			} else if max_hops > 0 && len(matches) > max_hops {
				// Regardless of whether we ended up
				// in the right place, it took too
				// many hops to get there.
				r.ExceededHops = append(
					r.ExceededHops,
					Mismatched{test, matches})
			} else if (test.code != finalMatch.code ||
				test.expected != finalMatch.match) {
				// There is at least one match, but
				// the final URL and code are not the
				// ones we expected.
				r.Mismatched = append(
					r.Mismatched,
					Mismatched{test, matches})
			} else {
				// Recognize that the first
				// rule was tested properly.
				used[matches[0].line_num] = true
			}
		}
	}

	for _, rule := range rules.rules {
		if used[rule.line_num] {
			r.Matched = append(r.Matched, rule)
		} else {
			r.Unmatched = append(r.Unmatched, rule)
		}
	}

	return &r, nil
}
