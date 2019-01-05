package gowhere

import (
	"fmt"
)

// Check represents a test for one Rule
type Check struct {
	// The line of the input file where the check was found
	LineNum int
	// The input to give to the RuleSet
	Input string
	// The expected HTTP response code
	Code string
	// The expected destination of the redirection
	Expected string
}

// NewCheck creates a Check from the strings on the input line
func NewCheck(lineNum int, params []string) (*Check, error) {
	var t Check

	t.LineNum = lineNum

	if len(params) == 3 {
		// input code expected
		t.Input = params[0]
		t.Code = params[1]
		t.Expected = params[2]
		return &t, nil
	}

	if len(params) == 2 {
		// input code
		// (no expected redirect)
		t.Input = params[0]
		t.Code = params[1]
		return &t, nil
	}

	return nil, fmt.Errorf("Could not understand check on line %d: %v",
		lineNum, params)
}
