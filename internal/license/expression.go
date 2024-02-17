package license

import (
	"slices"
	"strings"

	"github.com/github/go-spdx/expression"
)

func Satisfies(expressions []string, licenses []string) bool {
	for _, exp := range expressions {
		// check it there is an exact match (ignore case)
		if slices.ContainsFunc(licenses, func(s string) bool {
			return strings.EqualFold(s, exp)
		}) {
			return true
		}

		// check if an exp matches spdx standards
		if satisfies, err := expression.Satisfies(strings.ToUpper(exp), licenses); satisfies && err == nil {
			return true
		}
	}

	return false
}
