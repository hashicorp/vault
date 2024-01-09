// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/vault/helper/namespace"
)

type LockedUsersResponse struct {
	NamespaceID    string                    `json:"namespace_id" mapstructure:"namespace_id"`
	NamespacePath  string                    `json:"namespace_path" mapstructure:"namespace_path"`
	Counts         int                       `json:"counts" mapstructure:"counts"`
	MountAccessors []*ResponseMountAccessors `json:"mount_accessors" mapstructure:"mount_accessors"`
}

type ResponseMountAccessors struct {
	MountAccessor    string   `json:"mount_accessor" mapstructure:"mount_accessor"`
	Counts           int      `json:"counts" mapstructure:"counts"`
	AliasIdentifiers []string `json:"alias_identifiers" mapstructure:"alias_identifiers"`
}

// unlockUser deletes the entry for locked user from storage and userFailedLoginInfo map
func unlockUser(ctx context.Context, core *Core, mountAccessor string, aliasName string) error {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}

	lockedUserStoragePath := coreLockedUsersPath + ns.ID + "/" + mountAccessor + "/" + aliasName

	// remove entry for locked user from storage
	// if read only error, the error is handled by handleError in logical_system.go
	// this will be forwarded to the active node
	if err := core.barrier.Delete(ctx, lockedUserStoragePath); err != nil {
		return err
	}

	loginUserInfoKey := FailedLoginUser{
		aliasName:     aliasName,
		mountAccessor: mountAccessor,
	}

	// remove entry for locked user from userFailedLoginInfo map and storage
	if err := updateUserFailedLoginInfo(ctx, core, loginUserInfoKey, nil, true); err != nil {
		return err
	}

	return nil
}

// handleLockedUsersQuery reports the locked user metrics by namespace in the decreasing order
// of locked users
func (b *SystemBackend) handleLockedUsersQuery(ctx context.Context, mountAccessor string) (map[string]interface{}, error) {
	// Calculate the namespace response breakdowns of locked users for query namespace and child namespaces (if needed)
	totalCount, byNamespaceResponse, err := b.getLockedUsersResponses(ctx, mountAccessor)
	if err != nil {
		return nil, err
	}

	// Now populate the response based on breakdowns.
	responseData := make(map[string]interface{})
	responseData["by_namespace"] = byNamespaceResponse
	responseData["total"] = totalCount
	return responseData, nil
}

// getLockedUsersResponses returns the locked users
// for a particular mount_accessor if provided in request
// else returns it for the current namespace and all the child namespaces that has locked users
// they are sorted in the decreasing order of locked users count
func (b *SystemBackend) getLockedUsersResponses(ctx context.Context, mountAccessor string) (int, []*LockedUsersResponse, error) {
	lockedUsersResponse := make([]*LockedUsersResponse, 0)
	totalCounts := 0

	queryNS, err := namespace.FromContext(ctx)
	if err != nil {
		return 0, nil, err
	}

	if mountAccessor != "" {
		// get the locked user response for mount_accessor, here for mount_accessor in request
		totalCountForNSID, mountAccessorsResponse, err := b.getMountAccessorsLockedUsers(ctx, []string{mountAccessor + "/"},
			coreLockedUsersPath+queryNS.ID+"/")
		if err != nil {
			return 0, nil, err
		}

		totalCounts += totalCountForNSID
		lockedUsersResponse = append(lockedUsersResponse, &LockedUsersResponse{
			NamespaceID:    queryNS.ID,
			NamespacePath:  queryNS.Path,
			Counts:         totalCountForNSID,
			MountAccessors: mountAccessorsResponse,
		})
		return totalCounts, lockedUsersResponse, nil
	}

	// no mount_accessor is provided in request, get information for current namespace and its child namespaces

	// get all the namespaces of locked users
	nsIDs, err := b.Core.barrier.List(ctx, coreLockedUsersPath)
	if err != nil {
		return 0, nil, err
	}

	// identify if the namespaces must be included in response and get counts
	for _, nsID := range nsIDs {
		nsID = strings.TrimSuffix(nsID, "/")
		ns, err := NamespaceByID(ctx, nsID, b.Core)
		if err != nil {
			return 0, nil, err
		}

		if b.includeNSInLockedUsersResponse(queryNS, ns) {
			var displayPath string
			if ns == nil {
				// deleted namespace
				displayPath = fmt.Sprintf("deleted namespace %q", nsID)
			} else {
				displayPath = ns.Path
			}

			// get mount accessors of locked users for this namespace
			mountAccessors, err := b.Core.barrier.List(ctx, coreLockedUsersPath+nsID+"/")
			if err != nil {
				return 0, nil, err
			}

			// get the locked user response for mount_accessor list
			totalCountForNSID, mountAccessorsResponse, err := b.getMountAccessorsLockedUsers(ctx, mountAccessors, coreLockedUsersPath+nsID+"/")
			if err != nil {
				return 0, nil, err
			}

			totalCounts += totalCountForNSID
			lockedUsersResponse = append(lockedUsersResponse, &LockedUsersResponse{
				NamespaceID:    strings.TrimSuffix(nsID, "/"),
				NamespacePath:  displayPath,
				Counts:         totalCountForNSID,
				MountAccessors: mountAccessorsResponse,
			})

		}
	}

	// sort namespaces in response by decreasing order of counts
	sort.Slice(lockedUsersResponse, func(i, j int) bool {
		return lockedUsersResponse[i].Counts > lockedUsersResponse[j].Counts
	})

	return totalCounts, lockedUsersResponse, nil
}

// getMountAccessorsLockedUsers returns the locked users for all the mount_accessors of locked users for a namespace
// they are sorted in the decreasing order of locked users
// returns the total locked users for the namespace and  locked users response for every mount_accessor for a namespace that has locked users
func (b *SystemBackend) getMountAccessorsLockedUsers(ctx context.Context, mountAccessors []string, lockedUsersPath string) (int, []*ResponseMountAccessors, error) {
	byMountAccessorsResponse := make([]*ResponseMountAccessors, 0)
	totalCountForMountAccessors := 0

	for _, mountAccessor := range mountAccessors {
		// get the list of aliases of locked users for a mount accessor
		aliasIdentifiers, err := b.Core.barrier.List(ctx, lockedUsersPath+mountAccessor)
		if err != nil {
			return 0, nil, err
		}

		totalCountForMountAccessors += len(aliasIdentifiers)
		byMountAccessorsResponse = append(byMountAccessorsResponse, &ResponseMountAccessors{
			MountAccessor:    strings.TrimSuffix(mountAccessor, "/"),
			Counts:           len(aliasIdentifiers),
			AliasIdentifiers: aliasIdentifiers,
		})

	}

	// sort mount Accessors in response by decreasing order of counts
	sort.Slice(byMountAccessorsResponse, func(i, j int) bool {
		return byMountAccessorsResponse[i].Counts > byMountAccessorsResponse[j].Counts
	})

	return totalCountForMountAccessors, byMountAccessorsResponse, nil
}

// includeNSInLockedUsersResponse checks if the namespace is the child namespace of namespace in query
// if child namespace, it can be included in response
// locked users from deleted namespaces are listed under root namespace
func (b *SystemBackend) includeNSInLockedUsersResponse(query *namespace.Namespace, record *namespace.Namespace) bool {
	if record == nil {
		// Deleted namespace, only include in root queries
		return query.ID == namespace.RootNamespaceID
	}
	return record.HasParent(query)
}
