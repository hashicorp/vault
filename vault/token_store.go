package vault

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"regexp"
	"strings"
	"time"

	proto "github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	sockaddr "github.com/hashicorp/go-sockaddr"

	metrics "github.com/armon/go-metrics"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/helper/base62"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/locksutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/hashicorp/vault/logical/plugin/pb"
	"github.com/mitchellh/mapstructure"
)

const (
	// idPrefix is the prefix used to store tokens for their
	// primary ID based index
	idPrefix = "id/"

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

	// tokenRevocationPending indicates that the token should not be used
	// again. If this is encountered during an existing request flow, it means
	// that the token is but is currently fulfilling its final use; after this
	// request it will not be able to be looked up as being valid.
	tokenRevocationPending = -1
)

var (
	// TokenLength is the size of tokens we are currently generating, without
	// any namespace information
	TokenLength = 24

	// displayNameSanitize is used to sanitize a display name given to a token.
	displayNameSanitize = regexp.MustCompile("[^a-zA-Z0-9-]")

	// pathSuffixSanitize is used to ensure a path suffix in a role is valid.
	pathSuffixSanitize = regexp.MustCompile("\\w[\\w-.]+\\w")

	destroyCubbyhole = func(ctx context.Context, ts *TokenStore, te *logical.TokenEntry) error {
		if ts.cubbyholeBackend == nil {
			// Should only ever happen in testing
			return nil
		}

		if te == nil {
			return errors.New("nil token entry")
		}

		switch {
		case te.NamespaceID == namespace.RootNamespaceID && !strings.HasPrefix(te.ID, "s."):
			saltedID, err := ts.SaltID(ctx, te.ID)
			if err != nil {
				return err
			}
			return ts.cubbyholeBackend.revoke(ctx, salt.SaltID(ts.cubbyholeBackend.saltUUID, saltedID, salt.SHA1Hash))

		default:
			if te.CubbyholeID == "" {
				return fmt.Errorf("missing cubbyhole ID while destroying")
			}
			return ts.cubbyholeBackend.revoke(ctx, te.CubbyholeID)
		}
	}
)

func (ts *TokenStore) paths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "roles/?$",

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: ts.tokenStoreRoleList,
			},

			HelpSynopsis:    tokenListRolesHelp,
			HelpDescription: tokenListRolesHelp,
		},

		{
			Pattern: "accessors/$",

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: ts.tokenStoreAccessorList,
			},

			HelpSynopsis:    tokenListAccessorsHelp,
			HelpDescription: tokenListAccessorsHelp,
		},

		{
			Pattern: "roles/" + framework.GenericNameRegex("role_name"),
			Fields: map[string]*framework.FieldSchema{
				"role_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role",
				},

				"allowed_policies": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: tokenAllowedPoliciesHelp,
				},

				"disallowed_policies": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
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

				"bound_cidrs": &framework.FieldSchema{
					Type:        framework.TypeCommaStringSlice,
					Description: `Comma separated string or JSON list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.`,
				},

				"token_type": &framework.FieldSchema{
					Type:        framework.TypeString,
					Default:     "service",
					Description: "The type of token to generate, service or batch",
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   ts.tokenStoreRoleRead,
				logical.CreateOperation: ts.tokenStoreRoleCreateUpdate,
				logical.UpdateOperation: ts.tokenStoreRoleCreateUpdate,
				logical.DeleteOperation: ts.tokenStoreRoleDelete,
			},

			ExistenceCheck: ts.tokenStoreRoleExistenceCheck,
		},

		{
			Pattern: "create-orphan$",

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: ts.handleCreateOrphan,
			},

			HelpSynopsis:    strings.TrimSpace(tokenCreateOrphanHelp),
			HelpDescription: strings.TrimSpace(tokenCreateOrphanHelp),
		},

		{
			Pattern: "create/" + framework.GenericNameRegex("role_name"),

			Fields: map[string]*framework.FieldSchema{
				"role_name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the role",
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: ts.handleCreateAgainstRole,
			},

			HelpSynopsis:    strings.TrimSpace(tokenCreateRoleHelp),
			HelpDescription: strings.TrimSpace(tokenCreateRoleHelp),
		},

		{
			Pattern: "create$",

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: ts.handleCreate,
			},

			HelpSynopsis:    strings.TrimSpace(tokenCreateHelp),
			HelpDescription: strings.TrimSpace(tokenCreateHelp),
		},

		{
			Pattern: "lookup",

			Fields: map[string]*framework.FieldSchema{
				"token": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Token to lookup (POST request body)",
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   ts.handleLookup,
				logical.UpdateOperation: ts.handleLookup,
			},

			HelpSynopsis:    strings.TrimSpace(tokenLookupHelp),
			HelpDescription: strings.TrimSpace(tokenLookupHelp),
		},

		{
			Pattern: "lookup-accessor",

			Fields: map[string]*framework.FieldSchema{
				"accessor": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Accessor of the token to look up (request body)",
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: ts.handleUpdateLookupAccessor,
			},

			HelpSynopsis:    strings.TrimSpace(tokenLookupAccessorHelp),
			HelpDescription: strings.TrimSpace(tokenLookupAccessorHelp),
		},

		{
			Pattern: "lookup-self$",

			Fields: map[string]*framework.FieldSchema{
				"token": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Token to look up (unused, does not need to be set)",
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: ts.handleLookupSelf,
				logical.ReadOperation:   ts.handleLookupSelf,
			},

			HelpSynopsis:    strings.TrimSpace(tokenLookupHelp),
			HelpDescription: strings.TrimSpace(tokenLookupHelp),
		},

		{
			Pattern: "revoke-accessor",

			Fields: map[string]*framework.FieldSchema{
				"accessor": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Accessor of the token (request body)",
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: ts.handleUpdateRevokeAccessor,
			},

			HelpSynopsis:    strings.TrimSpace(tokenRevokeAccessorHelp),
			HelpDescription: strings.TrimSpace(tokenRevokeAccessorHelp),
		},

		{
			Pattern: "revoke-self$",

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: ts.handleRevokeSelf,
			},

			HelpSynopsis:    strings.TrimSpace(tokenRevokeSelfHelp),
			HelpDescription: strings.TrimSpace(tokenRevokeSelfHelp),
		},

		{
			Pattern: "revoke",

			Fields: map[string]*framework.FieldSchema{
				"token": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Token to revoke (request body)",
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: ts.handleRevokeTree,
			},

			HelpSynopsis:    strings.TrimSpace(tokenRevokeHelp),
			HelpDescription: strings.TrimSpace(tokenRevokeHelp),
		},

		{
			Pattern: "revoke-orphan",

			Fields: map[string]*framework.FieldSchema{
				"token": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Token to revoke (request body)",
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: ts.handleRevokeOrphan,
			},

			HelpSynopsis:    strings.TrimSpace(tokenRevokeOrphanHelp),
			HelpDescription: strings.TrimSpace(tokenRevokeOrphanHelp),
		},

		{
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
				logical.UpdateOperation: ts.handleRenewSelf,
			},

			HelpSynopsis:    strings.TrimSpace(tokenRenewSelfHelp),
			HelpDescription: strings.TrimSpace(tokenRenewSelfHelp),
		},

		{
			Pattern: "renew",

			Fields: map[string]*framework.FieldSchema{
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
				logical.UpdateOperation: ts.handleRenew,
			},

			HelpSynopsis:    strings.TrimSpace(tokenRenewHelp),
			HelpDescription: strings.TrimSpace(tokenRenewHelp),
		},

		{
			Pattern: "tidy$",

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: ts.handleTidy,
			},

			HelpSynopsis:    strings.TrimSpace(tokenTidyHelp),
			HelpDescription: strings.TrimSpace(tokenTidyDesc),
		},
	}
}

// LookupToken returns the properties of the token from the token store. This
// is particularly useful to fetch the accessor of the client token and get it
// populated in the logical request along with the client token. The accessor
// of the client token can get audit logged.
func (c *Core) LookupToken(ctx context.Context, token string) (*logical.TokenEntry, error) {
	if c.Sealed() {
		return nil, consts.ErrSealed
	}

	c.stateLock.RLock()
	defer c.stateLock.RUnlock()

	if c.standby && !c.perfStandby {
		return nil, consts.ErrStandby
	}

	// Many tests don't have a token store running
	if c.tokenStore == nil || c.tokenStore.expiration == nil {
		return nil, nil
	}

	return c.tokenStore.Lookup(ctx, token)
}

// TokenStore is used to manage client tokens. Tokens are used for
// clients to authenticate, and each token is mapped to an applicable
// set of policy which is used for authorization.
type TokenStore struct {
	*framework.Backend

	activeContext context.Context

	core *Core

	batchTokenEncryptor BarrierEncryptor

	baseBarrierView     *BarrierView
	idBarrierView       *BarrierView
	accessorBarrierView *BarrierView
	parentBarrierView   *BarrierView
	rolesBarrierView    *BarrierView

	expiration *ExpirationManager

	cubbyholeBackend *CubbyholeBackend

	tokenLocks []*locksutil.LockEntry

	// tokenPendingDeletion stores tokens that are being revoked. If the token is
	// not in the map, it means that there's no deletion in progress. If the value
	// is true it means deletion is in progress, and if false it means deletion
	// failed. Revocation needs to handle these states accordingly.
	tokensPendingDeletion *sync.Map

	cubbyholeDestroyer func(context.Context, *TokenStore, *logical.TokenEntry) error

	logger log.Logger

	saltLock sync.RWMutex
	salts    map[string]*salt.Salt

	tidyLock *uint32

	identityPoliciesDeriverFunc func(string) (*identity.Entity, []string, error)

	quitContext context.Context
}

// NewTokenStore is used to construct a token store that is
// backed by the given barrier view.
func NewTokenStore(ctx context.Context, logger log.Logger, core *Core, config *logical.BackendConfig) (*TokenStore, error) {
	// Create a sub-view
	view := core.systemBarrierView.SubView(tokenSubPath)

	// Initialize the store
	t := &TokenStore{
		activeContext:         ctx,
		core:                  core,
		batchTokenEncryptor:   core.barrier,
		baseBarrierView:       view,
		idBarrierView:         view.SubView(idPrefix),
		accessorBarrierView:   view.SubView(accessorPrefix),
		parentBarrierView:     view.SubView(parentPrefix),
		rolesBarrierView:      view.SubView(rolesPrefix),
		cubbyholeDestroyer:    destroyCubbyhole,
		logger:                logger,
		tokenLocks:            locksutil.CreateLocks(),
		tokensPendingDeletion: &sync.Map{},
		saltLock:              sync.RWMutex{},
		tidyLock:              new(uint32),
		quitContext:           core.activeContext,
		salts:                 make(map[string]*salt.Salt),
	}

	// Setup the framework endpoints
	t.Backend = &framework.Backend{
		AuthRenew: t.authRenew,

		PathsSpecial: &logical.Paths{
			Root: []string{
				"revoke-orphan/*",
				"accessors*",
			},

			// Most token store items are local since tokens are local, but a
			// notable exception is roles
			LocalStorage: []string{
				idPrefix,
				accessorPrefix,
				parentPrefix,
				salt.DefaultLocation,
			},
		},
		BackendType: logical.TypeCredential,
	}

	t.Backend.Paths = append(t.Backend.Paths, t.paths()...)

	t.Backend.Setup(ctx, config)

	return t, nil
}

func (ts *TokenStore) Invalidate(ctx context.Context, key string) {
	//ts.logger.Debug("invalidating key", "key", key)

	switch key {
	case tokenSubPath + salt.DefaultLocation:
		ts.saltLock.Lock()
		ts.salts = make(map[string]*salt.Salt)
		ts.saltLock.Unlock()
	}
}

func (ts *TokenStore) Salt(ctx context.Context) (*salt.Salt, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	ts.saltLock.RLock()
	if salt, ok := ts.salts[ns.ID]; ok {
		defer ts.saltLock.RUnlock()
		return salt, nil
	}
	ts.saltLock.RUnlock()
	ts.saltLock.Lock()
	defer ts.saltLock.Unlock()
	if salt, ok := ts.salts[ns.ID]; ok {
		return salt, nil
	}

	salt, err := salt.NewSalt(ctx, ts.baseView(ns), &salt.Config{
		HashFunc: salt.SHA1Hash,
		Location: salt.DefaultLocation,
	})
	if err != nil {
		return nil, err
	}
	ts.salts[ns.ID] = salt
	return salt, nil
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

	// The set of CIDRs that tokens generated using this role will be bound to
	BoundCIDRs []*sockaddr.SockAddrMarshaler `json:"bound_cidrs"`

	// The type of token this role should issue
	TokenType logical.TokenType `json:"token_type" mapstructure:"token_type"`
}

type accessorEntry struct {
	TokenID     string `json:"token_id"`
	AccessorID  string `json:"accessor_id"`
	NamespaceID string `json:"namespace_id"`
}

// SetExpirationManager is used to provide the token store with
// an expiration manager. This is used to manage prefix based revocation
// of tokens and to tidy entries when removed from the token store.
func (ts *TokenStore) SetExpirationManager(exp *ExpirationManager) {
	ts.expiration = exp
}

// SaltID is used to apply a salt and hash to an ID to make sure its not reversible
func (ts *TokenStore) SaltID(ctx context.Context, id string) (string, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return "", namespace.ErrNoNamespace
	}

	s, err := ts.Salt(ctx)
	if err != nil {
		return "", err
	}

	// For tokens of older format and belonging to the root namespace, use SHA1
	// hash for salting.
	if ns.ID == namespace.RootNamespaceID && !strings.Contains(id, ".") {
		return s.SaltID(id), nil
	}

	// For all other tokens, use SHA2-256 HMAC for salting. This includes
	// tokens of older format, but belonging to a namespace other than the root
	// namespace.
	return "h" + s.GetHMAC(id), nil
}

// rootToken is used to generate a new token with root privileges and no parent
func (ts *TokenStore) rootToken(ctx context.Context) (*logical.TokenEntry, error) {
	ctx = namespace.ContextWithNamespace(ctx, namespace.RootNamespace)
	te := &logical.TokenEntry{
		Policies:     []string{"root"},
		Path:         "auth/token/root",
		DisplayName:  "root",
		CreationTime: time.Now().Unix(),
		NamespaceID:  namespace.RootNamespaceID,
		Type:         logical.TokenTypeService,
	}
	if err := ts.create(ctx, te); err != nil {
		return nil, err
	}
	return te, nil
}

func (ts *TokenStore) tokenStoreAccessorList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	nsID := ns.ID

	entries, err := ts.accessorView(ns).List(ctx, "")
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{}

	ret := make([]string, 0, len(entries))
	for _, entry := range entries {
		aEntry, err := ts.lookupByAccessor(ctx, entry, true, false)
		if err != nil {
			resp.AddWarning(fmt.Sprintf("Found an accessor entry that could not be successfully decoded; associated error is %q", err.Error()))
			continue
		}

		if aEntry.TokenID == "" {
			resp.AddWarning(fmt.Sprintf("Found an accessor entry missing a token: %v", aEntry.AccessorID))
			continue
		}

		if aEntry.NamespaceID == nsID {
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
func (ts *TokenStore) createAccessor(ctx context.Context, entry *logical.TokenEntry) error {
	defer metrics.MeasureSince([]string{"token", "createAccessor"}, time.Now())

	var err error
	// Create a random accessor
	entry.Accessor, err = base62.Random(TokenLength)
	if err != nil {
		return err
	}

	tokenNS, err := NamespaceByID(ctx, entry.NamespaceID, ts.core)
	if err != nil {
		return err
	}
	if tokenNS == nil {
		return namespace.ErrNoNamespace
	}

	if tokenNS.ID != namespace.RootNamespaceID {
		entry.Accessor = fmt.Sprintf("%s.%s", entry.Accessor, tokenNS.ID)
	}

	// Create index entry, mapping the accessor to the token ID
	saltCtx := namespace.ContextWithNamespace(ctx, tokenNS)
	saltID, err := ts.SaltID(saltCtx, entry.Accessor)
	if err != nil {
		return err
	}

	aEntry := &accessorEntry{
		TokenID:     entry.ID,
		AccessorID:  entry.Accessor,
		NamespaceID: entry.NamespaceID,
	}

	aEntryBytes, err := jsonutil.EncodeJSON(aEntry)
	if err != nil {
		return errwrap.Wrapf("failed to marshal accessor index entry: {{err}}", err)
	}

	le := &logical.StorageEntry{Key: saltID, Value: aEntryBytes}
	if err := ts.accessorView(tokenNS).Put(ctx, le); err != nil {
		return errwrap.Wrapf("failed to persist accessor index entry: {{err}}", err)
	}
	return nil
}

// Create is used to create a new token entry. The entry is assigned
// a newly generated ID if not provided.
func (ts *TokenStore) create(ctx context.Context, entry *logical.TokenEntry) error {
	defer metrics.MeasureSince([]string{"token", "create"}, time.Now())

	tokenNS, err := NamespaceByID(ctx, entry.NamespaceID, ts.core)
	if err != nil {
		return err
	}
	if tokenNS == nil {
		return namespace.ErrNoNamespace
	}

	entry.Policies = policyutil.SanitizePolicies(entry.Policies, policyutil.DoNotAddDefaultPolicy)

	switch entry.Type {
	case logical.TokenTypeDefault, logical.TokenTypeService:
		// In case it was default, force to service
		entry.Type = logical.TokenTypeService

		// Generate an ID if necessary
		userSelectedID := true
		if entry.ID == "" {
			userSelectedID = false
			var err error
			entry.ID, err = base62.Random(TokenLength)
			if err != nil {
				return err
			}
		}

		if userSelectedID && strings.HasPrefix(entry.ID, "s.") {
			return fmt.Errorf("custom token ID cannot have the 's.' prefix")
		}

		if !userSelectedID {
			entry.ID = fmt.Sprintf("s.%s", entry.ID)
		}

		// Attach namespace ID for tokens that are not belonging to the root
		// namespace
		if tokenNS.ID != namespace.RootNamespaceID {
			entry.ID = fmt.Sprintf("%s.%s", entry.ID, tokenNS.ID)
		}

		if tokenNS.ID != namespace.RootNamespaceID || strings.HasPrefix(entry.ID, "s.") {
			if entry.CubbyholeID == "" {
				cubbyholeID, err := base62.Random(TokenLength)
				if err != nil {
					return err
				}
				entry.CubbyholeID = cubbyholeID
			}
		}

		// If the user didn't specifically pick the ID, e.g. because they were
		// sudo/root, check for collision; otherwise trust the process
		if userSelectedID {
			exist, _ := ts.lookupInternal(ctx, entry.ID, false, true)
			if exist != nil {
				return fmt.Errorf("cannot create a token with a duplicate ID")
			}
		}

		err = ts.createAccessor(ctx, entry)
		if err != nil {
			return err
		}

		return ts.storeCommon(ctx, entry, true)

	case logical.TokenTypeBatch:
		// Ensure fields we don't support/care about are nilled, proto marshal,
		// encrypt, skip persistence
		entry.ID = ""
		pEntry := &pb.TokenEntry{
			Parent:       entry.Parent,
			Policies:     entry.Policies,
			Path:         entry.Path,
			Meta:         entry.Meta,
			DisplayName:  entry.DisplayName,
			CreationTime: entry.CreationTime,
			TTL:          int64(entry.TTL),
			Role:         entry.Role,
			EntityID:     entry.EntityID,
			NamespaceID:  entry.NamespaceID,
			Type:         uint32(entry.Type),
		}

		boundCIDRs := make([]string, len(entry.BoundCIDRs))
		for i, cidr := range entry.BoundCIDRs {
			boundCIDRs[i] = cidr.String()
		}
		pEntry.BoundCIDRs = boundCIDRs

		mEntry, err := proto.Marshal(pEntry)
		if err != nil {
			return err
		}

		eEntry, err := ts.batchTokenEncryptor.Encrypt(ctx, "", mEntry)
		if err != nil {
			return err
		}

		bEntry := base64.RawURLEncoding.EncodeToString(eEntry)
		entry.ID = fmt.Sprintf("b.%s", bEntry)

		if tokenNS.ID != namespace.RootNamespaceID {
			entry.ID = fmt.Sprintf("%s.%s", entry.ID, tokenNS.ID)
		}

		return nil

	default:
		return fmt.Errorf("cannot create a token of type %d", entry.Type)
	}
}

// Store is used to store an updated token entry without writing the
// secondary index.
func (ts *TokenStore) store(ctx context.Context, entry *logical.TokenEntry) error {
	defer metrics.MeasureSince([]string{"token", "store"}, time.Now())
	return ts.storeCommon(ctx, entry, false)
}

// storeCommon handles the actual storage of an entry, possibly generating
// secondary indexes
func (ts *TokenStore) storeCommon(ctx context.Context, entry *logical.TokenEntry, writeSecondary bool) error {
	tokenNS, err := NamespaceByID(ctx, entry.NamespaceID, ts.core)
	if err != nil {
		return err
	}
	if tokenNS == nil {
		return namespace.ErrNoNamespace
	}

	saltCtx := namespace.ContextWithNamespace(ctx, tokenNS)
	saltedID, err := ts.SaltID(saltCtx, entry.ID)
	if err != nil {
		return err
	}

	// Marshal the entry
	enc, err := json.Marshal(entry)
	if err != nil {
		return errwrap.Wrapf("failed to encode entry: {{err}}", err)
	}

	if writeSecondary {
		// Write the secondary index if necessary. This is done before the
		// primary index because we'd rather have a dangling pointer with
		// a missing primary instead of missing the parent index and potentially
		// escaping the revocation chain.
		if entry.Parent != "" {
			// Ensure the parent exists
			parent, err := ts.Lookup(ctx, entry.Parent)
			if err != nil {
				return errwrap.Wrapf("failed to lookup parent: {{err}}", err)
			}
			if parent == nil {
				return fmt.Errorf("parent token not found")
			}

			parentNS, err := NamespaceByID(ctx, parent.NamespaceID, ts.core)
			if err != nil {
				return err
			}
			if parentNS == nil {
				return namespace.ErrNoNamespace
			}

			parentCtx := namespace.ContextWithNamespace(ctx, parentNS)

			// Create the index entry
			parentSaltedID, err := ts.SaltID(parentCtx, entry.Parent)
			if err != nil {
				return err
			}

			path := parentSaltedID + "/" + saltedID
			if tokenNS.ID != namespace.RootNamespaceID {
				path = fmt.Sprintf("%s.%s", path, tokenNS.ID)
			}

			le := &logical.StorageEntry{Key: path}
			if err := ts.parentView(parentNS).Put(ctx, le); err != nil {
				return errwrap.Wrapf("failed to persist entry: {{err}}", err)
			}
		}
	}

	// Write the primary ID
	le := &logical.StorageEntry{Key: saltedID, Value: enc}
	if len(entry.Policies) == 1 && entry.Policies[0] == "root" {
		le.SealWrap = true
	}
	if err := ts.idView(tokenNS).Put(ctx, le); err != nil {
		return errwrap.Wrapf("failed to persist entry: {{err}}", err)
	}
	return nil
}

// UseToken is used to manage restricted use tokens and decrement their
// available uses. Returns two values: a potentially updated entry or, if the
// token has been revoked, nil; and whether an error was encountered. The
// locking here isn't perfect, as other parts of the code may update an entry,
// but usually none after the entry is already created...so this is pretty
// good.
func (ts *TokenStore) UseToken(ctx context.Context, te *logical.TokenEntry) (*logical.TokenEntry, error) {
	if te == nil {
		return nil, fmt.Errorf("invalid token entry provided for use count decrementing")
	}

	// This case won't be hit with a token with restricted uses because we go
	// from 1 to -1. So it's a nice optimization to check this without a read
	// lock.
	if te.NumUses == 0 {
		return te, nil
	}

	// If we are attempting to unwrap a control group request, don't use the token.
	// It will be manually revoked by the handler.
	if len(te.Policies) == 1 && te.Policies[0] == controlGroupPolicyName {
		return te, nil
	}

	lock := locksutil.LockForKey(ts.tokenLocks, te.ID)
	lock.Lock()
	defer lock.Unlock()

	var err error
	te, err = ts.lookupInternal(ctx, te.ID, false, false)
	if err != nil {
		return nil, errwrap.Wrapf("failed to refresh entry: {{err}}", err)
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
		te.NumUses = tokenRevocationPending
	} else {
		te.NumUses--
	}

	err = ts.store(ctx, te)
	if err != nil {
		return nil, err
	}

	return te, nil
}

func (ts *TokenStore) UseTokenByID(ctx context.Context, id string) (*logical.TokenEntry, error) {
	te, err := ts.Lookup(ctx, id)
	if err != nil {
		return te, err
	}

	return ts.UseToken(ctx, te)
}

// Lookup is used to find a token given its ID. It acquires a read lock, then calls lookupInternal.
func (ts *TokenStore) Lookup(ctx context.Context, id string) (*logical.TokenEntry, error) {
	defer metrics.MeasureSince([]string{"token", "lookup"}, time.Now())
	if id == "" {
		return nil, fmt.Errorf("cannot lookup blank token")
	}

	// If it starts with "b." it's a batch token
	if len(id) > 2 && strings.HasPrefix(id, "b.") {
		return ts.lookupBatchToken(ctx, id)
	}

	lock := locksutil.LockForKey(ts.tokenLocks, id)
	lock.RLock()
	defer lock.RUnlock()

	return ts.lookupInternal(ctx, id, false, false)
}

// lookupTainted is used to find a token that may or may not be tainted given
// its ID. It acquires a read lock, then calls lookupInternal.
func (ts *TokenStore) lookupTainted(ctx context.Context, id string) (*logical.TokenEntry, error) {
	defer metrics.MeasureSince([]string{"token", "lookup"}, time.Now())
	if id == "" {
		return nil, fmt.Errorf("cannot lookup blank token")
	}

	lock := locksutil.LockForKey(ts.tokenLocks, id)
	lock.RLock()
	defer lock.RUnlock()

	return ts.lookupInternal(ctx, id, false, true)
}

func (ts *TokenStore) lookupBatchToken(ctx context.Context, id string) (*logical.TokenEntry, error) {
	// Strip the b. from the front and namespace ID from the back
	bEntry, _ := namespace.SplitIDFromString(id[2:])

	eEntry, err := base64.RawURLEncoding.DecodeString(bEntry)
	if err != nil {
		return nil, err
	}

	mEntry, err := ts.batchTokenEncryptor.Decrypt(ctx, "", eEntry)
	if err != nil {
		return nil, nil
	}

	pEntry := new(pb.TokenEntry)
	if err := proto.Unmarshal(mEntry, pEntry); err != nil {
		return nil, err
	}

	te, err := pb.ProtoTokenEntryToLogicalTokenEntry(pEntry)
	if err != nil {
		return nil, err
	}

	if time.Now().After(time.Unix(te.CreationTime, 0).Add(te.TTL)) {
		return nil, nil
	}

	if te.Parent != "" {
		pte, err := ts.Lookup(ctx, te.Parent)
		if err != nil {
			return nil, err
		}
		if pte == nil {
			return nil, nil
		}
	}

	te.ID = id
	return te, nil
}

// lookupInternal is used to find a token given its (possibly salted) ID. If
// tainted is true, entries that are in some revocation state (currently,
// indicated by num uses < 0), the entry will be returned anyways
func (ts *TokenStore) lookupInternal(ctx context.Context, id string, salted, tainted bool) (*logical.TokenEntry, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, errwrap.Wrapf("failed to find namespace in context: {{err}}", err)
	}

	// If it starts with "b." it's a batch token
	if len(id) > 2 && strings.HasPrefix(id, "b.") {
		return ts.lookupBatchToken(ctx, id)
	}

	var raw *logical.StorageEntry
	lookupID := id

	if !salted {
		// If possible, always use the token's namespace. If it doesn't match
		// the request namespace, ensure the request namespace is a child
		_, nsID := namespace.SplitIDFromString(id)
		if nsID != "" {
			tokenNS, err := NamespaceByID(ctx, nsID, ts.core)
			if err != nil {
				return nil, errwrap.Wrapf("failed to look up namespace from the token: {{err}}", err)
			}
			if tokenNS != nil {
				if tokenNS.ID != ns.ID {
					ns = tokenNS
					ctx = namespace.ContextWithNamespace(ctx, tokenNS)
				}
			}
		} else {
			// Any non-root-ns token should have an accessor and child
			// namespaces cannot have custom IDs. If someone omits or tampers
			// with it, the lookup in the root namespace simply won't work.
			ns = namespace.RootNamespace
			ctx = namespace.ContextWithNamespace(ctx, ns)
		}

		lookupID, err = ts.SaltID(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	raw, err = ts.idView(ns).Get(ctx, lookupID)
	if err != nil {
		return nil, errwrap.Wrapf("failed to read entry: {{err}}", err)
	}

	// Bail if not found
	if raw == nil {
		return nil, nil
	}

	// Unmarshal the token
	entry := new(logical.TokenEntry)
	if err := jsonutil.DecodeJSON(raw.Value, entry); err != nil {
		return nil, errwrap.Wrapf("failed to decode entry: {{err}}", err)
	}

	// This is a token that is awaiting deferred revocation or tainted
	if entry.NumUses < 0 && !tainted {
		return nil, nil
	}

	if entry.NamespaceID == "" {
		entry.NamespaceID = namespace.RootNamespaceID
	}

	// This will be the upgrade case
	if entry.Type == logical.TokenTypeDefault {
		entry.Type = logical.TokenTypeService
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

	// It's a root token with unlimited creation TTL (so never had an
	// expiration); this may or may not have a lease (based on when it was
	// generated, for later revocation purposes) but it doesn't matter, it's
	// allowed. Fast-path this.
	if len(entry.Policies) == 1 && entry.Policies[0] == "root" && entry.TTL == 0 {
		// If fields are getting upgraded, store the changes
		if persistNeeded {
			if err := ts.store(ctx, entry); err != nil {
				return nil, errwrap.Wrapf("failed to persist token upgrade: {{err}}", err)
			}
		}
		return entry, nil
	}

	// Perform these checks on upgraded fields, but before persisting

	// If we are still restoring the expiration manager, we want to ensure the
	// token is not expired
	if ts.expiration == nil {
		return nil, errors.New("expiration manager is nil on tokenstore")
	}
	le, err := ts.expiration.FetchLeaseTimesByToken(ctx, entry)
	if err != nil {
		return nil, errwrap.Wrapf("failed to fetch lease times: {{err}}", err)
	}

	var ret *logical.TokenEntry

	switch {
	// It's any kind of expiring token with no lease, immediately delete it
	case le == nil:
		tokenNS, err := NamespaceByID(ctx, entry.NamespaceID, ts.core)
		if err != nil {
			return nil, err
		}
		if tokenNS == nil {
			return nil, namespace.ErrNoNamespace
		}

		revokeCtx := namespace.ContextWithNamespace(ts.quitContext, tokenNS)
		leaseID, err := ts.expiration.CreateOrFetchRevocationLeaseByToken(revokeCtx, entry)
		if err != nil {
			return nil, err
		}

		err = ts.expiration.Revoke(revokeCtx, leaseID)
		if err != nil {
			return nil, err
		}

	// Only return if we're not past lease expiration (or if tainted is true),
	// otherwise assume expmgr is working on revocation
	default:
		if !le.ExpireTime.Before(time.Now()) || tainted {
			ret = entry
		}
	}

	// If fields are getting upgraded, store the changes
	if persistNeeded {
		if err := ts.store(ctx, entry); err != nil {
			return nil, errwrap.Wrapf("failed to persist token upgrade: {{err}}", err)
		}
	}

	return ret, nil
}

// Revoke is used to invalidate a given token, any child tokens
// will be orphaned.
func (ts *TokenStore) revokeOrphan(ctx context.Context, id string) error {
	defer metrics.MeasureSince([]string{"token", "revoke"}, time.Now())
	if id == "" {
		return fmt.Errorf("cannot revoke blank token")
	}

	saltedID, err := ts.SaltID(ctx, id)
	if err != nil {
		return err
	}

	return ts.revokeInternal(ctx, saltedID, false)
}

// revokeInternal is used to invalidate a given salted token, any child tokens
// will be orphaned unless otherwise specified. skipOrphan should be used
// whenever we are revoking the entire tree starting from a particular parent
// (e.g. revokeTreeInternal).
func (ts *TokenStore) revokeInternal(ctx context.Context, saltedID string, skipOrphan bool) (ret error) {
	// Check and set the token deletion state. We only proceed with the deletion
	// if we don't have a pending deletion (empty), or if the deletion previously
	// failed (state is false)
	state, loaded := ts.tokensPendingDeletion.LoadOrStore(saltedID, true)

	// If the entry was loaded and its state is true, we short-circuit
	if loaded && state == true {
		return nil
	}

	// The map check above should protect use from any concurrent revocations, so
	// we do another lookup here to make sure we have the right state
	entry, err := ts.lookupInternal(ctx, saltedID, true, true)
	if err != nil {
		return err
	}
	if entry == nil {
		return nil
	}

	if entry.NumUses != tokenRevocationPending {
		entry.NumUses = tokenRevocationPending
		if err := ts.store(ctx, entry); err != nil {
			// The only real reason for this is an underlying storage error
			// which also means that nothing else in this func or expmgr will
			// really work either. So we clear revocation state so the user can
			// try again.
			ts.logger.Error("failed to mark token as revoked")
			ts.tokensPendingDeletion.Store(entry.ID, false)
			return err
		}
	}

	tokenNS, err := NamespaceByID(ctx, entry.NamespaceID, ts.core)
	if err != nil {
		return err
	}
	if tokenNS == nil {
		return namespace.ErrNoNamespace
	}

	defer func() {
		// If we succeeded in all other revocation operations after this defer and
		// before we return, we can remove the token store entry
		if ret == nil {
			if err := ts.idView(tokenNS).Delete(ctx, saltedID); err != nil {
				ret = errwrap.Wrapf("failed to delete entry: {{err}}", err)
			}
		}

		// Check on ret again and update the sync.Map accordingly
		if ret != nil {
			// If we failed on any of the calls within, we store the state as false
			// so that the next call to revokeInternal will retry
			ts.tokensPendingDeletion.Store(saltedID, false)
		} else {
			ts.tokensPendingDeletion.Delete(saltedID)
		}
	}()

	// Destroy the token's cubby. This should go first as it's a
	// security-sensitive item.
	err = ts.cubbyholeDestroyer(ctx, ts, entry)
	if err != nil {
		return err
	}

	revokeCtx := namespace.ContextWithNamespace(ts.quitContext, tokenNS)
	if err := ts.expiration.RevokeByToken(revokeCtx, entry); err != nil {
		return err
	}

	// Clear the secondary index if any
	if entry.Parent != "" {
		_, parentNSID := namespace.SplitIDFromString(entry.Parent)
		parentCtx := revokeCtx
		parentNS := tokenNS

		if parentNSID != tokenNS.ID {
			switch {
			case parentNSID == "":
				parentNS = namespace.RootNamespace
			default:
				parentNS, err = NamespaceByID(ctx, parentNSID, ts.core)
				if err != nil {
					return errwrap.Wrapf("failed to get parent namespace: {{err}}", err)
				}
				if parentNS == nil {
					return namespace.ErrNoNamespace
				}
			}

			parentCtx = namespace.ContextWithNamespace(ctx, parentNS)
		}

		parentSaltedID, err := ts.SaltID(parentCtx, entry.Parent)
		if err != nil {
			return err
		}

		path := parentSaltedID + "/" + saltedID
		if tokenNS.ID != namespace.RootNamespaceID {
			path = fmt.Sprintf("%s.%s", path, tokenNS.ID)
		}

		if err = ts.parentView(parentNS).Delete(ctx, path); err != nil {
			return errwrap.Wrapf("failed to delete entry: {{err}}", err)
		}
	}

	// Clear the accessor index if any
	if entry.Accessor != "" {
		accessorSaltedID, err := ts.SaltID(revokeCtx, entry.Accessor)
		if err != nil {
			return err
		}

		if err = ts.accessorView(tokenNS).Delete(ctx, accessorSaltedID); err != nil {
			return errwrap.Wrapf("failed to delete entry: {{err}}", err)
		}
	}

	if !skipOrphan {
		// Mark all children token as orphan by removing
		// their parent index, and clear the parent entry.
		//
		// Marking the token as orphan should be skipped if it's called by
		// revokeTreeInternal to avoid unnecessary view.List operations. Since
		// the deletion occurs in a DFS fashion we don't need to perform a delete
		// on child prefixes as there will be none (as saltedID entry is a leaf node).
		children, err := ts.parentView(tokenNS).List(ctx, saltedID+"/")
		if err != nil {
			return errwrap.Wrapf("failed to scan for children: {{err}}", err)
		}
		for _, child := range children {
			var childNSID string
			childCtx := revokeCtx
			child, childNSID = namespace.SplitIDFromString(child)
			if childNSID != "" {
				childNS, err := NamespaceByID(ctx, childNSID, ts.core)
				if err != nil {
					return errwrap.Wrapf("failed to get child token: {{err}}", err)
				}
				if childNS == nil {
					return namespace.ErrNoNamespace
				}

				childCtx = namespace.ContextWithNamespace(ctx, childNS)
			}

			entry, err := ts.lookupInternal(childCtx, child, true, true)
			if err != nil {
				return errwrap.Wrapf("failed to get child token: {{err}}", err)
			}
			if entry == nil {
				// Seems it's already revoked, so nothing to do here except delete the index
				err = ts.parentView(tokenNS).Delete(ctx, child)
				if err != nil {
					return errwrap.Wrapf("failed to delete child entry: {{err}}", err)
				}
				continue
			}

			lock := locksutil.LockForKey(ts.tokenLocks, entry.ID)
			lock.Lock()

			entry.Parent = ""
			err = ts.store(childCtx, entry)
			if err != nil {
				lock.Unlock()
				return errwrap.Wrapf("failed to update child token: {{err}}", err)
			}
			lock.Unlock()

			// Delete the the child storage entry after we update the token entry Since
			// paths are not deeply nested (i.e. they are simply
			// parenPrefix/<parentID>/<childID>), we can simply call view.Delete instead
			// of logical.ClearView
			err = ts.parentView(tokenNS).Delete(ctx, child)
			if err != nil {
				return errwrap.Wrapf("failed to delete child entry: {{err}}", err)
			}
		}
	}

	return nil
}

// revokeTree is used to invalidate a given token and all
// child tokens.
func (ts *TokenStore) revokeTree(ctx context.Context, le *leaseEntry) error {
	defer metrics.MeasureSince([]string{"token", "revoke-tree"}, time.Now())
	// Verify the token is not blank
	if le.ClientToken == "" {
		return fmt.Errorf("cannot tree-revoke blank token")
	}

	// In case lookup fails for some reason for the token itself, set the
	// context for the next call from the lease entry's NS. This function is
	// only called when a lease for a given token is expiring, so it should run
	// in the context of the token namespace
	revCtx := namespace.ContextWithNamespace(ctx, le.namespace)

	saltedID, err := ts.SaltID(revCtx, le.ClientToken)
	if err != nil {
		return err
	}

	// Nuke the entire tree recursively
	return ts.revokeTreeInternal(revCtx, saltedID)
}

// revokeTreeInternal is used to invalidate a given token and all
// child tokens.
// Updated to be non-recursive and revoke child tokens
// before parent tokens(DFS).
func (ts *TokenStore) revokeTreeInternal(ctx context.Context, id string) error {
	dfs := []string{id}
	seenIDs := make(map[string]struct{})

	var ns *namespace.Namespace

	te, err := ts.lookupInternal(ctx, id, true, true)
	if err != nil {
		return err
	}
	if te == nil {
		ns, err = namespace.FromContext(ctx)
		if err != nil {
			return err
		}
	} else {
		ns, err = NamespaceByID(ctx, te.NamespaceID, ts.core)
		if err != nil {
			return err
		}
	}
	if ns == nil {
		return fmt.Errorf("failed to find namespace for token revocation")
	}

	for l := len(dfs); l > 0; l = len(dfs) {
		id := dfs[len(dfs)-1]
		seenIDs[id] = struct{}{}

		saltedCtx := ctx
		saltedNS := ns
		saltedID, saltedNSID := namespace.SplitIDFromString(id)
		if saltedNSID != "" {
			saltedNS, err = NamespaceByID(ctx, saltedNSID, ts.core)
			if err != nil {
				return errwrap.Wrapf("failed to find namespace for token revocation: {{err}}", err)
			}

			saltedCtx = namespace.ContextWithNamespace(ctx, saltedNS)
		}

		path := saltedID + "/"
		childrenRaw, err := ts.parentView(saltedNS).List(saltedCtx, path)
		if err != nil {
			return errwrap.Wrapf("failed to scan for children: {{err}}", err)
		}

		// Filter the child list to remove any items that have ever been in the dfs stack.
		// This is a robustness check, as a parent/child cycle can lead to an OOM crash.
		children := make([]string, 0, len(childrenRaw))
		for _, child := range childrenRaw {
			if _, seen := seenIDs[child]; !seen {
				children = append(children, child)
			} else {
				if err = ts.parentView(saltedNS).Delete(saltedCtx, path+child); err != nil {
					return errwrap.Wrapf("failed to delete entry: {{err}}", err)
				}

				ts.Logger().Warn("token cycle found", "token", child)
			}
		}

		// If the length of the children array is zero,
		// then we are at a leaf node.
		if len(children) == 0 {
			// Whenever revokeInternal is called, the token will be removed immediately and
			// any underlying secrets will be handed off to the expiration manager which will
			// take care of expiring them. If Vault is restarted, any revoked tokens
			// would have been deleted, and any pending leases for deletion will be restored
			// by the expiration manager.
			if err := ts.revokeInternal(saltedCtx, saltedID, true); err != nil {
				return errwrap.Wrapf("failed to revoke entry: {{err}}", err)
			}
			// If the length of l is equal to 1, then the last token has been deleted
			if l == 1 {
				return nil
			}
			dfs = dfs[:len(dfs)-1]
		} else {
			// If we make it here, there are children and they must be appended.
			dfs = append(dfs, children...)
		}
	}

	return nil
}

func (c *Core) IsBatchTokenCreationRequest(ctx context.Context, path string) (bool, error) {
	name := strings.TrimPrefix(path, "auth/token/create/")
	roleEntry, err := c.tokenStore.tokenStoreRole(ctx, name)
	if err != nil {
		return false, err
	}
	return roleEntry.TokenType == logical.TokenTypeBatch, nil
}

// handleCreateAgainstRole handles the auth/token/create path for a role
func (ts *TokenStore) handleCreateAgainstRole(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("role_name").(string)
	roleEntry, err := ts.tokenStoreRole(ctx, name)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown role %s", name)), nil
	}

	return ts.handleCreateCommon(ctx, req, d, false, roleEntry)
}

func (ts *TokenStore) lookupByAccessor(ctx context.Context, id string, salted, tainted bool) (accessorEntry, error) {
	var aEntry accessorEntry

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return aEntry, err
	}

	lookupID := id
	if !salted {
		_, nsID := namespace.SplitIDFromString(id)
		if nsID != "" {
			accessorNS, err := NamespaceByID(ctx, nsID, ts.core)
			if err != nil {
				return aEntry, err
			}
			if accessorNS != nil {
				if accessorNS.ID != ns.ID {
					ns = accessorNS
					ctx = namespace.ContextWithNamespace(ctx, accessorNS)
				}
			}
		} else {
			// Any non-root-ns token should have an accessor and child
			// namespaces cannot have custom IDs. If someone omits or tampers
			// with it, the lookup in the root namespace simply won't work.
			ns = namespace.RootNamespace
			ctx = namespace.ContextWithNamespace(ctx, ns)
		}

		lookupID, err = ts.SaltID(ctx, id)
		if err != nil {
			return aEntry, err
		}
	}

	entry, err := ts.accessorView(ns).Get(ctx, lookupID)

	if err != nil {
		return aEntry, errwrap.Wrapf("failed to read index using accessor: {{err}}", err)
	}
	if entry == nil {
		return aEntry, &logical.StatusBadRequest{Err: "invalid accessor"}
	}

	err = jsonutil.DecodeJSON(entry.Value, &aEntry)
	// If we hit an error, assume it's a pre-struct straight token ID
	if err != nil {
		te, err := ts.lookupInternal(ctx, string(entry.Value), false, tainted)
		if err != nil {
			return accessorEntry{}, errwrap.Wrapf("failed to look up token using accessor index: {{err}}", err)
		}
		// It's hard to reason about what to do here if te is nil -- it may be
		// that the token was revoked async, or that it's an old accessor index
		// entry that was somehow not cleared up, or or or. A nonexistent token
		// entry on lookup is nil, not an error, so we keep that behavior here
		// to be safe...the token ID is simply not filled in.
		if te != nil {
			aEntry.TokenID = te.ID
			aEntry.AccessorID = te.Accessor
			aEntry.NamespaceID = te.NamespaceID
		}
	}

	if aEntry.NamespaceID == "" {
		aEntry.NamespaceID = namespace.RootNamespaceID
	}

	return aEntry, nil
}

// handleTidy handles the cleaning up of leaked accessor storage entries and
// cleaning up of leases that are associated to tokens that are expired.
func (ts *TokenStore) handleTidy(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if !atomic.CompareAndSwapUint32(ts.tidyLock, 0, 1) {
		resp := &logical.Response{}
		resp.AddWarning("Tidy operation already in progress.")
		return resp, nil
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, errwrap.Wrapf("failed to get namespace from context: {{err}}", err)
	}

	go func() {
		defer atomic.StoreUint32(ts.tidyLock, 0)

		logger := ts.logger.Named("tidy")

		var tidyErrors *multierror.Error

		doTidy := func() error {

			ts.logger.Info("beginning tidy operation on tokens")
			defer ts.logger.Info("finished tidy operation on tokens")

			quitCtx := namespace.ContextWithNamespace(ts.quitContext, ns)

			// List out all the accessors
			saltedAccessorList, err := ts.accessorView(ns).List(quitCtx, "")
			if err != nil {
				return errwrap.Wrapf("failed to fetch accessor index entries: {{err}}", err)
			}

			// First, clean up secondary index entries that are no longer valid
			parentList, err := ts.parentView(ns).List(quitCtx, "")
			if err != nil {
				return errwrap.Wrapf("failed to fetch secondary index entries: {{err}}", err)
			}

			// List all the cubbyhole storage keys
			cubbyholeKeys, err := ts.cubbyholeBackend.storageView.List(quitCtx, "")
			if err != nil {
				return errwrap.Wrapf("failed to fetch cubbyhole storage keys: {{err}}", err)
			}

			var countParentEntries, deletedCountParentEntries, countParentList, deletedCountParentList int64

			// Scan through the secondary index entries; if there is an entry
			// with the token's salt ID at the end, remove it
			for _, parent := range parentList {
				countParentEntries++

				// Get the children
				children, err := ts.parentView(ns).List(quitCtx, parent)
				if err != nil {
					tidyErrors = multierror.Append(tidyErrors, errwrap.Wrapf("failed to read secondary index: {{err}}", err))
					continue
				}

				// First check if the salt ID of the parent exists, and if not mark this so
				// that deletion of children later with this loop below applies to all
				// children
				originalChildrenCount := int64(len(children))
				exists, _ := ts.lookupInternal(quitCtx, strings.TrimSuffix(parent, "/"), true, true)
				if exists == nil {
					ts.logger.Debug("deleting invalid parent prefix entry", "index", parentPrefix+parent)
				}

				var deletedChildrenCount int64
				for _, child := range children {
					countParentList++
					if countParentList%500 == 0 {
						ts.logger.Info("checking validity of tokens in secondary index list", "progress", countParentList)
					}

					// Look up tainted entries so we can be sure that if this isn't
					// found, it doesn't exist. Doing the following without locking
					// since appropriate locks cannot be held with salted token IDs.
					// Also perform deletion if the parent doesn't exist any more.
					te, _ := ts.lookupInternal(quitCtx, child, true, true)
					// If the child entry is not nil, but the parent doesn't exist, then turn
					// that child token into an orphan token. Theres no deletion in this case.
					if te != nil && exists == nil {
						lock := locksutil.LockForKey(ts.tokenLocks, te.ID)
						lock.Lock()

						te.Parent = ""
						err = ts.store(quitCtx, te)
						if err != nil {
							tidyErrors = multierror.Append(tidyErrors, errwrap.Wrapf("failed to convert child token into an orphan token: {{err}}", err))
						}
						lock.Unlock()
						continue
					}
					// Otherwise, if the entry doesn't exist, or if the parent doesn't exist go
					// on with the delete on the secondary index
					if te == nil || exists == nil {
						index := parent + child
						ts.logger.Debug("deleting invalid secondary index", "index", index)
						err = ts.parentView(ns).Delete(quitCtx, index)
						if err != nil {
							tidyErrors = multierror.Append(tidyErrors, errwrap.Wrapf("failed to delete secondary index: {{err}}", err))
							continue
						}
						deletedChildrenCount++
					}
				}
				// Add current children deleted count to the total count
				deletedCountParentList += deletedChildrenCount
				// N.B.: We don't call delete on the parent prefix since physical.Backend.Delete
				// implementations should be in charge of deleting empty prefixes.
				// If we deleted all the children, then add that to our deleted parent entries count.
				if originalChildrenCount == deletedChildrenCount {
					deletedCountParentEntries++
				}
			}

			var countAccessorList,
				countCubbyholeKeys,
				deletedCountAccessorEmptyToken,
				deletedCountAccessorInvalidToken,
				deletedCountInvalidTokenInAccessor,
				deletedCountInvalidCubbyholeKey int64

			validCubbyholeKeys := make(map[string]bool)

			// For each of the accessor, see if the token ID associated with it is
			// a valid one. If not, delete the leases associated with that token
			// and delete the accessor as well.
			for _, saltedAccessor := range saltedAccessorList {
				countAccessorList++
				if countAccessorList%500 == 0 {
					ts.logger.Info("checking if accessors contain valid tokens", "progress", countAccessorList)
				}

				accessorEntry, err := ts.lookupByAccessor(quitCtx, saltedAccessor, true, true)
				if err != nil {
					tidyErrors = multierror.Append(tidyErrors, errwrap.Wrapf("failed to read the accessor index: {{err}}", err))
					continue
				}

				// A valid accessor storage entry should always have a token ID
				// in it. If not, it is an invalid accessor entry and needs to
				// be deleted.
				if accessorEntry.TokenID == "" {
					// If deletion of accessor fails, move on to the next
					// item since this is just a best-effort operation
					err = ts.accessorView(ns).Delete(quitCtx, saltedAccessor)
					if err != nil {
						tidyErrors = multierror.Append(tidyErrors, errwrap.Wrapf("failed to delete the accessor index: {{err}}", err))
						continue
					}
					deletedCountAccessorEmptyToken++
				}

				lock := locksutil.LockForKey(ts.tokenLocks, accessorEntry.TokenID)
				lock.RLock()

				// Look up tainted variants so we only find entries that truly don't
				// exist
				te, err := ts.lookupInternal(quitCtx, accessorEntry.TokenID, false, true)
				if err != nil {
					tidyErrors = multierror.Append(tidyErrors, errwrap.Wrapf("failed to lookup tainted ID: {{err}}", err))
					lock.RUnlock()
					continue
				}

				lock.RUnlock()

				switch {
				case te == nil:
					// If token entry is not found assume that the token is not valid any
					// more and conclude that accessor, leases, and secondary index entries
					// for this token should not exist as well.

					ts.logger.Info("deleting token with nil entry referenced by accessor", "salted_accessor", saltedAccessor)

					// RevokeByToken expects a '*logical.TokenEntry'. For the
					// purposes of tidying, it is sufficient if the token
					// entry only has ID set.
					tokenEntry := &logical.TokenEntry{
						ID:          accessorEntry.TokenID,
						NamespaceID: accessorEntry.NamespaceID,
					}

					// Attempt to revoke the token. This will also revoke
					// the leases associated with the token.
					err = ts.expiration.RevokeByToken(quitCtx, tokenEntry)
					if err != nil {
						tidyErrors = multierror.Append(tidyErrors, errwrap.Wrapf("failed to revoke leases of expired token: {{err}}", err))
						continue
					}
					deletedCountInvalidTokenInAccessor++

					// If deletion of accessor fails, move on to the next item since
					// this is just a best-effort operation. We do this last so that on
					// next run if something above failed we still have the accessor
					// entry to try again.
					err = ts.accessorView(ns).Delete(quitCtx, saltedAccessor)
					if err != nil {
						tidyErrors = multierror.Append(tidyErrors, errwrap.Wrapf("failed to delete accessor entry: {{err}}", err))
						continue
					}
					deletedCountAccessorInvalidToken++
				default:
					// Cache the cubbyhole storage key when the token is valid
					switch {
					case te.NamespaceID == namespace.RootNamespaceID && !strings.HasPrefix(te.ID, "s."):
						saltedID, err := ts.SaltID(quitCtx, te.ID)
						if err != nil {
							tidyErrors = multierror.Append(tidyErrors, errwrap.Wrapf("failed to create salted token id: {{err}}", err))
							continue
						}
						validCubbyholeKeys[salt.SaltID(ts.cubbyholeBackend.saltUUID, saltedID, salt.SHA1Hash)] = true
					default:
						if te.CubbyholeID == "" {
							tidyErrors = multierror.Append(tidyErrors, fmt.Errorf("missing cubbyhole ID for a valid token"))
							continue
						}
						validCubbyholeKeys[te.CubbyholeID] = true
					}
				}
			}

			// Revoke invalid cubbyhole storage keys
			for _, key := range cubbyholeKeys {
				countCubbyholeKeys++
				if countCubbyholeKeys%500 == 0 {
					ts.logger.Info("checking if there are invalid cubbyholes", "progress", countCubbyholeKeys)
				}

				key := strings.TrimSuffix(key, "/")
				if !validCubbyholeKeys[key] {
					ts.logger.Info("deleting invalid cubbyhole", "key", key)
					err = ts.cubbyholeBackend.revoke(quitCtx, key)
					if err != nil {
						tidyErrors = multierror.Append(tidyErrors, errwrap.Wrapf(fmt.Sprintf("failed to revoke cubbyhole key %q: {{err}}", key), err))
					}
					deletedCountInvalidCubbyholeKey++
				}
			}

			ts.logger.Info("number of entries scanned in parent prefix", "count", countParentEntries)
			ts.logger.Info("number of entries deleted in parent prefix", "count", deletedCountParentEntries)
			ts.logger.Info("number of tokens scanned in parent index list", "count", countParentList)
			ts.logger.Info("number of tokens revoked in parent index list", "count", deletedCountParentList)
			ts.logger.Info("number of accessors scanned", "count", countAccessorList)
			ts.logger.Info("number of deleted accessors which had empty tokens", "count", deletedCountAccessorEmptyToken)
			ts.logger.Info("number of revoked tokens which were invalid but present in accessors", "count", deletedCountInvalidTokenInAccessor)
			ts.logger.Info("number of deleted accessors which had invalid tokens", "count", deletedCountAccessorInvalidToken)
			ts.logger.Info("number of deleted cubbyhole keys that were invalid", "count", deletedCountInvalidCubbyholeKey)

			return tidyErrors.ErrorOrNil()
		}

		if err := doTidy(); err != nil {
			logger.Error("error running tidy", "error", err)
			return
		}
	}()

	resp := &logical.Response{}
	resp.AddWarning("Tidy operation successfully started. Any information from the operation will be printed to Vault's server logs.")
	return logical.RespondWithStatusCode(resp, req, http.StatusAccepted)
}

// handleUpdateLookupAccessor handles the auth/token/lookup-accessor path for returning
// the properties of the token associated with the accessor
func (ts *TokenStore) handleUpdateLookupAccessor(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	accessor := data.Get("accessor").(string)
	if accessor == "" {
		return nil, &logical.StatusBadRequest{Err: "missing accessor"}
	}

	aEntry, err := ts.lookupByAccessor(ctx, accessor, false, false)
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
	resp, err := ts.handleLookup(ctx, req, d)
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

	return resp, nil
}

// handleUpdateRevokeAccessor handles the auth/token/revoke-accessor path for revoking
// the token associated with the accessor
func (ts *TokenStore) handleUpdateRevokeAccessor(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	accessor := data.Get("accessor").(string)
	if accessor == "" {
		return nil, &logical.StatusBadRequest{Err: "missing accessor"}
	}

	aEntry, err := ts.lookupByAccessor(ctx, accessor, false, true)
	if err != nil {
		return nil, err
	}

	te, err := ts.Lookup(ctx, aEntry.TokenID)
	if err != nil {
		return nil, err
	}
	if te == nil {
		return logical.ErrorResponse("token not found"), logical.ErrInvalidRequest
	}

	tokenNS, err := NamespaceByID(ctx, te.NamespaceID, ts.core)
	if err != nil {
		return nil, err
	}
	if tokenNS == nil {
		return nil, namespace.ErrNoNamespace
	}

	revokeCtx := namespace.ContextWithNamespace(ts.quitContext, tokenNS)
	leaseID, err := ts.expiration.CreateOrFetchRevocationLeaseByToken(revokeCtx, te)
	if err != nil {
		return nil, err
	}

	err = ts.expiration.Revoke(revokeCtx, leaseID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// handleCreate handles the auth/token/create path for creation of new orphan
// tokens
func (ts *TokenStore) handleCreateOrphan(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return ts.handleCreateCommon(ctx, req, d, true, nil)
}

// handleCreate handles the auth/token/create path for creation of new non-orphan
// tokens
func (ts *TokenStore) handleCreate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return ts.handleCreateCommon(ctx, req, d, false, nil)
}

// handleCreateCommon handles the auth/token/create path for creation of new tokens
func (ts *TokenStore) handleCreateCommon(ctx context.Context, req *logical.Request, d *framework.FieldData, orphan bool, role *tsRoleEntry) (*logical.Response, error) {
	// Read the parent policy
	parent, err := ts.Lookup(ctx, req.ClientToken)
	if err != nil {
		return nil, errwrap.Wrapf("parent token lookup failed: {{err}}", err)
	}
	if parent == nil {
		return logical.ErrorResponse("parent token lookup failed: no parent found"), logical.ErrInvalidRequest
	}
	if parent.Type == logical.TokenTypeBatch {
		return logical.ErrorResponse("batch tokens cannot create more tokens"), nil
	}

	// A token with a restricted number of uses cannot create a new token
	// otherwise it could escape the restriction count.
	if parent.NumUses > 0 {
		return logical.ErrorResponse("restricted use token cannot generate child tokens"),
			logical.ErrInvalidRequest
	}

	// Check if the client token has sudo/root privileges for the requested path
	isSudo := ts.System().SudoPrivilege(ctx, req.MountPoint+req.Path, req.ClientToken)

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
		Type            string `mapstructure:"type"`
	}
	if err := mapstructure.WeakDecode(req.Data, &data); err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error decoding request: %s", err)), logical.ErrInvalidRequest
	}

	// If the context's namespace is different from the parent and this is an
	// orphan token creation request, then this is an admin token generation for
	// the namespace
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if ns.ID != parent.NamespaceID {
		parentNS, err := NamespaceByID(ctx, parent.NamespaceID, ts.core)
		if err != nil {
			ts.logger.Error("error looking up parent namespace", "error", err, "parent_namespace", parent.NamespaceID)
			return nil, ErrInternalError
		}
		if parentNS == nil {
			ts.logger.Error("could not find information for parent namespace", "parent_namespace", parent.NamespaceID)
			return nil, ErrInternalError
		}

		if !isSudo {
			return logical.ErrorResponse("root or sudo privileges required to directly generate a token in a child namespace"), logical.ErrInvalidRequest
		}

		if strutil.StrListContains(data.Policies, "root") {
			return logical.ErrorResponse("root tokens may not be created from a parent namespace"), logical.ErrInvalidRequest
		}
	}

	renewable := true
	if data.Renewable != nil {
		renewable = *data.Renewable
	}

	tokenType := logical.TokenTypeService
	tokenTypeStr := data.Type
	if role != nil {
		switch role.TokenType {
		case logical.TokenTypeDefault, logical.TokenTypeDefaultService:
			// Use the user-given value, but fall back to service
		case logical.TokenTypeDefaultBatch:
			// Use the user-given value, but fall back to batch
			if tokenTypeStr == "" {
				tokenTypeStr = logical.TokenTypeBatch.String()
			}
		case logical.TokenTypeService:
			tokenTypeStr = logical.TokenTypeService.String()
		case logical.TokenTypeBatch:
			tokenTypeStr = logical.TokenTypeBatch.String()
		default:
			return logical.ErrorResponse(fmt.Sprintf("role being used for token creation contains invalid token type %q", role.TokenType.String())), nil
		}
	}
	switch tokenTypeStr {
	case "", "service":
	case "batch":
		var badReason string
		switch {
		case data.ExplicitMaxTTL != "":
			dur, err := parseutil.ParseDurationSecond(data.ExplicitMaxTTL)
			if err != nil {
				return logical.ErrorResponse(`"explicit_max_ttl" value could not be parsed`), nil
			}
			if dur != 0 {
				badReason = "explicit_max_ttl"
			}
		case data.NumUses != 0:
			badReason = "num_uses"
		case data.Period != "":
			dur, err := parseutil.ParseDurationSecond(data.Period)
			if err != nil {
				return logical.ErrorResponse(`"period" value could not be parsed`), nil
			}
			if dur != 0 {
				badReason = "period"
			}
		}
		if badReason != "" {
			return logical.ErrorResponse(fmt.Sprintf("batch tokens cannot have %q set", badReason)), nil
		}
		tokenType = logical.TokenTypeBatch
		renewable = false
	default:
		return logical.ErrorResponse("invalid 'token_type' value"), logical.ErrInvalidRequest
	}

	// Verify the number of uses is positive
	if data.NumUses < 0 {
		return logical.ErrorResponse("number of uses cannot be negative"),
			logical.ErrInvalidRequest
	}

	// Setup the token entry
	te := logical.TokenEntry{
		Parent: req.ClientToken,

		// The mount point is always the same since we have only one token
		// store; using req.MountPoint causes trouble in tests since they don't
		// have an official mount
		Path: fmt.Sprintf("auth/token/%s", req.Path),

		Meta:         data.Metadata,
		DisplayName:  "token",
		NumUses:      data.NumUses,
		CreationTime: time.Now().Unix(),
		NamespaceID:  ns.ID,
		Type:         tokenType,
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
		if ns.ID != namespace.RootNamespaceID {
			return logical.ErrorResponse("token IDs can only be manually specified in the root namespace"),
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
					return logical.ErrorResponse(fmt.Sprintf("token policies (%q) must be subset of the role's allowed policies (%q)", finalPolicies, sanitizedRolePolicies)), logical.ErrInvalidRequest
				}
			}
		} else {
			// Assign parent policies if none have been requested. As this is a
			// role, add default unless explicitly disabled.
			if len(finalPolicies) == 0 {
				finalPolicies = policyutil.SanitizePolicies(parent.Policies, localAddDefault)
			}
		}

		if len(role.DisallowedPolicies) > 0 {
			// We don't add the default here because we only want to disallow it if it's explicitly set
			sanitizedRolePolicies = strutil.RemoveDuplicates(role.DisallowedPolicies, true)

			for _, finalPolicy := range finalPolicies {
				if strutil.StrListContains(sanitizedRolePolicies, finalPolicy) {
					return logical.ErrorResponse(fmt.Sprintf("token policy %q is disallowed by this role", finalPolicy)), logical.ErrInvalidRequest
				}
			}
		}

		data.Policies = finalPolicies

	// We are creating a token from a parent namespace. We should only use the input
	// policies.
	case ns.ID != parent.NamespaceID:
		addDefault = !data.NoDefaultPolicy

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

	if strutil.StrListContains(te.Policies, "root") {
		// Prevent attempts to create a root token without an actual root token as parent.
		// This is to thwart privilege escalation by tokens having 'sudo' privileges.
		if !strutil.StrListContains(parent.Policies, "root") {
			return logical.ErrorResponse("root tokens may not be created without parent token being root"), logical.ErrInvalidRequest
		}

		if te.Type == logical.TokenTypeBatch {
			// Batch tokens cannot be revoked so we should never have root batch tokens
			return logical.ErrorResponse("batch tokens cannot be root tokens"), nil
		}
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

		if len(role.BoundCIDRs) > 0 {
			te.BoundCIDRs = role.BoundCIDRs
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

	// At this point, it is clear whether the token is going to be an orphan or
	// not. If the token is not going to be an orphan, inherit the parent's
	// entity identifier into the child token.
	if te.Parent != "" {
		te.EntityID = parent.EntityID

		// If the parent has bound CIDRs, copy those into the child. We don't
		// do this if role is not nil because then we always use the role's
		// bound CIDRs; roles allow escalation of privilege in proper
		// circumstances.
		if role == nil {
			te.BoundCIDRs = parent.BoundCIDRs
		}
	}

	var explicitMaxTTLToUse time.Duration
	if data.ExplicitMaxTTL != "" {
		dur, err := parseutil.ParseDurationSecond(data.ExplicitMaxTTL)
		if err != nil {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
		if dur < 0 {
			return logical.ErrorResponse("explicit_max_ttl must be positive"), logical.ErrInvalidRequest
		}
		te.ExplicitMaxTTL = dur
		explicitMaxTTLToUse = dur
	}

	var periodToUse time.Duration
	if data.Period != "" {
		dur, err := parseutil.ParseDurationSecond(data.Period)
		if err != nil {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}

		switch {
		case dur < 0:
			return logical.ErrorResponse("period must be positive"), logical.ErrInvalidRequest
		case dur == 0:
		default:
			if !isSudo {
				return logical.ErrorResponse("root or sudo privileges required to create periodic token"),
					logical.ErrInvalidRequest
			}
			te.Period = dur
			periodToUse = dur
		}
	}

	// Parse the TTL/lease if any
	if data.TTL != "" {
		dur, err := parseutil.ParseDurationSecond(data.TTL)
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

	// Set the lesser period/explicit max TTL if defined both in arguments and
	// in role. Batch tokens will error out if not set via role, but here we
	// need to explicitly check
	if role != nil && te.Type != logical.TokenTypeBatch {
		if role.ExplicitMaxTTL != 0 {
			switch {
			case explicitMaxTTLToUse == 0:
				explicitMaxTTLToUse = role.ExplicitMaxTTL
			default:
				if role.ExplicitMaxTTL < explicitMaxTTLToUse {
					explicitMaxTTLToUse = role.ExplicitMaxTTL
				}
				resp.AddWarning(fmt.Sprintf("Explicit max TTL specified both during creation call and in role; using the lesser value of %d seconds", int64(explicitMaxTTLToUse.Seconds())))
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

	// Only calculate a TTL if you are A) periodic, B) have a TTL, C) do not have a TTL and are not a root token
	if periodToUse > 0 || te.TTL > 0 || (te.TTL == 0 && !strutil.StrListContains(te.Policies, "root")) {
		ttl, warnings, err := framework.CalculateTTL(sysView, 0, te.TTL, periodToUse, 0, explicitMaxTTLToUse, time.Unix(te.CreationTime, 0))
		if err != nil {
			return nil, err
		}
		for _, warning := range warnings {
			resp.AddWarning(warning)
		}
		te.TTL = ttl
	}

	// Root tokens are still bound by explicit max TTL
	if te.TTL == 0 && explicitMaxTTLToUse > 0 {
		te.TTL = explicitMaxTTLToUse
	}

	// Don't advertise non-expiring root tokens as renewable, as attempts to
	// renew them are denied. Don't CIDR-restrict these either.
	if te.TTL == 0 {
		if parent.TTL != 0 {
			return logical.ErrorResponse("expiring root tokens cannot create non-expiring root tokens"), logical.ErrInvalidRequest
		}
		renewable = false
		te.BoundCIDRs = nil
	}

	if te.ID != "" {
		resp.AddWarning("Supplying a custom ID for the token uses the weaker SHA1 hashing instead of the more secure SHA2-256 HMAC for token obfuscation. SHA1 hashed tokens on the wire leads to less secure lookups.")
	}

	// Create the token
	if err := ts.create(ctx, &te); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Generate the response
	resp.Auth = &logical.Auth{
		NumUses:     te.NumUses,
		DisplayName: te.DisplayName,
		Policies:    te.Policies,
		Metadata:    te.Meta,
		LeaseOptions: logical.LeaseOptions{
			TTL:       te.TTL,
			Renewable: renewable,
		},
		ClientToken:    te.ID,
		Accessor:       te.Accessor,
		EntityID:       te.EntityID,
		Period:         periodToUse,
		ExplicitMaxTTL: explicitMaxTTLToUse,
		CreationPath:   te.Path,
		TokenType:      te.Type,
		Orphan:         te.Parent == "",
	}

	for _, p := range te.Policies {
		policy, err := ts.core.policyStore.GetPolicy(ctx, p, PolicyTypeToken)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("could not look up policy %s", p)), nil
		}
		if policy == nil {
			resp.AddWarning(fmt.Sprintf("Policy %q does not exist", p))
		}
	}

	return resp, nil
}

// handleRevokeSelf handles the auth/token/revoke-self path for revocation of tokens
// in a way that revokes all child tokens. Normally, using sys/revoke/leaseID will revoke
// the token and all children anyways, but that is only available when there is a lease.
func (ts *TokenStore) handleRevokeSelf(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return ts.revokeCommon(ctx, req, data, req.ClientToken)
}

// handleRevokeTree handles the auth/token/revoke/id path for revocation of tokens
// in a way that revokes all child tokens. Normally, using sys/revoke/leaseID will revoke
// the token and all children anyways, but that is only available when there is a lease.
func (ts *TokenStore) handleRevokeTree(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	id := data.Get("token").(string)
	if id == "" {
		return logical.ErrorResponse("missing token ID"), logical.ErrInvalidRequest
	}

	if resp, err := ts.revokeCommon(ctx, req, data, id); resp != nil || err != nil {
		return resp, err
	}

	return nil, nil
}

func (ts *TokenStore) revokeCommon(ctx context.Context, req *logical.Request, data *framework.FieldData, id string) (*logical.Response, error) {
	te, err := ts.Lookup(ctx, id)
	if err != nil {
		return nil, err
	}
	if te == nil {
		return nil, nil
	}

	if te.Type == logical.TokenTypeBatch {
		return logical.ErrorResponse("batch tokens cannot be revoked"), nil
	}

	tokenNS, err := NamespaceByID(ctx, te.NamespaceID, ts.core)
	if err != nil {
		return nil, err
	}
	if tokenNS == nil {
		return nil, namespace.ErrNoNamespace
	}

	revokeCtx := namespace.ContextWithNamespace(ts.quitContext, tokenNS)
	leaseID, err := ts.expiration.CreateOrFetchRevocationLeaseByToken(revokeCtx, te)
	if err != nil {
		return nil, err
	}

	err = ts.expiration.Revoke(revokeCtx, leaseID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// handleRevokeOrphan handles the auth/token/revoke-orphan/id path for revocation of tokens
// in a way that leaves child tokens orphaned. Normally, using sys/revoke/leaseID will revoke
// the token and all children.
func (ts *TokenStore) handleRevokeOrphan(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Parse the id
	id := data.Get("token").(string)
	if id == "" {
		return logical.ErrorResponse("missing token ID"), logical.ErrInvalidRequest
	}

	// Check if the client token has sudo/root privileges for the requested path
	isSudo := ts.System().SudoPrivilege(ctx, req.MountPoint+req.Path, req.ClientToken)

	if !isSudo {
		return logical.ErrorResponse("root or sudo privileges required to revoke and orphan"),
			logical.ErrInvalidRequest
	}

	// Do a lookup. Among other things, that will ensure that this is either
	// running in the same namespace or a parent.
	te, err := ts.Lookup(ctx, id)
	if err != nil {
		return nil, errwrap.Wrapf("error when looking up token to revoke: {{err}}", err)
	}
	if te == nil {
		return logical.ErrorResponse("token to revoke not found"), logical.ErrInvalidRequest
	}

	if te.Type == logical.TokenTypeBatch {
		return logical.ErrorResponse("batch tokens cannot be revoked"), nil
	}

	// Revoke and orphan
	if err := ts.revokeOrphan(ctx, id); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	return nil, nil
}

func (ts *TokenStore) handleLookupSelf(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	data.Raw["token"] = req.ClientToken
	return ts.handleLookup(ctx, req, data)
}

// handleLookup handles the auth/token/lookup/id path for querying information about
// a particular token. This can be used to see which policies are applicable.
func (ts *TokenStore) handleLookup(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	id := data.Get("token").(string)
	if id == "" {
		id = req.ClientToken
	}
	if id == "" {
		return logical.ErrorResponse("missing token ID"), logical.ErrInvalidRequest
	}

	lock := locksutil.LockForKey(ts.tokenLocks, id)
	lock.RLock()
	defer lock.RUnlock()

	out, err := ts.lookupInternal(ctx, id, false, true)
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
			"expire_time":      nil,
			"ttl":              int64(0),
			"explicit_max_ttl": int64(out.ExplicitMaxTTL.Seconds()),
			"entity_id":        out.EntityID,
			"type":             out.Type.String(),
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

	if len(out.BoundCIDRs) > 0 {
		resp.Data["bound_cidrs"] = out.BoundCIDRs
	}

	tokenNS, err := NamespaceByID(ctx, out.NamespaceID, ts.core)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	if tokenNS == nil {
		return nil, namespace.ErrNoNamespace
	}

	if out.NamespaceID != namespace.RootNamespaceID {
		resp.Data["namespace_path"] = tokenNS.Path
	}

	// Fetch the last renewal time
	leaseTimes, err := ts.expiration.FetchLeaseTimesByToken(ctx, out)
	if err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	if leaseTimes != nil {
		if !leaseTimes.LastRenewalTime.IsZero() {
			resp.Data["last_renewal_time"] = leaseTimes.LastRenewalTime.Unix()
			resp.Data["last_renewal"] = leaseTimes.LastRenewalTime
		}
		if !leaseTimes.ExpireTime.IsZero() {
			resp.Data["expire_time"] = leaseTimes.ExpireTime
			resp.Data["ttl"] = leaseTimes.ttl()
		}
		renewable, _ := leaseTimes.renewable()
		resp.Data["renewable"] = renewable
		resp.Data["issue_time"] = leaseTimes.IssueTime
	}

	if out.EntityID != "" {
		_, identityPolicies, err := ts.core.fetchEntityAndDerivedPolicies(ctx, tokenNS, out.EntityID)
		if err != nil {
			return nil, err
		}
		if len(identityPolicies) != 0 {
			resp.Data["identity_policies"] = identityPolicies[out.NamespaceID]
			delete(identityPolicies, out.NamespaceID)
			resp.Data["external_namespace_policies"] = identityPolicies
		}
	}

	return resp, nil
}

func (ts *TokenStore) handleRenewSelf(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	data.Raw["token"] = req.ClientToken
	return ts.handleRenew(ctx, req, data)
}

// handleRenew handles the auth/token/renew/id path for renewal of tokens.
// This is used to prevent token expiration and revocation.
func (ts *TokenStore) handleRenew(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	id := data.Get("token").(string)
	if id == "" {
		return logical.ErrorResponse("missing token ID"), logical.ErrInvalidRequest
	}
	incrementRaw := data.Get("increment").(int)

	// Convert the increment
	increment := time.Duration(incrementRaw) * time.Second

	// Lookup the token
	te, err := ts.Lookup(ctx, id)
	if err != nil {
		return nil, errwrap.Wrapf("error looking up token to renew: {{err}}", err)
	}
	if te == nil {
		return logical.ErrorResponse("token not found"), logical.ErrInvalidRequest
	}

	var resp *logical.Response

	if te.Type == logical.TokenTypeBatch {
		return logical.ErrorResponse("batch tokens cannot be renewed"), nil
	}

	// Renew the token and its children
	resp, err = ts.expiration.RenewToken(ctx, req, te, increment)

	return resp, err
}

func (ts *TokenStore) authRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if req.Auth == nil {
		return nil, fmt.Errorf("request auth is nil")
	}

	te, err := ts.Lookup(ctx, req.Auth.ClientToken)
	if err != nil {
		return nil, errwrap.Wrapf("error looking up token: {{err}}", err)
	}
	if te == nil {
		return nil, fmt.Errorf("no token entry found during lookup")
	}

	if te.Role == "" {
		req.Auth.Period = te.Period
		req.Auth.ExplicitMaxTTL = te.ExplicitMaxTTL
		return &logical.Response{Auth: req.Auth}, nil
	}

	role, err := ts.tokenStoreRole(ctx, te.Role)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("error looking up role %q: {{err}}", te.Role), err)
	}
	if role == nil {
		return nil, fmt.Errorf("original token role %q could not be found, not renewing", te.Role)
	}

	req.Auth.Period = role.Period
	req.Auth.ExplicitMaxTTL = role.ExplicitMaxTTL
	return &logical.Response{Auth: req.Auth}, nil
}

func (ts *TokenStore) tokenStoreRole(ctx context.Context, name string) (*tsRoleEntry, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	entry, err := ts.rolesView(ns).Get(ctx, name)
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

	if result.TokenType == logical.TokenTypeDefault {
		result.TokenType = logical.TokenTypeDefaultService
	}

	return &result, nil
}

func (ts *TokenStore) tokenStoreRoleList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	entries, err := ts.rolesView(ns).List(ctx, "")
	if err != nil {
		return nil, err
	}

	ret := make([]string, len(entries))
	for i, entry := range entries {
		ret[i] = strings.TrimPrefix(entry, rolesPrefix)
	}

	return logical.ListResponse(ret), nil
}

func (ts *TokenStore) tokenStoreRoleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = ts.rolesView(ns).Delete(ctx, data.Get("role_name").(string))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ts *TokenStore) tokenStoreRoleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	role, err := ts.tokenStoreRole(ctx, data.Get("role_name").(string))
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
			"token_type":          role.TokenType.String(),
		},
	}

	if len(role.BoundCIDRs) > 0 {
		resp.Data["bound_cidrs"] = role.BoundCIDRs
	}

	return resp, nil
}

func (ts *TokenStore) tokenStoreRoleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	name := data.Get("role_name").(string)
	if name == "" {
		return false, fmt.Errorf("role name cannot be empty")
	}
	role, err := ts.tokenStoreRole(ctx, name)
	if err != nil {
		return false, err
	}

	return role != nil, nil
}

func (ts *TokenStore) tokenStoreRoleCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("role_name").(string)
	if name == "" {
		return logical.ErrorResponse("role name cannot be empty"), nil
	}
	entry, err := ts.tokenStoreRole(ctx, name)
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

	boundCIDRsRaw, ok := data.GetOk("bound_cidrs")
	if ok {
		boundCIDRs := boundCIDRsRaw.([]string)
		if len(boundCIDRs) == 0 {
			entry.BoundCIDRs = nil
		} else {
			var parsedCIDRs []*sockaddr.SockAddrMarshaler
			for _, v := range boundCIDRs {
				parsedCIDR, err := sockaddr.NewSockAddr(v)
				if err != nil {
					return logical.ErrorResponse(errwrap.Wrapf(fmt.Sprintf("invalid value %q when parsing bound cidrs: {{err}}", v), err).Error()), nil
				}
				parsedCIDRs = append(parsedCIDRs, &sockaddr.SockAddrMarshaler{parsedCIDR})
			}
			entry.BoundCIDRs = parsedCIDRs
		}
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
				int64(entry.ExplicitMaxTTL.Seconds()), int64(sysView.MaxLeaseTTL().Seconds())))
		}
	}

	pathSuffixInt, ok := data.GetOk("path_suffix")
	if ok {
		pathSuffix := pathSuffixInt.(string)
		switch {
		case pathSuffix != "":
			matched := pathSuffixSanitize.MatchString(pathSuffix)
			if !matched {
				return logical.ErrorResponse(fmt.Sprintf(
					"given role path suffix contains invalid characters; must match %s",
					pathSuffixSanitize.String())), nil
			}
		}
		entry.PathSuffix = pathSuffix
	} else if req.Operation == logical.CreateOperation {
		entry.PathSuffix = data.Get("path_suffix").(string)
	}

	if strings.Contains(entry.PathSuffix, "..") {
		return logical.ErrorResponse(fmt.Sprintf("error registering path suffix: %s", consts.ErrPathContainsParentReferences)), nil
	}

	allowedPoliciesRaw, ok := data.GetOk("allowed_policies")
	if ok {
		entry.AllowedPolicies = policyutil.SanitizePolicies(allowedPoliciesRaw.([]string), policyutil.DoNotAddDefaultPolicy)
	} else if req.Operation == logical.CreateOperation {
		entry.AllowedPolicies = policyutil.SanitizePolicies(data.Get("allowed_policies").([]string), policyutil.DoNotAddDefaultPolicy)
	}

	disallowedPoliciesRaw, ok := data.GetOk("disallowed_policies")
	if ok {
		entry.DisallowedPolicies = strutil.RemoveDuplicates(disallowedPoliciesRaw.([]string), true)
	} else if req.Operation == logical.CreateOperation {
		entry.DisallowedPolicies = strutil.RemoveDuplicates(data.Get("disallowed_policies").([]string), true)
	}

	tokenType := entry.TokenType
	if tokenType == logical.TokenTypeDefault {
		tokenType = logical.TokenTypeDefaultService
	}
	tokenTypeRaw, ok := data.GetOk("token_type")
	if ok {
		tokenTypeStr := tokenTypeRaw.(string)
		switch tokenTypeStr {
		case "service":
			tokenType = logical.TokenTypeService
		case "batch":
			tokenType = logical.TokenTypeBatch
		case "default-service":
			tokenType = logical.TokenTypeDefaultService
		case "default-batch":
			tokenType = logical.TokenTypeDefaultBatch
		default:
			return logical.ErrorResponse(fmt.Sprintf("invalid 'token_type' value %q", tokenTypeStr)), nil
		}
	} else if req.Operation == logical.CreateOperation {
		tokenType = logical.TokenTypeDefaultService
	}
	entry.TokenType = tokenType

	if entry.TokenType == logical.TokenTypeBatch {
		if !entry.Orphan {
			return logical.ErrorResponse("'token_type' cannot be 'batch' when role is set to generate non-orphan tokens"), nil
		}
		if entry.Period != 0 {
			return logical.ErrorResponse("'token_type' cannot be 'batch' when role is set to generate periodic tokens"), nil
		}
		if entry.Renewable {
			return logical.ErrorResponse("'token_type' cannot be 'batch' when role is set to generate renewable tokens"), nil
		}
		if entry.ExplicitMaxTTL != 0 {
			return logical.ErrorResponse("'token_type' cannot be 'batch' when role is set to generate tokens with an explicit max TTL"), nil
		}
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Store it
	jsonEntry, err := logical.StorageEntryJSON(name, entry)
	if err != nil {
		return nil, err
	}
	if err := ts.rolesView(ns).Put(ctx, jsonEntry); err != nil {
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
be used to iterate and discover their properties
or revoke them. Because this can be used to
cause a denial of service, this endpoint
requires 'sudo' capability in addition to
'list'.`
)
