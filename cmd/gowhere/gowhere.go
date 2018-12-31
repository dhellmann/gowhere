package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dhellmann/gowhere/pkg/gowhere"
)
var ignoreUntested = flag.Bool("ignore-untested", false,
	"ignore untested rules")
var errorUntested = flag.Bool("error-untested", false,
	"error if there are untested rules")
var maxHops = flag.Int("max-hops", 0, "how many hops are allowed")
var verbose = flag.Bool("v", false, "turn on verbose output")

func showCheckAndMatches(msg string, check *gowhere.Check, matches []gowhere.Match) {
	fmt.Printf("%s on line %d: '%s' should produce %s '%s'\n",
		msg, check.LineNum, check.Input, check.Code, check.Expected)
	for _, m := range matches {
		fmt.Printf("    %s -> %s %s [line %d]\n",
			check.Input, m.Code, m.Match, m.LineNum)
	}
}

func main() {
	flag.Parse()
	remaining := flag.Args()
	if len(remaining) < 2 {
		fmt.Fprintf(os.Stderr,
			"please specify htaccess file and test file\n")
		return
	}
	if len(remaining) > 2 {
		fmt.Fprintf(os.Stderr,
			"unrecognized arguments: %v\n", remaining[2:])
		return
	}

	htaccessFile, err := os.Open(remaining[0])
	defer htaccessFile.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read htaccess file %s: %v\n",
			remaining[0], err)
		return
	}
	rules, err := gowhere.ParseRules(htaccessFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse htaccess file %s: %v\n",
			remaining[0], err)
		return
	}

	testFile, err := os.Open(remaining[1])
	defer testFile.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read test file %s: %v\n",
			remaining[1], err)
		return
	}
	checks, err := gowhere.ParseChecks(testFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse test file %s: %v\n",
			remaining[1], err)
		return
	}

	settings := gowhere.Settings{*verbose, *maxHops}
	results := gowhere.ProcessChecks(rules, checks, settings)

	failures := 0

	if *verbose {
		fmt.Println("")
	}

	for _, item := range results.Mismatched {
		failures++
		if len(item.Matches) > 0 {
			showCheckAndMatches("Unexpected rule matched check",
				&(item.Check), item.Matches)
		} else {
			showCheckAndMatches("No rule matched check",
				&(item.Check), item.Matches)
		}
	}

	for _, item := range results.Cycles {
		failures++
		showCheckAndMatches("Cycle found from rule",
			&(item.Check), item.Matches)
	}

	for _, item := range results.ExceededHops {
		failures++
		showCheckAndMatches("Excessive redirects found from rule",
			&(item.Check), item.Matches)
	}

	for _, item := range results.Unmatched {
		if *errorUntested {
			failures++
		}
		fmt.Printf("Untested rule %s\n", item.String())
	}

	if failures > 0 {
		fmt.Fprintf(os.Stderr, "\n%d failures\n", failures)
		os.Exit(1)
	}
}
