package main

import (
	"flag"
	"fmt"
	"os"
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
	htaccess_file := remaining[0]
	test_file := remaining[1]
	fmt.Println(*ignore_untested, *error_untested,
		htaccess_file, test_file)
}
