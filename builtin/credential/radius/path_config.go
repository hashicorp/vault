package radius

import (
	"strings"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"host": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "RADIUS server host",
			},

			"port": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Default:     1812,
				Description: "RADIUS server port (default: 1812)",
			},
			"secret": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Secret shared with the RADIUS server",
			},
			"unregistered_user_policies": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: "Comma-separated list of policies to grant upon successful RADIUS authentication of an unregisted user (default: emtpy)",
			},
			"dial_timeout": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Default:     10,
				Description: "Number of seconds before connect times out (default: 10)",
			},
			"read_timeout": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Default:     10,
				Description: "Number of seconds before response times out (default: 10). Note: kept for backwards compatibility, currently unused.",
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
	if cfg.Host == "" {
		return logical.ErrorResponse("config parameter `host` cannot be empty"), nil
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
	if cfg.Secret == "" {
		return logical.ErrorResponse("config parameter `secret` cannot be empty"), nil
	}

	policies := make([]string, 0)
	unregisteredUserPoliciesRaw, ok := d.GetOk("unregistered_user_policies")
	if ok {
		unregisteredUserPoliciesStr := unregisteredUserPoliciesRaw.(string)
		if strings.TrimSpace(unregisteredUserPoliciesStr) != "" {
			policies = strings.Split(unregisteredUserPoliciesStr, ",")
			for _, policy := range policies {
				if policy == "root" {
					return logical.ErrorResponse("root policy cannot be granted by an authentication backend"), nil
				}
			}
		}
		cfg.UnregisteredUserPolicies = policies
	} else if req.Operation == logical.CreateOperation {
		cfg.UnregisteredUserPolicies = policies
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
	Host                     string   `json:"host" structs:"host" mapstructure:"host"`
	Port                     int      `json:"port" structs:"port" mapstructure:"port"`
	Secret                   string   `json:"secret" structs:"secret" mapstructure:"secret"`
	UnregisteredUserPolicies []string `json:"unregistered_user_policies" structs:"unregistered_user_policies" mapstructure:"unregistered_user_policies"`
	DialTimeout              int      `json:"dial_timeout" structs:"dial_timeout" mapstructure:"dial_timeout"`
	ReadTimeout              int      `json:"read_timeout" structs:"read_timeout" mapstructure:"read_timeout"`
	NasPort                  int      `json:"nas_port" structs:"nas_port" mapstructure:"nas_port"`
}

const pathConfigHelpSyn = `
Configure the RADIUS server to connect to, along with its options.
`

const pathConfigHelpDesc = `
This endpoint allows you to configure the RADIUS server to connect to and its
configuration options.
`
