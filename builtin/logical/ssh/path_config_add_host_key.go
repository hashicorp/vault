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

func (b *backend) pathAddHostKeyWrite(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Vishal: ssh.pathAddHostKeyWrite\n")
	username := d.Get("username").(string)
	ip := d.Get("ip").(string)
	//TODO: parse ip into ipv4 address and validate it
	key := d.Get("key").(string)
	//log.Printf("Vishal: ssh.pathAddHostKeyWrite username:%#v ip:%#v key:%#v\n", username, ip, key)

	hostKeyPath := "hosts/" + ip + "/" + username
	log.Printf("Vishal: hostKeyPath: %#v\n", hostKeyPath)
	entry, err := logical.StorageEntryJSON("hosts/"+ip+"/"+username, &sshHostKey{
		Key: key,
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}
	return nil, nil
}

type sshHostKey struct {
	Key string
}

const pathConfigAddHostKeySyn = `
pathConfigAddHostKeySyn
`

const pathConfigAddHostKeyDesc = `
pathConfigAddHostKeyDesc
`
