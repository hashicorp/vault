// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package openldap

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
)

type dynamicRole struct {
	// required fields
	Name         string `json:"name"          mapstructure:"name"`
	CreationLDIF string `json:"creation_ldif" mapstructure:"creation_ldif"`
	DeletionLDIF string `json:"deletion_ldif" mapstructure:"deletion_ldif"`

	// optional fields
	RollbackLDIF     string        `json:"rollback_ldif"               mapstructure:"rollback_ldif,omitempty"`
	UsernameTemplate string        `json:"username_template,omitempty" mapstructure:"username_template,omitempty"`
	DefaultTTL       time.Duration `json:"default_ttl,omitempty"       mapstructure:"default_ttl,omitempty"`
	MaxTTL           time.Duration `json:"max_ttl,omitempty"           mapstructure:"max_ttl,omitempty"`
}

func retrieveDynamicRole(ctx context.Context, s logical.Storage, roleName string) (*dynamicRole, error) {
	entry, err := s.Get(ctx, path.Join(dynamicRolePath, roleName))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	result := new(dynamicRole)
	if err := entry.DecodeJSON(result); err != nil {
		return nil, err
	}

	return result, nil
}

func storeDynamicRole(ctx context.Context, s logical.Storage, role *dynamicRole) error {
	if role.Name == "" {
		return fmt.Errorf("missing role name")
	}
	entry, err := logical.StorageEntryJSON(path.Join(dynamicRolePath, role.Name), role)
	if err != nil {
		return fmt.Errorf("unable to marshal storage entry: %w", err)
	}

	err = s.Put(ctx, entry)
	if err != nil {
		return fmt.Errorf("failed to store dynamic role: %w", err)
	}
	return nil
}

func deleteDynamicRole(ctx context.Context, s logical.Storage, roleName string) error {
	if roleName == "" {
		return fmt.Errorf("missing role name")
	}
	return s.Delete(ctx, path.Join(dynamicRolePath, roleName))
}
