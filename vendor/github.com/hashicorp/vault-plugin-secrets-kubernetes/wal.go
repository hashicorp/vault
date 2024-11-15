// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kubesecrets

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

const (
	walRoleKind    = "role"
	walBindingKind = "roleBinding"
)

// Eventually expire the WAL if for some reason the rollback operation consistently fails
var maxWALAge = 24 * time.Hour

func (b *backend) walRollback(ctx context.Context, req *logical.Request, kind string, data interface{}) error {
	switch kind {
	case walRoleKind:
		return b.rollbackRoleWAL(ctx, req, data)
	case walBindingKind:
		return b.rollbackRoleBindingWAL(ctx, req, data)
	default:
		return fmt.Errorf("unknown rollback type %q", kind)
	}
}

type walRole struct {
	Namespace  string
	Name       string
	RoleType   string
	Expiration time.Time
}

// rollbackRoleWAL uses the info in a walRole entry to delete a Role/ClusterRole
// from Kubernetes. We're relying on Kubernetes garbage collection to delete the
// other related objects (RoleBinding/ClusterRoleBinding and ServiceAccount)
// since they should have an owner reference to the Role/ClusterRole
func (b *backend) rollbackRoleWAL(ctx context.Context, req *logical.Request, data interface{}) error {
	// Decode the WAL data
	var entry walRole
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeHookFunc(time.RFC3339),
		Result:     &entry,
	})
	if err != nil {
		return err
	}
	err = d.Decode(data)
	if err != nil {
		return err
	}

	client, err := b.getClient(ctx, req.Storage)
	if err != nil {
		return err
	}

	b.Logger().Debug("rolling back", "role", entry.RoleType, "namespace", entry.Namespace, "name", entry.Name)

	// Attempt to delete the Role. If we don't succeed within maxWALAge (e.g.
	// client creds are somehow incorrect and the delete will never succeed),
	// unconditionally remove the WAL.
	if err := client.deleteRole(ctx, entry.Namespace, entry.Name, entry.RoleType); err != nil {
		b.Logger().Warn("rollback error deleting", "roleType", entry.RoleType, "namespace", entry.Namespace, "name", entry.Name, "err", err)

		if time.Now().After(entry.Expiration) {
			b.Logger().Warn("giving up deleting", "roleType", entry.RoleType, "namespace", entry.Namespace, "name", entry.Name)
			return nil
		}
		return err
	}

	return nil
}

type walRoleBinding struct {
	Namespace  string
	Name       string
	IsCluster  bool
	Expiration time.Time
}

// rollbackRoleBindingWAL uses the info in a walRole entry to delete a
// Role/ClusterRole from Kubernetes. We're relying on Kubernetes garbage
// collection to delete the related ServiceAccount since it should have an owner
// reference to the RoleBinding/ClusterRoleBinding
func (b *backend) rollbackRoleBindingWAL(ctx context.Context, req *logical.Request, data interface{}) error {
	// Decode the WAL data
	var entry walRoleBinding
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeHookFunc(time.RFC3339),
		Result:     &entry,
	})
	if err != nil {
		return err
	}
	err = d.Decode(data)
	if err != nil {
		return err
	}

	client, err := b.getClient(ctx, req.Storage)
	if err != nil {
		return err
	}

	b.Logger().Debug("rolling back role binding", "isClusterRoleBinding", entry.IsCluster, "namespace", entry.Namespace, "name", entry.Name)

	// Attempt to delete the RoleBinding. If we don't succeed within maxWALAge
	// (e.g. client creds are somehow incorrect and the delete will never
	// succeed), unconditionally remove the WAL.
	if err := client.deleteRoleBinding(ctx, entry.Namespace, entry.Name, entry.IsCluster); err != nil {
		b.Logger().Warn("rollback error deleting role binding", "isClusterRoleBinding", entry.IsCluster, "namespace", entry.Namespace, "name", entry.Name, "err", err)

		if time.Now().After(entry.Expiration) {
			b.Logger().Warn("giving up deleting role binding", "isClusterRoleBinding", entry.IsCluster, "namespace", entry.Namespace, "name", entry.Name)
			return nil
		}
		return err
	}

	return nil
}
