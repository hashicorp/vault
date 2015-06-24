package ssh

import (
	"log"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathRoles(b *backend) *framework.Path {
	log.Printf("Vishal: ssh.pathRoles\n")
	return &framework.Path{
		Pattern: "roles/(?P<name>\\w+)",
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
	log.Printf("Vishal: ssh.pathRoleWrite\n")

	roleName := d.Get("name").(string)
	keyName := d.Get("key").(string)
	adminUser := d.Get("admin_user").(string)
	defaultUser := d.Get("default_user").(string)
	cidr := d.Get("cidr").(string)

	log.Printf("Vishal: name[%s] key[%s] admin_user[%s] default_user[%s] cidr[%s]\n", roleName, keyName, adminUser, defaultUser, cidr)

	rolePath := "policy/" + roleName

	entry, err := logical.StorageEntryJSON(rolePath, sshRole{
		KeyName:     keyName,
		AdminUser:   adminUser,
		DefaultUser: defaultUser,
		CIDR:        cidr,
	})

	if err != nil {
		return nil, err
	}

	log.Printf("Vishal: entryJSON:%s\n", entry.Value)
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathRoleRead(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	log.Printf("Vishal: ssh.pathRoleRead\n")
	roleName := d.Get("name").(string)
	rolePath := "policy/" + roleName
	entry, err := req.Storage.Get(rolePath)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"policy": string(entry.Value),
		},
	}, nil
}

func (b *backend) pathRoleDelete(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	log.Printf("Vishal: ssh.pathRoleDelete\n")
	roleName := d.Get("name").(string)
	rolePath := "policy/" + roleName
	err := req.Storage.Delete(rolePath)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

type sshRole struct {
	KeyName     string `json:"key"`
	AdminUser   string `json:"admin_user"`
	DefaultUser string `json:"default_user"`
	CIDR        string `json: "cidr"`
}

const pathRoleHelpSyn = `
Manage the roles that can be created with this backend.
`

const pathRoleHelpDesc = `
This path lets you manage the roles that can be created with this backend.
`
