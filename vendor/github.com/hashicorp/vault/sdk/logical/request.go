package logical

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// RequestWrapInfo is a struct that stores information about desired response
// and seal wrapping behavior
type RequestWrapInfo struct {
	// Setting to non-zero specifies that the response should be wrapped.
	// Specifies the desired TTL of the wrapping token.
	TTL time.Duration `json:"ttl" structs:"ttl" mapstructure:"ttl" sentinel:""`

	// The format to use for the wrapped response; if not specified it's a bare
	// token
	Format string `json:"format" structs:"format" mapstructure:"format" sentinel:""`

	// A flag to conforming backends that data for a given request should be
	// seal wrapped
	SealWrap bool `json:"seal_wrap" structs:"seal_wrap" mapstructure:"seal_wrap" sentinel:""`
}

func (r *RequestWrapInfo) SentinelGet(key string) (interface{}, error) {
	if r == nil {
		return nil, nil
	}
	switch key {
	case "ttl":
		return r.TTL, nil
	case "ttl_seconds":
		return int64(r.TTL.Seconds()), nil
	}

	return nil, nil
}

func (r *RequestWrapInfo) SentinelKeys() []string {
	return []string{
		"ttl",
		"ttl_seconds",
	}
}

type ClientTokenSource uint32

const (
	NoClientToken ClientTokenSource = iota
	ClientTokenFromVaultHeader
	ClientTokenFromAuthzHeader
)

// Request is a struct that stores the parameters and context of a request
// being made to Vault. It is used to abstract the details of the higher level
// request protocol from the handlers.
//
// Note: Many of these have Sentinel disabled because they are values populated
// by the router after policy checks; the token namespace would be the right
// place to access them via Sentinel
type Request struct {
	// Id is the uuid associated with each request
	ID string `json:"id" structs:"id" mapstructure:"id" sentinel:""`

	// If set, the name given to the replication secondary where this request
	// originated
	ReplicationCluster string `json:"replication_cluster" structs:"replication_cluster" mapstructure:"replication_cluster" sentinel:""`

	// Operation is the requested operation type
	Operation Operation `json:"operation" structs:"operation" mapstructure:"operation"`

	// Path is the part of the request path not consumed by the
	// routing. As an example, if the original request path is "prod/aws/foo"
	// and the AWS logical backend is mounted at "prod/aws/", then the
	// final path is "foo" since the mount prefix is trimmed.
	Path string `json:"path" structs:"path" mapstructure:"path" sentinel:""`

	// Request data is an opaque map that must have string keys.
	Data map[string]interface{} `json:"map" structs:"data" mapstructure:"data"`

	// Storage can be used to durably store and retrieve state.
	Storage Storage `json:"-" sentinel:""`

	// Secret will be non-nil only for Revoke and Renew operations
	// to represent the secret that was returned prior.
	Secret *Secret `json:"secret" structs:"secret" mapstructure:"secret" sentinel:""`

	// Auth will be non-nil only for Renew operations
	// to represent the auth that was returned prior.
	Auth *Auth `json:"auth" structs:"auth" mapstructure:"auth" sentinel:""`

	// Headers will contain the http headers from the request. This value will
	// be used in the audit broker to ensure we are auditing only the allowed
	// headers.
	Headers map[string][]string `json:"headers" structs:"headers" mapstructure:"headers" sentinel:""`

	// Connection will be non-nil only for credential providers to
	// inspect the connection information and potentially use it for
	// authentication/protection.
	Connection *Connection `json:"connection" structs:"connection" mapstructure:"connection"`

	// ClientToken is provided to the core so that the identity
	// can be verified and ACLs applied. This value is passed
	// through to the logical backends but after being salted and
	// hashed.
	ClientToken string `json:"client_token" structs:"client_token" mapstructure:"client_token" sentinel:""`

	// ClientTokenAccessor is provided to the core so that the it can get
	// logged as part of request audit logging.
	ClientTokenAccessor string `json:"client_token_accessor" structs:"client_token_accessor" mapstructure:"client_token_accessor" sentinel:""`

	// DisplayName is provided to the logical backend to help associate
	// dynamic secrets with the source entity. This is not a sensitive
	// name, but is useful for operators.
	DisplayName string `json:"display_name" structs:"display_name" mapstructure:"display_name" sentinel:""`

	// MountPoint is provided so that a logical backend can generate
	// paths relative to itself. The `Path` is effectively the client
	// request path with the MountPoint trimmed off.
	MountPoint string `json:"mount_point" structs:"mount_point" mapstructure:"mount_point" sentinel:""`

	// MountType is provided so that a logical backend can make decisions
	// based on the specific mount type (e.g., if a mount type has different
	// aliases, generating different defaults depending on the alias)
	MountType string `json:"mount_type" structs:"mount_type" mapstructure:"mount_type" sentinel:""`

	// MountAccessor is provided so that identities returned by the authentication
	// backends can be tied to the mount it belongs to.
	MountAccessor string `json:"mount_accessor" structs:"mount_accessor" mapstructure:"mount_accessor" sentinel:""`

	// WrapInfo contains requested response wrapping parameters
	WrapInfo *RequestWrapInfo `json:"wrap_info" structs:"wrap_info" mapstructure:"wrap_info" sentinel:""`

	// ClientTokenRemainingUses represents the allowed number of uses left on the
	// token supplied
	ClientTokenRemainingUses int `json:"client_token_remaining_uses" structs:"client_token_remaining_uses" mapstructure:"client_token_remaining_uses"`

	// EntityID is the identity of the caller extracted out of the token used
	// to make this request
	EntityID string `json:"entity_id" structs:"entity_id" mapstructure:"entity_id" sentinel:""`

	// PolicyOverride indicates that the requestor wishes to override
	// soft-mandatory Sentinel policies
	PolicyOverride bool `json:"policy_override" structs:"policy_override" mapstructure:"policy_override"`

	// Whether the request is unauthenticated, as in, had no client token
	// attached. Useful in some situations where the client token is not made
	// accessible.
	Unauthenticated bool `json:"unauthenticated" structs:"unauthenticated" mapstructure:"unauthenticated"`

	// MFACreds holds the parsed MFA information supplied over the API as part of
	// X-Vault-MFA header
	MFACreds MFACreds `json:"mfa_creds" structs:"mfa_creds" mapstructure:"mfa_creds" sentinel:""`

	// Cached token entry. This avoids another lookup in request handling when
	// we've already looked it up at http handling time. Note that this token
	// has not been "used", as in it will not properly take into account use
	// count limitations. As a result this field should only ever be used for
	// transport to a function that would otherwise do a lookup and then
	// properly use the token.
	tokenEntry *TokenEntry

	// For replication, contains the last WAL on the remote side after handling
	// the request, used for best-effort avoidance of stale read-after-write
	lastRemoteWAL uint64

	// ControlGroup holds the authorizations that have happened on this
	// request
	ControlGroup *ControlGroup `json:"control_group" structs:"control_group" mapstructure:"control_group" sentinel:""`

	// ClientTokenSource tells us where the client token was sourced from, so
	// we can delete it before sending off to plugins
	ClientTokenSource ClientTokenSource

	// HTTPRequest, if set, can be used to access fields from the HTTP request
	// that generated this logical.Request object, such as the request body.
	HTTPRequest *http.Request `json:"-" sentinel:""`

	// ResponseWriter if set can be used to stream a response value to the http
	// request that generated this logical.Request object.
	ResponseWriter *HTTPResponseWriter `json:"-" sentinel:""`
}

// Get returns a data field and guards for nil Data
func (r *Request) Get(key string) interface{} {
	if r.Data == nil {
		return nil
	}
	return r.Data[key]
}

// GetString returns a data field as a string
func (r *Request) GetString(key string) string {
	raw := r.Get(key)
	s, _ := raw.(string)
	return s
}

func (r *Request) GoString() string {
	return fmt.Sprintf("*%#v", *r)
}

func (r *Request) SentinelGet(key string) (interface{}, error) {
	switch key {
	case "path":
		// Sanitize it here so that it's consistent in policies
		return strings.TrimPrefix(r.Path, "/"), nil

	case "wrapping", "wrap_info":
		// If the pointer is nil accessing the wrap info is considered
		// "undefined" so this allows us to instead discover a TTL of zero
		if r.WrapInfo == nil {
			return &RequestWrapInfo{}, nil
		}
		return r.WrapInfo, nil
	}

	return nil, nil
}

func (r *Request) SentinelKeys() []string {
	return []string{
		"path",
		"wrapping",
		"wrap_info",
	}
}

func (r *Request) LastRemoteWAL() uint64 {
	return r.lastRemoteWAL
}

func (r *Request) SetLastRemoteWAL(last uint64) {
	r.lastRemoteWAL = last
}

func (r *Request) TokenEntry() *TokenEntry {
	return r.tokenEntry
}

func (r *Request) SetTokenEntry(te *TokenEntry) {
	r.tokenEntry = te
}

// RenewRequest creates the structure of the renew request.
func RenewRequest(path string, secret *Secret, data map[string]interface{}) *Request {
	return &Request{
		Operation: RenewOperation,
		Path:      path,
		Data:      data,
		Secret:    secret,
	}
}

// RenewAuthRequest creates the structure of the renew request for an auth.
func RenewAuthRequest(path string, auth *Auth, data map[string]interface{}) *Request {
	return &Request{
		Operation: RenewOperation,
		Path:      path,
		Data:      data,
		Auth:      auth,
	}
}

// RevokeRequest creates the structure of the revoke request.
func RevokeRequest(path string, secret *Secret, data map[string]interface{}) *Request {
	return &Request{
		Operation: RevokeOperation,
		Path:      path,
		Data:      data,
		Secret:    secret,
	}
}

// RollbackRequest creates the structure of the revoke request.
func RollbackRequest(path string) *Request {
	return &Request{
		Operation: RollbackOperation,
		Path:      path,
		Data:      make(map[string]interface{}),
	}
}

// Operation is an enum that is used to specify the type
// of request being made
type Operation string

const (
	// The operations below are called per path
	CreateOperation         Operation = "create"
	ReadOperation                     = "read"
	UpdateOperation                   = "update"
	DeleteOperation                   = "delete"
	ListOperation                     = "list"
	HelpOperation                     = "help"
	AliasLookaheadOperation           = "alias-lookahead"

	// The operations below are called globally, the path is less relevant.
	RevokeOperation   Operation = "revoke"
	RenewOperation              = "renew"
	RollbackOperation           = "rollback"
)

type MFACreds map[string][]string

// InitializationRequest stores the parameters and context of an Initialize()
// call being made to a logical.Backend.
type InitializationRequest struct {

	// Storage can be used to durably store and retrieve state.
	Storage Storage
}
