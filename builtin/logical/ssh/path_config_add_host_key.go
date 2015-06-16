package ssh

import (
	"log"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigAddHostKey(b *backend) *framework.Path {
	log.Printf("Vishal: ssh.pathConfigAddHostKey\n")
	return &framework.Path{
		Pattern: "config/addhostkey",
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
			logical.WriteOperation: b.pathAddHostKeyWrite,
		},
		HelpSynopsis:    pathConfigAddHostKeySyn,
		HelpDescription: pathConfigAddHostKeyDesc,
	}
}

func (b *backend) pathAddHostKeyWrite(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	log.Printf("Vishal: ssh.pathAddHostKeyWrite\n")
	return nil, nil
}

const pathConfigAddHostKeySyn = `
pathConfigAddHostKeySyn
`

const pathConfigAddHostKeyDesc = `
pathConfigAddHostKeyDesc
`
