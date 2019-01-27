package keysutil

import (
	"context"
	"encoding/base64"
	"errors"
	"math/big"
	paths "path"
	"sort"
	"strings"

	lru "github.com/hashicorp/golang-lru"
	"github.com/hashicorp/vault/logical"
)

const (
	// DefaultCacheSize is used if no cache size is specified for
	// NewEncryptedKeyStorage. This value is the number of cache entries to
	// store, not the size in bytes of the cache.
	DefaultCacheSize = 16 * 1024

	// DefaultPrefix is used if no prefix is specified for
	// NewEncryptedKeyStorage. Prefix must be defined so we can provide context
	// for the base folder.
	DefaultPrefix = "encryptedkeys/"

	// EncryptedKeyPolicyVersionTpl is a template that can be used to minimize
	// the amount of data that's stored with the ciphertext.
	EncryptedKeyPolicyVersionTpl = "{{version}}:"
)

var (
	// ErrPolicyDerivedKeys is returned if the provided policy does not use
	// derived keys. This is a requirement for this storage implementation.
	ErrPolicyDerivedKeys = errors.New("key policy must use derived keys")

	// ErrPolicyConvergentEncryption is returned if the provided policy does not use
	// convergent encryption. This is a requirement for this storage implementation.
	ErrPolicyConvergentEncryption = errors.New("key policy must use convergent encryption")

	// ErrPolicyConvergentVersion is returned if the provided policy does not use
	// a new enough convergent version. This is a requirement for this storage
	// implementation.
	ErrPolicyConvergentVersion = errors.New("key policy must use convergent version > 2")

	// ErrNilStorage is returned if the provided storage is nil.
	ErrNilStorage = errors.New("nil storage provided")

	// ErrNilPolicy is returned if the provided policy is nil.
	ErrNilPolicy = errors.New("nil policy provided")
)

// EncryptedKeyStorageConfig is used to configure an EncryptedKeyStorage object.
type EncryptedKeyStorageConfig struct {
	// Policy is the key policy to use to encrypt the key paths.
	Policy *Policy

	// Prefix is the storage prefix for this instance of the EncryptedKeyStorage
	// object. This is stored in plaintext. If not set the DefaultPrefix will be
	// used.
	Prefix string

	// CacheSize is the number of elements to cache. If not set the
	// DetaultCacheSize will be used.
	CacheSize int
}

// NewEncryptedKeyStorageWrapper takes an EncryptedKeyStorageConfig and returns a new
// EncryptedKeyStorage object.
func NewEncryptedKeyStorageWrapper(config EncryptedKeyStorageConfig) (*EncryptedKeyStorageWrapper, error) {
	if config.Policy == nil {
		return nil, ErrNilPolicy
	}

	if !config.Policy.Derived {
		return nil, ErrPolicyDerivedKeys
	}

	if !config.Policy.ConvergentEncryption {
		return nil, ErrPolicyConvergentEncryption
	}

	if config.Prefix == "" {
		config.Prefix = DefaultPrefix
	}

	if !strings.HasSuffix(config.Prefix, "/") {
		config.Prefix += "/"
	}

	size := config.CacheSize
	if size <= 0 {
		size = DefaultCacheSize
	}

	cache, err := lru.New2Q(size)
	if err != nil {
		return nil, err
	}

	return &EncryptedKeyStorageWrapper{
		policy: config.Policy,
		prefix: config.Prefix,
		lru:    cache,
	}, nil
}

type EncryptedKeyStorageWrapper struct {
	policy *Policy
	lru    *lru.TwoQueueCache
	prefix string
}

func (f *EncryptedKeyStorageWrapper) Wrap(s logical.Storage) logical.Storage {
	return &encryptedKeyStorage{
		policy: f.policy,
		s:      s,
		prefix: f.prefix,
		lru:    f.lru,
	}
}

// EncryptedKeyStorage implements the logical.Storage interface and ensures the
// storage paths are encrypted in the underlying storage.
type encryptedKeyStorage struct {
	policy *Policy
	s      logical.Storage
	lru    *lru.TwoQueueCache

	prefix string
}

func ensureTailingSlash(path string) string {
	if !strings.HasSuffix(path, "/") {
		return path + "/"
	}
	return path
}

// List implements the logical.Storage List method, and decrypts all the items
// in a path prefix. This can only operate on full folder structures so the
// prefix should end in a "/".
func (s *encryptedKeyStorage) List(ctx context.Context, prefix string) ([]string, error) {
	var decoder big.Int

	encPrefix, err := s.encryptPath(prefix)
	if err != nil {
		return nil, err
	}

	keys, err := s.s.List(ctx, ensureTailingSlash(encPrefix))
	if err != nil {
		return keys, err
	}

	decryptedKeys := make([]string, len(keys))

	// The context for the decryption operations will be the object's prefix
	// joined with the provided prefix. Join cleans the path ensuring there
	// isn't a trailing "/".
	context := []byte(paths.Join(s.prefix, prefix))

	for i, k := range keys {
		raw, ok := s.lru.Get(k)
		if ok {
			// cache HIT, we can bail early and skip the decode & decrypt operations.
			decryptedKeys[i] = raw.(string)
			continue
		}

		// If a folder is included in the keys it will have a trailing "/".
		// We need to remove this before decoding/decrypting and add it back
		// later.
		appendSlash := strings.HasSuffix(k, "/")
		if appendSlash {
			k = strings.TrimSuffix(k, "/")
		}

		decoder.SetString(k, 62)
		decoded := decoder.Bytes()
		if len(decoded) == 0 {
			return nil, errors.New("could not decode key")
		}

		// Decrypt the data with the object's key policy.
		encodedPlaintext, err := s.policy.Decrypt(context, nil, string(decoded[:]))
		if err != nil {
			return nil, err
		}

		// The plaintext is still base64 encoded, decode it.
		decoded, err = base64.StdEncoding.DecodeString(encodedPlaintext)
		if err != nil {
			return nil, err
		}

		plaintext := string(decoded[:])

		// Add the slash back to the plaintext value
		if appendSlash {
			plaintext += "/"
			k += "/"
		}

		// We want to store the unencoded version of the key in the cache.
		// This will make it more performent when it's a HIT.
		s.lru.Add(k, plaintext)

		decryptedKeys[i] = plaintext
	}

	sort.Strings(decryptedKeys)
	return decryptedKeys, nil
}

// Get implements the logical.Storage Get method.
func (s *encryptedKeyStorage) Get(ctx context.Context, path string) (*logical.StorageEntry, error) {
	encPath, err := s.encryptPath(path)
	if err != nil {
		return nil, err
	}

	return s.s.Get(ctx, encPath)
}

// Put implements the logical.Storage Put method.
func (s *encryptedKeyStorage) Put(ctx context.Context, entry *logical.StorageEntry) error {
	encPath, err := s.encryptPath(entry.Key)
	if err != nil {
		return err
	}
	e := &logical.StorageEntry{}
	*e = *entry

	e.Key = encPath

	return s.s.Put(ctx, e)
}

// Delete implements the logical.Storage Delete method.
func (s *encryptedKeyStorage) Delete(ctx context.Context, path string) error {
	encPath, err := s.encryptPath(path)
	if err != nil {
		return err
	}

	return s.s.Delete(ctx, encPath)
}

// encryptPath takes a plaintext path and encrypts each path section (separated
// by "/") with the object's key policy. The context for each encryption is the
// plaintext path prefix for the key.
func (s *encryptedKeyStorage) encryptPath(path string) (string, error) {
	var encoder big.Int

	if path == "" || path == "/" {
		return s.prefix, nil
	}

	path = paths.Clean(path)

	// Trim the prefix if it starts with a "/"
	path = strings.TrimPrefix(path, "/")

	parts := strings.Split(path, "/")

	encPath := s.prefix
	context := strings.TrimSuffix(s.prefix, "/")
	for _, p := range parts {
		encoded := base64.StdEncoding.EncodeToString([]byte(p))
		ciphertext, err := s.policy.Encrypt(0, []byte(context), nil, encoded)
		if err != nil {
			return "", err
		}

		encoder.SetBytes([]byte(ciphertext))
		encPath = paths.Join(encPath, encoder.Text(62))
		context = paths.Join(context, p)
	}

	return encPath, nil
}
