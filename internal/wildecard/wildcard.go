package wildecard

import "path/filepath"

func Match(s string, patterns []string) bool {
	for _, p := range patterns {
		match, err := filepath.Match(p, s)
		if err != nil {
			return false
		}

		if match {
			return true
		}
	}

	return false
}
