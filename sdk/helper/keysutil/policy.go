// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keysutil

import (
	"bytes"
	"context"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	stdlibEd25519 "crypto/ed25519"
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
	"hash"
	"io"
	"math/big"
	"path"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/kdf"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/tink-crypto/tink-go/v2/kwp/subtle"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/hkdf"
)

// Careful with iota; don't put anything before it in this const block because
// we need the default of zero to be the old-style KDF
const (
	Kdf_hmac_sha256_counter = iota // built-in helper
	Kdf_hkdf_sha256                // golang.org/x/crypto/hkdf

	HmacMinKeySize = 256 / 8
	HmacMaxKeySize = 4096 / 8
)

// Or this one...we need the default of zero to be the original AES256-GCM96
const (
	KeyType_AES256_GCM96 = iota
	KeyType_ECDSA_P256
	KeyType_ED25519
	KeyType_RSA2048
	KeyType_RSA4096
	KeyType_ChaCha20_Poly1305
	KeyType_ECDSA_P384
	KeyType_ECDSA_P521
	KeyType_AES128_GCM96
	KeyType_RSA3072
	KeyType_MANAGED_KEY
	KeyType_HMAC
	KeyType_AES128_CMAC
	KeyType_AES256_CMAC
	KeyType_ML_DSA
	KeyType_HYBRID
	// If adding to this list please update allTestKeyTypes in policy_test.go
)

const (
	ParameterSet_ML_DSA_44 = "44"
	ParameterSet_ML_DSA_65 = "65"
	ParameterSet_ML_DSA_87 = "87"
)

const (
	// ErrTooOld is returned whtn the ciphertext or signatures's key version is
	// too old.
	ErrTooOld = "ciphertext or signature version is disallowed by policy (too old)"

	// DefaultVersionTemplate is used when no version template is provided.
	DefaultVersionTemplate = "vault:v{{version}}:"
)

type PaddingScheme string

const (
	PaddingScheme_OAEP     = PaddingScheme("oaep")
	PaddingScheme_PKCS1v15 = PaddingScheme("pkcs1v15")
)

var genEd25519Options = func(hashAlgorithm HashType, signatureContext string) (*stdlibEd25519.Options, error) {
	if signatureContext != "" {
		return nil, fmt.Errorf("signature context is not supported feature")
	}

	if hashAlgorithm == HashTypeSHA2512 {
		return nil, fmt.Errorf("hash algorithm of SHA2 512 is not supported feature")
	}

	return &stdlibEd25519.Options{
		Hash: crypto.Hash(0),
	}, nil
}

func (p PaddingScheme) String() string {
	return string(p)
}

// ParsePaddingScheme expects a lower case string that can be directly compared to
// a defined padding scheme or returns an error.
func ParsePaddingScheme(s string) (PaddingScheme, error) {
	switch s {
	case PaddingScheme_OAEP.String():
		return PaddingScheme_OAEP, nil
	case PaddingScheme_PKCS1v15.String():
		return PaddingScheme_PKCS1v15, nil
	default:
		return "", fmt.Errorf("unknown padding scheme: %s", s)
	}
}

type AEADFactory interface {
	GetAEAD(iv []byte) (cipher.AEAD, error)
}

type AssociatedDataFactory interface {
	GetAssociatedData() ([]byte, error)
}

type ManagedKeyFactory interface {
	GetManagedKeyParameters() ManagedKeyParameters
}

type RestoreInfo struct {
	Time    time.Time `json:"time"`
	Version int       `json:"version"`
}

type BackupInfo struct {
	Time    time.Time `json:"time"`
	Version int       `json:"version"`
}

type SigningOptions struct {
	HashAlgorithm    HashType
	Marshaling       MarshalingType
	SaltLength       int
	SigAlgorithm     string
	SigContext       string // Provide a context for Ed25519ctx signatures
	ManagedKeyParams ManagedKeyParameters
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
	case KeyType_AES128_GCM96, KeyType_AES256_GCM96, KeyType_ChaCha20_Poly1305, KeyType_RSA2048, KeyType_RSA3072, KeyType_RSA4096, KeyType_MANAGED_KEY:
		return true
	}
	return false
}

func (kt KeyType) DecryptionSupported() bool {
	switch kt {
	case KeyType_AES128_GCM96, KeyType_AES256_GCM96, KeyType_ChaCha20_Poly1305, KeyType_RSA2048, KeyType_RSA3072, KeyType_RSA4096, KeyType_MANAGED_KEY:
		return true
	}
	return false
}

func (kt KeyType) SigningSupported() bool {
	switch kt {
	case KeyType_ECDSA_P256, KeyType_ECDSA_P384, KeyType_ECDSA_P521, KeyType_ED25519, KeyType_RSA2048, KeyType_RSA3072, KeyType_RSA4096, KeyType_MANAGED_KEY, KeyType_ML_DSA, KeyType_HYBRID:
		return true
	}
	return false
}

func (kt KeyType) HashSignatureInput() bool {
	switch kt {
	case KeyType_ECDSA_P256, KeyType_ECDSA_P384, KeyType_ECDSA_P521, KeyType_RSA2048, KeyType_RSA3072, KeyType_RSA4096, KeyType_MANAGED_KEY:
		return true
	}
	return false
}

func (kt KeyType) DerivationSupported() bool {
	switch kt {
	case KeyType_AES128_GCM96, KeyType_AES256_GCM96, KeyType_ChaCha20_Poly1305, KeyType_ED25519:
		return true
	}
	return false
}

func (kt KeyType) AssociatedDataSupported() bool {
	switch kt {
	case KeyType_AES128_GCM96, KeyType_AES256_GCM96, KeyType_ChaCha20_Poly1305, KeyType_MANAGED_KEY:
		return true
	}
	return false
}

func (kt KeyType) CMACSupported() bool {
	switch kt {
	case KeyType_AES128_CMAC, KeyType_AES256_CMAC:
		return true
	default:
		return false
	}
}

func (kt KeyType) HMACSupported() bool {
	switch {
	case kt.CMACSupported():
		return false
	case kt == KeyType_MANAGED_KEY:
		return false
	default:
		return true
	}
}

func (kt KeyType) IsPQC() bool {
	switch kt {
	case KeyType_ML_DSA, KeyType_HYBRID:
		return true
	default:
		return false
	}
}

func (kt KeyType) ImportPublicKeySupported() bool {
	switch kt {
	case KeyType_RSA2048, KeyType_RSA3072, KeyType_RSA4096, KeyType_ECDSA_P256, KeyType_ECDSA_P384, KeyType_ECDSA_P521, KeyType_ED25519:
		return true
	}
	return false
}

func (kt KeyType) PaddingSchemesSupported() bool {
	switch kt {
	case KeyType_RSA2048, KeyType_RSA3072, KeyType_RSA4096:
		return true
	default:
		return false
	}
}

func (kt KeyType) String() string {
	switch kt {
	case KeyType_AES128_GCM96:
		return "aes128-gcm96"
	case KeyType_AES256_GCM96:
		return "aes256-gcm96"
	case KeyType_ChaCha20_Poly1305:
		return "chacha20-poly1305"
	case KeyType_ECDSA_P256:
		return "ecdsa-p256"
	case KeyType_ECDSA_P384:
		return "ecdsa-p384"
	case KeyType_ECDSA_P521:
		return "ecdsa-p521"
	case KeyType_ED25519:
		return "ed25519"
	case KeyType_RSA2048:
		return "rsa-2048"
	case KeyType_RSA3072:
		return "rsa-3072"
	case KeyType_RSA4096:
		return "rsa-4096"
	case KeyType_HMAC:
		return "hmac"
	case KeyType_MANAGED_KEY:
		return "managed_key"
	case KeyType_AES128_CMAC:
		return "aes128-cmac"
	case KeyType_AES256_CMAC:
		return "aes256-cmac"
	case KeyType_ML_DSA:
		return "ml-dsa"
	case KeyType_HYBRID:
		return "hybrid"
	}

	return "[unknown]"
}

type KeyData struct {
	Policy       *Policy       `json:"policy"`
	ArchivedKeys *archivedKeys `json:"archived_keys"`
}

// KeyEntry stores the key and metadata
type KeyEntry struct {
	entKeyEntry

	// AES or some other kind that is a pure byte slice like ED25519
	Key []byte `json:"key"`

	// Key used for HMAC functions
	HMACKey []byte `json:"hmac_key"`

	// Time of creation
	CreationTime time.Time `json:"time"`

	EC_X *big.Int `json:"ec_x,omitempty"`
	EC_Y *big.Int `json:"ec_y,omitempty"`
	EC_D *big.Int `json:"ec_d,omitempty"`

	RSAKey       *rsa.PrivateKey `json:"rsa_key,omitempty"`
	RSAPublicKey *rsa.PublicKey  `json:"rsa_public_key,omitempty"`

	// The public key in an appropriate format for the type of key
	FormattedPublicKey string `json:"public_key,omitempty"`

	// If convergent is enabled, the version (falling back to what's in the
	// policy)
	ConvergentVersion int `json:"convergent_version,omitempty"`

	// This is deprecated (but still filled) in favor of the value above which
	// is more precise
	DeprecatedCreationTime int64 `json:"creation_time"`

	ManagedKeyUUID string `json:"managed_key_id,omitempty"`

	// Key entry certificate chain. If set, leaf certificate key matches the
	// KeyEntry key
	CertificateChain [][]byte `json:"certificate_chain,omitempty"`
}

func (ke *KeyEntry) IsPrivateKeyMissing() bool {
	if ke.RSAKey != nil || ke.EC_D != nil || len(ke.Key) != 0 || len(ke.ManagedKeyUUID) != 0 || !ke.IsEntPrivateKeyMissing() {
		return false
	}

	return true
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

	// ParameterSet indicates the parameter set to use with ML-DSA and SLH-DSA keys
	ParameterSet string
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
		ParameterSet:         config.ParameterSet,
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

	// Migrate RSA private keys to include their private counterpart. This lets
	// us reference RSAPublicKey whenever we need to, without necessarily
	// needing the private key handy, synchronizing the behavior with EC and
	// Ed25519 key pairs.
	switch policy.Type {
	case KeyType_RSA2048, KeyType_RSA3072, KeyType_RSA4096:
		for _, entry := range policy.Keys {
			if entry.RSAPublicKey == nil && entry.RSAKey != nil {
				entry.RSAPublicKey = entry.RSAKey.Public().(*rsa.PublicKey)
			}
		}
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

	Name    string      `json:"name"`
	Key     []byte      `json:"key,omitempty"`      // DEPRECATED
	KeySize int         `json:"key_size,omitempty"` // For algorithms with variable key sizes
	Keys    keyEntryMap `json:"keys"`

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

	// AutoRotatePeriod defines how frequently the key should automatically
	// rotate. Setting this to zero disables automatic rotation for the key.
	AutoRotatePeriod time.Duration `json:"auto_rotate_period"`

	// versionPrefixCache stores caches of version prefix strings and the split
	// version template.
	versionPrefixCache sync.Map

	// Imported indicates whether the key was generated by Vault or imported
	// from an external source
	Imported bool

	// AllowImportedKeyRotation indicates whether an imported key may be rotated by Vault
	AllowImportedKeyRotation bool

	// ParameterSet indicates the parameter set to use with ML-DSA and SLH-DSA keys
	ParameterSet string

	// HybridConfig contains the key types and parameters for hybrid keys
	HybridConfig HybridKeyConfig
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

	if p.Type.HMACSupported() {
		if p.Keys[strconv.Itoa(p.LatestVersion)].HMACKey == nil || len(p.Keys[strconv.Itoa(p.LatestVersion)].HMACKey) == 0 {
			return true
		}
	}

	return false
}

func (p *Policy) Upgrade(ctx context.Context, storage logical.Storage, randReader io.Reader) (retErr error) {
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

	if p.Type.HMACSupported() {
		if p.Keys[strconv.Itoa(p.LatestVersion)].HMACKey == nil || len(p.Keys[strconv.Itoa(p.LatestVersion)].HMACKey) == 0 {
			entry := p.Keys[strconv.Itoa(p.LatestVersion)]
			hmacKey, err := uuid.GenerateRandomBytesWithReader(32, randReader)
			if err != nil {
				return err
			}
			entry.HMACKey = hmacKey
			p.Keys[strconv.Itoa(p.LatestVersion)] = entry
			persistNeeded = true

			if p.Type == KeyType_HMAC {
				entry.HMACKey = entry.Key
			}
		}
	}

	if persistNeeded {
		err := p.Persist(ctx, storage)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetKey is used to derive the encryption key that should be used depending
// on the policy. If derivation is disabled the raw key is used and no context
// is required, otherwise the KDF mode is used with the context to derive the
// proper key.
func (p *Policy) GetKey(context []byte, ver, numBytes int) ([]byte, error) {
	// Fast-path non-derived keys
	if !p.Derived {
		keyEntry, err := p.safeGetKeyEntry(ver)
		if err != nil {
			return nil, err
		}

		return keyEntry.Key, nil
	}

	return p.DeriveKey(context, nil, ver, numBytes)
}

// DeriveKey is used to derive a symmetric key given a context and salt.  This does not
// check the policies Derived flag, but just implements the derivation logic.  GetKey
// is responsible for switching on the policy config.
func (p *Policy) DeriveKey(context, salt []byte, ver int, numBytes int) ([]byte, error) {
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

	keyEntry, err := p.safeGetKeyEntry(ver)
	if err != nil {
		return nil, err
	}

	switch p.KDF {
	case Kdf_hmac_sha256_counter:
		prf := kdf.HMACSHA256PRF
		prfLen := kdf.HMACSHA256PRFLen
		return kdf.CounterMode(prf, prfLen, keyEntry.Key, append(context, salt...), 256)

	case Kdf_hkdf_sha256:
		reader := hkdf.New(sha256.New, keyEntry.Key, salt, context)
		derBytes := bytes.NewBuffer(nil)
		derBytes.Grow(numBytes)
		limReader := &io.LimitedReader{
			R: reader,
			N: int64(numBytes),
		}

		switch p.Type {
		case KeyType_AES128_GCM96, KeyType_AES256_GCM96, KeyType_ChaCha20_Poly1305:
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

func (p *Policy) safeGetKeyEntry(ver int) (KeyEntry, error) {
	keyVerStr := strconv.Itoa(ver)
	keyEntry, ok := p.Keys[keyVerStr]
	if !ok {
		return keyEntry, errutil.UserError{Err: "no such key version"}
	}
	return keyEntry, nil
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
	return p.EncryptWithFactory(ver, context, nonce, value, nil)
}

func (p *Policy) Decrypt(context, nonce []byte, value string) (string, error) {
	return p.DecryptWithFactory(context, nonce, value, nil)
}

func (p *Policy) DecryptWithFactory(context, nonce []byte, value string, factories ...any) (string, error) {
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
	case KeyType_AES128_GCM96, KeyType_AES256_GCM96, KeyType_ChaCha20_Poly1305:
		numBytes := 32
		if p.Type == KeyType_AES128_GCM96 {
			numBytes = 16
		}

		encKey, err := p.GetKey(context, ver, numBytes)
		if err != nil {
			return "", err
		}

		if len(encKey) != numBytes {
			return "", errutil.InternalError{Err: "could not derive enc key, length not correct"}
		}

		symopts := SymmetricOpts{
			Convergent:        p.ConvergentEncryption,
			ConvergentVersion: p.ConvergentVersion,
		}
		for index, rawFactory := range factories {
			if rawFactory == nil {
				continue
			}
			switch factory := rawFactory.(type) {
			case AEADFactory:
				symopts.AEADFactory = factory
			case AssociatedDataFactory:
				symopts.AdditionalData, err = factory.GetAssociatedData()
				if err != nil {
					return "", errutil.InternalError{Err: fmt.Sprintf("unable to get associated_data/additional_data from factory[%d]: %v", index, err)}
				}
			case ManagedKeyFactory:
			default:
				return "", errutil.InternalError{Err: fmt.Sprintf("unknown type of factory[%d]: %T", index, rawFactory)}
			}
		}

		plain, err = p.SymmetricDecryptRaw(encKey, decoded, symopts)
		if err != nil {
			return "", err
		}
	case KeyType_RSA2048, KeyType_RSA3072, KeyType_RSA4096:
		paddingScheme, err := getPaddingScheme(factories)
		if err != nil {
			return "", err
		}
		keyEntry, err := p.safeGetKeyEntry(ver)
		if err != nil {
			return "", err
		}
		key := keyEntry.RSAKey

		switch paddingScheme {
		case PaddingScheme_PKCS1v15:
			plain, err = rsa.DecryptPKCS1v15(rand.Reader, key, decoded)
		case PaddingScheme_OAEP:
			plain, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, key, decoded, nil)
		default:
			return "", errutil.InternalError{Err: fmt.Sprintf("unsupported RSA padding scheme %s", paddingScheme)}
		}
		if err != nil {
			return "", errutil.InternalError{Err: fmt.Sprintf("failed to RSA decrypt the ciphertext: %v", err)}
		}
	case KeyType_MANAGED_KEY:
		keyEntry, err := p.safeGetKeyEntry(ver)
		if err != nil {
			return "", err
		}
		var aad []byte
		var managedKeyFactory ManagedKeyFactory
		for _, f := range factories {
			switch factory := f.(type) {
			case AssociatedDataFactory:
				aad, err = factory.GetAssociatedData()
				if err != nil {
					return "", err
				}
			case ManagedKeyFactory:
				managedKeyFactory = factory
			}
		}

		if managedKeyFactory == nil {
			return "", errors.New("key type is managed_key, but managed key parameters were not provided")
		}

		plain, err = p.decryptWithManagedKey(managedKeyFactory.GetManagedKeyParameters(), keyEntry, decoded, nonce, aad)
		if err != nil {
			return "", err
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
	keyEntry, err := p.safeGetKeyEntry(version)
	if err != nil {
		return nil, err
	}

	if p.Type == KeyType_HMAC {
		return keyEntry.Key, nil
	}
	if keyEntry.HMACKey == nil {
		return nil, fmt.Errorf("no HMAC key exists for that key version")
	}
	return keyEntry.HMACKey, nil
}

func (p *Policy) CMACKey(version int) ([]byte, error) {
	switch {
	case version < 0:
		return nil, fmt.Errorf("key version does not exist (cannot be negative)")
	case version > p.LatestVersion:
		return nil, fmt.Errorf("key version does not exist; latest key version is %d", p.LatestVersion)
	}
	keyEntry, err := p.safeGetKeyEntry(version)
	if err != nil {
		return nil, err
	}

	if p.Type.CMACSupported() {
		return keyEntry.Key, nil
	}

	return nil, fmt.Errorf("key type %s does not support CMAC operations", p.Type)
}

func (p *Policy) Sign(ver int, context, input []byte, hashAlgorithm HashType, sigAlgorithm string, marshaling MarshalingType) (*SigningResult, error) {
	return p.SignWithOptions(ver, context, input, &SigningOptions{
		HashAlgorithm: hashAlgorithm,
		Marshaling:    marshaling,
		SaltLength:    rsa.PSSSaltLengthAuto,
		SigAlgorithm:  sigAlgorithm,
	})
}

func (p *Policy) minRSAPSSSaltLength() int {
	// https://cs.opensource.google/go/go/+/refs/tags/go1.19:src/crypto/rsa/pss.go;l=247
	return rsa.PSSSaltLengthEqualsHash
}

func (p *Policy) maxRSAPSSSaltLength(keyBitLen int, hash crypto.Hash) int {
	// https://cs.opensource.google/go/go/+/refs/tags/go1.19:src/crypto/rsa/pss.go;l=288
	return (keyBitLen-1+7)/8 - 2 - hash.Size()
}

func (p *Policy) validRSAPSSSaltLength(keyBitLen int, hash crypto.Hash, saltLength int) bool {
	return p.minRSAPSSSaltLength() <= saltLength && saltLength <= p.maxRSAPSSSaltLength(keyBitLen, hash)
}

func (p *Policy) SignWithOptions(ver int, context, input []byte, options *SigningOptions) (*SigningResult, error) {
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
	keyParams, err := p.safeGetKeyEntry(ver)
	if err != nil {
		return nil, err
	}

	// Before signing, check if key has its private part, if not return error
	if keyParams.IsPrivateKeyMissing() {
		return nil, errutil.UserError{Err: "requested version for signing does not contain a private part"}
	}

	hashAlgorithm := options.HashAlgorithm
	marshaling := options.Marshaling
	saltLength := options.SaltLength
	sigAlgorithm := options.SigAlgorithm

	switch p.Type {
	case KeyType_ECDSA_P256, KeyType_ECDSA_P384, KeyType_ECDSA_P521:
		sig, err = signWithECDSA(p.Type, keyParams, input, marshaling)
	case KeyType_ED25519:
		sig, pubKey, err = p.signWithEd25519(ver, input, context, options, keyParams)
		if err != nil {
			return nil, err
		}
	case KeyType_RSA2048, KeyType_RSA3072, KeyType_RSA4096:
		key := keyParams.RSAKey

		algo, ok := CryptoHashMap[hashAlgorithm]
		if !ok {
			return nil, errutil.InternalError{Err: "unsupported hash algorithm"}
		}

		if sigAlgorithm == "" {
			sigAlgorithm = "pss"
		}

		switch sigAlgorithm {
		case "pss":
			if !p.validRSAPSSSaltLength(key.N.BitLen(), algo, saltLength) {
				return nil, errutil.UserError{Err: fmt.Sprintf("requested salt length %d is invalid", saltLength)}
			}
			sig, err = rsa.SignPSS(rand.Reader, key, algo, input, &rsa.PSSOptions{SaltLength: saltLength})
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
		sig, err = entSignWithOptions(p, input, context, ver, hashAlgorithm, options)
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

func (p *Policy) signWithEd25519(ver int, input []byte, context []byte, options *SigningOptions, keyParams KeyEntry) ([]byte, []byte, error) {
	var key ed25519.PrivateKey
	var pubKey []byte
	if p.Derived {
		// Derive the key that should be used
		var err error
		key, err = p.GetKey(context, ver, 32)
		if err != nil {
			return nil, nil, errutil.InternalError{Err: fmt.Sprintf("error deriving key: %v", err)}
		}
		pubKey = key.Public().(ed25519.PublicKey)
	} else {
		key = ed25519.PrivateKey(keyParams.Key)
	}

	opts, err := genEd25519Options(options.HashAlgorithm, options.SigContext)
	if err != nil {
		return nil, nil, errutil.UserError{Err: fmt.Sprintf("error generating Ed25519 options: %v", err)}
	}

	sig, err := key.Sign(rand.Reader, input, opts)
	if err != nil {
		return nil, nil, err
	}
	return sig, pubKey, nil
}

func signWithECDSA(keyType KeyType, keyParams KeyEntry, input []byte, marshaling MarshalingType) ([]byte, error) {
	var curveBits int
	var curve elliptic.Curve
	switch keyType {
	case KeyType_ECDSA_P256:
		curveBits = 256
		curve = elliptic.P256()
	case KeyType_ECDSA_P384:
		curveBits = 384
		curve = elliptic.P384()
	case KeyType_ECDSA_P521:
		curveBits = 521
		curve = elliptic.P521()
	default:
		return nil, fmt.Errorf("invalid key type %s for ECDSA", keyType)
	}

	key := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     keyParams.EC_X,
			Y:     keyParams.EC_Y,
		},
		D: keyParams.EC_D,
	}

	r, s, err := ecdsa.Sign(rand.Reader, key, input)
	if err != nil {
		return nil, err
	}

	var sig []byte
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

	return sig, nil
}

func (p *Policy) VerifySignature(context, input []byte, hashAlgorithm HashType, sigAlgorithm string, marshaling MarshalingType, sig string) (bool, error) {
	return p.VerifySignatureWithOptions(context, input, sig, &SigningOptions{
		HashAlgorithm: hashAlgorithm,
		Marshaling:    marshaling,
		SaltLength:    rsa.PSSSaltLengthAuto,
		SigAlgorithm:  sigAlgorithm,
	})
}

func (p *Policy) VerifySignatureWithOptions(context, input []byte, sig string, options *SigningOptions) (bool, error) {
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

	hashAlgorithm := options.HashAlgorithm
	marshaling := options.Marshaling
	saltLength := options.SaltLength
	sigAlgorithm := options.SigAlgorithm

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
	case KeyType_ECDSA_P256, KeyType_ECDSA_P384, KeyType_ECDSA_P521:
		key, err := p.safeGetKeyEntry(ver)
		if err != nil {
			return false, err
		}
		return verifyWithECDSA(p.Type, key, input, sigBytes, marshaling)
	case KeyType_ED25519:
		return p.verifyEd25519WithOptions(ver, input, context, options, sigBytes)
	case KeyType_RSA2048, KeyType_RSA3072, KeyType_RSA4096:
		keyEntry, err := p.safeGetKeyEntry(ver)
		if err != nil {
			return false, err
		}

		algo, ok := CryptoHashMap[hashAlgorithm]
		if !ok {
			return false, errutil.InternalError{Err: "unsupported hash algorithm"}
		}

		if sigAlgorithm == "" {
			sigAlgorithm = "pss"
		}

		switch sigAlgorithm {
		case "pss":
			publicKey := keyEntry.RSAPublicKey
			if !keyEntry.IsPrivateKeyMissing() {
				publicKey = &keyEntry.RSAKey.PublicKey
			}
			if !p.validRSAPSSSaltLength(publicKey.N.BitLen(), algo, saltLength) {
				return false, errutil.UserError{Err: fmt.Sprintf("requested salt length %d is invalid", saltLength)}
			}
			err = rsa.VerifyPSS(publicKey, algo, input, sigBytes, &rsa.PSSOptions{SaltLength: saltLength})
		case "pkcs1v15":
			publicKey := keyEntry.RSAPublicKey
			if !keyEntry.IsPrivateKeyMissing() {
				publicKey = &keyEntry.RSAKey.PublicKey
			}
			err = rsa.VerifyPKCS1v15(publicKey, algo, input, sigBytes)
		default:
			return false, errutil.InternalError{Err: fmt.Sprintf("unsupported rsa signature algorithm %s", sigAlgorithm)}
		}

		return err == nil, nil

	default:
		return entVerifySignatureWithOptions(p, input, context, sigBytes, ver, options)
	}
}

func verifyWithECDSA(keyType KeyType, keyParams KeyEntry, input, sigBytes []byte, marshaling MarshalingType) (bool, error) {
	var curve elliptic.Curve
	switch keyType {
	case KeyType_ECDSA_P256:
		curve = elliptic.P256()
	case KeyType_ECDSA_P384:
		curve = elliptic.P384()
	case KeyType_ECDSA_P521:
		curve = elliptic.P521()
	default:
		return false, fmt.Errorf("invalid key type %s for ECDSA", keyType)
	}

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

	key := &ecdsa.PublicKey{
		Curve: curve,
		X:     keyParams.EC_X,
		Y:     keyParams.EC_Y,
	}

	return ecdsa.Verify(key, input, ecdsaSig.R, ecdsaSig.S), nil
}

func (p *Policy) verifyEd25519WithOptions(ver int, input []byte, context []byte, options *SigningOptions, sigBytes []byte) (bool, error) {
	var pub ed25519.PublicKey
	if p.Derived {
		// Derive the key that should be used
		key, err := p.GetKey(context, ver, 32)
		if err != nil {
			return false, errutil.InternalError{Err: fmt.Sprintf("error deriving key: %v", err)}
		}
		pub = ed25519.PrivateKey(key).Public().(ed25519.PublicKey)
	} else {
		keyEntry, err := p.safeGetKeyEntry(ver)
		if err != nil {
			return false, err
		}

		raw, err := base64.StdEncoding.DecodeString(keyEntry.FormattedPublicKey)
		if err != nil {
			return false, err
		}

		pub = ed25519.PublicKey(raw)
	}

	return p.verifyEd25519WithPublicKey(input, sigBytes, pub, options)
}

func (p *Policy) verifyEd25519WithPublicKey(input []byte, sigBytes []byte, pub ed25519.PublicKey, options *SigningOptions) (bool, error) {
	opts, err := genEd25519Options(options.HashAlgorithm, options.SigContext)
	if err != nil {
		return false, errutil.UserError{Err: fmt.Sprintf("error generating Ed25519 options: %v", err)}
	}
	if pub == nil {
		return false, errutil.InternalError{Err: "no Ed25519 public key on policy"}
	}
	if err := stdlibEd25519.VerifyWithOptions(pub, input, sigBytes, opts); err != nil {
		// We drop the error, just report back that we failed signature verification
		return false, nil
	}

	return true, nil
}

func (p *Policy) Import(ctx context.Context, storage logical.Storage, key []byte, randReader io.Reader) error {
	return p.ImportPublicOrPrivate(ctx, storage, key, true, randReader)
}

func (p *Policy) ImportPublicOrPrivate(ctx context.Context, storage logical.Storage, key []byte, isPrivateKey bool, randReader io.Reader) error {
	now := time.Now()
	entry := KeyEntry{
		CreationTime:           now,
		DeprecatedCreationTime: now.Unix(),
	}

	// Before we insert this entry, check if the latest version is incomplete
	// and this entry matches the current version; if so, return without
	// updating to the next version.
	if p.LatestVersion > 0 {
		latestKey := p.Keys[strconv.Itoa(p.LatestVersion)]
		if latestKey.IsPrivateKeyMissing() && isPrivateKey {
			if err := p.ImportPrivateKeyForVersion(ctx, storage, p.LatestVersion, key); err == nil {
				return nil
			}
		}
	}

	if p.Type != KeyType_HMAC {
		hmacKey, err := uuid.GenerateRandomBytesWithReader(32, randReader)
		if err != nil {
			return err
		}
		entry.HMACKey = hmacKey
	}

	if p.Type == KeyType_ED25519 && p.Derived && !isPrivateKey {
		return fmt.Errorf("unable to import only public key for derived Ed25519 key: imported key should not be an Ed25519 key pair but is instead an HKDF key")
	}

	if ((p.Type == KeyType_AES128_GCM96 || p.Type == KeyType_AES128_CMAC) && len(key) != 16) ||
		((p.Type == KeyType_AES256_GCM96 || p.Type == KeyType_ChaCha20_Poly1305 || p.Type == KeyType_AES256_CMAC) && len(key) != 32) ||
		(p.Type == KeyType_HMAC && (len(key) < HmacMinKeySize || len(key) > HmacMaxKeySize)) {
		return fmt.Errorf("invalid key size %d bytes for key type %s", len(key), p.Type)
	}

	if p.Type == KeyType_AES128_GCM96 || p.Type == KeyType_AES256_GCM96 || p.Type == KeyType_ChaCha20_Poly1305 || p.Type == KeyType_HMAC || p.Type == KeyType_AES128_CMAC || p.Type == KeyType_AES256_CMAC {
		entry.Key = key
		if p.Type == KeyType_HMAC {
			p.KeySize = len(key)
			entry.HMACKey = key
		}
	} else {
		var parsedKey any
		var err error
		if isPrivateKey {
			parsedKey, err = x509.ParsePKCS8PrivateKey(key)
			if err != nil {
				if strings.Contains(err.Error(), "unknown elliptic curve") {
					var edErr error
					parsedKey, edErr = ParsePKCS8Ed25519PrivateKey(key)
					if edErr != nil {
						return fmt.Errorf("error parsing asymmetric key:\n - assuming contents are an ed25519 private key: %s\n - original error: %v", edErr, err)
					}

					// Parsing as Ed25519-in-PKCS8-ECPrivateKey succeeded!
				} else if strings.Contains(err.Error(), oidSignatureRSAPSS.String()) {
					var rsaErr error
					parsedKey, rsaErr = ParsePKCS8RSAPSSPrivateKey(key)
					if rsaErr != nil {
						return fmt.Errorf("error parsing asymmetric key:\n - assuming contents are an RSA/PSS private key: %v\n - original error: %w", rsaErr, err)
					}

					// Parsing as RSA-PSS in PKCS8 succeeded!
				} else {
					return fmt.Errorf("error parsing asymmetric key: %s", err)
				}
			}
		} else {
			pemBlock, _ := pem.Decode(key)
			if pemBlock == nil {
				return fmt.Errorf("error parsing public key: not in PEM format")
			}

			parsedKey, err = x509.ParsePKIXPublicKey(pemBlock.Bytes)
			if err != nil {
				return fmt.Errorf("error parsing public key: %w", err)
			}
		}

		err = entry.parseFromKey(p.Type, parsedKey)
		if err != nil {
			return err
		}
	}

	p.LatestVersion += 1

	if p.Keys == nil {
		// This is an initial key rotation when generating a new policy. We
		// don't need to call migrate here because if we've called getPolicy to
		// get the policy in the first place it will have been run.
		p.Keys = keyEntryMap{}
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

// Rotate rotates the policy and persists it to storage.
// If the rotation partially fails, the policy state will be restored.
func (p *Policy) Rotate(ctx context.Context, storage logical.Storage, randReader io.Reader) (retErr error) {
	priorLatestVersion := p.LatestVersion
	priorMinDecryptionVersion := p.MinDecryptionVersion
	var priorKeys keyEntryMap

	if p.Imported && !p.AllowImportedKeyRotation {
		return fmt.Errorf("imported key %s does not allow rotation within Vault", p.Name)
	}

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

	if err := p.RotateInMemory(randReader); err != nil {
		return err
	}

	p.Imported = false
	return p.Persist(ctx, storage)
}

// RotateInMemory rotates the policy but does not persist it to storage.
func (p *Policy) RotateInMemory(randReader io.Reader) (retErr error) {
	now := time.Now()
	entry := KeyEntry{
		CreationTime:           now,
		DeprecatedCreationTime: now.Unix(),
	}

	if p.Type != KeyType_AES128_CMAC && p.Type != KeyType_AES256_CMAC && p.Type != KeyType_HMAC {
		hmacKey, err := uuid.GenerateRandomBytesWithReader(32, randReader)
		if err != nil {
			return err
		}
		entry.HMACKey = hmacKey
	}

	var err error
	switch p.Type {
	case KeyType_AES128_GCM96, KeyType_AES256_GCM96, KeyType_ChaCha20_Poly1305, KeyType_HMAC, KeyType_AES128_CMAC, KeyType_AES256_CMAC:
		// Default to 256 bit key
		numBytes := 32
		if p.Type == KeyType_AES128_GCM96 || p.Type == KeyType_AES128_CMAC {
			numBytes = 16
		} else if p.Type == KeyType_HMAC {
			numBytes = p.KeySize
			if numBytes < HmacMinKeySize || numBytes > HmacMaxKeySize {
				return fmt.Errorf("invalid key size for HMAC key, must be between %d and %d bytes", HmacMinKeySize, HmacMaxKeySize)
			}
		}
		newKey, err := uuid.GenerateRandomBytesWithReader(numBytes, randReader)
		if err != nil {
			return err
		}
		entry.Key = newKey

		if p.Type == KeyType_HMAC {
			// To avoid causing problems, ensure HMACKey = Key.
			entry.HMACKey = newKey
		}

	case KeyType_ECDSA_P256, KeyType_ECDSA_P384, KeyType_ECDSA_P521:
		if err = generateECDSAKey(p.Type, &entry); err != nil {
			return err
		}

	case KeyType_ED25519:
		err := generateEd25519Key(randReader, &entry.Key, &entry.FormattedPublicKey)
		if err != nil {
			return err
		}
	case KeyType_RSA2048, KeyType_RSA3072, KeyType_RSA4096:
		bitSize := 2048
		if p.Type == KeyType_RSA3072 {
			bitSize = 3072
		}
		if p.Type == KeyType_RSA4096 {
			bitSize = 4096
		}

		entry.RSAKey, err = cryptoutil.GenerateRSAKey(randReader, bitSize)
		if err != nil {
			return err
		}

		entry.RSAPublicKey = entry.RSAKey.Public().(*rsa.PublicKey)

	default:
		if err := entRotateInMemory(p, &entry, randReader); err != nil {
			return err
		}
	}

	if p.ConvergentEncryption {
		if p.ConvergentVersion == -1 || p.ConvergentVersion > 1 {
			entry.ConvergentVersion = currentConvergentVersion
		}
	}

	p.LatestVersion += 1

	if p.Keys == nil {
		// This is an initial key rotation when generating a new policy. We
		// don't need to call migrate here because if we've called getPolicy to
		// get the policy in the first place it will have been run.
		p.Keys = keyEntryMap{}
	}
	p.Keys[strconv.Itoa(p.LatestVersion)] = entry

	// This ensures that with new key creations min decryption version is set
	// to 1 rather than the int default of 0, since keys start at 1 (either
	// fresh or after migration to the key map)
	if p.MinDecryptionVersion == 0 {
		p.MinDecryptionVersion = 1
	}

	return nil
}

func generateEd25519Key(randReader io.Reader, private *[]byte, public *string) error {
	// Go uses a 64-byte private key for Ed25519 keys (private+public, each
	// 32-bytes long). When we do Key derivation, we still generate a 32-byte
	// random value (and compute the corresponding Ed25519 public key), but
	// use this entire 64-byte key as if it was an HKDF key. The corresponding
	// underlying public key is never returned (which is probably good, because
	// doing so would leak half of our HKDF key...), but means we cannot import
	// derived-enabled Ed25519 public key components.
	pub, pri, err := ed25519.GenerateKey(randReader)
	if err != nil {
		return err
	}
	*private = pri
	*public = base64.StdEncoding.EncodeToString(pub)
	return nil
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

	prefix := strings.ReplaceAll(template, "{{version}}", strconv.Itoa(ver))
	p.versionPrefixCache.Store(ver, prefix)

	return prefix
}

// SymmetricOpts are the arguments to symmetric operations that are "optional", e.g.
// not always used.  This improves the aesthetics of calls to those functions.
type SymmetricOpts struct {
	// Whether to use convergent encryption
	Convergent bool
	// The version of the convergent encryption scheme
	ConvergentVersion int
	// The nonce, if not randomly generated
	Nonce []byte
	// Additional data to include in AEAD authentication
	AdditionalData []byte
	// The HMAC key, for generating IVs in convergent encryption
	HMACKey []byte
	// Allows an external provider of the AEAD, for e.g. managed keys
	AEADFactory AEADFactory
}

// Symmetrically encrypt a plaintext given the convergence configuration and appropriate keys
func (p *Policy) SymmetricEncryptRaw(ver int, encKey, plaintext []byte, opts SymmetricOpts) ([]byte, error) {
	var aead cipher.AEAD
	var err error
	nonce := opts.Nonce

	switch p.Type {
	case KeyType_AES128_GCM96, KeyType_AES256_GCM96:
		// Setup the cipher
		aesCipher, err := aes.NewCipher(encKey)
		if err != nil {
			return nil, errutil.InternalError{Err: err.Error()}
		}

		// Setup the GCM AEAD
		gcm, err := cipher.NewGCM(aesCipher)
		if err != nil {
			return nil, errutil.InternalError{Err: err.Error()}
		}

		aead = gcm

	case KeyType_ChaCha20_Poly1305:
		cha, err := chacha20poly1305.New(encKey)
		if err != nil {
			return nil, errutil.InternalError{Err: err.Error()}
		}

		aead = cha
	case KeyType_MANAGED_KEY:
		if opts.Convergent || len(opts.Nonce) != 0 {
			return nil, errutil.UserError{Err: "cannot use convergent encryption or provide a nonce to managed-key backed encryption"}
		}
		if opts.AEADFactory == nil {
			return nil, errors.New("expected AEAD factory from managed key, none provided")
		}
		aead, err = opts.AEADFactory.GetAEAD(nonce)
		if err != nil {
			return nil, err
		}
	}

	if opts.Convergent {
		convergentVersion := p.convergentVersion(ver)
		switch convergentVersion {
		case 1:
			if len(opts.Nonce) != aead.NonceSize() {
				return nil, errutil.UserError{Err: fmt.Sprintf("base64-decoded nonce must be %d bytes long when using convergent encryption with this key", aead.NonceSize())}
			}
		case 2, 3:
			if len(opts.HMACKey) == 0 {
				return nil, errutil.InternalError{Err: fmt.Sprintf("invalid hmac key length of zero")}
			}
			nonceHmac := hmac.New(sha256.New, opts.HMACKey)
			nonceHmac.Write(plaintext)
			nonceSum := nonceHmac.Sum(nil)
			nonce = nonceSum[:aead.NonceSize()]
		default:
			return nil, errutil.InternalError{Err: fmt.Sprintf("unhandled convergent version %d", convergentVersion)}
		}
	} else if len(nonce) == 0 {
		// Compute random nonce
		nonce, err = uuid.GenerateRandomBytes(aead.NonceSize())
		if err != nil {
			return nil, errutil.InternalError{Err: err.Error()}
		}
	} else if len(nonce) != aead.NonceSize() {
		return nil, errutil.UserError{Err: fmt.Sprintf("base64-decoded nonce must be %d bytes long but given %d bytes", aead.NonceSize(), len(nonce))}
	}

	// Encrypt and tag with AEAD
	ciphertext := aead.Seal(nil, nonce, plaintext, opts.AdditionalData)

	// Place the encrypted data after the nonce
	if !opts.Convergent || p.convergentVersion(ver) > 1 {
		ciphertext = append(nonce, ciphertext...)
	}
	return ciphertext, nil
}

// Symmetrically decrypt a ciphertext given the convergence configuration and appropriate keys
func (p *Policy) SymmetricDecryptRaw(encKey, ciphertext []byte, opts SymmetricOpts) ([]byte, error) {
	var aead cipher.AEAD
	var err error
	var nonce []byte

	switch p.Type {
	case KeyType_AES128_GCM96, KeyType_AES256_GCM96:
		// Setup the cipher
		aesCipher, err := aes.NewCipher(encKey)
		if err != nil {
			return nil, errutil.InternalError{Err: err.Error()}
		}

		// Setup the GCM AEAD
		gcm, err := cipher.NewGCM(aesCipher)
		if err != nil {
			return nil, errutil.InternalError{Err: err.Error()}
		}

		aead = gcm

	case KeyType_ChaCha20_Poly1305:
		cha, err := chacha20poly1305.New(encKey)
		if err != nil {
			return nil, errutil.InternalError{Err: err.Error()}
		}

		aead = cha
	case KeyType_MANAGED_KEY:
		aead, err = opts.AEADFactory.GetAEAD(nonce)
		if err != nil {
			return nil, err
		}
	}

	if len(ciphertext) < aead.NonceSize() {
		return nil, errutil.UserError{Err: "invalid ciphertext length"}
	}

	// Extract the nonce and ciphertext
	var trueCT []byte
	if opts.Convergent && opts.ConvergentVersion == 1 {
		trueCT = ciphertext
	} else {
		nonce = ciphertext[:aead.NonceSize()]
		trueCT = ciphertext[aead.NonceSize():]
	}

	// Verify and Decrypt
	plain, err := aead.Open(nil, nonce, trueCT, opts.AdditionalData)
	if err != nil {
		return nil, errutil.UserError{Err: err.Error()}
	}
	return plain, nil
}

func (p *Policy) EncryptWithFactory(ver int, context []byte, nonce []byte, value string, factories ...any) (string, error) {
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
	case KeyType_AES128_GCM96, KeyType_AES256_GCM96, KeyType_ChaCha20_Poly1305:
		hmacKey := context

		var encKey []byte
		var deriveHMAC bool

		encBytes := 32
		hmacBytes := 0
		convergentVersion := p.convergentVersion(ver)
		if convergentVersion > 2 {
			deriveHMAC = true
			hmacBytes = 32
			if len(nonce) > 0 {
				return "", errutil.UserError{Err: "nonce provided when not allowed"}
			}
		} else if len(nonce) > 0 && (!p.ConvergentEncryption || convergentVersion != 1) {
			return "", errutil.UserError{Err: "nonce provided when not allowed"}
		}
		if p.Type == KeyType_AES128_GCM96 {
			encBytes = 16
		}

		key, err := p.GetKey(context, ver, encBytes+hmacBytes)
		if err != nil {
			return "", err
		}

		if len(key) < encBytes+hmacBytes {
			return "", errutil.InternalError{Err: "could not derive key, length too small"}
		}

		encKey = key[:encBytes]
		if len(encKey) != encBytes {
			return "", errutil.InternalError{Err: "could not derive enc key, length not correct"}
		}
		if deriveHMAC {
			hmacKey = key[encBytes:]
			if len(hmacKey) != hmacBytes {
				return "", errutil.InternalError{Err: "could not derive hmac key, length not correct"}
			}
		}

		symopts := SymmetricOpts{
			Convergent: p.ConvergentEncryption,
			HMACKey:    hmacKey,
			Nonce:      nonce,
		}
		for index, rawFactory := range factories {
			if rawFactory == nil {
				continue
			}
			switch factory := rawFactory.(type) {
			case AEADFactory:
				symopts.AEADFactory = factory
			case AssociatedDataFactory:
				symopts.AdditionalData, err = factory.GetAssociatedData()
				if err != nil {
					return "", errutil.InternalError{Err: fmt.Sprintf("unable to get associated_data/additional_data from factory[%d]: %v", index, err)}
				}
			case ManagedKeyFactory:
			default:
				return "", errutil.InternalError{Err: fmt.Sprintf("unknown type of factory[%d]: %T", index, rawFactory)}
			}
		}

		ciphertext, err = p.SymmetricEncryptRaw(ver, encKey, plaintext, symopts)
		if err != nil {
			return "", err
		}
	case KeyType_RSA2048, KeyType_RSA3072, KeyType_RSA4096:
		paddingScheme, err := getPaddingScheme(factories)
		if err != nil {
			return "", err
		}
		keyEntry, err := p.safeGetKeyEntry(ver)
		if err != nil {
			return "", err
		}
		var publicKey *rsa.PublicKey
		if keyEntry.RSAKey != nil {
			publicKey = &keyEntry.RSAKey.PublicKey
		} else {
			publicKey = keyEntry.RSAPublicKey
		}
		switch paddingScheme {
		case PaddingScheme_PKCS1v15:
			ciphertext, err = rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintext)
		case PaddingScheme_OAEP:
			ciphertext, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, plaintext, nil)
		default:
			return "", errutil.InternalError{Err: fmt.Sprintf("unsupported RSA padding scheme %s", paddingScheme)}
		}

		if err != nil {
			return "", errutil.InternalError{Err: fmt.Sprintf("failed to RSA encrypt the plaintext: %v", err)}
		}
	case KeyType_MANAGED_KEY:
		keyEntry, err := p.safeGetKeyEntry(ver)
		if err != nil {
			return "", err
		}

		var aad []byte
		var managedKeyFactory ManagedKeyFactory
		for _, f := range factories {
			switch factory := f.(type) {
			case AssociatedDataFactory:
				aad, err = factory.GetAssociatedData()
				if err != nil {
					return "", nil
				}
			case ManagedKeyFactory:
				managedKeyFactory = factory
			}
		}

		if managedKeyFactory == nil {
			return "", errors.New("key type is managed_key, but managed key parameters were not provided")
		}

		ciphertext, err = p.encryptWithManagedKey(managedKeyFactory.GetManagedKeyParameters(), keyEntry, plaintext, nonce, aad)
		if err != nil {
			return "", err
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

func getPaddingScheme(factories []any) (PaddingScheme, error) {
	for _, rawFactory := range factories {
		if rawFactory == nil {
			continue
		}

		if p, ok := rawFactory.(PaddingScheme); ok && p != "" {
			return p, nil
		}
	}
	return PaddingScheme_OAEP, nil
}

func (p *Policy) KeyVersionCanBeUpdated(keyVersion int, isPrivateKey bool) error {
	keyEntry, err := p.safeGetKeyEntry(keyVersion)
	if err != nil {
		return err
	}

	if !p.Type.ImportPublicKeySupported() {
		return errors.New("provided type does not support importing key versions")
	}

	isPrivateKeyMissing := keyEntry.IsPrivateKeyMissing()
	if isPrivateKeyMissing && !isPrivateKey {
		return errors.New("cannot add a public key to a key version that already has a public key set")
	}

	if !isPrivateKeyMissing {
		return errors.New("private key imported, key version cannot be updated")
	}

	return nil
}

func (p *Policy) ImportPrivateKeyForVersion(ctx context.Context, storage logical.Storage, keyVersion int, key []byte) error {
	keyEntry, err := p.safeGetKeyEntry(keyVersion)
	if err != nil {
		return err
	}

	// Parse key
	parsedPrivateKey, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		if strings.Contains(err.Error(), "unknown elliptic curve") {
			var edErr error
			parsedPrivateKey, edErr = ParsePKCS8Ed25519PrivateKey(key)
			if edErr != nil {
				return fmt.Errorf("error parsing asymmetric key:\n - assuming contents are an ed25519 private key: %s\n - original error: %v", edErr, err)
			}

			// Parsing as Ed25519-in-PKCS8-ECPrivateKey succeeded!
		} else if strings.Contains(err.Error(), oidSignatureRSAPSS.String()) {
			var rsaErr error
			parsedPrivateKey, rsaErr = ParsePKCS8RSAPSSPrivateKey(key)
			if rsaErr != nil {
				return fmt.Errorf("error parsing asymmetric key:\n - assuming contents are an RSA/PSS private key: %v\n - original error: %w", rsaErr, err)
			}

			// Parsing as RSA-PSS in PKCS8 succeeded!
		} else {
			return fmt.Errorf("error parsing asymmetric key: %s", err)
		}
	}

	switch parsedPrivateKey.(type) {
	case *ecdsa.PrivateKey:
		ecdsaKey := parsedPrivateKey.(*ecdsa.PrivateKey)
		pemBlock, _ := pem.Decode([]byte(keyEntry.FormattedPublicKey))
		if pemBlock == nil {
			return fmt.Errorf("failed to parse key entry public key: invalid PEM blob")
		}
		publicKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
		if err != nil || publicKey == nil {
			return fmt.Errorf("failed to parse key entry public key: %v", err)
		}
		if !publicKey.(*ecdsa.PublicKey).Equal(&ecdsaKey.PublicKey) {
			return fmt.Errorf("cannot import key, key pair does not match")
		}
	case *rsa.PrivateKey:
		rsaKey := parsedPrivateKey.(*rsa.PrivateKey)
		if !rsaKey.PublicKey.Equal(keyEntry.RSAPublicKey) {
			return fmt.Errorf("cannot import key, key pair does not match")
		}
	case ed25519.PrivateKey:
		ed25519Key := parsedPrivateKey.(ed25519.PrivateKey)
		publicKey, err := base64.StdEncoding.DecodeString(keyEntry.FormattedPublicKey)
		if err != nil {
			return fmt.Errorf("failed to parse key entry public key: %v", err)
		}
		if !ed25519.PublicKey(publicKey).Equal(ed25519Key.Public()) {
			return fmt.Errorf("cannot import key, key pair does not match")
		}
	}

	err = keyEntry.parseFromKey(p.Type, parsedPrivateKey)
	if err != nil {
		return err
	}

	p.Keys[strconv.Itoa(keyVersion)] = keyEntry

	return p.Persist(ctx, storage)
}

func (ke *KeyEntry) parseFromKey(PolKeyType KeyType, parsedKey any) error {
	switch parsedKey.(type) {
	case *ecdsa.PrivateKey, *ecdsa.PublicKey:
		if PolKeyType != KeyType_ECDSA_P256 && PolKeyType != KeyType_ECDSA_P384 && PolKeyType != KeyType_ECDSA_P521 {
			return fmt.Errorf("invalid key type: expected %s, got %T", PolKeyType, parsedKey)
		}

		curve := elliptic.P256()
		if PolKeyType == KeyType_ECDSA_P384 {
			curve = elliptic.P384()
		} else if PolKeyType == KeyType_ECDSA_P521 {
			curve = elliptic.P521()
		}

		var derBytes []byte
		var err error
		ecdsaKey, ok := parsedKey.(*ecdsa.PrivateKey)
		if ok {

			if ecdsaKey.Curve != curve {
				return fmt.Errorf("invalid curve: expected %s, got %s", curve.Params().Name, ecdsaKey.Curve.Params().Name)
			}

			ke.EC_D = ecdsaKey.D
			ke.EC_X = ecdsaKey.X
			ke.EC_Y = ecdsaKey.Y

			derBytes, err = x509.MarshalPKIXPublicKey(ecdsaKey.Public())
			if err != nil {
				return errwrap.Wrapf("error marshaling public key: {{err}}", err)
			}
		} else {
			ecdsaKey := parsedKey.(*ecdsa.PublicKey)

			if ecdsaKey.Curve != curve {
				return fmt.Errorf("invalid curve: expected %s, got %s", curve.Params().Name, ecdsaKey.Curve.Params().Name)
			}

			ke.EC_X = ecdsaKey.X
			ke.EC_Y = ecdsaKey.Y

			derBytes, err = x509.MarshalPKIXPublicKey(ecdsaKey)
			if err != nil {
				return errwrap.Wrapf("error marshaling public key: {{err}}", err)
			}
		}

		pemBlock := &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: derBytes,
		}
		pemBytes := pem.EncodeToMemory(pemBlock)
		if pemBytes == nil || len(pemBytes) == 0 {
			return fmt.Errorf("error PEM-encoding public key")
		}
		ke.FormattedPublicKey = string(pemBytes)
	case ed25519.PrivateKey, ed25519.PublicKey:
		if PolKeyType != KeyType_ED25519 {
			return fmt.Errorf("invalid key type: expected %s, got %T", PolKeyType, parsedKey)
		}

		privateKey, ok := parsedKey.(ed25519.PrivateKey)
		if ok {
			ke.Key = privateKey
			publicKey := privateKey.Public().(ed25519.PublicKey)
			ke.FormattedPublicKey = base64.StdEncoding.EncodeToString(publicKey)
		} else {
			publicKey := parsedKey.(ed25519.PublicKey)
			ke.FormattedPublicKey = base64.StdEncoding.EncodeToString(publicKey)
		}
	case *rsa.PrivateKey, *rsa.PublicKey:
		if PolKeyType != KeyType_RSA2048 && PolKeyType != KeyType_RSA3072 && PolKeyType != KeyType_RSA4096 {
			return fmt.Errorf("invalid key type: expected %s, got %T", PolKeyType, parsedKey)
		}

		keyBytes := 256
		if PolKeyType == KeyType_RSA3072 {
			keyBytes = 384
		} else if PolKeyType == KeyType_RSA4096 {
			keyBytes = 512
		}

		rsaKey, ok := parsedKey.(*rsa.PrivateKey)
		if ok {
			if rsaKey.Size() != keyBytes {
				return fmt.Errorf("invalid key size: expected %d bytes, got %d bytes", keyBytes, rsaKey.Size())
			}
			ke.RSAKey = rsaKey
			ke.RSAPublicKey = rsaKey.Public().(*rsa.PublicKey)
		} else {
			rsaKey := parsedKey.(*rsa.PublicKey)
			if rsaKey.Size() != keyBytes {
				return fmt.Errorf("invalid key size: expected %d bytes, got %d bytes", keyBytes, rsaKey.Size())
			}
			ke.RSAPublicKey = rsaKey
		}
	default:
		return fmt.Errorf("invalid key type: expected %s, got %T", PolKeyType, parsedKey)
	}

	return nil
}

func (p *Policy) WrapKey(ver int, targetKey any, targetKeyType KeyType, hash hash.Hash) (string, error) {
	if !p.Type.SigningSupported() {
		return "", fmt.Errorf("message signing not supported for key type %v", p.Type)
	}

	switch {
	case ver == 0:
		ver = p.LatestVersion
	case ver < 0:
		return "", errutil.UserError{Err: "requested version for key wrapping is negative"}
	case ver > p.LatestVersion:
		return "", errutil.UserError{Err: "requested version for key wrapping is higher than the latest key version"}
	case p.MinEncryptionVersion > 0 && ver < p.MinEncryptionVersion:
		return "", errutil.UserError{Err: "requested version for key wrapping is less than the minimum encryption key version"}
	}

	keyEntry, err := p.safeGetKeyEntry(ver)
	if err != nil {
		return "", err
	}

	return keyEntry.WrapKey(targetKey, targetKeyType, hash)
}

func (ke *KeyEntry) WrapKey(targetKey any, targetKeyType KeyType, hash hash.Hash) (string, error) {
	// Presently this method implements a CKM_RSA_AES_KEY_WRAP-compatible
	// wrapping interface and only works on RSA keyEntries as a result.
	if ke.RSAPublicKey == nil {
		return "", fmt.Errorf("unsupported key type in use; must be a rsa key")
	}

	var preppedTargetKey []byte
	switch targetKeyType {
	case KeyType_AES128_GCM96, KeyType_AES256_GCM96, KeyType_ChaCha20_Poly1305, KeyType_HMAC, KeyType_AES128_CMAC, KeyType_AES256_CMAC:
		var ok bool
		preppedTargetKey, ok = targetKey.([]byte)
		if !ok {
			return "", fmt.Errorf("failed to wrap target key for import: symmetric key not provided in byte format (%T)", targetKey)
		}
	default:
		var err error
		preppedTargetKey, err = x509.MarshalPKCS8PrivateKey(targetKey)
		if err != nil {
			return "", fmt.Errorf("failed to wrap target key for import: %w", err)
		}
	}

	result, err := wrapTargetPKCS8ForImport(ke.RSAPublicKey, preppedTargetKey, hash)
	if err != nil {
		return result, fmt.Errorf("failed to wrap target key for import: %w", err)
	}

	return result, nil
}

func wrapTargetPKCS8ForImport(wrappingKey *rsa.PublicKey, preppedTargetKey []byte, hash hash.Hash) (string, error) {
	// Generate an ephemeral AES-256 key
	ephKey, err := uuid.GenerateRandomBytes(32)
	if err != nil {
		return "", fmt.Errorf("failed to generate an ephemeral AES wrapping key: %w", err)
	}

	// Wrap ephemeral AES key with public wrapping key
	ephKeyWrapped, err := rsa.EncryptOAEP(hash, rand.Reader, wrappingKey, ephKey, []byte{} /* label */)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt ephemeral wrapping key with public key: %w", err)
	}

	// Create KWP instance for wrapping target key
	kwp, err := subtle.NewKWP(ephKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate new KWP from AES key: %w", err)
	}

	// Wrap target key with KWP
	targetKeyWrapped, err := kwp.Wrap(preppedTargetKey)
	if err != nil {
		return "", fmt.Errorf("failed to wrap target key with KWP: %w", err)
	}

	// Combined wrapped keys into a single blob and base64 encode
	wrappedKeys := append(ephKeyWrapped, targetKeyWrapped...)
	return base64.StdEncoding.EncodeToString(wrappedKeys), nil
}

func (p *Policy) CreateCsr(keyVersion int, csrTemplate *x509.CertificateRequest) ([]byte, error) {
	if !p.Type.SigningSupported() {
		return nil, errutil.UserError{Err: fmt.Sprintf("key type '%s' does not support signing", p.Type)}
	}

	keyEntry, err := p.safeGetKeyEntry(keyVersion)
	if err != nil {
		return nil, err
	}

	if keyEntry.IsPrivateKeyMissing() {
		return nil, errutil.UserError{Err: "private key not imported for key version selected"}
	}

	csrTemplate.Signature = nil
	csrTemplate.SignatureAlgorithm = x509.UnknownSignatureAlgorithm

	var key crypto.Signer
	switch p.Type {
	case KeyType_ECDSA_P256, KeyType_ECDSA_P384, KeyType_ECDSA_P521:
		var curve elliptic.Curve
		switch p.Type {
		case KeyType_ECDSA_P384:
			curve = elliptic.P384()
		case KeyType_ECDSA_P521:
			curve = elliptic.P521()
		default:
			curve = elliptic.P256()
		}

		key = &ecdsa.PrivateKey{
			PublicKey: ecdsa.PublicKey{
				Curve: curve,
				X:     keyEntry.EC_X,
				Y:     keyEntry.EC_Y,
			},
			D: keyEntry.EC_D,
		}

	case KeyType_ED25519:
		if p.Derived {
			return nil, errutil.UserError{Err: "operation not supported on keys with derivation enabled"}
		}
		key = ed25519.PrivateKey(keyEntry.Key)

	case KeyType_RSA2048, KeyType_RSA3072, KeyType_RSA4096:
		key = keyEntry.RSAKey

	default:
		return nil, errutil.InternalError{Err: fmt.Sprintf("selected key type '%s' does not support signing", p.Type.String())}
	}
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, key)
	if err != nil {
		return nil, fmt.Errorf("could not create the cerfificate request: %w", err)
	}

	pemCsr := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	})

	return pemCsr, nil
}

func (p *Policy) ValidateLeafCertKeyMatch(keyVersion int, certPublicKeyAlgorithm x509.PublicKeyAlgorithm, certPublicKey any) (bool, error) {
	if !p.Type.SigningSupported() {
		return false, errutil.UserError{Err: fmt.Sprintf("key type '%s' does not support signing", p.Type)}
	}

	var keyTypeMatches bool
	switch p.Type {
	case KeyType_ECDSA_P256, KeyType_ECDSA_P384, KeyType_ECDSA_P521:
		if certPublicKeyAlgorithm == x509.ECDSA {
			keyTypeMatches = true
		}
	case KeyType_ED25519:
		if certPublicKeyAlgorithm == x509.Ed25519 {
			keyTypeMatches = true
		}
	case KeyType_RSA2048, KeyType_RSA3072, KeyType_RSA4096:
		if certPublicKeyAlgorithm == x509.RSA {
			keyTypeMatches = true
		}
	}
	if !keyTypeMatches {
		return false, errutil.UserError{Err: fmt.Sprintf("provided leaf certificate public key algorithm '%s' does not match the transit key type '%s'",
			certPublicKeyAlgorithm, p.Type)}
	}

	keyEntry, err := p.safeGetKeyEntry(keyVersion)
	if err != nil {
		return false, err
	}

	switch certPublicKeyAlgorithm {
	case x509.ECDSA:
		certPublicKey := certPublicKey.(*ecdsa.PublicKey)
		var curve elliptic.Curve
		switch p.Type {
		case KeyType_ECDSA_P384:
			curve = elliptic.P384()
		case KeyType_ECDSA_P521:
			curve = elliptic.P521()
		default:
			curve = elliptic.P256()
		}

		publicKey := &ecdsa.PublicKey{
			Curve: curve,
			X:     keyEntry.EC_X,
			Y:     keyEntry.EC_Y,
		}

		return publicKey.Equal(certPublicKey), nil

	case x509.Ed25519:
		if p.Derived {
			return false, errutil.UserError{Err: "operation not supported on keys with derivation enabled"}
		}
		certPublicKey := certPublicKey.(ed25519.PublicKey)

		raw, err := base64.StdEncoding.DecodeString(keyEntry.FormattedPublicKey)
		if err != nil {
			return false, err
		}
		publicKey := ed25519.PublicKey(raw)

		return publicKey.Equal(certPublicKey), nil

	case x509.RSA:
		certPublicKey := certPublicKey.(*rsa.PublicKey)
		publicKey := keyEntry.RSAKey.PublicKey
		return publicKey.Equal(certPublicKey), nil

	case x509.UnknownPublicKeyAlgorithm:
		return false, errutil.InternalError{Err: fmt.Sprint("certificate signed with an unknown algorithm")}
	}

	return false, nil
}

func (p *Policy) ValidateAndPersistCertificateChain(ctx context.Context, keyVersion int, certChain []*x509.Certificate, storage logical.Storage) error {
	if len(certChain) == 0 {
		return errutil.UserError{Err: "expected at least one certificate in the parsed certificate chain"}
	}

	if certChain[0].BasicConstraintsValid && certChain[0].IsCA {
		return errutil.UserError{Err: "certificate in the first position is not a leaf certificate"}
	}

	for _, cert := range certChain[1:] {
		if cert.BasicConstraintsValid && !cert.IsCA {
			return errutil.UserError{Err: "provided certificate chain contains more than one leaf certificate"}
		}
	}

	valid, err := p.ValidateLeafCertKeyMatch(keyVersion, certChain[0].PublicKeyAlgorithm, certChain[0].PublicKey)
	if err != nil {
		prefixedErr := fmt.Errorf("could not validate key match between leaf certificate key and key version in transit: %w", err)
		switch err.(type) {
		case errutil.UserError:
			return errutil.UserError{Err: prefixedErr.Error()}
		default:
			return prefixedErr
		}
	}
	if !valid {
		return fmt.Errorf("leaf certificate public key does match the key version selected")
	}

	keyEntry, err := p.safeGetKeyEntry(keyVersion)
	if err != nil {
		return err
	}

	// Convert the certificate chain to DER format
	derCertificates := make([][]byte, len(certChain))
	for i, cert := range certChain {
		derCertificates[i] = cert.Raw
	}

	keyEntry.CertificateChain = derCertificates

	p.Keys[strconv.Itoa(keyVersion)] = keyEntry
	return p.Persist(ctx, storage)
}

func generateECDSAKey(keyType KeyType, entry *KeyEntry) error {
	var curve elliptic.Curve
	switch keyType {
	case KeyType_ECDSA_P256:
		curve = elliptic.P256()
	case KeyType_ECDSA_P384:
		curve = elliptic.P384()
	case KeyType_ECDSA_P521:
		curve = elliptic.P521()
	default:
		return fmt.Errorf("invalid key type %s for ECDSA", keyType)
	}

	privKey, err := ecdsa.GenerateKey(curve, rand.Reader)
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

	return nil
}
