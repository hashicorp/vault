package nomad

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const configAccessKey = "config/access"

func pathConfigAccess(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/access",
		Fields: map[string]*framework.FieldSchema{
			"address": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Nomad server address",
			},

			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Token for API calls",
			},

			"max_token_length": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: "Max length for generated Nomad tokens",
				// Default length is 256 as of
				// https://github.com/hashicorp/nomad/blob/21682427f3474f92cc589832efe72850a61c83a7/nomad/structs/structs.go#L116
				Default: maxTokenNameLength,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConfigAccessRead,
			logical.CreateOperation: b.pathConfigAccessWrite,
			logical.UpdateOperation: b.pathConfigAccessWrite,
			logical.DeleteOperation: b.pathConfigAccessDelete,
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
		return nil, errwrap.Wrapf("error reading nomad access configuration: {{err}}", err)
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
			"address":          conf.Address,
			"max_token_length": conf.MaxTokenLength,
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

	// max_token_length has default of 256
	conf.MaxTokenLength = data.Get("max_token_length").(int)
	envMaxTokenLength := os.Getenv("NOMAD_MAX_TOKEN_LENGTH")
	if envMaxTokenLength != "" {
		// if we find NOMAD_MAX_max_token_length in the env and can parse it, override
		// the default length
		i, err := strconv.Atoi(envMaxTokenLength)
		if err != nil {
			log.Printf("[WARN] error parsing NOMAD_MAX_TOKEN_LENGTH, using default 256")
		} else {
			conf.MaxTokenLength = i
		}
	}

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
	Address        string `json:"address"`
	Token          string `json:"token"`
	MaxTokenLength int    `json:"max_token_length"`
}
