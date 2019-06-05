package transit

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathCacheConfig() *framework.Path {
	return &framework.Path{
		Pattern: "cache-config",
		Fields: map[string]*framework.FieldSchema{
			"size": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Required:    false,
				Default:     0,
				Description: `Size of cache, use 0 for an unlimited cache size, defaults to 0`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathCacheConfigRead,
				Summary:  "Returns the size of the active cache",
			},

			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathCacheConfigWrite,
				Summary:  "Configures a new cache of the specified size",
			},

			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathCacheConfigWrite,
				Summary:  "Configures a new cache of the specified size",
			},
		},

		HelpSynopsis:    pathCacheConfigHelpSyn,
		HelpDescription: pathCacheConfigHelpDesc,
	}
}

func (b *backend) pathCacheConfigWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// get target size
	cacheSize := d.Get("size").(int)
	if cacheSize < 0 {
		return logical.ErrorResponse("size must be greater or equal to 0"), logical.ErrInvalidRequest
	}

	// store cache size
	entry, err := logical.StorageEntryJSON("config/cache", &configCache{
		Size: cacheSize,
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Warnings: []string{"cache configurations will be applied when this backend is restarted"},
	}

	return resp, nil
}

type configCache struct {
	Size int `json:"size"`
}

func (b *backend) pathCacheConfigRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// error if no cache is configured
	if !b.lm.GetUseCache() {
		return nil, errors.New(
			"caching is disabled for this transit mount",
		)
	}

	// Compare current and stored cache sizes. If they are different warn the user.
	currentCacheSize := b.lm.GetCacheSize()
	storedCacheSize, err := GetCacheSizeFromStorage(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"size": storedCacheSize,
		},
	}

	if currentCacheSize != storedCacheSize {
		resp.Warnings = []string{"This cache size will not be applied until the transit mount is reloaded"}
	}

	return resp, nil
}

const pathCacheConfigHelpSyn = `Configure caching strategy`

const pathCacheConfigHelpDesc = `
This path is used to configure and query the cache size of the active cache, a size of 0 means unlimited.
`
