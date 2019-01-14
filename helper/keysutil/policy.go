package keysutil

import (
	"bytes"
	"context"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"path"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/hkdf"

	"github.com/hashicorp/errwrap"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/kdf"
	"github.com/hashicorp/vault/logical"
)

// Careful with iota; don't put anything before it in this const block because
// we need the default of zero to be the old-style KDF
const (
	Kdf_hmac_sha256_counter = iota // built-in helper
	Kdf_hkdf_sha256                // golang.org/x/crypto/hkdf
)

// Or this one...we need the default of zero to be the original AES256-GCM96
const (
	KeyType_AES256_GCM96 = iota
	KeyType_ECDSA_P256
	KeyType_ED25519
	KeyType_RSA2048
	KeyType_RSA4096
	KeyType_ChaCha20_Poly1305
)

const (
	// ErrTooOld is returned whtn the ciphertext or signatures's key version is
	// too old.
	ErrTooOld = "ciphertext or signature version is disallowed by policy (too old)"

	// DefaultVersionTemplate is used when no version template is provided.
	DefaultVersionTemplate = "vault:v{{version}}:"
)

type RestoreInfo struct {
	Time    time.Time `json:"time"`
	Version int       `json:"version"`
}

type BackupInfo struct {
	Time    time.Time `json:"time"`
	Version int       `json:"version"`
}

type SigningResult struct {
	Signature string
	PublicKey []byte
}

type ecdsaSignature struct {
	R, S *big.Int
}

type KeyType int

func (kt KeyType) EncryptionSupported() bool {
	switch kt {
	case KeyType_AES256_GCM96, KeyType_ChaCha20_Poly1305, KeyType_RSA2048, KeyType_RSA4096:
		return true
	}
	return false
}

func (kt KeyType) DecryptionSupported() bool {
	switch kt {
	case KeyType_AES256_GCM96, KeyType_ChaCha20_Poly1305, KeyType_RSA2048, KeyType_RSA4096:
		return true
	}
	return false
}

func (kt KeyType) SigningSupported() bool {
	switch kt {
	case KeyType_ECDSA_P256, KeyType_ED25519, KeyType_RSA2048, KeyType_RSA4096:
		return true
	}
	return false
}

func (kt KeyType) HashSignatureInput() bool {
	switch kt {
	case KeyType_ECDSA_P256, KeyType_RSA2048, KeyType_RSA4096:
		return true
	}
	return false
}

func (kt KeyType) DerivationSupported() bool {
	switch kt {
	case KeyType_AES256_GCM96, KeyType_ChaCha20_Poly1305, KeyType_ED25519:
		return true
	}
	return false
}

func (kt KeyType) String() string {
	switch kt {
	case KeyType_AES256_GCM96:
		return "aes256-gcm96"
	case KeyType_ChaCha20_Poly1305:
		return "chacha20-poly1305"
	case KeyType_ECDSA_P256:
		return "ecdsa-p256"
	case KeyType_ED25519:
		return "ed25519"
	case KeyType_RSA2048:
		return "rsa-2048"
	case KeyType_RSA4096:
		return "rsa-4096"
	}

	return "[unknown]"
}

type KeyData struct {
	Policy       *Policy       `json:"policy"`
	ArchivedKeys *archivedKeys `json:"archived_keys"`
}

// KeyEntry stores the key and metadata
type KeyEntry struct {
	// AES or some other kind that is a pure byte slice like ED25519
	Key []byte `json:"key"`

	// Key used for HMAC functions
	HMACKey []byte `json:"hmac_key"`

	// Time of creation
	CreationTime time.Time `json:"time"`

	EC_X *big.Int `json:"ec_x"`
	EC_Y *big.Int `json:"ec_y"`
	EC_D *big.Int `json:"ec_d"`

	RSAKey *rsa.PrivateKey `json:"rsa_key"`

	// The public key in an appropriate format for the type of key
	FormattedPublicKey string `json:"public_key"`

	// If convergent is enabled, the version (falling back to what's in the
	// policy)
	ConvergentVersion int `json:"convergent_version"`

	// This is deprecated (but still filled) in favor of the value above which
	// is more precise
	DeprecatedCreationTime int64 `json:"creation_time"`
}

// deprecatedKeyEntryMap is used to allow JSON marshal/unmarshal
type deprecatedKeyEntryMap map[int]KeyEntry

// MarshalJSON implements JSON marshaling
func (kem deprecatedKeyEntryMap) MarshalJSON() ([]byte, error) {
	intermediate := map[string]KeyEntry{}
	for k, v := range kem {
		intermediate[strconv.Itoa(k)] = v
	}
	return json.Marshal(&intermediate)
}

// MarshalJSON implements JSON unmarshalling
func (kem deprecatedKeyEntryMap) UnmarshalJSON(data []byte) error {
	intermediate := map[string]KeyEntry{}
	if err := jsonutil.DecodeJSON(data, &intermediate); err != nil {
		return err
	}
	for k, v := range intermediate {
		keyval, err := strconv.Atoi(k)
		if err != nil {
			return err
		}
		kem[keyval] = v
	}

	return nil
}

// keyEntryMap is used to allow JSON marshal/unmarshal
type keyEntryMap map[string]KeyEntry

// PolicyConfig is used to create a new policy
type PolicyConfig struct {
	// The name of the policy
	Name string `json:"name"`

	// The type of key
	Type KeyType

	// Derived keys MUST provide a context and the master underlying key is
	// never used.
	Derived              bool
	KDF                  int
	ConvergentEncryption bool

	// Whether the key is exportable
	Exportable bool

	// Whether the key is allowed to be deleted
	DeletionAllowed bool

	// AllowPlaintextBackup allows taking backup of the policy in plaintext
	AllowPlaintextBackup bool

	// VersionTemplate is used to prefix the ciphertext with information about
	// the key version. It must inclide {{version}} and a delimiter between the
	// version prefix and the ciphertext.
	VersionTemplate string

	// StoragePrefix is used to add a prefix when storing and retrieving the
	// policy object.
	StoragePrefix string
}

// NewPolicy takes a policy config and returns a Policy with those settings.
func NewPolicy(config PolicyConfig) *Policy {
	return &Policy{
		l:                    new(sync.RWMutex),
		Name:                 config.Name,
		Type:                 config.Type,
		Derived:              config.Derived,
		KDF:                  config.KDF,
		ConvergentEncryption: config.ConvergentEncryption,
		ConvergentVersion:    -1,
		Exportable:           config.Exportable,
		DeletionAllowed:      config.DeletionAllowed,
		AllowPlaintextBackup: config.AllowPlaintextBackup,
		VersionTemplate:      config.VersionTemplate,
		StoragePrefix:        config.StoragePrefix,
	}
}

// LoadPolicy will load a policy from the provided storage path and set the
// necessary un-exported variables. It is particularly useful when accessing a
// policy without the lock manager.
func LoadPolicy(ctx context.Context, s logical.Storage, path string) (*Policy, error) {
	raw, err := s.Get(ctx, path)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}

	var policy Policy
	err = jsonutil.DecodeJSON(raw.Value, &policy)
	if err != nil {
		return nil, err
	}

	policy.l = new(sync.RWMutex)

	return &policy, nil
}

// Policy is the struct used to store metadata
type Policy struct {
	// This is a pointer on purpose: if we are running with cache disabled we
	// need to actually swap in the lock manager's lock for this policy with
	// the local lock.
	l *sync.RWMutex
	// writeLocked allows us to implement Lock() and Unlock()
	writeLocked bool
	// Stores whether it's been deleted. This acts as a guard for operations
	// that may write data, e.g. if one request rotates and that request is
	// served after a delete.
	deleted uint32

	Name string      `json:"name"`
	Key  []byte      `json:"key,omitempty"` //DEPRECATED
	Keys keyEntryMap `json:"keys"`

	// Derived keys MUST provide a context and the master underlying key is
	// never used. If convergent encryption is true, the context will be used
	// as the nonce as well.
	Derived              bool `json:"derived"`
	KDF                  int  `json:"kdf"`
	ConvergentEncryption bool `json:"convergent_encryption"`

	// Whether the key is exportable
	Exportable bool `json:"exportable"`

	// The minimum version of the key allowed to be used for decryption
	MinDecryptionVersion int `json:"min_decryption_version"`

	// The minimum version of the key allowed to be used for encryption
	MinEncryptionVersion int `json:"min_encryption_version"`

	// The latest key version in this policy
	LatestVersion int `json:"latest_version"`

	// The latest key version in the archive. We never delete these, so this is
	// a max.
	ArchiveVersion int `json:"archive_version"`

	// ArchiveMinVersion is the minimum version of the key in the archive.
	ArchiveMinVersion int `json:"archive_min_version"`

	// MinAvailableVersion is the minimum version of the key present. All key
	// versions before this would have been deleted.
	MinAvailableVersion int `json:"min_available_version"`

	// Whether the key is allowed to be deleted
	DeletionAllowed bool `json:"deletion_allowed"`

	// The version of the convergent nonce to use
	ConvergentVersion int `json:"convergent_version"`

	// The type of key
	Type KeyType `json:"type"`

	// BackupInfo indicates the information about the backup action taken on
	// this policy
	BackupInfo *BackupInfo `json:"backup_info"`

	// RestoreInfo indicates the information about the restore action taken on
	// this policy
	RestoreInfo *RestoreInfo `json:"restore_info"`

	// AllowPlaintextBackup allows taking backup of the policy in plaintext
	AllowPlaintextBackup bool `json:"allow_plaintext_backup"`

	// VersionTemplate is used to prefix the ciphertext with information about
	// the key version. It must inclide {{version}} and a delimiter between the
	// version prefix and the ciphertext.
	VersionTemplate string `json:"version_template"`

	// StoragePrefix is used to add a prefix when storing and retrieving the
	// policy object.
	StoragePrefix string `json:"storage_prefix"`

	// versionPrefixCache stores caches of version prefix strings and the split
	// version template.
	versionPrefixCache sync.Map
}

func (p *Policy) Lock(exclusive bool) {
	if exclusive {
		p.l.Lock()
		p.writeLocked = true
	} else {
		p.l.RLock()
	}
}

func (p *Policy) Unlock() {
	if p.writeLocked {
		p.writeLocked = false
		p.l.Unlock()
	} else {
		p.l.RUnlock()
	}
}

// ArchivedKeys stores old keys. This is used to keep the key loading time sane
// when there are huge numbers of rotations.
type archivedKeys struct {
	Keys []KeyEntry `json:"keys"`
}

func (p *Policy) LoadArchive(ctx context.Context, storage logical.Storage) (*archivedKeys, error) {
	archive := &archivedKeys{}

	raw, err := storage.Get(ctx, path.Join(p.StoragePrefix, "archive", p.Name))
	if err != nil {
		return nil, err
	}
	if raw == nil {
		archive.Keys = make([]KeyEntry, 0)
		return archive, nil
	}

	if err := jsonutil.DecodeJSON(raw.Value, archive); err != nil {
		return nil, err
	}

	return archive, nil
}

func (p *Policy) storeArchive(ctx context.Context, storage logical.Storage, archive *archivedKeys) error {
	// Encode the policy
	buf, err := json.Marshal(archive)
	if err != nil {
		return err
	}

	// Write the policy into storage
	err = storage.Put(ctx, &logical.StorageEntry{
		Key:   path.Join(p.StoragePrefix, "archive", p.Name),
		Value: buf,
	})
	if err != nil {
		return err
	}

	return nil
}

// handleArchiving manages the movement of keys to and from the policy archive.
// This should *ONLY* be called from Persist() since it assumes that the policy
// will be persisted afterwards.
func (p *Policy) handleArchiving(ctx context.Context, storage logical.Storage) error {
	// We need to move keys that are no longer accessible to archivedKeys, and keys
	// that now need to be accessible back here.
	//
	// For safety, because there isn't really a good reason to, we never delete
	// keys from the archive even when we move them back.

	// Check if we have the latest minimum version in the current set of keys
	_, keysContainsMinimum := p.Keys[strconv.Itoa(p.MinDecryptionVersion)]

	// Sanity checks
	switch {
	case p.MinDecryptionVersion < 1:
		return fmt.Errorf("minimum decryption version of %d is less than 1", p.MinDecryptionVersion)
	case p.LatestVersion < 1:
		return fmt.Errorf("latest version of %d is less than 1", p.LatestVersion)
	case !keysContainsMinimum && p.ArchiveVersion != p.LatestVersion:
		return fmt.Errorf("need to move keys from archive but archive version not up-to-date")
	case p.ArchiveVersion > p.LatestVersion:
		return fmt.Errorf("archive version of %d is greater than the latest version %d",
			p.ArchiveVersion, p.LatestVersion)
	case p.MinEncryptionVersion > 0 && p.MinEncryptionVersion < p.MinDecryptionVersion:
		return fmt.Errorf("minimum decryption version of %d is greater than minimum encryption version %d",
			p.MinDecryptionVersion, p.MinEncryptionVersion)
	case p.MinDecryptionVersion > p.LatestVersion:
		return fmt.Errorf("minimum decryption version of %d is greater than the latest version %d",
			p.MinDecryptionVersion, p.LatestVersion)
	}

	archive, err := p.LoadArchive(ctx, storage)
	if err != nil {
		return err
	}

	if !keysContainsMinimum {
		// Need to move keys *from* archive
		for i := p.MinDecryptionVersion; i <= p.LatestVersion; i++ {
			p.Keys[strconv.Itoa(i)] = archive.Keys[i-p.MinAvailableVersion]
		}

		return nil
	}

	// Need to move keys *to* archive

	// We need a size that is equivalent to the latest version (number of keys)
	// but adding one since slice numbering starts at 0 and we're indexing by
	// key version
	if len(archive.Keys)+p.MinAvailableVersion < p.LatestVersion+1 {
		// Increase the size of the archive slice
		newKeys := make([]KeyEntry, p.LatestVersion-p.MinAvailableVersion+1)
		copy(newKeys, archive.Keys)
		archive.Keys = newKeys
	}

	// We are storing all keys in the archive, so we ensure that it is up to
	// date up to p.LatestVersion
	for i := p.ArchiveVersion + 1; i <= p.LatestVersion; i++ {
		archive.Keys[i-p.MinAvailableVersion] = p.Keys[strconv.Itoa(i)]
		p.ArchiveVersion = i
	}

	// Trim the keys if required
	if p.ArchiveMinVersion < p.MinAvailableVersion {
		archive.Keys = archive.Keys[p.MinAvailableVersion-p.ArchiveMinVersion:]
		p.ArchiveMinVersion = p.MinAvailableVersion
	}

	err = p.storeArchive(ctx, storage, archive)
	if err != nil {
		return err
	}

	// Perform deletion afterwards so that if there is an error saving we
	// haven't messed with the current policy
	for i := p.LatestVersion - len(p.Keys) + 1; i < p.MinDecryptionVersion; i++ {
		delete(p.Keys, strconv.Itoa(i))
	}

	return nil
}

func (p *Policy) Persist(ctx context.Context, storage logical.Storage) (retErr error) {
	if atomic.LoadUint32(&p.deleted) == 1 {
		return errors.New("key has been deleted, not persisting")
	}

	// Other functions will take care of restoring other values; this is just
	// responsible for archiving and keys since the archive function can modify
	// keys. At the moment one of the other functions calling persist will also
	// roll back keys, but better safe than sorry and this doesn't happen
	// enough to worry about the speed tradeoff.
	priorArchiveVersion := p.ArchiveVersion
	var priorKeys keyEntryMap

	if p.Keys != nil {
		priorKeys = keyEntryMap{}
		for k, v := range p.Keys {
			priorKeys[k] = v
		}
	}

	defer func() {
		if retErr != nil {
			p.ArchiveVersion = priorArchiveVersion
			p.Keys = priorKeys
		}
	}()

	err := p.handleArchiving(ctx, storage)
	if err != nil {
		return err
	}

	// Encode the policy
	buf, err := p.Serialize()
	if err != nil {
		return err
	}

	// Write the policy into storage
	err = storage.Put(ctx, &logical.StorageEntry{
		Key:   path.Join(p.StoragePrefix, "policy", p.Name),
		Value: buf,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *Policy) Serialize() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Policy) NeedsUpgrade() bool {
	// Ensure we've moved from Key -> Keys
	if p.Key != nil && len(p.Key) > 0 {
		return true
	}

	// With archiving, past assumptions about the length of the keys map are no
	// longer valid
	if p.LatestVersion == 0 && len(p.Keys) != 0 {
		return true
	}

	// We disallow setting the version to 0, since they start at 1 since moving
	// to rotate-able keys, so update if it's set to 0
	if p.MinDecryptionVersion == 0 {
		return true
	}

	// On first load after an upgrade, copy keys to the archive
	if p.ArchiveVersion == 0 {
		return true
	}

	// Need to write the version if zero; for version 3 on we set this to -1 to
	// ignore it since we store this information in each key entry
	if p.ConvergentEncryption && p.ConvergentVersion == 0 {
		return true
	}

	if p.Keys[strconv.Itoa(p.LatestVersion)].HMACKey == nil || len(p.Keys[strconv.Itoa(p.LatestVersion)].HMACKey) == 0 {
		return true
	}

	return false
}

func (p *Policy) Upgrade(ctx context.Context, storage logical.Storage) (retErr error) {
	priorKey := p.Key
	priorLatestVersion := p.LatestVersion
	priorMinDecryptionVersion := p.MinDecryptionVersion
	priorConvergentVersion := p.ConvergentVersion
	var priorKeys keyEntryMap

	if p.Keys != nil {
		priorKeys = keyEntryMap{}
		for k, v := range p.Keys {
			priorKeys[k] = v
		}
	}

	defer func() {
		if retErr != nil {
			p.Key = priorKey
			p.LatestVersion = priorLatestVersion
			p.MinDecryptionVersion = priorMinDecryptionVersion
			p.ConvergentVersion = priorConvergentVersion
			p.Keys = priorKeys
		}
	}()

	persistNeeded := false
	// Ensure we've moved from Key -> Keys
	if p.Key != nil && len(p.Key) > 0 {
		p.MigrateKeyToKeysMap()
		persistNeeded = true
	}

	// With archiving, past assumptions about the length of the keys map are no
	// longer valid
	if p.LatestVersion == 0 && len(p.Keys) != 0 {
		p.LatestVersion = len(p.Keys)
		persistNeeded = true
	}

	// We disallow setting the version to 0, since they start at 1 since moving
	// to rotate-able keys, so update if it's set to 0
	if p.MinDecryptionVersion == 0 {
		p.MinDecryptionVersion = 1
		persistNeeded = true
	}

	// On first load after an upgrade, copy keys to the archive
	if p.ArchiveVersion == 0 {
		persistNeeded = true
	}

	if p.ConvergentEncryption && p.ConvergentVersion == 0 {
		p.ConvergentVersion = 1
		persistNeeded = true
	}

	if p.Keys[strconv.Itoa(p.LatestVersion)].HMACKey == nil || len(p.Keys[strconv.Itoa(p.LatestVersion)].HMACKey) == 0 {
		entry := p.Keys[strconv.Itoa(p.LatestVersion)]
		hmacKey, err := uuid.GenerateRandomBytes(32)
		if err != nil {
			return err
		}
		entry.HMACKey = hmacKey
		p.Keys[strconv.Itoa(p.LatestVersion)] = entry
		persistNeeded = true
	}

	if persistNeeded {
		err := p.Persist(ctx, storage)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeriveKey is used to derive the encryption key that should be used depending
// on the policy. If derivation is disabled the raw key is used and no context
// is required, otherwise the KDF mode is used with the context to derive the
// proper key.
func (p *Policy) DeriveKey(context []byte, ver, numBytes int) ([]byte, error) {
	// Fast-path non-derived keys
	if !p.Derived {
		return p.Keys[strconv.Itoa(ver)].Key, nil
	}

	if !p.Type.DerivationSupported() {
		return nil, errutil.UserError{Err: fmt.Sprintf("derivation not supported for key type %v", p.Type)}
	}

	if p.Keys == nil || p.LatestVersion == 0 {
		return nil, errutil.InternalError{Err: "unable to access the key; no key versions found"}
	}

	if ver <= 0 || ver > p.LatestVersion {
		return nil, errutil.UserError{Err: "invalid key version"}
	}

	// Ensure a context is provided
	if len(context) == 0 {
		return nil, errutil.UserError{Err: "missing 'context' for key derivation; the key was created using a derived key, which means additional, per-request information must be included in order to perform operations with the key"}
	}

	switch p.KDF {
	case Kdf_hmac_sha256_counter:
		prf := kdf.HMACSHA256PRF
		prfLen := kdf.HMACSHA256PRFLen
		return kdf.CounterMode(prf, prfLen, p.Keys[strconv.Itoa(ver)].Key, context, 256)

	case Kdf_hkdf_sha256:
		reader := hkdf.New(sha256.New, p.Keys[strconv.Itoa(ver)].Key, nil, context)
		derBytes := bytes.NewBuffer(nil)
		derBytes.Grow(numBytes)
		limReader := &io.LimitedReader{
			R: reader,
			N: int64(numBytes),
		}

		switch p.Type {
		case KeyType_AES256_GCM96, KeyType_ChaCha20_Poly1305:
			n, err := derBytes.ReadFrom(limReader)
			if err != nil {
				return nil, errutil.InternalError{Err: fmt.Sprintf("error reading returned derived bytes: %v", err)}
			}
			if n != int64(numBytes) {
				return nil, errutil.InternalError{Err: fmt.Sprintf("unable to read enough derived bytes, needed %d, got %d", numBytes, n)}
			}
			return derBytes.Bytes(), nil

		case KeyType_ED25519:
			// We use the limited reader containing the derived bytes as the
			// "random" input to the generation function
			_, pri, err := ed25519.GenerateKey(limReader)
			if err != nil {
				return nil, errutil.InternalError{Err: fmt.Sprintf("error generating derived key: %v", err)}
			}
			return pri, nil

		default:
			return nil, errutil.InternalError{Err: "unsupported key type for derivation"}
		}

	default:
		return nil, errutil.InternalError{Err: "unsupported key derivation mode"}
	}
}

func (p *Policy) convergentVersion(ver int) int {
	if !p.ConvergentEncryption {
		return 0
	}

	convergentVersion := p.ConvergentVersion
	if convergentVersion == 0 {
		// For some reason, not upgraded yet
		convergentVersion = 1
	}
	currKey := p.Keys[strconv.Itoa(ver)]
	if currKey.ConvergentVersion != 0 {
		convergentVersion = currKey.ConvergentVersion
	}

	return convergentVersion
}

func (p *Policy) Encrypt(ver int, context, nonce []byte, value string) (string, error) {
	if !p.Type.EncryptionSupported() {
		return "", errutil.UserError{Err: fmt.Sprintf("message encryption not supported for key type %v", p.Type)}
	}

	// Decode the plaintext value
	plaintext, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", errutil.UserError{Err: err.Error()}
	}

	switch {
	case ver == 0:
		ver = p.LatestVersion
	case ver < 0:
		return "", errutil.UserError{Err: "requested version for encryption is negative"}
	case ver > p.LatestVersion:
		return "", errutil.UserError{Err: "requested version for encryption is higher than the latest key version"}
	case ver < p.MinEncryptionVersion:
		return "", errutil.UserError{Err: "requested version for encryption is less than the minimum encryption key version"}
	}

	var ciphertext []byte

	switch p.Type {
	case KeyType_AES256_GCM96, KeyType_ChaCha20_Poly1305:
		hmacKey := context

		var aead cipher.AEAD
		var encKey []byte
		var deriveHMAC bool

		numBytes := 32
		if p.convergentVersion(ver) > 2 {
			deriveHMAC = true
			numBytes = 64
		}
		key, err := p.DeriveKey(context, ver, numBytes)
		if err != nil {
			return "", err
		}

		if len(key) < numBytes {
			return "", errutil.InternalError{Err: "could not derive key, length too small"}
		}

		encKey = key[:32]
		if len(encKey) != 32 {
			return "", errutil.InternalError{Err: "could not derive enc key, length not correct"}
		}
		if deriveHMAC {
			hmacKey = key[32:]
			if len(hmacKey) != 32 {
				return "", errutil.InternalError{Err: "could not derive hmac key, length not correct"}
			}
		}

		switch p.Type {
		case KeyType_AES256_GCM96:
			// Setup the cipher
			aesCipher, err := aes.NewCipher(encKey)
			if err != nil {
				return "", errutil.InternalError{Err: err.Error()}
			}

			// Setup the GCM AEAD
			gcm, err := cipher.NewGCM(aesCipher)
			if err != nil {
				return "", errutil.InternalError{Err: err.Error()}
			}

			aead = gcm

		case KeyType_ChaCha20_Poly1305:
			cha, err := chacha20poly1305.New(encKey)
			if err != nil {
				return "", errutil.InternalError{Err: err.Error()}
			}

			aead = cha
		}

		if p.ConvergentEncryption {
			convergentVersion := p.convergentVersion(ver)
			switch convergentVersion {
			case 1:
				if len(nonce) != aead.NonceSize() {
					return "", errutil.UserError{Err: fmt.Sprintf("base64-decoded nonce must be %d bytes long when using convergent encryption with this key", aead.NonceSize())}
				}
			case 2, 3:
				if len(hmacKey) == 0 {
					return "", errutil.InternalError{Err: fmt.Sprintf("invalid hmac key length of zero")}
				}
				nonceHmac := hmac.New(sha256.New, hmacKey)
				nonceHmac.Write(plaintext)
				nonceSum := nonceHmac.Sum(nil)
				nonce = nonceSum[:aead.NonceSize()]
			default:
				return "", errutil.InternalError{Err: fmt.Sprintf("unhandled convergent version %d", convergentVersion)}
			}
		} else {
			// Compute random nonce
			nonce, err = uuid.GenerateRandomBytes(aead.NonceSize())
			if err != nil {
				return "", errutil.InternalError{Err: err.Error()}
			}
		}

		// Encrypt and tag with AEAD
		ciphertext = aead.Seal(nil, nonce, plaintext, nil)

		// Place the encrypted data after the nonce
		if !p.ConvergentEncryption || p.convergentVersion(ver) > 1 {
			ciphertext = append(nonce, ciphertext...)
		}

	case KeyType_RSA2048, KeyType_RSA4096:
		key := p.Keys[strconv.Itoa(ver)].RSAKey
		ciphertext, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, &key.PublicKey, plaintext, nil)
		if err != nil {
			return "", errutil.InternalError{Err: fmt.Sprintf("failed to RSA encrypt the plaintext: %v", err)}
		}

	default:
		return "", errutil.InternalError{Err: fmt.Sprintf("unsupported key type %v", p.Type)}
	}

	// Convert to base64
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	// Prepend some information
	encoded = p.getVersionPrefix(ver) + encoded

	return encoded, nil
}

func (p *Policy) Decrypt(context, nonce []byte, value string) (string, error) {
	if !p.Type.DecryptionSupported() {
		return "", errutil.UserError{Err: fmt.Sprintf("message decryption not supported for key type %v", p.Type)}
	}

	tplParts, err := p.getTemplateParts()
	if err != nil {
		return "", err
	}

	// Verify the prefix
	if !strings.HasPrefix(value, tplParts[0]) {
		return "", errutil.UserError{Err: "invalid ciphertext: no prefix"}
	}

	splitVerCiphertext := strings.SplitN(strings.TrimPrefix(value, tplParts[0]), tplParts[1], 2)
	if len(splitVerCiphertext) != 2 {
		return "", errutil.UserError{Err: "invalid ciphertext: wrong number of fields"}
	}

	ver, err := strconv.Atoi(splitVerCiphertext[0])
	if err != nil {
		return "", errutil.UserError{Err: "invalid ciphertext: version number could not be decoded"}
	}

	if ver == 0 {
		// Compatibility mode with initial implementation, where keys start at
		// zero
		ver = 1
	}

	if ver > p.LatestVersion {
		return "", errutil.UserError{Err: "invalid ciphertext: version is too new"}
	}

	if p.MinDecryptionVersion > 0 && ver < p.MinDecryptionVersion {
		return "", errutil.UserError{Err: ErrTooOld}
	}

	convergentVersion := p.convergentVersion(ver)
	if convergentVersion == 1 && (nonce == nil || len(nonce) == 0) {
		return "", errutil.UserError{Err: "invalid convergent nonce supplied"}
	}

	// Decode the base64
	decoded, err := base64.StdEncoding.DecodeString(splitVerCiphertext[1])
	if err != nil {
		return "", errutil.UserError{Err: "invalid ciphertext: could not decode base64"}
	}

	var plain []byte

	switch p.Type {
	case KeyType_AES256_GCM96, KeyType_ChaCha20_Poly1305:
		var aead cipher.AEAD

		encKey, err := p.DeriveKey(context, ver, 32)
		if err != nil {
			return "", err
		}

		if len(encKey) != 32 {
			return "", errutil.InternalError{Err: "could not derive enc key, length not correct"}
		}

		switch p.Type {
		case KeyType_AES256_GCM96:
			// Setup the cipher
			aesCipher, err := aes.NewCipher(encKey)
			if err != nil {
				return "", errutil.InternalError{Err: err.Error()}
			}

			// Setup the GCM AEAD
			gcm, err := cipher.NewGCM(aesCipher)
			if err != nil {
				return "", errutil.InternalError{Err: err.Error()}
			}

			aead = gcm

		case KeyType_ChaCha20_Poly1305:
			cha, err := chacha20poly1305.New(encKey)
			if err != nil {
				return "", errutil.InternalError{Err: err.Error()}
			}

			aead = cha
		}

		if len(decoded) < aead.NonceSize() {
			return "", errutil.UserError{Err: "invalid ciphertext length"}
		}

		// Extract the nonce and ciphertext
		var ciphertext []byte
		if p.ConvergentEncryption && convergentVersion == 1 {
			ciphertext = decoded
		} else {
			nonce = decoded[:aead.NonceSize()]
			ciphertext = decoded[aead.NonceSize():]
		}

		// Verify and Decrypt
		plain, err = aead.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return "", errutil.UserError{Err: "invalid ciphertext: unable to decrypt"}
		}

	case KeyType_RSA2048, KeyType_RSA4096:
		key := p.Keys[strconv.Itoa(ver)].RSAKey
		plain, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, key, decoded, nil)
		if err != nil {
			return "", errutil.InternalError{Err: fmt.Sprintf("failed to RSA decrypt the ciphertext: %v", err)}
		}

	default:
		return "", errutil.InternalError{Err: fmt.Sprintf("unsupported key type %v", p.Type)}
	}

	return base64.StdEncoding.EncodeToString(plain), nil
}

func (p *Policy) HMACKey(version int) ([]byte, error) {
	switch {
	case version < 0:
		return nil, fmt.Errorf("key version does not exist (cannot be negative)")
	case version > p.LatestVersion:
		return nil, fmt.Errorf("key version does not exist; latest key version is %d", p.LatestVersion)
	}

	if p.Keys[strconv.Itoa(version)].HMACKey == nil {
		return nil, fmt.Errorf("no HMAC key exists for that key version")
	}

	return p.Keys[strconv.Itoa(version)].HMACKey, nil
}

func (p *Policy) Sign(ver int, context, input []byte, hashAlgorithm HashType, sigAlgorithm string, marshaling MarshalingType) (*SigningResult, error) {
	if !p.Type.SigningSupported() {
		return nil, fmt.Errorf("message signing not supported for key type %v", p.Type)
	}

	switch {
	case ver == 0:
		ver = p.LatestVersion
	case ver < 0:
		return nil, errutil.UserError{Err: "requested version for signing is negative"}
	case ver > p.LatestVersion:
		return nil, errutil.UserError{Err: "requested version for signing is higher than the latest key version"}
	case p.MinEncryptionVersion > 0 && ver < p.MinEncryptionVersion:
		return nil, errutil.UserError{Err: "requested version for signing is less than the minimum encryption key version"}
	}

	var sig []byte
	var pubKey []byte
	var err error
	switch p.Type {
	case KeyType_ECDSA_P256:
		curveBits := 256
		keyParams := p.Keys[strconv.Itoa(ver)]
		key := &ecdsa.PrivateKey{
			PublicKey: ecdsa.PublicKey{
				Curve: elliptic.P256(),
				X:     keyParams.EC_X,
				Y:     keyParams.EC_Y,
			},
			D: keyParams.EC_D,
		}

		r, s, err := ecdsa.Sign(rand.Reader, key, input)
		if err != nil {
			return nil, err
		}

		switch marshaling {
		case MarshalingTypeASN1:
			// This is used by openssl and X.509
			sig, err = asn1.Marshal(ecdsaSignature{
				R: r,
				S: s,
			})
			if err != nil {
				return nil, err
			}

		case MarshalingTypeJWS:
			// This is used by JWS

			// First we have to get the length of the curve in bytes. Although
			// we only support 256 now, we'll do this in an agnostic way so we
			// can reuse this marshaling if we support e.g. 521. Getting the
			// number of bytes without rounding up would be 65.125 so we need
			// to add one in that case.
			keyLen := curveBits / 8
			if curveBits%8 > 0 {
				keyLen++
			}

			// Now create the output array
			sig = make([]byte, keyLen*2)
			rb := r.Bytes()
			sb := s.Bytes()
			copy(sig[keyLen-len(rb):], rb)
			copy(sig[2*keyLen-len(sb):], sb)

		default:
			return nil, errutil.UserError{Err: "requested marshaling type is invalid"}
		}

	case KeyType_ED25519:
		var key ed25519.PrivateKey

		if p.Derived {
			// Derive the key that should be used
			var err error
			key, err = p.DeriveKey(context, ver, 32)
			if err != nil {
				return nil, errutil.InternalError{Err: fmt.Sprintf("error deriving key: %v", err)}
			}
			pubKey = key.Public().(ed25519.PublicKey)
		} else {
			key = ed25519.PrivateKey(p.Keys[strconv.Itoa(ver)].Key)
		}

		// Per docs, do not pre-hash ed25519; it does two passes and performs
		// its own hashing
		sig, err = key.Sign(rand.Reader, input, crypto.Hash(0))
		if err != nil {
			return nil, err
		}

	case KeyType_RSA2048, KeyType_RSA4096:
		key := p.Keys[strconv.Itoa(ver)].RSAKey

		var algo crypto.Hash
		switch hashAlgorithm {
		case HashTypeSHA1:
			algo = crypto.SHA1
		case HashTypeSHA2224:
			algo = crypto.SHA224
		case HashTypeSHA2256:
			algo = crypto.SHA256
		case HashTypeSHA2384:
			algo = crypto.SHA384
		case HashTypeSHA2512:
			algo = crypto.SHA512
		default:
			return nil, errutil.InternalError{Err: "unsupported hash algorithm"}
		}

		if sigAlgorithm == "" {
			sigAlgorithm = "pss"
		}

		switch sigAlgorithm {
		case "pss":
			sig, err = rsa.SignPSS(rand.Reader, key, algo, input, nil)
			if err != nil {
				return nil, err
			}
		case "pkcs1v15":
			sig, err = rsa.SignPKCS1v15(rand.Reader, key, algo, input)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errutil.InternalError{Err: fmt.Sprintf("unsupported rsa signature algorithm %s", sigAlgorithm)}
		}

	default:
		return nil, fmt.Errorf("unsupported key type %v", p.Type)
	}

	// Convert to base64
	var encoded string
	switch marshaling {
	case MarshalingTypeASN1:
		encoded = base64.StdEncoding.EncodeToString(sig)
	case MarshalingTypeJWS:
		encoded = base64.RawURLEncoding.EncodeToString(sig)
	}
	res := &SigningResult{
		Signature: p.getVersionPrefix(ver) + encoded,
		PublicKey: pubKey,
	}

	return res, nil
}

func (p *Policy) VerifySignature(context, input []byte, hashAlgorithm HashType, sigAlgorithm string, marshaling MarshalingType, sig string) (bool, error) {
	if !p.Type.SigningSupported() {
		return false, errutil.UserError{Err: fmt.Sprintf("message verification not supported for key type %v", p.Type)}
	}

	tplParts, err := p.getTemplateParts()
	if err != nil {
		return false, err
	}

	// Verify the prefix
	if !strings.HasPrefix(sig, tplParts[0]) {
		return false, errutil.UserError{Err: "invalid signature: no prefix"}
	}

	splitVerSig := strings.SplitN(strings.TrimPrefix(sig, tplParts[0]), tplParts[1], 2)
	if len(splitVerSig) != 2 {
		return false, errutil.UserError{Err: "invalid signature: wrong number of fields"}
	}

	ver, err := strconv.Atoi(splitVerSig[0])
	if err != nil {
		return false, errutil.UserError{Err: "invalid signature: version number could not be decoded"}
	}

	if ver > p.LatestVersion {
		return false, errutil.UserError{Err: "invalid signature: version is too new"}
	}

	if p.MinDecryptionVersion > 0 && ver < p.MinDecryptionVersion {
		return false, errutil.UserError{Err: ErrTooOld}
	}

	var sigBytes []byte
	switch marshaling {
	case MarshalingTypeASN1:
		sigBytes, err = base64.StdEncoding.DecodeString(splitVerSig[1])
	case MarshalingTypeJWS:
		sigBytes, err = base64.RawURLEncoding.DecodeString(splitVerSig[1])
	default:
		return false, errutil.UserError{Err: "requested marshaling type is invalid"}
	}
	if err != nil {
		return false, errutil.UserError{Err: "invalid base64 signature value"}
	}

	switch p.Type {
	case KeyType_ECDSA_P256:
		var ecdsaSig ecdsaSignature

		switch marshaling {
		case MarshalingTypeASN1:
			rest, err := asn1.Unmarshal(sigBytes, &ecdsaSig)
			if err != nil {
				return false, errutil.UserError{Err: "supplied signature is invalid"}
			}
			if rest != nil && len(rest) != 0 {
				return false, errutil.UserError{Err: "supplied signature contains extra data"}
			}

		case MarshalingTypeJWS:
			paramLen := len(sigBytes) / 2
			rb := sigBytes[:paramLen]
			sb := sigBytes[paramLen:]
			ecdsaSig.R = new(big.Int)
			ecdsaSig.R.SetBytes(rb)
			ecdsaSig.S = new(big.Int)
			ecdsaSig.S.SetBytes(sb)
		}

		keyParams := p.Keys[strconv.Itoa(ver)]
		key := &ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     keyParams.EC_X,
			Y:     keyParams.EC_Y,
		}

		return ecdsa.Verify(key, input, ecdsaSig.R, ecdsaSig.S), nil

	case KeyType_ED25519:
		var key ed25519.PrivateKey

		if p.Derived {
			// Derive the key that should be used
			var err error
			key, err = p.DeriveKey(context, ver, 32)
			if err != nil {
				return false, errutil.InternalError{Err: fmt.Sprintf("error deriving key: %v", err)}
			}
		} else {
			key = ed25519.PrivateKey(p.Keys[strconv.Itoa(ver)].Key)
		}

		return ed25519.Verify(key.Public().(ed25519.PublicKey), input, sigBytes), nil

	case KeyType_RSA2048, KeyType_RSA4096:
		key := p.Keys[strconv.Itoa(ver)].RSAKey

		var algo crypto.Hash
		switch hashAlgorithm {
		case HashTypeSHA1:
			algo = crypto.SHA1
		case HashTypeSHA2224:
			algo = crypto.SHA224
		case HashTypeSHA2256:
			algo = crypto.SHA256
		case HashTypeSHA2384:
			algo = crypto.SHA384
		case HashTypeSHA2512:
			algo = crypto.SHA512
		default:
			return false, errutil.InternalError{Err: "unsupported hash algorithm"}
		}

		if sigAlgorithm == "" {
			sigAlgorithm = "pss"
		}

		switch sigAlgorithm {
		case "pss":
			err = rsa.VerifyPSS(&key.PublicKey, algo, input, sigBytes, nil)
		case "pkcs1v15":
			err = rsa.VerifyPKCS1v15(&key.PublicKey, algo, input, sigBytes)
		default:
			return false, errutil.InternalError{Err: fmt.Sprintf("unsupported rsa signature algorithm %s", sigAlgorithm)}
		}

		return err == nil, nil

	default:
		return false, errutil.InternalError{Err: fmt.Sprintf("unsupported key type %v", p.Type)}
	}
}

func (p *Policy) Rotate(ctx context.Context, storage logical.Storage) (retErr error) {
	priorLatestVersion := p.LatestVersion
	priorMinDecryptionVersion := p.MinDecryptionVersion
	var priorKeys keyEntryMap

	if p.Keys != nil {
		priorKeys = keyEntryMap{}
		for k, v := range p.Keys {
			priorKeys[k] = v
		}
	}

	defer func() {
		if retErr != nil {
			p.LatestVersion = priorLatestVersion
			p.MinDecryptionVersion = priorMinDecryptionVersion
			p.Keys = priorKeys
		}
	}()

	if p.Keys == nil {
		// This is an initial key rotation when generating a new policy. We
		// don't need to call migrate here because if we've called getPolicy to
		// get the policy in the first place it will have been run.
		p.Keys = keyEntryMap{}
	}

	p.LatestVersion += 1
	now := time.Now()
	entry := KeyEntry{
		CreationTime:           now,
		DeprecatedCreationTime: now.Unix(),
	}

	hmacKey, err := uuid.GenerateRandomBytes(32)
	if err != nil {
		return err
	}
	entry.HMACKey = hmacKey

	switch p.Type {
	case KeyType_AES256_GCM96, KeyType_ChaCha20_Poly1305:
		// Generate a 256bit key
		newKey, err := uuid.GenerateRandomBytes(32)
		if err != nil {
			return err
		}
		entry.Key = newKey

	case KeyType_ECDSA_P256:
		privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return err
		}
		entry.EC_D = privKey.D
		entry.EC_X = privKey.X
		entry.EC_Y = privKey.Y
		derBytes, err := x509.MarshalPKIXPublicKey(privKey.Public())
		if err != nil {
			return errwrap.Wrapf("error marshaling public key: {{err}}", err)
		}
		pemBlock := &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: derBytes,
		}
		pemBytes := pem.EncodeToMemory(pemBlock)
		if pemBytes == nil || len(pemBytes) == 0 {
			return fmt.Errorf("error PEM-encoding public key")
		}
		entry.FormattedPublicKey = string(pemBytes)

	case KeyType_ED25519:
		pub, pri, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return err
		}
		entry.Key = pri
		entry.FormattedPublicKey = base64.StdEncoding.EncodeToString(pub)

	case KeyType_RSA2048, KeyType_RSA4096:
		bitSize := 2048
		if p.Type == KeyType_RSA4096 {
			bitSize = 4096
		}

		entry.RSAKey, err = rsa.GenerateKey(rand.Reader, bitSize)
		if err != nil {
			return err
		}
	}

	if p.ConvergentEncryption {
		if p.ConvergentVersion == -1 || p.ConvergentVersion > 1 {
			entry.ConvergentVersion = currentConvergentVersion
		}
	}

	p.Keys[strconv.Itoa(p.LatestVersion)] = entry

	// This ensures that with new key creations min decryption version is set
	// to 1 rather than the int default of 0, since keys start at 1 (either
	// fresh or after migration to the key map)
	if p.MinDecryptionVersion == 0 {
		p.MinDecryptionVersion = 1
	}

	return p.Persist(ctx, storage)
}

func (p *Policy) MigrateKeyToKeysMap() {
	now := time.Now()
	p.Keys = keyEntryMap{
		"1": KeyEntry{
			Key:                    p.Key,
			CreationTime:           now,
			DeprecatedCreationTime: now.Unix(),
		},
	}
	p.Key = nil
}

// Backup should be called with an exclusive lock held on the policy
func (p *Policy) Backup(ctx context.Context, storage logical.Storage) (out string, retErr error) {
	if !p.Exportable {
		return "", fmt.Errorf("exporting is disallowed on the policy")
	}

	if !p.AllowPlaintextBackup {
		return "", fmt.Errorf("plaintext backup is disallowed on the policy")
	}

	priorBackupInfo := p.BackupInfo

	defer func() {
		if retErr != nil {
			p.BackupInfo = priorBackupInfo
		}
	}()

	// Create a record of this backup operation in the policy
	p.BackupInfo = &BackupInfo{
		Time:    time.Now(),
		Version: p.LatestVersion,
	}
	err := p.Persist(ctx, storage)
	if err != nil {
		return "", errwrap.Wrapf("failed to persist policy with backup info: {{err}}", err)
	}

	// Load the archive only after persisting the policy as the archive can get
	// adjusted while persisting the policy
	archivedKeys, err := p.LoadArchive(ctx, storage)
	if err != nil {
		return "", err
	}

	keyData := &KeyData{
		Policy:       p,
		ArchivedKeys: archivedKeys,
	}

	encodedBackup, err := jsonutil.EncodeJSON(keyData)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encodedBackup), nil
}

func (p *Policy) getTemplateParts() ([]string, error) {
	partsRaw, ok := p.versionPrefixCache.Load("template-parts")
	if ok {
		return partsRaw.([]string), nil
	}

	template := p.VersionTemplate
	if template == "" {
		template = DefaultVersionTemplate
	}

	tplParts := strings.Split(template, "{{version}}")
	if len(tplParts) != 2 {
		return nil, errutil.InternalError{Err: "error parsing version template"}
	}

	p.versionPrefixCache.Store("template-parts", tplParts)
	return tplParts, nil
}

func (p *Policy) getVersionPrefix(ver int) string {
	prefixRaw, ok := p.versionPrefixCache.Load(ver)
	if ok {
		return prefixRaw.(string)
	}

	template := p.VersionTemplate
	if template == "" {
		template = DefaultVersionTemplate
	}

	prefix := strings.Replace(template, "{{version}}", strconv.Itoa(ver), -1)
	p.versionPrefixCache.Store(ver, prefix)

	return prefix
}
