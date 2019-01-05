package gowhere

import (
	"fmt"
)

// RuleSet holds a group of Rules to be applied together
type RuleSet struct {
	rules []Rule
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
