package gowhere

import (
	"bytes"
	"testing"
)

func TestRuleSetFirstMatchRedirect(t *testing.T) {
	r, _ := NewRule(1, []string{"redirect", "301",
		"/project/def/new_page.html",
		"/project/def/other_page.html"})
	rs := RuleSet{[]Rule{*r}}

	m := rs.firstMatch("/project/def/new_page.html", true)
	if m == nil {
		t.Error("got nil instead of a match")
	}
	if m.Match != "/project/def/other_page.html" {
		t.Errorf("match is %s instead of /project/def/other_page.html",
			m.Match)
	}

	m = rs.firstMatch("/project/def/same_page.html", true)
	if m != nil {
		t.Errorf("got match for %s instead of nil", m.Match)
	}
}

func TestRuleSetFirstMatchRedirectMatch(t *testing.T) {
	r, _ := NewRule(1, []string{"redirectmatch", "301",
		"^/project/([^/]+)/old_page.html$",
		"/project/$1/new_page.html"})
	rs := RuleSet{[]Rule{*r}}

	m := rs.firstMatch("/project/def/old_page.html", true)
	if m == nil {
		t.Error("got nil instead of a match")
	}
	if m.Match != "/project/def/new_page.html" {
		t.Errorf("match is %s instead of /project/def/new_page.html",
			m.Match)
	}

	m = rs.firstMatch("/project/def/same_page.html", true)
	if m != nil {
		t.Errorf("got match for %s instead of nil", m.Match)
	}
}

func TestRuleSetFindMatchesNone(t *testing.T) {
	rs := RuleSet{[]Rule{}}
	c := Check{
		LineNum:  1,
		Input:    "/project/def/old_page.html",
		Code:     "301",
		Expected: "/project/def/new_page.html",
	}
	s := Settings{Verbose: true, MaxHops: 0}
	matches := rs.FindMatches(&c, s)
	if len(matches) != 0 {
		t.Errorf("found %d matches instead of 0: %v",
			len(matches), matches)
	}
}

func TestRuleSetFindMatchesOne(t *testing.T) {
	r, _ := NewRule(1, []string{"redirectmatch", "301",
		"^/project/([^/]+)/old_page.html$",
		"/project/$1/new_page.html"})
	rs := RuleSet{[]Rule{*r}}
	c := Check{
		LineNum:  1,
		Input:    "/project/def/old_page.html",
		Code:     "301",
		Expected: "/project/def/new_page.html",
	}
	s := Settings{Verbose: true, MaxHops: 0}
	matches := rs.FindMatches(&c, s)
	if len(matches) != 1 {
		t.Errorf("found %d matches instead of 1: %v",
			len(matches), matches)
	}
	m := matches[0]
	if m.Match != c.Expected {
		t.Errorf("found match %s instead of %s",
			m.Match, c.Expected)
	}
}

func TestRuleSetFindMatchesTooManyHops(t *testing.T) {
	data := []byte(`redirectmatch 301 ^/renamed/old/ /renamed/new1/
redirectmatch 301 ^/renamed/new1/ /renamed/new2/
redirectmatch 301 ^/renamed/new2/ /renamed/new3/
`)
	input := bytes.NewReader(data)
	rs, _ := ParseRules(input)
	c := Check{
		LineNum:  1,
		Input:    "/renamed/old/",
		Code:     "301",
		Expected: "/renamed/new3/",
	}
	s := Settings{Verbose: true, MaxHops: 2}
	matches := rs.FindMatches(&c, s)
	// The redirect that puts us over the MaxHops value is included in
	// the return set.
	if len(matches) != 3 {
		t.Errorf("found %d matches instead of 3: %v",
			len(matches), matches)
	}
}

func TestRuleSetFindMatchesCycle(t *testing.T) {
	data := []byte(`redirect 301 /renamed/old/ /renamed/new1/
redirect 301 /renamed/new1/ /renamed/new2/
redirect 301 /renamed/new2/ /renamed/old/
`)
	input := bytes.NewReader(data)
	rs, _ := ParseRules(input)
	c := Check{
		LineNum:  1,
		Input:    "/renamed/old/",
		Code:     "301",
		Expected: "/renamed/new3/",
	}
	s := Settings{Verbose: true, MaxHops: 0}
	matches := rs.FindMatches(&c, s)
	if len(matches) != 3 {
		t.Errorf("found %d matches instead of 3: %v",
			len(matches), matches)
	}
	if matches[0].Pattern != matches[2].Match {
		t.Errorf("start %s and end %s do not match",
			matches[0].Pattern, matches[2].Match)
	}
}
