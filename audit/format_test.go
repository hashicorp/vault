// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package audit

import (
	"context"
	"io"
	"testing"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/copystructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testingFormatWriter struct {
	salt         *salt.Salt
	lastRequest  *AuditRequestEntry
	lastResponse *AuditResponseEntry
}

func (fw *testingFormatWriter) WriteRequest(_ io.Writer, entry *AuditRequestEntry) error {
	fw.lastRequest = entry
	return nil
}

func (fw *testingFormatWriter) WriteResponse(_ io.Writer, entry *AuditResponseEntry) error {
	fw.lastResponse = entry
	return nil
}

func (fw *testingFormatWriter) Salt(ctx context.Context) (*salt.Salt, error) {
	if fw.salt != nil {
		return fw.salt, nil
	}
	var err error
	fw.salt, err = salt.NewSalt(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	return fw.salt, nil
}

// hashExpectedValueForComparison replicates enough of the audit HMAC process on a piece of expected data in a test,
// so that we can use assert.Equal to compare the expected and output values.
func (fw *testingFormatWriter) hashExpectedValueForComparison(input map[string]interface{}) map[string]interface{} {
	// Copy input before modifying, since we may re-use the same data in another test
	copied, err := copystructure.Copy(input)
	if err != nil {
		panic(err)
	}
	copiedAsMap := copied.(map[string]interface{})

	salter, err := fw.Salt(context.Background())
	if err != nil {
		panic(err)
	}

	err = hashMap(salter.GetIdentifiedHMAC, copiedAsMap, nil)
	if err != nil {
		panic(err)
	}

	return copiedAsMap
}

func TestFormatRequestErrors(t *testing.T) {
	config := FormatterConfig{}
	formatter := AuditFormatter{
		AuditFormatWriter: &testingFormatWriter{},
	}

	if err := formatter.FormatRequest(context.Background(), io.Discard, config, &logical.LogInput{}); err == nil {
		t.Fatal("expected error due to nil request")
	}

	in := &logical.LogInput{
		Request: &logical.Request{},
	}
	if err := formatter.FormatRequest(context.Background(), nil, config, in); err == nil {
		t.Fatal("expected error due to nil writer")
	}
}

func TestFormatResponseErrors(t *testing.T) {
	config := FormatterConfig{}
	formatter := AuditFormatter{
		AuditFormatWriter: &testingFormatWriter{},
	}

	if err := formatter.FormatResponse(context.Background(), io.Discard, config, &logical.LogInput{}); err == nil {
		t.Fatal("expected error due to nil request")
	}

	in := &logical.LogInput{
		Request: &logical.Request{},
	}
	if err := formatter.FormatResponse(context.Background(), nil, config, in); err == nil {
		t.Fatal("expected error due to nil writer")
	}
}

func TestElideListResponses(t *testing.T) {
	tfw := testingFormatWriter{}
	formatter := AuditFormatter{&tfw}
	ctx := namespace.RootContext(context.Background())

	type test struct {
		name         string
		inputData    map[string]interface{}
		expectedData map[string]interface{}
	}

	tests := []test{
		{
			"nil data",
			nil,
			nil,
		},
		{
			"Normal list (keys only)",
			map[string]interface{}{
				"keys": []string{"foo", "bar", "baz"},
			},
			map[string]interface{}{
				"keys": 3,
			},
		},
		{
			"Enhanced list (has key_info)",
			map[string]interface{}{
				"keys": []string{"foo", "bar", "baz", "quux"},
				"key_info": map[string]interface{}{
					"foo":  "alpha",
					"bar":  "beta",
					"baz":  "gamma",
					"quux": "delta",
				},
			},
			map[string]interface{}{
				"keys":     4,
				"key_info": 4,
			},
		},
		{
			"Unconventional other values in a list response are not touched",
			map[string]interface{}{
				"keys":           []string{"foo", "bar"},
				"something_else": "baz",
			},
			map[string]interface{}{
				"keys":           2,
				"something_else": "baz",
			},
		},
		{
			"Conventional values in a list response are not elided if their data types are unconventional",
			map[string]interface{}{
				"keys": map[string]interface{}{
					"You wouldn't expect keys to be a map": nil,
				},
				"key_info": []string{
					"You wouldn't expect key_info to be a slice",
				},
			},
			map[string]interface{}{
				"keys": map[string]interface{}{
					"You wouldn't expect keys to be a map": nil,
				},
				"key_info": []string{
					"You wouldn't expect key_info to be a slice",
				},
			},
		},
	}
	oneInterestingTestCase := tests[2]

	formatResponse := func(
		t *testing.T,
		config FormatterConfig,
		operation logical.Operation,
		inputData map[string]interface{},
	) {
		err := formatter.FormatResponse(ctx, io.Discard, config, &logical.LogInput{
			Request:  &logical.Request{Operation: operation},
			Response: &logical.Response{Data: inputData},
		})
		require.Nil(t, err)
	}

	t.Run("Default case", func(t *testing.T) {
		config := FormatterConfig{ElideListResponses: true}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				formatResponse(t, config, logical.ListOperation, tc.inputData)
				assert.Equal(t, tfw.hashExpectedValueForComparison(tc.expectedData), tfw.lastResponse.Response.Data)
			})
		}
	})

	t.Run("When Operation is not list, eliding does not happen", func(t *testing.T) {
		config := FormatterConfig{ElideListResponses: true}
		tc := oneInterestingTestCase
		formatResponse(t, config, logical.ReadOperation, tc.inputData)
		assert.Equal(t, tfw.hashExpectedValueForComparison(tc.inputData), tfw.lastResponse.Response.Data)
	})

	t.Run("When ElideListResponses is false, eliding does not happen", func(t *testing.T) {
		config := FormatterConfig{ElideListResponses: false}
		tc := oneInterestingTestCase
		formatResponse(t, config, logical.ListOperation, tc.inputData)
		assert.Equal(t, tfw.hashExpectedValueForComparison(tc.inputData), tfw.lastResponse.Response.Data)
	})

	t.Run("When Raw is true, eliding still happens", func(t *testing.T) {
		config := FormatterConfig{ElideListResponses: true, Raw: true}
		tc := oneInterestingTestCase
		formatResponse(t, config, logical.ListOperation, tc.inputData)
		assert.Equal(t, tc.expectedData, tfw.lastResponse.Response.Data)
	})
}
