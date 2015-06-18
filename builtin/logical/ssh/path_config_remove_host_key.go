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
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathRemoveHostKeyWrite,
		},
		HelpSynopsis:    pathConfigRemoveHostKeySyn,
		HelpDescription: pathConfigRemoveHostKeyDesc,
	}
}

func (b *backend) pathRemoveHostKeyWrite(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	log.Printf("Vishal: ssh.pathRemoveHostKeyWrite\n")
	username := d.Get("username").(string)
	ip := d.Get("ip").(string)
	//TODO: parse ip into ipv4 address and validate it
	log.Printf("Vishal: ssh.pathRemoveHostKeyWrite username:%#v ip:%#v\n", username, ip)
	err := req.Storage.Delete("hosts/" + ip + "/" + username)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

const pathConfigRemoveHostKeySyn = `
pathConfigRemoveHostKeySyn
`

const pathConfigRemoveHostKeyDesc = `
pathConfigRemoveHostKeyDesc
`
