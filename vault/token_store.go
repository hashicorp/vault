package vault

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/fatih/structs"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/salt"
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

	// prefixSanitize is used to ensure a prefix in a role is valid.
	prefixSanitize = regexp.MustCompile("\\w[\\w-.]+\\w")
)

// TokenStore is used to manage client tokens. Tokens are used for
// clients to authenticate, and each token is mapped to an applicable
// set of policy which is used for authorization.
type TokenStore struct {
	*framework.Backend

	view *BarrierView
	salt *salt.Salt

	expiration *ExpirationManager

	cubbyholeBackend *CubbyholeBackend

	policyLookupFunc func(string) (*Policy, error)
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

	if c.policyStore != nil {
		t.policyLookupFunc = c.policyStore.GetPolicy
	}

	// Setup the salt
	salt, err := salt.NewSalt(view, &salt.Config{
		HashFunc: salt.SHA1Hash,
	})
	if err != nil {
		return nil, err
	}
	t.salt = salt

	// Setup the framework endpoints
	t.Backend = &framework.Backend{
		AuthRenew: t.authRenew,

		PathsSpecial: &logical.Paths{
			Root: []string{
				"revoke-prefix/*",
				"revoke-orphan/*",
			},
		},

		Paths: []*framework.Path{
			&framework.Path{
				Pattern: "roles/?$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ListOperation: t.tokenStoreRoleList,
				},

				//				HelpSynopsis:    pathRoleHelpSyn,
				//				HelpDescription: pathRoleHelpDesc,
			},

			&framework.Path{
				Pattern: "roles/" + framework.GenericNameRegex("role_name"),
				Fields: map[string]*framework.FieldSchema{
					"role_name": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Name of the role",
					},

					"allowed_policies": &framework.FieldSchema{
						Type:    framework.TypeString,
						Default: "",
						Description: `If set, tokens created via this role
can be created with any subset of this list,
rather than the normal semantics of a subset
of the client token's policies. This
parameter should be sent as a comma-delimited
string.`,
					},

					"orphan": &framework.FieldSchema{
						Type:    framework.TypeBool,
						Default: false,
						Description: `If true, tokens created via this role
will be orphan tokens (have no parent)`,
					},

					"period": &framework.FieldSchema{
						Type:    framework.TypeDurationSecond,
						Default: 0,
						Description: `If set, tokens created via this role
will have no max lifetime; instead, their
renewal period will be fixed to this value.
This takes an integer number of seconds,
or a string duration (e.g. "24h").`,
					},

					"prefix": &framework.FieldSchema{
						Type:    framework.TypeString,
						Default: "",
						Description: `If set, tokens created via this role
will contain the given prefix as a part of
their path. This can be used to assist use
of the 'revoke-prefix' endpoint later on.
The given prefix must match the regular
expression ` + prefixSanitize.String(),
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation:   t.tokenStoreRoleRead,
					logical.UpdateOperation: t.tokenStoreRoleCreate,
					logical.DeleteOperation: t.tokenStoreRoleDelete,
				},

				//				HelpSynopsis:    pathRoleHelpSyn,
				//				HelpDescription: pathRoleHelpDesc,
			},

			&framework.Path{
				Pattern: "create-orphan/$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: t.handleCreateOrphan,
				},

				HelpSynopsis:    strings.TrimSpace(tokenCreateOrphanHelp),
				HelpDescription: strings.TrimSpace(tokenCreateOrphanHelp),
			},

			&framework.Path{
				Pattern: "create/" + framework.GenericNameRegex("role_name"),

				Fields: map[string]*framework.FieldSchema{
					"role_name": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Name of the role",
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: t.handleCreateRole,
				},

				//HelpSynopsis:    strings.TrimSpace(tokenCreateOrphanHelp),
				//HelpDescription: strings.TrimSpace(tokenCreateOrphanHelp),
			},

			&framework.Path{
				Pattern: "create$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: t.handleCreate,
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
				Pattern: "revoke-self$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: t.handleRevokeSelf,
				},

				HelpSynopsis:    strings.TrimSpace(tokenRevokeSelfHelp),
				HelpDescription: strings.TrimSpace(tokenRevokeSelfHelp),
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
					logical.UpdateOperation: t.handleRevokeTree,
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
					logical.UpdateOperation: t.handleRevokeOrphan,
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
					logical.UpdateOperation: t.handleRevokePrefix,
				},

				HelpSynopsis:    strings.TrimSpace(tokenRevokePrefixHelp),
				HelpDescription: strings.TrimSpace(tokenRevokePrefixHelp),
			},

			&framework.Path{
				Pattern: "renew-self$",

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
					logical.UpdateOperation: t.handleRenewSelf,
				},

				HelpSynopsis:    strings.TrimSpace(tokenRenewSelfHelp),
				HelpDescription: strings.TrimSpace(tokenRenewSelfHelp),
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
					logical.UpdateOperation: t.handleRenew,
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
	ID           string            // ID of this entry, generally a random UUID
	Parent       string            // Parent token, used for revocation trees
	Policies     []string          // Which named policies should be used
	Path         string            // Used for audit trails, this is something like "auth/user/login"
	Meta         map[string]string // Used for auditing. This could include things like "source", "user", "ip"
	DisplayName  string            // Used for operators to be able to associate with the source
	NumUses      int               // Used to restrict the number of uses (zero is unlimited). This is to support one-time-tokens (generalized).
	CreationTime int64             // Time of token creation
	TTL          time.Duration     // Duration set when token was created
	Role         string            // If set, the role that was used for parameters at creation time
}

// tsRoleEntry contains token store role information
type tsRoleEntry struct {
	// The name of the role. Embedded so it can be used for pathing
	Name string `json:"name" mapstructure:"name" structs:"name"`

	// The policies that creation functions using this role can assign to a token,
	// escaping or further locking down normal subset checking
	AllowedPolicies []string `json:"allowed_policies" mapstructure:"allowed_policies" structs:"allowed_policies"`

	// If true, tokens created using this role will be orphans
	Orphan bool `json:"orphan" mapstructure:"orphan" structs:"orphan"`

	// If non-zero, tokens created using this role will be able to be renewed
	// forever, but will have a fixed renewal period of this value
	Period time.Duration `json:"period" mapstructure:"period" structs:"period"`

	// If set, a prefix will be set on leases, making it easier to revoke using
	// 'revoke-prefix'.
	Prefix string `json:"prefix" mapstructure:"prefix" structs:"prefix"`
}

// SetExpirationManager is used to provide the token store with
// an expiration manager. This is used to manage prefix based revocation
// of tokens and to cleanup entries when removed from the token store.
func (ts *TokenStore) SetExpirationManager(exp *ExpirationManager) {
	ts.expiration = exp
}

// SaltID is used to apply a salt and hash to an ID to make sure its not reversable
func (ts *TokenStore) SaltID(id string) string {
	return ts.salt.SaltID(id)
}

// RootToken is used to generate a new token with root privileges and no parent
func (ts *TokenStore) rootToken() (*TokenEntry, error) {
	te := &TokenEntry{
		Policies:     []string{"root"},
		Path:         "auth/token/root",
		DisplayName:  "root",
		CreationTime: time.Now().Unix(),
	}
	if err := ts.create(te); err != nil {
		return nil, err
	}
	return te, nil
}

// Create is used to create a new token entry. The entry is assigned
// a newly generated ID if not provided.
func (ts *TokenStore) create(entry *TokenEntry) error {
	defer metrics.MeasureSince([]string{"token", "create"}, time.Now())
	// Generate an ID if necessary
	if entry.ID == "" {
		entryUUID, err := uuid.GenerateUUID()
		if err != nil {
			return err
		}
		entry.ID = entryUUID
	}

	return ts.storeCommon(entry, true)
}

// Store is used to store an updated token entry without writing the
// secondary index.
func (ts *TokenStore) store(entry *TokenEntry) error {
	defer metrics.MeasureSince([]string{"token", "store"}, time.Now())
	return ts.storeCommon(entry, false)
}

// storeCommon handles the actual storage of an entry, possibly generating
// secondary indexes
func (ts *TokenStore) storeCommon(entry *TokenEntry, writeSecondary bool) error {
	saltedId := ts.SaltID(entry.ID)

	// Marshal the entry
	enc, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to encode entry: %v", err)
	}

	if writeSecondary {
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

	// Destroy the cubby space
	err = ts.destroyCubbyhole(saltedId)
	if err != nil {
		return err
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

// handleCreateRole handles the auth/token/create path for a role
func (ts *TokenStore) handleCreateRole(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("role_name").(string)
	roleEntry, err := ts.getTokenStoreRole(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown role %s", name)), nil
	}

	return ts.handleCreateCommon(req, d, false, roleEntry)
}

// handleCreate handles the auth/token/create path for creation of new orphan
// tokens
func (ts *TokenStore) handleCreateOrphan(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return ts.handleCreateCommon(req, d, true, nil)
}

// handleCreate handles the auth/token/create path for creation of new non-orphan
// tokens
func (ts *TokenStore) handleCreate(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return ts.handleCreateCommon(req, d, false, nil)
}

// handleCreateCommon handles the auth/token/create path for creation of new tokens
func (ts *TokenStore) handleCreateCommon(
	req *logical.Request, d *framework.FieldData, orphan bool, role *tsRoleEntry) (*logical.Response, error) {
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

	// Check if the client token has sudo/root privileges for the requested path
	isSudo := ts.System().SudoPrivilege(req.MountPoint+req.Path, req.ClientToken)

	// Read and parse the fields
	var data struct {
		ID              string
		Policies        []string
		Metadata        map[string]string `mapstructure:"meta"`
		NoParent        bool              `mapstructure:"no_parent"`
		NoDefaultPolicy bool              `mapstructure:"no_default_policy"`
		Lease           string
		TTL             string
		DisplayName     string `mapstructure:"display_name"`
		NumUses         int    `mapstructure:"num_uses"`
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
		Parent:       req.ClientToken,
		Path:         "auth/token/create",
		Meta:         data.Metadata,
		DisplayName:  "token",
		NumUses:      data.NumUses,
		CreationTime: time.Now().Unix(),
	}

	if role != nil {
		te.Role = role.Name

		te.Path = fmt.Sprintf("%s/%s", te.Path, role.Name)
		if role.Prefix != "" {
			te.Path = fmt.Sprintf("%s/%s", te.Path, role.Prefix)
		}
	}

	// Attach the given display name if any
	if data.DisplayName != "" {
		full := "token-" + data.DisplayName
		full = displayNameSanitize.ReplaceAllString(full, "-")
		full = strings.TrimSuffix(full, "-")
		te.DisplayName = full
	}

	// Allow specifying the ID of the token if the client has root or sudo privileges
	if data.ID != "" {
		if !isSudo {
			return logical.ErrorResponse("root or sudo privileges required to specify token id"),
				logical.ErrInvalidRequest
		}
		te.ID = data.ID
	}

	switch {
	// If we have a role, we don't even consider parent policies; the role
	// allowed policies trumps all
	case role != nil:
		if !strListSubset(role.AllowedPolicies, data.Policies) {
			return logical.ErrorResponse("token policies must be subset of the role's allowed policies"), logical.ErrInvalidRequest
		}

	case len(data.Policies) == 0:
		data.Policies = parent.Policies

	// When a role is not in use, only permit policies to be a subset unless
	// the client has root or sudo privileges
	case !isSudo && !strListSubset(parent.Policies, data.Policies):
		return logical.ErrorResponse("child policies must be subset of parent"), logical.ErrInvalidRequest
	}

	// Use a map to filter out/prevent duplicates
	policyMap := map[string]bool{}
	for _, policy := range data.Policies {
		if policy == "" {
			// Don't allow a policy with no name, even though it is a valid
			// slice member
			continue
		}
		policyMap[policy] = true
	}
	if !policyMap["root"] &&
		!data.NoDefaultPolicy {
		policyMap["default"] = true
	}

	for k, _ := range policyMap {
		te.Policies = append(te.Policies, k)
	}
	sort.Strings(te.Policies)

	switch {
	case role != nil:
		if role.Orphan {
			te.Parent = ""
		}

	case data.NoParent:
		// Only allow an orphan token if the client has sudo policy
		if !isSudo {
			return logical.ErrorResponse("root or sudo privileges required to create orphan token"),
				logical.ErrInvalidRequest
		}

		te.Parent = ""

	default:
		// This comes from create-orphan, which can be properly ACLd
		if orphan {
			te.Parent = ""
		}
	}

	if role != nil && role.Period > 0 {
		te.TTL = role.Period
	} else {
		// Parse the TTL/lease if any
		if data.TTL != "" {
			dur, err := time.ParseDuration(data.TTL)
			if err != nil {
				return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
			}
			if dur < 0 {
				return logical.ErrorResponse("ttl must be positive"), logical.ErrInvalidRequest
			}
			te.TTL = dur
		} else if data.Lease != "" {
			dur, err := time.ParseDuration(data.Lease)
			if err != nil {
				return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
			}
			if dur < 0 {
				return logical.ErrorResponse("lease must be positive"), logical.ErrInvalidRequest
			}
			te.TTL = dur
		}

		sysView := ts.System()

		// Set the default lease if non-provided, root tokens are exempt
		if te.TTL == 0 && !strListContains(te.Policies, "root") {
			te.TTL = sysView.DefaultLeaseTTL()
		}

		// Limit the lease duration
		if te.TTL > sysView.MaxLeaseTTL() {
			te.TTL = sysView.MaxLeaseTTL()
		}
	}

	// Create the token
	if err := ts.create(&te); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Generate the response
	resp := &logical.Response{
		Auth: &logical.Auth{
			DisplayName: te.DisplayName,
			Policies:    te.Policies,
			Metadata:    te.Meta,
			LeaseOptions: logical.LeaseOptions{
				TTL:       te.TTL,
				Renewable: true,
			},
			ClientToken: te.ID,
			InternalData: map[string]interface{}{
				"id": te.ID,
			},
		},
	}

	if ts.policyLookupFunc != nil {
		for _, p := range te.Policies {
			policy, err := ts.policyLookupFunc(p)
			if err != nil {
				return logical.ErrorResponse(fmt.Sprintf("could not look up policy %s", p)), nil
			}
			if policy == nil {
				resp.AddWarning(fmt.Sprintf("policy \"%s\" does not exist", p))
			}
		}
	}

	return resp, nil
}

// handleRevokeSelf handles the auth/token/revoke-self path for revocation of tokens
// in a way that revokes all child tokens. Normally, using sys/revoke/leaseID will revoke
// the token and all children anyways, but that is only available when there is a lease.
func (ts *TokenStore) handleRevokeSelf(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Revoke the token and its children
	if err := ts.RevokeTree(req.ClientToken); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
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

	parent, err := ts.Lookup(req.ClientToken)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("parent token lookup failed: %s", err.Error())), logical.ErrInvalidRequest
	}
	if parent == nil {
		return logical.ErrorResponse("parent token lookup failed"), logical.ErrInvalidRequest
	}

	// Check if the client token has sudo/root privileges for the requested path
	isSudo := ts.System().SudoPrivilege(req.MountPoint+req.Path, req.ClientToken)

	if !isSudo {
		return logical.ErrorResponse("root or sudo privileges required to revoke and orphan"),
			logical.ErrInvalidRequest
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
	// you could escalate your privileges.
	resp := &logical.Response{
		Data: map[string]interface{}{
			"id":            out.ID,
			"policies":      out.Policies,
			"path":          out.Path,
			"meta":          out.Meta,
			"display_name":  out.DisplayName,
			"num_uses":      out.NumUses,
			"orphan":        false,
			"creation_time": int64(out.CreationTime),
			"creation_ttl":  int64(out.TTL.Seconds()),
			"ttl":           int64(0),
			"role":          out.Role,
		},
	}

	if out.Parent == "" {
		resp.Data["orphan"] = true
	}

	// Fetch the last renewal time
	leaseTimes, err := ts.expiration.FetchLeaseTimesByToken(out.Path, out.ID)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	if leaseTimes != nil {
		if !leaseTimes.LastRenewalTime.IsZero() {
			resp.Data["last_renewal_time"] = leaseTimes.LastRenewalTime.Unix()
		}
		if !leaseTimes.ExpireTime.IsZero() {
			resp.Data["ttl"] = int64(leaseTimes.ExpireTime.Sub(time.Now().Round(time.Second)).Seconds())
		}
	}

	return resp, nil
}

func (ts *TokenStore) handleRenewSelf(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	data.Raw["token"] = req.ClientToken
	return ts.handleRenew(req, data)
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
	te, err := ts.Lookup(id)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Verify the token exists
	if te == nil {
		return logical.ErrorResponse("token not found"), logical.ErrInvalidRequest
	}

	// Renew the token and its children
	auth, err := ts.expiration.RenewToken(req, te.Path, te.ID, increment)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Generate the response
	resp := &logical.Response{
		Auth: auth,
	}
	return resp, nil
}

func (ts *TokenStore) destroyCubbyhole(saltedID string) error {
	if ts.cubbyholeBackend == nil {
		// Should only ever happen in testing
		return nil
	}
	return ts.cubbyholeBackend.revoke(salt.SaltID(ts.cubbyholeBackend.saltUUID, saltedID, salt.SHA1Hash))
}

func (ts *TokenStore) authRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if req.Auth == nil {
		return nil, fmt.Errorf("request auth is nil")
	}

	f := framework.LeaseExtend(0, 0, ts.System())

	idInt, ok := req.Auth.InternalData["id"]
	if !ok {
		// Fall back here; this is pre-roles so there are no stored IDs, so use previous behavior
		return f(req, d)
	}

	id, ok := idInt.(string)
	if !ok {
		return nil, fmt.Errorf("found id in internal data but could not interpret as string")
	}

	te, err := ts.Lookup(id)
	if err != nil {
		return nil, fmt.Errorf("error looking up token: %s", err)
	}
	if te == nil {
		return nil, fmt.Errorf("no token entry found during lookup")
	}

	// No role? Use normal LeaseExtend semantics
	if te.Role == "" {
		return f(req, d)
	}

	role, err := ts.getTokenStoreRole(req.Storage, te.Role)
	if err != nil {
		return nil, fmt.Errorf("error looking up role %s: %s", te.Role, err)
	}

	if role == nil {
		resp, err := f(req, d)
		if resp != nil {
			resp.AddWarning(fmt.Sprintf("The token was created via role %s, but that role could no longer be found. Renewal fell back to normal token renewal semantics; if this token was a periodic token, this could mean that it can no longer be renewed.", te.Role))
		}
		return resp, err
	}

	if role.Period != 0 {
		req.Auth.TTL = role.Period
		return &logical.Response{Auth: req.Auth}, nil
	}

	return f(req, d)
}

func (ts *TokenStore) getTokenStoreRole(s logical.Storage, n string) (*tsRoleEntry, error) {
	entry, err := s.Get("roles/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result tsRoleEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (ts *TokenStore) tokenStoreRoleList(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List("roles/")
	if err != nil {
		return nil, err
	}

	ret := make([]string, len(entries))
	for i, entry := range entries {
		ret[i] = strings.TrimPrefix(entry, "roles/")
	}

	return logical.ListResponse(ret), nil
}

func (ts *TokenStore) tokenStoreRoleDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete("roles/" + data.Get("role_name").(string))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ts *TokenStore) tokenStoreRoleRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	role, err := ts.getTokenStoreRole(req.Storage, data.Get("role_name").(string))
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: structs.New(role).Map(),
	}

	return resp, nil
}

func (ts *TokenStore) tokenStoreRoleCreate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("role_name").(string)
	if name == "" {
		return logical.ErrorResponse("role name cannot be empty"), nil
	}

	prefix := data.Get("prefix").(string)
	if prefix != "" {
		matched := prefixSanitize.MatchString(prefix)
		if !matched {
			return logical.ErrorResponse(fmt.Sprintf("given role prefix contains invalid characters; must match %s", prefixSanitize.String())), nil
		}
	}

	entry := &tsRoleEntry{
		Name:   name,
		Orphan: data.Get("orphan").(bool),
		Period: time.Second * time.Duration(data.Get("period").(int)),
		Prefix: prefix,
	}

	allowedPolicies := data.Get("allowed_policies").(string)
	if allowedPolicies != "" {
		entry.AllowedPolicies = strings.Split(allowedPolicies, ",")
	}

	// Store it
	jsonEntry, err := logical.StorageEntryJSON("roles/"+name, entry)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(jsonEntry); err != nil {
		return nil, err
	}

	return nil, nil
}

const (
	tokenBackendHelp = `The token credential backend is always enabled and builtin to Vault.
Client tokens are used to identify a client and to allow Vault to associate policies and ACLs
which are enforced on every request. This backend also allows for generating sub-tokens as well
as revocation of tokens. The tokens are renewable if associated with a lease.`
	tokenCreateHelp       = `The token create path is used to create new tokens.`
	tokenCreateOrphanHelp = `The token create path is used to create new orphan tokens.`
	tokenLookupHelp       = `This endpoint will lookup a token and its properties.`
	tokenRevokeHelp       = `This endpoint will delete the given token and all of its child tokens.`
	tokenRevokeSelfHelp   = `This endpoint will delete the token used to call it and all of its child tokens.`
	tokenRevokeOrphanHelp = `This endpoint will delete the token and orphan its child tokens.`
	tokenRevokePrefixHelp = `This endpoint will delete all tokens generated under a prefix with their child tokens.`
	tokenRenewHelp        = `This endpoint will renew the given token and prevent expiration.`
	tokenRenewSelfHelp    = `This endpoint will renew the token used to call it and prevent expiration.`
)
