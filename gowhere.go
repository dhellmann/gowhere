package main

import (
	"flag"
	"fmt"
	"os"

	"gowhere"
)

var ignore_untested = flag.Bool("ignore-untested", false,
	"ignore untested rules")
var error_untested = flag.Bool("error-untested", false,
	"error if there are untested rules")
var max_hops = flag.Int("max-hops", 0, "how many hops are allowed")
var verbose = flag.Bool("v", false, "turn on verbose output")

func showTestAndMatches(msg string, test *gowhere.RuleTest, matches []gowhere.Match) {
	fmt.Printf("%s on line %d: '%s' should produce %s '%s'\n",
		msg, test.LineNum, test.Input, test.Code, test.Expected)
	for _, m := range matches {
		fmt.Printf("    %s -> %s %s [line %d]\n",
			test.Input, m.Code, m.Match, m.LineNum)
	}
}

func main() {
	flag.Parse()
	remaining := flag.Args()
	if len(remaining) < 2 {
		fmt.Fprintf(os.Stderr,
			"please specify htaccess_file and test_file\n")
		return
	}
	if len(remaining) > 2 {
		fmt.Fprintf(os.Stderr,
			"unrecognized arguments: %v\n", remaining[2:])
		return
	}

	htaccess_file, err := os.Open(remaining[0])
	defer htaccess_file.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read htaccess file %s: %v\n",
			remaining[0], err)
		return
	}
	rules, err := gowhere.ParseRules(htaccess_file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse htaccess file %s: %v\n",
			remaining[0], err)
		return
	}

	test_file, err := os.Open(remaining[1])
	defer test_file.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read test file %s: %v\n",
			remaining[1], err)
		return
	}
	tests, err := gowhere.ParseTests(test_file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse test file %s: %v\n",
			remaining[1], err)
		return
	}

	fmt.Printf("%v\n", *rules)
	fmt.Printf("%v\n", *tests)

	results, err := gowhere.ProcessTests(rules, tests, *max_hops)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Processing failure: %v\n", err)
		return
	}

	failures := 0

	fmt.Println("")
	for _, item := range results.Mismatched {
		failures++
		if len(item.Matches) > 0 {
			showTestAndMatches("Unexpected rule matched test",
				&(item.Test), item.Matches)
		} else {
			showTestAndMatches("No rule matched test",
				&(item.Test), item.Matches)
		}
	}

	for _, item := range results.Cycles {
		failures++
		showTestAndMatches("Cycle found from rule",
			&(item.Test), item.Matches)
	}

	for _, item := range results.ExceededHops {
		failures++
		showTestAndMatches("Excessive redirects found from rule",
			&(item.Test), item.Matches)
	}

	for _, item := range results.Unmatched {
		if *error_untested {
			failures++
		}
		fmt.Printf("Untested rule %s\n", item.String())
	}

	if failures > 0 {
		fmt.Fprintf(os.Stderr, "\n%d failures\n", failures)
	}
}
