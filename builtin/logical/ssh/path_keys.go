package ssh

import (
	"fmt"

	"golang.org/x/crypto/ssh"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

type sshHostKey struct {
	Key string `json:"key"`
}

func pathKeys(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("key_name"),
		Fields: map[string]*framework.FieldSchema{
			"key_name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "[Required] Name of the key",
			},
			"key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "[Required] SSH private key with super user privileges in host",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathKeysWrite,
			logical.DeleteOperation: b.pathKeysDelete,
		},
		HelpSynopsis:    pathKeysSyn,
		HelpDescription: pathKeysDesc,
	}
}

func (b *backend) getKey(s logical.Storage, n string) (*sshHostKey, error) {
	entry, err := s.Get("keys/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result sshHostKey
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (b *backend) pathKeysDelete(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	keyName := d.Get("key_name").(string)
	keyPath := fmt.Sprintf("keys/%s", keyName)
	err := req.Storage.Delete(keyPath)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) pathKeysWrite(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	keyName := d.Get("key_name").(string)
	if keyName == "" {
		return logical.ErrorResponse("Missing key_name"), nil
	}

	keyString := d.Get("key").(string)

	// Check if the key provided is infact a private key
	signer, err := ssh.ParsePrivateKey([]byte(keyString))
	if err != nil || signer == nil {
		return logical.ErrorResponse("Invalid key"), nil
	}

	if keyString == "" {
		return logical.ErrorResponse("Missing key"), nil
	}

	keyPath := fmt.Sprintf("keys/%s", keyName)

	// Store the key
	entry, err := logical.StorageEntryJSON(keyPath, map[string]interface{}{
		"key": keyString,
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}
	return nil, nil
}

const pathKeysSyn = `
Register a shared private key with Vault.
`

const pathKeysDesc = `
Vault uses this key to install and uninstall dynamic keys in remote hosts. This
key should have sudoer privileges in remote hosts. This enables installing keys
for unprivileged usernames.

If this backend is mounted as "ssh", then the endpoint for registering shared
key is "ssh/keys/<name>". The name given here can be associated with any number
of roles via the endpoint "ssh/roles/".
`
