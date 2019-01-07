package tokenhelper

import (
	"errors"
	"fmt"
	"time"

	sockaddr "github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

type TokenParams struct {
	// The set of CIDRs that tokens generated using this role will be bound to
	BoundCIDRs []*sockaddr.SockAddrMarshaler `json:"bound_cidrs"`

	// If set, the token entry will have an explicit maximum TTL set, rather
	// than deferring to role/mount values
	ExplicitMaxTTL time.Duration `json:"explicit_max_ttl" mapstructure:"explicit_max_ttl"`

	MaxTTL time.Duration `json:"max_ttl" mapstructure:"max_ttl"`

	// If non-zero, tokens created using this role will be able to be renewed
	// forever, but will have a fixed renewal period of this value
	Period time.Duration `json:"period" mapstructure:"period"`

	Policies []string `json:"policies" mapstructure:"policies"`

	// If set, controls whether created tokens are marked as being renewable
	Renewable bool `json:"renewable" mapstructure:"renewable"`

	// The type of token this role should issue
	TokenType logical.TokenType `json:"token_type" mapstructure:"token_type"`

	TTL time.Duration `json:"ttl" mapstructure:"ttl"`
}

// AddTokenFields adds fields to an existing role. It panics if it would
// overwrite an existing field.
func AddTokenFields(m map[string]*framework.FieldSchema) {
	AddTokenFieldsWithAllowList(m, nil)
}

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

func TokenFields() map[string]*framework.FieldSchema {
	return map[string]*framework.FieldSchema{
		"bound_cidrs": &framework.FieldSchema{
			Type:        framework.TypeCommaStringSlice,
			Description: `Comma separated string or JSON list of CIDR blocks. If set, specifies the blocks of IP addresses which are allowed to use the generated token.`,
		},

		"explicit_max_ttl": &framework.FieldSchema{
			Type:        framework.TypeDurationSecond,
			Description: tokenExplicitMaxTTLHelp,
		},

		"max_ttl": &framework.FieldSchema{
			Type:        framework.TypeDurationSecond,
			Description: "The maximum lifetime of the generated token",
		},

		"period": &framework.FieldSchema{
			Type:        framework.TypeDurationSecond,
			Description: tokenPeriodHelp,
		},

		"policies": &framework.FieldSchema{
			Type:        framework.TypeCommaStringSlice,
			Description: "Comma-separated list of policies",
		},

		"renewable": &framework.FieldSchema{
			Type:        framework.TypeBool,
			Default:     true,
			Description: tokenRenewableHelp,
		},

		"token_type": &framework.FieldSchema{
			Type:        framework.TypeString,
			Description: "The type of token to generate, service or batch",
		},

		"ttl": &framework.FieldSchema{
			Type:        framework.TypeDurationSecond,
			Description: "The initial ttl of the token to generate",
		},
	}
}

func (t *TokenParams) ParseTokenFields(req *logical.Request, d *framework.FieldData) error {
	if boundCIDRsRaw, ok := d.GetOk("bound_cidrs"); ok {
		boundCIDRs, err := parseutil.ParseAddrs(boundCIDRsRaw.([]string))
		if err != nil {
			return err
		}
		t.BoundCIDRs = boundCIDRs
	}

	if explicitMaxTTLRaw, ok := d.GetOk("explicit_max_ttl"); ok {
		t.ExplicitMaxTTL = time.Duration(explicitMaxTTLRaw.(int)) * time.Second
	}

	if maxTTLRaw, ok := d.GetOk("max_ttl"); ok {
		t.MaxTTL = time.Duration(maxTTLRaw.(int)) * time.Second
	} else if maxTTLRaw, ok := d.GetOk("lease_max"); ok {
		t.MaxTTL = time.Duration(maxTTLRaw.(int)) * time.Second
	}
	if t.MaxTTL < 0 {
		return errors.New("'max_ttl' cannot be negative")
	}

	if periodRaw, ok := d.GetOk("period"); ok {
		t.Period = time.Duration(periodRaw.(int)) * time.Second
	}
	if t.Period < 0 {
		return errors.New("'period' cannot be negative")
	}

	if policiesRaw, ok := d.GetOk("policies"); ok {
		t.Policies = policiesRaw.([]string)
	}

	if renewableRaw, ok := d.GetOk("renewable"); ok {
		t.Renewable = renewableRaw.(bool)
	}

	if tokenTypeRaw, ok := d.GetOk("token_type"); ok {
		var tokenType logical.TokenType
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
			return fmt.Errorf("invalid 'token_type' value %q", tokenTypeStr)
		}
		t.TokenType = tokenType
	}

	if ttlRaw, ok := d.GetOk("ttl"); ok {
		t.TTL = time.Duration(ttlRaw.(int)) * time.Second
	} else if ttlRaw, ok := d.GetOk("lease"); ok {
		t.TTL = time.Duration(ttlRaw.(int)) * time.Second
	}
	if t.TTL < 0 {
		return errors.New("'ttl' cannot be negative")
	}

	return nil
}

func (t TokenParams) PopulateTokenData(m map[string]interface{}) {
	m["bound_cidrs"] = t.BoundCIDRs
	m["explicit_max_ttl"] = t.ExplicitMaxTTL.Seconds()
	m["max_ttl"] = t.MaxTTL.Seconds()
	m["period"] = t.Period.Seconds()
	m["policies"] = t.Policies
	m["renewable"] = t.Renewable
	m["token_type"] = t.TokenType.String()
	m["ttl"] = t.TTL.Seconds()
}

func (t TokenParams) PopulateTokenAuth(auth *logical.Auth) {
	auth.BoundCIDRs = t.BoundCIDRs
	auth.ExplicitMaxTTL = t.ExplicitMaxTTL
	auth.MaxTTL = t.MaxTTL
	auth.Period = t.Period
	auth.Policies = t.Policies
	auth.Renewable = t.Renewable
	auth.TokenType = t.TokenType
	auth.TTL = t.TTL
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
	tokenRenewableHelp = `Tokens created via this role will be
renewable or not according to this value.
Defaults to "true".`
)
