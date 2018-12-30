package gowhere

import (
	"fmt"
)

// Results when a Check produces unexpected matches
type Mismatched struct {
	Check    Check
	Matches []Match
}

// Processing results
type Results struct {
	// inputs that did not match the expected value
	Mismatched []Mismatched
	// rules that result in too many hops
	ExceededHops []Mismatched
	// inputs that result in redirect cycles
	Cycles []Mismatched
	// rules that never matched
	Unmatched []Rule
	// rules that were matched properly
	Matched []Rule
}

// Processing input settings
type Settings struct {
	Verbose bool
	MaxHops int
}

// Run all of the rules against the checks and produce a results set.
func ProcessChecks(rules *RuleSet, checks []Check, settings Settings) (*Results, error) {
	r := Results{}
	used := make(map[int]bool)

	for _, check := range checks {
		if settings.Verbose {
			fmt.Printf("\ncheck: %v\n", check)
		}
		matches := rules.FindMatches(&check, settings)
		if settings.Verbose {
			fmt.Printf("found %d matches: %v\n", len(matches), matches)
		}
		if len(matches) == 0 {
			if check.Code == "200" {
				// The check is ensuring that a URL
				// does *not* redirect, so the check is
				// passing.
			} else {
				// The check did not match any rules,
				// so record the mismatch as having
				// not redirected.
				r.Mismatched = append(
					r.Mismatched,
					Mismatched{check, matches})
			}
		} else {
			// Record only the first match as used,
			// encouraging individual checks for each rule.
			used[matches[0].LineNum] = true

			// Look for cycles, mismatches, etc.
			finalMatch := matches[len(matches)-1]
			if check.Input == finalMatch.Match {
				// The matches resulted in going back to
				// the starting point, so we have a cycle
				r.Cycles = append(r.Cycles,
					Mismatched{check, matches})
			} else if settings.MaxHops > 0 && len(matches) > settings.MaxHops {
				// Regardless of whether we ended up
				// in the right place, it took too
				// many hops to get there.
				r.ExceededHops = append(
					r.ExceededHops,
					Mismatched{check, matches})
			} else if check.Code != finalMatch.Code ||
				check.Expected != finalMatch.Match {
				// There is at least one match, but
				// the final URL and code are not the
				// ones we expected.
				r.Mismatched = append(
					r.Mismatched,
					Mismatched{check, matches})
			} else {
				// Recognize that the first
				// rule was checked properly.
				used[matches[0].LineNum] = true
			}
		}
	}

	for _, rule := range rules.rules {
		if used[rule.LineNum] {
			r.Matched = append(r.Matched, rule)
		} else {
			r.Unmatched = append(r.Unmatched, rule)
		}
	}

	return &r, nil
}
