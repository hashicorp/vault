package pki

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

type caChainConfig struct {
	CAChain string `json:"ca_chain" mapstructure:"ca_chain" structs:"ca_chain"`
}

func pathConfigChain(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/chain",
		Fields: map[string]*framework.FieldSchema{
			"ca_chain": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `TODO`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathChainWrite,
			logical.ReadOperation:   b.pathChainRead,
			logical.DeleteOperation: b.pathChainDelete,
		},

		HelpSynopsis:    pathConfigChainHelpSyn,
		HelpDescription: pathConfigChainHelpDesc,
	}
}

func (b *backend) pathChainWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	caChain := data.Get("ca_chain").(string)

	cb := &certutil.CertBundle{}
	entry, err := req.Storage.Get("config/ca_bundle")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return logical.ErrorResponse("could not find any existing ca entry"), nil
	}

	cb.IssuingCAChain = caChain

	err = entry.DecodeJSON(cb)
	if err != nil {
		return nil, err
	}

	parsedCB, err := cb.ToParsedCertBundle()
	if err != nil {
		return nil, fmt.Errorf("error parsing cert bundle: %s", err)
	}

	cb, err = parsedCB.ToCertBundle()
	if err != nil {
		return nil, fmt.Errorf("error converting raw values into cert bundle: %s", err)
	}

	entry, err = logical.StorageEntryJSON("config/ca_bundle", cb)
	if err != nil {
		return nil, err
	}
	err = req.Storage.Put(entry)
	if err != nil {
		return nil, err
	}

	return nil, err
}

func getCAChain(req *logical.Request) (*caChainConfig, error) {
	entry, err := req.Storage.Get("config/ca_bundle")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	cb := &certutil.CertBundle{}
	if err := entry.DecodeJSON(&cb); err != nil {
		return nil, err
	}

	chainConfig := &caChainConfig{
		CAChain: cb.IssuingCAChain,
	}

	return chainConfig, nil
}

func (b *backend) pathChainRead(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := getCAChain(req)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: structs.New(entry).Map(),
	}

	return resp, nil
}

func (b *backend) pathChainDelete(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	cb := &certutil.CertBundle{}
	entry, err := req.Storage.Get("config/ca_bundle")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return logical.ErrorResponse("could not find any existing entry with a private key"), nil
	}

	err = entry.DecodeJSON(cb)
	if err != nil {
		return nil, err
	}

	cb.IssuingCAChain = ""

	entry, err = logical.StorageEntryJSON("config/ca_bundle", cb)
	if err != nil {
		return nil, err
	}
	err = req.Storage.Put(entry)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

const pathConfigChainHelpSyn = `
TODO
`

const pathConfigChainHelpDesc = `
TODO
`
