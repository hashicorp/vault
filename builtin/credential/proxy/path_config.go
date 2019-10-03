package proxy

import (
	"context"
	"fmt"
	"net/textproto"

	"github.com/hashicorp/errwrap"
	sockaddr "github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	configStoragePath     string = "config"
	configBoundCidrsField        = "bound_cidrs"
	configUserHeaderField        = "user_header"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			configUserHeaderField: &framework.FieldSchema{
				Type:     framework.TypeString,
				Required: true,
				Description: `The name of the header that contains the authenticated user’s` +
					`username.  Case insensitive.`,
			},
			configBoundCidrsField: &framework.FieldSchema{
				Type: framework.TypeCommaStringSlice,
				Description: `Comma separated string or list of CIDR blocks. If ` +
					`set, restricts the blocks of IP addresses which can perform the ` +
					`login operation.`,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigRead,
				Summary:  "Read the current proxy authentication backend configuration.",
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathConfigCreateUpdate,
				Summary:  "Set the current proxy authentication backend configuration.",
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigCreateUpdate,
				Summary:  "Update the current proxy authentication backend configuration.",
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathConfigDelete,
				Summary:  "Delete the current proxy authentication backend configuration.",
			},
		},

		ExistenceCheck: b.pathConfigExistenceCheck,

		HelpSynopsis: "Manage the proxy authentication backend configuration",
	}
}

type proxyConfig struct {
	UserHeader string                        `json:"user_header"`
	BoundCIDRs []*sockaddr.SockAddrMarshaler `json:"bound_cidrs"`
}

// config fetchs the proxyConfig from the storage backend
func (b *backend) config(ctx context.Context, s logical.Storage) (*proxyConfig, error) {
	entry, err := s.Get(ctx, configStoragePath)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	config := proxyConfig{}
	if err := entry.DecodeJSON(&config); err != nil {
		return nil, errwrap.Wrapf("error reading proxy backend configuration: {{err}}", err)
	}
	return &config, nil
}

func (b *backend) pathConfigExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return false, err
	}

	return config != nil, nil
}

func (b *backend) pathConfigCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if config == nil {
		config = &proxyConfig{}
	}

	if val, ok := data.GetOk(configUserHeaderField); ok {
		config.UserHeader = textproto.CanonicalMIMEHeaderKey(val.(string))
	}
	if config.UserHeader == "" {
		return logical.ErrorResponse(fmt.Sprintf("%s must be set", configUserHeaderField)), nil
	}

	if val, ok := data.GetOk(configBoundCidrsField); ok {
		parsedCIDRs, err := parseutil.ParseAddrs(val)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("%s is invalid: %s", configBoundCidrsField, err.Error())), nil
		}
		config.BoundCIDRs = parsedCIDRs
	}

	entry, err := logical.StorageEntryJSON(configStoragePath, config)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			configUserHeaderField: config.UserHeader,
			configBoundCidrsField: config.BoundCIDRs,
		},
	}
	return resp, nil
}

func (b *backend) pathConfigDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, configStoragePath)
	return nil, err
}
