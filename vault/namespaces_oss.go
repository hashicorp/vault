// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/helper/namespace"
)

func createNamespace(_ context.Context, _ *Core, _ string, _ map[string]string) (*namespace.Namespace, error) {
	return nil, errors.New("namespaces are a Vault Enterprise feature")
}

func (c *Core) NamespaceByID(ctx context.Context, nsID string) (*namespace.Namespace, error) {
	return namespaceByID(ctx, nsID, c)
}

func (c *Core) ListNamespaces(includePath bool) []*namespace.Namespace {
	return []*namespace.Namespace{namespace.RootNamespace}
}

func (c *Core) resetNamespaceCache() {}
