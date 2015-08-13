package ssh

import (
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const KeyTypeOTP = "otp"
const KeyTypeDynamic = "dynamic"

func pathRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/(?P<role>[-\\w]+)",
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
				[Required for dynamic type] [Not applicable for otp type]
				Name of the registered key in Vault. Before creating the role, use the
				'keys/' endpoint to create a named key.`,
			},
			"admin_user": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Required for dynamic type] [Not applicable for otp type]
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
				[Required for both types]
				Comma separated list of CIDR blocks for which the role is applicable for.
				CIDR blocks can belong to more than one role.`,
			},
			"port": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `
				[Optional for both types]
				Port number for SSH connection. Default is '22'.`,
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
				[Optional for dynamic type] [Not applicable for otp type]
				Length of the RSA dynamic key in bits. It can be one of 1024, 2048 or 4096.`,
			},
			"install_script": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `
				[Optional for dynamic type][Not-applicable for otp type]
				Script used to install and uninstall public keys in the target machine.
				The inbuilt default install script will be for Linux hosts. For sample
				script, refer the project's documentation website.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathRoleRead,
			logical.WriteOperation:  b.pathRoleWrite,
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

	defaultUser := d.Get("default_user").(string)
	if defaultUser == "" {
		return logical.ErrorResponse("Missing default user"), nil
	}

	cidrList := d.Get("cidr_list").(string)
	if cidrList == "" {
		return logical.ErrorResponse("Missing CIDR blocks"), nil
	}
	for _, item := range strings.Split(cidrList, ",") {
		_, _, err := net.ParseCIDR(item)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Invalid CIDR list entry '%s'", item)), nil
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

	var err error
	var roleEntry sshRole
	if keyType == KeyTypeOTP {
		adminUser := d.Get("admin_user").(string)
		if adminUser != "" {
			return logical.ErrorResponse("Admin user not required for OTP type"), nil
		}

		roleEntry = sshRole{
			DefaultUser: defaultUser,
			CIDRList:    cidrList,
			KeyType:     KeyTypeOTP,
			Port:        port,
		}
	} else if keyType == KeyTypeDynamic {
		keyName := d.Get("key").(string)
		if keyName == "" {
			return logical.ErrorResponse("Missing key name"), nil
		}
		keyEntry, err := req.Storage.Get(fmt.Sprintf("keys/%s", keyName))
		if err != nil || keyEntry == nil {
			return logical.ErrorResponse(fmt.Sprintf("Invalid 'key': '%s'", keyName)), nil
		}

		installScript := d.Get("install_script").(string)
		if installScript == "" {
			return logical.ErrorResponse("Missing install script"), nil
		}

		adminUser := d.Get("admin_user").(string)
		if adminUser == "" {
			return logical.ErrorResponse("Missing admin username"), nil
		}

		keyBits := d.Get("key_bits").(int)
		if keyBits != 0 && keyBits != 1024 && keyBits != 2048 && keyBits != 4096 {
			return logical.ErrorResponse("Invalid key_bits field"), nil
		}
		if keyBits == 0 {
			keyBits = 2048
		}

		roleEntry = sshRole{
			KeyName:       keyName,
			AdminUser:     adminUser,
			DefaultUser:   defaultUser,
			CIDRList:      cidrList,
			Port:          port,
			KeyType:       KeyTypeDynamic,
			KeyBits:       keyBits,
			InstallScript: installScript,
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

func (b *backend) pathRoleRead(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	role, err := b.getRole(req.Storage, d.Get("role").(string))
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	if role.KeyType == KeyTypeOTP {
		return &logical.Response{
			Data: map[string]interface{}{
				"default_user": role.DefaultUser,
				"cidr_list":    role.CIDRList,
				"port":         role.Port,
				"key_type":     role.KeyType,
			},
		}, nil
	} else {
		return &logical.Response{
			Data: map[string]interface{}{
				"key":          role.KeyName,
				"admin_user":   role.AdminUser,
				"default_user": role.DefaultUser,
				"cidr_list":    role.CIDRList,
				"port":         role.Port,
				"key_type":     role.KeyType,
			},
		}, nil
	}
}

func (b *backend) pathRoleDelete(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleName := d.Get("role").(string)
	err := req.Storage.Delete(fmt.Sprintf("roles/%s", roleName))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

type sshRole struct {
	KeyType       string `mapstructure:"key_type" json:"key_type"`
	KeyName       string `mapstructure:"key" json:"key"`
	KeyBits       int    `mapstructure:"key_bits" json:"key_bits"`
	AdminUser     string `mapstructure:"admin_user" json:"admin_user"`
	DefaultUser   string `mapstructure:"default_user" json:"default_user"`
	CIDRList      string `mapstructure:"cidr_list" json:"cidr_list"`
	Port          int    `mapstructure:"port" json:"port"`
	InstallScript string `mapstructure:"install_script" json:"install_script"`
}

const pathRoleHelpSyn = `
Manage the 'roles' that can be created with this backend.
`

const pathRoleHelpDesc = `
This path allows you to manage the roles that are used to generate 
credentials. These roles will be having privileged access to all
the hosts mentioned by CIDR blocks. For example, if the backend
is mounted at "ssh" and the role is created at "ssh/roles/web",
then a user could request for a new key at "ssh/creds/web" for the
supplied username and IP address.

The 'cidr_list' field takes comma seperated CIDR blocks. The 'admin_user'
should have root access in all the hosts represented by the 'cidr_list'
field. When the user requests key for an IP, the key will be installed
for the user mentioned by 'default_user' field. The 'key' field takes
a named key which can be configured by 'ssh/keys/' endpoint.
`
