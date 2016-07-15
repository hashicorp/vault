package transit

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/kdf"
	"github.com/hashicorp/vault/logical"
)

const (
	// kdfMode is the only KDF mode currently supported
	kdfMode = "hmac-sha256-counter"

	ErrTooOld = "ciphertext version is disallowed by policy (too old)"
)

// KeyEntry stores the key and metadata
type KeyEntry struct {
	Key          []byte `json:"key"`
	CreationTime int64  `json:"creation_time"`
}

// KeyEntryMap is used to allow JSON marshal/unmarshal
type KeyEntryMap map[int]KeyEntry

// MarshalJSON implements JSON marshaling
func (kem KeyEntryMap) MarshalJSON() ([]byte, error) {
	intermediate := map[string]KeyEntry{}
	for k, v := range kem {
		intermediate[strconv.Itoa(k)] = v
	}
	return json.Marshal(&intermediate)
}

// MarshalJSON implements JSON unmarshaling
func (kem KeyEntryMap) UnmarshalJSON(data []byte) error {
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

// Policy is the struct used to store metadata
type Policy struct {
	Name       string      `json:"name"`
	Key        []byte      `json:"key,omitempty"` //DEPRECATED
	Keys       KeyEntryMap `json:"keys"`
	CipherMode string      `json:"cipher"`

	// Derived keys MUST provide a context and the master underlying key is
	// never used. If convergent encryption is true, the context will be used
	// as the nonce as well.
	Derived              bool   `json:"derived"`
	KDFMode              string `json:"kdf_mode"`
	ConvergentEncryption bool   `json:"convergent_encryption"`

	// The minimum version of the key allowed to be used
	// for decryption
	MinDecryptionVersion int `json:"min_decryption_version"`

	// The latest key version in this policy
	LatestVersion int `json:"latest_version"`

	// The latest key version in the archive. We never delete these, so this is
	// a max.
	ArchiveVersion int `json:"archive_version"`

	// Whether the key is allowed to be deleted
	DeletionAllowed bool `json:"deletion_allowed"`
}

// ArchivedKeys stores old keys. This is used to keep the key loading time sane
// when there are huge numbers of rotations.
type ArchivedKeys struct {
	Keys []KeyEntry `json:"keys"`
}

func (p *Policy) loadArchive(storage logical.Storage) (*ArchivedKeys, error) {
	archive := &ArchivedKeys{}

	raw, err := storage.Get("archive/" + p.Name)
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

func (p *Policy) storeArchive(archive *ArchivedKeys, storage logical.Storage) error {
	// Encode the policy
	buf, err := json.Marshal(archive)
	if err != nil {
		return err
	}

	// Write the policy into storage
	err = storage.Put(&logical.StorageEntry{
		Key:   "archive/" + p.Name,
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
func (p *Policy) handleArchiving(storage logical.Storage) error {
	// We need to move keys that are no longer accessible to ArchivedKeys, and keys
	// that now need to be accessible back here.
	//
	// For safety, because there isn't really a good reason to, we never delete
	// keys from the archive even when we move them back.

	// Check if we have the latest minimum version in the current set of keys
	_, keysContainsMinimum := p.Keys[p.MinDecryptionVersion]

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
	case p.MinDecryptionVersion > p.LatestVersion:
		return fmt.Errorf("minimum decryption version of %d is greater than the latest version %d",
			p.MinDecryptionVersion, p.LatestVersion)
	}

	archive, err := p.loadArchive(storage)
	if err != nil {
		return err
	}

	if !keysContainsMinimum {
		// Need to move keys *from* archive

		for i := p.MinDecryptionVersion; i <= p.LatestVersion; i++ {
			p.Keys[i] = archive.Keys[i]
		}

		return nil
	}

	// Need to move keys *to* archive

	// We need a size that is equivalent to the latest version (number of keys)
	// but adding one since slice numbering starts at 0 and we're indexing by
	// key version
	if len(archive.Keys) < p.LatestVersion+1 {
		// Increase the size of the archive slice
		newKeys := make([]KeyEntry, p.LatestVersion+1)
		copy(newKeys, archive.Keys)
		archive.Keys = newKeys
	}

	// We are storing all keys in the archive, so we ensure that it is up to
	// date up to p.LatestVersion
	for i := p.ArchiveVersion + 1; i <= p.LatestVersion; i++ {
		archive.Keys[i] = p.Keys[i]
		p.ArchiveVersion = i
	}

	err = p.storeArchive(archive, storage)
	if err != nil {
		return err
	}

	// Perform deletion afterwards so that if there is an error saving we
	// haven't messed with the current policy
	for i := p.LatestVersion - len(p.Keys) + 1; i < p.MinDecryptionVersion; i++ {
		delete(p.Keys, i)
	}

	return nil
}

func (p *Policy) Persist(storage logical.Storage) error {
	err := p.handleArchiving(storage)
	if err != nil {
		return err
	}

	// Encode the policy
	buf, err := p.Serialize()
	if err != nil {
		return err
	}

	// Write the policy into storage
	err = storage.Put(&logical.StorageEntry{
		Key:   "policy/" + p.Name,
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

func (p *Policy) needsUpgrade() bool {
	// Ensure we've moved from Key -> Keys
	if p.Key != nil && len(p.Key) > 0 {
		return true
	}

	// With archiving, past assumptions about the length of the keys map are no longer valid
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

	return false
}

func (p *Policy) upgrade(storage logical.Storage) error {
	persistNeeded := false
	// Ensure we've moved from Key -> Keys
	if p.Key != nil && len(p.Key) > 0 {
		p.migrateKeyToKeysMap()
		persistNeeded = true
	}

	// With archiving, past assumptions about the length of the keys map are no longer valid
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

	if persistNeeded {
		err := p.Persist(storage)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeriveKey is used to derive the encryption key that should
// be used depending on the policy. If derivation is disabled the
// raw key is used and no context is required, otherwise the KDF
// mode is used with the context to derive the proper key.
func (p *Policy) DeriveKey(context []byte, ver int) ([]byte, error) {
	if p.Keys == nil || p.LatestVersion == 0 {
		return nil, certutil.InternalError{Err: "unable to access the key; no key versions found"}
	}

	if p.LatestVersion == 0 {
		return nil, certutil.InternalError{Err: "unable to access the key; no key versions found"}
	}

	if ver <= 0 || ver > p.LatestVersion {
		return nil, certutil.UserError{Err: "invalid key version"}
	}

	// Fast-path non-derived keys
	if !p.Derived {
		return p.Keys[ver].Key, nil
	}

	// Ensure a context is provided
	if len(context) == 0 {
		return nil, certutil.UserError{Err: "missing 'context' for key deriviation. The key was created using a derived key, which means additional, per-request information must be included in order to encrypt or decrypt information"}
	}

	switch p.KDFMode {
	case kdfMode:
		prf := kdf.HMACSHA256PRF
		prfLen := kdf.HMACSHA256PRFLen
		return kdf.CounterMode(prf, prfLen, p.Keys[ver].Key, context, 256)
	default:
		return nil, certutil.InternalError{Err: "unsupported key derivation mode"}
	}
}

func (p *Policy) Encrypt(context []byte, value string) (string, error) {
	// Decode the plaintext value
	plaintext, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", certutil.UserError{Err: "failed to decode plaintext as base64"}
	}

	// Derive the key that should be used
	key, err := p.DeriveKey(context, p.LatestVersion)
	if err != nil {
		return "", err
	}

	// Guard against a potentially invalid cipher-mode
	switch p.CipherMode {
	case "aes-gcm":
	default:
		return "", certutil.InternalError{Err: "unsupported cipher mode"}
	}

	// Setup the cipher
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return "", certutil.InternalError{Err: err.Error()}
	}

	// Setup the GCM AEAD
	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return "", certutil.InternalError{Err: err.Error()}
	}

	if p.ConvergentEncryption && len(context) != gcm.NonceSize() {
		return "", certutil.UserError{Err: fmt.Sprintf("base64-decoded context must be %d bytes long when using convergent encryption with this key", gcm.NonceSize())}
	}

	// Compute random nonce
	var nonce []byte
	if p.ConvergentEncryption {
		nonce = context
	} else {
		nonce = make([]byte, gcm.NonceSize())
		_, err = rand.Read(nonce)
		if err != nil {
			return "", certutil.InternalError{Err: err.Error()}
		}
	}

	// Encrypt and tag with GCM
	out := gcm.Seal(nil, nonce, plaintext, nil)

	// Place the encrypted data after the nonce
	full := append(nonce, out...)

	// Convert to base64
	encoded := base64.StdEncoding.EncodeToString(full)

	// Prepend some information
	encoded = "vault:v" + strconv.Itoa(p.LatestVersion) + ":" + encoded

	return encoded, nil
}

func (p *Policy) Decrypt(context []byte, value string) (string, error) {
	// Verify the prefix
	if !strings.HasPrefix(value, "vault:v") {
		return "", certutil.UserError{Err: "invalid ciphertext: no prefix"}
	}

	splitVerCiphertext := strings.SplitN(strings.TrimPrefix(value, "vault:v"), ":", 2)
	if len(splitVerCiphertext) != 2 {
		return "", certutil.UserError{Err: "invalid ciphertext: wrong number of fields"}
	}

	ver, err := strconv.Atoi(splitVerCiphertext[0])
	if err != nil {
		return "", certutil.UserError{Err: "invalid ciphertext: version number could not be decoded"}
	}

	if ver == 0 {
		// Compatibility mode with initial implementation, where keys start at
		// zero
		ver = 1
	}

	if ver > p.LatestVersion {
		return "", certutil.UserError{Err: "invalid ciphertext: version is too new"}
	}

	if p.MinDecryptionVersion > 0 && ver < p.MinDecryptionVersion {
		return "", certutil.UserError{Err: ErrTooOld}
	}

	// Derive the key that should be used
	key, err := p.DeriveKey(context, ver)
	if err != nil {
		return "", err
	}

	// Guard against a potentially invalid cipher-mode
	switch p.CipherMode {
	case "aes-gcm":
	default:
		return "", certutil.InternalError{Err: "unsupported cipher mode"}
	}

	// Decode the base64
	decoded, err := base64.StdEncoding.DecodeString(splitVerCiphertext[1])
	if err != nil {
		return "", certutil.UserError{Err: "invalid ciphertext: could not decode base64"}
	}

	// Setup the cipher
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return "", certutil.InternalError{Err: err.Error()}
	}

	// Setup the GCM AEAD
	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return "", certutil.InternalError{Err: err.Error()}
	}

	// Extract the nonce and ciphertext
	nonce := decoded[:gcm.NonceSize()]
	ciphertext := decoded[gcm.NonceSize():]

	// Verify and Decrypt
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", certutil.UserError{Err: "invalid ciphertext: unable to decrypt"}
	}

	return base64.StdEncoding.EncodeToString(plain), nil
}

func (p *Policy) rotate(storage logical.Storage) error {
	if p.Keys == nil {
		// This is an initial key rotation when generating a new policy. We
		// don't need to call migrate here because if we've called getPolicy to
		// get the policy in the first place it will have been run.
		p.Keys = KeyEntryMap{}
	}

	// Generate a 256bit key
	newKey := make([]byte, 32)
	_, err := rand.Read(newKey)
	if err != nil {
		return err
	}

	p.LatestVersion += 1

	p.Keys[p.LatestVersion] = KeyEntry{
		Key:          newKey,
		CreationTime: time.Now().Unix(),
	}

	// This ensures that with new key creations min decryption version is set
	// to 1 rather than the int default of 0, since keys start at 1 (either
	// fresh or after migration to the key map)
	if p.MinDecryptionVersion == 0 {
		p.MinDecryptionVersion = 1
	}

	//fmt.Printf("policy %s rotated to %d\n", p.Name, p.LatestVersion)

	return p.Persist(storage)
}

func (p *Policy) migrateKeyToKeysMap() {
	p.Keys = KeyEntryMap{
		1: KeyEntry{
			Key:          p.Key,
			CreationTime: time.Now().Unix(),
		},
	}
	p.Key = nil
}
