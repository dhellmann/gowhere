package gowhere

import (
	"bufio"
	"io"
	"strings"
)

func ParseRules(fd io.Reader) (*RuleSet, error) {
	var rules RuleSet
	line_num := 0
	input := bufio.NewScanner(fd)
	for input.Scan() {
		line := strings.Trim(input.Text(), " \t\n")
		line_num++

		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}

		r, err := NewRule(line_num, strings.Fields(line))
		if err != nil {
			return &rules, err
		}
		rules.rules = append(rules.rules, *r)
	}
	return &rules, nil
}

func ParseChecks(fd io.Reader) ([]Check, error) {
	var checks []Check
	line_num := 0
	input := bufio.NewScanner(fd)
	for input.Scan() {
		line := strings.Trim(input.Text(), " \t\n")
		line_num++

		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}

		t, err := NewCheck(line_num, strings.Fields(line))
		if err != nil {
			return checks, err
		}
		checks = append(checks, *t)
	}
	return checks, nil
}
