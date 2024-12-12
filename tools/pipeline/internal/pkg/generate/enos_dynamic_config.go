// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package generate

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/Masterminds/semver"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/metadata"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	slogctx "github.com/veqryn/slog-context"
)

// EnosDynamicConfigReq is a request to generate dynamic enos configuration
type EnosDynamicConfigReq struct {
	VaultEdition string
	VaultVersion string
	EnosDir      string
	FileName     string
	Skip         []string
	NMinus       uint
	releases.VersionLister
}

// EnosDynamicConfigRes is a response from a request to generate dynamic enos configuration
type EnosDynamicConfigRes struct {
	Globals *Globals `json:"globals,omitempty" hcl:"globals,block" cty:"globals"`
}

// Globals are our dynamic globals
type Globals struct {
	SampleAttributes *SampleAttrs `json:"sample_attributes,omitempty" hcl:"sample_attributes" cty:"sample_attributes"`
}

// SampleAttrs are the dynamic sample attributes that we'll write as globals
type SampleAttrs struct {
	AWSRegion             []string `json:"aws_region,omitempty" hcl:"aws_region" cty:"aws_region"`
	DistroVersionAmzn     []string `json:"distro_version_amzn,omitempty" hcl:"distro_version_amzn" cty:"distro_version_amzn"`
	DistroVersionLeap     []string `json:"distro_version_leap,omitempty" hcl:"distro_version_leap" cty:"distro_version_leap"`
	DistroVersionRhel     []string `json:"distro_version_rhel,omitempty" hcl:"distro_version_rhel" cty:"distro_version_rhel"`
	DistroVersionSles     []string `json:"distro_version_sles,omitempty" hcl:"distro_version_sles" cty:"distro_version_sles"`
	DistroVersionUbuntu   []string `json:"distro_version_ubuntu,omitempty" hcl:"distro_version_ubuntu" cty:"distro_version_ubuntu"`
	UpgradeInitialVersion []string `json:"upgrade_initial_version,omitempty" hcl:"upgrade_initial_version" cty:"upgrade_initial_version"`
}

// Validate validates the request parameters
func (e *EnosDynamicConfigReq) Validate(ctx context.Context) error {
	if e == nil {
		return errors.New("enos dynamic config req: validate: uninitialized")
	}

	slog.Default().DebugContext(ctx, "validating enos dynamic config request")

	if e.FileName == "" {
		return errors.New("no destination file name set")
	}

	if e.VersionLister == nil {
		return errors.New("no version lister set")
	}

	if !slices.Contains(metadata.Editions, e.VaultEdition) {
		return fmt.Errorf("unknown edition: %s", e.VaultEdition)
	}

	_, err := semver.NewVersion(e.VaultVersion)
	if err != nil {
		return fmt.Errorf("invalid version: %s: %w", e.VaultVersion, err)
	}

	s, err := os.Stat(e.EnosDir)
	if err != nil {
		return fmt.Errorf("invalid enos dir: %s: %w", e.EnosDir, err)
	}

	if !s.IsDir() {
		return fmt.Errorf("invalid enos dir: %s is not a directory", e.EnosDir)
	}

	return nil
}

// Run runs the dynamic configuration request
func (e *EnosDynamicConfigReq) Run(ctx context.Context) (*EnosDynamicConfigRes, error) {
	if e == nil {
		return nil, fmt.Errorf("enos-dynamic-config-req uninitialized")
	}

	ctx = slogctx.Append(ctx,
		slog.String("vault-edition", e.VaultEdition),
		slog.String("vault-version", e.VaultVersion),
		slog.String("file-name", e.FileName),
		slog.String("dir", e.EnosDir),
		slog.Uint64("n-minnux", uint64(e.NMinus)),
		"skip", e.Skip,
	)
	slog.Default().DebugContext(ctx, "running enos dynamic config request")

	err := e.Validate(ctx)
	if err != nil {
		return nil, err
	}

	res := &EnosDynamicConfigRes{}
	res.Globals, err = e.getGlobals(ctx)
	if err != nil {
		return nil, err
	}

	return res, e.writeFile(ctx, res)
}

func (e *EnosDynamicConfigReq) getGlobals(ctx context.Context) (*Globals, error) {
	var err error
	res := &Globals{}
	res.SampleAttributes, err = e.getSampleAttrs(ctx)

	return res, err
}

func (e *EnosDynamicConfigReq) getSampleAttrs(ctx context.Context) (*SampleAttrs, error) {
	// Create our HCL body
	attrs := &SampleAttrs{
		// Use the cheapest regions
		AWSRegion: []string{"us-east-1", "us-west-2"},
		// Current distro defaults
		DistroVersionAmzn:   []string{"2023"},
		DistroVersionLeap:   []string{"15.6"},
		DistroVersionRhel:   []string{"8.10", "9.4"},
		DistroVersionSles:   []string{"15.6"},
		DistroVersionUbuntu: []string{"20.04", "24.04"},
	}

	// Create our initial upgrade version list. We'll find all released versions between N-3 -> Current
	// version, minus any explicitly skipped versions that have been set. Since CE and Ent do not share
	// the same version lineage now we'll also have to figure that in as well.
	versionReq := &releases.ListVersionsReq{
		VersionLister: e.VersionLister,
		LicenseClass:  e.VaultEdition,
		UpperBound:    e.VaultVersion,
		NMinus:        e.NMinus,
		Skip:          e.Skip,
	}

	versionRes, err := versionReq.Run(ctx)
	if err != nil {
		return nil, err
	}
	attrs.UpgradeInitialVersion = versionRes.Versions

	return attrs, nil
}

// writeFile creates the dynamic config file and writes the dynamic data into it
func (e *EnosDynamicConfigReq) writeFile(ctx context.Context, res *EnosDynamicConfigRes) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	slog.Default().DebugContext(ctx, "writing enos dynamic config request")

	// Make sure our path is valid
	path, err := filepath.Abs(filepath.Join(e.EnosDir, e.FileName))
	if err != nil {
		return fmt.Errorf("expanding path dynamic config request path: %w", err)
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("opening dynamic config request destination file: %w", err)
	}
	defer f.Close()

	if _, err = f.WriteString(`# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

# Code generated by pipeline generate enos-dynamic-config DO NOT EDIT.

# This file is overwritten in CI as it contains branch specific and sometimes ever-changing values.
# It's checked in here so that enos samples and scenarios can be performed, just be aware that this
# might change out from under you.
`); err != nil {
		return err
	}

	hf := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(res, hf.Body())
	bytes := hclwrite.Format(hf.Bytes())

	slog.Default().InfoContext(ctx, "writing enos dynamic config request",
		"path", path,
		"hcl", string(bytes),
	)
	_, err = f.Write(bytes)

	return err
}
