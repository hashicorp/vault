package pki

import (
	"crypto/x509"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathListRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},

		HelpSynopsis:    pathListRolesHelpSyn,
		HelpDescription: pathListRolesHelpDesc,
	}
}

func pathRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role",
			},

			"ttl": &framework.FieldSchema{
				Type:    framework.TypeDurationSecond,
				Default: "",
				Description: `The lease duration if no specific lease duration is
requested. The lease duration controls the expiration
of certificates issued by this backend. Defaults to
the value of max_ttl.`,
			},

			"max_ttl": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "",
				Description: "The maximum allowed lease duration",
			},

			"allow_localhost": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: true,
				Description: `Whether to allow "localhost" as a valid common
name in a request`,
			},

			"allowed_domains": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "",
				Description: `If set, clients can request certificates for
subdomains directly beneath these domains, including
the wildcard subdomains. See the documentation for more
information. This parameter accepts a comma-separated list
of domains.`,
			},

			"allow_bare_domains": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: false,
				Description: `If set, clients can request certificates
for the base domains themselves, e.g. "example.com".
This is a separate option as in some cases this can
be considered a security threat.`,
			},

			"allow_subdomains": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: false,
				Description: `If set, clients can request certificates for
subdomains of the CNs allowed by the other role options,
including wildcard subdomains. See the documentation for
more information.`,
			},

			"allow_glob_domains": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: false,
				Description: `If set, domains specified in "allowed_domains"
can include glob patterns, e.g. "ftp*.example.com". See
the documentation for more information.`,
			},

			"allow_any_name": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: false,
				Description: `If set, clients can request certificates for
any CN they like. See the documentation for more
information.`,
			},

			"enforce_hostnames": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: true,
				Description: `If set, only valid host names are allowed for
CN and SANs. Defaults to true.`,
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
				Description: `If set, certificates are flagged for server auth use.
Defaults to true.`,
			},

			"client_flag": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: true,
				Description: `If set, certificates are flagged for client auth use.
Defaults to true.`,
			},

			"code_signing_flag": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: false,
				Description: `If set, certificates are flagged for code signing
use. Defaults to false.`,
			},

			"email_protection_flag": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: false,
				Description: `If set, certificates are flagged for email
protection use. Defaults to false.`,
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

			"key_usage": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "DigitalSignature,KeyAgreement,KeyEncipherment",
				Description: `A comma-separated set of key usages (not extended
key usages). Valid values can be found at
https://golang.org/pkg/crypto/x509/#KeyUsage
-- simply drop the "KeyUsage" part of the name.
To remove all key usages from being set, set
this value to an empty string.`,
			},

			"use_csr_common_name": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: true,
				Description: `If set, when used with a signing profile,
the common name in the CSR will be used. This
does *not* include any requested Subject Alternative
Names. Defaults to true.`,
			},

			"use_csr_sans": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: true,
				Description: `If set, when used with a signing profile,
the SANs in the CSR will be used. This does *not*
include the Common Name (cn). Defaults to true.`,
			},

			"ou": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "",
				Description: `If set, the OU (OrganizationalUnit) will be set to
this value in certificates issued by this role.`,
			},

			"organization": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "",
				Description: `If set, the O (Organization) will be set to
this value in certificates issued by this role.`,
			},

			"generate_lease": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: false,
				Description: `
If set, certificates issued/signed against this role will have Vault leases
attached to them. Defaults to "false". Certificates can be added to the CRL by
"vault revoke <lease_id>" when certificates are associated with leases.  It can
also be done using the "pki/revoke" endpoint. However, when lease generation is
disabled, invoking "pki/revoke" would be the only way to add the certificates
to the CRL.  When large number of certificates are generated with long
lifetimes, it is recommended that lease generation be disabled, as large amount of
leases adversely affect the startup time of Vault.`,
			},
			"no_store": &framework.FieldSchema{
				Type:    framework.TypeBool,
				Default: false,
				Description: `
If set, certificates issued/signed against this role will not be stored in the
in the storage backend. This can improve performance when issuing large numbers
of certificates. However, certificates issued in this way cannot be enumerated
or revoked, so this option is recommended only for certificates that are
non-sensitive, or extremely short-lived. This option implies a value of "false"
for "generate_lease".`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathRoleRead,
			logical.UpdateOperation: b.pathRoleCreate,
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

	// Migrate existing saved entries and save back if changed
	modified := false
	if len(result.TTL) == 0 && len(result.Lease) != 0 {
		result.TTL = result.Lease
		result.Lease = ""
		modified = true
	}
	if len(result.MaxTTL) == 0 && len(result.LeaseMax) != 0 {
		result.MaxTTL = result.LeaseMax
		result.LeaseMax = ""
		modified = true
	}
	if result.AllowBaseDomain {
		result.AllowBaseDomain = false
		result.AllowBareDomains = true
		modified = true
	}
	if result.AllowedBaseDomain != "" {
		found := false
		allowedDomains := strings.Split(result.AllowedDomains, ",")
		if len(allowedDomains) != 0 {
			for _, v := range allowedDomains {
				if v == result.AllowedBaseDomain {
					found = true
					break
				}
			}
		}
		if !found {
			if result.AllowedDomains == "" {
				result.AllowedDomains = result.AllowedBaseDomain
			} else {
				result.AllowedDomains += "," + result.AllowedBaseDomain
			}
		}
		result.AllowedBaseDomain = ""
		modified = true
	}

	// Upgrade generate_lease in role
	if result.GenerateLease == nil {
		// All the new roles will have GenerateLease always set to a value. A
		// nil value indicates that this role needs an upgrade. Set it to
		// `true` to not alter its current behavior.
		result.GenerateLease = new(bool)
		*result.GenerateLease = true
		modified = true
	}

	if modified {
		jsonEntry, err := logical.StorageEntryJSON("role/"+n, &result)
		if err != nil {
			return nil, err
		}
		if err := s.Put(jsonEntry); err != nil {
			return nil, err
		}
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
	roleName := data.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role name"), nil
	}

	role, err := b.getRole(req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	hasMax := true
	if len(role.MaxTTL) == 0 {
		role.MaxTTL = "(system default)"
		hasMax = false
	}
	if len(role.TTL) == 0 {
		if hasMax {
			role.TTL = "(system default, capped to role max)"
		} else {
			role.TTL = "(system default)"
		}
	}

	resp := &logical.Response{
		Data: structs.New(role).Map(),
	}

	if resp.Data == nil {
		return nil, fmt.Errorf("error converting role data to response")
	}

	// These values are deprecated and the entries are migrated on read
	delete(resp.Data, "lease")
	delete(resp.Data, "lease_max")
	delete(resp.Data, "allowed_base_domain")
	delete(resp.Data, "allow_base_domain")
	delete(resp.Data, "AllowExpirationPastCA")

	return resp, nil
}

func (b *backend) pathRoleList(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List("role/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func (b *backend) pathRoleCreate(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var err error
	name := data.Get("name").(string)

	entry := &roleEntry{
		MaxTTL:              data.Get("max_ttl").(string),
		TTL:                 (time.Duration(data.Get("ttl").(int)) * time.Second).String(),
		AllowLocalhost:      data.Get("allow_localhost").(bool),
		AllowedDomains:      data.Get("allowed_domains").(string),
		AllowBareDomains:    data.Get("allow_bare_domains").(bool),
		AllowSubdomains:     data.Get("allow_subdomains").(bool),
		AllowGlobDomains:    data.Get("allow_glob_domains").(bool),
		AllowAnyName:        data.Get("allow_any_name").(bool),
		EnforceHostnames:    data.Get("enforce_hostnames").(bool),
		AllowIPSANs:         data.Get("allow_ip_sans").(bool),
		ServerFlag:          data.Get("server_flag").(bool),
		ClientFlag:          data.Get("client_flag").(bool),
		CodeSigningFlag:     data.Get("code_signing_flag").(bool),
		EmailProtectionFlag: data.Get("email_protection_flag").(bool),
		KeyType:             data.Get("key_type").(string),
		KeyBits:             data.Get("key_bits").(int),
		UseCSRCommonName:    data.Get("use_csr_common_name").(bool),
		UseCSRSANs:          data.Get("use_csr_sans").(bool),
		KeyUsage:            data.Get("key_usage").(string),
		OU:                  data.Get("ou").(string),
		Organization:        data.Get("organization").(string),
		GenerateLease:       new(bool),
		NoStore:             data.Get("no_store").(bool),
	}

	// no_store implies generate_lease := false
	if entry.NoStore {
		*entry.GenerateLease = false
	} else {
		*entry.GenerateLease = data.Get("generate_lease").(bool)
	}

	if entry.KeyType == "rsa" && entry.KeyBits < 2048 {
		return logical.ErrorResponse("RSA keys < 2048 bits are unsafe and not supported"), nil
	}

	var maxTTL time.Duration
	maxSystemTTL := b.System().MaxLeaseTTL()
	if len(entry.MaxTTL) == 0 {
		maxTTL = maxSystemTTL
	} else {
		maxTTL, err = parseutil.ParseDurationSecond(entry.MaxTTL)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
				"Invalid max ttl: %s", err)), nil
		}
	}
	if maxTTL > maxSystemTTL {
		return logical.ErrorResponse("Requested max TTL is higher than backend maximum"), nil
	}

	ttl := b.System().DefaultLeaseTTL()
	if len(entry.TTL) != 0 {
		ttl, err = parseutil.ParseDurationSecond(entry.TTL)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
				"Invalid ttl: %s", err)), nil
		}
	}
	if ttl > maxTTL {
		// If they are using the system default, cap it to the role max;
		// if it was specified on the command line, make it an error
		if len(entry.TTL) == 0 {
			ttl = maxTTL
		} else {
			return logical.ErrorResponse(
				`"ttl" value must be less than "max_ttl" and/or backend default max lease TTL value`,
			), nil
		}
	}

	// Persist clamped TTLs
	entry.TTL = ttl.String()
	entry.MaxTTL = maxTTL.String()

	if errResp := validateKeyTypeLength(entry.KeyType, entry.KeyBits); errResp != nil {
		return errResp, nil
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

func parseKeyUsages(input string) int {
	var parsedKeyUsages x509.KeyUsage
	splitKeyUsage := strings.Split(input, ",")
	for _, k := range splitKeyUsage {
		switch strings.ToLower(strings.TrimSpace(k)) {
		case "digitalsignature":
			parsedKeyUsages |= x509.KeyUsageDigitalSignature
		case "contentcommitment":
			parsedKeyUsages |= x509.KeyUsageContentCommitment
		case "keyencipherment":
			parsedKeyUsages |= x509.KeyUsageKeyEncipherment
		case "dataencipherment":
			parsedKeyUsages |= x509.KeyUsageDataEncipherment
		case "keyagreement":
			parsedKeyUsages |= x509.KeyUsageKeyAgreement
		case "certsign":
			parsedKeyUsages |= x509.KeyUsageCertSign
		case "crlsign":
			parsedKeyUsages |= x509.KeyUsageCRLSign
		case "encipheronly":
			parsedKeyUsages |= x509.KeyUsageEncipherOnly
		case "decipheronly":
			parsedKeyUsages |= x509.KeyUsageDecipherOnly
		}
	}

	return int(parsedKeyUsages)
}

type roleEntry struct {
	LeaseMax              string `json:"lease_max" structs:"lease_max" mapstructure:"lease_max"`
	Lease                 string `json:"lease" structs:"lease" mapstructure:"lease"`
	MaxTTL                string `json:"max_ttl" structs:"max_ttl" mapstructure:"max_ttl"`
	TTL                   string `json:"ttl" structs:"ttl" mapstructure:"ttl"`
	AllowLocalhost        bool   `json:"allow_localhost" structs:"allow_localhost" mapstructure:"allow_localhost"`
	AllowedBaseDomain     string `json:"allowed_base_domain" structs:"allowed_base_domain" mapstructure:"allowed_base_domain"`
	AllowedDomains        string `json:"allowed_domains" structs:"allowed_domains" mapstructure:"allowed_domains"`
	AllowBaseDomain       bool   `json:"allow_base_domain" structs:"allow_base_domain" mapstructure:"allow_base_domain"`
	AllowBareDomains      bool   `json:"allow_bare_domains" structs:"allow_bare_domains" mapstructure:"allow_bare_domains"`
	AllowTokenDisplayName bool   `json:"allow_token_displayname" structs:"allow_token_displayname" mapstructure:"allow_token_displayname"`
	AllowSubdomains       bool   `json:"allow_subdomains" structs:"allow_subdomains" mapstructure:"allow_subdomains"`
	AllowGlobDomains      bool   `json:"allow_glob_domains" structs:"allow_glob_domains" mapstructure:"allow_glob_domains"`
	AllowAnyName          bool   `json:"allow_any_name" structs:"allow_any_name" mapstructure:"allow_any_name"`
	EnforceHostnames      bool   `json:"enforce_hostnames" structs:"enforce_hostnames" mapstructure:"enforce_hostnames"`
	AllowIPSANs           bool   `json:"allow_ip_sans" structs:"allow_ip_sans" mapstructure:"allow_ip_sans"`
	ServerFlag            bool   `json:"server_flag" structs:"server_flag" mapstructure:"server_flag"`
	ClientFlag            bool   `json:"client_flag" structs:"client_flag" mapstructure:"client_flag"`
	CodeSigningFlag       bool   `json:"code_signing_flag" structs:"code_signing_flag" mapstructure:"code_signing_flag"`
	EmailProtectionFlag   bool   `json:"email_protection_flag" structs:"email_protection_flag" mapstructure:"email_protection_flag"`
	UseCSRCommonName      bool   `json:"use_csr_common_name" structs:"use_csr_common_name" mapstructure:"use_csr_common_name"`
	UseCSRSANs            bool   `json:"use_csr_sans" structs:"use_csr_sans" mapstructure:"use_csr_sans"`
	KeyType               string `json:"key_type" structs:"key_type" mapstructure:"key_type"`
	KeyBits               int    `json:"key_bits" structs:"key_bits" mapstructure:"key_bits"`
	MaxPathLength         *int   `json:",omitempty" structs:"max_path_length,omitempty" mapstructure:"max_path_length"`
	KeyUsage              string `json:"key_usage" structs:"key_usage" mapstructure:"key_usage"`
	OU                    string `json:"ou" structs:"ou" mapstructure:"ou"`
	Organization          string `json:"organization" structs:"organization" mapstructure:"organization"`
	GenerateLease         *bool  `json:"generate_lease,omitempty" structs:"generate_lease,omitempty"`
	NoStore               bool   `json:"no_store" structs:"no_store" mapstructure:"no_store"`

	// Used internally for signing intermediates
	AllowExpirationPastCA bool
}

const pathListRolesHelpSyn = `List the existing roles in this backend`

const pathListRolesHelpDesc = `Roles will be listed by the role name.`

const pathRoleHelpSyn = `Manage the roles that can be created with this backend.`

const pathRoleHelpDesc = `This path lets you manage the roles that can be created with this backend.`
