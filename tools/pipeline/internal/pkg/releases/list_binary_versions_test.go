// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"strings"
	"testing"
)

func TestNewListBinaryVersionsReq(t *testing.T) {
	// Verify that the constructor returns a valid object
	req := NewListBinaryVersionsReq("1.2.3 4.5.6-ent")
	if req == nil {
		t.Fatal("returned nil")
	}

	// Ensure the input string is stored without modification
	if req.VersionsString != "1.2.3 4.5.6-ent" {
		t.Errorf("got %q", req.VersionsString)
	}
}

// Unit test: JSON output and String() summary formatting
func TestListBinaryVersionsRes_ToJSON_and_String(t *testing.T) {
	// Create a fake response object
	res := &ListBinaryVersionsRes{
		ValidVersions: map[string]struct {
			Status   string        `json:"status"`
			Variants []VariantInfo `json:"variants"`
		}{
			"1.17.0-ent": {
				Status: "valid",
				Variants: []VariantInfo{
					{Variant: "1.17.0+ent", OS: []string{"linux", "darwin"}},
					{Variant: "1.17.0+ent.hsm", OS: []string{"linux"}},
				},
			},
		},
		InvalidVersions: []string{"9.9.9"},
		AllVersions:     []string{"1.17.0-ent", "9.9.9"},
	}

	// Test: ToJSON output
	b, err := res.ToJSON()
	if err != nil {
		t.Fatal(err)
	}

	// Ensure the JSON contains known variant values
	if !strings.Contains(string(b), "1.17.0+ent.hsm") {
		t.Error("JSON missing variant")
	}

	// -------------------------
	// Test: String() summary
	// -------------------------
	// Expected values based on the test fixture:
	//   total requested: 2
	//   valid versions: 1
	//   total variants: 2
	//   missing: 1
	expected := "2 → 1 valid (2 variants), 1 missing"

	// Check that the summary contains the correct computed text.
	if !strings.Contains(res.String(), expected) {
		t.Errorf("String() = %q, want substring %q", res.String(), expected)
	}
}
