// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package experiments

const (
	VaultExperimentCoreAuditEventsAlpha1 = "core.audit.events.alpha1"
	VaultExperimentSecretsSyncAlpha1     = "secrets.sync.alpha1"

	// Unused experiments. We keep them so that we don't break users who include them in their
	// flags or configs, but they no longer have any effect.
	VaultExperimentEventsAlpha1 = "events.alpha1"
)

var validExperiments = []string{
	VaultExperimentEventsAlpha1,
	VaultExperimentCoreAuditEventsAlpha1,
	VaultExperimentSecretsSyncAlpha1,
}

// ValidExperiments exposes the list without exposing a mutable global variable.
// Experiments can only be enabled when starting a server, and will typically
// enable pre-GA API functionality.
func ValidExperiments() []string {
	result := make([]string, len(validExperiments))
	copy(result, validExperiments)
	return result
}
