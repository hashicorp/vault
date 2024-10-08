// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ssh

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/cidrutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/crypto/ssh"
)

const (
	// KeyTypeOTP is an key of type OTP
	KeyTypeOTP = "otp"
	// KeyTypeDynamic is dynamic key type; removed.
	KeyTypeDynamic = "dynamic"
	// KeyTypeCA is an key of type CA
	KeyTypeCA = "ca"

	// DefaultAlgorithmSigner is the default RSA signing algorithm
	DefaultAlgorithmSigner = "default"

	// Present version of the sshRole struct; when adding a new field or are
	// needing to perform a migration, increment this struct and read the note
	// in checkUpgrade(...).
	roleEntryVersion = 3
)

// Structure that represents a role in SSH backend. This is a common role structure
// for both OTP and CA roles. Not all the fields are mandatory for both type.
// Some are applicable for one and not for other. It doesn't matter.
type sshRole struct {
	KeyType                    string            `mapstructure:"key_type" json:"key_type"`
	DefaultUser                string            `mapstructure:"default_user" json:"default_user"`
	DefaultUserTemplate        bool              `mapstructure:"default_user_template" json:"default_user_template"`
	CIDRList                   string            `mapstructure:"cidr_list" json:"cidr_list"`
	ExcludeCIDRList            string            `mapstructure:"exclude_cidr_list" json:"exclude_cidr_list"`
	Port                       int               `mapstructure:"port" json:"port"`
	AllowedUsers               string            `mapstructure:"allowed_users" json:"allowed_users"`
	AllowedUsersTemplate       bool              `mapstructure:"allowed_users_template" json:"allowed_users_template"`
	AllowedDomains             string            `mapstructure:"allowed_domains" json:"allowed_domains"`
	AllowedDomainsTemplate     bool              `mapstructure:"allowed_domains_template" json:"allowed_domains_template"`
	MaxTTL                     string            `mapstructure:"max_ttl" json:"max_ttl"`
	TTL                        string            `mapstructure:"ttl" json:"ttl"`
	DefaultCriticalOptions     map[string]string `mapstructure:"default_critical_options" json:"default_critical_options"`
	DefaultExtensions          map[string]string `mapstructure:"default_extensions" json:"default_extensions"`
	DefaultExtensionsTemplate  bool              `mapstructure:"default_extensions_template" json:"default_extensions_template"`
	AllowedCriticalOptions     string            `mapstructure:"allowed_critical_options" json:"allowed_critical_options"`
	AllowedExtensions          string            `mapstructure:"allowed_extensions" json:"allowed_extensions"`
	AllowUserCertificates      bool              `mapstructure:"allow_user_certificates" json:"allow_user_certificates"`
	AllowHostCertificates      bool              `mapstructure:"allow_host_certificates" json:"allow_host_certificates"`
	AllowBareDomains           bool              `mapstructure:"allow_bare_domains" json:"allow_bare_domains"`
	AllowSubdomains            bool              `mapstructure:"allow_subdomains" json:"allow_subdomains"`
	AllowUserKeyIDs            bool              `mapstructure:"allow_user_key_ids" json:"allow_user_key_ids"`
	KeyIDFormat                string            `mapstructure:"key_id_format" json:"key_id_format"`
	OldAllowedUserKeyLengths   map[string]int    `mapstructure:"allowed_user_key_lengths" json:"allowed_user_key_lengths,omitempty"`
	AllowedUserKeyTypesLengths map[string][]int  `mapstructure:"allowed_user_key_types_lengths" json:"allowed_user_key_types_lengths"`
	AlgorithmSigner            string            `mapstructure:"algorithm_signer" json:"algorithm_signer"`
	Version                    int               `mapstructure:"role_version" json:"role_version"`
	NotBeforeDuration          time.Duration     `mapstructure:"not_before_duration" json:"not_before_duration"`
	AllowEmptyPrincipals       bool              `mapstructure:"allow_empty_principals" json:"allow_empty_principals"`
}

func pathListRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixSSH,
			OperationSuffix: "roles",
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},

		HelpSynopsis:    pathRoleHelpSyn,
		HelpDescription: pathRoleHelpDesc,
	}
}

func pathRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameWithAtRegex("role"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixSSH,
			OperationSuffix: "role",
		},

		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type: framework.TypeString,
				Description: `
				[Required for all types]
				Name of the role being created.`,
			},
			"default_user": {
				Type: framework.TypeString,
				Description: `
				[Required for OTP type] [Optional for CA type]
				Default username for which a credential will be generated.
				When the endpoint 'creds/' is used without a username, this
				value will be used as default username.`,
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Default Username",
				},
			},
			"default_user_template": {
				Type: framework.TypeBool,
				Description: `
				[Not applicable for OTP type] [Optional for CA type]
				If set, Default user can be specified using identity template policies.
				Non-templated users are also permitted.
				`,
				Default: false,
			},
			"cidr_list": {
				Type: framework.TypeString,
				Description: `
				[Optional for OTP type] [Not applicable for CA type]
				Comma separated list of CIDR blocks for which the role is applicable for.
				CIDR blocks can belong to more than one role.`,
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "CIDR List",
				},
			},
			"exclude_cidr_list": {
				Type: framework.TypeString,
				Description: `
				[Optional for OTP type] [Not applicable for CA type]
				Comma separated list of CIDR blocks. IP addresses belonging to these blocks are not
				accepted by the role. This is particularly useful when big CIDR blocks are being used
				by the role and certain parts of it needs to be kept out.`,
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Exclude CIDR List",
				},
			},
			"port": {
				Type: framework.TypeInt,
				Description: `
				[Optional for OTP type] [Not applicable for CA type]
				Port number for SSH connection. Default is '22'. Port number does not
				play any role in creation of OTP. For 'otp' type, this is just a way
				to inform client about the port number to use. Port number will be
				returned to client by Vault server along with OTP.`,
				DisplayAttrs: &framework.DisplayAttributes{
					Value: 22,
				},
			},
			"key_type": {
				Type: framework.TypeString,
				Description: `
				[Required for all types]
				Type of key used to login to hosts. It can be either 'otp' or 'ca'.
				'otp' type requires agent to be installed in remote hosts.`,
				AllowedValues: []interface{}{"otp", "ca"},
				DisplayAttrs: &framework.DisplayAttributes{
					Value: "ca",
				},
			},
			"allowed_users": {
				Type: framework.TypeString,
				Description: `
				[Optional for all types] [Works differently for CA type]
				If this option is not specified, or is '*', client can request a
				credential for any valid user at the remote host, including the
				admin user. If only certain usernames are to be allowed, then
				this list enforces it. If this field is set, then credentials
				can only be created for default_user and usernames present in
				this list. Setting this option will enable all the users with
				access to this role to fetch credentials for all other usernames
				in this list. Use with caution. N.B.: with the CA type, an empty
				list means that no users are allowed; explicitly specify '*' to
				allow any user.
				`,
			},
			"allowed_users_template": {
				Type: framework.TypeBool,
				Description: `
				[Not applicable for OTP type] [Optional for CA type]
				If set, Allowed users can be specified using identity template policies.
				Non-templated users are also permitted.
				`,
				Default: false,
			},
			"allowed_domains": {
				Type: framework.TypeString,
				Description: `
				[Not applicable for OTP type] [Optional for CA type]
				If this option is not specified, client can request for a signed certificate for any
				valid host. If only certain domains are allowed, then this list enforces it.
				`,
			},
			"allowed_domains_template": {
				Type: framework.TypeBool,
				Description: `
				[Not applicable for OTP type] [Optional for CA type]
				If set, Allowed domains can be specified using identity template policies.
				Non-templated domains are also permitted.
				`,
				Default: false,
			},
			"ttl": {
				Type: framework.TypeDurationSecond,
				Description: `
				[Not applicable for OTP type] [Optional for CA type]
				The lease duration if no specific lease duration is
				requested. The lease duration controls the expiration
				of certificates issued by this backend. Defaults to
				the value of max_ttl.`,
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "TTL",
				},
			},
			"max_ttl": {
				Type: framework.TypeDurationSecond,
				Description: `
				[Not applicable for OTP type] [Optional for CA type]
				The maximum allowed lease duration
				`,
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Max TTL",
				},
			},
			"allowed_critical_options": {
				Type: framework.TypeString,
				Description: `
				[Not applicable for OTP type] [Optional for CA type]
				A comma-separated list of critical options that certificates can have when signed.
 				To allow any critical options, set this to an empty string.
				 `,
			},
			"allowed_extensions": {
				Type: framework.TypeString,
				Description: `
				[Not applicable for OTP type] [Optional for CA type]
				A comma-separated list of extensions that certificates can have when signed.
				An empty list means that no extension overrides are allowed by an end-user; explicitly
				specify '*' to allow any extensions to be set.
				`,
			},
			"default_critical_options": {
				Type: framework.TypeMap,
				Description: `
				[Not applicable for OTP type] [Optional for CA type]
				Critical options certificates should
				have if none are provided when signing. This field takes in key
				value pairs in JSON format.  Note that these are not restricted
				by "allowed_critical_options". Defaults to none.
				`,
			},
			"default_extensions": {
				Type: framework.TypeMap,
				Description: `
				[Not applicable for OTP type] [Optional for CA type]
				Extensions certificates should have if
				none are provided when signing. This field takes in key value
				pairs in JSON format. Note that these are not restricted by
				"allowed_extensions". Defaults to none.
				`,
			},
			"default_extensions_template": {
				Type: framework.TypeBool,
				Description: `
				[Not applicable for OTP type] [Optional for CA type]
				If set, Default extension values can be specified using identity template policies.
				Non-templated extension values are also permitted.
				`,
				Default: false,
			},
			"allow_user_certificates": {
				Type: framework.TypeBool,
				Description: `
				[Not applicable for OTP type] [Optional for CA type]
				If set, certificates are allowed to be signed for use as a 'user'.
				`,
				Default: false,
			},
			"allow_host_certificates": {
				Type: framework.TypeBool,
				Description: `
				[Not applicable for OTP type] [Optional for CA type]
				If set, certificates are allowed to be signed for use as a 'host'.
				`,
				Default: false,
			},
			"allow_bare_domains": {
				Type: framework.TypeBool,
				Description: `
				[Not applicable for OTP type] [Optional for CA type]
				If set, host certificates that are requested are allowed to use the base domains listed in
				"allowed_domains", e.g. "example.com".
				This is a separate option as in some cases this can be considered a security threat.
				`,
			},
			"allow_subdomains": {
				Type: framework.TypeBool,
				Description: `
				[Not applicable for OTP type] [Optional for CA type]
				If set, host certificates that are requested are allowed to use subdomains of those listed in "allowed_domains".
				`,
			},
			"allow_user_key_ids": {
				Type: framework.TypeBool,
				Description: `
				[Not applicable for OTP type] [Optional for CA type]
				If true, users can override the key ID for a signed certificate with the "key_id" field.
				When false, the key ID will always be the token display name.
				The key ID is logged by the SSH server and can be useful for auditing.
				`,
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Allow User Key IDs",
				},
			},
			"key_id_format": {
				Type: framework.TypeString,
				Description: `
				[Not applicable for OTP type] [Optional for CA type]
				When supplied, this value specifies a custom format for the key id of a signed certificate.
				The following variables are available for use: '{{token_display_name}}' - The display name of
				the token used to make the request. '{{role_name}}' - The name of the role signing the request.
				'{{public_key_hash}}' - A SHA256 checksum of the public key that is being signed.
				`,
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Key ID Format",
				},
			},
			"allowed_user_key_lengths": {
				Type: framework.TypeMap,
				Description: `
				[Not applicable for OTP type] [Optional for CA type]
				If set, allows the enforcement of key types and minimum key sizes to be signed.
				`,
			},
			"algorithm_signer": {
				Type: framework.TypeString,
				Description: `
				[Not applicable for OTP type] [Optional for CA type]
				When supplied, this value specifies a signing algorithm for the key. Possible values:
				ssh-rsa, rsa-sha2-256, rsa-sha2-512, default, or the empty string.
				`,
				AllowedValues: []interface{}{"", DefaultAlgorithmSigner, ssh.SigAlgoRSA, ssh.SigAlgoRSASHA2256, ssh.SigAlgoRSASHA2512},
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Signing Algorithm",
				},
			},
			"not_before_duration": {
				Type:    framework.TypeDurationSecond,
				Default: 30,
				Description: `
				[Not applicable for OTP type] [Optional for CA type]
   				The duration that the SSH certificate should be backdated by at issuance.`,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:  "Not before duration",
					Value: 30,
				},
			},
			"allow_empty_principals": {
				Type:        framework.TypeBool,
				Description: `Whether to allow issuing certificates with no valid principals (meaning any valid principal).  Exists for backwards compatibility only, the default of false is highly recommended.`,
				Default:     false,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathRoleRead,
			logical.UpdateOperation: b.pathRoleWrite,
			logical.DeleteOperation: b.pathRoleDelete,
		},

		HelpSynopsis:    pathRoleHelpSyn,
		HelpDescription: pathRoleHelpDesc,
	}
}

func (b *backend) pathRoleWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleName := d.Get("role").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role name"), nil
	}

	// Allowed users is an optional field, applicable for both OTP and CA types.
	allowedUsers := d.Get("allowed_users").(string)

	// Validate the CIDR blocks
	cidrList := d.Get("cidr_list").(string)
	if cidrList != "" {
		valid, err := cidrutil.ValidateCIDRListString(cidrList, ",")
		if err != nil {
			return nil, fmt.Errorf("failed to validate cidr_list: %w", err)
		}
		if !valid {
			return logical.ErrorResponse("failed to validate cidr_list"), nil
		}
	}

	// Validate the excluded CIDR blocks
	excludeCidrList := d.Get("exclude_cidr_list").(string)
	if excludeCidrList != "" {
		valid, err := cidrutil.ValidateCIDRListString(excludeCidrList, ",")
		if err != nil {
			return nil, fmt.Errorf("failed to validate exclude_cidr_list entry: %w", err)
		}
		if !valid {
			return logical.ErrorResponse(fmt.Sprintf("failed to validate exclude_cidr_list entry: %v", err)), nil
		}
	}

	port := d.Get("port").(int)
	if port == 0 {
		port = 22
	}

	keyType := d.Get("key_type").(string)
	if keyType == "" {
		return logical.ErrorResponse("missing key type"), nil
	}
	keyType = strings.ToLower(keyType)

	var roleEntry sshRole
	if keyType == KeyTypeOTP {
		defaultUser := d.Get("default_user").(string)
		if defaultUser == "" {
			return logical.ErrorResponse("missing default user"), nil
		}

		// Below are the only fields used from the role structure for OTP type.
		roleEntry = sshRole{
			DefaultUser:     defaultUser,
			CIDRList:        cidrList,
			ExcludeCIDRList: excludeCidrList,
			KeyType:         KeyTypeOTP,
			Port:            port,
			AllowedUsers:    allowedUsers,
			Version:         roleEntryVersion,
		}
	} else if keyType == KeyTypeDynamic {
		return logical.ErrorResponse("dynamic key type roles are no longer supported"), nil
	} else if keyType == KeyTypeCA {
		algorithmSigner := DefaultAlgorithmSigner
		algorithmSignerRaw, ok := d.GetOk("algorithm_signer")
		if ok {
			algorithmSigner = algorithmSignerRaw.(string)
			switch algorithmSigner {
			case ssh.SigAlgoRSA, ssh.SigAlgoRSASHA2256, ssh.SigAlgoRSASHA2512:
			case "", DefaultAlgorithmSigner:
				// This case is valid, and the sign operation will use the signer's
				// default algorithm. Explicitly reset the value to the default value
				// rather than use the more vague implicit empty string.
				algorithmSigner = DefaultAlgorithmSigner
			default:
				return nil, fmt.Errorf("unknown algorithm signer %q", algorithmSigner)
			}
		}

		role, errorResponse := b.createCARole(allowedUsers, d.Get("default_user").(string), algorithmSigner, d)
		if errorResponse != nil {
			return errorResponse, nil
		}
		roleEntry = *role
	} else {
		return logical.ErrorResponse("invalid key type"), nil
	}

	entry, err := logical.StorageEntryJSON(fmt.Sprintf("roles/%s", roleName), roleEntry)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) createCARole(allowedUsers, defaultUser, signer string, data *framework.FieldData) (*sshRole, *logical.Response) {
	ttl := time.Duration(data.Get("ttl").(int)) * time.Second
	maxTTL := time.Duration(data.Get("max_ttl").(int)) * time.Second
	role := &sshRole{
		AllowedCriticalOptions:    data.Get("allowed_critical_options").(string),
		AllowedExtensions:         data.Get("allowed_extensions").(string),
		AllowUserCertificates:     data.Get("allow_user_certificates").(bool),
		AllowHostCertificates:     data.Get("allow_host_certificates").(bool),
		AllowedUsers:              allowedUsers,
		AllowedUsersTemplate:      data.Get("allowed_users_template").(bool),
		AllowedDomains:            data.Get("allowed_domains").(string),
		AllowedDomainsTemplate:    data.Get("allowed_domains_template").(bool),
		DefaultUser:               defaultUser,
		DefaultUserTemplate:       data.Get("default_user_template").(bool),
		AllowBareDomains:          data.Get("allow_bare_domains").(bool),
		AllowSubdomains:           data.Get("allow_subdomains").(bool),
		AllowUserKeyIDs:           data.Get("allow_user_key_ids").(bool),
		DefaultExtensionsTemplate: data.Get("default_extensions_template").(bool),
		KeyIDFormat:               data.Get("key_id_format").(string),
		KeyType:                   KeyTypeCA,
		AlgorithmSigner:           signer,
		Version:                   roleEntryVersion,
		NotBeforeDuration:         time.Duration(data.Get("not_before_duration").(int)) * time.Second,
		AllowEmptyPrincipals:      data.Get("allow_empty_principals").(bool),
	}

	if !role.AllowUserCertificates && !role.AllowHostCertificates {
		return nil, logical.ErrorResponse("Either 'allow_user_certificates' or 'allow_host_certificates' must be set to 'true'")
	}

	defaultCriticalOptions := convertMapToStringValue(data.Get("default_critical_options").(map[string]interface{}))
	defaultExtensions := convertMapToStringValue(data.Get("default_extensions").(map[string]interface{}))
	allowedUserKeyLengths, err := convertMapToIntSlice(data.Get("allowed_user_key_lengths").(map[string]interface{}))
	if err != nil {
		return nil, logical.ErrorResponse(fmt.Sprintf("error processing allowed_user_key_lengths: %s", err.Error()))
	}

	if ttl != 0 && maxTTL != 0 && ttl > maxTTL {
		return nil, logical.ErrorResponse(
			`"ttl" value must be less than "max_ttl" when both are specified`)
	}

	// Persist TTLs
	role.TTL = ttl.String()
	role.MaxTTL = maxTTL.String()
	role.DefaultCriticalOptions = defaultCriticalOptions
	role.DefaultExtensions = defaultExtensions
	role.AllowedUserKeyTypesLengths = allowedUserKeyLengths

	return role, nil
}

func (b *backend) getRole(ctx context.Context, s logical.Storage, n string) (*sshRole, error) {
	entry, err := s.Get(ctx, "roles/"+n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result sshRole
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	if err := b.checkUpgrade(ctx, s, n, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) checkUpgrade(ctx context.Context, s logical.Storage, n string, result *sshRole) error {
	modified := false

	// NOTE: When introducing a new migration, increment roleEntryVersion and
	// check if the version is less than the version this change was introduced
	// at and perform the change. At the end, set modified and update the
	// version to the version this migration was introduced at! Additionally,
	// add new migrations after all existing migrations.
	//
	// Otherwise, past or future migrations may not execute!
	if result.Version == roleEntryVersion {
		return nil
	}

	// Role version introduced at version 1, migrating OldAllowedUserKeyLengths
	// to the newer AllowedUserKeyTypesLengths field.
	if result.Version < 1 {
		// Only migrate if we have old data and no new data to avoid clobbering.
		//
		// This change introduced the first role version, value of 1.
		if len(result.OldAllowedUserKeyLengths) > 0 && len(result.AllowedUserKeyTypesLengths) == 0 {
			result.AllowedUserKeyTypesLengths = make(map[string][]int)
			for k, v := range result.OldAllowedUserKeyLengths {
				result.AllowedUserKeyTypesLengths[k] = []int{v}
			}
			result.OldAllowedUserKeyLengths = nil
		}

		result.Version = 1
		modified = true
	}

	// Role version 2 migrates an empty AlgorithmSigner to an explicit ssh-rsa
	// value WHEN the SSH CA key is a RSA key.
	if result.Version < 2 {
		// In order to perform the version 2 upgrade, we need knowledge of the
		// signing key type as we want to make ssh-rsa an explicitly notated
		// algorithm choice.
		var publicKey ssh.PublicKey
		publicKeyEntry, err := caKey(ctx, s, caPublicKey)
		if err != nil {
			b.Logger().Debug(fmt.Sprintf("failed to load public key entry while attempting to migrate: %v", err))
			goto SKIPVERSION2
		}
		if publicKeyEntry == nil || publicKeyEntry.Key == "" {
			b.Logger().Debug(fmt.Sprintf("got empty public key entry while attempting to migrate"))
			goto SKIPVERSION2
		}

		publicKey, err = parsePublicSSHKey(publicKeyEntry.Key)
		if err == nil {
			// Move an empty signing algorithm to an explicit ssh-rsa (SHA-1)
			// if this key is of type RSA. This isn't a secure default but
			// exists for backwards compatibility with existing versions of
			// Vault. By making it explicit, operators can see that this is
			// the value and move it to a newer algorithm in the future.
			if publicKey.Type() == ssh.KeyAlgoRSA && result.AlgorithmSigner == "" {
				result.AlgorithmSigner = ssh.SigAlgoRSA
			}

			result.Version = 2
			modified = true
		}

	SKIPVERSION2:
		err = nil
	}

	if result.Version < 3 {
		modified = true
		result.NotBeforeDuration = 30 * time.Second
		result.Version = 3
	}

	// Add new migrations just before here.
	//
	// Condition copied from PKI builtin.
	if modified && (b.System().LocalMount() || !b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary)) {
		jsonEntry, err := logical.StorageEntryJSON("roles/"+n, &result)
		if err != nil {
			return err
		}
		if err := s.Put(ctx, jsonEntry); err != nil {
			// Only perform upgrades on replication primary
			if !strings.Contains(err.Error(), logical.ErrReadOnly.Error()) {
				return err
			}
		}
	}

	return nil
}

// parseRole converts a sshRole object into its map[string]interface representation,
// with appropriate values for each KeyType. If the KeyType is invalid, it will return
// an error.
func (b *backend) parseRole(role *sshRole) (map[string]interface{}, error) {
	var result map[string]interface{}

	switch role.KeyType {
	case KeyTypeOTP:
		result = map[string]interface{}{
			"default_user":      role.DefaultUser,
			"cidr_list":         role.CIDRList,
			"exclude_cidr_list": role.ExcludeCIDRList,
			"key_type":          role.KeyType,
			"port":              role.Port,
			"allowed_users":     role.AllowedUsers,
		}
	case KeyTypeCA:
		ttl, err := parseutil.ParseDurationSecond(role.TTL)
		if err != nil {
			return nil, err
		}
		maxTTL, err := parseutil.ParseDurationSecond(role.MaxTTL)
		if err != nil {
			return nil, err
		}

		result = map[string]interface{}{
			"allowed_users":               role.AllowedUsers,
			"allowed_users_template":      role.AllowedUsersTemplate,
			"allowed_domains":             role.AllowedDomains,
			"allowed_domains_template":    role.AllowedDomainsTemplate,
			"default_user":                role.DefaultUser,
			"default_user_template":       role.DefaultUserTemplate,
			"ttl":                         int64(ttl.Seconds()),
			"max_ttl":                     int64(maxTTL.Seconds()),
			"allowed_critical_options":    role.AllowedCriticalOptions,
			"allowed_extensions":          role.AllowedExtensions,
			"allow_user_certificates":     role.AllowUserCertificates,
			"allow_host_certificates":     role.AllowHostCertificates,
			"allow_bare_domains":          role.AllowBareDomains,
			"allow_subdomains":            role.AllowSubdomains,
			"allow_user_key_ids":          role.AllowUserKeyIDs,
			"key_id_format":               role.KeyIDFormat,
			"key_type":                    role.KeyType,
			"default_critical_options":    role.DefaultCriticalOptions,
			"default_extensions":          role.DefaultExtensions,
			"default_extensions_template": role.DefaultExtensionsTemplate,
			"allowed_user_key_lengths":    role.AllowedUserKeyTypesLengths,
			"algorithm_signer":            role.AlgorithmSigner,
			"not_before_duration":         int64(role.NotBeforeDuration.Seconds()),
		}
	case KeyTypeDynamic:
		return nil, fmt.Errorf("dynamic key type roles are no longer supported")
	default:
		return nil, fmt.Errorf("invalid key type: %v", role.KeyType)
	}

	return result, nil
}

func (b *backend) pathRoleList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "roles/")
	if err != nil {
		return nil, err
	}

	keyInfo := map[string]interface{}{}
	for _, entry := range entries {
		role, err := b.getRole(ctx, req.Storage, entry)
		if err != nil {
			// On error, log warning and continue
			if b.Logger().IsWarn() {
				b.Logger().Warn("error getting role info", "role", entry, "error", err)
			}
			continue
		}
		if role == nil {
			// On empty role, log warning and continue
			if b.Logger().IsWarn() {
				b.Logger().Warn("no role info found", "role", entry)
			}
			continue
		}

		roleInfo, err := b.parseRole(role)
		if err != nil {
			if b.Logger().IsWarn() {
				b.Logger().Warn("error parsing role info", "role", entry, "error", err)
			}
			continue
		}

		if keyType, ok := roleInfo["key_type"]; ok {
			keyInfo[entry] = map[string]interface{}{
				"key_type": keyType,
			}
		}
	}

	return logical.ListResponseWithInfo(entries, keyInfo), nil
}

func (b *backend) pathRoleRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	role, err := b.getRole(ctx, req.Storage, d.Get("role").(string))
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	roleInfo, err := b.parseRole(role)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: roleInfo,
	}, nil
}

func (b *backend) pathRoleDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleName := d.Get("role").(string)

	// If the role was given privilege to accept any IP address, there will
	// be an entry for this role in zero-address roles list. Before the role
	// is removed, the entry in the list has to be removed.
	err := b.removeZeroAddressRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Delete(ctx, fmt.Sprintf("roles/%s", roleName))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

const pathRoleHelpSyn = `
Manage the 'roles' that can be created with this backend.
`

const pathRoleHelpDesc = `
This path allows you to manage the roles that are used to generate credentials.

Role takes a 'key_type' parameter that decides what type of credential this role
can generate. If remote hosts have Vault SSH Agent installed, an 'otp' type can
be used, otherwise 'dynamic' type can be used.

If the backend is mounted at "ssh" and the role is created at "ssh/roles/web",
then a user could request for a credential at "ssh/creds/web" for an IP that
belongs to the role. The credential will be for the 'default_user' registered
with the role. There is also an optional parameter 'username' for 'creds/' endpoint.
`
