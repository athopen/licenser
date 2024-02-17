package license

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestSatisfies(t *testing.T) {
	cases := []struct {
		expressions []string
		licenses    []string
		result      bool
	}{
		{
			expressions: []string{"APACHE-2.0"},
			licenses:    []string{"APACHE-2.0"},
			result:      true,
		},
		{
			expressions: []string{"Apache-2.0"},
			licenses:    []string{"APACHE-2.0"},
			result:      true,
		},
		{
			expressions: []string{"(LGPL-2.1-only or GPL-3.0-or-later)"},
			licenses:    []string{"LGPL-2.1-only"},
			result:      true,
		},
		{
			expressions: []string{"UNLICENSE"},
			licenses:    []string{"Apache-2.0"},
			result:      false,
		},
		{
			expressions: []string{"Apache-2.0"},
			licenses:    []string{"UNKNOWN-License"},
			result:      false,
		},
	}

	for _, tc := range cases {
		assert.Equal(t, Satisfies(tc.expressions, tc.licenses), tc.result)
	}
}
