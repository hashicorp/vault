// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) getReadLockedPolicy(ctx context.Context, s logical.Storage, name string) (*keysutil.Policy, error) {
	p, _, err := b.GetPolicy(ctx, keysutil.PolicyRequest{
		Storage: s,
		Name:    name,
	}, b.GetRandomReader())
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("%w: key %s not found", logical.ErrInvalidRequest, name)
	}
	if !b.System().CachingDisabled() {
		p.Lock(false)
	}
	return p, nil
}

// runWithReadLockedPolicy runs a function passing in the policy specified by keyName that has been
// locked in a read only fashion without the ability to upsert the policy
func (b *backend) runWithReadLockedPolicy(ctx context.Context, s logical.Storage, keyName string, f func(p *keysutil.Policy) (*logical.Response, error)) (*logical.Response, error) {
	p, err := b.getReadLockedPolicy(ctx, s, keyName)
	if err != nil {
		if errors.Is(err, logical.ErrInvalidRequest) {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		return nil, err
	}
	defer p.Unlock()
	return f(p)
}

// validateKeyVersion verifies that the passed in key version is valid for our
// current key policy, returning correct version to use within the policy.
func validateKeyVersion(p *keysutil.Policy, ver int) (int, error) {
	switch {
	case ver < 0:
		return 0, fmt.Errorf("cannot use negative key version %d", ver)
	case ver == 0:
		// Allowed, will use latest; set explicitly here to ensure the string
		// is generated properly
		ver = p.LatestVersion
	case ver == p.LatestVersion:
		// Allowed
	case p.MinEncryptionVersion > 0 && ver < p.MinEncryptionVersion:
		return 0, fmt.Errorf("cannot use key version %d: version is too old (disallowed by policy) for key %s", ver, p.Name)
	}
	return ver, nil
}
