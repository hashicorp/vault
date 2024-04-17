// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/logical"
)

// Our forwarded writer has two components: a reference to Core, allowing
// us to tap into the GRPC client and resolved paths, and lower storage
// layer to call upon when we don't wish to forward our writes.
//
// This implementation lives in OSS: while the GRPC connection isn't present
// on OSS, we need to ensure paths written to these forwarded nodes correctly
// template {{clusterId}} if they are later upgraded to Enterprise, and don't
// just write with the template sentinel still there.
//
// XXX: In the future, we'll need to support wrapping transactional storage.
type ForwardedWriter struct {
	core      *Core
	lower     logical.Storage
	clusterID string
}

func (c *Core) NewForwardedWriter(ctx context.Context, wrapped logical.Storage, _ bool /* local */) (logical.Storage, error) {
	// local is unused above on this OSS implementation: local mounts only
	// exist on Vault Enterprise.

	// Cache the cluster id; we assume we'll be recreated when plugins reload
	// if this changes, and should not change without reloading plugins.
	cluster, err := c.Cluster(ctx)
	if err != nil || cluster.ID == "" {
		return nil, fmt.Errorf("failed to fetch local cluster info: %v", err)
	}

	return &ForwardedWriter{
		core:      c,
		lower:     wrapped,
		clusterID: cluster.ID,
	}, nil
}

func (w *ForwardedWriter) List(ctx context.Context, path string) ([]string, error) {
	// storage.List(...) operations are always handled locally. However, we
	// may need to resolve any {{clusterId}} template sentinels if given to us
	// and we'd otherwise consider this a forwarded write operation.
	var err error
	path, err = w.resolvePathIfNecessary(path)
	if err != nil {
		return nil, fmt.Errorf("failed to do local cross-cluster list: failed to resolve path: %w", err)
	}

	return w.lower.List(ctx, path)
}

func (w *ForwardedWriter) Get(ctx context.Context, path string) (*logical.StorageEntry, error) {
	// See note in List(...)above.
	var err error
	path, err = w.resolvePathIfNecessary(path)
	if err != nil {
		return nil, fmt.Errorf("failed to do local cross-cluster read: failed to resolve path: %w", err)
	}

	return w.lower.Get(ctx, path)
}

func (w *ForwardedWriter) Put(ctx context.Context, entry *logical.StorageEntry) error {
	// See note above about List(...).
	var err error
	entry.Key, err = w.resolvePathIfNecessary(entry.Key)
	if err != nil {
		return fmt.Errorf("failed to do local cross-cluster write: failed to resolve path: %w", err)
	}

	return w.lower.Put(ctx, entry)
}

func (w *ForwardedWriter) Delete(ctx context.Context, path string) error {
	// See note above about List(...).
	var err error
	path, err = w.resolvePathIfNecessary(path)
	if err != nil {
		return fmt.Errorf("failed to do local cross-cluster delete: failed to resolve path: %w", err)
	}
	return w.lower.Delete(ctx, path)
}

func (w *ForwardedWriter) resolvePathIfNecessary(path string) (string, error) {
	// We should only resolve this path when we're going to be servicing
	// it locally.
	//
	// We don't bother checking if we're a perf primary or not, as even
	// perf secondaries could use locally serviced operations on these paths
	// (e.g., a storage.List(...)).
	forwardablePath := w.core.writeForwardedPaths.HasPath(path)

	if forwardablePath {
		return w.resolvePath(path)
	}

	return path, nil
}

func (w *ForwardedWriter) resolvePath(path string) (string, error) {
	// This is the source-agnostic path resolution helper. Here we ensure
	// we've got a forwarded path (one that contains the proper UUID
	// sentinel) and we fetch this cluster's UUID and update the path.
	if !strings.Contains(path, logical.PBPWFClusterSentinel) {
		return "", fmt.Errorf("invalid path: lacks '%v' sentinel for expansion", logical.PBPWFClusterSentinel)
	}

	return strings.Replace(path, logical.PBPWFClusterSentinel, w.clusterID, 1), nil
}
