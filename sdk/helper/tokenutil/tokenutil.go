package tokenutil

import (
	"errors"
	"fmt"
	"time"

	sockaddr "github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// TokenParams contains a set of common parameters that auth plugins can use
// for setting token behavior
type TokenParams struct {
	// The set of CIDRs that tokens generated using this role will be bound to
	TokenBoundCIDRs []*sockaddr.SockAddrMarshaler `json:"token_bound_cidrs"`

	// If set, the token entry will have an explicit maximum TTL set, rather
	// than deferring to role/mount values
	TokenExplicitMaxTTL time.Duration `json:"token_explicit_max_ttl" mapstructure:"token_explicit_max_ttl"`

	// The max TTL to use for the token
	TokenMaxTTL time.Duration `json:"token_max_ttl" mapstructure:"token_max_ttl"`

	// If set, core will not automatically add default to the policy list
	TokenNoDefaultPolicy bool `json:"token_no_default_policy" mapstructure:"token_no_default_policy"`

	// The maximum number of times a token issued from this role may be used.
	TokenNumUses int `json:"token_num_uses" mapstructure:"token_num_uses"`

	// If non-zero, tokens created using this role will be able to be renewed
	// forever, but will have a fixed renewal period of this value
	TokenPeriod time.Duration `json:"token_period" mapstructure:"token_period"`

	// The policies to set
	TokenPolicies []string `json:"token_policies" mapstructure:"token_policies"`

	// The type of token this role should issue
	TokenType logical.TokenType `json:"token_type" mapstructure:"token_type"`

	// The TTL to user for the token
	TokenTTL time.Duration `json:"token_ttl" mapstructure:"token_ttl"`
}

// AddTokenFields adds fields to an existing role. It panics if it would
// overwrite an existing field.
func AddTokenFields(m map[string]*framework.FieldSchema) {
	AddTokenFieldsWithAllowList(m, nil)
}

// AddTokenFields adds fields to an existing role. It panics if it would
// overwrite an existing field. Allowed can be use to restrict the set, e.g. if
// there would be conflicts.
func AddTokenFieldsWithAllowList(m map[string]*framework.FieldSchema, allowed []string) {
	r := TokenFields()
	for k, v := range r {
		if len(allowed) > 0 && !strutil.StrListContains(allowed, k) {
			continue
		}
		if _, has := m[k]; has {
			panic(fmt.Sprintf("adding role field %s would overwrite existing field", k))
		}
		m[k] = v
	}
}

// TokenFields provides a set of field schemas for the parameters
func TokenFields() map[string]*framework.FieldSchema {
	return map[string]*framework.FieldSchema{
		"token_bound_cidrs": &framework.FieldSchema{
			Type:        framework.TypeCommaStringSlice,
			Description: `Comma separated string or JSON list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.`,
		},

		"token_explicit_max_ttl": &framework.FieldSchema{
			Type:        framework.TypeDurationSecond,
			Description: tokenExplicitMaxTTLHelp,
		},

		"token_max_ttl": &framework.FieldSchema{
			Type:        framework.TypeDurationSecond,
			Description: "The maximum lifetime of the generated token",
		},

		"token_no_default_policy": &framework.FieldSchema{
			Type:        framework.TypeBool,
			Description: "If true, the 'default' policy will not automatically be added to generated tokens",
		},

		"token_period": &framework.FieldSchema{
			Type:        framework.TypeDurationSecond,
			Description: tokenPeriodHelp,
		},

		"token_policies": &framework.FieldSchema{
			Type:        framework.TypeCommaStringSlice,
			Description: "Comma-separated list of policies",
		},

		"token_type": &framework.FieldSchema{
			Type:        framework.TypeString,
			Default:     "default-service",
			Description: "The type of token to generate, service or batch",
		},

		"token_ttl": &framework.FieldSchema{
			Type:        framework.TypeDurationSecond,
			Description: "The initial ttl of the token to generate",
		},

		"token_num_uses": &framework.FieldSchema{
			Type:        framework.TypeInt,
			Description: "The maximum number of times a token may be used, a value of zero means unlimited",
		},
	}
}

// ParseTokenFields provides common field parsing functionality into a TokenFields struct
func (t *TokenParams) ParseTokenFields(req *logical.Request, d *framework.FieldData) error {
	if boundCIDRsRaw, ok := d.GetOk("token_bound_cidrs"); ok {
		boundCIDRs, err := parseutil.ParseAddrs(boundCIDRsRaw.([]string))
		if err != nil {
			return err
		}
		t.TokenBoundCIDRs = boundCIDRs
	}

	if explicitMaxTTLRaw, ok := d.GetOk("token_explicit_max_ttl"); ok {
		t.TokenExplicitMaxTTL = time.Duration(explicitMaxTTLRaw.(int)) * time.Second
	}

	if maxTTLRaw, ok := d.GetOk("token_max_ttl"); ok {
		t.TokenMaxTTL = time.Duration(maxTTLRaw.(int)) * time.Second
	}
	if t.TokenMaxTTL < 0 {
		return errors.New("'token_max_ttl' cannot be negative")
	}

	if noDefaultRaw, ok := d.GetOk("token_no_default_policy"); ok {
		t.TokenNoDefaultPolicy = noDefaultRaw.(bool)
	}

	if periodRaw, ok := d.GetOk("token_period"); ok {
		t.TokenPeriod = time.Duration(periodRaw.(int)) * time.Second
	}
	if t.TokenPeriod < 0 {
		return errors.New("'token_period' cannot be negative")
	}

	if policiesRaw, ok := d.GetOk("token_policies"); ok {
		t.TokenPolicies = policiesRaw.([]string)
	}

	if tokenTypeRaw, ok := d.GetOk("token_type"); ok {
		var tokenType logical.TokenType
		tokenTypeStr := tokenTypeRaw.(string)
		switch tokenTypeStr {
		case "service":
			tokenType = logical.TokenTypeService
		case "batch":
			tokenType = logical.TokenTypeBatch
		case "", "default", "default-service":
			tokenType = logical.TokenTypeDefaultService
		case "default-batch":
			tokenType = logical.TokenTypeDefaultBatch
		default:
			return fmt.Errorf("invalid 'token_type' value %q", tokenTypeStr)
		}
		t.TokenType = tokenType
	}

	if t.TokenType == logical.TokenTypeBatch || t.TokenType == logical.TokenTypeDefaultBatch {
		if t.TokenPeriod != 0 {
			return errors.New("'token_type' cannot be 'batch' or 'default_batch' when set to generate periodic tokens")
		}
		if t.TokenNumUses != 0 {
			return errors.New("'token_type' cannot be 'batch' or 'default_batch' when set to generate tokens with limited use count")
		}
	}

	if ttlRaw, ok := d.GetOk("token_ttl"); ok {
		t.TokenTTL = time.Duration(ttlRaw.(int)) * time.Second
	}
	if t.TokenTTL < 0 {
		return errors.New("'token_ttl' cannot be negative")
	}
	if t.TokenTTL > 0 && t.TokenMaxTTL > 0 && t.TokenTTL > t.TokenMaxTTL {
		return errors.New("'token_ttl' cannot be greater than 'token_max_ttl'")
	}

	if tokenNumUses, ok := d.GetOk("token_num_uses"); ok {
		t.TokenNumUses = tokenNumUses.(int)
	}
	if t.TokenNumUses < 0 {
		return errors.New("'token_num_uses' cannot be negative")
	}

	return nil
}

// PopulateTokenData adds information from TokenParams into the map
func (t *TokenParams) PopulateTokenData(m map[string]interface{}) {
	m["token_bound_cidrs"] = t.TokenBoundCIDRs
	m["token_explicit_max_ttl"] = int64(t.TokenExplicitMaxTTL.Seconds())
	m["token_max_ttl"] = int64(t.TokenMaxTTL.Seconds())
	m["token_no_default_policy"] = t.TokenNoDefaultPolicy
	m["token_period"] = int64(t.TokenPeriod.Seconds())
	m["token_policies"] = t.TokenPolicies
	m["token_type"] = t.TokenType.String()
	m["token_ttl"] = int64(t.TokenTTL.Seconds())
	m["token_num_uses"] = t.TokenNumUses

	if len(t.TokenPolicies) == 0 {
		m["token_policies"] = []string{}
	}

	if len(t.TokenBoundCIDRs) == 0 {
		m["token_bound_cidrs"] = []string{}
	}
}

// PopulateTokenAuth populates Auth with parameters
func (t *TokenParams) PopulateTokenAuth(auth *logical.Auth) {
	auth.BoundCIDRs = t.TokenBoundCIDRs
	auth.ExplicitMaxTTL = t.TokenExplicitMaxTTL
	auth.MaxTTL = t.TokenMaxTTL
	auth.NoDefaultPolicy = t.TokenNoDefaultPolicy
	auth.Period = t.TokenPeriod
	auth.Policies = t.TokenPolicies
	auth.TokenType = t.TokenType
	auth.TTL = t.TokenTTL
	auth.NumUses = t.TokenNumUses
}

func DeprecationText(param string) string {
	return fmt.Sprintf("Use %q instead. If this and %q are both specified, only %q will be used.", param, param, param)
}

const (
	tokenPeriodHelp = `If set, tokens created via this role
will have no max lifetime; instead, their
renewal period will be fixed to this value.
This takes an integer number of seconds,
or a string duration (e.g. "24h").`
	tokenExplicitMaxTTLHelp = `If set, tokens created via this role
carry an explicit maximum TTL. During renewal,
the current maximum TTL values of the role
and the mount are not checked for changes,
and any updates to these values will have
no effect on the token being renewed.`
)
