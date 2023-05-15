// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vault

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/vault/sdk/logical"
)

type acmeBillingSystemViewImpl struct {
	extendedSystemView
	logical.ManagedKeySystemView
	core *Core
}

var _ logical.ACMEBillingSystemView = (*acmeBillingSystemViewImpl)(nil)

func (c *Core) NewAcmeBillingSystemView(sysView interface{}) *acmeBillingSystemViewImpl {
	es := sysView.(extendedSystemView)
	managed, ok := sysView.(logical.ManagedKeySystemView)
	if !ok {
		return &acmeBillingSystemViewImpl{
			extendedSystemView: es,
			core:               c,
		}
	}

	return &acmeBillingSystemViewImpl{
		extendedSystemView:           es,
		ManagedKeySystemView:         managed,
		core:                         c,
	}
}

func (a *acmeBillingSystemViewImpl) CreateActivityCountEventForIdentifiers(ctx context.Context, identifiers []string) error {
	sort.Strings(identifiers)
	return fmt.Errorf("got identifiers but not implemented: %v", identifiers)
}
