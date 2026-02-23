// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"reflect"
	"testing"
)

func TestApplyRetentionPolicy_CoreRules_N_N1_N2_N3(t *testing.T) {
	req := &ListUpdatedVersionsReq{}

	// Setup: versions from 1.0.x to 1.3.x → latest = 1.3.x
	allVersions := map[string]*Version{
		"1.0.x": {CEActive: false, LTS: false}, // n-3 → REMOVED (not LTS)
		"1.1.x": {CEActive: false, LTS: true},  // n-2 → KEPT (always)
		"1.2.x": {CEActive: false, LTS: false}, // n-1 → KEPT
		"1.3.x": {CEActive: false, LTS: false}, // n → ce_active = true
	}

	newInputVersions := []string{"1.3.x"}

	got := req.applyRetentionPolicy(allVersions, newInputVersions)

	want := map[string]*Version{
		"1.3.x": {CEActive: true, LTS: false},  // n: ce_active = true
		"1.2.x": {CEActive: false, LTS: false}, // n-1: kept
		"1.1.x": {CEActive: false, LTS: true},  // n-2: kept (always)
		// "1.0.x": removed → not in output
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("core retention policy failed:\n got: %#v\nwant: %#v", got, want)
	}
}

func TestApplyRetentionPolicy_N3_Kept_WhenLTS(t *testing.T) {
	req := &ListUpdatedVersionsReq{}

	allVersions := map[string]*Version{
		"1.0.x": {CEActive: false, LTS: true},
		"1.1.x": {CEActive: false, LTS: false},
		"1.2.x": {CEActive: false, LTS: false},
		"1.3.x": {CEActive: false, LTS: false},
	}

	newInputVersions := []string{"1.3.x"}

	got := req.applyRetentionPolicy(allVersions, newInputVersions)

	want := map[string]*Version{
		"1.3.x": {CEActive: true, LTS: false},
		"1.2.x": {CEActive: false, LTS: false},
		"1.1.x": {CEActive: false, LTS: false},
		"1.0.x": {CEActive: false, LTS: true},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("n-3 LTS retention failed:\n got: %#v\nwant: %#v", got, want)
	}
}
