// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build testonly

package vault

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"unicode"

	"github.com/golang/protobuf/ptypes"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/storagepacker"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/protobuf/types/known/anypb"
)

// entityTestonlyPaths returns a list of testonly API endpoints supported to
// operate on entities in a way that is not supported by production Vault. These
// all generate duplicate identity resources IN STORAGE. MemDB won't reflect
// them until Vault has been sealed and unsealed again!
//
// Use of these endpoints is a bit nuanced as they are low level and do almost
// no validation. By design, they are allowing you to write invalid state into
// storage because that is what is needed to replicate some customer scenarios
// caused by historical bugs. Bear the following non-obvious things in mind if
// you use them.
//
//   - Very little validation is done. You can create state that in invalid in
//     ways that Vault, even with it's bugs, has never been able to create.
//   - These write the duplicates directly to storage without checking contents.
//     So if you call the same endpoint with the same name multiple times you
//     will end up with even more duplicates of the same name.
//   - Because they write direct to storage, they DON'T update MemDB so regular
//     API calls won't see the created resources until you seal and unseal.
func entityTestonlyPaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "duplicate/entity-aliases",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "entity-aliases",
				OperationVerb:   "create-duplicates",
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the entities to create",
				},
				"namespace_id": {
					Type:        framework.TypeString,
					Description: "NamespaceID of the entities to create",
				},
				"different_case": {
					Type:        framework.TypeBool,
					Description: "Create entities with different case variations",
				},
				"mount_accessor": {
					Type:        framework.TypeString,
					Description: "Mount accessor ID for the alias",
				},
				"metadata": {
					Type:        framework.TypeKVPairs,
					Description: "Metadata",
				},
				"count": {
					Type:        framework.TypeInt,
					Description: "Number of entity aliases to create",
				},
				"local": {
					Type:        framework.TypeBool,
					Description: "Local alias toggle",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                  i.createDuplicateEntityAliases(),
					ForwardPerformanceStandby: true,
					// Duplicate entity alias injector now forwards
					// CreateEntity calls when flags.Local is true.
					ForwardPerformanceSecondary: false,
				},
			},
		},
		{
			Pattern: "duplicate/local-entity-alias",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "entity-alias",
				OperationVerb:   "create-duplicates",
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the entities to create",
				},
				"namespace_id": {
					Type:        framework.TypeString,
					Description: "NamespaceID of the entities to create",
				},
				"canonical_id": {
					Type:        framework.TypeString,
					Description: "The canonical entity ID to attach the local alias to",
				},
				"mount_accessor": {
					Type:        framework.TypeString,
					Description: "Mount accessor ID for the alias",
				},
				"metadata": {
					Type:        framework.TypeKVPairs,
					Description: "Metadata",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                    i.createDuplicateLocalEntityAlias(),
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: false, // Allow this on a perf secondary.
				},
			},
		},
		{
			Pattern: "duplicate/entities",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "entities",
				OperationVerb:   "create-duplicates",
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the entities to create",
				},
				"namespace_id": {
					Type:        framework.TypeString,
					Description: "NamespaceID of the entities to create",
				},
				"different_case": {
					Type:        framework.TypeBool,
					Description: "Create entities with different case variations",
				},
				"metadata": {
					Type: framework.TypeKVPairs,
					Description: `Metadata to be associated with the entity.
In CLI, this parameter can be repeated multiple times, and it all gets merged together.
For example:
vault <command> <path> metadata=key1=value1 metadata=key2=value2
					`,
				},
				"policies": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Policies to be tied to the entity.",
				},
				"count": {
					Type:        framework.TypeInt,
					Description: "Number of entities to create",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                  i.createDuplicateEntities(),
					ForwardPerformanceStandby: true,
					// Writing global (non-local) state should be replicated.
					ForwardPerformanceSecondary: true,
				},
			},
		},
		{
			Pattern: "duplicate/groups",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "groups",
				OperationVerb:   "create-duplicates",
			},
			Fields: map[string]*framework.FieldSchema{
				"id": {
					Type:        framework.TypeString,
					Description: "ID of the group. If set, updates the corresponding existing group.",
				},
				"type": {
					Type:        framework.TypeString,
					Description: "Type of the group, 'internal' or 'external'. Defaults to 'internal'",
				},
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the group.",
				},
				"namespace_id": {
					Type:        framework.TypeString,
					Description: "NamespaceID of the entities to create",
				},
				"different_case": {
					Type:        framework.TypeBool,
					Description: "Create entities with different case variations",
				},
				"metadata": {
					Type: framework.TypeKVPairs,
					Description: `Metadata to be associated with the group.
In CLI, this parameter can be repeated multiple times, and it all gets merged together.
For example:
vault <command> <path> metadata=key1=value1 metadata=key2=value2
					`,
				},
				"policies": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Policies to be tied to the group.",
				},
				"member_group_ids": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Group IDs to be assigned as group members.",
				},
				"member_entity_ids": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Entity IDs to be assigned as group members.",
				},
				"count": {
					Type:        framework.TypeInt,
					Description: "Number of groups to create",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                  i.createDuplicateGroups(),
					ForwardPerformanceStandby: true,
					// Writing global (non-local) state should be replicated.
					ForwardPerformanceSecondary: true,
				},
			},
		},
		{
			Pattern: "entity/from-storage/?$",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "entity",
				OperationVerb:   "list",
				OperationSuffix: "from-storage",
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback:                  i.listEntitiesFromStorage(),
					ForwardPerformanceStandby: true,
					// Allow reading local cluster state
					ForwardPerformanceSecondary: false,
				},
			},
		},
		{
			Pattern: "group/from-storage/?$",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "group",
				OperationVerb:   "list",
				OperationSuffix: "from-storage",
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback:                  i.listGroupsFromStorage(),
					ForwardPerformanceStandby: true,
					// Allow reading local cluster state
					ForwardPerformanceSecondary: false,
				},
			},
		},
	}
}

type CommonDuplicateFlags struct {
	Name          string            `json:"name"`
	NamespaceID   string            `json:"namespace_id"`
	DifferentCase bool              `json:"different_case"`
	Metadata      map[string]string `json:"metadata"`
}

type CommonAliasFlags struct {
	MountAccessor string `json:"mount_accessor"`
	CanonicalID   string `json:"canonical_id"`
}

type DuplicateEntityFlags struct {
	CommonDuplicateFlags
	Policies []string `json:"policies"`
	Count    int      `json:"count"`
}

type DuplicateGroupFlags struct {
	CommonDuplicateFlags
	Type            string   `json:"type"`
	Policies        []string `json:"policies"`
	MemberGroupIDs  []string `json:"member_group_ids"`
	MemberEntityIDs []string `json:"member_entity_ids"`
	Count           int      `json:"count"`
}

type DuplicateEntityAliasFlags struct {
	CommonDuplicateFlags
	CommonAliasFlags
	Count int  `json:"count"`
	Local bool `json:"local"`
}

type DuplicateGroupAliasFlags struct {
	CommonAliasFlags
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (i *IdentityStore) createDuplicateEntities() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		metadata, ok := data.GetOk("metadata")
		if !ok {
			metadata = make(map[string]string)
		}
		flags := DuplicateEntityFlags{
			CommonDuplicateFlags: CommonDuplicateFlags{
				Name:          data.Get("name").(string),
				NamespaceID:   data.Get("namespace_id").(string),
				DifferentCase: data.Get("different_case").(bool),
				Metadata:      metadata.(map[string]string),
			},
			Policies: data.Get("policies").([]string),
			Count:    data.Get("count").(int),
		}

		if flags.Count < 1 {
			flags.Count = 2
		}

		ids, err := i.CreateDuplicateEntitiesInStorage(ctx, flags)
		if err != nil {
			i.logger.Error("error creating duplicate entities", "error", err)
			return logical.ErrorResponse("error creating duplicate entities"), err
		}

		return &logical.Response{
			Data: map[string]interface{}{
				"entity_ids": ids,
			},
		}, nil
	}
}

func (i *IdentityStore) createDuplicateEntityAliases() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		metadata, ok := data.GetOk("metadata")
		if !ok {
			metadata = make(map[string]string)
		}

		flags := DuplicateEntityAliasFlags{
			CommonDuplicateFlags: CommonDuplicateFlags{
				Name:          data.Get("name").(string),
				NamespaceID:   data.Get("namespace_id").(string),
				DifferentCase: data.Get("different_case").(bool),
				Metadata:      metadata.(map[string]string),
			},
			CommonAliasFlags: CommonAliasFlags{
				MountAccessor: data.Get("mount_accessor").(string),
			},
			Count: data.Get("count").(int),
			Local: data.Get("local").(bool),
		}

		if flags.Count < 1 {
			flags.Count = 2
		}

		ids, err := i.CreateDuplicateEntityAliasesInStorage(ctx, flags)
		if err != nil {
			i.logger.Error("error creating duplicate entity aliases", "error", err)
			return logical.ErrorResponse("error creating duplicate entity aliases"), err
		}

		return &logical.Response{
			Data: map[string]interface{}{
				"entity_ids": ids,
				"local":      flags.Local,
			},
		}, nil
	}
}

func (i *IdentityStore) createDuplicateLocalEntityAlias() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		metadata, ok := data.GetOk("metadata")
		if !ok {
			metadata = make(map[string]string)
		}

		flags := DuplicateEntityAliasFlags{
			CommonDuplicateFlags: CommonDuplicateFlags{
				Name:        data.Get("name").(string),
				NamespaceID: data.Get("namespace_id").(string),
				Metadata:    metadata.(map[string]string),
			},
			CommonAliasFlags: CommonAliasFlags{
				MountAccessor: data.Get("mount_accessor").(string),
				CanonicalID:   data.Get("canonical_id").(string),
			},
		}

		if flags.Name == "" {
			return logical.ErrorResponse("name is required"), nil
		}
		if flags.CanonicalID == "" {
			return logical.ErrorResponse("canonical_id is required"), nil
		}
		if flags.MountAccessor == "" {
			return logical.ErrorResponse("mount_accessor is required"), nil
		}

		ids, err := i.CreateDuplicateLocalEntityAliasInStorage(ctx, flags)
		if err != nil {
			i.logger.Error("error creating duplicate local alias", "error", err)
			return logical.ErrorResponse("error creating duplicate local alias"), err
		}

		return &logical.Response{
			Data: map[string]interface{}{
				"alias_ids": ids,
			},
		}, nil
	}
}

func (i *IdentityStore) createDuplicateGroups() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		metadata, ok := data.GetOk("metadata")
		if !ok {
			metadata = make(map[string]string)
		}

		flags := DuplicateGroupFlags{
			CommonDuplicateFlags: CommonDuplicateFlags{
				Name:          data.Get("name").(string),
				NamespaceID:   data.Get("namespace_id").(string),
				DifferentCase: data.Get("different_case").(bool),
				Metadata:      metadata.(map[string]string),
			},
			Type:            data.Get("type").(string),
			Policies:        data.Get("policies").([]string),
			MemberGroupIDs:  data.Get("member_group_ids").([]string),
			MemberEntityIDs: data.Get("member_entity_ids").([]string),
			Count:           data.Get("count").(int),
		}

		if flags.Count < 1 {
			flags.Count = 2
		}

		ids, err := i.CreateDuplicateGroupsInStorage(ctx, flags)
		if err != nil {
			i.logger.Error("error creating duplicate entities", "error", err)
			return logical.ErrorResponse("error creating duplicate entities"), err
		}

		return &logical.Response{
			Data: map[string]interface{}{
				"group_ids": ids,
			},
		}, nil
	}
}

func (i *IdentityStore) CreateDuplicateGroupsInStorage(ctx context.Context, flags DuplicateGroupFlags) ([]string, error) {
	var groupIDs []string
	if flags.NamespaceID == "" {
		flags.NamespaceID = namespace.RootNamespaceID
	}
	for d := 0; d < flags.Count; d++ {
		groupID, err := uuid.GenerateUUID()
		if err != nil {
			return nil, err
		}
		groupIDs = append(groupIDs, groupID)

		// Alias name is either exact match or different case
		groupName := flags.Name
		if flags.DifferentCase {
			groupName = randomCase(flags.Name)
		}

		g := &identity.Group{
			ID:              groupID,
			Name:            groupName,
			Policies:        flags.Policies,
			MemberEntityIDs: flags.MemberEntityIDs,
			ParentGroupIDs:  flags.MemberGroupIDs,
			Type:            flags.Type,
			NamespaceID:     flags.NamespaceID,
			BucketKey:       i.groupPacker.BucketKey(groupID),
		}

		group, err := ptypes.MarshalAny(g)
		if err != nil {
			return nil, err
		}
		item := &storagepacker.Item{
			ID:      g.ID,
			Message: group,
		}
		if err = i.groupPacker.PutItem(ctx, item); err != nil {
			return nil, err
		}
	}

	return groupIDs, nil
}

// CreateDuplicateEntityAliasesInStorage creates n entities with a duplicate
// alias in storage This should only be used in testing. This method can only
// create non-local aliases. Local aliases are stored differently.
//
// Pass in mount type and accessor to create the entities
func (i *IdentityStore) CreateDuplicateEntityAliasesInStorage(ctx context.Context, flags DuplicateEntityAliasFlags) ([]string, error) {
	var entityIDs []string
	if flags.NamespaceID == "" {
		flags.NamespaceID = namespace.RootNamespaceID
	}
	for d := 0; d < flags.Count; d++ {
		aliasID, err := uuid.GenerateUUID()
		if err != nil {
			return nil, err
		}

		// Alias name is either exact match or different case
		dupAliasName := flags.Name
		if flags.DifferentCase {
			dupAliasName = randomCase(flags.Name)
		}

		// In real life alias dupes are due to races where they were auto-created
		// along with the entities they point to. When we auto create entities like
		// this they get random names so never collide with each other directly so
		// don't create entities with duplicate names for this case as it doesn't
		// match what customers who get in this state see. Instead use the same code
		// path as CreateOfFetchEntity.
		e := new(identity.Entity)
		e.NamespaceID = flags.NamespaceID
		err = i.sanitizeEntity(ctx, e)
		if err != nil {
			return nil, fmt.Errorf("error sanitizing entity: %w", err)
		}
		entityIDs = append(entityIDs, e.ID)

		a := &identity.Alias{
			ID:            aliasID,
			NamespaceID:   flags.NamespaceID,
			CanonicalID:   e.ID,
			MountAccessor: flags.MountAccessor,
			Name:          dupAliasName,
			Local:         flags.Local,
		}

		persistEntity := func(ent *identity.Entity) error {
			entity, err := ptypes.MarshalAny(ent)
			if err != nil {
				return fmt.Errorf("error marhsaling entity: %w", err)
			}
			item := &storagepacker.Item{
				ID:      ent.ID,
				Message: entity,
			}
			if err = i.entityPacker.PutItem(ctx, item); err != nil {
				return err
			}

			return nil
		}

		if flags.Local {
			// Check to see if the entity creation should be forwarded.
			i.logger.Trace("forwarding entity creation for local alias cache")
			e, err := i.entityCreator.CreateEntity(ctx)
			if err != nil {
				return nil, err
			}

			localAliases, err := i.parseLocalAliases(e.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to parse local aliases from entity: %w", err)
			}
			if localAliases == nil {
				localAliases = &identity.LocalAliases{}
			}

			// Don't check if this is a duplicate, since we're allowing the developer to
			// create duplicates here.
			localAliases.Aliases = append(localAliases.Aliases, a)

			marshaledAliases, err := anypb.New(localAliases)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal local aliases: %w", err)
			}
			item := &storagepacker.Item{
				ID:      e.ID,
				Message: marshaledAliases,
			}
			if err := i.localAliasPacker.PutItem(ctx, item); err != nil {
				return nil, fmt.Errorf("failed to put item in local alias packer: %w", err)
			}
		} else {
			e.UpsertAlias(a)
			err := persistEntity(e)
			if err != nil {
				return nil, err
			}

		}
	}

	return entityIDs, nil
}

// CreateDuplicateLocalEntityAliasInStorage creates a single local entity alias
// directly in storage. This should only be used in testing. This method can
// only create local aliases and assumes that the entity is already created
// separately and it's ID passed as CanonicalID. No validation of the mounts or
// entity is done so if you need these to be realistic the caller must ensure
// the entity and mount exist and that the mount is a local auth method of the
// right type.
//
// Pass in mount type and accessor to create the entities
func (i *IdentityStore) CreateDuplicateLocalEntityAliasInStorage(ctx context.Context, flags DuplicateEntityAliasFlags) ([]string, error) {
	var aliasIDs []string
	if flags.NamespaceID == "" {
		flags.NamespaceID = namespace.RootNamespaceID
	}

	aliasID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	aliasIDs = append(aliasIDs, aliasID)

	a := &identity.Alias{
		ID:            aliasID,
		NamespaceID:   flags.NamespaceID,
		CanonicalID:   flags.CanonicalID,
		MountAccessor: flags.MountAccessor,
		Name:          flags.Name,
		Local:         true,
	}

	localAliases, err := i.parseLocalAliases(flags.CommonAliasFlags.CanonicalID)
	if err != nil {
		return nil, err
	}
	if localAliases == nil {
		localAliases = &identity.LocalAliases{}
	}

	// Don't check if this is a duplicate, since we're allowing the developer to
	// create duplicates here.
	localAliases.Aliases = append(localAliases.Aliases, a)

	marshaledAliases, err := anypb.New(localAliases)
	if err != nil {
		return nil, err
	}
	if err := i.localAliasPacker.PutItem(ctx, &storagepacker.Item{
		ID:      flags.CommonAliasFlags.CanonicalID,
		Message: marshaledAliases,
	}); err != nil {
		return nil, err
	}

	return aliasIDs, nil
}

func (i *IdentityStore) CreateDuplicateEntitiesInStorage(ctx context.Context, flags DuplicateEntityFlags) ([]string, error) {
	var entityIDs []string
	for d := 0; d < flags.Count; d++ {
		entityID, err := uuid.GenerateUUID()
		if err != nil {
			return nil, err
		}
		entityIDs = append(entityIDs, entityID)

		dupName := flags.Name
		if flags.DifferentCase {
			dupName = randomCase(flags.Name)
		}

		e := &identity.Entity{
			ID:          entityID,
			Name:        dupName,
			NamespaceID: flags.NamespaceID,
			BucketKey:   i.entityPacker.BucketKey(entityID),
		}

		entity, err := ptypes.MarshalAny(e)
		if err != nil {
			return nil, err
		}
		item := &storagepacker.Item{
			ID:      e.ID,
			Message: entity,
		}
		if err = i.entityPacker.PutItem(ctx, item); err != nil {
			return nil, err
		}
	}

	return entityIDs, nil
}

func randomCase(s string) string {
	return strings.Map(func(r rune) rune {
		if rand.Intn(2) == 0 {
			return unicode.ToUpper(r)
		}
		return unicode.ToLower(r)
	}, s)
}

func (i *IdentityStore) ListEntitiesFromStorage(ctx context.Context) ([]*identity.Entity, error) {
	// Get Existing Buckets
	existing, err := i.entityPacker.View().List(ctx, storagepacker.StoragePackerBucketsPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to scan for entity buckets: %w", err)
	}

	workerCount := 64
	entities := make([]*identity.Entity, 0)

	// Make channels for worker pool
	broker := make(chan string)
	quit := make(chan bool)

	errs := make(chan error, (len(existing)))
	result := make(chan *storagepacker.Bucket, len(existing))

	wg := &sync.WaitGroup{}

	// Stand up workers
	for j := 0; j < workerCount; j++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case key, ok := <-broker:
					if !ok {
						return
					}

					bucket, err := i.entityPacker.GetBucket(ctx, storagepacker.StoragePackerBucketsPrefix+key)
					if err != nil {
						errs <- err
						continue
					}

					result <- bucket

				case <-quit:
					return
				}
			}
		}()
	}

	// Distribute the collected keys to the workers in a go routine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for j, key := range existing {
			if j%500 == 0 {
				i.logger.Debug("entities loading", "progress", j)
			}

			select {
			case <-quit:
				return

			default:
				broker <- key
			}
		}

		// Close the broker, causing worker routines to exit
		close(broker)
	}()

	// Restore each key by pulling from the result chan
LOOP:
	for j := 0; j < len(existing); j++ {
		select {
		case err = <-errs:
			// Close all go routines
			close(quit)
			break LOOP

		case bucket := <-result:
			// If there is no entry, nothing to restore
			if bucket == nil {
				continue
			}

			for _, item := range bucket.Items {
				entity, err := i.parseEntityFromBucketItem(ctx, item)
				if err != nil {
					return nil, err
				}
				if entity == nil {
					continue
				}

				// Load local aliases for entity
				localAliases, err := i.parseLocalAliases(entity.ID)
				if err != nil {
					return nil, err
				}
				if localAliases != nil {
					entity.Aliases = append(entity.Aliases, localAliases.Aliases...)
				}

				entities = append(entities, entity)
			}
		}
	}

	// Let all go routines finish
	wg.Wait()
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (i *IdentityStore) listEntitiesFromStorage() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		entities, err := i.ListEntitiesFromStorage(ctx)
		if err != nil {
			i.logger.Error("error listing entities", "error", err)
			return logical.ErrorResponse("error listing entities"), err
		}
		resp := &logical.Response{
			Data: map[string]interface{}{
				"entities": entities,
			},
		}
		return resp, nil
	}
}

func (i *IdentityStore) ListGroupsFromStorage(ctx context.Context) ([]*identity.Group, error) {
	existing, err := i.groupPacker.View().List(ctx, groupBucketsPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to scan for groups: %w", err)
	}

	groups := make([]*identity.Group, 0)

	for _, key := range existing {
		bucket, err := i.groupPacker.GetBucket(ctx, groupBucketsPrefix+key)
		if err != nil {
			return nil, err
		}

		if bucket == nil {
			continue
		}

		for _, item := range bucket.Items {
			group, err := i.parseGroupFromBucketItem(item)
			if err != nil {
				return nil, err
			}
			if group == nil {
				continue
			}
			groups = append(groups, group)
		}
	}
	return groups, nil
}

func (i *IdentityStore) listGroupsFromStorage() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		groups, err := i.ListGroupsFromStorage(ctx)
		if err != nil {
			i.logger.Error("error listing groups", "error", err)
			return logical.ErrorResponse("error listing groups"), err
		}
		resp := &logical.Response{
			Data: map[string]interface{}{
				"groups": groups,
			},
		}
		return resp, nil
	}
}
