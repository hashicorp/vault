// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package database

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-uuid"
	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
)

// GetConnectionWithConfig gets or creates a database connection with the given config for community edition
func (b *databaseBackend) GetConnectionWithConfig(ctx context.Context, name string, config *DatabaseConfig) (*dbPluginInstance, error) {
	// fast path, reuse the existing connection
	dbi := b.connections.Get(name)
	if dbi != nil {
		return dbi, nil
	}

	// slow path, create a new connection
	// if we don't lock the rest of the operation, there is a race condition for multiple callers of this function
	b.createConnectionLock.Lock()
	defer b.createConnectionLock.Unlock()

	// check again in case we lost the race
	dbi = b.connections.Get(name)
	if dbi != nil {
		return dbi, nil
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	// Override the configured version if there is a pinned version.
	pinnedVersion, err := b.getPinnedVersion(ctx, config.PluginName)
	if err != nil {
		return nil, err
	}
	pluginVersion := config.PluginVersion
	if pinnedVersion != "" {
		pluginVersion = pinnedVersion
	}

	dbw, err := newDatabaseWrapper(ctx, config.PluginName, pluginVersion, b.System(), b.logger)
	if err != nil {
		return nil, fmt.Errorf("unable to create database instance: %w", err)
	}

	initReq := v5.InitializeRequest{
		Config:           config.ConnectionDetails,
		VerifyConnection: config.VerifyConnection,
	}
	_, err = dbw.Initialize(ctx, initReq)
	if err != nil {
		dbw.Close()
		return nil, err
	}

	dbi = &dbPluginInstance{
		database:             dbw,
		id:                   id,
		name:                 name,
		runningPluginVersion: pluginVersion,
	}
	conn, ok := b.connections.PutIfEmpty(name, dbi)
	if !ok {
		// this is a bug
		b.Logger().Warn("BUG: there was a race condition adding to the database connection map")
		// There was already an existing connection, so we will use that and close our new one to avoid a race condition.
		err := dbi.Close()
		if err != nil {
			b.Logger().Warn("Error closing new database connection", "error", err)
		}
	}
	return conn, nil
}
