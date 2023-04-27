// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package vault

import (
	"context"
	"path"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

func addPathCheckers(c *Core, entry *MountEntry, backend logical.Backend, viewPath string) {
	c.addBackendWriteForwardedPaths(backend, viewPath)
}

func removePathCheckers(c *Core, entry *MountEntry, viewPath string) {
	c.writeForwardedPaths.RemovePathPrefix(viewPath)
}

func addAuditPathChecker(*Core, *MountEntry, *BarrierView, string)            {}
func removeAuditPathChecker(*Core, *MountEntry)                               {}
func addFilterablePath(*Core, string)                                         {}
func preprocessMount(*Core, *MountEntry, *BarrierView) (bool, error)          { return false, nil }
func clearIgnoredPaths(context.Context, *Core, logical.Backend, string) error { return nil }
func addLicenseCallback(*Core, logical.Backend)                               {}
func runFilteredPathsEvaluation(context.Context, *Core) error                 { return nil }

// ViewPath returns storage prefix for the view
func (e *MountEntry) ViewPath() string {
	switch e.Type {
	case systemMountType:
		return systemBarrierPrefix
	case "token":
		return path.Join(systemBarrierPrefix, tokenSubPath) + "/"
	}

	switch e.Table {
	case mountTableType:
		return backendBarrierPrefix + e.UUID + "/"
	case credentialTableType:
		return credentialBarrierPrefix + e.UUID + "/"
	case auditTableType:
		return auditBarrierPrefix + e.UUID + "/"
	}

	panic("invalid mount entry")
}

func verifyNamespace(*Core, *namespace.Namespace, *MountEntry) error { return nil }

// mountEntrySysView creates a logical.SystemView from global and
// mount-specific entries; because this should be called when setting
// up a mountEntry, it doesn't check to ensure that me is not nil
func (c *Core) mountEntrySysView(entry *MountEntry) extendedSystemView {
	return extendedSystemViewImpl{
		dynamicSystemView{
			core:        c,
			mountEntry:  entry,
			perfStandby: c.perfStandby,
		},
	}
}
