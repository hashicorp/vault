// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package vault

import (
	"github.com/hashicorp/vault/helper/namespace"
)

func (i *IdentityStore) listNamespaces() []*namespace.Namespace {
	return []*namespace.Namespace{namespace.RootNamespace}
}
