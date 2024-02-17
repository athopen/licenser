package wildecard

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestMatch(t *testing.T) {
	cases := []struct {
		s        string
		patterns []string
		matches  bool
	}{
		{
			s:        "foo",
			patterns: []string{"foo"},
			matches:  true,
		},
		{
			s:        "foo",
			patterns: []string{"*"},
			matches:  true,
		},
		{
			s:        "foo/bar",
			patterns: []string{"foo/*"},
			matches:  true,
		},
		{
			s:        "foo/bar",
			patterns: []string{"foo/[a-z]ar"},
			matches:  true,
		},
		{
			s:        "foo/bar",
			patterns: []string{"foo/?ar"},
			matches:  true,
		},
		{
			s:        "foo/bar",
			patterns: []string{"foo/?az"},
			matches:  false,
		},
	}

	for _, tc := range cases {
		assert.Equal(t, Match(tc.s, tc.patterns), tc.matches)
	}
}
