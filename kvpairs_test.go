package main

import (
	"testing"

	"github.com/hashicorp/vault/sdk/framework"
)

func TestTypeKVPairs(t *testing.T) {
	fields := map[string]*framework.FieldSchema{
		"metadata": {
			Type: framework.TypeKVPairs,
		},
	}

	testCases := []struct {
		name          string
		input         string
		expectedCount int
		expectedPairs map[string]string
	}{
		{
			name:          "Single pair",
			input:         "A=a",
			expectedCount: 1,
			expectedPairs: map[string]string{"A": "a"},
		},
		{
			name:          "Multiple pairs - THIS WILL FAIL",
			input:         "A=a,B=b,C=c",
			expectedCount: 3,
			expectedPairs: map[string]string{
				"A": "a",
				"B": "b",
				"C": "c",
			},
		},
		{
			name:          "Two pairs",
			input:         "key1=value1,key2=value2",
			expectedCount: 2,
			expectedPairs: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			rawData := map[string]interface{}{
				"metadata": tc.input,
			}

			fieldData := &framework.FieldData{
				Raw:    rawData,
				Schema: fields,
			}

			result := fieldData.Get("metadata")
			kvMap, ok := result.(map[string]string)

			if !ok {
				t.Fatal("Failed to convert to map[string]string")
			}

			if len(kvMap) != tc.expectedCount {
				t.Errorf("Expected %d pairs, got %d", tc.expectedCount, len(kvMap))
				t.Logf("Actual map: %v", kvMap)
			}

			for expectedKey, expectedValue := range tc.expectedPairs {
				actualValue, exists := kvMap[expectedKey]
				if !exists {
					t.Errorf("Expected key '%s' not found in result", expectedKey)
				} else if actualValue != expectedValue {
					t.Errorf("For key '%s': expected value '%s', got '%s'",
						expectedKey, expectedValue, actualValue)

				}
			}
		})
	}
}
