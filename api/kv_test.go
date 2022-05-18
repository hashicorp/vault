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
				Version:        10,
				CreatedTime:    expectedCreatedTimeParsed,
				DeletionTime:   time.Time{},
				Destroyed:      false,
				CustomMetadata: nil,
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
				Version:        10,
				CreatedTime:    expectedCreatedTimeParsed,
				DeletionTime:   expectedDeletionTimeParsed,
				Destroyed:      false,
				CustomMetadata: nil,
			},
		},
		{
			name: "a secret with custom metadata",
			input: &Secret{
				Data: map[string]interface{}{
					"data": map[string]interface{}{
						"password": "Hashi123",
					},
					"metadata": map[string]interface{}{
						"version":       10,
						"created_time":  inputCreatedTimeStr,
						"deletion_time": "",
						"destroyed":     false,
						"custom_metadata": map[string]string{
							"foo": "abc",
							"bar": "def",
							"baz": "ghi",
						},
					},
				},
			},
			expected: &KVVersionMetadata{
				Version:      10,
				CreatedTime:  expectedCreatedTimeParsed,
				DeletionTime: time.Time{},
				Destroyed:    false,
				CustomMetadata: map[string]string{
					"foo": "abc",
					"bar": "def",
					"baz": "ghi",
				},
			},
		},
	}

	for _, tc := range testCases {
		versionMetadata, err := extractVersionMetadata(tc.input)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		if !reflect.DeepEqual(versionMetadata, tc.expected) {
			t.Fatalf("got\n%#v\nexpected\n%#v\n", versionMetadata, tc.expected)
		}
	}
}
