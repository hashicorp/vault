package config

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/activedirectory"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// operationHandler receives inbound calls from our API or CLI.
// It's a high-level controller that can access both the config cache,
// and underlying storage.
type operationHandler struct {
	logger hclog.Logger
	cache  *cache
}

func (h *operationHandler) Delete(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	if err := deleteConfig(ctx, req.Storage); err != nil {
		return nil, err
	}
	h.cache.Set(newUnsetEngineConf())
	return nil, nil
}

func (h *operationHandler) Read(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {

	engineConf, ok := h.cache.Get()
	if !ok {
		var err error
		engineConf, err = readConfig(ctx, req.Storage)
		if err != nil {
			return nil, err
		}
	}

	resp := &logical.Response{
		Data: engineConf.Map(),
	}
	resp.AddWarning("read access to this endpoint should be controlled via ACLs as it will return the configuration information as-is, including any passwords.")
	return resp, nil
}

func (h *operationHandler) Update(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (*logical.Response, error) {

	// Parse and validate the desired conf.
	activeDirectoryConf, err := activedirectory.NewConfiguration(h.logger, fieldData)
	if err != nil {
		return nil, err
	}

	passwordConf := newPasswordConfig(fieldData)

	engineConf, err := newEngineConf(passwordConf, activeDirectoryConf)
	if err != nil {
		return nil, err
	}

	// Write and cache it.
	if err := writeConfig(ctx, req.Storage, engineConf); err != nil {
		return nil, err
	}
	h.cache.Set(engineConf)

	// Respond.
	resp := &logical.Response{
		Data: engineConf.Map(),
	}
	resp.AddWarning("write access to this endpoint should be controlled via ACLs as it will return the configuration information as-is, including any passwords.")
	return resp, nil
}
