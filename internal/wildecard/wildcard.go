package wildecard

import "path/filepath"

type Matcher []string

func NewMatcher(patterns []string) Matcher {
	return patterns
}

func (m Matcher) Match(s string) (bool, error) {
	for _, r := range m {
		match, err := filepath.Match(r, s)
		if err != nil {
			return false, err
		}

		if match {
			return true, nil
		}
	}

	return false, nil
}
