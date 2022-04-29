package transit

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathImportVersion() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name") + "/import_version",
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The name of the key",
			},
			"ciphertext": {
				Type: framework.TypeString,
				Description: `The base64-encoded ciphertext of the keys. The AES key should be encrypted using OAEP 
with the wrapping key and then concatenated with the import key, wrapped by the AES key.`,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathImportVersionWrite,
		},
		HelpSynopsis:    pathImportVersionWriteSyn,
		HelpDescription: pathImportVersionWriteDesc,
	}
}

func (b *backend) pathImportVersionWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	ciphertextString := d.Get("ciphertext").(string)

	polReq := keysutil.PolicyRequest{
		Storage: req.Storage,
		Name:    name,
		Upsert:  false,
	}

	p, _, err := b.GetPolicy(ctx, polReq, b.GetRandomReader())
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("no key found with name %s; to import a new key, use the import/ endpoint", name)
	}
	if !p.Imported {
		return nil, errors.New("the import_version endpoint can only be used with an imported key")
	}
	if p.ConvergentEncryption {
		return nil, errors.New("import_version cannot be used on keys with convergent encryption enabled")
	}

	if !b.System().CachingDisabled() {
		p.Lock(true)
	}
	defer p.Unlock()

	ciphertext, err := base64.RawURLEncoding.DecodeString(ciphertextString)
	if err != nil {
		return nil, err
	}
	importKey, err := b.decryptImportedKey(ctx, req.Storage, ciphertext)
	err = p.Import(ctx, req.Storage, importKey, b.GetRandomReader())
	if err != nil {
		return nil, err
	}

	return nil, nil
}

const pathImportVersionWriteSyn = ""
const pathImportVersionWriteDesc = ""
