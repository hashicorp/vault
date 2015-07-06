package transit

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathRaw() *framework.Path {
	return &framework.Path{
		Pattern: `raw/(?P<name>\w+)`,
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the key",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: pathRawRead,
		},

		HelpSynopsis:    pathPolicyHelpSyn,
		HelpDescription: pathPolicyHelpDesc,
	}
}

func pathRawRead(
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
			"derived":     p.Derived,
		},
	}
	if p.Derived {
		resp.Data["kdf_mode"] = p.KDFMode
	}
	return resp, nil
}

const pathRawHelpSyn = `Fetch raw keys for named encrption keys`

const pathRawHelpDesc = `
This path is used to get the underlying encryption keys used for the
named keys that are available.
`
