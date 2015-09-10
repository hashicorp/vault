package vault

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/helper/uuid"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/mitchellh/mapstructure"
)

const (
	// lookupPrefix is the prefix used to store tokens for their
	// primary ID based index
	lookupPrefix = "id/"

	// parentPrefix is the prefix used to store tokens for their
	// secondar parent based index
	parentPrefix = "parent/"

	// tokenSubPath is the sub-path used for the token store
	// view. This is nested under the system view.
	tokenSubPath = "token/"
)

var (
	// displayNameSanitize is used to sanitize a display name given to a token.
	displayNameSanitize = regexp.MustCompile("[^a-zA-Z0-9-]")
)

// TokenStore is used to manage client tokens. Tokens are used for
// clients to authenticate, and each token is mapped to an applicable
// set of policy which is used for authorization.
type TokenStore struct {
	*framework.Backend

	view *BarrierView
	salt *salt.Salt

	expiration *ExpirationManager
}

// NewTokenStore is used to construct a token store that is
// backed by the given barrier view.
func NewTokenStore(c *Core, config *logical.BackendConfig) (*TokenStore, error) {
	// Create a sub-view
	view := c.systemBarrierView.SubView(tokenSubPath)

	// Initialize the store
	t := &TokenStore{
		view: view,
	}

	// Setup the salt
	salt, err := salt.NewSalt(view, nil)
	if err != nil {
		return nil, err
	}
	t.salt = salt

	// Setup the framework endpoints
	t.Backend = &framework.Backend{
		// Allow a token lease to be extended indefinitely, but each time for only
		// as much as the original lease allowed for. If the lease has a 1 hour expiration,
		// it can only be extended up to another hour each time this means.
		AuthRenew: framework.LeaseExtend(0, 0, true),

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
					"increment": &framework.FieldSchema{
						Type:        framework.TypeDurationSecond,
						Description: "The desired increment in seconds to the token expiration",
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

	t.Backend.Setup(config)

	return t, nil
}

// TokenEntry is used to represent a given token
type TokenEntry struct {
	ID          string            // ID of this entry, generally a random UUID
	Parent      string            // Parent token, used for revocation trees
	Policies    []string          // Which named policies should be used
	Path        string            // Used for audit trails, this is something like "auth/user/login"
	Meta        map[string]string // Used for auditing. This could include things like "source", "user", "ip"
	DisplayName string            // Used for operators to be able to associate with the source
	NumUses     int               // Used to restrict the number of uses (zero is unlimited). This is to support one-time-tokens (generalized).
}

// SetExpirationManager is used to provide the token store with
// an expiration manager. This is used to manage prefix based revocation
// of tokens and to cleanup entries when removed from the token store.
func (t *TokenStore) SetExpirationManager(exp *ExpirationManager) {
	t.expiration = exp
}

// SaltID is used to apply a salt and hash to an ID to make sure its not reversable
func (ts *TokenStore) SaltID(id string) string {
	return ts.salt.SaltID(id)
}

// RootToken is used to generate a new token with root privileges and no parent
func (ts *TokenStore) RootToken() (*TokenEntry, error) {
	te := &TokenEntry{
		Policies:    []string{"root"},
		Path:        "auth/token/root",
		DisplayName: "root",
	}
	if err := ts.Create(te); err != nil {
		return nil, err
	}
	return te, nil
}

// Create is used to create a new token entry. The entry is assigned
// a newly generated ID if not provided.
func (ts *TokenStore) Create(entry *TokenEntry) error {
	defer metrics.MeasureSince([]string{"token", "create"}, time.Now())
	// Generate an ID if necessary
	if entry.ID == "" {
		entry.ID = uuid.GenerateUUID()
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

// UseToken is used to manage restricted use tokens and decrement
// their available uses.
func (ts *TokenStore) UseToken(te *TokenEntry) error {
	// If the token is not restricted, there is nothing to do
	if te.NumUses == 0 {
		return nil
	}

	// Decrement the count
	te.NumUses -= 1

	// Revoke the token if there are no remaining uses.
	// XXX: There is a race condition here with parallel
	// requests using the same token. This would require
	// some global coordination to avoid, as we must ensure
	// no requests using the same restricted token are handled
	// in parallel.
	if te.NumUses == 0 {
		return ts.Revoke(te.ID)
	}

	// Marshal the entry
	enc, err := json.Marshal(te)
	if err != nil {
		return fmt.Errorf("failed to encode entry: %v", err)
	}

	// Write under the primary ID
	saltedId := ts.SaltID(te.ID)
	path := lookupPrefix + saltedId
	le := &logical.StorageEntry{Key: path, Value: enc}
	if err := ts.view.Put(le); err != nil {
		return fmt.Errorf("failed to persist entry: %v", err)
	}
	return nil
}

// Lookup is used to find a token given its ID
func (ts *TokenStore) Lookup(id string) (*TokenEntry, error) {
	defer metrics.MeasureSince([]string{"token", "lookup"}, time.Now())
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
	defer metrics.MeasureSince([]string{"token", "revoke"}, time.Now())
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

	// Revoke all secrets under this token
	if entry != nil {
		if err := ts.expiration.RevokeByToken(entry.ID); err != nil {
			return err
		}
	}
	return nil
}

// RevokeTree is used to invalide a given token and all
// child tokens.
func (ts *TokenStore) RevokeTree(id string) error {
	defer metrics.MeasureSince([]string{"token", "revoke-tree"}, time.Now())
	// Verify the token is not blank
	if id == "" {
		return fmt.Errorf("cannot revoke blank token")
	}

	// Get the salted ID
	saltedId := ts.SaltID(id)

	// Nuke the entire tree recursively
	if err := ts.revokeTreeSalted(saltedId); err != nil {
		return err
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
			return err
		}
	}

	// Revoke this entry
	if err := ts.revokeSalted(saltedId); err != nil {
		return fmt.Errorf("failed to revoke entry: %v", err)
	}
	return nil
}

// handleCreate handles the auth/token/create path for creation of new tokens
func (ts *TokenStore) handleCreate(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Read the parent policy
	parent, err := ts.Lookup(req.ClientToken)
	if err != nil || parent == nil {
		return logical.ErrorResponse("parent token lookup failed"), logical.ErrInvalidRequest
	}

	// A token with a restricted number of uses cannot create a new token
	// otherwise it could escape the restriction count.
	if parent.NumUses > 0 {
		return logical.ErrorResponse("restricted use token cannot generate child tokens"),
			logical.ErrInvalidRequest
	}

	// Check if the parent policy is root
	isRoot := strListContains(parent.Policies, "root")

	// Read and parse the fields
	var data struct {
		ID          string
		Policies    []string
		Metadata    map[string]string `mapstructure:"meta"`
		NoParent    bool              `mapstructure:"no_parent"`
		Lease       string
		DisplayName string `mapstructure:"display_name"`
		NumUses     int    `mapstructure:"num_uses"`
	}
	if err := mapstructure.WeakDecode(req.Data, &data); err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error decoding request: %s", err)), logical.ErrInvalidRequest
	}

	// Verify the number of uses is positive
	if data.NumUses < 0 {
		return logical.ErrorResponse("number of uses cannot be negative"),
			logical.ErrInvalidRequest
	}

	// Setup the token entry
	te := TokenEntry{
		Parent:      req.ClientToken,
		Path:        "auth/token/create",
		Meta:        data.Metadata,
		DisplayName: "token",
		NumUses:     data.NumUses,
	}

	// Attach the given display name if any
	if data.DisplayName != "" {
		full := "token-" + data.DisplayName
		full = displayNameSanitize.ReplaceAllString(full, "-")
		full = strings.TrimSuffix(full, "-")
		te.DisplayName = full
	}

	// Allow specifying the ID of the token if the client is root
	if data.ID != "" {
		if !isRoot {
			return logical.ErrorResponse("root required to specify token id"),
				logical.ErrInvalidRequest
		}
		te.ID = data.ID
	}

	// Only permit policies to be a subset unless the client is root
	if len(data.Policies) == 0 {
		data.Policies = parent.Policies
	}
	if !isRoot && !strListSubset(parent.Policies, data.Policies) {
		return logical.ErrorResponse("child policies must be subset of parent"), logical.ErrInvalidRequest
	}
	te.Policies = data.Policies

	// Only allow an orphan token if the client is root
	if data.NoParent {
		if !isRoot {
			return logical.ErrorResponse("root required to create orphan token"),
				logical.ErrInvalidRequest
		}

		te.Parent = ""
	}

	// Parse the lease if any
	var leaseDuration time.Duration
	if data.Lease != "" {
		dur, err := time.ParseDuration(data.Lease)
		if err != nil {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		if dur < 0 {
			return logical.ErrorResponse("lease must be positive"), logical.ErrInvalidRequest
		}
		leaseDuration = dur
	}

	// Create the token
	if err := ts.Create(&te); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Generate the response
	resp := &logical.Response{
		Auth: &logical.Auth{
			DisplayName: te.DisplayName,
			Policies:    te.Policies,
			Metadata:    te.Meta,
			LeaseOptions: logical.LeaseOptions{
				TTL:         leaseDuration,
				GracePeriod: leaseDuration / 10,
				Renewable:   leaseDuration > 0,
			},
			ClientToken: te.ID,
		},
	}

	return resp, nil
}

// handleRevokeTree handles the auth/token/revoke/id path for revocation of tokens
// in a way that revokes all child tokens. Normally, using sys/revoke/leaseID will revoke
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
// in a way that leaves child tokens orphaned. Normally, using sys/revoke/leaseID will revoke
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

	if out == nil {
		return logical.ErrorResponse("bad token"), logical.ErrPermissionDenied
	}

	// Generate a response. We purposely omit the parent reference otherwise
	// you could escalade your privileges.
	resp := &logical.Response{
		Data: map[string]interface{}{
			"id":           out.ID,
			"policies":     out.Policies,
			"path":         out.Path,
			"meta":         out.Meta,
			"display_name": out.DisplayName,
			"num_uses":     out.NumUses,
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
	incrementRaw := data.Get("increment").(int)

	// Convert the increment
	increment := time.Duration(incrementRaw) * time.Second

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
	auth, err := ts.expiration.RenewToken(out.Path, out.ID, increment)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Generate the response
	resp := &logical.Response{
		Auth: auth,
	}
	return resp, nil
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
