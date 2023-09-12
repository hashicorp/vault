// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"fmt"
	"strings"
)

type EabPolicyName string

const (
	eabPolicyNotRequired        EabPolicyName = "not-required"
	eabPolicyNewAccountRequired EabPolicyName = "new-account-required"
	eabPolicyAlwaysRequired     EabPolicyName = "always-required"
)

func getEabPolicyByString(name string) (EabPolicy, error) {
	lcName := strings.TrimSpace(strings.ToLower(name))
	switch lcName {
	case string(eabPolicyNotRequired):
		return getEabPolicyByName(eabPolicyNotRequired), nil
	case string(eabPolicyNewAccountRequired):
		return getEabPolicyByName(eabPolicyNewAccountRequired), nil
	case string(eabPolicyAlwaysRequired):
		return getEabPolicyByName(eabPolicyAlwaysRequired), nil
	default:
		return getEabPolicyByName(eabPolicyAlwaysRequired), fmt.Errorf("unknown eab policy name: %s", name)
	}
}

func getEabPolicyByName(name EabPolicyName) EabPolicy {
	return EabPolicy{Name: name}
}

type EabPolicy struct {
	Name EabPolicyName
}

// EnforceForNewAccount for new account creations, should we require an EAB.
func (ep EabPolicy) EnforceForNewAccount(eabData *eabType) error {
	if (ep.Name == eabPolicyAlwaysRequired || ep.Name == eabPolicyNewAccountRequired) && eabData == nil {
		return ErrExternalAccountRequired
	}

	return nil
}

// EnforceForExistingAccount for all operations within ACME, does the account being used require an EAB attached to it.
func (ep EabPolicy) EnforceForExistingAccount(account *acmeAccount) error {
	if ep.Name == eabPolicyAlwaysRequired && account.Eab == nil {
		return ErrExternalAccountRequired
	}

	return nil
}

// IsExternalAccountRequired for new accounts incoming does is an EAB required
func (ep EabPolicy) IsExternalAccountRequired() bool {
	return ep.Name == eabPolicyAlwaysRequired || ep.Name == eabPolicyNewAccountRequired
}

// OverrideEnvDisablingPublicAcme determines if ACME is enabled but the OS environment variable
// has said to disable public acme support, if we can override that environment variable to
// turn on ACME support
func (ep EabPolicy) OverrideEnvDisablingPublicAcme() bool {
	return ep.Name == eabPolicyAlwaysRequired
}
