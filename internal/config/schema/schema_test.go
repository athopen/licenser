package schema

import (
	"testing"

	"gotest.tools/v3/assert"
)

type dict map[string]interface{}

func TestValidate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		cfg := dict{
			"licenses": []string{
				"foo",
			},
			"excluded": []string{
				"bar",
			},
		}

		assert.NilError(t, Validate(cfg))
	})

	t.Run("undefined top level option", func(t *testing.T) {
		cfg := dict{
			"foo": dict{},
		}

		assert.ErrorContains(t, Validate(cfg), "config is not valid")
	})
}
