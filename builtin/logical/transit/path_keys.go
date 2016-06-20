package transit

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathKeys() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the key",
			},

			"derived": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Description: "Enables key derivation mode. This allows for per-transaction unique keys",
			},

			"convergent_encryption": &framework.FieldSchema{
				Type: framework.TypeBool,
				Description: `Whether to use convergent encryption.
This is only supported when using a key with
key derivation enabled and will require all
context values to be 96 bits (12 bytes) when
base64-decoded. This mode ensures that when
the same context is supplied, the same
ciphertext is emitted from the encryption
function. It is *very important* when using
this mode that you ensure that all contexts
are *globally unique*. Failing to do so will
severely impact the security of the key.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathPolicyWrite,
			logical.DeleteOperation: b.pathPolicyDelete,
			logical.ReadOperation:   b.pathPolicyRead,
		},

		HelpSynopsis:    pathPolicyHelpSyn,
		HelpDescription: pathPolicyHelpDesc,
	}
}

func (b *backend) pathPolicyWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	derived := d.Get("derived").(bool)
	convergent := d.Get("convergent_encryption").(bool)

	if !derived && convergent {
		return logical.ErrorResponse("convergent encryption requires derivation to be enabled"), nil
	}

	p, lock, upserted, err := b.lm.GetPolicyUpsert(req.Storage, name, derived, convergent)
	if lock != nil {
		defer lock.RUnlock()
	}
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("error generating key: returned policy was nil")
	}

	resp := &logical.Response{}
	if !upserted {
		resp.AddWarning(fmt.Sprintf("key %s already existed", name))
	}

	return nil, nil
}

func (b *backend) pathPolicyRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	p, lock, err := b.lm.GetPolicyShared(req.Storage, name)
	if lock != nil {
		defer lock.RUnlock()
	}
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}

	// Return the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"name":                   p.Name,
			"cipher_mode":            p.CipherMode,
			"derived":                p.Derived,
			"deletion_allowed":       p.DeletionAllowed,
			"min_decryption_version": p.MinDecryptionVersion,
			"latest_version":         p.LatestVersion,
		},
	}
	if p.Derived {
		resp.Data["kdf_mode"] = p.KDFMode
		resp.Data["convergent_encryption"] = p.ConvergentEncryption
	}

	retKeys := map[string]int64{}
	for k, v := range p.Keys {
		retKeys[strconv.Itoa(k)] = v.CreationTime
	}
	resp.Data["keys"] = retKeys

	return resp, nil
}

func (b *backend) pathPolicyDelete(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	// Delete does its own locking
	err := b.lm.DeletePolicy(req.Storage, name)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("error deleting policy %s: %s", name, err)), err
	}

	return nil, nil
}

const pathPolicyHelpSyn = `Managed named encryption keys`

const pathPolicyHelpDesc = `
This path is used to manage the named keys that are available.
Doing a write with no value against a new named key will create
it using a randomly generated key.
`
