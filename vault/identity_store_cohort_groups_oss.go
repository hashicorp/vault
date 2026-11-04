// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/vault/helper/identity"
)

// checkCohortGroupUpdate allows every update, since cohorts are enterprise-only.
func (i *IdentityStore) checkCohortGroupUpdate(context.Context, *identity.Group, []string, []string) error {
	return nil
}

// checkCohortGroupDelete allows every delete, since cohorts are enterprise-only.
func (i *IdentityStore) checkCohortGroupDelete(context.Context, *identity.Group) error {
	return nil
}

func (c *Core) agentRegistryCohortExists(context.Context, string, *identity.AgentRegistryCohort) (bool, error) {
	return false, nil
}
