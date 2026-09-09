// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/vault/helper/identity"
)

// SCIM client storage prefix
const scimClientStoragePrefix = "scim/client/"

func scimClientSchema(_ bool) *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: scimClientsTable,
		Indexes: map[string]*memdb.IndexSchema{
			"id": {
				Name:   "id",
				Unique: true,
				Indexer: &memdb.StringFieldIndex{
					Field: "ClientID",
				},
			},
			"client_name": {
				Name:   "client_name",
				Unique: true,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field: "NamespaceID",
						},
						&memdb.StringFieldIndex{
							Field: "ClientName",
						},
					},
				},
			},
			"namespace_id": {
				Name: "namespace_id",
				Indexer: &memdb.StringFieldIndex{
					Field: "NamespaceID",
				},
			},
			"access_grant_principal": {
				Name:   "access_grant_principal",
				Unique: true,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field: "NamespaceID",
						},
						&memdb.StringFieldIndex{
							Field: "AccessGrantPrincipal",
						},
					},
				},
			},
		},
	}
}

type scimClientRequest struct{}

func addSCIMClientIDToContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, scimClientRequest{}, id)
}

func scimClientIDFromContext(ctx context.Context) string {
	val := ctx.Value(scimClientRequest{})
	if val == nil {
		return ""
	}
	return val.(string)
}

// scimResourceCheck performs a number of checks on the resource to ensure that the operation is allowed.
// A new resource must only set a SCIM client ID if the request came through SCIM via the same client.
// An updated resource must:
// - not change the SCIM client ID
// - not be modified by a different SCIM client than the one that owns it
// - only update SCIM managed fields if the request came via the same SCIM client
func (i *IdentityStore) scimResourceCheck(ctx context.Context, resource scimManaged, originalSCIMID string, isCreate bool, modifiedFields []string) error {
	reqSCIMClientID := scimClientIDFromContext(ctx)
	resourceSCIMClientID := resource.SCIMClientID()

	switch isCreate {
	case true:
		// The request must have come via a SCIM API in order to set
		// the SCIM client ID
		if reqSCIMClientID == "" && resourceSCIMClientID != "" {
			return errors.New("cannot set scim_client_id")
		}
		if reqSCIMClientID != "" && resourceSCIMClientID == "" {
			// this shouldn't ever happen
			return errors.New("cannot create a resource via SCIM without a SCIM client ID")
		}
		if reqSCIMClientID != resourceSCIMClientID {
			// this also shouldn't ever happen
			return errors.New("cannot create resource via SCIM with a different SCIM client ID")
		}

	default:
		// if the resource is being updated, the SCIM client ID
		// cannot be modified
		if originalSCIMID != resourceSCIMClientID {
			return errors.New("cannot update scim_client_id")
		}
		if originalSCIMID != "" && reqSCIMClientID != "" && originalSCIMID != reqSCIMClientID {
			return errors.New("cannot update resource via SCIM with a different SCIM client ID")
		}
		if err := i.scimFieldGuard(ctx, resource, modifiedFields); err != nil {
			return err
		}
	}
	return nil
}

// scimFieldGuard checks whether an API request is allowed to modify the given
// fields on a SCIM-managed resource. If the resource is not SCIM-managed
// (scimClientID is empty), all modifications are allowed. If the request
// comes from the owning SCIM client, all modifications are allowed. Otherwise,
// only non-SCIM-managed fields may be modified.
func (i *IdentityStore) scimFieldGuard(ctx context.Context, resource scimManaged, modifiedFields []string) error {
	// If the resource is not SCIM-managed, no restrictions apply.
	scimClientID := resource.SCIMClientID()
	if scimClientID == "" {
		return nil
	}

	// If the request comes from the owning SCIM client, all modifications are allowed.
	if scimClientIDFromContext(ctx) == scimClientID {
		return nil
	}

	// The request is from the API (not SCIM). Check whether any of the
	// modified fields are SCIM-managed.
	scimManagedFields := resource.SCIMFields()
	managedSet := make(map[string]struct{}, len(scimManagedFields))
	for _, f := range scimManagedFields {
		managedSet[f] = struct{}{}
	}

	var blocked []string
	for _, f := range modifiedFields {
		if _, ok := managedSet[f]; ok {
			blocked = append(blocked, f)
		}
	}

	if len(blocked) > 0 {
		return fmt.Errorf("cannot modify SCIM-managed field(s) %q through the API", strings.Join(blocked, ", "))
	}

	return nil
}

type scimManaged interface {
	// SCIMClientID returns the SCIM client ID of the managed resource
	SCIMClientID() string
	// SCIMFields returns the list of fields that are managed by SCIM and cannot
	// be modified through the API. This list does not need to include `scim_client_id`
	// as that is verified separately
	SCIMFields() []string
}

var (
	_ scimManaged = (*identity.Entity)(nil)
	_ scimManaged = (*identity.Group)(nil)
	_ scimManaged = (*identity.Alias)(nil)
)

// scimClientIDSentinel is the index key used for groups with no SCIM owner
// (ScimClientID == ""). A plain StringFieldIndex returns (false, nil, nil) for
// empty strings, which causes CompoundIndex+AllowMissing to truncate the stored
// key to just the NamespaceID component. The radix-tree prefix iterator then
// returns keys strictly longer than the search prefix, so a prefix scan for
// (nsID, "") finds only owned/other-client groups — the unmanaged groups stored
// at the exact nsID\x00 node are never yielded. The sentinel solves this by
// ensuring every group's key is always fully qualified, making unmanaged groups
// reachable via a normal prefix scan on (nsID, scimClientIDSentinel).
const scimClientIDSentinel = "__unmanaged__"

// scimClientIDIndexer is a memdb.SingleIndexer + memdb.PrefixIndexer for the
// ScimClientID field of identity.Group. It maps "" to scimClientIDSentinel so
// that unmanaged groups receive a full, distinct index key rather than a
// truncated one, and are therefore reachable via a direct prefix scan.
type scimClientIDIndexer struct{}

func (scimClientIDIndexer) FromObject(raw interface{}) (bool, []byte, error) {
	g, ok := raw.(*identity.Group)
	if !ok {
		return false, nil, fmt.Errorf("scimClientIDIndexer: unexpected type %T", raw)
	}
	val := g.ScimClientID
	if val == "" {
		val = scimClientIDSentinel
	}
	return true, []byte(val + "\x00"), nil
}

func (scimClientIDIndexer) FromArgs(args ...interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("scimClientIDIndexer: expected 1 argument, got %d", len(args))
	}
	val, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("scimClientIDIndexer: argument must be a string, got %T", args[0])
	}
	if val == "" {
		val = scimClientIDSentinel
	}
	return []byte(val + "\x00"), nil
}

func (scimClientIDIndexer) PrefixFromArgs(args ...interface{}) ([]byte, error) {
	b, err := scimClientIDIndexer{}.FromArgs(args...)
	if err != nil {
		return nil, err
	}
	return b[:len(b)-1], nil // strip null terminator, leaving a prefix
}
