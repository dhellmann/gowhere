package gowhere

import (
	"bytes"
	"testing"
)

func TestParseRules(t *testing.T) {
	data := []byte("redirect 301 /project/def/new_page.html /project/def/other_page.html")
	input := bytes.NewReader(data)
	rs, err := ParseRules(input)
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if len(rs.rules) != 1 {
		t.Errorf("got %d rules expected 1", len(rs.rules))
	}
	r := rs.rules[0]
	if r.Directive != "redirect" {
		t.Errorf("got directive %s expected redirect", r.Directive)
	}
	if r.Pattern != "/project/def/new_page.html" {
		t.Errorf("got pattern %s expected /project/def/new_page.html", r.Pattern)
	}
	if r.Target != "/project/def/other_page.html" {
		t.Errorf("got target %s expected /project/def/other_page.html", r.Target)
	}
}

func TestParseRulesIgnoreComments(t *testing.T) {
	data := []byte("#redirect 301 /project/def/new_page.html /project/def/other_page.html")
	input := bytes.NewReader(data)
	rs, err := ParseRules(input)
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if len(rs.rules) != 0 {
		t.Errorf("got %d rules expected 0", len(rs.rules))
	}
}

func TestParseRulesIgnoreBlankLines(t *testing.T) {
	data := []byte("\nredirect 301 /pattern /target\n")
	input := bytes.NewReader(data)
	rs, err := ParseRules(input)
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if len(rs.rules) != 1 {
		t.Errorf("got %d rules expected 0", len(rs.rules))
	}
}
