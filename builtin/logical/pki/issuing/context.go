// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package issuing

import (
	"context"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

// IssuerRoleContext combines in a single struct an issuer and a role that we should
// leverage to issue a certificate along with the
type IssuerRoleContext struct {
	context.Context
	Role   *RoleEntry
	Issuer *IssuerEntry
}

func NewIssuerRoleContext(ctx context.Context, issuer *IssuerEntry, role *RoleEntry) IssuerRoleContext {
	return IssuerRoleContext{
		Context: ctx,
		Role:    role,
		Issuer:  issuer,
	}
}

// MountAttributionFromRequest builds a MountAttribution from a logical.Request and
// the context it was dispatched in. The Count field is left zero — callers set it
// to the per-certificate billable units via WithMountInfo on CertCountIncrementer.
func MountAttributionFromRequest(ctx context.Context, req *logical.Request, backendUUID string) logical.MountAttribution {
	attr := logical.MountAttribution{
		NamespaceID: namespace.RootNamespaceID,
	}
	if req != nil {
		attr.MountAccessor = req.MountAccessor
		attr.MountPath = req.MountPoint
		attr.MountType = req.MountType
		attr.BackendAwareUUID = backendUUID
	}
	if ns, err := namespace.FromContext(ctx); err == nil {
		attr.NamespaceID = ns.ID
		attr.NamespacePath = ns.Path
	}
	return attr
}
