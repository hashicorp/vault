// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import "github.com/hashicorp/vault/sdk/framework"

// entAddAuthTuneRequestFields is a stub implementation for CE
func entAddAuthTuneRequestFields(fields map[string]*framework.FieldSchema) {
	// No additional fields in CE
}

// entAddAuthTuneResponseFields is a stub implementation for CE
func entAddAuthTuneResponseFields(fields map[string]*framework.FieldSchema) {
	// No additional fields in CE
}

// entAddAuthRequestFields is a stub implementation for CE
func entAddAuthRequestFields(fields map[string]*framework.FieldSchema) {
	// No additional fields in CE
}

// entAddAuthResponseFields is a stub implementation for CE
func entAddAuthResponseFields(fields map[string]*framework.FieldSchema) {
	// No additional fields in CE
}

// entAddSecretsTuneRequestFields is a stub implementation for CE
func entAddSecretsTuneRequestFields(fields map[string]*framework.FieldSchema) {
	// No additional fields in CE
}

// entAddSecretsTuneResponseFields is a stub implementation for CE
func entAddSecretsTuneResponseFields(fields map[string]*framework.FieldSchema) {
	// No additional fields in CE
}

// entAddSecretsRequestFields is a stub implementation for CE
func entAddSecretsRequestFields(fields map[string]*framework.FieldSchema) {
	// No additional fields in CE
}

// entAddSecretsResponseFields is a stub implementation for CE
func entAddSecretsResponseFields(fields map[string]*framework.FieldSchema) {
	// No additional fields in CE
}
