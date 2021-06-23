package vault

import (
	"context"
	"fmt"
	"strings"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/compressutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// protectedPaths cannot be accessed via the raw APIs.
// This is both for security and to prevent disrupting Vault.
var protectedPaths = []string{
	keyringPath,
	// Changing the cluster info path can change the cluster ID which can be disruptive
	coreLocalClusterInfoPath,
}

type RawBackend struct {
	*framework.Backend
	barrier      SecurityBarrier
	logger       log.Logger
	checkRaw     func(path string) error
	recoveryMode bool
}

func NewRawBackend(core *Core) *RawBackend {
	r := &RawBackend{
		barrier: core.barrier,
		logger:  core.logger.Named("raw"),
		checkRaw: func(path string) error {
			return nil
		},
		recoveryMode: core.recoveryMode,
	}
	r.Backend = &framework.Backend{
		Paths: rawPaths("sys/", r),
	}
	return r
}

// handleRawRead is used to read directly from the barrier
func (b *RawBackend) handleRawRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)

	if b.recoveryMode {
		b.logger.Info("reading", "path", path)
	}

	// Prevent access of protected paths
	for _, p := range protectedPaths {
		if strings.HasPrefix(path, p) {
			err := fmt.Sprintf("cannot read '%s'", path)
			return logical.ErrorResponse(err), logical.ErrInvalidRequest
		}
	}

	// Run additional checks if needed
	if err := b.checkRaw(path); err != nil {
		b.logger.Warn(err.Error(), "path", path)
		return logical.ErrorResponse("cannot read '%s'", path), logical.ErrInvalidRequest
	}

	entry, err := b.barrier.Get(ctx, path)
	if err != nil {
		return handleErrorNoReadOnlyForward(err)
	}
	if entry == nil {
		return nil, nil
	}

	// Run this through the decompression helper to see if it's been compressed.
	// If the input contained the compression canary, `outputBytes` will hold
	// the decompressed data. If the input was not compressed, then `outputBytes`
	// will be nil.
	outputBytes, _, err := compressutil.Decompress(entry.Value)
	if err != nil {
		return handleErrorNoReadOnlyForward(err)
	}

	// `outputBytes` is nil if the input is uncompressed. In that case set it to the original input.
	if outputBytes == nil {
		outputBytes = entry.Value
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"value": string(outputBytes),
		},
	}
	return resp, nil
}

// handleRawWrite is used to write directly to the barrier
func (b *RawBackend) handleRawWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)

	if b.recoveryMode {
		b.logger.Info("writing", "path", path)
	}

	// Prevent access of protected paths
	for _, p := range protectedPaths {
		if strings.HasPrefix(path, p) {
			err := fmt.Sprintf("cannot write '%s'", path)
			return logical.ErrorResponse(err), logical.ErrInvalidRequest
		}
	}

	value := data.Get("value").(string)
	entry := &logical.StorageEntry{
		Key:   path,
		Value: []byte(value),
	}
	if err := b.barrier.Put(ctx, entry); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}
	return nil, nil
}

// handleRawDelete is used to delete directly from the barrier
func (b *RawBackend) handleRawDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)

	if b.recoveryMode {
		b.logger.Info("deleting", "path", path)
	}

	// Prevent access of protected paths
	for _, p := range protectedPaths {
		if strings.HasPrefix(path, p) {
			err := fmt.Sprintf("cannot delete '%s'", path)
			return logical.ErrorResponse(err), logical.ErrInvalidRequest
		}
	}

	if err := b.barrier.Delete(ctx, path); err != nil {
		return handleErrorNoReadOnlyForward(err)
	}
	return nil, nil
}

// handleRawList is used to list directly from the barrier
func (b *RawBackend) handleRawList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	if path != "" && !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	if b.recoveryMode {
		b.logger.Info("listing", "path", path)
	}

	// Prevent access of protected paths
	for _, p := range protectedPaths {
		if strings.HasPrefix(path, p) {
			err := fmt.Sprintf("cannot list '%s'", path)
			return logical.ErrorResponse(err), logical.ErrInvalidRequest
		}
	}

	// Run additional checks if needed
	if err := b.checkRaw(path); err != nil {
		b.logger.Warn(err.Error(), "path", path)
		return logical.ErrorResponse("cannot list '%s'", path), logical.ErrInvalidRequest
	}

	keys, err := b.barrier.List(ctx, path)
	if err != nil {
		return handleErrorNoReadOnlyForward(err)
	}
	return logical.ListResponse(keys), nil
}

func rawPaths(prefix string, r *RawBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: prefix + "(raw/?$|raw/(?P<path>.+))",

			Fields: map[string]*framework.FieldSchema{
				"path": {
					Type: framework.TypeString,
				},
				"value": {
					Type: framework.TypeString,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: r.handleRawRead,
					Summary:  "Read the value of the key at the given path.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: r.handleRawWrite,
					Summary:  "Update the value of the key at the given path.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: r.handleRawDelete,
					Summary:  "Delete the key with given path.",
				},
				logical.ListOperation: &framework.PathOperation{
					Callback: r.handleRawList,
					Summary:  "Return a list keys for a given path prefix.",
				},
			},

			HelpSynopsis:    strings.TrimSpace(sysHelp["raw"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["raw"][1]),
		},
	}
}
