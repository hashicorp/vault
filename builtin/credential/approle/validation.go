package approle

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// secretIDStorageEntry represents the information stored in storage
// when a SecretID is created. The structure of the SecretID storage
// entry is the same for all the types of SecretIDs generated.
type secretIDStorageEntry struct {
	// Accessor for the SecretID. It is a random UUID serving as
	// a secondary index for the SecretID. This uniquely identifies
	// the SecretID it belongs to, and hence can be used for listing
	// and deleting SecretIDs. Accessors cannot be used as valid
	// SecretIDs during login.
	SecretIDAccessor string `json:"secret_id_accessor" structs:"secret_id_accessor" mapstructure:"secret_id_accessor"`

	// Number of times this SecretID can be used to perform the login operation
	SecretIDNumUses int `json:"secret_id_num_uses" structs:"secret_id_num_uses" mapstructure:"secret_id_num_uses"`

	// Duration after which this SecretID should expire. This is
	// croleed by the backend mount's max TTL value.
	SecretIDTTL time.Duration `json:"secret_id_ttl" structs:"secret_id_ttl" mapstructure:"secret_id_ttl"`

	// The time when the SecretID was created
	CreationTime time.Time `json:"creation_time" structs:"creation_time" mapstructure:"creation_time"`

	// The time when the SecretID becomes eligible for tidy
	// operation. Tidying is performed by the PeriodicFunc of the
	// backend which is 1 minute apart.
	ExpirationTime time.Time `json:"expiration_time" structs:"expiration_time" mapstructure:"expiration_time"`

	// The time representing the last time this storage entry was modified
	LastUpdatedTime time.Time `json:"last_updated_time" structs:"last_updated_time" mapstructure:"last_updated_time"`

	// Metadata that belongs to the SecretID.
	Metadata map[string]string `json:"metadata" structs:"metadata" mapstructure:"metadata"`
}

// Represents the payload of the storage entry of the accessor that maps to a unique
// SecretID. Note that SecretIDs should never be stored in plaintext anywhere in the
// backend. SecretIDHMAC will be used as an index to fetch the properties of the
// SecretID and to delete the SecretID.
type secretIDAccessorStorageEntry struct {
	// Hash of the SecretID which can be used to find the storage index at which
	// properties of SecretID is stored.
	SecretIDHMAC string `json:"secret_id_hmac" structs:"secret_id_hmac" mapstructure:"secret_id_hmac"`
}

// Checks if all the CIDR blocks in the comma separated list are valid by parsing it.
func validateCIDRList(cidrList string) error {
	if cidrList == "" {
		return nil
	}

	cidrBlocks := strings.Split(cidrList, ",")
	for _, block := range cidrBlocks {
		if _, _, err := net.ParseCIDR(strings.TrimSpace(block)); err != nil {
			return err
		}
	}
	return nil
}

// Checks if the Role represented by the RoleID still exists
func (b *backend) validateRoleID(s logical.Storage, roleID string) (*roleStorageEntry, string, error) {
	// Look for the storage entry that maps the roleID to role
	roleIDIndex, err := b.roleIDEntry(s, roleID)
	if err != nil {
		return nil, "", err
	}
	if roleIDIndex == nil {
		return nil, "", fmt.Errorf("failed to find secondary index for role_id:%s\n", roleID)
	}

	role, err := b.roleEntry(s, roleIDIndex.Name)
	if err != nil {
		return nil, "", err
	}
	if role == nil {
		return nil, "", fmt.Errorf("role %s referred by the SecretID does not exist", roleIDIndex.Name)
	}

	return role, roleIDIndex.Name, nil
}

// Validates the supplied RoleID and SecretID
func (b *backend) validateCredentials(req *logical.Request, data *framework.FieldData) (*roleStorageEntry, string, map[string]string, error) {
	var metadata map[string]string
	// RoleID must be supplied during every login
	roleID := strings.TrimSpace(data.Get("role_id").(string))
	if roleID == "" {
		return nil, "", metadata, fmt.Errorf("missing role_id")
	}

	// Validate the RoleID and get the Role entry
	role, roleName, err := b.validateRoleID(req.Storage, roleID)
	if err != nil {
		return nil, "", metadata, err
	}
	if role == nil || roleName == "" {
		return nil, "", metadata, fmt.Errorf("failed to validate role_id")
	}

	// Calculate the TTL boundaries since this reflects the properties of the token issued
	if role.TokenTTL, role.TokenMaxTTL, err = b.SanitizeTTL(role.TokenTTL, role.TokenMaxTTL); err != nil {
		return nil, "", metadata, err
	}

	if role.BindSecretID {
		// If 'bind_secret_id' was set on role, look for the field 'secret_id'
		// to be specified and validate it.
		secretID := strings.TrimSpace(data.Get("secret_id").(string))
		if secretID == "" {
			return nil, "", metadata, fmt.Errorf("missing secret_id")
		}

		// Check if the SecretID supplied is valid. If use limit was specified
		// on the SecretID, it will be decremented in this call.
		var valid bool
		valid, metadata, err = b.validateBindSecretID(req.Storage, roleName, secretID, role.HMACKey)
		if err != nil {
			return nil, "", metadata, err
		}
		if !valid {
			return nil, "", metadata, fmt.Errorf("invalid secret_id: %s\n", secretID)
		}
	}

	if role.BoundCIDRList != "" {
		// If 'bound_cidr_list' was set, verify the CIDR restrictions
		cidrBlocks := strings.Split(role.BoundCIDRList, ",")
		for _, block := range cidrBlocks {
			_, cidr, err := net.ParseCIDR(block)
			if err != nil {
				return nil, "", metadata, fmt.Errorf("invalid cidr: %s", err)
			}

			var addr string
			if req.Connection != nil {
				addr = req.Connection.RemoteAddr
			}
			if addr == "" || !cidr.Contains(net.ParseIP(addr)) {
				return nil, "", metadata, fmt.Errorf("unauthorized source address")
			}
		}
	}

	return role, roleName, metadata, nil
}

// validateBindSecretID is used to determine if the given SecretID is a valid one.
func (b *backend) validateBindSecretID(s logical.Storage, roleName, secretID, hmacKey string) (bool, map[string]string, error) {
	secretIDHMAC, err := createHMAC(hmacKey, secretID)
	if err != nil {
		return false, nil, fmt.Errorf("failed to create HMAC of secret_id: %s", err)
	}

	roleNameHMAC, err := createHMAC(hmacKey, roleName)
	if err != nil {
		return false, nil, fmt.Errorf("failed to create HMAC of role_name: %s", err)
	}

	entryIndex := fmt.Sprintf("secret_id/%s/%s", roleNameHMAC, secretIDHMAC)

	// SecretID locks are always index based on secretIDHMACs. This helps
	// acquiring the locks when the SecretIDs are listed. This allows grabbing
	// the correct locks even if the SecretIDs are not known in plaintext.
	lock := b.secretIDLock(secretIDHMAC)
	lock.RLock()

	result := secretIDStorageEntry{}
	if entry, err := s.Get(entryIndex); err != nil {
		lock.RUnlock()
		return false, nil, err
	} else if entry == nil {
		lock.RUnlock()
		return false, nil, nil
	} else if err := entry.DecodeJSON(&result); err != nil {
		lock.RUnlock()
		return false, nil, err
	}

	// SecretIDNumUses will be zero only if the usage limit was not set at all,
	// in which case, the SecretID will remain to be valid as long as it is not
	// expired.
	if result.SecretIDNumUses == 0 {
		lock.RUnlock()
		return true, result.Metadata, nil
	}

	// If the SecretIDNumUses is non-zero, it means that its use-count should be updated
	// in the storage. Switch the lock from a `read` to a `write` and update
	// the storage entry.
	lock.RUnlock()

	lock.Lock()
	defer lock.Unlock()

	// Lock switching may change the data. Refresh the contents.
	result = secretIDStorageEntry{}
	if entry, err := s.Get(entryIndex); err != nil {
		return false, nil, err
	} else if entry == nil {
		return false, nil, nil
	} else if err := entry.DecodeJSON(&result); err != nil {
		return false, nil, err
	}

	// If there exists a single use left, delete the SecretID entry from
	// the storage but do not fail the validation request. Subsequest
	// requests to use the same SecretID will fail.
	if result.SecretIDNumUses == 1 {
		accessorEntryIndex := "accessor/" + b.salt.SaltID(result.SecretIDAccessor)
		if err := s.Delete(accessorEntryIndex); err != nil {
			return false, nil, fmt.Errorf("failed to delete accessor storage entry: %s", err)
		}
		if err := s.Delete(entryIndex); err != nil {
			return false, nil, fmt.Errorf("failed to delete SecretID: %s", err)
		}
	} else {
		// If the use count is greater than one, decrement it and update the last updated time.
		result.SecretIDNumUses -= 1
		result.LastUpdatedTime = time.Now()
		if entry, err := logical.StorageEntryJSON(entryIndex, &result); err != nil {
			return false, nil, fmt.Errorf("failed to decrement the use count for SecretID:%s", secretID)
		} else if err = s.Put(entry); err != nil {
			return false, nil, fmt.Errorf("failed to decrement the use count for SecretID:%s", secretID)
		}
	}

	return true, result.Metadata, nil
}

// Creates a SHA256 HMAC of the given 'value' using the given 'key'
// and returns a hex encoded string.
func createHMAC(key, value string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("invalid HMAC key")
	}
	hm := hmac.New(sha256.New, []byte(key))
	hm.Write([]byte(value))
	return hex.EncodeToString(hm.Sum(nil)), nil
}

// secretIDLock is used to get a lock from the pre-initialized map
// of locks. Map is indexed based on the first 2 characters of the
// secretIDHMAC. If the input is not hex encoded or if empty, a
// "custom" lock will be returned.
func (b *backend) secretIDLock(secretIDHMAC string) *sync.RWMutex {
	var lock *sync.RWMutex
	var ok bool
	if len(secretIDHMAC) >= 2 {
		lock, ok = b.secretIDLocksMap[secretIDHMAC[0:2]]
	}
	if !ok || lock == nil {
		// Fall back for custom SecretIDs
		lock = b.secretIDLocksMap["custom"]
	}
	return lock
}

// registerSecretIDEntry creates a new storage entry for the given SecretID.
func (b *backend) registerSecretIDEntry(s logical.Storage, roleName, secretID, hmacKey string, secretEntry *secretIDStorageEntry) (*secretIDStorageEntry, error) {
	secretIDHMAC, err := createHMAC(hmacKey, secretID)
	if err != nil {
		return nil, fmt.Errorf("failed to create HMAC of secret_id: %s", err)
	}
	roleNameHMAC, err := createHMAC(hmacKey, roleName)
	if err != nil {
		return nil, fmt.Errorf("failed to create HMAC of role_name: %s", err)
	}

	entryIndex := fmt.Sprintf("secret_id/%s/%s", roleNameHMAC, secretIDHMAC)

	lock := b.secretIDLock(secretIDHMAC)
	lock.RLock()

	entry, err := s.Get(entryIndex)
	if err != nil {
		lock.RUnlock()
		return nil, err
	}
	if entry != nil {
		lock.RUnlock()
		return nil, fmt.Errorf("SecretID is already registered")
	}

	// If there isn't an entry for the secretID already, switch the read lock
	// with a write lock and create an entry.
	lock.RUnlock()
	lock.Lock()
	defer lock.Unlock()

	// But before saving a new entry, check if the secretID entry was created during the lock switch.
	entry, err = s.Get(entryIndex)
	if err != nil {
		return nil, err
	}
	if entry != nil {
		return nil, fmt.Errorf("SecretID is already registered")
	}

	// Create a new entry for the SecretID

	// Set the creation time for the SecretID
	currentTime := time.Now()
	secretEntry.CreationTime = currentTime
	secretEntry.LastUpdatedTime = currentTime

	// If SecretIDTTL is not specified or if it crosses the backend mount's limit,
	// cap the expiration to backend's max. Otherwise, use it to determine the
	// expiration time.
	if secretEntry.SecretIDTTL < time.Duration(0) || secretEntry.SecretIDTTL > b.System().MaxLeaseTTL() {
		secretEntry.ExpirationTime = currentTime.Add(b.System().MaxLeaseTTL())
	} else if secretEntry.SecretIDTTL != time.Duration(0) {
		// Set the ExpirationTime only if SecretIDTTL was set. SecretIDs should not
		// expire by default.
		secretEntry.ExpirationTime = currentTime.Add(secretEntry.SecretIDTTL)
	}

	// Before storing the SecretID, store its accessor.
	if err := b.createAccessor(s, secretEntry, secretIDHMAC); err != nil {
		return nil, err
	}

	if entry, err := logical.StorageEntryJSON(entryIndex, secretEntry); err != nil {
		return nil, err
	} else if err = s.Put(entry); err != nil {
		return nil, err
	}

	return secretEntry, nil
}

// secretIDAccessorEntry is used to read the storage entry that maps an
// accessor to a secret_id. This method should be called when the lock
// for the corresponding SecretID is held.
func (b *backend) secretIDAccessorEntry(s logical.Storage, secretIDAccessor string) (*secretIDAccessorStorageEntry, error) {
	if secretIDAccessor == "" {
		return nil, fmt.Errorf("missing secretIDAccessor")
	}

	var result secretIDAccessorStorageEntry

	// Create index entry, mapping the accessor to the token ID
	entryIndex := "accessor/" + b.salt.SaltID(secretIDAccessor)

	if entry, err := s.Get(entryIndex); err != nil {
		return nil, err
	} else if entry == nil {
		return nil, nil
	} else if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// createAccessor creates an identifier for the SecretID. A storage index,
// mapping the accessor to the SecretID is also created. This method should
// be called when the lock for the corresponding SecretID is held.
func (b *backend) createAccessor(s logical.Storage, entry *secretIDStorageEntry, secretIDHMAC string) error {
	// Create a random accessor
	accessorUUID, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	entry.SecretIDAccessor = accessorUUID

	// Create index entry, mapping the accessor to the token ID
	entryIndex := "accessor/" + b.salt.SaltID(entry.SecretIDAccessor)
	if entry, err := logical.StorageEntryJSON(entryIndex, &secretIDAccessorStorageEntry{
		SecretIDHMAC: secretIDHMAC,
	}); err != nil {
		return err
	} else if err = s.Put(entry); err != nil {
		return fmt.Errorf("failed to persist accessor index entry: %s", err)
	}

	return nil
}

// flushRoleSecrets deletes all the SecretIDs that belong to the given
// RoleID.
func (b *backend) flushRoleSecrets(s logical.Storage, roleName, hmacKey string) error {
	roleNameHMAC, err := createHMAC(hmacKey, roleName)
	if err != nil {
		return fmt.Errorf("failed to create HMAC of role_name: %s", err)
	}

	// Acquire the custom lock to perform listing of SecretIDs
	customLock := b.secretIDLock("")
	customLock.RLock()
	defer customLock.RUnlock()

	secretIDHMACs, err := s.List(fmt.Sprintf("secret_id/%s/", roleNameHMAC))
	if err != nil {
		return err
	}
	for _, secretIDHMAC := range secretIDHMACs {
		// Acquire the lock belonging to the SecretID
		lock := b.secretIDLock(secretIDHMAC)
		lock.Lock()
		entryIndex := fmt.Sprintf("secret_id/%s/%s", roleNameHMAC, secretIDHMAC)
		if err := s.Delete(entryIndex); err != nil {
			lock.Unlock()
			return fmt.Errorf("error deleting SecretID %s from storage: %s", secretIDHMAC, err)
		}
		lock.Unlock()
	}
	return nil
}
