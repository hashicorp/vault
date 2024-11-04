// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package transit

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
)

// addEntSignFieldArgs adds the enterprise only fields to the field schema definition
func addEntSignFieldArgs(_ map[string]*framework.FieldSchema) {
	// Do nothing
}

// addEntVerifyFieldArgs adds the enterprise only fields to the field schema definition
func addEntVerifyFieldArgs(_ map[string]*framework.FieldSchema) {
	// Do nothing
}

// validateSignApiArgsVersionSpecific will perform a validation of the Sign API parameters
// from the Enterprise or CE point of view.
func validateSignApiArgsVersionSpecific(p *keysutil.Policy, apiArgs commonSignVerifyApiArgs) error {
	if err := _validateEntSpecificKeyType(p); err != nil {
		return err
	}

	if apiArgs.hashAlgorithm == keysutil.HashTypeNone {
		if !apiArgs.prehashed || apiArgs.sigAlgorithm != "pkcs1v15" {
			return fmt.Errorf("hash_algorithm=none requires both prehashed=true and signature_algorithm=pkcs1v15")
		}
	}

	return nil
}

// populateEntPolicySigning augments or tweaks the input parameters to the SDK policy.SignWithOptions for
// Enterprise usage.
func (b *backend) populateEntPolicySigningOptions(_ context.Context, _ *keysutil.Policy, _ signApiArgs, _ batchRequestSignItem, _ *policySignArgs) error {
	return nil
}

// populateEntPolicyVerifyOptions augments or tweaks the input parameters to the SDK policy.VerifyWithOptions for
// Enterprise usage.
func (b *backend) populateEntPolicyVerifyOptions(ctx context.Context, p *keysutil.Policy, args verifyApiArgs, item batchRequestVerifyItem, vsa *policyVerifyArgs) error {
	return nil
}

func _validateEntSpecificKeyType(p *keysutil.Policy) error {
	switch p.Type {
	case keysutil.KeyType_AES128_CMAC, keysutil.KeyType_AES256_CMAC, keysutil.KeyType_MANAGED_KEY:
		return fmt.Errorf("enterprise specific key type %q can not be used on CE", p.Type)
	default:
		return nil
	}
}
