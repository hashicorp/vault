package transit

import (
	"crypto/rand"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/vault/helper/kdf"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	// kdfMode is the only KDF mode currently supported
	kdfMode = "hmac-sha256-counter"
)

// Policy is the struct used to store metadata
type Policy struct {
	Name       string `json:"name"`
	Key        []byte `json:"key"`
	CipherMode string `json:"cipher"`

	// Derived keys MUST provide a context and the
	// master underlying key is never used.
	Derived bool   `json:"derived"`
	KDFMode string `json:"kdf_mode"`
}

func (p *Policy) Serialize() ([]byte, error) {
	return json.Marshal(p)
}

// DeriveKey is used to derive the encryption key that should
// be used depending on the policy. If derivation is disabled the
// raw key is used and no context is required, otherwise the KDF
// mode is used with the context to derive the proper key.
func (p *Policy) DeriveKey(context []byte) ([]byte, error) {
	// Fast-path non-derived keys
	if !p.Derived {
		return p.Key, nil
	}

	// Ensure a context is provided
	if len(context) == 0 {
		return nil, fmt.Errorf("missing 'context' for key deriviation. The key was created using a derived key, which means additional, per-request information must be included in order to encrypt or decrypt information.")
	}

	switch p.KDFMode {
	case kdfMode:
		prf := kdf.HMACSHA256PRF
		prfLen := kdf.HMACSHA256PRFLen
		return kdf.CounterMode(prf, prfLen, p.Key, context, 256)
	default:
		return nil, fmt.Errorf("unsupported key derivation mode")
	}
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

// generatePolicy is used to create a new named policy with
// a randomly generated key
func generatePolicy(storage logical.Storage, name string, derived bool) (*Policy, error) {
	// Create the policy object
	p := &Policy{
		Name:       name,
		CipherMode: "aes-gcm",
		Derived:    derived,
	}
	if derived {
		p.KDFMode = kdfMode
	}

	// Generate a 256bit key
	p.Key = make([]byte, 32)
	_, err := rand.Read(p.Key)
	if err != nil {
		return nil, err
	}

	// Encode the policy
	buf, err := p.Serialize()
	if err != nil {
		return nil, err
	}

	// Write the policy into storage
	err = storage.Put(&logical.StorageEntry{
		Key:   "policy/" + name,
		Value: buf,
	})
	if err != nil {
		return nil, err
	}

	// Return the policy
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

			"derived": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Description: "Enables key derivation mode. This allows for per-transaction unique keys",
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
	derived := d.Get("derived").(bool)

	// Check if the policy already exists
	existing, err := getPolicy(req, name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, nil
	}

	// Generate the policy
	_, err = generatePolicy(req.Storage, name, derived)
	return nil, err
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
			"cipher_mode": p.CipherMode,
			"derived":     p.Derived,
		},
	}
	if p.Derived {
		resp.Data["kdf_mode"] = p.KDFMode
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
