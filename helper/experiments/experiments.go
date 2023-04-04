// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package experiments

const VaultExperimentEventsAlpha1 = "events.alpha1"

var validExperiments = []string{
	VaultExperimentEventsAlpha1,
}

// ValidExperiments exposes the list without exposing a mutable global variable.
// Experiments can only be enabled when starting a server, and will typically
// enable pre-GA API functionality.
func ValidExperiments() []string {
	result := make([]string, len(validExperiments))
	copy(result, validExperiments)
	return result
}
