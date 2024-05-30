// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package experiments

import "slices"

const (
	VaultExperimentCoreAuditEventsAlpha1 = "core.audit.events.alpha1"
	VaultExperimentSecretsImport         = "secrets.import.alpha1"

	// Unused experiments. We keep them so that we don't break users who include them in their
	// flags or configs, but they no longer have any effect.
	VaultExperimentEventsAlpha1 = "events.alpha1"
)

var validExperiments = []string{
	VaultExperimentEventsAlpha1,
	VaultExperimentCoreAuditEventsAlpha1,
	VaultExperimentSecretsImport,
}

var unusedExperiments = []string{
	VaultExperimentEventsAlpha1,
}

// ValidExperiments exposes the list of valid experiments without exposing a mutable
// global variable. Experiments can only be enabled when starting a server, and will
// typically enable pre-GA API functionality.
func ValidExperiments() []string {
	return slices.Clone(validExperiments)
}

// IsUnused returns true if the given experiment is in the unused list.
func IsUnused(experiment string) bool {
	return slices.Contains(unusedExperiments, experiment)
}
