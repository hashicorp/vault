// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"slices"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/jedib0t/go-pretty/v6/table"
	slogctx "github.com/veqryn/slog-context"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type (
	VersionsConfig struct {
		Schema        int            `json:"schema" cty:"schema" hcl:"schema,optional"`
		ActiveVersion *ActiveVersion `json:"active_versions" cty:"active_versions" hcl:"active_versions"`
	}
	ActiveVersion struct {
		Versions map[string]*Version `json:"versions"`
	}
	Version struct {
		CEActive bool `json:"ce_active"`
		LTS      bool `json:"lts"`
	}
)

// ListActiveVersionsReq is a request to list the active branch versions from the
// .release/metadata file
type ListActiveVersionsReq struct {
	// ReleaseVersionConfigPath is the path to the .release/versions.hcl file
	ReleaseVersionConfigPath string
	// The depth to recursively search backwards for a .release/versions.hcl file
	Recurse uint
}

// ListActiveVersionsRes are the active versions and associated metadata for the repo
type ListActiveVersionsRes struct {
	VersionsConfig *VersionsConfig `json:"versions_config,omitempty"`
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

// Run runs the dynamic configuration request
func (l *ListActiveVersionsReq) Run(ctx context.Context) (*ListActiveVersionsRes, error) {
	if l == nil {
		return nil, fmt.Errorf("list active versions request is uninitialized")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var err error
	res := &ListActiveVersionsRes{}

	defer func() {
		if err != nil {
			err = fmt.Errorf("list active release versions: %w", err)
		}
	}()

	slog.Default().DebugContext(ctx, "running list active versions request")

	file, err := l.openReleaseVersions(ctx)
	if err != nil {
		return nil, err
	}

	var bytes []byte
	bytes, err = io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	res.VersionsConfig, err = l.unmarshalConfig(ctx, bytes)

	return res, err
}

// ToJSON marshals the response to JSON.
func (l *ListActiveVersionsRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(l)
	if err != nil {
		return nil, fmt.Errorf("marshaling list changed files to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table.
func (l *ListActiveVersionsRes) ToTable() string {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.AppendHeader(table.Row{"version", "ce active", "lts"})
	for _, version := range slices.Sorted(maps.Keys(l.VersionsConfig.ActiveVersion.Versions)) {
		values := l.VersionsConfig.ActiveVersion.Versions[version]
		t.AppendRow(table.Row{version, values.CEActive, values.LTS})
	}
	return t.Render()
}

// unmarshalConfig unmarshals the bytes of version.hcl into a *VersionsConfig
func (l *ListActiveVersionsReq) unmarshalConfig(ctx context.Context, bytes []byte) (*VersionsConfig, error) {
	var err error
	slog.Default().DebugContext(ctx, "unmarshaling versions.hcl")

	defer func() {
		if err != nil {
			err = fmt.Errorf("unmarhsal versions.hcl bytes: %w", err)
		}
	}()

	parser := hclparse.NewParser()
	var file *hcl.File
	var diags hcl.Diagnostics
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
		err = errors.New("no active_versions stanza found in decoded value")
	}

	for _, versions := range av.AsValueSlice() {
		for version, versionVal := range versions.AsValueMap() {
			v := &Version{}

			for attr, val := range versionVal.AsValueMap() {
				switch attr {
				case "ce_active":
					v.CEActive = val.True()
				case "lts":
					v.LTS = val.True()
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

// openReleaseVersions searches the current path and optionally recursively for
// .release/versions.hcl and returns a file handle to it.
func (l *ListActiveVersionsReq) openReleaseVersions(ctx context.Context) (*os.File, error) {
	var err error
	var file *os.File
	var path string

	defer func() {
		if err != nil {
			err = fmt.Errorf("open release versions: %w", err)
		}
	}()

	if l == nil {
		err = errors.New("uninitialized")
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	slog.Default().DebugContext(ctx, "open .release/versions.hcl")
	if l.ReleaseVersionConfigPath != "" {
		slog.Default().DebugContext(
			slogctx.Append(ctx, slog.String("path", l.ReleaseVersionConfigPath)),
			"attempting to open versions.hcl",
		)

		path, err = filepath.Abs(l.ReleaseVersionConfigPath)
		if err != nil {
			return nil, err
		}

		file, err = os.Open(path)
		return file, err
	}

	path, err = os.Getwd()
	if err != nil {
		return nil, err
	}

	for depth := uint(0); path != string(os.PathSeparator) && depth <= l.Recurse; depth++ {
		path = filepath.Join(path, ".release", "versions.hcl")
		slog.Default().DebugContext(
			slogctx.Append(ctx, slog.String("path", path), slog.Uint64("recurse", uint64(depth))),
			"attempting to open versions.hcl",
		)
		file, err = os.Open(path)
		if err == nil {
			return file, nil
		}
		path, err = filepath.Abs(filepath.Dir(filepath.Dir(filepath.Dir(path))))
		if err != nil {
			return nil, err
		}
	}
	err = errors.New("unable to locate .release/versions.hcl")

	return nil, err
}
