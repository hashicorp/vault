// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclparse"
	slogctx "github.com/veqryn/slog-context"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type VersionsConfig struct {
	Schema        int            `json:"schema" cty:"schema" hcl:"schema,optional"`
	ActiveVersion *ActiveVersion `json:"active_versions" cty:"active_versions" hcl:"active_versions"`
}

type ActiveVersion struct {
	Versions map[string]*Version `json:"versions"`
}

type Version struct {
	CEActive bool `json:"ce_active"`
	LTS      bool `json:"lts"`
}

// DecodeReq contains the parameters for decoding a pipeline configuration.
type DecodeReq struct {
	// Path is the path to the pipeline.hcl file that we wish to decode.
	Path string `json:"path,omitempty"`
}

// DecodeRes contains the result of decoding a pipeline configuration.
type DecodeRes struct {
	// Path is the evaluated path to the configuration file that was loaded. It
	// will always be the absolute path.
	Path string `json:"path,omitempty"`

	// Config is the decoded versions configuration.
	Config *VersionsConfig `json:"config,omitempty"`

	// Err contains the error if decoding failed.
	Err error

	// ErrStr contains the error message if decoding failed.
	ErrStr string `json:"error,omitempty"`
}

var v1Schema = hcldec.ObjectSpec{
	"schema": &hcldec.AttrSpec{
		Name:     "schema",
		Type:     cty.Number,
		Required: false,
	},
	"active_versions": &hcldec.BlockSpec{
		TypeName: "active_versions",
		Nested: hcldec.ObjectSpec{
			"version": &hcldec.BlockMapSpec{
				TypeName:   "version",
				LabelNames: []string{"name"},
				Nested: hcldec.ObjectSpec{
					"ce_active": &hcldec.AttrSpec{
						Name:     "ce_active",
						Type:     cty.Bool,
						Required: false,
					},
					"lts": &hcldec.AttrSpec{
						Name:     "lts",
						Type:     cty.Bool,
						Required: false,
					},
				},
			},
		},
	},
}

// Load takes a context and variable opts and locates the pipeline configuration
// file, parses it, and returns the Config.
func Load(ctx context.Context, path string) (*VersionsConfig, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	ctx = slogctx.Append(ctx,
		slog.String("path", path),
	)

	slog.Default().DebugContext(ctx, "loading versions.hcl")
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	return DecodeFile(path)
}

// Decode takes a DecodeReq, uses the embedded information to load and decode
// the configuration, and returns all of the data in a DecodeRes.
// If you do not need the surrounding DecodeRes struct you can use Load()
// instead.
func Decode(ctx context.Context, req *DecodeReq) *DecodeRes {
	res := &DecodeRes{}

	// Determine the path to use
	if req.Path == "" {
		res.Err = errors.New("no path or repo root provided")
		res.ErrStr = res.Err.Error()
	}

	// Get absolute path
	absPath, err := filepath.Abs(req.Path)
	if err != nil {
		res.Err = err
		res.ErrStr = err.Error()

		return res
	}
	res.Path = absPath

	// Load the configuration
	cfg, err := Load(ctx, absPath)
	if err != nil {
		res.Err = err
		res.ErrStr = err.Error()

		return res
	}
	res.Config = cfg

	return res
}

// DecodeFile decodes the configuration file at the given path into a Config.
func DecodeFile(path string) (*VersionsConfig, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return DecodeBytes(bytes)
}

// DecodeBytes decodes the configuration from bytes into a Config.
func DecodeBytes(bytes []byte) (*VersionsConfig, error) {
	var err error
	var file *hcl.File
	var diags hcl.Diagnostics
	parser := hclparse.NewParser()

	file, diags = parser.ParseHCL(bytes, "versions.hcl")
	if diags != nil && diags.HasErrors() {
		for _, diag := range diags {
			err = errors.Join(err, errors.New(diag.Error()))
		}

		return nil, err
	}

	val, moreDiags := hcldec.Decode(file.Body, v1Schema, nil)
	if moreDiags != nil && moreDiags.HasErrors() {
		for _, diag := range moreDiags {
			err = errors.Join(err, errors.New(diag.Error()))
		}
		return nil, err
	}

	res := &VersionsConfig{ActiveVersion: &ActiveVersion{Versions: map[string]*Version{}}}
	if !val.IsWhollyKnown() || !val.CanIterateElements() {
		err = fmt.Errorf("unexpected version type: %s", val.Type().GoString())
		return nil, err
	}

	schema, ok := val.AsValueMap()["schema"]
	if ok && schema.IsWhollyKnown() && schema.Type().Equals(cty.Number) {
		err = gocty.FromCtyValue(schema, &res.Schema)
		if err != nil {
			return nil, err
		}
	}

	av, ok := val.AsValueMap()["active_versions"]
	if !ok {
		return res, errors.New("no active_versions stanza has been defined")
	}

	if av.IsNull() {
		return res, errors.New("active_versions stanza cannot be empty")
	}

	for _, versions := range av.AsValueSlice() {
		for version, versionVal := range versions.AsValueMap() {
			v := &Version{}

			for attr, val := range versionVal.AsValueMap() {
				switch attr {
				case "ce_active":
					if !val.IsNull() {
						v.CEActive = val.True()
					}
				case "lts":
					if !val.IsNull() {
						v.LTS = val.True()
					}
				default:
					err = fmt.Errorf("unknown value: %s", attr)
					return nil, err
				}
			}

			res.ActiveVersion.Versions[version] = v
		}
	}

	return res, err
}

// Validate validates that the response is initialized and does not have
// an embedded error.
func (d *DecodeRes) Validate(ctx context.Context) error {
	if d == nil {
		return errors.New("versions.hcl decode response is unitialized")
	}

	if d.Err != nil {
		return d.Err
	}

	return nil
}
