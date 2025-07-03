// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package observations

import "github.com/hashicorp/go-hclog"

// ObservationSystemConfig serves as a definition of what an observation
// system needs from the Vault config file to start.
type ObservationSystemConfig struct {
	// LedgerPath is the path to the observation system's ledger.
	LedgerPath string `json:"ledger_path" hcl:"ledger_path"`

	// TypePrefixDenylist will deny any observations with types with matching prefixes
	// to be emitted to the ledger.
	TypePrefixDenylist []string `json:"type_prefix_denylist" hcl:"type_prefix_denylist"`

	// TypePrefixAllowlist will only allow observations with types with matching prefixes
	// to be emitted to the ledger.
	TypePrefixAllowlist []string `json:"type_prefix_allowlist" hcl:"type_prefix_allowlist"`

	// FileMode will attempt to open the ledger at the ledger path with the following
	// file mode. Specified as a string, but parsed as an octal, e.g. "0755".
	FileMode string `json:"file_mode" hcl:"file_mode"`
}

// NewObservationSystemConfig is the config for a new Observation System, provided
// to NewObservationSystem
type NewObservationSystemConfig struct {
	*ObservationSystemConfig
	LocalNodeId string
	Logger      hclog.Logger
}
