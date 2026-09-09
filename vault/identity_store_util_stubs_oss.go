// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/vault/helper/identity"
)

func (i *IdentityStore) MemDBAliasByIssuerAndExternalId(_, _, _ string, _ bool) (*identity.Alias, error) {
	return nil, nil
}

func (i *IdentityStore) MemDBAliasByIssuerAndExternalIdInTxn(_ *memdb.Txn, _, _, _ string, _ bool) (*identity.Alias, error) {
	return nil, nil
}
