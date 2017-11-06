package ssh

import (
	"fmt"
	"strings"

	"time"

	"github.com/hashicorp/vault/helper/cidrutil"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	KeyTypeOTP     = "otp"
	KeyTypeDynamic = "dynamic"
	KeyTypeCA      = "ca"
)

// Structure that represents a role in SSH backend. This is a common role structure
// for both OTP and Dynamic roles. Not all the fields are mandatory for both type.
// Some are applicable for one and not for other. It doesn't matter.
type sshRole struct {
	KeyType                string            `mapstructure:"key_type" json:"key_type"`
	KeyName                string            `mapstructure:"key" json:"key"`
	KeyBits                int               `mapstructure:"key_bits" json:"key_bits"`
	AdminUser              string            `mapstructure:"admin_user" json:"admin_user"`
	DefaultUser            string            `mapstructure:"default_user" json:"default_user"`
	CIDRList               string            `mapstructure:"cidr_list" json:"cidr_list"`
	ExcludeCIDRList        string            `mapstructure:"exclude_cidr_list" json:"exclude_cidr_list"`
	Port                   int               `mapstructure:"port" json:"port"`
	InstallScript          string            `mapstructure:"install_script" json:"install_script"`
	AllowedUsers           string            `mapstructure:"allowed_users" json:"allowed_users"`
	AllowedDomains         string            `mapstructure:"allowed_domains" json:"allowed_domains"`
	KeyOptionSpecs         string            `mapstructure:"key_option_specs" json:"key_option_specs"`
	MaxTTL                 string            `mapstructure:"max_ttl" json:"max_ttl"`
	TTL                    string            `mapstructure:"ttl" json:"ttl"`
	DefaultCriticalOptions map[string]string `mapstructure:"default_critical_options" json:"default_critical_options"`
	DefaultExtensions      map[string]string `mapstructure:"default_extensions" json:"default_extensions"`
	AllowedCriticalOptions string            `mapstructure:"allowed_critical_options" json:"allowed_critical_options"`
	AllowedExtensions      string            `mapstructure:"allowed_extensions" json:"allowed_extensions"`
	AllowUserCertificates  bool              `mapstructure:"allow_user_certificates" json:"allow_user_certificates"`
	AllowHostCertificates  bool              `mapstructure:"allow_host_certificates" json:"allow_host_certificates"`
	AllowBareDomains       bool              `mapstructure:"allow_bare_domains" json:"allow_bare_domains"`
	AllowSubdomains        bool              `mapstructure:"allow_subdomains" json:"allow_subdomains"`
	AllowUserKeyIDs        bool              `mapstructure:"allow_user_key_ids" json:"allow_user_key_ids"`
	KeyIDFormat            string            `mapstructure:"key_id_format" json:"key_id_format"`
}

func pathListRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},

		HelpSynopsis:    pathRoleHelpSyn,
		HelpDescription: pathRoleHelpDesc,
	}
}

func pathRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("role"),
		Fields: map[string]*framework.FieldSchema{
			"role": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Required for all types]
				Name of the role being created.`,
			},
			"key": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Required for Dynamic type] [Not applicable for OTP type] [Not applicable for CA type]
				Name of the registered key in Vault. Before creating the role, use the
				'keys/' endpoint to create a named key.`,
			},
			"admin_user": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Required for Dynamic type] [Not applicable for OTP type] [Not applicable for CA type]
				Admin user at remote host. The shared key being registered should be
				for this user and should have root privileges. Everytime a dynamic 
				credential is being generated for other users, Vault uses this admin
				username to login to remote host and install the generated credential
				for the other user.`,
			},
			"default_user": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Required for Dynamic type] [Required for OTP type] [Optional for CA type]
				Default username for which a credential will be generated.
				When the endpoint 'creds/' is used without a username, this
				value will be used as default username.`,
			},
			"cidr_list": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Optional for Dynamic type] [Optional for OTP type] [Not applicable for CA type]
				Comma separated list of CIDR blocks for which the role is applicable for.
				CIDR blocks can belong to more than one role.`,
			},
			"exclude_cidr_list": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Optional for Dynamic type] [Optional for OTP type] [Not applicable for CA type]
				Comma separated list of CIDR blocks. IP addresses belonging to these blocks are not
				accepted by the role. This is particularly useful when big CIDR blocks are being used
				by the role and certain parts of it needs to be kept out.`,
			},
			"port": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `
				[Optional for Dynamic type] [Optional for OTP type] [Not applicable for CA type]
				Port number for SSH connection. Default is '22'. Port number does not
				play any role in creation of OTP. For 'otp' type, this is just a way
				to inform client about the port number to use. Port number will be
				returned to client by Vault server along with OTP.`,
			},
			"key_type": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Required for all types]
				Type of key used to login to hosts. It can be either 'otp', 'dynamic' or 'ca'.
				'otp' type requires agent to be installed in remote hosts.`,
			},
			"key_bits": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `
				[Optional for Dynamic type] [Not applicable for OTP type] [Not applicable for CA type]
				Length of the RSA dynamic key in bits. It is 1024 by default or it can be 2048.`,
			},
			"install_script": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Optional for Dynamic type] [Not-applicable for OTP type] [Not applicable for CA type]
				Script used to install and uninstall public keys in the target machine.
				The inbuilt default install script will be for Linux hosts. For sample
				script, refer the project documentation website.`,
			},
			"allowed_users": &framework.FieldSchema{
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
			"allowed_domains": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Not applicable for Dynamic type] [Not applicable for OTP type] [Optional for CA type]
				If this option is not specified, client can request for a signed certificate for any
				valid host. If only certain domains are allowed, then this list enforces it.
				`,
			},
			"key_option_specs": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Optional for Dynamic type] [Not applicable for OTP type] [Not applicable for CA type]
				Comma separated option specifications which will be prefixed to RSA key in
				authorized_keys file. Options should be valid and comply with authorized_keys
				file format and should not contain spaces.
				`,
			},
			"ttl": &framework.FieldSchema{
				Type: framework.TypeDurationSecond,
				Description: `
				[Not applicable for Dynamic type] [Not applicable for OTP type] [Optional for CA type]
				The lease duration if no specific lease duration is
				requested. The lease duration controls the expiration
				of certificates issued by this backend. Defaults to
				the value of max_ttl.`,
			},
			"max_ttl": &framework.FieldSchema{
				Type: framework.TypeDurationSecond,
				Description: `
				[Not applicable for Dynamic type] [Not applicable for OTP type] [Optional for CA type]
				The maximum allowed lease duration
				`,
			},
			"allowed_critical_options": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Not applicable for Dynamic type] [Not applicable for OTP type] [Optional for CA type]
				A comma-separated list of critical options that certificates can have when signed.
 				To allow any critical options, set this to an empty string.
 				`,
			},
			"allowed_extensions": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Not applicable for Dynamic type] [Not applicable for OTP type] [Optional for CA type]
				A comma-separated list of extensions that certificates can have when signed.
				To allow any extensions, set this to an empty string.
				`,
			},
			"default_critical_options": &framework.FieldSchema{
				Type: framework.TypeMap,
				Description: `
				[Not applicable for Dynamic type] [Not applicable for OTP type]
				[Optional for CA type] Critical options certificates should
				have if none are provided when signing. This field takes in key
				value pairs in JSON format.  Note that these are not restricted
				by "allowed_critical_options". Defaults to none.
				`,
			},
			"default_extensions": &framework.FieldSchema{
				Type: framework.TypeMap,
				Description: `
				[Not applicable for Dynamic type] [Not applicable for OTP type]
				[Optional for CA type] Extensions certificates should have if
				none are provided when signing. This field takes in key value
				pairs in JSON format. Note that these are not restricted by
				"allowed_extensions". Defaults to none.
				`,
			},
			"allow_user_certificates": &framework.FieldSchema{
				Type: framework.TypeBool,
				Description: `
				[Not applicable for Dynamic type] [Not applicable for OTP type] [Optional for CA type]
				If set, certificates are allowed to be signed for use as a 'user'.
				`,
				Default: false,
			},
			"allow_host_certificates": &framework.FieldSchema{
				Type: framework.TypeBool,
				Description: `
				[Not applicable for Dynamic type] [Not applicable for OTP type] [Optional for CA type]
				If set, certificates are allowed to be signed for use as a 'host'.
				`,
				Default: false,
			},
			"allow_bare_domains": &framework.FieldSchema{
				Type: framework.TypeBool,
				Description: `
				[Not applicable for Dynamic type] [Not applicable for OTP type] [Optional for CA type]
				If set, host certificates that are requested are allowed to use the base domains listed in
				"allowed_domains", e.g. "example.com".
				This is a separate option as in some cases this can be considered a security threat.
				`,
			},
			"allow_subdomains": &framework.FieldSchema{
				Type: framework.TypeBool,
				Description: `
				[Not applicable for Dynamic type] [Not applicable for OTP type] [Optional for CA type]
				If set, host certificates that are requested are allowed to use subdomains of those listed in "allowed_domains".
				`,
			},
			"allow_user_key_ids": &framework.FieldSchema{
				Type: framework.TypeBool,
				Description: `
				[Not applicable for Dynamic type] [Not applicable for OTP type] [Optional for CA type]
				If true, users can override the key ID for a signed certificate with the "key_id" field.
				When false, the key ID will always be the token display name.
				The key ID is logged by the SSH server and can be useful for auditing.
				`,
			},
			"key_id_format": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Not applicable for Dynamic type] [Not applicable for OTP type] [Optional for CA type]
				When supplied, this value specifies a custom format for the key id of a signed certificate.
				The following variables are availble for use: '{{token_display_name}}' - The display name of
				the token used to make the request. '{{role_name}}' - The name of the role signing the request.
				'{{public_key_hash}}' - A SHA256 checksum of the public key that is being signed.
				`,
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

func (b *backend) pathRoleWrite(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleName := d.Get("role").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role name"), nil
	}

	// Allowed users is an optional field, applicable for both OTP and Dynamic types.
	allowedUsers := d.Get("allowed_users").(string)

	// Validate the CIDR blocks
	cidrList := d.Get("cidr_list").(string)
	if cidrList != "" {
		valid, err := cidrutil.ValidateCIDRListString(cidrList, ",")
		if err != nil {
			return nil, fmt.Errorf("failed to validate cidr_list: %v", err)
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
			return nil, fmt.Errorf("failed to validate exclude_cidr_list entry: %v", err)
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

		// Admin user is not used if OTP key type is used because there is
		// no need to login to remote machine.
		adminUser := d.Get("admin_user").(string)
		if adminUser != "" {
			return logical.ErrorResponse("admin user not required for OTP type"), nil
		}

		// Below are the only fields used from the role structure for OTP type.
		roleEntry = sshRole{
			DefaultUser:     defaultUser,
			CIDRList:        cidrList,
			ExcludeCIDRList: excludeCidrList,
			KeyType:         KeyTypeOTP,
			Port:            port,
			AllowedUsers:    allowedUsers,
		}
	} else if keyType == KeyTypeDynamic {
		defaultUser := d.Get("default_user").(string)
		if defaultUser == "" {
			return logical.ErrorResponse("missing default user"), nil
		}
		// Key name is required by dynamic type and not by OTP type.
		keyName := d.Get("key").(string)
		if keyName == "" {
			return logical.ErrorResponse("missing key name"), nil
		}
		keyEntry, err := req.Storage.Get(fmt.Sprintf("keys/%s", keyName))
		if err != nil || keyEntry == nil {
			return logical.ErrorResponse(fmt.Sprintf("invalid 'key': %q", keyName)), nil
		}

		installScript := d.Get("install_script").(string)
		keyOptionSpecs := d.Get("key_option_specs").(string)

		// Setting the default script here. The script will install the
		// generated public key in the authorized_keys file of linux host.
		if installScript == "" {
			installScript = DefaultPublicKeyInstallScript
		}

		adminUser := d.Get("admin_user").(string)
		if adminUser == "" {
			return logical.ErrorResponse("missing admin username"), nil
		}

		// This defaults to 1024 and it can also be 2048.
		keyBits := d.Get("key_bits").(int)
		if keyBits != 0 && keyBits != 1024 && keyBits != 2048 {
			return logical.ErrorResponse("invalid key_bits field"), nil
		}

		// If user has not set this field, default it to 1024
		if keyBits == 0 {
			keyBits = 1024
		}

		// Store all the fields required by dynamic key type
		roleEntry = sshRole{
			KeyName:         keyName,
			AdminUser:       adminUser,
			DefaultUser:     defaultUser,
			CIDRList:        cidrList,
			ExcludeCIDRList: excludeCidrList,
			Port:            port,
			KeyType:         KeyTypeDynamic,
			KeyBits:         keyBits,
			InstallScript:   installScript,
			AllowedUsers:    allowedUsers,
			KeyOptionSpecs:  keyOptionSpecs,
		}
	} else if keyType == KeyTypeCA {
		role, errorResponse := b.createCARole(allowedUsers, d.Get("default_user").(string), d)
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

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) createCARole(allowedUsers, defaultUser string, data *framework.FieldData) (*sshRole, *logical.Response) {
	ttl := time.Duration(data.Get("ttl").(int)) * time.Second
	maxTTL := time.Duration(data.Get("max_ttl").(int)) * time.Second
	role := &sshRole{
		AllowedCriticalOptions: data.Get("allowed_critical_options").(string),
		AllowedExtensions:      data.Get("allowed_extensions").(string),
		AllowUserCertificates:  data.Get("allow_user_certificates").(bool),
		AllowHostCertificates:  data.Get("allow_host_certificates").(bool),
		AllowedUsers:           allowedUsers,
		AllowedDomains:         data.Get("allowed_domains").(string),
		DefaultUser:            defaultUser,
		AllowBareDomains:       data.Get("allow_bare_domains").(bool),
		AllowSubdomains:        data.Get("allow_subdomains").(bool),
		AllowUserKeyIDs:        data.Get("allow_user_key_ids").(bool),
		KeyIDFormat:            data.Get("key_id_format").(string),
		KeyType:                KeyTypeCA,
	}

	if !role.AllowUserCertificates && !role.AllowHostCertificates {
		return nil, logical.ErrorResponse("Either 'allow_user_certificates' or 'allow_host_certificates' must be set to 'true'")
	}

	defaultCriticalOptions := convertMapToStringValue(data.Get("default_critical_options").(map[string]interface{}))
	defaultExtensions := convertMapToStringValue(data.Get("default_extensions").(map[string]interface{}))

	if ttl != 0 && maxTTL != 0 && ttl > maxTTL {
		return nil, logical.ErrorResponse(
			`"ttl" value must be less than "max_ttl" when both are specified`)
	}

	// Persist TTLs
	role.TTL = ttl.String()
	role.MaxTTL = maxTTL.String()
	role.DefaultCriticalOptions = defaultCriticalOptions
	role.DefaultExtensions = defaultExtensions

	return role, nil
}

func (b *backend) getRole(s logical.Storage, n string) (*sshRole, error) {
	entry, err := s.Get("roles/" + n)
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

	return &result, nil
}

// parseRole converts a sshRole object into its map[string]interface representation,
// with appropriate values for each KeyType. If the KeyType is invalid, it will retun
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
			"allowed_users":            role.AllowedUsers,
			"allowed_domains":          role.AllowedDomains,
			"default_user":             role.DefaultUser,
			"ttl":                      int64(ttl.Seconds()),
			"max_ttl":                  int64(maxTTL.Seconds()),
			"allowed_critical_options": role.AllowedCriticalOptions,
			"allowed_extensions":       role.AllowedExtensions,
			"allow_user_certificates":  role.AllowUserCertificates,
			"allow_host_certificates":  role.AllowHostCertificates,
			"allow_bare_domains":       role.AllowBareDomains,
			"allow_subdomains":         role.AllowSubdomains,
			"allow_user_key_ids":       role.AllowUserKeyIDs,
			"key_id_format":            role.KeyIDFormat,
			"key_type":                 role.KeyType,
			"default_critical_options": role.DefaultCriticalOptions,
			"default_extensions":       role.DefaultExtensions,
		}
	case KeyTypeDynamic:
		result = map[string]interface{}{
			"key":               role.KeyName,
			"admin_user":        role.AdminUser,
			"default_user":      role.DefaultUser,
			"cidr_list":         role.CIDRList,
			"exclude_cidr_list": role.ExcludeCIDRList,
			"port":              role.Port,
			"key_type":          role.KeyType,
			"key_bits":          role.KeyBits,
			"allowed_users":     role.AllowedUsers,
			"key_option_specs":  role.KeyOptionSpecs,
			// Returning install script will make the output look messy.
			// But this is one way for clients to see the script that is
			// being used to install the key. If there is some problem,
			// the script can be modified and configured by clients.
			"install_script": role.InstallScript,
		}
	default:
		return nil, fmt.Errorf("invalid key type: %v", role.KeyType)
	}

	return result, nil
}

func (b *backend) pathRoleList(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List("roles/")
	if err != nil {
		return nil, err
	}

	keyInfo := map[string]interface{}{}
	for _, entry := range entries {
		role, err := b.getRole(req.Storage, entry)
		if err != nil {
			// On error, log warning and continue
			if b.Logger().IsWarn() {
				b.Logger().Warn("ssh: error getting role info", "role", entry, "error", err)
			}
			continue
		}
		if role == nil {
			// On empty role, log warning and continue
			if b.Logger().IsWarn() {
				b.Logger().Warn("ssh: no role info found", "role", entry)
			}
			continue
		}

		roleInfo, err := b.parseRole(role)
		if err != nil {
			if b.Logger().IsWarn() {
				b.Logger().Warn("ssh: error parsing role info", "role", entry, "error", err)
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

func (b *backend) pathRoleRead(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	role, err := b.getRole(req.Storage, d.Get("role").(string))
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

func (b *backend) pathRoleDelete(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleName := d.Get("role").(string)

	// If the role was given privilege to accept any IP address, there will
	// be an entry for this role in zero-address roles list. Before the role
	// is removed, the entry in the list has to be removed.
	err := b.removeZeroAddressRole(req.Storage, roleName)
	if err != nil {
		return nil, err
	}

	err = req.Storage.Delete(fmt.Sprintf("roles/%s", roleName))
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
