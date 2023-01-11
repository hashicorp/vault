package testhelpers

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestValidateResponse(t *testing.T) {
	cases := map[string]struct {
		schema        *framework.Response
		response      *logical.Response
		strict        bool
		errorExpected bool
	}{
		"nil schema, nil response, strict": {
			schema:        nil,
			response:      nil,
			strict:        true,
			errorExpected: false,
		},

		"nil schema, nil response, not strict": {
			schema:        nil,
			response:      nil,
			strict:        false,
			errorExpected: false,
		},

		"nil schema, good response, strict": {
			schema: nil,
			response: &logical.Response{
				Data: map[string]interface{}{
					"foo": "bar",
				},
			},
			strict:        true,
			errorExpected: false,
		},

		"nil schema, good response, not strict": {
			schema: nil,
			response: &logical.Response{
				Data: map[string]interface{}{
					"foo": "bar",
				},
			},
			strict:        true,
			errorExpected: false,
		},

		"nil schema fields, good response, strict": {
			schema: &framework.Response{},
			response: &logical.Response{
				Data: map[string]interface{}{
					"foo": "bar",
				},
			},
			strict:        true,
			errorExpected: false,
		},

		"nil schema fields, good response, not strict": {
			schema: &framework.Response{},
			response: &logical.Response{
				Data: map[string]interface{}{
					"foo": "bar",
				},
			},
			strict:        true,
			errorExpected: false,
		},

		"string schema field, string response, strict": {
			schema: &framework.Response{
				Fields: map[string]*framework.FieldSchema{
					"foo": {
						Type: framework.TypeString,
					},
				},
			},
			response: &logical.Response{
				Data: map[string]interface{}{
					"foo": "bar",
				},
			},
			strict:        true,
			errorExpected: false,
		},

		"string schema field, string response, not strict": {
			schema: &framework.Response{
				Fields: map[string]*framework.FieldSchema{
					"foo": {
						Type: framework.TypeString,
					},
				},
			},
			response: &logical.Response{
				Data: map[string]interface{}{
					"foo": "bar",
				},
			},
			strict:        false,
			errorExpected: false,
		},

		"string schema not required field, empty response, strict": {
			schema: &framework.Response{
				Fields: map[string]*framework.FieldSchema{
					"foo": {
						Type:     framework.TypeString,
						Required: false,
					},
				},
			},
			response:      &logical.Response{},
			strict:        true,
			errorExpected: false,
		},

		"string schema required field, empty response, strict": {
			schema: &framework.Response{
				Fields: map[string]*framework.FieldSchema{
					"foo": {
						Type:     framework.TypeString,
						Required: true,
					},
				},
			},
			response:      &logical.Response{},
			strict:        true,
			errorExpected: true,
		},

		"string schema required field, empty response, not strict": {
			schema: &framework.Response{
				Fields: map[string]*framework.FieldSchema{
					"foo": {
						Type:     framework.TypeString,
						Required: true,
					},
				},
			},
			response:      &logical.Response{},
			strict:        false,
			errorExpected: false,
		},

		"string schema required field, nil response, strict": {
			schema: &framework.Response{
				Fields: map[string]*framework.FieldSchema{
					"foo": {
						Type:     framework.TypeString,
						Required: true,
					},
				},
			},
			response:      nil,
			strict:        true,
			errorExpected: true,
		},

		"string schema required field, nil response, not strict": {
			schema: &framework.Response{
				Fields: map[string]*framework.FieldSchema{
					"foo": {
						Type:     framework.TypeString,
						Required: true,
					},
				},
			},
			response:      nil,
			strict:        false,
			errorExpected: false,
		},

		"empty schema, string response, strict": {
			schema: &framework.Response{
				Fields: map[string]*framework.FieldSchema{},
			},
			response: &logical.Response{
				Data: map[string]interface{}{
					"foo": "bar",
				},
			},
			strict:        true,
			errorExpected: true,
		},

		"empty schema, string response, not strict": {
			schema: &framework.Response{
				Fields: map[string]*framework.FieldSchema{},
			},
			response: &logical.Response{
				Data: map[string]interface{}{
					"foo": "bar",
				},
			},
			strict:        false,
			errorExpected: false,
		},

		"time schema, string response, strict": {
			schema: &framework.Response{
				Fields: map[string]*framework.FieldSchema{
					"time": {
						Type:     framework.TypeTime,
						Required: true,
					},
				},
			},
			response: &logical.Response{
				Data: map[string]interface{}{
					"time": "2024-12-11T09:08:07Z",
				},
			},
			strict:        true,
			errorExpected: false,
		},

		"time schema, string response, not strict": {
			schema: &framework.Response{
				Fields: map[string]*framework.FieldSchema{
					"time": {
						Type:     framework.TypeTime,
						Required: true,
					},
				},
			},
			response: &logical.Response{
				Data: map[string]interface{}{
					"time": "2024-12-11T09:08:07Z",
				},
			},
			strict:        false,
			errorExpected: false,
		},

		"time schema, time response, strict": {
			schema: &framework.Response{
				Fields: map[string]*framework.FieldSchema{
					"time": {
						Type:     framework.TypeTime,
						Required: true,
					},
				},
			},
			response: &logical.Response{
				Data: map[string]interface{}{
					"time": time.Date(2024, 12, 11, 9, 8, 7, 0, time.UTC),
				},
			},
			strict:        true,
			errorExpected: false,
		},

		"time schema, time response, not strict": {
			schema: &framework.Response{
				Fields: map[string]*framework.FieldSchema{
					"time": {
						Type:     framework.TypeTime,
						Required: true,
					},
				},
			},
			response: &logical.Response{
				Data: map[string]interface{}{
					"time": time.Date(2024, 12, 11, 9, 8, 7, 0, time.UTC),
				},
			},
			strict:        false,
			errorExpected: false,
		},
	}

	for name, tc := range cases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := ValidateResponse(
				tc.schema,
				tc.response,
				tc.strict,
			)
			if err == nil && tc.errorExpected == true {
				t.Fatalf("expected an error, got nil")
			}
			if err != nil && tc.errorExpected == false {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
