package approle

import (
	"fmt"
	"sync/atomic"
	"time"

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
func (b *backend) tidySecretID(s logical.Storage) error {
	grabbed := atomic.CompareAndSwapUint32(&b.tidySecretIDCASGuard, 0, 1)
	if grabbed {
		defer atomic.StoreUint32(&b.tidySecretIDCASGuard, 0)
	} else {
		return fmt.Errorf("SecretID tidy operation already running")
	}

	roleNameHMACs, err := s.List("secret_id/")
	if err != nil {
		return err
	}

	var result error
	for _, roleNameHMAC := range roleNameHMACs {
		// roleNameHMAC will already have a '/' suffix. Don't append another one.
		secretIDHMACs, err := s.List(fmt.Sprintf("secret_id/%s", roleNameHMAC))
		if err != nil {
			return err
		}
		for _, secretIDHMAC := range secretIDHMACs {
			// In order to avoid lock swroleing in case there is need to delete,
			// grab the write lock.
			lock := b.secretIDLock(secretIDHMAC)
			lock.Lock()
			// roleNameHMAC will already have a '/' suffix. Don't append another one.
			entryIndex := fmt.Sprintf("secret_id/%s%s", roleNameHMAC, secretIDHMAC)
			secretIDEntry, err := s.Get(entryIndex)
			if err != nil {
				lock.Unlock()
				return fmt.Errorf("error fetching SecretID %s: %s", secretIDHMAC, err)
			}

			if secretIDEntry == nil {
				result = multierror.Append(result, fmt.Errorf("entry for SecretID %s is nil", secretIDHMAC))
				lock.Unlock()
				continue
			}

			if secretIDEntry.Value == nil || len(secretIDEntry.Value) == 0 {
				lock.Unlock()
				return fmt.Errorf("found entry for SecretID %s but actual SecretID is empty", secretIDHMAC)
			}

			var result secretIDStorageEntry
			if err := secretIDEntry.DecodeJSON(&result); err != nil {
				lock.Unlock()
				return err
			}

			// ExpirationTime not being set indicates non-expiring SecretIDs
			if !result.ExpirationTime.IsZero() && time.Now().After(result.ExpirationTime) {
				if err := s.Delete(entryIndex); err != nil {
					lock.Unlock()
					return fmt.Errorf("error deleting SecretID %s from storage: %s", secretIDHMAC, err)
				}
			}
			lock.Unlock()
		}
	}
	return result
}

// pathTidySecretIDUpdate is used to delete the expired SecretID entries
func (b *backend) pathTidySecretIDUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return nil, b.tidySecretID(req.Storage)
}

const pathTidySecretIDSyn = "Trigger the clean-up of expired SecretID entries."
const pathTidySecretIDDesc = `SecretIDs will have expiratin time attached to them. The periodic function
of the backend will look for expired entries and delete them. This happens once in a minute. Invoking
this endpoint will trigger the clean-up action, without waiting for the backend's periodic function.`
