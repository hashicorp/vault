package testhelpers

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// ValidateResponse validates whether the given response object conforms
// to the response schema (schema.Fields). It cycles through the data map and
// validates conversions in the schema. In "strict" mode, this function will
// also ensure that the data map has all schema-required fields and does not
// have any fields outside of the schema.
//
// This function is inefficient and is intended to be used in tests only.
func ValidateResponse(schema *framework.Response, response *logical.Response, strict bool) error {
	if response != nil {
		return ValidateResponseData(schema, response.Data, strict)
	}

	return ValidateResponseData(schema, nil, strict)
}

// ValidateResponseData validates whether the given response data map conforms
// to the response schema (schema.Fields). It cycles through the data map and
// validates conversions in the schema. In "strict" mode, this function will
// also ensure that the data map has all schema-required fields and does not
// have any fields outside of the schema.
//
// This function is inefficient and is intended to be used in tests only.
func ValidateResponseData(schema *framework.Response, data map[string]interface{}, strict bool) error {
	// nothing to validate
	if schema == nil {
		return nil
	}

	// Convert to json & back to coerse data entries into the final
	// output format expected by Validate() and ValidateStrict(). This is
	// not efficient and is done for testing purposes only.
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to convert input to json: %w", err)
	}

	var jsonData map[string]interface{}

	if err := json.Unmarshal(jsonBytes, &jsonData); err != nil {
		return fmt.Errorf("failed to unmashal data: %w", err)
	}

	// Validate
	fd := framework.FieldData{
		Raw:    jsonData,
		Schema: schema.Fields,
	}

	if strict {
		return fd.ValidateStrict()
	}

	return fd.Validate()
}
