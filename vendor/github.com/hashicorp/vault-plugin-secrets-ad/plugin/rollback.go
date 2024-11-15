// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

const (
	rotateCredentialWAL = "rotateCredentialWAL"
)

// rotateCredentialEntry is used to store information in a WAL that can retry a
// credential rotation in the event of partial failure.
type rotateCredentialEntry struct {
	LastVaultRotation  time.Time `json:"last_vault_rotation"`
	LastPassword       string    `json:"last_password"`
	CurrentPassword    string    `json:"current_password"`
	RoleName           string    `json:"name"`
	ServiceAccountName string    `json:"service_account_name"`
	TTL                int       `json:"ttl"`
}

func (b *backend) walRollback(ctx context.Context, req *logical.Request, kind string, data interface{}) error {
	switch kind {
	case rotateCredentialWAL:
		return b.handleRotateCredentialRollback(ctx, req.Storage, data)
	default:
		return fmt.Errorf("unknown WAL entry kind %q", kind)
	}
}

func (b *backend) handleRotateCredentialRollback(ctx context.Context, storage logical.Storage, data interface{}) error {
	var wal rotateCredentialEntry
	if err := mapstructure.WeakDecode(data, &wal); err != nil {
		return err
	}

	if wal.CurrentPassword == "" {
		b.Logger().Warn("WAL does not contain a password for service account")
		return nil
	}

	// Check creds for deltas. Exit if creds and WAL are the same.
	path := fmt.Sprintf("%s/%s", storageKey, wal.RoleName)
	credEntry, err := storage.Get(ctx, path)
	if err == nil && credEntry != nil {
		cred := make(map[string]interface{})
		err := credEntry.DecodeJSON(&cred)
		if err == nil && cred != nil {
			currentPassword := cred["current_password"]
			lastPassword := cred["last_password"]

			if currentPassword == wal.CurrentPassword && lastPassword == wal.LastPassword {
				return nil
			}
		}
	}

	role := &backendRole{
		ServiceAccountName: wal.ServiceAccountName,
		TTL:                wal.TTL,
		LastVaultRotation:  wal.LastVaultRotation,
	}

	if err := b.writeRoleToStorage(ctx, storage, wal.RoleName, role); err != nil {
		return err
	}

	// Cache the full role to minimize Vault storage calls.
	b.roleCache.SetDefault(wal.RoleName, role)

	conf, err := readConfig(ctx, storage)
	if err != nil {
		return err
	}
	if conf == nil {
		return errors.New("the config is currently unset")
	}

	if err := b.client.UpdatePassword(conf.ADConf, role.ServiceAccountName, wal.CurrentPassword); err != nil {
		return err
	}

	// Although a service account name is typically my_app@example.com,
	// the username it uses is just my_app, or everything before the @.
	username, err := getUsername(role.ServiceAccountName)
	if err != nil {
		return err
	}

	b.credLock.Lock()
	defer b.credLock.Unlock()

	cred := map[string]interface{}{
		"username":         username,
		"current_password": wal.CurrentPassword,
	}

	if wal.LastPassword != "" {
		cred["last_password"] = wal.LastPassword
	}

	// Cache and save the cred.
	path = fmt.Sprintf("%s/%s", storageKey, wal.RoleName)
	entry, err := logical.StorageEntryJSON(path, cred)
	if err != nil {
		return err
	}
	if err := storage.Put(ctx, entry); err != nil {
		return err
	}

	b.credCache.SetDefault(wal.RoleName, cred)

	return nil
}
