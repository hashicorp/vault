// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package loadedsnapshot

import (
	"github.com/hashicorp/vault/sdk/logical"
)

// Manager is responsible for managing loaded snapshots in the system. It
// handles the lifecycle of snapshots, including their loading, storing them on the disk, and
// expiration.
type Manager struct{}

// NewManager creates a new snapshot manager.
func NewManager(s logical.Storage, raftDataDirPath string, clusterID string) *Manager {
	return &Manager{}
}

// Shutdown is called to clean up any resources used by the manager before shutting Vault down.
func (m *Manager) Shutdown() {
	// Close background tasks here
}

// Start is called to initialize the manager and start any background tasks.
func (m *Manager) Start() {
	// Run background tasks here
}
