// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

// Package generate implements our configuration file generation. We use it to
// generate configuration files dynamically in CI.
package generate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"text/template"
	"time"

	"github.com/Masterminds/semver"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
)

// GenerateTemplateReq is our request to generate a template with pipeline
// context.
type GenerateTemplateReq struct {
	TemplatePath  string                 // Path to template or "-" for stdin
	OutputPath    string                 // Path to output or "" for stdout
	Version       string                 // Current product version to expose in templates (optional)
	Edition       string                 // Current product edition to expose in templates (optional)
	VersionLister releases.VersionLister // How to get version information
}

// GenerateTemplateRes is our GenerateTemplateReq response. It contains the
// bytes of the rendered template. NOTE: It does not write to OutputPath, that
// is up to the caller.
type GenerateTemplateRes struct {
	RenderedTemplate []byte
	TemplatePath     string // Path to template or "-" for stdin
	OutputPath       string // Path to output or "" for stdout
}

// TemplateContext is the root context object available in templates.
type TemplateContext struct {
	VersionLister releases.VersionLister
	GeneratedAt   time.Time // When the template was rendered
	Version       string    // Current version (optional, empty if not provided)
	Edition       string    // Current product edition (optional, empty if not provided)
}

// Run runs the template generation request.
func (r *GenerateTemplateReq) Run(ctx context.Context) (*GenerateTemplateRes, error) {
	res := &GenerateTemplateRes{
		TemplatePath: r.TemplatePath,
		OutputPath:   r.OutputPath,
	}

	tc := &TemplateContext{
		GeneratedAt:   time.Now().UTC(),
		VersionLister: r.VersionLister,
		Version:       r.Version,
		Edition:       r.Edition,
	}

	var err error
	res.RenderedTemplate, err = r.renderTemplate(ctx, tc)
	if err != nil {
		err = fmt.Errorf("rendering template: %w", err)
	}

	return res, err
}

// Render executes the template rendering process.
func (r *GenerateTemplateReq) renderTemplate(ctx context.Context, tc *TemplateContext) ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized")
	}

	var body []byte
	var err error
	if r.TemplatePath == "-" {
		body, err = io.ReadAll(os.Stdin)
	} else {
		body, err = os.ReadFile(r.TemplatePath)
	}
	if err != nil {
		return nil, fmt.Errorf("reading contents of template: %w", err)
	}

	template, err := template.New("template").Funcs(templateFuncsFor(ctx, tc)).Parse(string(body))
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	buf := &bytes.Buffer{}
	if err := template.Execute(buf, tc); err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}

	return buf.Bytes(), nil
}

// templateFuncsFor returns the function map for template execution
func templateFuncsFor(ctx context.Context, tc *TemplateContext) template.FuncMap {
	return template.FuncMap{
		// Version listing functions
		"VersionsNMinus":            versionsNMinusFunc(ctx, tc),
		"VersionsBounded":           versionsBoundedFunc(ctx, tc),
		"VersionsNMinusTransition":  versionsNMinusTransitionFunc(ctx, tc),
		"VersionsBoundedTransition": versionsBoundedTransitionFunc(ctx, tc),

		// Version utilities
		"ParseVersion":    parseVersionFunc,
		"CompareVersions": compareVersionsFunc,
		"FilterVersions":  filterVersionsFunc,
	}
}

// versionsNMinusFunc returns a function that lists released versions using
// n-minus.
func versionsNMinusFunc(ctx context.Context, tc *TemplateContext) func(
	product, edition, upperBound string,
	nminus uint,
	cadence releases.VersionCadence,
	skip ...string) ([]string, error) {
	return func(
		product, edition, upperBound string,
		nminus uint,
		cadence releases.VersionCadence,
		skip ...string,
	) ([]string, error) {
		req := &releases.ListVersionsReq{
			VersionLister: tc.VersionLister,
			Cadence:       cadence,
			ProductName:   product,
			LicenseClass:  edition,
			UpperBound:    upperBound,
			NMinus:        nminus,
			Skip:          skip,
		}

		res, err := req.Run(ctx)
		if err != nil {
			return nil, err
		}

		return res.Versions, nil
	}
}

// versionsBoundedFunc returns a function that lists released versions using
// explicit bounds.
func versionsBoundedFunc(ctx context.Context, tc *TemplateContext) func(
	product, edition, upperBound, lowerBound string,
	cadence releases.VersionCadence,
	skip ...string) ([]string, error) {
	return func(
		product, edition, upperBound, lowerBound string,
		cadence releases.VersionCadence,
		skip ...string,
	) ([]string, error) {
		req := &releases.ListVersionsReq{
			VersionLister: tc.VersionLister,
			Cadence:       cadence,
			ProductName:   product,
			LicenseClass:  edition,
			UpperBound:    upperBound,
			LowerBound:    lowerBound,
			Skip:          skip,
		}

		res, err := req.Run(ctx)
		if err != nil {
			return nil, err
		}

		return res.Versions, nil
	}
}

// versionsNMinusTransitionFuncction that lists released versions using n-minus
// with cadence support.
func versionsNMinusTransitionFunc(ctx context.Context, tc *TemplateContext) func(
	product, edition, upperBound string,
	nminus uint,
	cadence, transitionVersion, priorCadence string,
	skip ...string,
) ([]string, error) {
	return func(
		product, edition, upperBound string,
		nminus uint,
		cadence, transitionVersion, priorCadence string,
		skip ...string,
	) ([]string, error) {
		req := &releases.ListVersionsReq{
			VersionLister:     tc.VersionLister,
			ProductName:       product,
			LicenseClass:      edition,
			UpperBound:        upperBound,
			NMinus:            nminus,
			Cadence:           releases.VersionCadence(cadence),
			TransitionVersion: transitionVersion,
			PriorCadence:      releases.VersionCadence(priorCadence),
			Skip:              skip,
		}

		res, err := req.Run(ctx)
		if err != nil {
			return nil, err
		}

		return res.Versions, nil
	}
}

// versionsBoundedTransitionFunc returns a function that lists released versions
// using explicit bounds with cadence support.
func versionsBoundedTransitionFunc(ctx context.Context, tc *TemplateContext) func(
	product, edition, upperBound, lowerBound, cadence string,
	transitionVersion, priorCadence string,
	skip ...string,
) ([]string, error) {
	return func(
		product, edition, upperBound, lowerBound, cadence string,
		transitionVersion, priorCadence string,
		skip ...string,
	) ([]string, error) {
		req := &releases.ListVersionsReq{
			VersionLister:     tc.VersionLister,
			ProductName:       product,
			LicenseClass:      edition,
			UpperBound:        upperBound,
			LowerBound:        lowerBound,
			Cadence:           releases.VersionCadence(cadence),
			TransitionVersion: transitionVersion,
			PriorCadence:      releases.VersionCadence(priorCadence),
			Skip:              skip,
		}

		res, err := req.Run(ctx)
		if err != nil {
			return nil, err
		}

		return res.Versions, nil
	}
}

// parseVersionFunc parses a version string into a semver.Version
func parseVersionFunc(version string) (*semver.Version, error) {
	v, err := semver.NewVersion(version)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// compareVersionsFunc compares two version strings
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersionsFunc(v1, v2 string) (int, error) {
	ver1, err := semver.NewVersion(v1)
	if err != nil {
		return 0, err
	}

	ver2, err := semver.NewVersion(v2)
	if err != nil {
		return 0, err
	}

	return ver1.Compare(ver2), nil
}

// filterVersionsFunc filters versions based on a semver constraint
func filterVersionsFunc(versions []string, constraint string) ([]string, error) {
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return nil, err
	}

	var filtered []string
	for _, v := range versions {
		ver, err := semver.NewVersion(v)
		if err != nil {
			// Skip invalid versions
			continue
		}
		if c.Check(ver) {
			filtered = append(filtered, v)
		}
	}

	return filtered, nil
}
