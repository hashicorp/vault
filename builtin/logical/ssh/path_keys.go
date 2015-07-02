package ssh

import (
	"fmt"
	"log"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathKeys(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "keys/(?P<name>\\w+)",
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "IP address of host.",
			},
			"key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "SSH private key for host.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathKeysRead,
			logical.WriteOperation:  b.pathKeysWrite,
			logical.DeleteOperation: b.pathKeysDelete,
		},
		HelpSynopsis:    pathKeysSyn,
		HelpDescription: pathKeysDesc,
	}
}

func (b *backend) pathKeysRead(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	keyName := d.Get("name").(string)
	keyPath := fmt.Sprintf("keys/%s", keyName)
	entry, err := req.Storage.Get(keyPath)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"key": string(entry.Value),
		},
	}, nil
}

func (b *backend) pathKeysDelete(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	keyName := d.Get("name").(string)
	keyPath := fmt.Sprintf("keys/%s", keyName)
	err := req.Storage.Delete(keyPath)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) pathKeysWrite(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	keyName := d.Get("name").(string)
	keyString := d.Get("key").(string)

	if keyString == "" {
		return nil, fmt.Errorf("Invalid 'key'")
	}

	keyPath := fmt.Sprintf("keys/%s", keyName)

	entry, err := logical.StorageEntryJSON(keyPath, &sshHostKey{
		Key: keyString,
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

const pathKeysSyn = `
Register a shared key which can be used to install dynamic key
in remote machine.
`

const pathKeysDesc = `
The shared key registered will be used to install and uninstall
dynamic keys in remote machine. This key should have "root" 
privileges which enables installing keys for unprivileged usernames.
If this backend is mounted as "ssh", then the endpoint for registering
shared key is "ssh/keys/rack1", if "rack1" is the user coined 
name for the key. The name given here can be associated with any
number of roles via the endpoint "ssh/roles/".
`
