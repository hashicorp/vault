package testhelpers

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// ValidateResponse validates whether the given response object conforms
// to the response schema (schema.Fields). It cycles through the data map and
// validates conversions in the schema. In "strict" mode, this function will
// also ensure that the data map has all schema-required fields and does not
// have any fields outside of the schema.
//
// This function is intended to be used in tests only.
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
// This function is intended to be used in tests only.
func ValidateResponseData(schema *framework.Response, data map[string]interface{}, strict bool) error {
	// nothing to validate
	if schema == nil {
		return nil
	}

	fd := framework.FieldData{
		Raw:    data,
		Schema: schema.Fields,
	}

	if strict {
		return fd.ValidateStrict()
	}

	return fd.Validate()
}
