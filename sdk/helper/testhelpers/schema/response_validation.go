package schema

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// ValidateResponseData is a test helper that validates whether the given
// response data map conforms to the response schema (schema.Fields). It cycles
// through the data map and validates conversions in the schema. In "strict"
// mode, this function will also ensure that the data map has all schema's
// requred fields and does not have any fields outside of the schema.
func ValidateResponse(t *testing.T, schema *framework.Response, response *logical.Response, strict bool) {
	t.Helper()

	if response != nil {
		ValidateResponseData(t, schema, response.Data, strict)
	} else {
		ValidateResponseData(t, schema, nil, strict)
	}
}

// ValidateResponse is a test helper that validates whether the given response
// object conforms to the response schema (schema.Fields). It cycles through
// the data map and validates conversions in the schema. In "strict" mode, this
// function will also ensure that the data map has all schema-required fields
// and does not have any fields outside of the schema.
func ValidateResponseData(t *testing.T, schema *framework.Response, data map[string]interface{}, strict bool) {
	t.Helper()

	if err := validateResponseDataImpl(
		schema,
		data,
		strict,
	); err != nil {
		t.Fatalf("validation error: %v; response data: %#v", err, data)
	}
}

// validateResponseDataImpl is extracted so that it can be tested
func validateResponseDataImpl(schema *framework.Response, data map[string]interface{}, strict bool) error {
	// nothing to validate
	if schema == nil {
		return nil
	}

	// Marshal the data to JSON and back to convert the map's values into
	// JSON strings expected by Validate() and ValidateStrict(). This is
	// not efficient and is done for testing purposes only.
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to convert input to json: %w", err)
	}

	var dataWithStringValues map[string]interface{}
	if err := json.Unmarshal(
		jsonBytes,
		&dataWithStringValues,
	); err != nil {
		return fmt.Errorf("failed to unmashal data: %w", err)
	}

	// Validate
	fd := framework.FieldData{
		Raw:    dataWithStringValues,
		Schema: schema.Fields,
	}

	if strict {
		return fd.ValidateStrict()
	}

	return fd.Validate()
}

// FindResponseSchema is a test helper to extract the response schema from a given framework path / operation
func FindResponseSchema(t *testing.T, paths []*framework.Path, pathIdx int, operation logical.Operation) *framework.Response {
	t.Helper()

	if pathIdx >= len(paths) {
		t.Fatalf("path index %d is out of range", pathIdx)
	}

	schemaPath := paths[pathIdx]

	return GetResponseSchema(t, schemaPath, operation)
}

func GetResponseSchema(t *testing.T, path *framework.Path, operation logical.Operation) *framework.Response {
	t.Helper()

	schemaOperation, ok := path.Operations[operation]
	if !ok {
		t.Fatalf(
			"could not find response schema: %s: %q operation does not exist",
			path.Pattern,
			operation,
		)
	}

	var schemaResponses []framework.Response

	for _, status := range []int{
		http.StatusOK,        // 200
		http.StatusAccepted,  // 202
		http.StatusNoContent, // 204
	} {
		schemaResponses, ok = schemaOperation.Properties().Responses[status]
		if ok {
			break
		}
	}

	if len(schemaResponses) == 0 {
		t.Fatalf(
			"could not find response schema: %s: %q operation: no responses found",
			path.Pattern,
			operation,
		)
	}

	return &schemaResponses[0]
}

func ResponseValidatingCallback(t *testing.T) func(logical.Backend, *logical.Request, *logical.Response) {

	type PathRouter interface {
		Route(string) *framework.Path
	}

	return func(b logical.Backend, req *logical.Request, resp *logical.Response) {
		t.Helper()

		if b == nil {
			t.Fatalf("non-nil backend required")
		}
		backend, ok := b.(PathRouter)
		if !ok {
			t.Fatalf("could not cast %T to have `Route(string) *framework.Path`", b)
		}

		route := backend.Route(req.Path)
		if route == nil {
			t.Fatalf("backend %T could not find a route for %s", b, req.Path)
		}

		ValidateResponse(
			t,
			GetResponseSchema(t, route, req.Operation),
			resp,
			true,
		)
	}
}
