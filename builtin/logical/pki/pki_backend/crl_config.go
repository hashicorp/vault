// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki_backend

const latestCrlConfigVersion = 1

// CRLConfig holds basic CRL configuration information
type CrlConfig struct {
	Version                   int    `json:"version"`
	Expiry                    string `json:"expiry"`
	Disable                   bool   `json:"disable"`
	OcspDisable               bool   `json:"ocsp_disable"`
	AutoRebuild               bool   `json:"auto_rebuild"`
	AutoRebuildGracePeriod    string `json:"auto_rebuild_grace_period"`
	OcspExpiry                string `json:"ocsp_expiry"`
	EnableDelta               bool   `json:"enable_delta"`
	DeltaRebuildInterval      string `json:"delta_rebuild_interval"`
	UseGlobalQueue            bool   `json:"cross_cluster_revocation"`
	UnifiedCRL                bool   `json:"unified_crl"`
	UnifiedCRLOnExistingPaths bool   `json:"unified_crl_on_existing_paths"`
	MaxCRLEntries             int    `json:"max_crl_entries"`
}

// Implicit default values for the config if it does not exist.
var DefaultCrlConfig = CrlConfig{
	Version:                   latestCrlConfigVersion,
	Expiry:                    "72h",
	Disable:                   false,
	OcspDisable:               false,
	OcspExpiry:                "12h",
	AutoRebuild:               false,
	AutoRebuildGracePeriod:    "12h",
	EnableDelta:               false,
	DeltaRebuildInterval:      "15m",
	UseGlobalQueue:            false,
	UnifiedCRL:                false,
	UnifiedCRLOnExistingPaths: false,
	MaxCRLEntries:             100000,
}
