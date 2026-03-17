// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateReq_Run(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	for name, test := range map[string]struct {
		cfg      func(t *testing.T, ctx context.Context) *ValidateReq
		validate func(t *testing.T, res *ValidateRes, err error)
	}{
		"with provided config": {
			cfg: func(t *testing.T, ctx context.Context) *ValidateReq {
				cfg := &Config{}
				decodeRes := &DecodeRes{
					Config: cfg,
				}
				return &ValidateReq{
					DecodeRes: decodeRes,
				}
			},
			validate: func(t *testing.T, res *ValidateRes, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotNil(t, res.Config)
			},
		},
		"with valid path": {
			cfg: func(t *testing.T, ctx context.Context) *ValidateReq {
				// Use existing test fixture
				path := filepath.Join("fixtures", "pipeline.hcl")
				decodeRes := Decode(ctx, &DecodeReq{
					Path: path,
				})
				return &ValidateReq{
					DecodeRes: decodeRes,
				}
			},
			validate: func(t *testing.T, res *ValidateRes, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotNil(t, res.Config)
			},
		},
		"with invalid path": {
			cfg: func(t *testing.T, ctx context.Context) *ValidateReq {
				decodeRes := Decode(ctx, &DecodeReq{
					Path: "/nonexistent/path/pipeline.hcl",
				})
				return &ValidateReq{
					DecodeRes: decodeRes,
				}
			},
			validate: func(t *testing.T, res *ValidateRes, err error) {
				require.Error(t, err)
				require.Nil(t, res)
			},
		},
		"with no decode result": {
			cfg: func(t *testing.T, ctx context.Context) *ValidateReq {
				return &ValidateReq{}
			},
			validate: func(t *testing.T, res *ValidateRes, err error) {
				require.Error(t, err)
				require.Nil(t, res)
			},
		},
		"with nil request": {
			cfg: func(t *testing.T, ctx context.Context) *ValidateReq {
				return nil
			},
			validate: func(t *testing.T, res *ValidateRes, err error) {
				require.Error(t, err)
				require.Nil(t, res)
			},
		},
		"with malformed config file": {
			cfg: func(t *testing.T, ctx context.Context) *ValidateReq {
				// Create a temporary malformed config file
				tmpDir := t.TempDir()
				tmpFile := filepath.Join(tmpDir, "bad_pipeline.hcl")
				err := os.WriteFile(tmpFile, []byte("invalid { hcl syntax"), 0o644)
				require.NoError(t, err)

				decodeRes := Decode(ctx, &DecodeReq{
					Path: tmpFile,
				})
				return &ValidateReq{
					DecodeRes: decodeRes,
				}
			},
			validate: func(t *testing.T, res *ValidateRes, err error) {
				require.Error(t, err)
				require.Nil(t, res)
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := test.cfg(t, ctx)
			res, err := req.Run(ctx)
			test.validate(t, res, err)
		})
	}
}

func TestValidateRes_ToJSON(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		res   *ValidateRes
		valid bool
	}{
		"valid": {
			&ValidateRes{
				Config: &Config{},
			},
			true,
		},
		"nil": {
			nil,
			false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if test.valid {
				b, err := test.res.ToJSON()
				require.NoError(t, err)
				require.NotNil(t, b)
			} else {
				b, err := test.res.ToJSON()
				require.Error(t, err)
				require.Nil(t, b)
			}
		})
	}
}
