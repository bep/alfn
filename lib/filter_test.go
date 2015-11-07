package lib

import (
	rss "github.com/jteeuwen/go-pkg-rss"
	"regexp"
	"testing"
)

func TestAppyFilters(t *testing.T) {

	var unFiltered = []*rss.Item{
		&rss.Item{
			Title:       "AB",
			Description: "CD"},
		&rss.Item{
			Title:       "EF",
			Description: "GH"},
		&rss.Item{
			Title:       "IJ",
			Description: "KL"},
		&rss.Item{
			Content: &rss.Content{Text: "AB CD"}},
		&rss.Item{
			Title:       "MUSTCONTAIN CD",
			Description: "EF"},
		&rss.Item{
			Content: &rss.Content{Text: "MUSTCONTAIN CD EF"}},
		&rss.Item{
			Title:       "MUSTNOTCONTAIN NO",
			Description: "WAY"},
		&rss.Item{
			Content: &rss.Content{Text: "MUSTNOTCONTAIN NO WAY"}},
	}

	// only regular matchers
	var matchers1 = []matcher{
		regexpMatcher{Regexp: regexp.MustCompile("AB CD")},
		regexpMatcher{Regexp: regexp.MustCompile("IJ KL")},
	}

	filtered1 := appyFilters(unFiltered, matchers1)

	if len(filtered1) != 3 {
		t.Error("Got", len(filtered1))
	}

	// a mix of matchers
	var matchers2 = []matcher{
		regexpMatcher{Regexp: regexp.MustCompile(".*")},
		regexpMatcher{negate: true, matchBreaker: true, Regexp: regexp.MustCompile("MUSTNOTCONTAIN")},
		regexpMatcher{negate: false, matchBreaker: true, Regexp: regexp.MustCompile("MUSTCONTAIN")},
	}

	filtered2 := appyFilters(unFiltered, matchers2)

	if len(filtered2) != 2 {
		t.Error("Got", len(filtered2))
	}
}
