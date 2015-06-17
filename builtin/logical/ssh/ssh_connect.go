package ssh

import (
	"log"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func sshConnect(b *backend) *framework.Path {
	log.Printf("Vishal: ssh.sshConnect\n")
	return &framework.Path{
		Pattern: "connect",
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "username at SSH host",
			},
			"address": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "IPv4 address of SSH host",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.sshConnectWrite,
		},
		HelpSynopsis:    sshConnectHelpSyn,
		HelpDescription: sshConnectHelpDesc,
	}
}

func (b *backend) sshConnectWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	log.Printf("Vishal: ssh.sshConnectWrite username:%#v address:%#v\n", d.Get("username").(string), d.Get("address").(string))
	return &logical.Response{
		Data: map[string]interface{}{
			"key": "createdKey",
		},
	}, nil
}

const sshConnectHelpSyn = `
sshConnectionHelpSyn
`

const sshConnectHelpDesc = `
sshConnectionHelpDesc
`
