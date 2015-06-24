package ssh

import (
	"log"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathKeys(b *backend) *framework.Path {
	log.Printf("Vishal: ssh.pathConfigAddHostKey\n")
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
		HelpSynopsis:    pathConfigAddHostKeySyn,
		HelpDescription: pathConfigAddHostKeyDesc,
	}
}

func (b *backend) pathKeysRead(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	keyName := d.Get("name").(string)
	log.Printf("Vishal: ssh.pathKeysRead: keyName: %#v\n", keyName)
	keyPath := "keys/" + keyName
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
	log.Printf("Vishal: ssh.pathKeysDelete: keyName: %#v\n", keyName)
	keyPath := "keys/" + keyName
	err := req.Storage.Delete(keyPath)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) pathKeysWrite(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Vishal: ssh.pathKeysWrite\n")

	keyName := d.Get("name").(string)
	keyString := d.Get("key").(string)
	keyPath := "keys/" + keyName

	log.Printf("Vishal: ssh.path_keys.pathKeysWrite: keyPath: %#v\n", keyPath)
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

const pathConfigAddHostKeySyn = `
pathConfigAddHostKeySyn
`

const pathConfigAddHostKeyDesc = `
pathConfigAddHostKeyDesc
`
