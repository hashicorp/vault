// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package puller

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

var _ pluginSource = (*fsPluginSource)(nil)

type fsPluginSource struct {
	basePath string
}

func (p *fsPluginSource) listMetadata(ctx context.Context, plugin string) ([]metadata, error) {
	path := filepath.Join(p.basePath, "v1", "releases", plugin)
	return fsGet[[]metadata](ctx, path)
}

func (p *fsPluginSource) getMetadata(ctx context.Context, plugin, version string) (metadata, error) {
	path := filepath.Join(p.basePath, "v1", "releases", plugin, version)
	return fsGet[metadata](ctx, path)
}

func fsGet[T any](_ context.Context, path string) (T, error) {
	var t T
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return t, errors.Join(errNotFound, err)
		}

		return t, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	if err = dec.Decode(&t); err != nil {
		return t, err
	}

	return t, nil
}

func (p *fsPluginSource) getContentReader(ctx context.Context, path string) (reader io.ReadCloser, err error) {
	return os.Open(path)
}
