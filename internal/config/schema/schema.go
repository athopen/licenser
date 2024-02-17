package schema

import (
	_ "embed"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed spec.json
var Schema string

func Validate(dict map[string]interface{}) error {
	schemaLoader := gojsonschema.NewStringLoader(Schema)
	dataLoader := gojsonschema.NewGoLoader(dict)

	result, err := gojsonschema.Validate(schemaLoader, dataLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		return fmt.Errorf("config is not valid")
	}

	return nil
}
