package vault

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
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

	// tokenSubPath is the sub-path used for the token store
	// view. This is nested under the system view.
	tokenSubPath = "token/"
)

// TokenStore is used to manage client tokens. Tokens are used for
// clients to authenticate, and each token is mapped to an applicable
// set of policy which is used for authorization.
type TokenStore struct {
	*framework.Backend

	view *BarrierView
	salt string

	expiration *ExpirationManager
}

// NewTokenStore is used to construct a token store that is
// backed by the given barrier view.
func NewTokenStore(c *Core) (*TokenStore, error) {
	// Create a sub-view
	view := c.systemView.SubView(tokenSubPath)

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

	// Setup the framework endpoints
	t.Backend = &framework.Backend{
		PathsSpecial: &logical.Paths{
			Root: []string{
				"revoke-prefix/*",
			},
		},

		Paths: []*framework.Path{
			&framework.Path{
				Pattern: "create$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.WriteOperation: t.handleCreate,
				},

				HelpSynopsis:    strings.TrimSpace(tokenCreateHelp),
				HelpDescription: strings.TrimSpace(tokenCreateHelp),
			},

			&framework.Path{
				Pattern: "lookup/(?P<token>.+)",

				Fields: map[string]*framework.FieldSchema{
					"token": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Token to lookup",
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: t.handleLookup,
				},

				HelpSynopsis:    strings.TrimSpace(tokenLookupHelp),
				HelpDescription: strings.TrimSpace(tokenLookupHelp),
			},

			&framework.Path{
				Pattern: "lookup-self$",

				Fields: map[string]*framework.FieldSchema{
					"token": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Token to lookup",
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: t.handleLookup,
				},

				HelpSynopsis:    strings.TrimSpace(tokenLookupHelp),
				HelpDescription: strings.TrimSpace(tokenLookupHelp),
			},

			&framework.Path{
				Pattern: "revoke/(?P<token>.+)",

				Fields: map[string]*framework.FieldSchema{
					"token": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Token to revoke",
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.WriteOperation: t.handleRevokeTree,
				},

				HelpSynopsis:    strings.TrimSpace(tokenRevokeHelp),
				HelpDescription: strings.TrimSpace(tokenRevokeHelp),
			},

			&framework.Path{
				Pattern: "revoke-orphan/(?P<token>.+)",

				Fields: map[string]*framework.FieldSchema{
					"token": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Token to revoke",
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.WriteOperation: t.handleRevokeOrphan,
				},

				HelpSynopsis:    strings.TrimSpace(tokenRevokeOrphanHelp),
				HelpDescription: strings.TrimSpace(tokenRevokeOrphanHelp),
			},

			&framework.Path{
				Pattern: "revoke-prefix/(?P<prefix>.+)",

				Fields: map[string]*framework.FieldSchema{
					"prefix": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Token source prefix to revoke",
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.WriteOperation: t.handleRevokePrefix,
				},

				HelpSynopsis:    strings.TrimSpace(tokenRevokePrefixHelp),
				HelpDescription: strings.TrimSpace(tokenRevokePrefixHelp),
			},

			&framework.Path{
				Pattern: "renew/(?P<token>.+)",

				Fields: map[string]*framework.FieldSchema{
					"token": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Token to renew",
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.WriteOperation: t.handleRenew,
				},

				HelpSynopsis:    strings.TrimSpace(tokenRenewHelp),
				HelpDescription: strings.TrimSpace(tokenRenewHelp),
			},
		},
	}

	return t, nil
}

// TokenEntry is used to represent a given token
type TokenEntry struct {
	ID       string            // ID of this entry, generally a random UUID
	Parent   string            // Parent token, used for revocation trees
	Policies []string          // Which named policies should be used
	Path     string            // Used for audit trails, this is something like "auth/user/login"
	Meta     map[string]string // Used for auditing. This could include things like "source", "user", "ip"
}

// SetExpirationManager is used to provide the token store with
// an expiration manager. This is used to manage prefix based revocation
// of tokens and to cleanup entries when removed from the token store.
func (t *TokenStore) SetExpirationManager(exp *ExpirationManager) {
	t.expiration = exp
}

// SaltID is used to apply a salt and hash to an ID to make sure its not reversable
func (ts *TokenStore) SaltID(id string) string {
	comb := ts.salt + id
	hash := sha1.Sum([]byte(comb))
	return hex.EncodeToString(hash[:])
}

// RootToken is used to generate a new token with root privileges and no parent
func (ts *TokenStore) RootToken() (*TokenEntry, error) {
	te := &TokenEntry{
		Policies: []string{"root"},
		Path:     "sys/root",
	}
	if err := ts.Create(te); err != nil {
		return nil, err
	}
	return te, nil
}

// Create is used to create a new token entry. The entry is assigned
// a newly generated ID if not provided.
func (ts *TokenStore) Create(entry *TokenEntry) error {
	// Generate an ID if necessary
	if entry.ID == "" {
		entry.ID = generateUUID()
	}
	saltedId := ts.SaltID(entry.ID)

	// Marshal the entry
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
		path := parentPrefix + ts.SaltID(entry.Parent) + "/" + saltedId
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
	if id == "" {
		return nil, fmt.Errorf("cannot lookup blank token")
	}
	return ts.lookupSalted(ts.SaltID(id))
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
	if id == "" {
		return fmt.Errorf("cannot revoke blank token")
	}
	return ts.revokeSalted(ts.SaltID(id))
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
		path := parentPrefix + ts.SaltID(entry.Parent) + "/" + saltedId
		if ts.view.Delete(path); err != nil {
			return fmt.Errorf("failed to delete entry: %v", err)
		}
	}
	return nil
}

// RevokeTree is used to invalide a given token and all
// child tokens.
func (ts *TokenStore) RevokeTree(id string) error {
	// Verify the token is not blank
	if id == "" {
		return fmt.Errorf("cannot revoke blank token")
	}

	// Get the salted ID
	saltedId := ts.SaltID(id)

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
		path := parentPrefix + ts.SaltID(entry.Parent) + "/" + saltedId
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
func (ts *TokenStore) RevokeAll() error {
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

// handleCreate handles the auth/token/create path for creation of new tokens
func (ts *TokenStore) handleCreate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Read the parent policy
	parent, err := ts.Lookup(req.ClientToken)
	if err != nil || parent == nil {
		return logical.ErrorResponse("parent token lookup failed"), logical.ErrInvalidRequest
	}

	// Check if the parent policy is root
	isRoot := strListContains(parent.Policies, "root")

	// Read and parse the fields
	idRaw, _ := req.Data["id"]
	policiesRaw, _ := req.Data["policies"]
	metaRaw, _ := req.Data["meta"]
	noParentRaw, _ := req.Data["no_parent"]
	leaseRaw, _ := req.Data["lease"]

	// Setup the token entry
	te := TokenEntry{
		Parent: req.ClientToken,
		Path:   "auth/token/create",
	}

	// Allow specifying the ID of the token if the client is root
	if id, ok := idRaw.(string); ok {
		if !isRoot {
			return logical.ErrorResponse("root required to specify token id"),
				logical.ErrInvalidRequest
		}
		te.ID = id
	}

	// Only permit policies to be a subset unless the client is root
	if policies, ok := policiesRaw.([]string); ok {
		if !isRoot && !strListSubset(parent.Policies, policies) {
			return logical.ErrorResponse("child policies must be subset of parent"), logical.ErrInvalidRequest
		}
		te.Policies = policies
	}

	// Ensure is some associated policy
	if len(te.Policies) == 0 {
		return logical.ErrorResponse("token must have at least one policy"), logical.ErrInvalidRequest
	}

	// Only allow an orphan token if the client is root
	if noParent, _ := noParentRaw.(bool); noParent {
		if !isRoot {
			return logical.ErrorResponse("root required to create orphan token"),
				logical.ErrInvalidRequest
		}
		te.Parent = ""
	}

	// Parse any metadata associated with the token
	if meta, ok := metaRaw.(map[string]string); ok {
		te.Meta = meta
	}

	// Parse the lease if any
	var secret *logical.Secret
	if lease, ok := leaseRaw.(string); ok {
		dur, err := time.ParseDuration(lease)
		if err != nil {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		if dur < 0 {
			return logical.ErrorResponse("lease must be positive"), logical.ErrInvalidRequest
		}
		secret = &logical.Secret{
			Lease:            dur,
			LeaseGracePeriod: dur / 10, // Provide a 10% grace buffer
			Renewable:        true,
		}
	}

	// Create the token
	if err := ts.Create(&te); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Generate the response
	resp := &logical.Response{
		Secret: secret,
		Auth: &logical.Auth{
			ClientToken: te.ID,
		},
	}

	return resp, nil
}

// handleRevokeTree handles the auth/token/revoke/id path for revocation of tokens
// in a way that revokes all child tokens. Normally, using sys/revoke/vaultID will revoke
// the token and all children anyways, but that is only available when there is a lease.
func (ts *TokenStore) handleRevokeTree(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	id := data.Get("token").(string)
	if id == "" {
		return logical.ErrorResponse("missing token ID"), logical.ErrInvalidRequest
	}

	// Revoke the token and its children
	if err := ts.RevokeTree(id); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
}

// handleRevokeOrphan handles the auth/token/revoke-orphan/id path for revocation of tokens
// in a way that leaves child tokens orphaned. Normally, using sys/revoke/vaultID will revoke
// the token and all children.
func (ts *TokenStore) handleRevokeOrphan(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Parse the id
	id := data.Get("token").(string)
	if id == "" {
		return logical.ErrorResponse("missing token ID"), logical.ErrInvalidRequest
	}

	// Revoke and orphan
	if err := ts.Revoke(id); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
}

// handleRevokePrefix handles the auth/token/revoke-prefix/path for revocation of tokens
// generated by a given path.
func (ts *TokenStore) handleRevokePrefix(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Parse the prefix
	prefix := data.Get("prefix").(string)
	if prefix == "" {
		return logical.ErrorResponse("missing source prefix"), logical.ErrInvalidRequest
	}

	// Revoke using the prefix
	if err := ts.expiration.RevokePrefix(prefix); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
}

// handleLookup handles the auth/token/lookup/id path for querying information about
// a particular token. This can be used to see which policies are applicable.
func (ts *TokenStore) handleLookup(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	id := data.Get("token").(string)
	if id == "" {
		id = req.ClientToken
	}
	if id == "" {
		return logical.ErrorResponse("missing token ID"), logical.ErrInvalidRequest
	}

	// Lookup the token
	out, err := ts.Lookup(id)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Fast-path the not found case
	if out == nil {
		return nil, nil
	}

	// Generate a response. We purposely omit the parent reference otherwise
	// you could escalade your privileges.
	resp := &logical.Response{
		Data: map[string]interface{}{
			"id":       out.ID,
			"policies": out.Policies,
			"path":     out.Path,
			"meta":     out.Meta,
		},
	}
	return resp, nil
}

// handleRenew handles the auth/token/renew/id path for renewal of tokens.
// This is used to prevent token expiration and revocation.
func (ts *TokenStore) handleRenew(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	id := data.Get("token").(string)
	if id == "" {
		return logical.ErrorResponse("missing token ID"), logical.ErrInvalidRequest
	}

	// Lookup the token
	out, err := ts.Lookup(id)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Verify the token exists
	if out == nil {
		return logical.ErrorResponse("token not found"), logical.ErrInvalidRequest
	}

	// Revoke the token and its children
	if err := ts.expiration.RenewToken(out.Path, out.ID); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
}

const (
	tokenBackendHelp = `The token credential backend is always enabled and builtin to Vault.
Client tokens are used to identify a client and to allow Vault to associate policies and ACLs
which are enforced on every request. This backend also allows for generating sub-tokens as well
as revocation of tokens.`
	tokenCreateHelp       = `The token create path is used to create new tokens.`
	tokenLookupHelp       = `This endpoint will lookup a token and its properties.`
	tokenRevokeHelp       = `This endpoint will delete the token and all of its child tokens.`
	tokenRevokeOrphanHelp = `This endpoint will delete the token and orphan its child tokens.`
	tokenRevokePrefixHelp = `This endpoint will delete all tokens generated under a prefix with their child tokens.`
	tokenRenewHelp        = `This endpoint will renew the token and prevent expiration.`
)
