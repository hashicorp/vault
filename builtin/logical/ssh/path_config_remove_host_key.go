package ssh

import (
	"log"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigRemoveHostKey(b *backend) *framework.Path {
	log.Printf("Vishal: ssh.pathConfigRemoveHostKey\n")
	return &framework.Path{
		Pattern: "config/removehostkey",
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Username in host.",
			},
			"ip": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "IP address of host.",
			},
			"key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "SSH private key for host.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathRemoveHostKeyWrite,
		},
		HelpSynopsis:    pathConfigRemoveHostKeySyn,
		HelpDescription: pathConfigRemoveHostKeyDesc,
	}
}

func (b *backend) pathRemoveHostKeyWrite(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	log.Printf("Vishal: ssh.pathRemoveHostKeyWrite\n")
	return nil, nil
}

const pathConfigRemoveHostKeySyn = `
pathConfigRemoveHostKeySyn
`

const pathConfigRemoveHostKeyDesc = `
pathConfigRemoveHostKeyDesc
`
