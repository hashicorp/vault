// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"time"
)

type backendRole struct {
	ServiceAccountName string    `json:"service_account_name"`
	TTL                int       `json:"ttl"`
	LastVaultRotation  time.Time `json:"last_vault_rotation"`
	PasswordLastSet    time.Time `json:"password_last_set"`
}

func (r *backendRole) Map() map[string]interface{} {
	m := map[string]interface{}{
		"service_account_name": r.ServiceAccountName,
		"ttl":                  r.TTL,
	}

	var unset time.Time
	if r.LastVaultRotation != unset {
		m["last_vault_rotation"] = r.LastVaultRotation
	}
	if r.PasswordLastSet != unset {
		m["password_last_set"] = r.PasswordLastSet
	}
	return m
}
