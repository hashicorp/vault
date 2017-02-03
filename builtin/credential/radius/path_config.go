package radius

import (
	"strings"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `config`,
		Fields: map[string]*framework.FieldSchema{
			"host": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "radius host to connect",
			},

			"port": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Default:     1812,
				Description: "radius port (default: 1812)",
			},
			"secret": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "secret shared with the radius server",
			},
			"allow_unknown_users": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     false,
				Description: "enable granting default policy upon successful RADIUS authentication (default: false)",
			},
			"reauth_on_renew": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Default:     false,
				Description: "attempt reauthentication with backend before granting token renewal (default: false)",
			},
			"dial_timeout": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Default:     10,
				Description: "number of seconds before connect timeouts (default: 10)",
			},
			"read_timeout": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Default:     10,
				Description: "number of seconds before response timeouts (default: 10)",
			},
			"nas_port": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Default:     10,
				Description: "RADIUS NAS port field (default: 10)",
			},
		},

		ExistenceCheck: b.configExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConfigRead,
			logical.CreateOperation: b.pathConfigCreateUpdate,
			logical.UpdateOperation: b.pathConfigCreateUpdate,
		},

		HelpSynopsis:    pathConfigHelpSyn,
		HelpDescription: pathConfigHelpDesc,
	}
}

// Establishes dichotomy of request operation between CreateOperation and UpdateOperation.
// Returning 'true' forces an UpdateOperation, CreateOperation otherwise.
func (b *backend) configExistenceCheck(req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := b.Config(req)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

/*
 * Construct ConfigEntry struct using stored configuration.
 */
func (b *backend) Config(req *logical.Request) (*ConfigEntry, error) {

	storedConfig, err := req.Storage.Get("config")
	if err != nil {
		return nil, err
	}

	if storedConfig == nil {
		// No user overrides, return default configuration
		return nil, nil
	}

	var result ConfigEntry

	if err := storedConfig.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathConfigRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	cfg, err := b.Config(req)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: structs.New(cfg).Map(),
	}
	resp.AddWarning("Read access to this endpoint should be controlled via ACLs as it will return the configuration information as-is, including any secrets.")
	return resp, nil
}

func (b *backend) pathConfigCreateUpdate(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	// Build a ConfigEntry struct out of the supplied FieldData
	cfg, err := b.Config(req)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = &ConfigEntry{}
	}

	host, ok := d.GetOk("host")
	if ok {
		cfg.Host = strings.ToLower(host.(string))
	} else if req.Operation == logical.CreateOperation {
		cfg.Host = strings.ToLower(d.Get("host").(string))
	}

	port, ok := d.GetOk("port")
	if ok {
		cfg.Port = port.(int)
	} else if req.Operation == logical.CreateOperation {
		cfg.Port = d.Get("port").(int)
	}

	secret, ok := d.GetOk("secret")
	if ok {
		cfg.Secret = secret.(string)
	} else if req.Operation == logical.CreateOperation {
		cfg.Secret = d.Get("secret").(string)
	}

	allowUnknownUsers, ok := d.GetOk("allow_unknown_users")
	if ok {
		cfg.AllowUnknownUsers = allowUnknownUsers.(bool)
	} else if req.Operation == logical.CreateOperation {
		cfg.AllowUnknownUsers = d.Get("allow_unknown_users").(bool)
	}

	reauthOnRenew, ok := d.GetOk("reauth_on_renew")
	if ok {
		cfg.ReauthOnRenew = reauthOnRenew.(bool)
	} else if req.Operation == logical.CreateOperation {
		cfg.ReauthOnRenew = d.Get("reauth_on_renew").(bool)
	}

	dialTimeout, ok := d.GetOk("dial_timeout")
	if ok {
		cfg.DialTimeout = dialTimeout.(int)
	} else if req.Operation == logical.CreateOperation {
		cfg.DialTimeout = d.Get("dial_timeout").(int)
	}

	readTimeout, ok := d.GetOk("read_timeout")
	if ok {
		cfg.ReadTimeout = readTimeout.(int)
	} else if req.Operation == logical.CreateOperation {
		cfg.ReadTimeout = d.Get("read_timeout").(int)
	}

	nasPort, ok := d.GetOk("nas_port")
	if ok {
		cfg.NasPort = nasPort.(int)
	} else if req.Operation == logical.CreateOperation {
		cfg.NasPort = d.Get("nas_port").(int)
	}

	entry, err := logical.StorageEntryJSON("config", cfg)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

type ConfigEntry struct {
	Host              string `json:"host" structs:"host" mapstructure:"host"`
	Port              int    `json:"port" structs:"port" mapstructure:"port"`
	Secret            string `json:"secret" structs:"secret" mapstructure:"secret"`
	AllowUnknownUsers bool   `json:"allow_unknown_users" structs:"allow_unknown_users" mapstructure:"allow_unknown_users"`
	ReauthOnRenew     bool   `json:"reauth_on_renew" structs:"reauth_on_renew" mapstructure:"reauth_on_renew"`
	DialTimeout       int    `json:"dial_timeout" structs:"dial_timeout" mapstructure:"dial_timeout"`
	ReadTimeout       int    `json:"read_timeout" structs:"read_timeout" mapstructure:"read_timeout"`
	NasPort           int    `json:"nas_port" structs:"nas_port" mapstructure:"nas_port"`
}

const pathConfigHelpSyn = `
Configure the RADIUS server to connect to, along with its options.
`

const pathConfigHelpDesc = `
This endpoint allows you to configure the RADIUS server to connect to and its
configuration options.
`
