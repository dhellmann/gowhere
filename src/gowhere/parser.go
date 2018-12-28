package gowhere

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func ParseRules(fd io.Reader) (*RuleSet, error) {
	var rules RuleSet
	line_num := 1
	input := bufio.NewScanner(fd)
	for input.Scan() {
		params := strings.Split(input.Text(), " ")
		err := rules.Add(line_num, params)
		if err != nil {
			return nil, fmt.Errorf("Could not parse rule on line %d (%s): %s", 
				line_num, input.Text(), err)
		}
		line_num++
	}
	return &rules, nil
}
