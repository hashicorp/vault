// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"

	"github.com/hashicorp/vault/helper/namespace"
)

var NamespaceByID func(context.Context, string, *Core) (*namespace.Namespace, error) = namespaceByID

func namespaceByID(ctx context.Context, nsID string, c *Core) (*namespace.Namespace, error) {
	if nsID == namespace.RootNamespaceID {
		return namespace.RootNamespace, nil
	}
	return nil, namespace.ErrNoNamespace
}

var NamespaceRegister func(context.Context, *namespace.Namespace, *Core) error = namespaceRegister

func namespaceRegister(ctx context.Context, ns *namespace.Namespace, c *Core) error {
	return nil
}
