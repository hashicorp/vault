// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package approle

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathTidySecretID(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "tidy/secret-id$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAppRole,
			OperationSuffix: "secret-id",
			OperationVerb:   "tidy",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathTidySecretIDUpdate,
				Responses: map[int][]framework.Response{
					http.StatusAccepted: {{
						Description: http.StatusText(http.StatusAccepted),
					}},
				},
			},
		},

		HelpSynopsis:    pathTidySecretIDSyn,
		HelpDescription: pathTidySecretIDDesc,
	}
}

// tidySecretID is used to delete entries in the whitelist that are expired.
func (b *backend) tidySecretID(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	// If we are a performance standby forward the request to the active node
	if b.System().ReplicationState().HasState(consts.ReplicationPerformanceStandby) {
		return nil, logical.ErrReadOnly
	}

	if !atomic.CompareAndSwapUint32(b.tidySecretIDCASGuard, 0, 1) {
		resp := &logical.Response{}
		resp.AddWarning("Tidy operation already in progress.")
		return resp, nil
	}

	go b.tidySecretIDinternal(req.Storage)

	resp := &logical.Response{}
	resp.AddWarning("Tidy operation successfully started. Any information from the operation will be printed to Vault's server logs.")
	return logical.RespondWithStatusCode(resp, req, http.StatusAccepted)
}

type tidyHelperSecretIDAccessor struct {
	secretIDAccessorStorageEntry
	saltedSecretIDAccessor string
}

func (b *backend) tidySecretIDinternal(s logical.Storage) {
	defer atomic.StoreUint32(b.tidySecretIDCASGuard, 0)

	logger := b.Logger().Named("tidy")

	checkCount := 0

	defer func() {
		logger.Trace("done checking entries", "num_entries", checkCount)
	}()

	// Don't cancel when the original client request goes away
	ctx := context.Background()

	salt, err := b.Salt(ctx)
	if err != nil {
		logger.Error("error tidying secret IDs", "error", err)
		return
	}

	tidyFunc := func(secretIDPrefixToUse, accessorIDPrefixToUse string) error {
		logger.Trace("listing accessors", "prefix", accessorIDPrefixToUse)

		// List all the accessors and add them all to a map
		// These hashes are the result of salting the accessor id.
		accessorHashes, err := s.List(ctx, accessorIDPrefixToUse)
		if err != nil {
			return err
		}
		skipHashes := make(map[string]bool, len(accessorHashes))
		accHashesByLockID := make([][]tidyHelperSecretIDAccessor, 256)
		for _, accessorHash := range accessorHashes {
			var entry secretIDAccessorStorageEntry
			entryIndex := accessorIDPrefixToUse + accessorHash
			se, err := s.Get(ctx, entryIndex)
			if err != nil {
				return err
			}
			if se == nil {
				continue
			}
			err = se.DecodeJSON(&entry)
			if err != nil {
				return err
			}

			lockIdx := locksutil.LockIndexForKey(entry.SecretIDHMAC)
			accHashesByLockID[lockIdx] = append(accHashesByLockID[lockIdx], tidyHelperSecretIDAccessor{
				secretIDAccessorStorageEntry: entry,
				saltedSecretIDAccessor:       accessorHash,
			})
		}

		secretIDCleanupFunc := func(secretIDHMAC, roleNameHMAC, secretIDPrefixToUse string) error {
			checkCount++
			lock := b.secretIDLock(secretIDHMAC)
			lock.Lock()
			defer lock.Unlock()

			entryIndex := fmt.Sprintf("%s%s%s", secretIDPrefixToUse, roleNameHMAC, secretIDHMAC)
			secretIDEntry, err := s.Get(ctx, entryIndex)
			if err != nil {
				return fmt.Errorf("error fetching SecretID %q: %w", secretIDHMAC, err)
			}

			if secretIDEntry == nil {
				logger.Error("entry for secret id was nil", "secret_id_hmac", secretIDHMAC)
				return nil
			}

			if secretIDEntry.Value == nil || len(secretIDEntry.Value) == 0 {
				return fmt.Errorf("found entry for SecretID %q but actual SecretID is empty", secretIDHMAC)
			}

			var result secretIDStorageEntry
			if err := secretIDEntry.DecodeJSON(&result); err != nil {
				return err
			}

			// If a secret ID entry does not have a corresponding accessor
			// entry, revoke the secret ID immediately
			accessorEntry, err := b.secretIDAccessorEntry(ctx, s, result.SecretIDAccessor, secretIDPrefixToUse)
			if err != nil {
				return fmt.Errorf("failed to read secret ID accessor entry: %w", err)
			}
			if accessorEntry == nil {
				logger.Trace("found nil accessor")
				if err := s.Delete(ctx, entryIndex); err != nil {
					return fmt.Errorf("error deleting secret ID %q from storage: %w", secretIDHMAC, err)
				}
				return nil
			}

			// ExpirationTime not being set indicates non-expiring SecretIDs
			if !result.ExpirationTime.IsZero() && time.Now().After(result.ExpirationTime) {
				logger.Trace("found expired secret ID")
				// Clean up the accessor of the secret ID first
				err = b.deleteSecretIDAccessorEntry(ctx, s, result.SecretIDAccessor, secretIDPrefixToUse)
				if err != nil {
					return fmt.Errorf("failed to delete secret ID accessor entry: %w", err)
				}

				if err := s.Delete(ctx, entryIndex); err != nil {
					return fmt.Errorf("error deleting SecretID %q from storage: %w", secretIDHMAC, err)
				}

				return nil
			}

			// At this point, the secret ID is not expired and is valid. Flag
			// the corresponding accessor as not needing attention.
			skipHashes[salt.SaltID(result.SecretIDAccessor)] = true

			return nil
		}

		logger.Trace("listing role HMACs", "prefix", secretIDPrefixToUse)

		roleNameHMACs, err := s.List(ctx, secretIDPrefixToUse)
		if err != nil {
			return err
		}

		for _, roleNameHMAC := range roleNameHMACs {
			logger.Trace("listing secret ID HMACs", "role_hmac", roleNameHMAC)
			secretIDHMACs, err := s.List(ctx, fmt.Sprintf("%s%s", secretIDPrefixToUse, roleNameHMAC))
			if err != nil {
				return err
			}
			for _, secretIDHMAC := range secretIDHMACs {
				err = secretIDCleanupFunc(secretIDHMAC, roleNameHMAC, secretIDPrefixToUse)
				if err != nil {
					return err
				}
			}
		}

		// Accessor indexes were not getting cleaned up until 0.9.3. This is a fix
		// to clean up the dangling accessor entries.
		if len(accessorHashes) > len(skipHashes) {
			// There is some raciness here because we're querying secretids for
			// roles without having a lock while doing so.  Because
			// accHashesByLockID was populated previously, at worst this may
			// mean that we fail to clean up something we ought to.
			allSecretIDHMACs := make(map[string]struct{})
			for _, roleNameHMAC := range roleNameHMACs {
				secretIDHMACs, err := s.List(ctx, secretIDPrefixToUse+roleNameHMAC)
				if err != nil {
					return err
				}
				for _, v := range secretIDHMACs {
					allSecretIDHMACs[v] = struct{}{}
				}
			}

			tidyEntries := func(entries []tidyHelperSecretIDAccessor) error {
				for _, entry := range entries {
					// Don't clean up accessor index entry if secretid cleanup func
					// determined that it should stay.
					if _, ok := skipHashes[entry.saltedSecretIDAccessor]; ok {
						continue
					}

					// Don't clean up accessor index entry if referenced in role.
					if _, ok := allSecretIDHMACs[entry.SecretIDHMAC]; ok {
						continue
					}

					if err := s.Delete(context.Background(), accessorIDPrefixToUse+entry.saltedSecretIDAccessor); err != nil {
						return err
					}
				}
				return nil
			}

			for lockIdx, entries := range accHashesByLockID {
				// Ideally, locking on accessors should be performed here too
				// but for that, accessors are required in plaintext, which are
				// not available.
				// ...
				// The lock is held when writing accessor/secret so if we have
				// the lock we know we're not in a
				// wrote-accessor-but-not-yet-secret case, which can be racy.
				b.secretIDLocks[lockIdx].Lock()
				err = tidyEntries(entries)
				b.secretIDLocks[lockIdx].Unlock()
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	err = tidyFunc(secretIDPrefix, secretIDAccessorPrefix)
	if err != nil {
		logger.Error("error tidying global secret IDs", "error", err)
		return
	}
	err = tidyFunc(secretIDLocalPrefix, secretIDAccessorLocalPrefix)
	if err != nil {
		logger.Error("error tidying local secret IDs", "error", err)
		return
	}
}

// pathTidySecretIDUpdate is used to delete the expired SecretID entries
func (b *backend) pathTidySecretIDUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.tidySecretID(ctx, req)
}

const (
	pathTidySecretIDSyn  = "Trigger the clean-up of expired SecretID entries."
	pathTidySecretIDDesc = `SecretIDs will have expiration time attached to them. The periodic function
of the backend will look for expired entries and delete them. This happens once in a minute. Invoking
this endpoint will trigger the clean-up action, without waiting for the backend's periodic function.`
)
