package ssh

import (
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/(?P<name>[-\\w]+)",
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role",
			},
			"key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Named key in Vault",
			},
			"admin_user": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Admin user at target address",
			},
			"default_user": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Default user to whom the dynamic key is installed",
			},
			"cidr": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "CIDR blocks and IP addresses",
			},
			"port": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Port number for SSH connection",
			},
			"key_type": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "one-time-password or dynamic-key",
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

func createOTPRole(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleName := d.Get("name").(string)
	adminUser := d.Get("admin_user").(string)
	defaultUser := d.Get("default_user").(string)
	cidr := d.Get("cidr").(string)
	port := d.Get("port").(string)
	if roleName == "" {
		return logical.ErrorResponse("Missing role name"), nil
	}
	if adminUser == "" {
		return logical.ErrorResponse("Missing admin username"), nil
	}
	if cidr == "" {
		return logical.ErrorResponse("Missing cidr blocks"), nil
	}
	for _, item := range strings.Split(cidr, ",") {
		_, _, err := net.ParseCIDR(item)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Invalid cidr entry '%s'", item)), nil
		}
	}
	if defaultUser == "" {
		defaultUser = adminUser
	}
	if port == "" {
		port = "22"
	}
	entry, err := logical.StorageEntryJSON(fmt.Sprintf("policy/%s", roleName), sshRole{
		AdminUser:   adminUser,
		DefaultUser: defaultUser,
		CIDR:        cidr,
		Port:        port,
		KeyType:     KeyTypeOTP,
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func createDynamicKeyRole(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleName := d.Get("name").(string)
	keyName := d.Get("key").(string)
	adminUser := d.Get("admin_user").(string)
	defaultUser := d.Get("default_user").(string)
	cidr := d.Get("cidr").(string)
	port := d.Get("port").(string)

	// Input validations
	if roleName == "" {
		return logical.ErrorResponse("Missing role name"), nil
	}
	if keyName == "" {
		return logical.ErrorResponse("Missing key name"), nil
	}
	if adminUser == "" {
		return logical.ErrorResponse("Missing admin username"), nil
	}
	if cidr == "" {
		return logical.ErrorResponse("Missing cidr blocks"), nil
	}

	for _, item := range strings.Split(cidr, ",") {
		_, _, err := net.ParseCIDR(item)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Invalid cidr entry '%s'", item)), nil
		}
	}

	keyEntry, err := req.Storage.Get(fmt.Sprintf("keys/%s", keyName))
	if err != nil || keyEntry == nil {
		return logical.ErrorResponse(fmt.Sprintf("Invalid 'key': '%s'", keyName)), nil
	}

	if defaultUser == "" {
		defaultUser = adminUser
	}
	if port == "" {
		port = "22"
	}

	entry, err := logical.StorageEntryJSON(fmt.Sprintf("policy/%s", roleName), sshRole{
		KeyName:     keyName,
		AdminUser:   adminUser,
		DefaultUser: defaultUser,
		CIDR:        cidr,
		Port:        port,
		KeyType:     KeyTypeDynamic,
	})

	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathRoleWrite(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	keyType := d.Get("key_type").(string)
	if keyType == "" {
		return logical.ErrorResponse("Missing key type"), nil
	}
	keyType = strings.ToLower(keyType)
	if keyType == KeyTypeOTP {
		return createOTPRole(req, d)
	} else if keyType == KeyTypeDynamic {
		return createDynamicKeyRole(req, d)
	} else {
		return logical.ErrorResponse("Invalid key type"), nil
	}
}

func (b *backend) pathRoleRead(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleName := d.Get("name").(string)
	roleEntry, err := req.Storage.Get(fmt.Sprintf("policy/%s", roleName))
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return nil, nil
	}
	var role sshRole
	if err := roleEntry.DecodeJSON(&role); err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"key":          role.KeyName,
			"admin_user":   role.AdminUser,
			"default_user": role.DefaultUser,
			"cidr":         role.CIDR,
			"port":         role.Port,
			"key_type":     role.KeyType,
		},
	}, nil
}

func (b *backend) pathRoleDelete(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleName := d.Get("name").(string)
	err := req.Storage.Delete(fmt.Sprintf("policy/%s", roleName))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

type sshRole struct {
	KeyType     string `json:"key_type"`
	KeyName     string `json:"key"`
	AdminUser   string `json:"admin_user"`
	DefaultUser string `json:"default_user"`
	CIDR        string `json:"cidr"`
	Port        string `json:"port"`
}

const KeyTypeOTP = "otp"
const KeyTypeDynamic = "dynamic"

const pathRoleHelpSyn = `
Manage the 'roles' that can be created with this backend.
`

const pathRoleHelpDesc = `
This path allows you to manage the roles that are used to create
dynamic keys. These roles will be having privileged access to all
the hosts mentioned by CIDR blocks. For example, if the backend
is mounted at "ssh" and the role is created at "ssh/roles/web",
then a user could request for a new key at "ssh/creds/web" for the
supplied username and IP address.

The 'cidr' field takes comma seperated CIDR blocks. The 'admin_user'
should have root access in all the hosts represented by the 'cidr'
field. When the user requests key for an IP, the key will be installed
for the user mentioned by 'default_user' field. The 'key' field takes
a named key which can be configured by 'ssh/keys/' endpoint.
`
