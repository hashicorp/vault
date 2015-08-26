package aws

import (
	"testing"
)

func TestNormalizeDisplayName(t *testing.T) {
	invalidName := "^#$test name\nshould be normalized)(*"
	expectedName := "___test_name_should_be_normalized___"
	normalizedName := normalizeDisplayName(invalidName)
	if normalizedName != expectedName {
		t.Fatalf("normalizeDisplayName does not normalize AWS name correctly: %s", normalizedName)
	}

}
