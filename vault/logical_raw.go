// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
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

	// Preserve pre-existing behavior to decompress if `compressed` is missing
	compressed := true
	if d, ok := data.GetOk("compressed"); ok {
		compressed = d.(bool)
	}

	encoding := data.Get("encoding").(string)
	if encoding != "" && encoding != "base64" {
		return logical.ErrorResponse("invalid encoding %q", encoding), logical.ErrInvalidRequest
	}

	if b.recoveryMode {
		b.logger.Info("reading", "path", path)
	}

	// Prevent access of protected paths
	for _, p := range protectedPaths {
		if strings.HasPrefix(path, p) {
			err := fmt.Sprintf("cannot read %q", path)
			return logical.ErrorResponse(err), logical.ErrInvalidRequest
		}
	}

	// Run additional checks if needed
	if err := b.checkRaw(path); err != nil {
		b.logger.Warn(err.Error(), "path", path)
		return logical.ErrorResponse("cannot read %q", path), logical.ErrInvalidRequest
	}

	entry, err := b.barrier.Get(ctx, path)
	if err != nil {
		return handleErrorNoReadOnlyForward(err)
	}
	if entry == nil {
		return nil, nil
	}

	valueBytes := entry.Value
	if compressed {
		// Run this through the decompression helper to see if it's been compressed.
		// If the input contained the compression canary, `valueBytes` will hold
		// the decompressed data. If the input was not compressed, then `valueBytes`
		// will be nil.
		valueBytes, _, err = compressutil.Decompress(entry.Value)
		if err != nil {
			return handleErrorNoReadOnlyForward(err)
		}

		// `valueBytes` is nil if the input is uncompressed. In that case set it to the original input.
		if valueBytes == nil {
			valueBytes = entry.Value
		}
	}

	var value interface{} = string(valueBytes)
	// Golang docs (https://pkg.go.dev/encoding/json#Marshal), []byte encodes as a base64-encoded string
	if encoding == "base64" {
		value = valueBytes
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"value": value,
		},
	}
	return resp, nil
}

// handleRawWrite is used to write directly to the barrier
func (b *RawBackend) handleRawWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	path := data.Get("path").(string)
	compressionType := ""
	c, compressionTypeOk := data.GetOk("compression_type")
	if compressionTypeOk {
		compressionType = c.(string)
	}

	encoding := data.Get("encoding").(string)
	if encoding != "" && encoding != "base64" {
		return logical.ErrorResponse("invalid encoding %q", encoding), logical.ErrInvalidRequest
	}

	if b.recoveryMode {
		b.logger.Info("writing", "path", path)
	}

	// Prevent access of protected paths
	for _, p := range protectedPaths {
		if strings.HasPrefix(path, p) {
			err := fmt.Sprintf("cannot write %q", path)
			return logical.ErrorResponse(err), logical.ErrInvalidRequest
		}
	}

	v := data.Get("value").(string)
	value := []byte(v)
	if encoding == "base64" {
		var err error
		value, err = base64.StdEncoding.DecodeString(v)
		if err != nil {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
	}

	if req.Operation == logical.UpdateOperation {
		// Check if this is an existing value with compression applied, if so, use the same compression (or no compression)
		entry, err := b.barrier.Get(ctx, path)
		if err != nil {
			return handleErrorNoReadOnlyForward(err)
		}
		if entry == nil {
			err := fmt.Sprintf("cannot figure out compression type because entry does not exist")
			return logical.ErrorResponse(err), logical.ErrInvalidRequest
		}

		// For cases where DecompressWithCanary errored, treat entry as non-compressed data.
		_, existingCompressionType, _, _ := compressutil.DecompressWithCanary(entry.Value)

		// Ensure compression_type matches existing entries' compression
		// except allow writing non-compressed data over compressed data
		if existingCompressionType != compressionType && compressionType != "" {
			err := fmt.Sprintf("the entry uses a different compression scheme then compression_type")
			return logical.ErrorResponse(err), logical.ErrInvalidRequest
		}

		if !compressionTypeOk {
			compressionType = existingCompressionType
		}
	}

	if compressionType != "" {
		var config *compressutil.CompressionConfig
		switch compressionType {
		case compressutil.CompressionTypeLZ4:
			config = &compressutil.CompressionConfig{
				Type: compressutil.CompressionTypeLZ4,
			}
			break
		case compressutil.CompressionTypeLZW:
			config = &compressutil.CompressionConfig{
				Type: compressutil.CompressionTypeLZW,
			}
			break
		case compressutil.CompressionTypeGzip:
			config = &compressutil.CompressionConfig{
				Type:                 compressutil.CompressionTypeGzip,
				GzipCompressionLevel: gzip.BestCompression,
			}
			break
		case compressutil.CompressionTypeSnappy:
			config = &compressutil.CompressionConfig{
				Type: compressutil.CompressionTypeSnappy,
			}
			break
		default:
			err := fmt.Sprintf("invalid compression type %q", compressionType)
			return logical.ErrorResponse(err), logical.ErrInvalidRequest
		}

		var err error
		value, err = compressutil.Compress(value, config)
		if err != nil {
			return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
		}
	}

	entry := &logical.StorageEntry{
		Key:   path,
		Value: value,
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
			err := fmt.Sprintf("cannot delete %q", path)
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
			err := fmt.Sprintf("cannot list %q", path)
			return logical.ErrorResponse(err), logical.ErrInvalidRequest
		}
	}

	// Run additional checks if needed
	if err := b.checkRaw(path); err != nil {
		b.logger.Warn(err.Error(), "path", path)
		return logical.ErrorResponse("cannot list %q", path), logical.ErrInvalidRequest
	}

	keys, err := b.barrier.List(ctx, path)
	if err != nil {
		return handleErrorNoReadOnlyForward(err)
	}
	return logical.ListResponse(keys), nil
}

// existenceCheck checks if entry exists, used in handleRawWrite for update or create operations
func (b *RawBackend) existenceCheck(ctx context.Context, request *logical.Request, data *framework.FieldData) (bool, error) {
	path := data.Get("path").(string)
	entry, err := b.barrier.Get(ctx, path)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

func rawPaths(prefix string, r *RawBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: prefix + "raw/" + framework.MatchAllRegex("path"),

			Fields: map[string]*framework.FieldSchema{
				"path": {
					Type: framework.TypeString,
				},
				"value": {
					Type: framework.TypeString,
				},
				"compressed": {
					Type: framework.TypeBool,
				},
				"encoding": {
					Type: framework.TypeString,
				},
				"compression_type": {
					Type: framework.TypeString,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: r.handleRawRead,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationPrefix: "raw",
						OperationVerb:   "read",
					},
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
							Fields: map[string]*framework.FieldSchema{
								"value": {
									Type:     framework.TypeString,
									Required: true,
								},
							},
						}},
					},
					Summary: "Read the value of the key at the given path.",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: r.handleRawWrite,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationPrefix: "raw",
						OperationVerb:   "write",
					},
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
						}},
					},
					Summary: "Update the value of the key at the given path.",
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: r.handleRawWrite,
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
					Summary: "Create a key with value at the given path.",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: r.handleRawDelete,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationPrefix: "raw",
						OperationVerb:   "delete",
					},
					Responses: map[int][]framework.Response{
						http.StatusNoContent: {{
							Description: "OK",
						}},
					},
					Summary: "Delete the key with given path.",
				},
				logical.ListOperation: &framework.PathOperation{
					Callback: r.handleRawList,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationPrefix: "raw",
						OperationVerb:   "list",
					},
					Summary: "Return a list keys for a given path prefix.",
				},
			},

			ExistenceCheck:  r.existenceCheck,
			HelpSynopsis:    strings.TrimSpace(sysHelp["raw"][0]),
			HelpDescription: strings.TrimSpace(sysHelp["raw"][1]),
		},
	}
}
