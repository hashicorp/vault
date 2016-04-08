package ssh

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	KeyTypeOTP     = "otp"
	KeyTypeDynamic = "dynamic"
)

// Structure that represents a role in SSH backend. This is a common role structure
// for both OTP and Dynamic roles. Not all the fields are mandatory for both type.
// Some are applicable for one and not for other. It doesn't matter.
type sshRole struct {
	KeyType         string `mapstructure:"key_type" json:"key_type"`
	KeyName         string `mapstructure:"key" json:"key"`
	KeyBits         int    `mapstructure:"key_bits" json:"key_bits"`
	AdminUser       string `mapstructure:"admin_user" json:"admin_user"`
	DefaultUser     string `mapstructure:"default_user" json:"default_user"`
	CIDRList        string `mapstructure:"cidr_list" json:"cidr_list"`
	ExcludeCIDRList string `mapstructure:"exclude_cidr_list" json:"exclude_cidr_list"`
	Port            int    `mapstructure:"port" json:"port"`
	InstallScript   string `mapstructure:"install_script" json:"install_script"`
	AllowedUsers    string `mapstructure:"allowed_users" json:"allowed_users"`
	KeyOptionSpecs  string `mapstructure:"key_option_specs" json:"key_option_specs"`
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
				[Required for both types]
				Name of the role being created.`,
			},
			"key": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Required for Dynamic type] [Not applicable for OTP type]
				Name of the registered key in Vault. Before creating the role, use the
				'keys/' endpoint to create a named key.`,
			},
			"admin_user": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Required for Dynamic type] [Not applicable for OTP type]
				Admin user at remote host. The shared key being registered should be
				for this user and should have root privileges. Everytime a dynamic 
				credential is being generated for other users, Vault uses this admin
				username to login to remote host and install the generated credential
				for the other user.`,
			},
			"default_user": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Required for both types]
				Default username for which a credential will be generated.
				When the endpoint 'creds/' is used without a username, this
				value will be used as default username.`,
			},
			"cidr_list": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Optional for both types]
				Comma separated list of CIDR blocks for which the role is applicable for.
				CIDR blocks can belong to more than one role.`,
			},
			"exclude_cidr_list": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Optional for both types]
				Comma separated list of CIDR blocks. IP addresses belonging to these blocks are not
				accepted by the role. This is particularly useful when big CIDR blocks are being used
				by the role and certain parts of it needs to be kept out.`,
			},
			"port": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `
				[Optional for both types]
				Port number for SSH connection. Default is '22'. Port number does not
				play any role in creation of OTP. For 'otp' type, this is just a way
				to inform client about the port number to use. Port number will be
				returned to client by Vault server along with OTP.`,
			},
			"key_type": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Required for both types] 
				Type of key used to login to hosts. It can be either 'otp' or 'dynamic'.
				'otp' type requires agent to be installed in remote hosts.`,
			},
			"key_bits": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `
				[Optional for Dynamic type] [Not applicable for OTP type]
				Length of the RSA dynamic key in bits. It is 1024 by default or it can be 2048.`,
			},
			"install_script": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Optional for Dynamic type] [Not-applicable for OTP type]
				Script used to install and uninstall public keys in the target machine.
				The inbuilt default install script will be for Linux hosts. For sample
				script, refer the project documentation website.`,
			},
			"allowed_users": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Optional for both types]
				If this option is not specified, client can request for a credential for
				any valid user at the remote host, including the admin user. If only certain
				usernames are to be allowed, then this list enforces it. If this field is
				set, then credentials can only be created for default_user and usernames
				present in this list.
				`,
			},
			"key_option_specs": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Optional for Dynamic type] [Not applicable for OTP type]
				Comma separated option specifications which will be prefixed to RSA key in
				authorized_keys file. Options should be valid and comply with authorized_keys
				file format and should not contain spaces.
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
		return logical.ErrorResponse("Missing role name"), nil
	}

	// Allowed users is an optional field, applicable for both OTP and Dynamic types.
	allowedUsers := d.Get("allowed_users").(string)

	defaultUser := d.Get("default_user").(string)
	if defaultUser == "" {
		return logical.ErrorResponse("Missing default user"), nil
	}

	cidrList := d.Get("cidr_list").(string)

	// Check if all the CIDR blocks are infact valid entries and they don't conflict
	// with each other
	if len(cidrList) != 0 {
		overlaps, err := validateCIDRList(cidrList)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Invalid cidr_list entry. %s", err)), nil
		}
		if len(overlaps) != 0 {
			return logical.ErrorResponse(fmt.Sprintf("CIDR blocks conflicting: %s", overlaps)), nil
		}
	}

	excludeCidrList := d.Get("exclude_cidr_list").(string)

	// Check if all the CIDR blocks are infact valid entries
	if len(excludeCidrList) != 0 {
		_, err := validateCIDRList(excludeCidrList)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Invalid exclude_cidr_list entry. %s", err)), nil
		}
	}

	port := d.Get("port").(int)
	if port == 0 {
		port = 22
	}

	keyType := d.Get("key_type").(string)
	if keyType == "" {
		return logical.ErrorResponse("Missing key type"), nil
	}
	keyType = strings.ToLower(keyType)

	var roleEntry sshRole
	if keyType == KeyTypeOTP {
		// Admin user is not used if OTP key type is used because there is
		// no need to login to remote machine.
		adminUser := d.Get("admin_user").(string)
		if adminUser != "" {
			return logical.ErrorResponse("Admin user not required for OTP type"), nil
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
		// Key name is required by dynamic type and not by OTP type.
		keyName := d.Get("key").(string)
		if keyName == "" {
			return logical.ErrorResponse("Missing key name"), nil
		}
		keyEntry, err := req.Storage.Get(fmt.Sprintf("keys/%s", keyName))
		if err != nil || keyEntry == nil {
			return logical.ErrorResponse(fmt.Sprintf("Invalid 'key': '%s'", keyName)), nil
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
			return logical.ErrorResponse("Missing admin username"), nil
		}

		// This defaults to 1024 and it can also be 2048.
		keyBits := d.Get("key_bits").(int)
		if keyBits != 0 && keyBits != 1024 && keyBits != 2048 {
			return logical.ErrorResponse("Invalid key_bits field"), nil
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
	} else {
		return logical.ErrorResponse("Invalid key type"), nil
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

func (b *backend) pathRoleList(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List("roles/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func (b *backend) pathRoleRead(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	role, err := b.getRole(req.Storage, d.Get("role").(string))
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	// Return information should be based on the key type of the role
	if role.KeyType == KeyTypeOTP {
		return &logical.Response{
			Data: map[string]interface{}{
				"default_user":      role.DefaultUser,
				"cidr_list":         role.CIDRList,
				"exclude_cidr_list": role.ExcludeCIDRList,
				"key_type":          role.KeyType,
				"port":              role.Port,
				"allowed_users":     role.AllowedUsers,
			},
		}, nil
	} else {
		return &logical.Response{
			Data: map[string]interface{}{
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
			},
		}, nil
	}
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
