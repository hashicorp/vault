// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package openldap

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

const (
	// operationPrefixLDAP/LDAPLibrary are used as prefixes for OpenAPI operation id's.
	operationPrefixLDAP        = "ldap"
	operationPrefixLDAPLibrary = "ldap-library"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	ldapClient := NewClient(conf.Logger)
	b := Backend(ldapClient)
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}

	return b, nil
}

func Backend(client ldapClient) *backend {
	b := &backend{
		client:            client,
		credRotationQueue: queue.New(),
		roleLocks:         locksutil.CreateLocks(),
		checkOutLocks:     locksutil.CreateLocks(),
		managedUsers:      make(map[string]struct{}),
	}

	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),
		PathsSpecial: &logical.Paths{
			LocalStorage: []string{
				framework.WALPrefix,
			},
			SealWrapStorage: []string{
				configPath,
				staticRolePath + "*",
			},
		},
		Paths: framework.PathAppend(
			// These paths must be at the top of the list since their regex
			// Patterns are the most specific. Otherwise, a more generic regex
			// will swallow the request because role and set names can contain
			// arbitrary numbers of slashes.
			// For example, a request to `library/:set_name/check-in` could be
			// swallowed by the regex for `library/:set_name`.
			b.pathSetManageCheckIn(),
			b.pathSetCheckIn(),
			b.pathSetCheckOut(),
			b.pathSetStatus(),

			// These paths are more generic than the above. They must be
			// appended last.
			b.pathConfig(),
			b.pathDynamicRoles(),
			b.pathDynamicCredsCreate(),
			b.pathStaticRoles(),
			b.pathStaticCredsCreate(),
			b.pathListStaticRoles(),
			b.pathRotateCredentials(),
			b.pathSets(),
			b.pathListSets(),
		),
		InitializeFunc: b.initialize,
		Secrets: []*framework.Secret{
			dynamicSecretCreds(b),
			checkoutSecretCreds(b),
		},
		Clean:       b.clean,
		BackendType: logical.TypeLogical,
	}

	return b
}

func (b *backend) initialize(ctx context.Context, initRequest *logical.InitializationRequest) error {
	// Load managed LDAP users into memory from storage
	staticRoles, err := b.loadManagedUsers(ctx, initRequest.Storage)
	if err != nil {
		return err
	}

	// Create a context with a cancel method for processing any WAL entries and
	// populating the queue
	ictx, cancel := context.WithCancel(context.Background())
	b.cancelQueue = cancel

	// Load static role queue and kickoff new periodic ticker
	go b.initQueue(ictx, initRequest, staticRoles)

	return nil
}

func (b *backend) clean(_ context.Context) {
	b.invalidateQueue()
}

// invalidateQueue cancels any background queue loading and destroys the queue.
func (b *backend) invalidateQueue() {
	b.Lock()
	defer b.Unlock()

	if b.cancelQueue != nil {
		b.cancelQueue()
	}
	b.credRotationQueue = nil
}

type backend struct {
	*framework.Backend
	sync.RWMutex
	// CredRotationQueue is an in-memory priority queue used to track Static Roles
	// that require periodic rotation. Backends will have a PriorityQueue
	// initialized on setup, but only backends that are mounted by a primary
	// server or mounted as a local mount will perform the rotations.
	//
	// cancelQueue is used to remove the priority queue and terminate the
	// background ticker.
	credRotationQueue *queue.PriorityQueue
	cancelQueue       context.CancelFunc

	// roleLocks is used to lock modifications to roles in the queue, to ensure
	// concurrent requests are not modifying the same role and possibly causing
	// issues with the priority queue.
	roleLocks []*locksutil.LockEntry
	client    ldapClient

	// managedUsers contains the set of LDAP usernames managed by the secrets engine
	// static role and check-in/check-out systems. It is used to ensure that users
	// are exclusively managed by one system and not both. Access to managedUsers is
	// synchronized by the managedUserLock.
	managedUsers    map[string]struct{}
	managedUserLock sync.Mutex

	// checkOutLocks are used for avoiding races when working with library sets
	// in the check-in/check-out system.
	checkOutLocks []*locksutil.LockEntry
}

// walkfunc type takes a storage path argument and returns true if a storage
// entry exists, false otherwise
type walkFunc func(string) (bool, error)

// walkStoragePath performs a non-recursive breadth-first search of the given
// storage path and applies the walkfunc to each storage entry
func walkStoragePath(ctx context.Context, s logical.Storage, path string, walker walkFunc) error {
	keys, err := s.List(ctx, path)
	if err != nil {
		return fmt.Errorf("unable to list keys: %w", err)
	}

	// Storage entries can be defined with hierarchical paths, e.g.
	// "foo/bar/baz". But the storage.List() call to the top-level key will
	// only return the top-level keys in the hierarchy. So we perform a
	// non-recursive breadth-first search through all the keys returned from
	// storage.List() and apply the walkfunc.
	for i := 0; i < len(keys); i++ {
		key := keys[i]

		entryExists, err := walker(key)
		if err != nil {
			return fmt.Errorf("unable to read entry %q: %w", key, err)
		}

		if !entryExists && strings.HasSuffix(key, "/") {
			// this is a directory
			subKeys, err := s.List(ctx, path+key)
			if err != nil {
				return fmt.Errorf("unable to list keys: %w", err)
			}

			// append to the keys slice to continue search in the sub-directory
			for _, subKey := range subKeys {
				// prevent infinite loop but this should never happen
				if subKey == "" {
					continue
				}
				subKey = fmt.Sprintf("%s%s", key, subKey)
				keys = append(keys, subKey)
			}
		}
	}
	return nil
}

// loadManagedUsers loads users managed by the secrets engine from storage into
// the backend's managedUsers set. Users are loaded from both the static role and
// check-in/check-out systems. Returns an error if one occurs during loading.
func (b *backend) loadManagedUsers(ctx context.Context, s logical.Storage) (map[string]*roleEntry, error) {
	log := b.Logger()
	b.managedUserLock.Lock()
	defer b.managedUserLock.Unlock()

	// Clear managed users before loading to ensure that the full set is
	// loaded from storage. This is important during initialization after
	// leadership changes to keep the set consistent with storage.
	b.managedUsers = make(map[string]struct{})

	// Load users managed under static roles
	roles := map[string]*roleEntry{}
	roleFunc := func(roleName string) (bool, error) {
		entry, err := b.staticRole(ctx, s, roleName)
		if err != nil {
			log.Warn("unable to read static role", "error", err, "role", roleName)
			return false, err
		}
		entryExists := entry != nil && !strings.HasSuffix(roleName, "/")
		if entryExists {
			roles[roleName] = entry
		}
		return entryExists, nil
	}
	err := walkStoragePath(ctx, s, staticRolePath, roleFunc)
	if err != nil {
		return nil, err
	}

	for roleName, role := range roles {
		if role == nil || role.StaticAccount == nil {
			// This indicates that a static role returned from the list operation was
			// deleted before the read operation in this loop. This shouldn't happen
			// at this point in the plugin lifecycle, so we'll log if it does.
			b.Logger().Debug("unexpected nil static role found while loading managed users",
				"name", roleName)
			continue
		}

		// Add the static role user to the managed user set
		b.managedUsers[role.StaticAccount.Username] = struct{}{}
	}

	// Load users managed under library sets
	librarySets := map[string]*librarySet{}
	setFunc := func(setName string) (bool, error) {
		entry, err := readSet(ctx, s, setName)
		if err != nil {
			log.Warn("unable to read library set", "error", err, "set", setName)
			return false, err
		}
		entryExists := entry != nil && !strings.HasSuffix(setName, "/")
		if entryExists {
			librarySets[setName] = entry
		}
		return entryExists, nil
	}
	err = walkStoragePath(ctx, s, libraryPrefix, setFunc)
	if err != nil {
		return nil, err
	}
	for setName, set := range librarySets {
		if set == nil {
			// This indicates that a library set returned from the list operation was
			// deleted before the read operation in this loop. This shouldn't happen
			// at this point in the plugin lifecycle, so we'll log if it does.
			b.Logger().Debug("unexpected nil library set found while loading managed users",
				"name", setName)
			continue
		}

		// Add the service account names to the managed user set
		for _, name := range set.ServiceAccountNames {
			b.managedUsers[name] = struct{}{}
		}
	}

	return roles, nil
}

const backendHelp = `
The LDAP backend supports managing existing LDAP entry passwords by providing:

 * end points to add entries
 * manual rotation of entry passwords
 * auto rotation of entry passwords
 * check-in/check-out for libraries of entries
 
The LDAP secret engine supports OpenLDAP, Active Directory, and IBM RACF 
implementations via schema configuration.

After mounting this secret backend, configure it using the "ldap/config" path.
`
