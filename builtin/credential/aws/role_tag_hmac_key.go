package aws

import (
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/logical"
)

const (
	RoleTagHMACKeyLocation = "role_tag_hmac_key"
)

// hmacKey returns the key to HMAC the RoleTag value. The key is valid per backend mount.
// If a key is not created for the mount, a new key will be created.
func hmacKey(s logical.Storage) (string, error) {
	raw, err := s.Get(RoleTagHMACKeyLocation)
	if err != nil {
		return "", fmt.Errorf("failed to read key: %v", err)
	}

	key := ""
	if raw != nil {
		key = string(raw.Value)
	}

	if key == "" {
		key, err = uuid.GenerateUUID()
		if err != nil {
			return "", fmt.Errorf("failed to generate uuid: %v", err)
		}
		if s != nil {
			entry := &logical.StorageEntry{
				Key:   RoleTagHMACKeyLocation,
				Value: []byte(key),
			}
			if err := s.Put(entry); err != nil {
				return "", fmt.Errorf("failed to persist key: %v", err)
			}
		}
	}

	return key, nil
}
