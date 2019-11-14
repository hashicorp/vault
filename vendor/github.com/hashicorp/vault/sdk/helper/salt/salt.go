package salt

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"

	"github.com/hashicorp/errwrap"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	// DefaultLocation is the path in the view we store our key salt
	// if no other path is provided.
	DefaultLocation = "salt"
)

// Salt is used to manage a persistent salt key which is used to
// hash values. This allows keys to be generated and recovered
// using the global salt. Primarily, this allows paths in the storage
// backend to be obfuscated if they may contain sensitive information.
type Salt struct {
	config    *Config
	salt      string
	generated bool
}

type HashFunc func([]byte) []byte

// Config is used to parameterize the Salt
type Config struct {
	// Location is the path in the storage backend for the
	// salt. Uses DefaultLocation if not specified.
	Location string

	// HashFunc is the hashing function to use for salting.
	// Defaults to SHA1 if not provided.
	HashFunc HashFunc

	// HMAC allows specification of a hash function to use for
	// the HMAC helpers
	HMAC func() hash.Hash

	// String prepended to HMAC strings for identification.
	// Required if using HMAC
	HMACType string
}

// NewSalt creates a new salt based on the configuration
func NewSalt(ctx context.Context, view logical.Storage, config *Config) (*Salt, error) {
	// Setup the configuration
	if config == nil {
		config = &Config{}
	}
	if config.Location == "" {
		config.Location = DefaultLocation
	}
	if config.HashFunc == nil {
		config.HashFunc = SHA256Hash
	}
	if config.HMAC == nil {
		config.HMAC = sha256.New
		config.HMACType = "hmac-sha256"
	}

	// Create the salt
	s := &Salt{
		config: config,
	}

	// Look for the salt
	var raw *logical.StorageEntry
	var err error
	if view != nil {
		raw, err = view.Get(ctx, config.Location)
		if err != nil {
			return nil, errwrap.Wrapf("failed to read salt: {{err}}", err)
		}
	}

	// Restore the salt if it exists
	if raw != nil {
		s.salt = string(raw.Value)
	}

	// Generate a new salt if necessary
	if s.salt == "" {
		s.salt, err = uuid.GenerateUUID()
		if err != nil {
			return nil, errwrap.Wrapf("failed to generate uuid: {{err}}", err)
		}
		s.generated = true
		if view != nil {
			raw := &logical.StorageEntry{
				Key:   config.Location,
				Value: []byte(s.salt),
			}
			if err := view.Put(ctx, raw); err != nil {
				return nil, errwrap.Wrapf("failed to persist salt: {{err}}", err)
			}
		}
	}

	if config.HMAC != nil {
		if len(config.HMACType) == 0 {
			return nil, fmt.Errorf("HMACType must be defined")
		}
	}

	return s, nil
}

// SaltID is used to apply a salt and hash function to an ID to make sure
// it is not reversible
func (s *Salt) SaltID(id string) string {
	return SaltID(s.salt, id, s.config.HashFunc)
}

// GetHMAC is used to apply a salt and hash function to data to make sure it is
// not reversible, with an additional HMAC
func (s *Salt) GetHMAC(data string) string {
	hm := hmac.New(s.config.HMAC, []byte(s.salt))
	hm.Write([]byte(data))
	return hex.EncodeToString(hm.Sum(nil))
}

// GetIdentifiedHMAC is used to apply a salt and hash function to data to make
// sure it is not reversible, with an additional HMAC, and ID prepended
func (s *Salt) GetIdentifiedHMAC(data string) string {
	return s.config.HMACType + ":" + s.GetHMAC(data)
}

// DidGenerate returns true if the underlying salt value was generated
// on initialization.
func (s *Salt) DidGenerate() bool {
	return s.generated
}

// SaltIDHashFunc uses the supplied hash function instead of the configured
// hash func in the salt.
func (s *Salt) SaltIDHashFunc(id string, hashFunc HashFunc) string {
	return SaltID(s.salt, id, hashFunc)
}

// SaltID is used to apply a salt and hash function to an ID to make sure
// it is not reversible
func SaltID(salt, id string, hash HashFunc) string {
	comb := salt + id
	hashVal := hash([]byte(comb))
	return hex.EncodeToString(hashVal)
}

func HMACValue(salt, val string, hashFunc func() hash.Hash) string {
	hm := hmac.New(hashFunc, []byte(salt))
	hm.Write([]byte(val))
	return hex.EncodeToString(hm.Sum(nil))
}

func HMACIdentifiedValue(salt, val, hmacType string, hashFunc func() hash.Hash) string {
	return hmacType + ":" + HMACValue(salt, val, hashFunc)
}

// SHA1Hash returns the SHA1 of the input
func SHA1Hash(inp []byte) []byte {
	hashed := sha1.Sum(inp)
	return hashed[:]
}

// SHA256Hash returns the SHA256 of the input
func SHA256Hash(inp []byte) []byte {
	hashed := sha256.Sum256(inp)
	return hashed[:]
}
