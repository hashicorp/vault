// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/snapshots"
)

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

type entSystemBackend struct{}

func entUnauthenticatedPaths() []string {
	return []string{}
}

func (s *SystemBackend) entInit() {}

func (s *SystemBackend) makeSnapshotSource(ctx context.Context, _ *framework.FieldData) (snapshots.Source, error) {
	body, ok := logical.ContextOriginalBodyValue(ctx)
	if !ok {
		return nil, errors.New("no reader for request")
	}
	return snapshots.NewManualSnapshotSource(body), nil
}

func entWrappedPluginsCRUDPath(b *SystemBackend) []*framework.Path {
	return []*framework.Path{b.pluginsCatalogCRUDPath()}
}
