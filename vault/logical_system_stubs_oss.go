// Copyright IBM Corp. 2016, 2025
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

type entSystemBackend struct{}

func entUnauthenticatedPaths() []string {
	return []string{}
}

func entBinaryPaths() []string {
	return []string{}
}

func entLocalStoragePaths() []string {
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

// mountInfo returns a map of information about the given mount entry
// Enterprise-specific fields are added in the enterprise version of this method.
func (b *SystemBackend) mountInfo(ctx context.Context, entry *MountEntry, legacyTTLFormat bool) map[string]interface{} {
	info, entryConfig := b.internalMountInfoCommon(ctx, entry, legacyTTLFormat)
	info["config"] = entryConfig
	return info
}

func (b *SystemBackend) callUnsyncMountHelper(ctx context.Context, path string) error {
	return nil
}
