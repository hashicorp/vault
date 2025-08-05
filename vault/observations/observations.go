// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package observations

// ObservationSystemConfig serves as a definition of what an observation
// system needs from the Vault config file to start.
type ObservationSystemConfig struct {
	// LedgerPath is the path to the observation system's ledger.
	LedgerPath string `json:"ledger_path" hcl:"ledger_path"`
}
