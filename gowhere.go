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
	fmt.Printf("%v\n", *results)
}
