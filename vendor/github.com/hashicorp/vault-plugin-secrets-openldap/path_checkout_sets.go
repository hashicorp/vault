// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package openldap

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	libraryPrefix       = "library/"
	libraryManagePrefix = "library/manage/"
)

type librarySet struct {
	ServiceAccountNames       []string      `json:"service_account_names"`
	TTL                       time.Duration `json:"ttl"`
	MaxTTL                    time.Duration `json:"max_ttl"`
	DisableCheckInEnforcement bool          `json:"disable_check_in_enforcement"`
}

// Validate ensures that a set meets our code assumptions that TTLs are set in
// a way that makes sense, and that there's at least one service account.
func (l *librarySet) Validate() error {
	if len(l.ServiceAccountNames) < 1 {
		return fmt.Errorf("at least one service account must be configured")
	}
	for _, name := range l.ServiceAccountNames {
		if name == "" {
			return fmt.Errorf("service account name must not be empty")
		}
	}
	if l.MaxTTL > 0 {
		if l.MaxTTL < l.TTL {
			return fmt.Errorf("max_ttl (%d seconds) may not be less than ttl (%d seconds)", l.MaxTTL, l.TTL)
		}
	}
	return nil
}

func (b *backend) pathListSets() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: strings.TrimSuffix(libraryPrefix, "/") + optionalGenericNameWithForwardSlashListRegex("path"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixLDAPLibrary,
				OperationVerb:   "list",
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.listSetsOperation,
				},
			},
			Fields: map[string]*framework.FieldSchema{
				"path": {
					Type:        framework.TypeLowerCaseString,
					Description: "Path of sets to list",
				},
			},
			HelpSynopsis:    pathListSetsHelpSyn,
			HelpDescription: pathListSetsHelpDesc,
		},
	}
}

func (b *backend) listSetsOperation(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	setPath := data.Get("path").(string)
	keys, err := req.Storage.List(ctx, libraryPrefix+setPath)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(keys), nil
}

func (b *backend) pathSets() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: strings.TrimSuffix(libraryPrefix, "/") + genericNameWithForwardSlashRegex("name"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixLDAPLibrary,
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeLowerCaseString,
					Description: "Name of the set.",
					Required:    true,
				},
				"service_account_names": {
					Type:        framework.TypeCommaStringSlice,
					Description: "The username/logon name for the service accounts with which this set will be associated.",
				},
				"ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "In seconds, the amount of time a check-out should last. Defaults to 24 hours.",
					Default:     24 * 60 * 60, // 24 hours
				},
				"max_ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "In seconds, the max amount of time a check-out's renewals should last. Defaults to 24 hours.",
					Default:     24 * 60 * 60, // 24 hours
				},
				"disable_check_in_enforcement": {
					Type:        framework.TypeBool,
					Description: "Disable the default behavior of requiring that check-ins are performed by the entity that checked them out.",
					Default:     false,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.operationSetCreate,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "configure",
					},
					Summary:                     "Create a library set.",
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.operationSetUpdate,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "configure",
					},
					Summary:                     "Update a library set.",
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.operationSetRead,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "read",
					},
					Summary: "Read a library set.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.operationSetDelete,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "delete",
					},
					Summary:                     "Delete a library set.",
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
			},
			ExistenceCheck:  b.operationSetExistenceCheck,
			HelpSynopsis:    setHelpSynopsis,
			HelpDescription: setHelpDescription,
		},
	}
}

func (b *backend) operationSetExistenceCheck(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (bool, error) {
	setName := fieldData.Get("name").(string)

	lock := locksutil.LockForKey(b.checkOutLocks, setName)
	lock.RLock()
	defer lock.RUnlock()

	set, err := readSet(ctx, req.Storage, setName)
	if err != nil {
		return false, err
	}
	return set != nil, nil
}

func (b *backend) operationSetCreate(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (*logical.Response, error) {
	setName := fieldData.Get("name").(string)

	lock := locksutil.LockForKey(b.checkOutLocks, setName)
	lock.Lock()
	defer lock.Unlock()

	config, err := readConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return logical.ErrorResponse("config must exist to create a library"), nil
	}

	serviceAccountNames := fieldData.Get("service_account_names").([]string)
	ttl := time.Duration(fieldData.Get("ttl").(int)) * time.Second
	maxTTL := time.Duration(fieldData.Get("max_ttl").(int)) * time.Second
	disableCheckInEnforcement := fieldData.Get("disable_check_in_enforcement").(bool)

	if len(serviceAccountNames) == 0 {
		return logical.ErrorResponse(`"service_account_names" must be provided`), nil
	}

	b.managedUserLock.Lock()
	defer b.managedUserLock.Unlock()

	// Ensure the given service accounts aren't already managed by
	// another library set or static role.
	for _, serviceAccountName := range serviceAccountNames {
		if _, exists := b.managedUsers[serviceAccountName]; exists {
			return logical.ErrorResponse("%q is already managed by the secrets engine",
				serviceAccountName), nil
		}
	}

	set := &librarySet{
		ServiceAccountNames:       serviceAccountNames,
		TTL:                       ttl,
		MaxTTL:                    maxTTL,
		DisableCheckInEnforcement: disableCheckInEnforcement,
	}
	if err := set.Validate(); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	for _, serviceAccountName := range serviceAccountNames {
		if err := b.CheckIn(ctx, req.Storage, serviceAccountName); err != nil {
			return nil, err
		}
	}
	if err := storeSet(ctx, req.Storage, setName, set); err != nil {
		return nil, err
	}

	// Add the service account names to the managed user set
	for _, name := range set.ServiceAccountNames {
		b.managedUsers[name] = struct{}{}
	}

	return nil, nil
}

func (b *backend) operationSetUpdate(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (*logical.Response, error) {
	setName := fieldData.Get("name").(string)

	lock := locksutil.LockForKey(b.checkOutLocks, setName)
	lock.Lock()
	defer lock.Unlock()

	config, err := readConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return logical.ErrorResponse("config must exist to update a library"), nil
	}

	newServiceAccountNamesRaw, newServiceAccountNamesSent := fieldData.GetOk("service_account_names")
	var newServiceAccountNames []string
	if newServiceAccountNamesSent {
		newServiceAccountNames = newServiceAccountNamesRaw.([]string)
	}

	ttlRaw, ttlSent := fieldData.GetOk("ttl")
	if !ttlSent {
		ttlRaw = fieldData.Schema["ttl"].Default
	}
	ttl := time.Duration(ttlRaw.(int)) * time.Second

	maxTTLRaw, maxTTLSent := fieldData.GetOk("max_ttl")
	if !maxTTLSent {
		maxTTLRaw = fieldData.Schema["max_ttl"].Default
	}
	maxTTL := time.Duration(maxTTLRaw.(int)) * time.Second

	disableCheckInEnforcementRaw, enforcementSent := fieldData.GetOk("disable_check_in_enforcement")
	if !enforcementSent {
		disableCheckInEnforcementRaw = false
	}
	disableCheckInEnforcement := disableCheckInEnforcementRaw.(bool)

	set, err := readSet(ctx, req.Storage, setName)
	if err != nil {
		return nil, err
	}
	if set == nil {
		return logical.ErrorResponse(fmt.Sprintf(`%q doesn't exist`, setName)), nil
	}

	b.managedUserLock.Lock()
	defer b.managedUserLock.Unlock()

	var beingAdded []string
	var beingDeleted []string
	if newServiceAccountNamesSent {
		// For new service accounts we receive, before we check them in, ensure they're not in another set.
		beingAdded = strutil.Difference(newServiceAccountNames, set.ServiceAccountNames, true)
		for _, newServiceAccountName := range beingAdded {
			if _, exists := b.managedUsers[newServiceAccountName]; exists {
				return logical.ErrorResponse("%q is already managed by the secrets engine",
					newServiceAccountName), nil
			}
		}

		// For service accounts we won't be handling anymore, before we delete them, ensure they're not checked out.
		beingDeleted = strutil.Difference(set.ServiceAccountNames, newServiceAccountNames, true)
		for _, prevServiceAccountName := range beingDeleted {
			checkOut, err := b.LoadCheckOut(ctx, req.Storage, prevServiceAccountName)
			if err != nil {
				if err == errNotFound {
					// Nothing else to do here.
					continue
				}
				return nil, err
			}
			if !checkOut.IsAvailable {
				return logical.ErrorResponse(fmt.Sprintf(`"%s" can't be deleted because it is currently checked out'`, prevServiceAccountName)), nil
			}
		}
		set.ServiceAccountNames = newServiceAccountNames
	}

	if ttlSent {
		set.TTL = ttl
	}
	if maxTTLSent {
		set.MaxTTL = maxTTL
	}
	if enforcementSent {
		set.DisableCheckInEnforcement = disableCheckInEnforcement
	}

	if err := set.Validate(); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// Now that we know we can take all these actions, let's take them.
	for _, newServiceAccountName := range beingAdded {
		if err := b.CheckIn(ctx, req.Storage, newServiceAccountName); err != nil {
			return nil, err
		}
	}
	for _, prevServiceAccountName := range beingDeleted {
		if err := b.DeleteCheckout(ctx, req.Storage, prevServiceAccountName); err != nil {
			return nil, err
		}
	}
	if err := storeSet(ctx, req.Storage, setName, set); err != nil {
		return nil, err
	}

	// Delete the prior service account names from the managed user set
	for _, name := range beingDeleted {
		delete(b.managedUsers, name)
	}

	// Add the new service account names to the managed user set
	for _, name := range beingAdded {
		b.managedUsers[name] = struct{}{}
	}

	return nil, nil
}

func (b *backend) operationSetRead(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (*logical.Response, error) {
	setName := fieldData.Get("name").(string)

	lock := locksutil.LockForKey(b.checkOutLocks, setName)
	lock.RLock()
	defer lock.RUnlock()

	set, err := readSet(ctx, req.Storage, setName)
	if err != nil {
		return nil, err
	}
	if set == nil {
		return nil, nil
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"service_account_names":        set.ServiceAccountNames,
			"ttl":                          int64(set.TTL.Seconds()),
			"max_ttl":                      int64(set.MaxTTL.Seconds()),
			"disable_check_in_enforcement": set.DisableCheckInEnforcement,
		},
	}, nil
}

func (b *backend) operationSetDelete(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (*logical.Response, error) {
	setName := fieldData.Get("name").(string)

	lock := locksutil.LockForKey(b.checkOutLocks, setName)
	lock.Lock()
	defer lock.Unlock()

	set, err := readSet(ctx, req.Storage, setName)
	if err != nil {
		return nil, err
	}
	if set == nil {
		return nil, nil
	}

	// We need to remove all the items we'd stored for these service accounts.
	for _, serviceAccountName := range set.ServiceAccountNames {
		checkOut, err := b.LoadCheckOut(ctx, req.Storage, serviceAccountName)
		if err != nil {
			if err == errNotFound {
				// Nothing else to do here.
				continue
			}
			return nil, err
		}
		if !checkOut.IsAvailable {
			return logical.ErrorResponse(fmt.Sprintf(`"%s" can't be deleted because it is currently checked out'`, serviceAccountName)), nil
		}
	}
	for _, serviceAccountName := range set.ServiceAccountNames {
		if err := b.DeleteCheckout(ctx, req.Storage, serviceAccountName); err != nil {
			return nil, err
		}
	}
	if err := req.Storage.Delete(ctx, libraryPrefix+setName); err != nil {
		return nil, err
	}

	// Delete the service account names from the managed user set
	b.managedUserLock.Lock()
	defer b.managedUserLock.Unlock()
	for _, name := range set.ServiceAccountNames {
		delete(b.managedUsers, name)
	}

	return nil, nil
}

// readSet is a helper method for reading a set from storage by name.
// It's intended to be used anywhere in the plugin. It may return nil, nil if
// a librarySet doesn't currently exist for a given setName.
func readSet(ctx context.Context, storage logical.Storage, setName string) (*librarySet, error) {
	entry, err := storage.Get(ctx, libraryPrefix+setName)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	set := &librarySet{}
	if err := entry.DecodeJSON(set); err != nil {
		return nil, err
	}
	return set, nil
}

// storeSet stores a librarySet.
func storeSet(ctx context.Context, storage logical.Storage, setName string, set *librarySet) error {
	entry, err := logical.StorageEntryJSON(libraryPrefix+setName, set)
	if err != nil {
		return err
	}
	return storage.Put(ctx, entry)
}

const (
	setHelpSynopsis = `
Build a library of service accounts that can be checked out.
`
	setHelpDescription = `
This endpoint allows you to read, write, and delete individual sets of service accounts for check-out.
Deleting a set of service accounts can only be performed if all its accounts are currently checked in.
`
	pathListSetsHelpSyn = `
List the name of each set of service accounts currently stored.
`
	pathListSetsHelpDesc = `
To learn which service accounts are being managed by Vault, list the set names using
this endpoint. Then read any individual set by name to learn more.
`
)
