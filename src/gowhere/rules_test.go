package gowhere

import (
	"bytes"
	"fmt"
	"testing"
)

type Want struct {
	directive string
	code      string
	pattern   string
	target    string
	re        bool
}

func TestNewRule(t *testing.T) {
	var tests = []struct {
		input []string
		want  Want
	}{

		{
			[]string{"redirect", "301",
				"/project/def/new_page.html",
				"/project/def/other_page.html"},
			Want{
				directive: "redirect",
				code:      "301",
				pattern:   "/project/def/new_page.html",
				target:    "/project/def/other_page.html",
				re:        false,
			},
		},

		{
			[]string{"redirect", "410",
				"/project/def/new_page.html"},
			Want{
				directive: "redirect",
				code:      "410",
				pattern:   "/project/def/new_page.html",
				target:    "",
				re:        false,
			},
		},

		{
			[]string{"redirect",
				"/project/def/new_page.html",
				"/project/def/other_page.html"},
			Want{
				directive: "redirect",
				code:      "301",
				pattern:   "/project/def/new_page.html",
				target:    "/project/def/other_page.html",
				re:        false,
			},
		},

		{
			[]string{"redirectmatch",
				"301",
				"^/project/([^/]+)/old_page.html$",
				"/project/$1/new_page.html"},
			Want{
				directive: "redirectmatch",
				code:      "301",
				pattern:   "^/project/([^/]+)/old_page.html$",
				target:    "/project/$1/new_page.html",
				re:        true,
			},
		},
	}

	for n, test := range tests {
		fmt.Printf("Test %d: %v\n", n, test)
		r, err := NewRule(1, test.input)
		fmt.Printf("Rule %d: %v\n", n, r)
		if err != nil {
			t.Errorf("test %d: should not have an error: %v", n, err)
		}
		if test.want.re {
			if r.re == nil {
				t.Errorf("test %d: should have a regexp", n)
			}
		} else {
			if r.re != nil {
				t.Errorf("test %d: should not have a regexp", n)
			}
		}
		if r.LineNum != 1 {
			t.Errorf("test %d: r.LineNum == %d, expected 1", n, r.LineNum)
		}
		if r.Code != test.want.code {
			t.Errorf("test %d: r.Code == %s, expected %s",
				n, r.Code, test.want.code)
		}
		if r.Target != test.want.target {
			t.Errorf("test %d: r.Target == %s, expected %s",
				n, r.Target, test.want.target)
		}
	}
}

func TestNewRuleTooFewParams(t *testing.T) {
	params := []string{
		"redirect",
		"410",
	}
	r, err := NewRule(1, params)
	if err == nil {
		t.Errorf("should have an error: %v", r)
	}
}

func TestRuleMatchString(t *testing.T) {
	r, _ := NewRule(1, []string{"redirect", "301",
		"/project/def/new_page.html",
		"/project/def/other_page.html"})
	s := r.Match("/project/def/new_page.html")
	if s != "/project/def/other_page.html" {
		t.Errorf("received %s instead of /project/def/other_page.html", s)
	}
	s = r.Match("/project/def/no_match.html")
	if s != "" {
		t.Errorf("received %s instead of empty string", s)
	}
}

func TestRuleMatchRegexp(t *testing.T) {
	r, _ := NewRule(1, []string{"redirectmatch", "301",
		"^/project/([^/]+)/old_page.html$",
		"/project/$1/new_page.html"})
	s := r.Match("/project/def/old_page.html")
	if s != "/project/def/new_page.html" {
		t.Errorf("received %s instead of /project/def/new_page.html", s)
	}
	s = r.Match("/project/def/no_match.html")
	if s != "" {
		t.Errorf("received %s instead of empty string", s)
	}
}

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
