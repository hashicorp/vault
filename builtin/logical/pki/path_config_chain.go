package pki

import (
	"fmt"

	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigChain(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/chain",
		Fields: map[string]*framework.FieldSchema{
			"pem_bundle": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `PEM-format, concatenated certificates
for the CA trust chain.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathChainWrite,
			logical.ReadOperation:   b.pathChainRead,
		},

		HelpSynopsis:    pathConfigChainHelpSyn,
		HelpDescription: pathConfigChainHelpDesc,
	}
}

func (b *backend) pathChainWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	pemBundle := data.Get("pem_bundle").(string)

	parsedCAChain, err := certutil.ParsePEMBundle(pemBundle)
	if err != nil {
		switch err.(type) {
		case errutil.InternalError:
			return nil, err
		default:
			return logical.ErrorResponse(err.Error()), nil
		}
	}

	caChain, err := parsedCAChain.ToCertBundle()
	if err != nil {
		return nil, fmt.Errorf("error converting raw values into cert bundle: %s", err)
	}

	caBundle, err := fetchCABundle(req)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch CA bundle: %v", err)}
	}

	caBundle.IssuingCA = caChain.Certificate
	if len(caChain.IssuingCAChain) > 0 {
		caBundle.IssuingCAChain = fmt.Sprintf("%s/n%s", caChain.IssuingCA, caChain.IssuingCAChain)
	} else {
		caBundle.IssuingCAChain = caChain.IssuingCA
	}

	parsedBundle, err := caBundle.ToParsedCertBundle()
	if err != nil {
		return nil, errutil.InternalError{Err: err.Error()}
	}

	if err := parsedBundle.Verify(); err != nil {
		return nil, fmt.Errorf("verification of parsed bundle failed: %s", err)
	}

	cb, err := parsedBundle.ToCertBundle()
	if err != nil {
		return nil, fmt.Errorf("error converting raw values into cert bundle: %s", err)
	}

	entry, err := logical.StorageEntryJSON("config/ca_bundle", cb)
	if err != nil {
		return nil, err
	}
	err = req.Storage.Put(entry)
	if err != nil {
		return nil, err
	}

	return nil, err
}

func (b *backend) pathChainRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	caBundle, err := fetchCABundle(req)
	if err != nil {
		return nil, err
	}
	if caBundle == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"issuing_ca":       caBundle.IssuingCA,
			"issuing_ca_chain": caBundle.IssuingCAChain,
		},
	}

	if s, ok := resp.Data["issuing_ca_chain"].(string); ok && len(s) == 0 {
		delete(resp.Data, "issuing_ca_chain")
	}

	return resp, nil
}

const pathConfigChainHelpSyn = `
Configure the certificate authority trust chain.
`

const pathConfigChainHelpDesc = `
This endpoint allows configuration of the trust chain for the certificate
authority.  By populating the trust chain, this information will be returned
when issuing certificates and will be returned when requesting pem bundles.

Multiple certificates can be concatenated into a single file in order from the
issuing certificate authority.  Because certificate validation requires that
root keys be distributed independently, the root certificate authority should
be omitted from the trust chain.
`
