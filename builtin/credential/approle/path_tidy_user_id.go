package approle

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-multierror"
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
func (b *backend) tidySecretID(ctx context.Context, s logical.Storage) (*logical.Response, error) {
	if !atomic.CompareAndSwapUint32(b.tidySecretIDCASGuard, 0, 1) {
		resp := &logical.Response{}
		resp.AddWarning("Tidy operation already in progress.")
		return resp, nil
	}

	go func() {
		defer atomic.StoreUint32(b.tidySecretIDCASGuard, 0)

		var result error

		// Don't cancel when the original client request goes away
		ctx = context.Background()

		logger := b.Logger().Named("tidy")

		tidyFunc := func(secretIDPrefixToUse, accessorIDPrefixToUse string) error {
			roleNameHMACs, err := s.List(ctx, secretIDPrefixToUse)
			if err != nil {
				return err
			}

			// List all the accessors and add them all to a map
			accessorHashes, err := s.List(ctx, accessorIDPrefixToUse)
			if err != nil {
				return err
			}
			accessorMap := make(map[string]bool, len(accessorHashes))
			for _, accessorHash := range accessorHashes {
				accessorMap[accessorHash] = true
			}

			secretIDCleanupFunc := func(secretIDHMAC, roleNameHMAC, secretIDPrefixToUse string) error {
				lock := b.secretIDLock(secretIDHMAC)
				lock.Lock()
				defer lock.Unlock()

				entryIndex := fmt.Sprintf("%s%s%s", secretIDPrefixToUse, roleNameHMAC, secretIDHMAC)
				secretIDEntry, err := s.Get(ctx, entryIndex)
				if err != nil {
					return errwrap.Wrapf(fmt.Sprintf("error fetching SecretID %q: {{err}}", secretIDHMAC), err)
				}

				if secretIDEntry == nil {
					result = multierror.Append(result, fmt.Errorf("entry for SecretID %q is nil", secretIDHMAC))
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
					if err := s.Delete(ctx, entryIndex); err != nil {
						return errwrap.Wrapf(fmt.Sprintf("error deleting secret ID %q from storage: {{err}}", secretIDHMAC), err)
					}
					return nil
				}

				// ExpirationTime not being set indicates non-expiring SecretIDs
				if !result.ExpirationTime.IsZero() && time.Now().After(result.ExpirationTime) {
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
			for accessorHash, _ := range accessorMap {
				// Ideally, locking should be performed here. But for that, accessors
				// are required in plaintext, which are not available. Hence performing
				// a racy cleanup.
				err = s.Delete(ctx, secretIDAccessorPrefix+accessorHash)
				if err != nil {
					return err
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
	return resp, nil
}

// pathTidySecretIDUpdate is used to delete the expired SecretID entries
func (b *backend) pathTidySecretIDUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return b.tidySecretID(ctx, req.Storage)
}

const pathTidySecretIDSyn = "Trigger the clean-up of expired SecretID entries."
const pathTidySecretIDDesc = `SecretIDs will have expiration time attached to them. The periodic function
of the backend will look for expired entries and delete them. This happens once in a minute. Invoking
this endpoint will trigger the clean-up action, without waiting for the backend's periodic function.`
