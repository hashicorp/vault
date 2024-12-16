// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/identity"
)

var errDuplicateIdentityName = errors.New("duplicate identity name")

// ConflictResolver defines the interface for resolving conflicts between
// entities, groups, and aliases. All methods should implement a check for
// existing=nil. This is an intentional design choice to allow the caller to
// search for extra information if necessary.
type ConflictResolver interface {
	ResolveEntities(ctx context.Context, existing, duplicate *identity.Entity) error
	ResolveGroups(ctx context.Context, existing, duplicate *identity.Group) error
	ResolveAliases(ctx context.Context, parent *identity.Entity, existing, duplicate *identity.Alias) error
}

// The errorResolver is a ConflictResolver that logs a warning message when a
// pre-existing Identity artifact is found with the same factors as a new one.
type errorResolver struct {
	logger log.Logger
}

// ResolveEntities logs a warning message when a pre-existing Entity is found
// and returns a duplicate name error, which should be handled by the caller by
// putting the system in case-sensitive mode.
func (r *errorResolver) ResolveEntities(ctx context.Context, existing, duplicate *identity.Entity) error {
	if existing == nil {
		return nil
	}

	r.logger.Warn(errDuplicateIdentityName.Error(),
		"entity_name", duplicate.Name,
		"duplicate_of_name", existing.Name,
		"duplicate_of_id", existing.ID,
		"action", "merge the duplicate entities into one")

	return errDuplicateIdentityName
}

// ResolveGroups logs a warning message when a pre-existing Group is found and
// returns a duplicate name error, which should be handled by the caller by
// putting the system in case-sensitive mode.
func (r *errorResolver) ResolveGroups(ctx context.Context, existing, duplicate *identity.Group) error {
	if existing == nil {
		return nil
	}

	r.logger.Warn(errDuplicateIdentityName.Error(),
		"group_name", duplicate.Name,
		"duplicate_of_name", existing.Name,
		"duplicate_of_id", existing.ID,
		"action", "merge the contents of duplicated groups into one and delete the other")

	return errDuplicateIdentityName
}

// ResolveAliases logs a warning message when a pre-existing Alias is found and
// returns a duplicate name error, which should be handled by the caller by
// putting the system in case-sensitive mode.
func (r *errorResolver) ResolveAliases(ctx context.Context, parent *identity.Entity, existing, duplicate *identity.Alias) error {
	if existing == nil {
		return nil
	}

	r.logger.Warn(errDuplicateIdentityName.Error(),
		"alias_name", duplicate.Name,
		"mount_accessor", duplicate.MountAccessor,
		"local", duplicate.Local,
		"entity_name", parent.Name,
		"alias_canonical_id", duplicate.CanonicalID,
		"duplicate_of_name", existing.Name,
		"duplicate_of_canonical_id", existing.CanonicalID,
		"action", "merge the canonical entity IDs into one")

	return errDuplicateIdentityName
}
