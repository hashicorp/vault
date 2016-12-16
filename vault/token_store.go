package vault

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/duration"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/locksutil"
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/mitchellh/mapstructure"
)

const (
	// lookupPrefix is the prefix used to store tokens for their
	// primary ID based index
	lookupPrefix = "id/"

	// accessorPrefix is the prefix used to store the index from
	// Accessor to Token ID
	accessorPrefix = "accessor/"

	// parentPrefix is the prefix used to store tokens for their
	// secondar parent based index
	parentPrefix = "parent/"

	// tokenSubPath is the sub-path used for the token store
	// view. This is nested under the system view.
	tokenSubPath = "token/"

	// rolesPrefix is the prefix used to store role information
	rolesPrefix = "roles/"

	// tokenRevocationDeferred indicates that the token should not be used
	// again but is currently fulfilling its final use
	tokenRevocationDeferred = -1

	// tokenRevocationInProgress indicates that revocation of that token/its
	// leases is ongoing
	tokenRevocationInProgress = -2

	// tokenRevocationFailed indicates that revocation failed; the entry is
	// kept around so that when the tidy function is run it can be tried
	// again (or when the revocation function is run again), but all other uses
	// will report the token invalid
	tokenRevocationFailed = -3
)

var (
	// displayNameSanitize is used to sanitize a display name given to a token.
	displayNameSanitize = regexp.MustCompile("[^a-zA-Z0-9-]")

	// pathSuffixSanitize is used to ensure a path suffix in a role is valid.
	pathSuffixSanitize = regexp.MustCompile("\\w[\\w-.]+\\w")

	destroyCubbyhole = func(ts *TokenStore, saltedID string) error {
		if ts.cubbyholeBackend == nil {
			// Should only ever happen in testing
			return nil
		}
		return ts.cubbyholeBackend.revoke(salt.SaltID(ts.cubbyholeBackend.saltUUID, saltedID, salt.SHA1Hash))
	}
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

	tokenLocks map[string]*sync.RWMutex

	cubbyholeDestroyer func(*TokenStore, string) error
}

// NewTokenStore is used to construct a token store that is
// backed by the given barrier view.
func NewTokenStore(c *Core, config *logical.BackendConfig) (*TokenStore, error) {
	// Create a sub-view
	view := c.systemBarrierView.SubView(tokenSubPath)

	// Initialize the store
	t := &TokenStore{
		view:               view,
		cubbyholeDestroyer: destroyCubbyhole,
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

	t.tokenLocks = map[string]*sync.RWMutex{}

	// Create 256 locks
	if err = locksutil.CreateLocks(t.tokenLocks, 256); err != nil {
		return nil, fmt.Errorf("failed to create locks: %v", err)
	}

	t.tokenLocks["custom"] = &sync.RWMutex{}

	// Setup the framework endpoints
	t.Backend = &framework.Backend{
		AuthRenew: t.authRenew,

		PathsSpecial: &logical.Paths{
			Root: []string{
				"revoke-orphan/*",
				"accessors*",
			},
		},

		Paths: []*framework.Path{
			&framework.Path{
				Pattern: "roles/?$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ListOperation: t.tokenStoreRoleList,
				},

				HelpSynopsis:    tokenListRolesHelp,
				HelpDescription: tokenListRolesHelp,
			},

			&framework.Path{
				Pattern: "accessors/?$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ListOperation: t.tokenStoreAccessorList,
				},

				HelpSynopsis:    tokenListAccessorsHelp,
				HelpDescription: tokenListAccessorsHelp,
			},

			&framework.Path{
				Pattern: "roles/" + framework.GenericNameRegex("role_name"),
				Fields: map[string]*framework.FieldSchema{
					"role_name": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Name of the role",
					},

					"allowed_policies": &framework.FieldSchema{
						Type:        framework.TypeString,
						Default:     "",
						Description: tokenAllowedPoliciesHelp,
					},

					"disallowed_policies": &framework.FieldSchema{
						Type:        framework.TypeString,
						Default:     "",
						Description: tokenDisallowedPoliciesHelp,
					},

					"orphan": &framework.FieldSchema{
						Type:        framework.TypeBool,
						Default:     false,
						Description: tokenOrphanHelp,
					},

					"period": &framework.FieldSchema{
						Type:        framework.TypeDurationSecond,
						Default:     0,
						Description: tokenPeriodHelp,
					},

					"path_suffix": &framework.FieldSchema{
						Type:        framework.TypeString,
						Default:     "",
						Description: tokenPathSuffixHelp + pathSuffixSanitize.String(),
					},

					"explicit_max_ttl": &framework.FieldSchema{
						Type:        framework.TypeDurationSecond,
						Default:     0,
						Description: tokenExplicitMaxTTLHelp,
					},

					"renewable": &framework.FieldSchema{
						Type:        framework.TypeBool,
						Default:     true,
						Description: tokenRenewableHelp,
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation:   t.tokenStoreRoleRead,
					logical.CreateOperation: t.tokenStoreRoleCreateUpdate,
					logical.UpdateOperation: t.tokenStoreRoleCreateUpdate,
					logical.DeleteOperation: t.tokenStoreRoleDelete,
				},

				ExistenceCheck: t.tokenStoreRoleExistenceCheck,

				HelpSynopsis:    tokenPathRolesHelp,
				HelpDescription: tokenPathRolesHelp,
			},

			&framework.Path{
				Pattern: "create-orphan$",

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
					logical.UpdateOperation: t.handleCreateAgainstRole,
				},

				HelpSynopsis:    strings.TrimSpace(tokenCreateRoleHelp),
				HelpDescription: strings.TrimSpace(tokenCreateRoleHelp),
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
				Pattern: "lookup" + framework.OptionalParamRegex("urltoken"),

				Fields: map[string]*framework.FieldSchema{
					"urltoken": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Token to lookup (URL parameter)",
					},
					"token": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Token to lookup (POST request body)",
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation:   t.handleLookup,
					logical.UpdateOperation: t.handleLookup,
				},

				HelpSynopsis:    strings.TrimSpace(tokenLookupHelp),
				HelpDescription: strings.TrimSpace(tokenLookupHelp),
			},

			&framework.Path{
				Pattern: "lookup-accessor" + framework.OptionalParamRegex("urlaccessor"),

				Fields: map[string]*framework.FieldSchema{
					"urlaccessor": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Accessor of the token to look up (URL parameter)",
					},
					"accessor": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Accessor of the token to look up (request body)",
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: t.handleUpdateLookupAccessor,
				},

				HelpSynopsis:    strings.TrimSpace(tokenLookupAccessorHelp),
				HelpDescription: strings.TrimSpace(tokenLookupAccessorHelp),
			},

			&framework.Path{
				Pattern: "lookup-self$",

				Fields: map[string]*framework.FieldSchema{
					"token": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Token to look up (unused, does not need to be set)",
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: t.handleLookupSelf,
					logical.ReadOperation:   t.handleLookupSelf,
				},

				HelpSynopsis:    strings.TrimSpace(tokenLookupHelp),
				HelpDescription: strings.TrimSpace(tokenLookupHelp),
			},

			&framework.Path{
				Pattern: "revoke-accessor" + framework.OptionalParamRegex("urlaccessor"),

				Fields: map[string]*framework.FieldSchema{
					"urlaccessor": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Accessor of the token (URL parameter)",
					},
					"accessor": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Accessor of the token (request body)",
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: t.handleUpdateRevokeAccessor,
				},

				HelpSynopsis:    strings.TrimSpace(tokenRevokeAccessorHelp),
				HelpDescription: strings.TrimSpace(tokenRevokeAccessorHelp),
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
				Pattern: "revoke" + framework.OptionalParamRegex("urltoken"),

				Fields: map[string]*framework.FieldSchema{
					"urltoken": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Token to revoke (URL parameter)",
					},
					"token": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Token to revoke (request body)",
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: t.handleRevokeTree,
				},

				HelpSynopsis:    strings.TrimSpace(tokenRevokeHelp),
				HelpDescription: strings.TrimSpace(tokenRevokeHelp),
			},

			&framework.Path{
				Pattern: "revoke-orphan" + framework.OptionalParamRegex("urltoken"),

				Fields: map[string]*framework.FieldSchema{
					"urltoken": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Token to revoke (URL parameter)",
					},
					"token": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Token to revoke (request body)",
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: t.handleRevokeOrphan,
				},

				HelpSynopsis:    strings.TrimSpace(tokenRevokeOrphanHelp),
				HelpDescription: strings.TrimSpace(tokenRevokeOrphanHelp),
			},

			&framework.Path{
				Pattern: "renew-self$",

				Fields: map[string]*framework.FieldSchema{
					"token": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Token to renew (unused, does not need to be set)",
					},
					"increment": &framework.FieldSchema{
						Type:        framework.TypeDurationSecond,
						Default:     0,
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
				Pattern: "renew" + framework.OptionalParamRegex("urltoken"),

				Fields: map[string]*framework.FieldSchema{
					"urltoken": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Token to renew (URL parameter)",
					},
					"token": &framework.FieldSchema{
						Type:        framework.TypeString,
						Description: "Token to renew (request body)",
					},
					"increment": &framework.FieldSchema{
						Type:        framework.TypeDurationSecond,
						Default:     0,
						Description: "The desired increment in seconds to the token expiration",
					},
				},

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: t.handleRenew,
				},

				HelpSynopsis:    strings.TrimSpace(tokenRenewHelp),
				HelpDescription: strings.TrimSpace(tokenRenewHelp),
			},

			&framework.Path{
				Pattern: "tidy$",

				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: t.handleTidy,
				},

				HelpSynopsis:    strings.TrimSpace(tokenTidyHelp),
				HelpDescription: strings.TrimSpace(tokenTidyDesc),
			},
		},
	}

	t.Backend.Setup(config)

	return t, nil
}

// TokenEntry is used to represent a given token
type TokenEntry struct {
	// ID of this entry, generally a random UUID
	ID string `json:"id" mapstructure:"id" structs:"id"`

	// Accessor for this token, a random UUID
	Accessor string `json:"accessor" mapstructure:"accessor" structs:"accessor"`

	// Parent token, used for revocation trees
	Parent string `json:"parent" mapstructure:"parent" structs:"parent"`

	// Which named policies should be used
	Policies []string `json:"policies" mapstructure:"policies" structs:"policies"`

	// Used for audit trails, this is something like "auth/user/login"
	Path string `json:"path" mapstructure:"path" structs:"path"`

	// Used for auditing. This could include things like "source", "user", "ip"
	Meta map[string]string `json:"meta" mapstructure:"meta" structs:"meta"`

	// Used for operators to be able to associate with the source
	DisplayName string `json:"display_name" mapstructure:"display_name" structs:"display_name"`

	// Used to restrict the number of uses (zero is unlimited). This is to
	// support one-time-tokens (generalized). There are a few special values:
	// if it's -1 it has run through its use counts and is executing its final
	// use; if it's -2 it is tainted, which means revocation is currently
	// running on it; and if it's -3 it's also tainted but revocation
	// previously ran and failed, so this hints the tidy function to try it
	// again.
	NumUses int `json:"num_uses" mapstructure:"num_uses" structs:"num_uses"`

	// Time of token creation
	CreationTime int64 `json:"creation_time" mapstructure:"creation_time" structs:"creation_time"`

	// Duration set when token was created
	TTL time.Duration `json:"ttl" mapstructure:"ttl" structs:"ttl"`

	// Explicit maximum TTL on the token
	ExplicitMaxTTL time.Duration `json:"explicit_max_ttl" mapstructure:"explicit_max_ttl" structs:"explicit_max_ttl"`

	// If set, the role that was used for parameters at creation time
	Role string `json:"role" mapstructure:"role" structs:"role"`

	// If set, the period of the token. This is only used when created directly
	// through the create endpoint; periods managed by roles or other auth
	// backends are subject to those renewal rules.
	Period time.Duration `json:"period" mapstructure:"period" structs:"period"`

	// These are the deprecated fields
	DisplayNameDeprecated    string        `json:"DisplayName" mapstructure:"DisplayName" structs:"DisplayName"`
	NumUsesDeprecated        int           `json:"NumUses" mapstructure:"NumUses" structs:"NumUses"`
	CreationTimeDeprecated   int64         `json:"CreationTime" mapstructure:"CreationTime" structs:"CreationTime"`
	ExplicitMaxTTLDeprecated time.Duration `json:"ExplicitMaxTTL" mapstructure:"ExplicitMaxTTL" structs:"ExplicitMaxTTL"`
}

// tsRoleEntry contains token store role information
type tsRoleEntry struct {
	// The name of the role. Embedded so it can be used for pathing
	Name string `json:"name" mapstructure:"name" structs:"name"`

	// The policies that creation functions using this role can assign to a token,
	// escaping or further locking down normal subset checking
	AllowedPolicies []string `json:"allowed_policies" mapstructure:"allowed_policies" structs:"allowed_policies"`

	// List of policies to be not allowed during token creation using this role
	DisallowedPolicies []string `json:"disallowed_policies" mapstructure:"disallowed_policies" structs:"disallowed_policies"`

	// If true, tokens created using this role will be orphans
	Orphan bool `json:"orphan" mapstructure:"orphan" structs:"orphan"`

	// If non-zero, tokens created using this role will be able to be renewed
	// forever, but will have a fixed renewal period of this value
	Period time.Duration `json:"period" mapstructure:"period" structs:"period"`

	// If set, a suffix will be set on the token path, making it easier to
	// revoke using 'revoke-prefix'
	PathSuffix string `json:"path_suffix" mapstructure:"path_suffix" structs:"path_suffix"`

	// If set, controls whether created tokens are marked as being renewable
	Renewable bool `json:"renewable" mapstructure:"renewable" structs:"renewable"`

	// If set, the token entry will have an explicit maximum TTL set, rather
	// than deferring to role/mount values
	ExplicitMaxTTL time.Duration `json:"explicit_max_ttl" mapstructure:"explicit_max_ttl" structs:"explicit_max_ttl"`
}

type accessorEntry struct {
	TokenID    string `json:"token_id"`
	AccessorID string `json:"accessor_id"`
}

// SetExpirationManager is used to provide the token store with
// an expiration manager. This is used to manage prefix based revocation
// of tokens and to tidy entries when removed from the token store.
func (ts *TokenStore) SetExpirationManager(exp *ExpirationManager) {
	ts.expiration = exp
}

// SaltID is used to apply a salt and hash to an ID to make sure its not reversible
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

func (ts *TokenStore) tokenStoreAccessorList(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := ts.view.List(accessorPrefix)
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{}

	ret := make([]string, 0, len(entries))
	for _, entry := range entries {
		aEntry, err := ts.lookupBySaltedAccessor(entry)
		if err != nil {
			resp.AddWarning("Found an accessor entry that could not be successfully decoded")
			continue
		}
		if aEntry.TokenID == "" {
			resp.AddWarning(fmt.Sprintf("Found an accessor entry missing a token: %v", aEntry.AccessorID))
		} else {
			ret = append(ret, aEntry.AccessorID)
		}
	}

	resp.Data = map[string]interface{}{
		"keys": ret,
	}
	return resp, nil
}

// createAccessor is used to create an identifier for the token ID.
// A storage index, mapping the accessor to the token ID is also created.
func (ts *TokenStore) createAccessor(entry *TokenEntry) error {
	defer metrics.MeasureSince([]string{"token", "createAccessor"}, time.Now())

	// Create a random accessor
	accessorUUID, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	entry.Accessor = accessorUUID

	// Create index entry, mapping the accessor to the token ID
	path := accessorPrefix + ts.SaltID(entry.Accessor)

	aEntry := &accessorEntry{
		TokenID:    entry.ID,
		AccessorID: entry.Accessor,
	}
	aEntryBytes, err := jsonutil.EncodeJSON(aEntry)
	if err != nil {
		return fmt.Errorf("failed to marshal accessor index entry: %v", err)
	}

	le := &logical.StorageEntry{Key: path, Value: aEntryBytes}
	if err := ts.view.Put(le); err != nil {
		return fmt.Errorf("failed to persist accessor index entry: %v", err)
	}
	return nil
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

	entry.Policies = policyutil.SanitizePolicies(entry.Policies, policyutil.DoNotAddDefaultPolicy)

	err := ts.createAccessor(entry)
	if err != nil {
		return err
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

func (ts *TokenStore) getTokenLock(id string) *sync.RWMutex {
	// Find our multilevel lock, or fall back to global
	var lock *sync.RWMutex
	var ok bool
	if len(id) >= 2 {
		lock, ok = ts.tokenLocks[id[0:2]]
	}
	if !ok || lock == nil {
		// Fall back for custom token IDs
		lock = ts.tokenLocks["custom"]
	}

	return lock
}

// UseToken is used to manage restricted use tokens and decrement their
// available uses. Returns two values: a potentially updated entry or, if the
// token has been revoked, nil; and whether an error was encountered. The
// locking here isn't perfect, as other parts of the code may update an entry,
// but usually none after the entry is already created...so this is pretty
// good.
func (ts *TokenStore) UseToken(te *TokenEntry) (*TokenEntry, error) {
	if te == nil {
		return nil, fmt.Errorf("invalid token entry provided for use count decrementing")
	}

	// This case won't be hit with a token with restricted uses because we go
	// from 1 to -1. So it's a nice optimization to check this without a read
	// lock.
	if te.NumUses == 0 {
		return te, nil
	}

	lock := ts.getTokenLock(te.ID)

	lock.Lock()
	defer lock.Unlock()

	// Call lookupSalted instead of Lookup to avoid deadlocking since Lookup grabs a read lock
	te, err := ts.lookupSalted(ts.SaltID(te.ID), false)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh entry: %v", err)
	}
	// If it can't be found we shouldn't be trying to use it, so if we get nil
	// back, it is because it has been revoked in the interim or will be
	// revoked (NumUses is -1)
	if te == nil {
		return nil, fmt.Errorf("token not found or fully used already")
	}

	// Decrement the count. If this is our last use count, we need to indicate
	// that this is no longer valid, but revocation is deferred to the end of
	// the call, so this will make sure that any Lookup that happens doesn't
	// return an entry. This essentially acts as a write-ahead lock and is
	// especially useful since revocation can end up (via the expiration
	// manager revoking children) attempting to acquire the same lock
	// repeatedly.
	if te.NumUses == 1 {
		te.NumUses = -1
	} else {
		te.NumUses -= 1
	}

	err = ts.storeCommon(te, false)
	if err != nil {
		return nil, err
	}

	return te, nil
}

func (ts *TokenStore) UseTokenByID(id string) (*TokenEntry, error) {
	te, err := ts.Lookup(id)
	if err != nil {
		return te, err
	}

	return ts.UseToken(te)
}

// Lookup is used to find a token given its ID. It acquires a read lock, then calls lookupSalted.
func (ts *TokenStore) Lookup(id string) (*TokenEntry, error) {
	defer metrics.MeasureSince([]string{"token", "lookup"}, time.Now())
	if id == "" {
		return nil, fmt.Errorf("cannot lookup blank token")
	}

	lock := ts.getTokenLock(id)
	lock.RLock()
	defer lock.RUnlock()

	return ts.lookupSalted(ts.SaltID(id), false)
}

// lookupSalted is used to find a token given its salted ID. If tainted is
// true, entries that are in some revocation state (currently, indicated by num
// uses < 0), the entry will be returned anyways
func (ts *TokenStore) lookupSalted(saltedId string, tainted bool) (*TokenEntry, error) {
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
	if err := jsonutil.DecodeJSON(raw.Value, entry); err != nil {
		return nil, fmt.Errorf("failed to decode entry: %v", err)
	}

	// This is a token that is awaiting deferred revocation or tainted
	if entry.NumUses < 0 && !tainted {
		return nil, nil
	}

	persistNeeded := false

	// Upgrade the deprecated fields
	if entry.DisplayNameDeprecated != "" {
		if entry.DisplayName == "" {
			entry.DisplayName = entry.DisplayNameDeprecated
		}
		entry.DisplayNameDeprecated = ""
		persistNeeded = true
	}

	if entry.CreationTimeDeprecated != 0 {
		if entry.CreationTime == 0 {
			entry.CreationTime = entry.CreationTimeDeprecated
		}
		entry.CreationTimeDeprecated = 0
		persistNeeded = true
	}

	if entry.ExplicitMaxTTLDeprecated != 0 {
		if entry.ExplicitMaxTTL == 0 {
			entry.ExplicitMaxTTL = entry.ExplicitMaxTTLDeprecated
		}
		entry.ExplicitMaxTTLDeprecated = 0
		persistNeeded = true
	}

	if entry.NumUsesDeprecated != 0 {
		if entry.NumUses == 0 || entry.NumUsesDeprecated < entry.NumUses {
			entry.NumUses = entry.NumUsesDeprecated
		}
		entry.NumUsesDeprecated = 0
		persistNeeded = true
	}

	// If fields are getting upgraded, store the changes
	if persistNeeded {
		if err := ts.storeCommon(entry, false); err != nil {
			return nil, fmt.Errorf("failed to persist token upgrade: %v", err)
		}
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
func (ts *TokenStore) revokeSalted(saltedId string) (ret error) {
	// Protect the entry lookup/writing with locks. The rub here is that we
	// don't know the ID until we look it up once, so first we look it up, then
	// do a locked lookup.
	entry, err := ts.lookupSalted(saltedId, true)
	if err != nil {
		return err
	}
	if entry == nil {
		return nil
	}

	lock := ts.getTokenLock(entry.ID)
	lock.Lock()

	// Lookup the token first
	entry, err = ts.lookupSalted(saltedId, true)
	if err != nil {
		lock.Unlock()
		return err
	}

	if entry == nil {
		lock.Unlock()
		return nil
	}

	// On failure we write -3, so if we hit -2 here we're already running a
	// revocation operation. This can happen due to e.g. recursion into this
	// function via the expiration manager's RevokeByToken.
	if entry.NumUses == tokenRevocationInProgress {
		lock.Unlock()
		return nil
	}

	// This acts as a WAL. lookupSalted will no longer return this entry,
	// so the token cannot be used, but this way we can keep the entry
	// around until after the rest of this function is attempted, and a
	// tidy function can key off of this value to try again.
	entry.NumUses = tokenRevocationInProgress
	err = ts.storeCommon(entry, false)
	lock.Unlock()
	if err != nil {
		return err
	}

	// If we are returning an error, mark the entry with -3 to indicate
	// failed revocation. This way we don't try to clean up during active
	// revocation (-2).
	defer func() {
		if ret != nil {
			lock.Lock()
			defer lock.Unlock()

			// Lookup the token again to make sure something else didn't
			// revoke in the interim
			entry, err := ts.lookupSalted(saltedId, true)
			if err != nil {
				return
			}

			// If it exists just taint to -3 rather than trying to figure
			// out what it means if it's already -3 after the -2 above
			if entry != nil {
				entry.NumUses = tokenRevocationFailed
				ts.storeCommon(entry, false)
			}
		}
	}()

	// Destroy the token's cubby. This should go first as it's a
	// security-sensitive item.
	err = ts.cubbyholeDestroyer(ts, saltedId)
	if err != nil {
		return err
	}

	// Revoke all secrets under this token. This should go first as it's a
	// security-sensitive item.
	if err := ts.expiration.RevokeByToken(entry); err != nil {
		return err
	}

	// Clear the secondary index if any
	if entry.Parent != "" {
		path := parentPrefix + ts.SaltID(entry.Parent) + "/" + saltedId
		if err = ts.view.Delete(path); err != nil {
			return fmt.Errorf("failed to delete entry: %v", err)
		}
	}

	// Clear the accessor index if any
	if entry.Accessor != "" {
		path := accessorPrefix + ts.SaltID(entry.Accessor)
		if err = ts.view.Delete(path); err != nil {
			return fmt.Errorf("failed to delete entry: %v", err)
		}
	}

	// Now that the entry is not usable for any revocation tasks, nuke it
	path := lookupPrefix + saltedId
	if err = ts.view.Delete(path); err != nil {
		return fmt.Errorf("failed to delete entry: %v", err)
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

// handleCreateAgainstRole handles the auth/token/create path for a role
func (ts *TokenStore) handleCreateAgainstRole(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("role_name").(string)
	roleEntry, err := ts.tokenStoreRole(name)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown role %s", name)), nil
	}

	return ts.handleCreateCommon(req, d, false, roleEntry)
}

func (ts *TokenStore) lookupByAccessor(accessor string) (accessorEntry, error) {
	return ts.lookupBySaltedAccessor(ts.SaltID(accessor))
}

func (ts *TokenStore) lookupBySaltedAccessor(saltedAccessor string) (accessorEntry, error) {
	entry, err := ts.view.Get(accessorPrefix + saltedAccessor)
	var aEntry accessorEntry

	if err != nil {
		return aEntry, fmt.Errorf("failed to read index using accessor: %s", err)
	}
	if entry == nil {
		return aEntry, &StatusBadRequest{Err: "invalid accessor"}
	}

	err = jsonutil.DecodeJSON(entry.Value, &aEntry)
	// If we hit an error, assume it's a pre-struct straight token ID
	if err != nil {
		aEntry.TokenID = string(entry.Value)
		te, err := ts.lookupSalted(ts.SaltID(aEntry.TokenID), false)
		if err != nil {
			return accessorEntry{}, fmt.Errorf("failed to look up token using accessor index: %s", err)
		}
		// It's hard to reason about what to do here -- it may be that the
		// token was revoked async, or that it's an old accessor index entry
		// that was somehow not cleared up, or or or. A nonexistent token entry
		// on lookup is nil, not an error, so we keep that behavior here to be
		// safe...the token ID is simply not filled in.
		if te != nil {
			aEntry.AccessorID = te.Accessor
		}
	}

	return aEntry, nil
}

// handleTidy handles the cleaning up of leaked accessor storage entries and
// cleaning up of leases that are associated to tokens that are expired.
func (ts *TokenStore) handleTidy(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// List out all the accessors
	saltedAccessorList, err := ts.view.List(accessorPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch accessor entries: %v", err)
	}

	var tidyErrors *multierror.Error

	// First, clean up secondary index entries that are no longer valid
	parentList, err := ts.view.List(parentPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch secondary index entries: %v", err)
	}

	// Scan through the secondary index entries; if there is an entry
	// with the token's salt ID at the end, remove it
	for _, parent := range parentList {
		children, err := ts.view.List(parentPrefix + parent)
		if err != nil {
			tidyErrors = multierror.Append(tidyErrors, fmt.Errorf("failed to read child index entry: %v", err))
			continue
		}

		for _, child := range children {
			// Look up tainted entries so we can be sure that if this isn't
			// found, it doesn't exist
			te, _ := ts.lookupSalted(child, true)
			if te == nil {
				err = ts.view.Delete(parentPrefix + parent + "/" + child)
				if err != nil {
					tidyErrors = multierror.Append(tidyErrors, fmt.Errorf("failed to delete secondary index entry: %v", err))
				}
			}
		}
	}

	// For each of the accessor, see if the token ID associated with it is
	// a valid one. If not, delete the leases associated with that token
	// and delete the accessor as well.
	for _, saltedAccessor := range saltedAccessorList {
		accessorEntry, err := ts.lookupBySaltedAccessor(saltedAccessor)
		if err != nil {
			tidyErrors = multierror.Append(tidyErrors, fmt.Errorf("failed to read the accessor entry: %v", err))
			continue
		}

		// A valid accessor storage entry should always have a token ID
		// in it. If not, it is an invalid accessor entry and needs to
		// be deleted.
		if accessorEntry.TokenID == "" {
			// If deletion of accessor fails, move on to the next
			// item since this is just a best-effort operation
			err = ts.view.Delete(accessorPrefix + saltedAccessor)
			if err != nil {
				tidyErrors = multierror.Append(tidyErrors, fmt.Errorf("failed to delete the accessor entry: %v", err))
				continue
			}
		}

		saltedId := ts.SaltID(accessorEntry.TokenID)

		// Look up tainted variants so we only find entries that truly don't
		// exist
		te, err := ts.lookupSalted(saltedId, true)

		// If token entry is not found assume that the token is not valid any
		// more and conclude that accessor, leases, and secondary index entries
		// for this token should not exist as well.
		if te == nil {
			// RevokeByToken expects a '*TokenEntry'. For the
			// purposes of tidying, it is sufficient if the token
			// entry only has ID set.
			tokenEntry := &TokenEntry{
				ID: accessorEntry.TokenID,
			}

			// Attempt to revoke the token. This will also revoke
			// the leases associated with the token.
			err := ts.expiration.RevokeByToken(tokenEntry)
			if err != nil {
				tidyErrors = multierror.Append(tidyErrors, fmt.Errorf("failed to revoke leases of expired token: %v", err))
				continue
			}

			// If deletion of accessor fails, move on to the next item since
			// this is just a best-effort operation. We do this last so that on
			// next run if something above failed we still have the accessor
			// entry to try again.
			err = ts.view.Delete(accessorPrefix + saltedAccessor)
			if err != nil {
				tidyErrors = multierror.Append(tidyErrors, fmt.Errorf("failed to delete accessor entry: %v", err))
				continue
			}
		}
	}

	// Later request handling code seems to check if the type is multierror so
	// if we haven't added any errors we need to just return a normal nil error
	if tidyErrors == nil {
		return nil, nil
	}

	return nil, tidyErrors
}

// handleUpdateLookupAccessor handles the auth/token/lookup-accessor path for returning
// the properties of the token associated with the accessor
func (ts *TokenStore) handleUpdateLookupAccessor(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var urlaccessor bool
	accessor := data.Get("accessor").(string)
	if accessor == "" {
		accessor = data.Get("urlaccessor").(string)
		if accessor == "" {
			return nil, &StatusBadRequest{Err: "missing accessor"}
		}
		urlaccessor = true
	}

	aEntry, err := ts.lookupByAccessor(accessor)
	if err != nil {
		return nil, err
	}

	// Prepare the field data required for a lookup call
	d := &framework.FieldData{
		Raw: map[string]interface{}{
			"token": aEntry.TokenID,
		},
		Schema: map[string]*framework.FieldSchema{
			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Token to lookup",
			},
		},
	}
	resp, err := ts.handleLookup(req, d)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("failed to lookup the token")
	}
	if resp.IsError() {
		return resp, nil

	}

	// Remove the token ID from the response
	if resp.Data != nil {
		resp.Data["id"] = ""
	}

	if urlaccessor {
		resp.AddWarning(`Using an accessor in the path is unsafe as the accessor can be logged in many places. Please use POST or PUT with the accessor passed in via the "accessor" parameter.`)
	}

	return resp, nil
}

// handleUpdateRevokeAccessor handles the auth/token/revoke-accessor path for revoking
// the token associated with the accessor
func (ts *TokenStore) handleUpdateRevokeAccessor(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var urlaccessor bool
	accessor := data.Get("accessor").(string)
	if accessor == "" {
		accessor = data.Get("urlaccessor").(string)
		if accessor == "" {
			return nil, &StatusBadRequest{Err: "missing accessor"}
		}
		urlaccessor = true
	}

	aEntry, err := ts.lookupByAccessor(accessor)
	if err != nil {
		return nil, err
	}

	// Revoke the token and its children
	if err := ts.RevokeTree(aEntry.TokenID); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	if urlaccessor {
		resp := &logical.Response{}
		resp.AddWarning(`Using an accessor in the path is unsafe as the accessor can be logged in many places. Please use POST or PUT with the accessor passed in via the "accessor" parameter.`)
		return resp, nil
	}

	return nil, nil
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
		Renewable       *bool
		ExplicitMaxTTL  string `mapstructure:"explicit_max_ttl"`
		DisplayName     string `mapstructure:"display_name"`
		NumUses         int    `mapstructure:"num_uses"`
		Period          string
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
		Parent: req.ClientToken,

		// The mount point is always the same since we have only one token
		// store; using req.MountPoint causes trouble in tests since they don't
		// have an official mount
		Path: fmt.Sprintf("auth/token/%s", req.Path),

		Meta:         data.Metadata,
		DisplayName:  "token",
		NumUses:      data.NumUses,
		CreationTime: time.Now().Unix(),
	}

	renewable := true
	if data.Renewable != nil {
		renewable = *data.Renewable
	}

	// If the role is not nil, we add the role name as part of the token's
	// path. This makes it much easier to later revoke tokens that were issued
	// by a role (using revoke-prefix). Users can further specify a PathSuffix
	// in the role; that way they can use something like "v1", "v2" to indicate
	// role revisions, and revoke only tokens issued with a previous revision.
	if role != nil {
		te.Role = role.Name

		// If renewable hasn't been disabled in the call and the role has
		// renewability disabled, set renewable false
		if renewable && !role.Renewable {
			renewable = false
		}

		if role.PathSuffix != "" {
			te.Path = fmt.Sprintf("%s/%s", te.Path, role.PathSuffix)
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

	resp := &logical.Response{}

	var addDefault bool

	// N.B.: The logic here uses various calculations as to whether default
	// should be added. In the end we decided that if NoDefaultPolicy is set it
	// should be stripped out regardless, *but*, the logic of when it should
	// and shouldn't be added is kept because we want to do subset comparisons
	// based on adding default when it's correct to do so.
	switch {
	case role != nil && (len(role.AllowedPolicies) > 0 || len(role.DisallowedPolicies) > 0):
		// Holds the final set of policies as they get munged
		var finalPolicies []string

		// We don't make use of the global one because roles with allowed or
		// disallowed set do their own policy rules
		var localAddDefault bool

		// If the request doesn't say not to add "default" and if "default"
		// isn't in the disallowed list, add it. This is in line with the idea
		// that roles, when allowed/disallowed ar set, allow a subset of
		// policies to be set disjoint from the parent token's policies.
		if !data.NoDefaultPolicy && !strutil.StrListContains(role.DisallowedPolicies, "default") {
			localAddDefault = true
		}

		// Start with passed-in policies as a baseline, if they exist
		if len(data.Policies) > 0 {
			finalPolicies = policyutil.SanitizePolicies(data.Policies, localAddDefault)
		}

		var sanitizedRolePolicies []string

		// First check allowed policies; if policies are specified they will be
		// checked, otherwise if an allowed set exists that will be the set
		// that is used
		if len(role.AllowedPolicies) > 0 {
			// Note that if "default" is already in allowed, and also in
			// disallowed, this will still result in an error later since this
			// doesn't strip out default
			sanitizedRolePolicies = policyutil.SanitizePolicies(role.AllowedPolicies, localAddDefault)

			if len(finalPolicies) == 0 {
				finalPolicies = sanitizedRolePolicies
			} else {
				if !strutil.StrListSubset(sanitizedRolePolicies, finalPolicies) {
					return logical.ErrorResponse(fmt.Sprintf("token policies (%v) must be subset of the role's allowed policies (%v)", finalPolicies, sanitizedRolePolicies)), logical.ErrInvalidRequest
				}
			}
		} else {
			// Check against parent policies, or assign parent policies. As
			// this is a role, add default unless explicitly disabled.
			if len(finalPolicies) == 0 {
				finalPolicies = policyutil.SanitizePolicies(parent.Policies, localAddDefault)
			} else {
				// If we added default based on the fact that this is using a
				// role, we need to add it here too to ensure that the subset
				// matching works.
				sanitizedParentPolicies := policyutil.SanitizePolicies(parent.Policies, localAddDefault)
				if !strutil.StrListSubset(sanitizedParentPolicies, finalPolicies) {
					return logical.ErrorResponse("child policies must be subset of parent when role contains no allowed_policies"), logical.ErrInvalidRequest
				}
			}
		}

		if len(role.DisallowedPolicies) > 0 {
			// We don't add the default here because we only want to disallow it if it's explicitly set
			sanitizedRolePolicies = policyutil.SanitizePolicies(role.DisallowedPolicies, policyutil.DoNotAddDefaultPolicy)

			for _, finalPolicy := range finalPolicies {
				if strutil.StrListContains(sanitizedRolePolicies, finalPolicy) {
					return logical.ErrorResponse(fmt.Sprintf("token policy %q is disallowed by this role", finalPolicy)), logical.ErrInvalidRequest
				}
			}
		}

		data.Policies = finalPolicies

	// No policies specified, inherit parent
	case len(data.Policies) == 0:
		// Only inherit "default" if the parent already has it, so don't touch addDefault here
		data.Policies = policyutil.SanitizePolicies(parent.Policies, policyutil.DoNotAddDefaultPolicy)

	// When a role is not in use or does not specify allowed/disallowed, only
	// permit policies to be a subset unless the client has root or sudo
	// privileges. Default is added in this case if the parent has it, unless
	// the client specified for it not to be added.
	case !isSudo:
		// Sanitize passed-in and parent policies before comparison
		sanitizedInputPolicies := policyutil.SanitizePolicies(data.Policies, policyutil.DoNotAddDefaultPolicy)
		sanitizedParentPolicies := policyutil.SanitizePolicies(parent.Policies, policyutil.DoNotAddDefaultPolicy)

		if !strutil.StrListSubset(sanitizedParentPolicies, sanitizedInputPolicies) {
			return logical.ErrorResponse("child policies must be subset of parent"), logical.ErrInvalidRequest
		}

		// If the parent has default, and they haven't requested not to get it,
		// add it. Note that if they have explicitly put "default" in
		// data.Policies it will still be added because NoDefaultPolicy
		// controls *automatic* adding.
		if !data.NoDefaultPolicy && strutil.StrListContains(parent.Policies, "default") {
			addDefault = true
		}

	// Add default by default in this case unless requested not to
	case isSudo:
		addDefault = !data.NoDefaultPolicy
	}

	te.Policies = policyutil.SanitizePolicies(data.Policies, addDefault)

	// Yes, this is a little inefficient to do it like this, but meh
	if data.NoDefaultPolicy {
		te.Policies = strutil.StrListDelete(te.Policies, "default")
	}

	// Prevent internal policies from being assigned to tokens
	for _, policy := range te.Policies {
		if strutil.StrListContains(nonAssignablePolicies, policy) {
			return logical.ErrorResponse(fmt.Sprintf("cannot assign policy %q", policy)), nil
		}
	}

	// Prevent attempts to create a root token without an actual root token as parent.
	// This is to thwart privilege escalation by tokens having 'sudo' privileges.
	if strutil.StrListContains(data.Policies, "root") && !strutil.StrListContains(parent.Policies, "root") {
		return logical.ErrorResponse("root tokens may not be created without parent token being root"), logical.ErrInvalidRequest
	}

	//
	// NOTE: Do not modify policies below this line. We need the checks above
	// to be the last checks as they must look at the final policy set.
	//

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

	if data.ExplicitMaxTTL != "" {
		dur, err := duration.ParseDurationSecond(data.ExplicitMaxTTL)
		if err != nil {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		if dur < 0 {
			return logical.ErrorResponse("explicit_max_ttl must be positive"), logical.ErrInvalidRequest
		}
		te.ExplicitMaxTTL = dur
	}

	var periodToUse time.Duration
	if data.Period != "" {
		if !isSudo {
			return logical.ErrorResponse("root or sudo privileges required to create periodic token"),
				logical.ErrInvalidRequest
		}
		dur, err := duration.ParseDurationSecond(data.Period)
		if err != nil {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		if dur < 0 {
			return logical.ErrorResponse("period must be positive"), logical.ErrInvalidRequest
		}
		te.Period = dur
		periodToUse = dur
	}

	// Parse the TTL/lease if any
	if data.TTL != "" {
		dur, err := duration.ParseDurationSecond(data.TTL)
		if err != nil {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		if dur < 0 {
			return logical.ErrorResponse("ttl must be positive"), logical.ErrInvalidRequest
		}
		te.TTL = dur
	} else if data.Lease != "" {
		// This block is compatibility
		dur, err := time.ParseDuration(data.Lease)
		if err != nil {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		if dur < 0 {
			return logical.ErrorResponse("lease must be positive"), logical.ErrInvalidRequest
		}
		te.TTL = dur
	}

	// Set the lesser period/explicit max TTL if defined both in arguments and in role
	if role != nil {
		if role.ExplicitMaxTTL != 0 {
			switch {
			case te.ExplicitMaxTTL == 0:
				te.ExplicitMaxTTL = role.ExplicitMaxTTL
			default:
				if role.ExplicitMaxTTL < te.ExplicitMaxTTL {
					te.ExplicitMaxTTL = role.ExplicitMaxTTL
				}
				resp.AddWarning(fmt.Sprintf("Explicit max TTL specified both during creation call and in role; using the lesser value of %d seconds", int64(te.ExplicitMaxTTL.Seconds())))
			}
		}
		if role.Period != 0 {
			switch {
			case periodToUse == 0:
				periodToUse = role.Period
			default:
				if role.Period < periodToUse {
					periodToUse = role.Period
				}
				resp.AddWarning(fmt.Sprintf("Period specified both during creation call and in role; using the lesser value of %d seconds", int64(periodToUse.Seconds())))
			}
		}
	}

	sysView := ts.System()

	if periodToUse > 0 {
		te.TTL = periodToUse
	} else {
		// Set the default lease if not provided, root tokens are exempt
		if te.TTL == 0 && !strutil.StrListContains(te.Policies, "root") {
			te.TTL = sysView.DefaultLeaseTTL()
		}

		// Limit the lease duration
		if te.TTL > sysView.MaxLeaseTTL() && sysView.MaxLeaseTTL() != 0 {
			te.TTL = sysView.MaxLeaseTTL()
		}
	}

	// Run some bounding checks if the explicit max TTL is set; we do not check
	// period as it's defined to escape the max TTL
	if te.ExplicitMaxTTL > 0 {
		// Limit the lease duration, except for periodic tokens -- in that case the explicit max limits the period, which itself can escape normal max
		if sysView.MaxLeaseTTL() != 0 && te.ExplicitMaxTTL > sysView.MaxLeaseTTL() && periodToUse == 0 {
			resp.AddWarning(fmt.Sprintf(
				"Explicit max TTL of %d seconds is greater than system/mount allowed value; value is being capped to %d seconds",
				int64(te.ExplicitMaxTTL.Seconds()), int64(sysView.MaxLeaseTTL().Seconds())))
			te.ExplicitMaxTTL = sysView.MaxLeaseTTL()
		}

		if te.TTL == 0 {
			// This won't be the case if it's periodic -- it will be set above
			te.TTL = te.ExplicitMaxTTL
		} else {
			// Limit even in the periodic case
			if te.TTL > te.ExplicitMaxTTL {
				resp.AddWarning(fmt.Sprintf(
					"Requested TTL of %d seconds higher than explicit max TTL; value being capped to %d seconds",
					int64(te.TTL.Seconds()), int64(te.ExplicitMaxTTL.Seconds())))
				te.TTL = te.ExplicitMaxTTL
			}
		}
	}

	// Don't advertise non-expiring root tokens as renewable, as attempts to renew them are denied
	if te.TTL == 0 {
		if parent.TTL != 0 {
			return logical.ErrorResponse("expiring root tokens cannot create non-expiring root tokens"), logical.ErrInvalidRequest
		}
		renewable = false
	}

	// Create the token
	if err := ts.create(&te); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Generate the response
	resp.Auth = &logical.Auth{
		DisplayName: te.DisplayName,
		Policies:    te.Policies,
		Metadata:    te.Meta,
		LeaseOptions: logical.LeaseOptions{
			TTL:       te.TTL,
			Renewable: renewable,
		},
		ClientToken: te.ID,
		Accessor:    te.Accessor,
	}

	if ts.policyLookupFunc != nil {
		for _, p := range te.Policies {
			policy, err := ts.policyLookupFunc(p)
			if err != nil {
				return logical.ErrorResponse(fmt.Sprintf("could not look up policy %s", p)), nil
			}
			if policy == nil {
				resp.AddWarning(fmt.Sprintf("Policy %q does not exist", p))
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
	var urltoken bool
	id := data.Get("token").(string)
	if id == "" {
		id = data.Get("urltoken").(string)
		if id == "" {
			return logical.ErrorResponse("missing token ID"), logical.ErrInvalidRequest
		}
		urltoken = true
	}

	// Revoke the token and its children
	if err := ts.RevokeTree(id); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	if urltoken {
		resp := &logical.Response{}
		resp.AddWarning(`Using a token in the path is unsafe as the token can be logged in many places. Please use POST or PUT with the token passed in via the "token" parameter.`)
		return resp, nil
	}

	return nil, nil
}

// handleRevokeOrphan handles the auth/token/revoke-orphan/id path for revocation of tokens
// in a way that leaves child tokens orphaned. Normally, using sys/revoke/leaseID will revoke
// the token and all children.
func (ts *TokenStore) handleRevokeOrphan(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var urltoken bool
	// Parse the id
	id := data.Get("token").(string)
	if id == "" {
		id = data.Get("urltoken").(string)
		if id == "" {
			return logical.ErrorResponse("missing token ID"), logical.ErrInvalidRequest
		}
		urltoken = true
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

	if urltoken {
		resp := &logical.Response{}
		resp.AddWarning(`Using a token in the path is unsafe as the token can be logged in many places. Please use POST or PUT with the token passed in via the "token" parameter.`)
		return resp, nil
	}

	return nil, nil
}

func (ts *TokenStore) handleLookupSelf(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	data.Raw["token"] = req.ClientToken
	return ts.handleLookup(req, data)
}

// handleLookup handles the auth/token/lookup/id path for querying information about
// a particular token. This can be used to see which policies are applicable.
func (ts *TokenStore) handleLookup(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var urltoken bool
	id := data.Get("token").(string)
	if id == "" {
		id = data.Get("urltoken").(string)
		if id != "" {
			urltoken = true
		}
	}
	if id == "" {
		id = req.ClientToken
	}
	if id == "" {
		return logical.ErrorResponse("missing token ID"), logical.ErrInvalidRequest
	}

	lock := ts.getTokenLock(id)
	lock.RLock()
	defer lock.RUnlock()

	// Lookup the token
	saltedId := ts.SaltID(id)
	out, err := ts.lookupSalted(saltedId, true)

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
			"id":               out.ID,
			"accessor":         out.Accessor,
			"policies":         out.Policies,
			"path":             out.Path,
			"meta":             out.Meta,
			"display_name":     out.DisplayName,
			"num_uses":         out.NumUses,
			"orphan":           false,
			"creation_time":    int64(out.CreationTime),
			"creation_ttl":     int64(out.TTL.Seconds()),
			"ttl":              int64(0),
			"explicit_max_ttl": int64(out.ExplicitMaxTTL.Seconds()),
		},
	}

	if out.Parent == "" {
		resp.Data["orphan"] = true
	}

	if out.Role != "" {
		resp.Data["role"] = out.Role
	}
	if out.Period != 0 {
		resp.Data["period"] = int64(out.Period.Seconds())
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
		if err := leaseTimes.renewable(); err == nil {
			resp.Data["renewable"] = true
		} else {
			resp.Data["renewable"] = false
		}
	}

	if urltoken {
		resp.AddWarning(`Using a token in the path is unsafe as the token can be logged in many places. Please use POST or PUT with the token passed in via the "token" parameter.`)
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
	var urltoken bool
	id := data.Get("token").(string)
	if id == "" {
		id = data.Get("urltoken").(string)
		if id == "" {
			return logical.ErrorResponse("missing token ID"), logical.ErrInvalidRequest
		}
		urltoken = true
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
	resp, err := ts.expiration.RenewToken(req, te.Path, te.ID, increment)

	if urltoken {
		resp.AddWarning(`Using a token in the path is unsafe as the token can be logged in many places. Please use POST or PUT with the token passed in via the "token" parameter.`)
	}

	return resp, err
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

	te, err := ts.Lookup(req.Auth.ClientToken)
	if err != nil {
		return nil, fmt.Errorf("error looking up token: %s", err)
	}
	if te == nil {
		return nil, fmt.Errorf("no token entry found during lookup")
	}

	f := framework.LeaseExtend(req.Auth.Increment, te.ExplicitMaxTTL, ts.System())

	// If (te/role).Period is not zero, this is a periodic token. The TTL for a
	// periodic token is always the same (the period value). It is not subject
	// to normal maximum TTL checks that would come from calling LeaseExtend,
	// so we fast path it.
	//
	// The one wrinkle here is if the token has an explicit max TTL. If both
	// are set, we treat it as a regular token and use the periodic value as
	// the increment.

	// No role? Use normal LeaseExtend semantics, taking into account
	// TokenEntry properties
	if te.Role == "" {
		//Explicit max TTL overrides the period, if both are set
		if te.Period != 0 {
			if te.ExplicitMaxTTL == 0 {
				req.Auth.TTL = te.Period
				return &logical.Response{Auth: req.Auth}, nil
			} else {
				maxTime := time.Unix(te.CreationTime, 0).Add(te.ExplicitMaxTTL)
				if time.Now().Add(te.Period).After(maxTime) {
					req.Auth.TTL = maxTime.Sub(time.Now())
				} else {
					req.Auth.TTL = te.Period
				}
				return &logical.Response{Auth: req.Auth}, nil
			}
		}
		return f(req, d)
	}

	role, err := ts.tokenStoreRole(te.Role)
	if err != nil {
		return nil, fmt.Errorf("error looking up role %s: %s", te.Role, err)
	}

	if role == nil {
		return nil, fmt.Errorf("original token role (%s) could not be found, not renewing", te.Role)
	}

	// Same deal here, but using the role period
	if role.Period != 0 {
		periodToUse := role.Period
		if te.Period > 0 && te.Period < role.Period {
			periodToUse = te.Period
		}
		if te.ExplicitMaxTTL == 0 {
			req.Auth.TTL = periodToUse
			return &logical.Response{Auth: req.Auth}, nil
		} else {
			maxTime := time.Unix(te.CreationTime, 0).Add(te.ExplicitMaxTTL)
			if time.Now().Add(periodToUse).After(maxTime) {
				req.Auth.TTL = maxTime.Sub(time.Now())
			} else {
				req.Auth.TTL = periodToUse
			}
			return &logical.Response{Auth: req.Auth}, nil
		}
	}

	return f(req, d)
}

func (ts *TokenStore) tokenStoreRole(name string) (*tsRoleEntry, error) {
	entry, err := ts.view.Get(fmt.Sprintf("%s%s", rolesPrefix, name))
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
	entries, err := ts.view.List(rolesPrefix)
	if err != nil {
		return nil, err
	}

	ret := make([]string, len(entries))
	for i, entry := range entries {
		ret[i] = strings.TrimPrefix(entry, rolesPrefix)
	}

	return logical.ListResponse(ret), nil
}

func (ts *TokenStore) tokenStoreRoleDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	err := ts.view.Delete(fmt.Sprintf("%s%s", rolesPrefix, data.Get("role_name").(string)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ts *TokenStore) tokenStoreRoleRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	role, err := ts.tokenStoreRole(data.Get("role_name").(string))
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"period":              int64(role.Period.Seconds()),
			"explicit_max_ttl":    int64(role.ExplicitMaxTTL.Seconds()),
			"disallowed_policies": role.DisallowedPolicies,
			"allowed_policies":    role.AllowedPolicies,
			"name":                role.Name,
			"orphan":              role.Orphan,
			"path_suffix":         role.PathSuffix,
			"renewable":           role.Renewable,
		},
	}

	return resp, nil
}

func (ts *TokenStore) tokenStoreRoleExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	name := data.Get("role_name").(string)
	if name == "" {
		return false, fmt.Errorf("role name cannot be empty")
	}
	role, err := ts.tokenStoreRole(name)
	if err != nil {
		return false, err
	}

	return role != nil, nil
}

func (ts *TokenStore) tokenStoreRoleCreateUpdate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("role_name").(string)
	if name == "" {
		return logical.ErrorResponse("role name cannot be empty"), nil
	}
	entry, err := ts.tokenStoreRole(name)
	if err != nil {
		return nil, err
	}

	// Due to the existence check, entry will only be nil if it's a create
	// operation, so just create a new one
	if entry == nil {
		entry = &tsRoleEntry{
			Name: name,
		}
	}

	// In this series of blocks, if we do not find a user-provided value and
	// it's a creation operation, we call data.Get to get the appropriate
	// default

	orphanInt, ok := data.GetOk("orphan")
	if ok {
		entry.Orphan = orphanInt.(bool)
	} else if req.Operation == logical.CreateOperation {
		entry.Orphan = data.Get("orphan").(bool)
	}

	periodInt, ok := data.GetOk("period")
	if ok {
		entry.Period = time.Second * time.Duration(periodInt.(int))
	} else if req.Operation == logical.CreateOperation {
		entry.Period = time.Second * time.Duration(data.Get("period").(int))
	}

	renewableInt, ok := data.GetOk("renewable")
	if ok {
		entry.Renewable = renewableInt.(bool)
	} else if req.Operation == logical.CreateOperation {
		entry.Renewable = data.Get("renewable").(bool)
	}

	var resp *logical.Response

	explicitMaxTTLInt, ok := data.GetOk("explicit_max_ttl")
	if ok {
		entry.ExplicitMaxTTL = time.Second * time.Duration(explicitMaxTTLInt.(int))
	} else if req.Operation == logical.CreateOperation {
		entry.ExplicitMaxTTL = time.Second * time.Duration(data.Get("explicit_max_ttl").(int))
	}
	if entry.ExplicitMaxTTL != 0 {
		sysView := ts.System()

		if sysView.MaxLeaseTTL() != time.Duration(0) && entry.ExplicitMaxTTL > sysView.MaxLeaseTTL() {
			if resp == nil {
				resp = &logical.Response{}
			}
			resp.AddWarning(fmt.Sprintf(
				"Given explicit max TTL of %d is greater than system/mount allowed value of %d seconds; until this is fixed attempting to create tokens against this role will result in an error",
				entry.ExplicitMaxTTL.Seconds(), sysView.MaxLeaseTTL().Seconds()))
		}
	}

	pathSuffixInt, ok := data.GetOk("path_suffix")
	if ok {
		pathSuffix := pathSuffixInt.(string)
		if pathSuffix != "" {
			matched := pathSuffixSanitize.MatchString(pathSuffix)
			if !matched {
				return logical.ErrorResponse(fmt.Sprintf(
					"given role path suffix contains invalid characters; must match %s",
					pathSuffixSanitize.String())), nil
			}
			entry.PathSuffix = pathSuffix
		}
	} else if req.Operation == logical.CreateOperation {
		entry.PathSuffix = data.Get("path_suffix").(string)
	}

	allowedPoliciesStr, ok := data.GetOk("allowed_policies")
	if ok {
		entry.AllowedPolicies = policyutil.SanitizePolicies(strings.Split(allowedPoliciesStr.(string), ","), policyutil.DoNotAddDefaultPolicy)
	} else if req.Operation == logical.CreateOperation {
		entry.AllowedPolicies = policyutil.SanitizePolicies(strings.Split(data.Get("allowed_policies").(string), ","), policyutil.DoNotAddDefaultPolicy)
	}

	disallowedPoliciesStr, ok := data.GetOk("disallowed_policies")
	if ok {
		entry.DisallowedPolicies = policyutil.SanitizePolicies(strings.Split(disallowedPoliciesStr.(string), ","), policyutil.DoNotAddDefaultPolicy)
	} else if req.Operation == logical.CreateOperation {
		entry.DisallowedPolicies = policyutil.SanitizePolicies(strings.Split(data.Get("disallowed_policies").(string), ","), policyutil.DoNotAddDefaultPolicy)
	}

	// Store it
	jsonEntry, err := logical.StorageEntryJSON(fmt.Sprintf("%s%s", rolesPrefix, name), entry)
	if err != nil {
		return nil, err
	}
	if err := ts.view.Put(jsonEntry); err != nil {
		return nil, err
	}

	return resp, nil
}

const (
	tokenTidyHelp = `
This endpoint performs cleanup tasks that can be run if certain error
conditions have occurred.
`
	tokenTidyDesc = `
This endpoint performs cleanup tasks that can be run to clean up token and
lease entries after certain error conditions. Usually running this is not
necessary, and is only required if upgrade notes or support personnel suggest
it.
`
	tokenBackendHelp = `The token credential backend is always enabled and builtin to Vault.
Client tokens are used to identify a client and to allow Vault to associate policies and ACLs
which are enforced on every request. This backend also allows for generating sub-tokens as well
as revocation of tokens. The tokens are renewable if associated with a lease.`
	tokenCreateHelp          = `The token create path is used to create new tokens.`
	tokenCreateOrphanHelp    = `The token create path is used to create new orphan tokens.`
	tokenCreateRoleHelp      = `This token create path is used to create new tokens adhering to the given role.`
	tokenListRolesHelp       = `This endpoint lists configured roles.`
	tokenLookupAccessorHelp  = `This endpoint will lookup a token associated with the given accessor and its properties. Response will not contain the token ID.`
	tokenLookupHelp          = `This endpoint will lookup a token and its properties.`
	tokenPathRolesHelp       = `This endpoint allows creating, reading, and deleting roles.`
	tokenRevokeAccessorHelp  = `This endpoint will delete the token associated with the accessor and all of its child tokens.`
	tokenRevokeHelp          = `This endpoint will delete the given token and all of its child tokens.`
	tokenRevokeSelfHelp      = `This endpoint will delete the token used to call it and all of its child tokens.`
	tokenRevokeOrphanHelp    = `This endpoint will delete the token and orphan its child tokens.`
	tokenRenewHelp           = `This endpoint will renew the given token and prevent expiration.`
	tokenRenewSelfHelp       = `This endpoint will renew the token used to call it and prevent expiration.`
	tokenAllowedPoliciesHelp = `If set, tokens can be created with any subset of the policies in this
list, rather than the normal semantics of tokens being a subset of the
calling token's policies. The parameter is a comma-delimited string of
policy names.`
	tokenDisallowedPoliciesHelp = `If set, successful token creation via this role will require that
no policies in the given list are requested. The parameter is a comma-delimited string of policy names.`
	tokenOrphanHelp = `If true, tokens created via this role
will be orphan tokens (have no parent)`
	tokenPeriodHelp = `If set, tokens created via this role
will have no max lifetime; instead, their
renewal period will be fixed to this value.
This takes an integer number of seconds,
or a string duration (e.g. "24h").`
	tokenPathSuffixHelp = `If set, tokens created via this role
will contain the given suffix as a part of
their path. This can be used to assist use
of the 'revoke-prefix' endpoint later on.
The given suffix must match the regular
expression.`
	tokenExplicitMaxTTLHelp = `If set, tokens created via this role
carry an explicit maximum TTL. During renewal,
the current maximum TTL values of the role
and the mount are not checked for changes,
and any updates to these values will have
no effect on the token being renewed.`
	tokenRenewableHelp = `Tokens created via this role will be
renewable or not according to this value.
Defaults to "true".`
	tokenListAccessorsHelp = `List token accessors, which can then be
be used to iterate and discover their properities
or revoke them. Because this can be used to
cause a denial of service, this endpoint
requires 'sudo' capability in addition to
'list'.`
)
