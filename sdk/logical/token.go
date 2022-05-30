package logical

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"time"

	sockaddr "github.com/hashicorp/go-sockaddr"
)

type TokenType uint8

const (
	// TokenTypeDefault means "use the default, if any, that is currently set
	// on the mount". If not set, results in a Service token.
	TokenTypeDefault TokenType = iota

	// TokenTypeService is a "normal" Vault token for long-lived services
	TokenTypeService

	// TokenTypeBatch is a batch token
	TokenTypeBatch

	// TokenTypeDefaultService configured on a mount, means that if
	// TokenTypeDefault is sent back by the mount, create Service tokens
	TokenTypeDefaultService

	// TokenTypeDefaultBatch configured on a mount, means that if
	// TokenTypeDefault is sent back by the mount, create Batch tokens
	TokenTypeDefaultBatch

	// ClientIDTWEDelimiter Delimiter between the string fields used to generate a client
	// ID for tokens without entities. This is the 0 character, which
	// is a non-printable string. Please see unicode.IsPrint for details.
	ClientIDTWEDelimiter = rune('\x00')

	// SortedPoliciesTWEDelimiter Delimiter between each policy in the sorted policies used to
	// generate a client ID for tokens without entities. This is the 127
	// character, which is a non-printable string. Please see unicode.IsPrint
	// for details.
	SortedPoliciesTWEDelimiter = rune('\x7F')
)

func (t *TokenType) UnmarshalJSON(b []byte) error {
	if len(b) == 1 {
		*t = TokenType(b[0] - '0')
		return nil
	}

	// Handle upgrade from pre-1.2 where we were serialized as string:
	s := string(b)
	switch s {
	case `"default"`, `""`:
		*t = TokenTypeDefault
	case `"service"`:
		*t = TokenTypeService
	case `"batch"`:
		*t = TokenTypeBatch
	case `"default-service"`:
		*t = TokenTypeDefaultService
	case `"default-batch"`:
		*t = TokenTypeDefaultBatch
	default:
		return fmt.Errorf("unknown token type %q", s)
	}
	return nil
}

func (t TokenType) String() string {
	switch t {
	case TokenTypeDefault:
		return "default"
	case TokenTypeService:
		return "service"
	case TokenTypeBatch:
		return "batch"
	case TokenTypeDefaultService:
		return "default-service"
	case TokenTypeDefaultBatch:
		return "default-batch"
	default:
		panic("unreachable")
	}
}

// TokenEntry is used to represent a given token
type TokenEntry struct {
	Type TokenType `json:"type" mapstructure:"type" structs:"type" sentinel:""`

	// ID of this entry, generally a random UUID
	ID string `json:"id" mapstructure:"id" structs:"id" sentinel:""`

	// ExternalID is the ID of a newly created service
	// token that will be returned to a user
	ExternalID string `json:"-"`

	// Accessor for this token, a random UUID
	Accessor string `json:"accessor" mapstructure:"accessor" structs:"accessor" sentinel:""`

	// Parent token, used for revocation trees
	Parent string `json:"parent" mapstructure:"parent" structs:"parent" sentinel:""`

	// Which named policies should be used
	Policies []string `json:"policies" mapstructure:"policies" structs:"policies"`

	// InlinePolicy specifies ACL rules to be applied to this token entry.
	InlinePolicy string `json:"inline_policy" mapstructure:"inline_policy" structs:"inline_policy"`

	// Used for audit trails, this is something like "auth/user/login"
	Path string `json:"path" mapstructure:"path" structs:"path"`

	// Used for auditing. This could include things like "source", "user", "ip"
	Meta map[string]string `json:"meta" mapstructure:"meta" structs:"meta" sentinel:"meta"`

	// InternalMeta is used to store internal metadata. This metadata will not be audit logged or returned from lookup APIs.
	InternalMeta map[string]string `json:"internal_meta" mapstructure:"internal_meta" structs:"internal_meta"`

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
	CreationTime int64 `json:"creation_time" mapstructure:"creation_time" structs:"creation_time" sentinel:""`

	// Duration set when token was created
	TTL time.Duration `json:"ttl" mapstructure:"ttl" structs:"ttl" sentinel:""`

	// Explicit maximum TTL on the token
	ExplicitMaxTTL time.Duration `json:"explicit_max_ttl" mapstructure:"explicit_max_ttl" structs:"explicit_max_ttl" sentinel:""`

	// If set, the role that was used for parameters at creation time
	Role string `json:"role" mapstructure:"role" structs:"role"`

	// If set, the period of the token. This is only used when created directly
	// through the create endpoint; periods managed by roles or other auth
	// backends are subject to those renewal rules.
	Period time.Duration `json:"period" mapstructure:"period" structs:"period" sentinel:""`

	// These are the deprecated fields
	DisplayNameDeprecated    string        `json:"DisplayName" mapstructure:"DisplayName" structs:"DisplayName" sentinel:""`
	NumUsesDeprecated        int           `json:"NumUses" mapstructure:"NumUses" structs:"NumUses" sentinel:""`
	CreationTimeDeprecated   int64         `json:"CreationTime" mapstructure:"CreationTime" structs:"CreationTime" sentinel:""`
	ExplicitMaxTTLDeprecated time.Duration `json:"ExplicitMaxTTL" mapstructure:"ExplicitMaxTTL" structs:"ExplicitMaxTTL" sentinel:""`

	// EntityID is the ID of the entity associated with this token.
	EntityID string `json:"entity_id" mapstructure:"entity_id" structs:"entity_id"`

	// If NoIdentityPolicies is true, the token will not inherit
	// identity policies from the associated EntityID.
	NoIdentityPolicies bool `json:"no_identity_policies" mapstructure:"no_identity_policies" structs:"no_identity_policies"`

	// The set of CIDRs that this token can be used with
	BoundCIDRs []*sockaddr.SockAddrMarshaler `json:"bound_cidrs" sentinel:""`

	// NamespaceID is the identifier of the namespace to which this token is
	// confined to. Do not return this value over the API when the token is
	// being looked up.
	NamespaceID string `json:"namespace_id" mapstructure:"namespace_id" structs:"namespace_id" sentinel:""`

	// CubbyholeID is the identifier of the cubbyhole storage belonging to this
	// token
	CubbyholeID string `json:"cubbyhole_id" mapstructure:"cubbyhole_id" structs:"cubbyhole_id" sentinel:""`
}

// CreateClientID returns the client ID, and a boolean which is false if the clientID
// has an entity, and true otherwise
func (te *TokenEntry) CreateClientID() (string, bool) {
	var clientIDInputBuilder strings.Builder

	// if entry has an associated entity ID, return it
	if te.EntityID != "" {
		return te.EntityID, false
	}

	// The entry is associated with a TWE (token without entity). In this case
	// we must create a client ID by calculating the following formula:
	// clientID = SHA256(sorted policies + namespace)

	// Step 1: Copy entry policies to a new struct
	sortedPolicies := make([]string, len(te.Policies))
	copy(sortedPolicies, te.Policies)

	// Step 2: Sort and join copied policies
	sort.Strings(sortedPolicies)
	for _, pol := range sortedPolicies {
		clientIDInputBuilder.WriteRune(SortedPoliciesTWEDelimiter)
		clientIDInputBuilder.WriteString(pol)
	}

	// Step 3: Add namespace ID
	clientIDInputBuilder.WriteRune(ClientIDTWEDelimiter)
	clientIDInputBuilder.WriteString(te.NamespaceID)

	if clientIDInputBuilder.Len() == 0 {
		return "", true
	}
	// Step 4: Remove the first character in the string, as it's an unnecessary delimiter
	clientIDInput := clientIDInputBuilder.String()[1:]

	// Step 5: Hash the sum
	hashed := sha256.Sum256([]byte(clientIDInput))
	return base64.StdEncoding.EncodeToString(hashed[:]), true
}

func (te *TokenEntry) SentinelGet(key string) (interface{}, error) {
	if te == nil {
		return nil, nil
	}
	switch key {
	case "policies":
		return te.Policies, nil

	case "path":
		return te.Path, nil

	case "display_name":
		return te.DisplayName, nil

	case "num_uses":
		return te.NumUses, nil

	case "role":
		return te.Role, nil

	case "entity_id":
		return te.EntityID, nil

	case "period":
		return te.Period, nil

	case "period_seconds":
		return int64(te.Period.Seconds()), nil

	case "explicit_max_ttl":
		return te.ExplicitMaxTTL, nil

	case "explicit_max_ttl_seconds":
		return int64(te.ExplicitMaxTTL.Seconds()), nil

	case "creation_ttl":
		return te.TTL, nil

	case "creation_ttl_seconds":
		return int64(te.TTL.Seconds()), nil

	case "creation_time":
		return time.Unix(te.CreationTime, 0).Format(time.RFC3339Nano), nil

	case "creation_time_unix":
		return time.Unix(te.CreationTime, 0), nil

	case "meta", "metadata":
		return te.Meta, nil

	case "type":
		teType := te.Type
		switch teType {
		case TokenTypeBatch, TokenTypeService:
		case TokenTypeDefault:
			teType = TokenTypeService
		default:
			return "unknown", nil
		}
		return teType.String(), nil
	}

	return nil, nil
}

func (te *TokenEntry) SentinelKeys() []string {
	return []string{
		"period",
		"period_seconds",
		"explicit_max_ttl",
		"explicit_max_ttl_seconds",
		"creation_ttl",
		"creation_ttl_seconds",
		"creation_time",
		"creation_time_unix",
		"meta",
		"metadata",
		"type",
	}
}

// IsRoot returns false if the token is not root (or doesn't exist)
func (te *TokenEntry) IsRoot() bool {
	if te == nil {
		return false
	}

	return len(te.Policies) == 1 && te.Policies[0] == "root"
}
