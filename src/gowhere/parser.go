package gowhere

import (
	"bufio"
	"fmt"
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
			return &rules, fmt.Errorf(
				"Could not parse rule on line %d (%s): %s",
				line_num, input.Text(), err)
		}
		rules.rules = append(rules.rules, *r)
	}
	return &rules, nil
}
