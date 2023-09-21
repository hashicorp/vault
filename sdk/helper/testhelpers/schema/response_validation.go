// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// ValidateResponse is a test helper that validates whether the given response
// object conforms to the response schema (schema.Fields). It cycles through
// the data map and validates conversions in the schema. In "strict" mode, this
// function will also ensure that the data map has all schema-required fields
// and does not have any fields outside of the schema.
func ValidateResponse(t *testing.T, schema *framework.Response, response *logical.Response, strict bool) {
	t.Helper()

	if response != nil {
		ValidateResponseData(t, schema, response.Data, strict)
	} else {
		ValidateResponseData(t, schema, nil, strict)
	}
}

// ValidateResponseData is a test helper that validates whether the given
// response data map conforms to the response schema (schema.Fields). It cycles
// through the data map and validates conversions in the schema. In "strict"
// mode, this function will also ensure that the data map has all schema's
// requred fields and does not have any fields outside of the schema.
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

	// Certain responses may come through with non-2xx status codes. While
	// these are not always errors (e.g. 3xx redirection codes), we don't
	// consider them for the purposes of schema validation
	if status, exists := data[logical.HTTPStatusCode]; exists {
		s, ok := status.(int)
		if ok && (s < 200 || s > 299) {
			return nil
		}
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

	// these are special fields that will not show up in the final response and
	// should be ignored
	for _, field := range []string{
		logical.HTTPContentType,
		logical.HTTPRawBody,
		logical.HTTPStatusCode,
		logical.HTTPRawBodyAlreadyJSONDecoded,
		logical.HTTPCacheControlHeader,
		logical.HTTPPragmaHeader,
		logical.HTTPWWWAuthenticateHeader,
	} {
		delete(dataWithStringValues, field)

		if _, ok := schema.Fields[field]; ok {
			return fmt.Errorf("encountered a reserved field in response schema: %s", field)
		}
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

// FindResponseSchema is a test helper to extract response schema from the
// given framework path / operation.
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
		// ListOperations have a default response schema that is implicit unless overridden
		if operation == logical.ListOperation {
			return &framework.Response{
				Description: "OK",
				Fields: map[string]*framework.FieldSchema{
					"keys": {
						Type: framework.TypeStringSlice,
					},
				},
			}
		}

		t.Fatalf(
			"could not find response schema: %s: %q operation: no responses found",
			path.Pattern,
			operation,
		)
	}

	return &schemaResponses[0]
}

// ResponseValidatingCallback can be used in setting up a [vault.TestCluster]
// that validates every response against the openapi specifications.
//
// [vault.TestCluster]: https://pkg.go.dev/github.com/hashicorp/vault/vault#TestCluster
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

		// The full request path includes the backend but when passing to the
		// backend, we have to trim the mount point:
		//   `sys/mounts/secret` -> `mounts/secret`
		//   `auth/token/create` -> `create`
		requestPath := strings.TrimPrefix(req.Path, req.MountPoint)

		route := backend.Route(requestPath)
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
