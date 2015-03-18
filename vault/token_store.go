package vault

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/vault/logical"
)

const (
	// lookupPrefix is the prefix used to store tokens for their
	// primary ID based index
	lookupPrefix = "id/"

	// parentPrefix is the prefix used to store tokens for their
	// secondar parent based index
	parentPrefix = "parent/"

	// tokenSaltLocation is the path in the view we store our key salt.
	// This is used to ensure the paths we write out are obfuscated so
	// that token names cannot be guessed as that would compromise their
	// use.
	tokenSaltLocation = "salt"
)

// TokenStore is used to manage client tokens. Tokens are used for
// clients to authenticate, and each token is mapped to an applicable
// set of policy which is used for authorization.
type TokenStore struct {
	view *BarrierView
	salt string
}

// NewTokenStore is used to construct a token store that is
// backed by the given barrier view.
func NewTokenStore(view *BarrierView) (*TokenStore, error) {
	// Initialize the store
	t := &TokenStore{
		view: view,
	}

	// Look for the salt
	raw, err := view.Get(tokenSaltLocation)
	if err != nil {
		return nil, fmt.Errorf("failed to read salt: %v", err)
	}

	// Restore the salt if it exists
	if raw != nil {
		t.salt = string(raw.Value)
	}

	// Generate a new salt if necessary
	if t.salt == "" {
		t.salt = generateUUID()
		raw = &logical.StorageEntry{Key: tokenSaltLocation, Value: []byte(t.salt)}
		if err := view.Put(raw); err != nil {
			return nil, fmt.Errorf("failed to persist salt: %v", err)
		}
	}
	return t, nil
}

// TokenEntry is used to represent a given token
type TokenEntry struct {
	ID       string   // ID of this entry, generally a random UUID
	Parent   string   // Parent token, used for revocation trees
	Source   string   // Used for audit trails, this is something like "source:github.com user:armon"
	Policies []string // Which named policies should be used
}

// saltID is used to apply a salt and hash to an ID to make sure its not reversable
func (ts *TokenStore) saltID(id string) string {
	comb := ts.salt + id
	hash := sha1.Sum([]byte(comb))
	return hex.EncodeToString(hash[:])
}

// Create is used to create a new token entry. The entry is assigned
// a newly generated ID
func (ts *TokenStore) Create(entry *TokenEntry) error {
	// Marshal the entry
	entry.ID = generateUUID()
	saltedId := ts.saltID(entry.ID)
	enc, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to encode entry: %v", err)
	}

	// Write the secondary index if necessary. This is done before the
	// primary index because we'd rather have a dangling pointer with
	// a missing primary instead of missing the parent index and potentially
	// escaping the revocation chain.
	if entry.Parent != "" {
		// Ensure the parent exists
		parent, err := ts.Lookup(entry.Parent)
		if err != nil {
			return fmt.Errorf("failed to lookup parent: %v", err)
		}
		if parent == nil {
			return fmt.Errorf("parent token not found")
		}

		// Create the index entry
		path := parentPrefix + ts.saltID(entry.Parent) + "/" + saltedId
		le := &logical.StorageEntry{Key: path}
		if err := ts.view.Put(le); err != nil {
			return fmt.Errorf("failed to persist entry: %v", err)
		}
	}

	// Write the primary ID
	path := lookupPrefix + saltedId
	le := &logical.StorageEntry{Key: path, Value: enc}
	if err := ts.view.Put(le); err != nil {
		return fmt.Errorf("failed to persist entry: %v", err)
	}
	return nil
}

// Lookup is used to find a token given its ID
func (ts *TokenStore) Lookup(id string) (*TokenEntry, error) {
	return ts.lookupSalted(ts.saltID(id))
}

// lookupSlated is used to find a token given its salted ID
func (ts *TokenStore) lookupSalted(saltedId string) (*TokenEntry, error) {
	// Lookup token
	path := lookupPrefix + saltedId
	raw, err := ts.view.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read entry: %v", err)
	}

	// Bail if not found
	if raw == nil {
		return nil, nil
	}

	// Unmarshal the token
	entry := new(TokenEntry)
	if err := json.Unmarshal(raw.Value, entry); err != nil {
		return nil, fmt.Errorf("failed to decode entry: %v", err)
	}
	return entry, nil
}

// Revoke is used to invalidate a given token, any child tokens
// will be orphaned.
func (ts *TokenStore) Revoke(id string) error {
	return ts.revokeSalted(ts.saltID(id))
}

// revokeSalted is used to invalidate a given salted token,
// any child tokens will be orphaned.
func (ts *TokenStore) revokeSalted(saltedId string) error {
	// Lookup the token first
	entry, err := ts.lookupSalted(saltedId)
	if err != nil {
		return err
	}

	// Nuke the primary key first
	path := lookupPrefix + saltedId
	if ts.view.Delete(path); err != nil {
		return fmt.Errorf("failed to delete entry: %v", err)
	}

	// Clear the secondary index if any
	if entry != nil && entry.Parent != "" {
		path := parentPrefix + ts.saltID(entry.Parent) + "/" + saltedId
		if ts.view.Delete(path); err != nil {
			return fmt.Errorf("failed to delete entry: %v", err)
		}
	}
	return nil
}

// RevokeTree is used to invalide a given token and all
// child tokens.
func (ts *TokenStore) RevokeTree(id string) error {
	// Get the salted ID
	saltedId := ts.saltID(id)

	// Lookup the token first
	entry, err := ts.lookupSalted(saltedId)
	if err != nil {
		return err
	}

	// Nuke the child entries recursively
	if err := ts.revokeTreeSalted(saltedId); err != nil {
		return err
	}

	// Clear the secondary index if any
	if entry != nil && entry.Parent != "" {
		path := parentPrefix + ts.saltID(entry.Parent) + "/" + saltedId
		if ts.view.Delete(path); err != nil {
			return fmt.Errorf("failed to delete entry: %v", err)
		}
	}
	return nil
}

// revokeTreeSalted is used to invalide a given token and all
// child tokens using a saltedID.
func (ts *TokenStore) revokeTreeSalted(saltedId string) error {
	// Scan for child tokens
	path := parentPrefix + saltedId + "/"
	children, err := ts.view.List(path)
	if err != nil {
		return fmt.Errorf("failed to scan for children: %v", err)
	}

	// Recursively nuke the children. The subtle nuance here is that
	// we don't have the acutal ID of the child, but we have the salted
	// value. Turns out, this is good enough!
	for _, child := range children {
		if err := ts.revokeTreeSalted(child); err != nil {
			return fmt.Errorf("failed to revoke child: %v", err)
		}
		childIndex := path + child
		if err := ts.view.Delete(childIndex); err != nil {
			return fmt.Errorf("failed to delete child index: %v", err)
		}
	}

	// Nuke the primary key
	path = lookupPrefix + saltedId
	if ts.view.Delete(path); err != nil {
		return fmt.Errorf("failed to delete entry: %v", err)
	}
	return nil
}

// RevokeAll is used to invalidate all generated tokens.
func (ts *TokenStore) RevokeAll(entry *TokenEntry) error {
	// Collect all the tokens
	sub := ts.view.SubView(lookupPrefix)
	tokens, err := CollectKeys(sub)
	if err != nil {
		return fmt.Errorf("failed to scan tokens: %v", err)
	}

	// Invalidate them all, note that the keys we get back from the
	// sub-view are all salted
	for idx, token := range tokens {
		if err := ts.revokeSalted(token); err != nil {
			return fmt.Errorf("failed to revoke '%s' (%d / %d): %v",
				token, idx+1, len(tokens), err)
		}
	}
	return nil
}
