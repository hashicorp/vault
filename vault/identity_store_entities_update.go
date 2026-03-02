// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// EntityBuilder is used to construct or update an identity.Entity.
type EntityBuilder struct {
	store  *IdentityStore
	entity *identity.Entity
	isNew  bool
	err    error
}

// NewEntityBuilder creates a new builder instance.
func NewEntityBuilder(store *IdentityStore) *EntityBuilder {
	return &EntityBuilder{
		store:  store,
		entity: new(identity.Entity),
		isNew:  true, // Assume new until an existing entity is loaded
	}
}

// WithExistingEntity allows you to pass in a preloaded entity.
func (b *EntityBuilder) WithExistingEntity(entity *identity.Entity) *EntityBuilder {
	if b.err != nil || entity == nil {
		return b
	}
	b.entity = entity
	b.isNew = false
	return b
}

// WithID attempts to load an entity by its ID.
func (b *EntityBuilder) WithID(id string) *EntityBuilder {
	if b.err != nil || id == "" {
		return b
	}

	entity, err := b.store.MemDBEntityByID(id, true)
	if err != nil {
		b.err = err
		return b
	}
	if entity == nil {
		b.err = fmt.Errorf("entity not found from id: %s", id)
		return b
	}

	b.entity = entity
	b.isNew = false
	return b
}

// WithExternalID handles logic related to the external_id.
func (b *EntityBuilder) WithExternalID(ctx context.Context, externalID string) *EntityBuilder {
	if b.err != nil || externalID == "" {
		return b
	}

	entityByExternalID, err := b.store.MemDBEntityByExternalID(ctx, externalID, true)
	if err != nil {
		b.err = err
		return b
	}

	if entityByExternalID != nil {
		// An entity with this external ID already exists, so we'll update it.
		b.entity = entityByExternalID
		b.isNew = false
	} else {
		// No entity found, so we're just setting the external ID on the current one.
		b.entity.ExternalID = externalID
	}

	return b
}

// WithName handles the complex logic for finding an entity by name and checking for conflicts.
func (b *EntityBuilder) WithName(ctx context.Context, name string) *EntityBuilder {
	if b.err != nil || name == "" {
		return b
	}

	entityByName, err := b.store.MemDBEntityByName(ctx, name, true)
	if err != nil {
		b.err = err
		return b
	}

	switch {
	case entityByName == nil:
		// Not found, safe to use this name.
	case b.isNew:
		// We haven't loaded an entity yet, but one with this name exists. Let's update it.
		b.entity = entityByName
		b.isNew = false
	case b.entity.ID == entityByName.ID:
		// The loaded entity and the one found by name are the same. No-op.
	default:
		// A different entity already has this name, which is a conflict.
		b.err = fmt.Errorf("entity name '%s' is already in use", name)
		return b
	}

	b.entity.Name = name
	return b
}

// WithPolicies sets the policies for the entity.
func (b *EntityBuilder) WithPolicies(policies []string) *EntityBuilder {
	if b.err != nil {
		return b
	}
	if strutil.StrListContainsCaseInsensitive(policies, "root") {
		b.err = fmt.Errorf("policies cannot contain root")
		return b
	}
	b.entity.Policies = strutil.RemoveDuplicates(policies, false)
	return b
}

// WithDisabled sets the disabled status of the entity.
func (b *EntityBuilder) WithDisabled(disabled bool) *EntityBuilder {
	if b.err != nil {
		return b
	}
	b.entity.Disabled = disabled
	return b
}

// WithMetadata sets the metadata for the entity.
// The original entity's value for duplicate_of_canonical_id will be preserved.
func (b *EntityBuilder) WithMetadata(metadata map[string]string) *EntityBuilder {
	if b.err != nil {
		return b
	}
	if value, ok := b.entity.Metadata[duplicateCanonicalIDMetadataKey]; ok {
		metadata[duplicateCanonicalIDMetadataKey] = value
	}
	b.entity.Metadata = metadata
	return b
}

// WithSCIMClientID sets the SCIM client ID.
func (b *EntityBuilder) WithSCIMClientID(scimClientID string) *EntityBuilder {
	if b.err != nil {
		return b
	}
	b.entity.ScimClientID = scimClientID
	return b
}

// Build finalizes the entity creation/update process.
func (b *EntityBuilder) Build(ctx context.Context) (*logical.Response, error) {
	// If any previous step set an error, return it immediately.
	if b.err != nil {
		return logical.ErrorResponse(b.err.Error()), nil
	}

	// Sanitize and persist the entity
	if err := b.store.sanitizeEntity(ctx, b.entity); err != nil {
		return nil, err
	}
	if err := b.store.upsertEntity(ctx, b.entity, nil, true); err != nil {
		return nil, err
	}

	// If this was an update to an existing entity, return 204 No Content
	if !b.isNew {
		return &logical.Response{}, nil
	}

	// This was a new entity, so prepare the response data
	respData := map[string]interface{}{
		"id":   b.entity.ID,
		"name": b.entity.Name,
	}
	var aliasIDs []string
	for _, alias := range b.entity.Aliases {
		aliasIDs = append(aliasIDs, alias.ID)
	}
	respData["aliases"] = aliasIDs

	return &logical.Response{
		Data: respData,
	}, nil
}

// FromFieldData is a convenience method to populate the builder from a FieldData object.
func (b *EntityBuilder) FromFieldData(ctx context.Context, d *framework.FieldData) *EntityBuilder {
	if b.err != nil {
		return b
	}

	if id, ok := d.GetOk("id"); ok {
		b.WithID(id.(string))
	}
	if externalID, ok := d.GetOk("external_id"); ok {
		b.WithExternalID(ctx, externalID.(string))
	}
	if name, ok := d.GetOk("name"); ok {
		b.WithName(ctx, name.(string))
	}
	if scimClientID, ok := d.GetOk("scim_client_id"); ok {
		b.WithSCIMClientID(scimClientID.(string))
	}
	if policies, ok := d.GetOk("policies"); ok {
		b.WithPolicies(policies.([]string))
	}
	if disabled, ok := d.GetOk("disabled"); ok {
		b.WithDisabled(disabled.(bool))
	}
	if metadata, ok, err := d.GetOkErr("metadata"); err != nil {
		b.err = fmt.Errorf("failed to parse metadata: %v", err)
	} else if ok {
		b.WithMetadata(metadata.(map[string]string))
	}

	return b
}

// Upsert finalizes the entity update process by persisting it to storage.
func (b *EntityBuilder) Upsert(ctx context.Context) (*identity.Entity, error) {
	if b.err != nil {
		return nil, b.err
	}

	// Sanitize and persist the entity
	if err := b.store.sanitizeEntity(ctx, b.entity); err != nil {
		return nil, err
	}
	if err := b.store.upsertEntity(ctx, b.entity, nil, true); err != nil {
		return nil, err
	}

	return b.entity, nil
}

func (i *IdentityStore) EntityUpdateCommon(ctx context.Context, d *framework.FieldData) (*logical.Response, error) {
	return NewEntityBuilder(i).
		FromFieldData(ctx, d).
		Build(ctx)
}
