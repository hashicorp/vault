// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package consts

const (
	// This PEN is assigned to HashiCorp
	hashicorpPEN = "1.3.6.1.4.1.55813"
	// Reserve everything under .1 for Vault
	vaultOIDBase = hashicorpPEN + ".1"
	// TPMAuthOIDBase is the OID prefix reserved for the tpm auth method
	TPMAuthOIDBase     = vaultOIDBase + ".1"
	TPMAuthOIDTPMID    = TPMAuthOIDBase + ".1"
	TPMAuthOIDRoleName = TPMAuthOIDBase + ".2"
)
