// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package transit

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/helper/keysutil"
)

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
func (b *backend) populateEntPolicySigningOptions(_ context.Context, p *keysutil.Policy, args signApiArgs, item batchRequestSignItem, _ *policySignArgs) error {
	return _forbidEd25519EntBehavior(p, args.commonSignVerifyApiArgs, item["signature_context"])
}

// populateEntPolicyVerifyOptions augments or tweaks the input parameters to the SDK policy.VerifyWithOptions for
// Enterprise usage.
func (b *backend) populateEntPolicyVerifyOptions(_ context.Context, p *keysutil.Policy, args verifyApiArgs, item batchRequestVerifyItem, _ *policyVerifyArgs) error {
	sigContext, err := _validateString(item, "signature_context")
	if err != nil {
		return err
	}
	return _forbidEd25519EntBehavior(p, args.commonSignVerifyApiArgs, sigContext)
}

func _validateString(item batchRequestVerifyItem, key string) (string, error) {
	if itemVal, exists := item[key]; exists {
		if itemStrVal, ok := itemVal.(string); ok {
			return itemStrVal, nil
		}
		return "", fmt.Errorf("expected string for key=%q, got=%q", key, itemVal)
	}
	return "", nil
}

func _forbidEd25519EntBehavior(p *keysutil.Policy, apiArgs commonSignVerifyApiArgs, sigContext string) error {
	if p.Type != keysutil.KeyType_ED25519 {
		return nil
	}

	switch {
	case apiArgs.prehashed:
		return fmt.Errorf("only Pure Ed25519 signatures supported, prehashed must be false")
	case apiArgs.hashAlgorithm == keysutil.HashTypeSHA2512:
		return fmt.Errorf("only Pure Ed25519 signatures supported, hash_alogithm should not be set")
	case sigContext != "":
		return fmt.Errorf("only Pure Ed25519 signatures supported, signature_context must be empty")
	}

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
