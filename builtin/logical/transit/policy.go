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
	"github.com/hashicorp/vault/helper/kdf"
	"github.com/hashicorp/vault/logical"
)

const (
	// kdfMode is the only KDF mode currently supported
	kdfMode = "hmac-sha256-counter"
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
	err := json.Unmarshal(data, &intermediate)
	if err != nil {
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

	// Derived keys MUST provide a context and the
	// master underlying key is never used.
	Derived bool   `json:"derived"`
	KDFMode string `json:"kdf_mode"`

	// The minimum version of the key allowed to be used
	// for decryption
	MinDecryptionVersion int `json:"min_decryption_version"`

	// The latest key version in this policy
	LatestVersion int `json:"latest_version"`

	// The latest key version in the archive. We never delete these, so this is a max.
	ArchiveVersion int `json:"archive_version"`

	// Whether the key is allowed to be deleted
	DeletionAllowed bool `json:"deletion_allowed"`
}

// ArchivedKeys stores old keys. This is used to keep the key loading time sane when
// there are huge numbers of rotations.
type ArchivedKeys struct {
	Keys []KeyEntry `json:"keys"`
}

func (p *Policy) loadArchive(storage logical.Storage, name string) (*ArchivedKeys, error) {
	archive := &ArchivedKeys{}

	raw, err := storage.Get("policy/" + name + "/archive")
	if err != nil {
		return nil, err
	}
	if raw == nil {
		archive.Keys = make([]KeyEntry, 0)
		return archive, nil
	}

	if err := json.Unmarshal(raw.Value, archive); err != nil {
		return nil, err
	}

	return archive, nil
}

func (p *Policy) storeArchive(archive *ArchivedKeys, storage logical.Storage, name string) error {
	// Encode the policy
	buf, err := json.Marshal(archive)
	if err != nil {
		return err
	}

	// Write the policy into storage
	err = storage.Put(&logical.StorageEntry{
		Key:   "policy/" + name + "/archive",
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
func (p *Policy) handleArchiving(storage logical.Storage, name string) error {
	// We need to move keys that are no longer accessible to ArchivedKeys, and keys
	// that now need to be accessible back here.
	//
	// For safety, because there isn't really a good reason to, we never delete
	// keys from the archive even when we move them back.

	// 0/1 are aliases, so don't deal with this code path unless we're past that
	if p.MinDecryptionVersion < 2 {
		return nil
	}

	// Check if we have the latest minimum version in the current set of keys
	_, keysContainsMinimum := p.Keys[p.MinDecryptionVersion]

	// If keys contains the minimum value, we are moving keys *to* the archive,
	// but we only need to do this if the archive doesn't contain those key
	// versions, since we don't remove key versions from the archive.
	if keysContainsMinimum &&
		p.ArchiveVersion >= p.MinDecryptionVersion-1 {
		return nil
	}

	archive, err := p.loadArchive(storage, name)
	if err != nil {
		return err
	}

	if keysContainsMinimum {
		// Need to move keys *to* archive

		if len(archive.Keys) < p.MinDecryptionVersion-1 {
			// Increase the size of the archive slice. We need a size that is
			// equivalent to the minimum decryption version minus 1, but adding
			// one since slice numbering starts at 0 and we're indexing by key
			// version
			newKeys := make([]KeyEntry, p.MinDecryptionVersion)
			copy(newKeys, archive.Keys)
			archive.Keys = newKeys
		}

		// As we are archiving progressively, we should only have to archive
		// from the min version down to the latest version minus however many
		// keys are in the policy's map. For example, if we have never
		// archived, the latest version is 10, and we move the min decryption
		// version to 5, we will archive from 4 down to (10-10) non-inclusive.
		// If the latest version now becomes 8, we will archive from 7 down to
		// (10-6) non-inclusive, e.g. keys 5, 6, and 7.
		for i := p.LatestVersion - len(p.Keys) + 1; i < p.MinDecryptionVersion; i++ {
			archive.Keys[i] = p.Keys[i]
		}

		err = p.storeArchive(archive, storage, name)
		if err != nil {
			return err
		}

		// Perform deletion afterwards so that if there is an error saving we
		// haven't messed with the current policy
		for i := p.LatestVersion - len(p.Keys) + 1; i < p.MinDecryptionVersion; i++ {
			delete(p.Keys, i)
		}

		// Update the archive max key version. This also corresponds to the
		// maximum safe index into the slice. Continuing our example from
		// before, p.ArchiveVersion will now be 7.
		p.ArchiveVersion = p.MinDecryptionVersion - 1

	} else {
		// Need to move keys *from* archive

		// If we've been archiving, keys should have been archived
		// sequentially. So we can perform a sanity check. First test the
		// actual latest version in the policy, so continuing the previous
		// example, if the key version is 10 and the minimum was 8, p.Keys
		// should hold 8, 9, and 10. Now if we move the minimum back, e.g. to
		// 5, we need to load keys 5, 6, and 7, so should load everything up to
		// (10-3), inclusive. If p.ArchiveVersion is less than this (which it
		// shouldn't be, as set earlier in the example), we have a problem.
		// Also, we should never have a situation where the Archive version is
		// less than the minimum decryption version but we also do not have the
		// minimum version in p.Keys (which is the only way we'd be in this
		// code path to begin with). That's also a problem.
		//
		// Note that we *should never have these problems*. If we do it's
		// serious.
		if p.ArchiveVersion < p.LatestVersion-len(p.Keys) ||
			p.ArchiveVersion < p.MinDecryptionVersion {
			return fmt.Errorf("latest archived key version not high enough to satisfy request")
		}

		for i := p.MinDecryptionVersion; i <= p.LatestVersion-len(p.Keys); i++ {
			_, ok := p.Keys[i]
			if ok {
				// We hit the beginning of the values currently in the keyset,
				// so break
				break
			}
			p.Keys[i] = archive.Keys[i]
		}
	}

	return nil
}

func (p *Policy) Persist(storage logical.Storage, name string) error {
	err := p.handleArchiving(storage, name)
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
		Key:   "policy/" + name,
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

// DeriveKey is used to derive the encryption key that should
// be used depending on the policy. If derivation is disabled the
// raw key is used and no context is required, otherwise the KDF
// mode is used with the context to derive the proper key.
func (p *Policy) DeriveKey(context []byte, ver int) ([]byte, error) {
	if p.Keys == nil || p.LatestVersion == 0 {
		if p.Key == nil || len(p.Key) == 0 {
			return nil, certutil.InternalError{Err: "unable to access the key; no key versions found"}
		}
		p.migrateKeyToKeysMap()
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
		return "", certutil.InternalError{Err: err.Error()}
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

	// Compute random nonce
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return "", certutil.InternalError{Err: err.Error()}
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
		return "", certutil.UserError{Err: "invalid ciphertext"}
	}

	splitVerCiphertext := strings.SplitN(strings.TrimPrefix(value, "vault:v"), ":", 2)
	if len(splitVerCiphertext) != 2 {
		return "", certutil.UserError{Err: "invalid ciphertext"}
	}

	ver, err := strconv.Atoi(splitVerCiphertext[0])
	if err != nil {
		return "", certutil.UserError{Err: "invalid ciphertext"}
	}

	if ver == 0 {
		// Compatibility mode with initial implementation, where keys start at zero
		ver = 1
	}

	if p.MinDecryptionVersion > 0 && ver < p.MinDecryptionVersion {
		return "", certutil.UserError{Err: "ciphertext version is disallowed by policy (too old)"}
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
		return "", certutil.UserError{Err: "invalid ciphertext"}
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
		return "", certutil.UserError{Err: "invalid ciphertext"}
	}

	return base64.StdEncoding.EncodeToString(plain), nil
}

func (p *Policy) rotate(storage logical.Storage) error {
	if p.Keys == nil {
		p.migrateKeyToKeysMap()
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

	return p.Persist(storage, p.Name)
}

func (p *Policy) migrateKeyToKeysMap() {
	if p.Key == nil || len(p.Key) == 0 {
		p.Key = nil
		p.Keys = KeyEntryMap{}
		return
	}

	p.Keys = KeyEntryMap{
		1: KeyEntry{
			Key:          p.Key,
			CreationTime: time.Now().Unix(),
		},
	}
	p.Key = nil
}

func deserializePolicy(buf []byte) (*Policy, error) {
	p := &Policy{
		Keys: KeyEntryMap{},
	}
	if err := json.Unmarshal(buf, p); err != nil {
		return nil, err
	}

	return p, nil
}

func getPolicy(req *logical.Request, name string) (*Policy, error) {
	// Check if the policy already exists
	raw, err := req.Storage.Get("policy/" + name)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}

	// Decode the policy
	p, err := deserializePolicy(raw.Value)
	if err != nil {
		return nil, err
	}

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

	if persistNeeded {
		err = p.Persist(req.Storage, name)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

// generatePolicy is used to create a new named policy with
// a randomly generated key
func generatePolicy(storage logical.Storage, name string, derived bool) (*Policy, error) {
	// Create the policy object
	p := &Policy{
		Name:       name,
		CipherMode: "aes-gcm",
		Derived:    derived,
	}
	if derived {
		p.KDFMode = kdfMode
	}

	err := p.rotate(storage)
	if err != nil {
		return nil, err
	}

	// Return the policy
	return p, nil
}
