package gowhere

import (
	"bufio"
	"io"
	"strings"
)

// ParseRules reads the redirect rules (such as from an htaccess file)
// and returns a RuleSet containing all of them. Stops on the first
// error parsing the file.
func ParseRules(fd io.Reader) (*RuleSet, error) {
	var rules RuleSet
	lineNum := 0
	input := bufio.NewScanner(fd)
	for input.Scan() {
		line := strings.Trim(input.Text(), " \t\n")
		lineNum++

		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}

		r, err := NewRule(lineNum, strings.Fields(line))
		if err != nil {
			return &rules, err
		}
		rules.rules = append(rules.rules, *r)
	}
	return &rules, nil
}

// ParseChecks reads the rule checks and returns a slice of Check
// objects. Stops on the first error parsing the file.
func ParseChecks(fd io.Reader) ([]Check, error) {
	var checks []Check
	lineNum := 0
	input := bufio.NewScanner(fd)
	for input.Scan() {
		line := strings.Trim(input.Text(), " \t\n")
		lineNum++

		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}

		t, err := NewCheck(lineNum, strings.Fields(line))
		if err != nil {
			return checks, err
		}
		checks = append(checks, *t)
	}
	return checks, nil
}
