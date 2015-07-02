package pki

import (
	"fmt"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `roles/(?P<name>\w[\w-]+\w)`,
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role",
			},

			"lease": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "",
				Description: `The lease length if no specific lease length is
requested. The lease length controls the expiration
of certificates issued by this backend. Defaults to
the value of lease_max.`,
			},

			"lease_max": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: "The maximum allowed lease length",
			},

			"allow_localhost": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: true,
				Description: `Whether to allow "localhost" as a valid common
name in a request`,
			},

			"allowed_base_domain": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "",
				Description: `If set, clients can request certificates for
subdomains directly beneath this base domain, including
the wildcard subdomain. See the documentation for more
information.`,
			},

			"allow_token_displayname": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: false,
				Description: `If set, clients can request certificates for
matching the value of the Display Name on the requesting
token. See the documentation for more information.`,
			},

			"allow_subdomains": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: false,
				Description: `If set, clients can request certificates for
subdomains of the CNs allowed by the other role options,
including wildcard subdomains. See the documentation for
more information.`,
			},

			"allow_any_name": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: false,
				Description: `If set, clients can request certificates for
any CN they like. See the documentation for more
information.`,
			},

			"allow_ip_sans": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: true,
				Description: `If set, IP Subject Alternative Names are allowed.
Any valid IP is accepted.`,
			},

			"server_flag": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: true,
				Description: `If set, certificates are flagged for server use.
Defaults to true.`,
			},

			"client_flag": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: true,
				Description: `If set, certificates are flagged for client use.
Defaults to true.`,
			},

			"code_signing_flag": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: false,
				Description: `If set, certificates are flagged for code signing
use. Defaults to false.`,
			},

			"key_type": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "rsa",
				Description: `The type of key to use; defaults to RSA. "rsa"
and "ec" are the only valid values.`,
			},

			"key_bits": &framework.FieldSchema{
				Type:    framework.TypeInt,
				Default: 2048,
				Description: `The number of bits to use. You will almost
certainly want to change this if you adjust
the key_type.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathRoleRead,
			logical.WriteOperation:  b.pathRoleCreate,
			logical.DeleteOperation: b.pathRoleDelete,
		},

		HelpSynopsis:    pathRoleHelpSyn,
		HelpDescription: pathRoleHelpDesc,
	}
}

func (b *backend) getRole(s logical.Storage, n string) (*roleEntry, error) {
	entry, err := s.Get("role/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result roleEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathRoleDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete("role/" + data.Get("name").(string))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathRoleRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	role, err := b.getRole(req.Storage, data.Get("name").(string))
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

func (b *backend) pathRoleCreate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)

	entry := &roleEntry{
		LeaseMax:              data.Get("lease_max").(string),
		Lease:                 data.Get("lease").(string),
		AllowLocalhost:        data.Get("allow_localhost").(bool),
		AllowedBaseDomain:     data.Get("allowed_base_domain").(string),
		AllowTokenDisplayName: data.Get("allow_token_displayname").(bool),
		AllowSubdomains:       data.Get("allow_subdomains").(bool),
		AllowAnyName:          data.Get("allow_any_name").(bool),
		AllowIPSANs:           data.Get("allow_ip_sans").(bool),
		ServerFlag:            data.Get("server_flag").(bool),
		ClientFlag:            data.Get("client_flag").(bool),
		CodeSigningFlag:       data.Get("code_signing_flag").(bool),
		KeyType:               data.Get("key_type").(string),
		KeyBits:               data.Get("key_bits").(int),
	}

	if len(entry.LeaseMax) == 0 {
		return logical.ErrorResponse("\"lease_max\" value must be supplied"), nil
	}

	leaseMax, err := time.ParseDuration(entry.LeaseMax)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Invalid lease: %s", err)), nil
	}

	switch len(entry.Lease) {
	case 0:
		entry.Lease = entry.LeaseMax
	default:
		lease, err := time.ParseDuration(entry.Lease)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
				"Invalid lease: %s", err)), nil
		}
		if lease > leaseMax {
			return logical.ErrorResponse("\"lease\" value must be less than \"lease_max\" value"), nil
		}
	}

	if len(entry.KeyType) == 0 {
		entry.KeyType = "rsa"
	}
	if entry.KeyBits == 0 {
		entry.KeyBits = 2048
	}

	switch entry.KeyType {
	case "rsa":
	case "ec":
		switch entry.KeyBits {
		case 224:
		case 256:
		case 384:
		case 521:
		default:
			return logical.ErrorResponse(fmt.Sprintf("Unsupported bit length for EC key: %d", entry.KeyBits)), nil
		}
	default:
		return logical.ErrorResponse(fmt.Sprintf("Unknown key type %s", entry.KeyType)), nil
	}

	// Store it
	jsonEntry, err := logical.StorageEntryJSON("role/"+name, entry)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(jsonEntry); err != nil {
		return nil, err
	}

	return nil, nil
}

type roleEntry struct {
	LeaseMax              string `json:"lease_max" structs:"lease_max" mapstructure:"lease_max"`
	Lease                 string `json:"lease" structs:"lease" mapstructure:"lease"`
	AllowLocalhost        bool   `json:"allow_localhost" structs:"allow_localhost" mapstructure:"allow_localhost"`
	AllowedBaseDomain     string `json:"allowed_base_domain" structs:"allowed_base_domain" mapstructure:"allowed_base_domain"`
	AllowTokenDisplayName bool   `json:"allow_token_displayname" structs:"allow_token_displayname" mapstructure:"allow_token_displayname"`
	AllowSubdomains       bool   `json:"allow_subdomains" structs:"allow_subdomains" mapstructure:"allow_subdomains"`
	AllowAnyName          bool   `json:"allow_any_name" structs:"allow_any_name" mapstructure:"allow_any_name"`
	AllowIPSANs           bool   `json:"allow_ip_sans" structs:"allow_ip_sans" mapstructure:"allow_ip_sans"`
	ServerFlag            bool   `json:"server_flag" structs:"server_flag" mapstructure:"server_flag"`
	ClientFlag            bool   `json:"client_flag" structs:"client_flag" mapstructure:"client_flag"`
	CodeSigningFlag       bool   `json:"code_signing_flag" structs:"code_signing_flag" mapstructure:"code_signing_flag"`
	KeyType               string `json:"key_type" structs:"key_type" mapstructure:"key_type"`
	KeyBits               int    `json:"key_bits" structs:"key_bits" mapstructure:"key_bits"`
}

const pathRoleHelpSyn = `
Manage the roles that can be created with this backend.
`

const pathRoleHelpDesc = `
This path lets you manage the roles that can be created with this backend.
`
