// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package nomad

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const configAccessKey = "config/access"

func pathConfigAccess(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/access",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixNomad,
		},

		Fields: map[string]*framework.FieldSchema{
			"address": {
				Type:        framework.TypeString,
				Description: "Nomad server address",
			},

			"token": {
				Type:        framework.TypeString,
				Description: "Token for API calls",
			},

			"max_token_name_length": {
				Type:        framework.TypeInt,
				Description: "Max length for name of generated Nomad tokens",
			},
			"ca_cert": {
				Type: framework.TypeString,
				Description: `CA certificate to use when verifying Nomad server certificate,
must be x509 PEM encoded.`,
			},
			"client_cert": {
				Type: framework.TypeString,
				Description: `Client certificate used for Nomad's TLS communication,
must be x509 PEM encoded and if this is set you need to also set client_key.`,
			},
			"client_key": {
				Type: framework.TypeString,
				Description: `Client key used for Nomad's TLS communication,
must be x509 PEM encoded and if this is set you need to also set client_cert.`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigAccessRead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "read",
					OperationSuffix: "access-configuration",
				},
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathConfigAccessWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "access",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigAccessWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "access",
				},
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathConfigAccessDelete,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "delete",
					OperationSuffix: "access-configuration",
				},
			},
		},

		ExistenceCheck: b.configExistenceCheck,
	}
}

func (b *backend) configExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := b.readConfigAccess(ctx, req.Storage)
	if err != nil {
		return false, err
	}

	return entry != nil, nil
}

func (b *backend) readConfigAccess(ctx context.Context, storage logical.Storage) (*accessConfig, error) {
	entry, err := storage.Get(ctx, configAccessKey)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	conf := &accessConfig{}
	if err := entry.DecodeJSON(conf); err != nil {
		return nil, fmt.Errorf("error reading nomad access configuration: %w", err)
	}

	return conf, nil
}

func (b *backend) pathConfigAccessRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	conf, err := b.readConfigAccess(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if conf == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"address":               conf.Address,
			"max_token_name_length": conf.MaxTokenNameLength,
			"ca_cert":               conf.CACert,
			"client_cert":           conf.ClientCert,
		},
	}, nil
}

func (b *backend) pathConfigAccessWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	conf, err := b.readConfigAccess(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if conf == nil {
		conf = &accessConfig{}
	}

	address, ok := data.GetOk("address")
	if ok {
		conf.Address = address.(string)
	}
	token, ok := data.GetOk("token")
	if ok {
		conf.Token = token.(string)
	}
	caCert, ok := data.GetOk("ca_cert")
	if ok {
		conf.CACert = caCert.(string)
	}
	clientCert, ok := data.GetOk("client_cert")
	if ok {
		conf.ClientCert = clientCert.(string)
	}
	clientKey, ok := data.GetOk("client_key")
	if ok {
		conf.ClientKey = clientKey.(string)
	}

	if conf.Token == "" {
		client, err := clientFromConfig(conf)
		if err != nil {
			return logical.ErrorResponse("Token not provided and failed to constuct client"), err
		}
		token, _, err := client.ACLTokens().Bootstrap(nil)
		if err != nil {
			return logical.ErrorResponse("Token not provided and failed to bootstrap ACLs"), err
		}
		conf.Token = token.SecretID
	}

	conf.MaxTokenNameLength = data.Get("max_token_name_length").(int)

	entry, err := logical.StorageEntryJSON("config/access", conf)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathConfigAccessDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, configAccessKey); err != nil {
		return nil, err
	}
	return nil, nil
}

type accessConfig struct {
	Address            string `json:"address"`
	Token              string `json:"token"`
	MaxTokenNameLength int    `json:"max_token_name_length"`
	CACert             string `json:"ca_cert"`
	ClientCert         string `json:"client_cert"`
	ClientKey          string `json:"client_key"`
}
