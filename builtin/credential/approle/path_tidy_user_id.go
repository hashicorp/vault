package approle

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathTidySecretID(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "tidy/secret-id$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathTidySecretIDUpdate,
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

	s := req.Storage

	go func() {
		defer atomic.StoreUint32(b.tidySecretIDCASGuard, 0)

		logger := b.Logger().Named("tidy")

		checkCount := 0

		defer func() {
			if b.testTidyDelay > 0 {
				logger.Trace("done checking entries", "num_entries", checkCount)
			}
		}()

		// Don't cancel when the original client request goes away
		ctx = context.Background()

		tidyFunc := func(secretIDPrefixToUse, accessorIDPrefixToUse string) error {
			logger.Trace("listing role HMACs", "prefix", secretIDPrefixToUse)

			roleNameHMACs, err := s.List(ctx, secretIDPrefixToUse)
			if err != nil {
				return err
			}

			logger.Trace("listing accessors", "prefix", accessorIDPrefixToUse)

			// List all the accessors and add them all to a map
			accessorHashes, err := s.List(ctx, accessorIDPrefixToUse)
			if err != nil {
				return err
			}
			accessorMap := make(map[string]bool, len(accessorHashes))
			for _, accessorHash := range accessorHashes {
				accessorMap[accessorHash] = true
			}

			time.Sleep(b.testTidyDelay)

			secretIDCleanupFunc := func(secretIDHMAC, roleNameHMAC, secretIDPrefixToUse string) error {
				checkCount++
				lock := b.secretIDLock(secretIDHMAC)
				lock.Lock()
				defer lock.Unlock()

				entryIndex := fmt.Sprintf("%s%s%s", secretIDPrefixToUse, roleNameHMAC, secretIDHMAC)
				secretIDEntry, err := s.Get(ctx, entryIndex)
				if err != nil {
					return errwrap.Wrapf(fmt.Sprintf("error fetching SecretID %q: {{err}}", secretIDHMAC), err)
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
					return errwrap.Wrapf("failed to read secret ID accessor entry: {{err}}", err)
				}
				if accessorEntry == nil {
					logger.Trace("found nil accessor")
					if err := s.Delete(ctx, entryIndex); err != nil {
						return errwrap.Wrapf(fmt.Sprintf("error deleting secret ID %q from storage: {{err}}", secretIDHMAC), err)
					}
					return nil
				}

				// ExpirationTime not being set indicates non-expiring SecretIDs
				if !result.ExpirationTime.IsZero() && time.Now().After(result.ExpirationTime) {
					logger.Trace("found expired secret ID")
					// Clean up the accessor of the secret ID first
					err = b.deleteSecretIDAccessorEntry(ctx, s, result.SecretIDAccessor, secretIDPrefixToUse)
					if err != nil {
						return errwrap.Wrapf("failed to delete secret ID accessor entry: {{err}}", err)
					}

					if err := s.Delete(ctx, entryIndex); err != nil {
						return errwrap.Wrapf(fmt.Sprintf("error deleting SecretID %q from storage: {{err}}", secretIDHMAC), err)
					}

					return nil
				}

				// At this point, the secret ID is not expired and is valid. Delete
				// the corresponding accessor from the accessorMap. This will leave
				// only the dangling accessors in the map which can then be cleaned
				// up later.
				salt, err := b.Salt(ctx)
				if err != nil {
					return err
				}
				delete(accessorMap, salt.SaltID(result.SecretIDAccessor))

				return nil
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
			if len(accessorMap) > 0 {
				for _, lock := range b.secretIDLocks {
					lock.Lock()
					defer lock.Unlock()
				}
				for accessorHash, _ := range accessorMap {
					logger.Trace("found dangling accessor, verifying")
					// Ideally, locking on accessors should be performed here too
					// but for that, accessors are required in plaintext, which are
					// not available. The code above helps but it may still be
					// racy.
					// ...
					// Look up the secret again now that we have all the locks. The
					// lock is held when writing accessor/secret so if we have the
					// lock we know we're not in a
					// wrote-accessor-but-not-yet-secret case, which can be racy.
					var entry secretIDAccessorStorageEntry
					entryIndex := accessorIDPrefixToUse + accessorHash
					se, err := s.Get(ctx, entryIndex)
					if err != nil {
						return err
					}
					if se != nil {
						err = se.DecodeJSON(&entry)
						if err != nil {
							return err
						}

						// The storage entry doesn't store the role ID, so we have
						// to go about this the long way; fortunately we shouldn't
						// actually hit this very often
						var found bool
					searchloop:
						for _, roleNameHMAC := range roleNameHMACs {
							secretIDHMACs, err := s.List(ctx, fmt.Sprintf("%s%s", secretIDPrefixToUse, roleNameHMAC))
							if err != nil {
								return err
							}
							for _, v := range secretIDHMACs {
								if v == entry.SecretIDHMAC {
									found = true
									logger.Trace("accessor verified, not removing")
									break searchloop
								}
							}
						}
						if !found {
							logger.Trace("could not verify dangling accessor, removing")
							err = s.Delete(ctx, entryIndex)
							if err != nil {
								return err
							}
						}
					}
				}
			}

			return nil
		}

		err := tidyFunc(secretIDPrefix, secretIDAccessorPrefix)
		if err != nil {
			logger.Error("error tidying global secret IDs", "error", err)
			return
		}
		err = tidyFunc(secretIDLocalPrefix, secretIDAccessorLocalPrefix)
		if err != nil {
			logger.Error("error tidying local secret IDs", "error", err)
			return
		}
	}()

	resp := &logical.Response{}
	resp.AddWarning("Tidy operation successfully started. Any information from the operation will be printed to Vault's server logs.")
	return logical.RespondWithStatusCode(resp, req, http.StatusAccepted)
}

// pathTidySecretIDUpdate is used to delete the expired SecretID entries
func (b *backend) pathTidySecretIDUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.tidySecretID(ctx, req)
}

const pathTidySecretIDSyn = "Trigger the clean-up of expired SecretID entries."
const pathTidySecretIDDesc = `SecretIDs will have expiration time attached to them. The periodic function
of the backend will look for expired entries and delete them. This happens once in a minute. Invoking
this endpoint will trigger the clean-up action, without waiting for the backend's periodic function.`
