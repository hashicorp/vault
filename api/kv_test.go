// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"reflect"
	"testing"
	"time"
)

func TestExtractVersionMetadata(t *testing.T) {
	t.Parallel()

	inputCreatedTimeStr := "2022-05-06T23:02:04.865025Z"
	inputDeletionTimeStr := "2022-06-17T01:15:03.279013Z"
	expectedCreatedTimeParsed, err := time.Parse(time.RFC3339, inputCreatedTimeStr)
	if err != nil {
		t.Fatalf("unable to parse expected created time: %v", err)
	}
	expectedDeletionTimeParsed, err := time.Parse(time.RFC3339, inputDeletionTimeStr)
	if err != nil {
		t.Fatalf("unable to parse expected created time: %v", err)
	}

	testCases := []struct {
		name     string
		input    *Secret
		expected *KVVersionMetadata
	}{
		{
			name: "a secret",
			input: &Secret{
				Data: map[string]interface{}{
					"data": map[string]interface{}{
						"password": "Hashi123",
					},
					"metadata": map[string]interface{}{
						"version":         10,
						"created_time":    inputCreatedTimeStr,
						"deletion_time":   "",
						"destroyed":       false,
						"custom_metadata": nil,
					},
				},
			},
			expected: &KVVersionMetadata{
				Version:      10,
				CreatedTime:  expectedCreatedTimeParsed,
				DeletionTime: time.Time{},
				Destroyed:    false,
			},
		},
		{
			name: "a secret that has been deleted",
			input: &Secret{
				Data: map[string]interface{}{
					"data": map[string]interface{}{
						"password": "Hashi123",
					},
					"metadata": map[string]interface{}{
						"version":         10,
						"created_time":    inputCreatedTimeStr,
						"deletion_time":   inputDeletionTimeStr,
						"destroyed":       false,
						"custom_metadata": nil,
					},
				},
			},
			expected: &KVVersionMetadata{
				Version:      10,
				CreatedTime:  expectedCreatedTimeParsed,
				DeletionTime: expectedDeletionTimeParsed,
				Destroyed:    false,
			},
		},
		{
			name: "a response from a Write operation",
			input: &Secret{
				Data: map[string]interface{}{
					"version":         10,
					"created_time":    inputCreatedTimeStr,
					"deletion_time":   "",
					"destroyed":       false,
					"custom_metadata": nil,
				},
			},
			expected: &KVVersionMetadata{
				Version:      10,
				CreatedTime:  expectedCreatedTimeParsed,
				DeletionTime: time.Time{},
				Destroyed:    false,
			},
		},
	}

	for _, tc := range testCases {
		versionMetadata, err := extractVersionMetadata(tc.input)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		if !reflect.DeepEqual(versionMetadata, tc.expected) {
			t.Fatalf("%s: got\n%#v\nexpected\n%#v\n", tc.name, versionMetadata, tc.expected)
		}
	}
}

func TestExtractDataAndVersionMetadata(t *testing.T) {
	t.Parallel()

	inputCreatedTimeStr := "2022-05-06T23:02:04.865025Z"
	inputDeletionTimeStr := "2022-06-17T01:15:03.279013Z"
	expectedCreatedTimeParsed, err := time.Parse(time.RFC3339, inputCreatedTimeStr)
	if err != nil {
		t.Fatalf("unable to parse expected created time: %v", err)
	}
	expectedDeletionTimeParsed, err := time.Parse(time.RFC3339, inputDeletionTimeStr)
	if err != nil {
		t.Fatalf("unable to parse expected created time: %v", err)
	}

	readResp := &Secret{
		Data: map[string]interface{}{
			"data": map[string]interface{}{
				"password": "Hashi123",
			},
			"metadata": map[string]interface{}{
				"version":         10,
				"created_time":    inputCreatedTimeStr,
				"deletion_time":   "",
				"destroyed":       false,
				"custom_metadata": nil,
			},
		},
	}

	readRespDeleted := &Secret{
		Data: map[string]interface{}{
			"data": nil,
			"metadata": map[string]interface{}{
				"version":         10,
				"created_time":    inputCreatedTimeStr,
				"deletion_time":   inputDeletionTimeStr,
				"destroyed":       false,
				"custom_metadata": nil,
			},
		},
	}

	testCases := []struct {
		name     string
		input    *Secret
		expected *KVSecret
	}{
		{
			name:  "a response from a Read operation",
			input: readResp,
			expected: &KVSecret{
				Data: map[string]interface{}{
					"password": "Hashi123",
				},
				VersionMetadata: &KVVersionMetadata{
					Version:      10,
					CreatedTime:  expectedCreatedTimeParsed,
					DeletionTime: time.Time{},
					Destroyed:    false,
				},
				// it's tempting to test some Secrets with custom_metadata but
				// we can't in this test because it isn't until we call the
				// extractCustomMetadata function that the custom metadata
				// gets added onto the struct. See TestExtractCustomMetadata.
				CustomMetadata: nil,
				Raw:            readResp,
			},
		},
		{
			name:  "a secret that has been deleted and thus has nil data",
			input: readRespDeleted,
			expected: &KVSecret{
				Data: nil,
				VersionMetadata: &KVVersionMetadata{
					Version:      10,
					CreatedTime:  expectedCreatedTimeParsed,
					DeletionTime: expectedDeletionTimeParsed,
					Destroyed:    false,
				},
				CustomMetadata: nil,
				Raw:            readRespDeleted,
			},
		},
	}

	for _, tc := range testCases {
		dvm, err := extractDataAndVersionMetadata(tc.input)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		if !reflect.DeepEqual(dvm, tc.expected) {
			t.Fatalf("%s: got\n%#v\nexpected\n%#v\n", tc.name, dvm, tc.expected)
		}
	}
}

func TestExtractFullMetadata(t *testing.T) {
	inputCreatedTimeStr := "2022-05-20T00:51:49.419794Z"
	expectedCreatedTimeParsed, err := time.Parse(time.RFC3339, inputCreatedTimeStr)
	if err != nil {
		t.Fatalf("unable to parse expected created time: %v", err)
	}

	inputUpdatedTimeStr := "2022-05-20T20:23:43.284488Z"
	expectedUpdatedTimeParsed, err := time.Parse(time.RFC3339, inputUpdatedTimeStr)
	if err != nil {
		t.Fatalf("unable to parse expected updated time: %v", err)
	}

	inputDeletedTimeStr := "2022-05-21T00:05:49.521697Z"
	expectedDeletedTimeParsed, err := time.Parse(time.RFC3339, inputDeletedTimeStr)
	if err != nil {
		t.Fatalf("unable to parse expected deletion time: %v", err)
	}

	metadataResp := &Secret{
		Data: map[string]interface{}{
			"cas_required":    true,
			"created_time":    inputCreatedTimeStr,
			"current_version": 2,
			"custom_metadata": map[string]interface{}{
				"org": "eng",
			},
			"delete_version_after": "200s",
			"max_versions":         3,
			"oldest_version":       1,
			"updated_time":         inputUpdatedTimeStr,
			"versions": map[string]interface{}{
				"2": map[string]interface{}{
					"created_time":  inputUpdatedTimeStr,
					"deletion_time": "",
					"destroyed":     false,
				},
				"1": map[string]interface{}{
					"created_time":  inputCreatedTimeStr,
					"deletion_time": inputDeletedTimeStr,
					"destroyed":     false,
				},
			},
		},
	}

	testCases := []struct {
		name     string
		input    *Secret
		expected *KVMetadata
	}{
		{
			name:  "a metadata response",
			input: metadataResp,
			expected: &KVMetadata{
				CASRequired:    true,
				CreatedTime:    expectedCreatedTimeParsed,
				CurrentVersion: 2,
				CustomMetadata: map[string]interface{}{
					"org": "eng",
				},
				DeleteVersionAfter: time.Duration(200 * time.Second),
				MaxVersions:        3,
				OldestVersion:      1,
				UpdatedTime:        expectedUpdatedTimeParsed,
				Versions: map[string]KVVersionMetadata{
					"2": {
						Version:      2,
						CreatedTime:  expectedUpdatedTimeParsed,
						DeletionTime: time.Time{},
					},
					"1": {
						Version:      1,
						CreatedTime:  expectedCreatedTimeParsed,
						DeletionTime: expectedDeletedTimeParsed,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		md, err := extractFullMetadata(tc.input)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		if !reflect.DeepEqual(md, tc.expected) {
			t.Fatalf("%s: got\n%#v\nexpected\n%#v\n", tc.name, md, tc.expected)
		}
	}
}

func TestExtractCustomMetadata(t *testing.T) {
	testCases := []struct {
		name         string
		inputAPIResp *Secret
		expected     map[string]interface{}
	}{
		{
			name: "a read response with some custom metadata",
			inputAPIResp: &Secret{
				Data: map[string]interface{}{
					"metadata": map[string]interface{}{
						"custom_metadata": map[string]interface{}{"org": "eng"},
					},
				},
			},
			expected: map[string]interface{}{"org": "eng"},
		},
		{
			name: "a write response with some (pre-existing) custom metadata",
			inputAPIResp: &Secret{
				Data: map[string]interface{}{
					"custom_metadata": map[string]interface{}{"org": "eng"},
				},
			},
			expected: map[string]interface{}{"org": "eng"},
		},
		{
			name: "a read response with no custom metadata from a pre-1.9 Vault server",
			inputAPIResp: &Secret{
				Data: map[string]interface{}{
					"metadata": map[string]interface{}{},
				},
			},
			expected: map[string]interface{}(nil),
		},
		{
			name: "a write response with no custom metadata from a pre-1.9 Vault server",
			inputAPIResp: &Secret{
				Data: map[string]interface{}{},
			},
			expected: map[string]interface{}(nil),
		},
		{
			name: "a read response with no custom metadata from a post-1.9 Vault server",
			inputAPIResp: &Secret{
				Data: map[string]interface{}{
					"metadata": map[string]interface{}{
						"custom_metadata": nil,
					},
				},
			},
			expected: map[string]interface{}(nil),
		},
		{
			name: "a write response with no custom metadata from a post-1.9 Vault server",
			inputAPIResp: &Secret{
				Data: map[string]interface{}{
					"custom_metadata": nil,
				},
			},
			expected: map[string]interface{}(nil),
		},
		{
			name: "a read response where custom metadata was deleted",
			inputAPIResp: &Secret{
				Data: map[string]interface{}{
					"metadata": map[string]interface{}{
						"custom_metadata": map[string]interface{}{},
					},
				},
			},
			expected: map[string]interface{}{},
		},
		{
			name: "a write response where custom metadata was deleted",
			inputAPIResp: &Secret{
				Data: map[string]interface{}{
					"custom_metadata": map[string]interface{}{},
				},
			},
			expected: map[string]interface{}{},
		},
	}

	for _, tc := range testCases {
		cm := extractCustomMetadata(tc.inputAPIResp)

		if !reflect.DeepEqual(cm, tc.expected) {
			t.Fatalf("%s: got\n%#v\nexpected\n%#v\n", tc.name, cm, tc.expected)
		}
	}
}
