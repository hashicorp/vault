// +build !enterprise

package vault

import (
	"context"
	"path"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/logical"
)

func addPathCheckers(*Core, *MountEntry, logical.Backend, string)             {}
func removePathCheckers(*Core, *MountEntry, string)                           {}
func addAuditPathChecker(*Core, *MountEntry, *BarrierView, string)            {}
func removeAuditPathChecker(*Core, *MountEntry)                               {}
func addFilterablePath(*Core, string)                                         {}
func preprocessMount(*Core, *MountEntry, *BarrierView) (bool, error)          { return false, nil }
func clearIgnoredPaths(context.Context, *Core, logical.Backend, string) error { return nil }

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
