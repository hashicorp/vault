package transit

import (
	"crypto/rand"
	"encoding/json"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// Policy is the struct used to store metadata
type Policy struct {
	Name       string `json:"name"`
	Key        []byte `json:"key"`
	CipherMode string `json:"cipher"`
}

func (p *Policy) Serialize() ([]byte, error) {
	return json.Marshal(p)
}

func DeserializePolicy(buf []byte) (*Policy, error) {
	p := new(Policy)
	if err := json.Unmarshal(buf, p); err != nil {
		return nil, err
	}
	return p, nil
}

func getPolicy(req *logical.Request, name string) (*Policy, error) {
	// Check if the policy already exists
	raw, err := req.Storage.Get("policy/" + name)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}

	// Decode the policy
	p, err := DeserializePolicy(raw.Value)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func pathKeys() *framework.Path {
	return &framework.Path{
		Pattern: `keys/(?P<name>\w+)`,
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the key",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation:  pathPolicyWrite,
			logical.DeleteOperation: pathPolicyDelete,
			logical.ReadOperation:   pathPolicyRead,
		},

		HelpSynopsis:    pathPolicyHelpSyn,
		HelpDescription: pathPolicyHelpDesc,
	}
}

func pathPolicyWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	// Check if the policy already exists
	existing, err := getPolicy(req, name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, nil
	}

	// Create the policy object
	p := &Policy{
		Name:       name,
		CipherMode: "aes-gcm",
	}

	// Generate a 256bit key
	p.Key = make([]byte, 32)
	_, err = rand.Read(p.Key)
	if err != nil {
		return nil, err
	}

	// Encode the policy
	buf, err := p.Serialize()
	if err != nil {
		return nil, err
	}

	// Write the policy into storage
	err = req.Storage.Put(&logical.StorageEntry{
		Key:   "policy/" + name,
		Value: buf,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func pathPolicyRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	p, err := getPolicy(req, name)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}

	// Return the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"name":        p.Name,
			"key":         p.Key,
			"cipher_mode": p.CipherMode,
		},
	}
	return resp, nil
}

func pathPolicyDelete(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	err := req.Storage.Delete("policy/" + name)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

const pathPolicyHelpSyn = `Managed named encrption keys`

const pathPolicyHelpDesc = `
This path is used to manage the named keys that are available.
Doing a write with no value against a new named key will create
it using a randomly generated key.
`
