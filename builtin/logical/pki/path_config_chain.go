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
		},

		HelpSynopsis:    pathConfigChainHelpSyn,
		HelpDescription: pathConfigChainHelpDesc,
	}
}

func (b *backend) pathChainWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	caChain := data.Get("ca_chain").(string)

	cb := &certutil.CertBundle{
		IssuingCAChain: caChain,
	}
	_, err := cb.ToParsedCertBundle()
	if err != nil {
		return nil, fmt.Errorf("Unable to parse CA chain: %s", err)
	}

	config := &caChainConfig{
		CAChain: caChain,
	}

	entry, err := logical.StorageEntryJSON("config/ca_chain", config)
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
	entry, err := req.Storage.Get("config/ca_chain")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var chainConfig caChainConfig
	if err := entry.DecodeJSON(&chainConfig); err != nil {
		return nil, err
	}

	return &chainConfig, nil
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

const pathConfigChainHelpSyn = `
TODO
`

const pathConfigChainHelpDesc = `
TODO
`
