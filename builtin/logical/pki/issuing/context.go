// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package issuing

import "context"

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
