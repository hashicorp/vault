// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/generate"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/spf13/cobra"
)

var generateTemplateReq = &generate.GenerateTemplateReq{
	VersionLister: releases.NewClient(),
}

func newGenerateTemplateCmd() *cobra.Command {
	templateCmd := &cobra.Command{
		Use:   "template [template-path] [output-path]",
		Short: "Generate output from Go templates with pipeline data",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  runGenerateTemplateCmd,
		Long: `Generate output by rendering Go templates with access to pipeline data. It's currently limited to release data but the idea is to eventually extend it to any request types that would be useful to embed into templates via template functions.

The template has access to a context object with the following fields:
  - .Version: Current version (set via --version flag)
  - .Edition: Current edition (set via --edition flag)
  - .GeneratedAt: Timestamp when template was rendered

Available template functions:

  - VersionsNMinus(product, edition, upperBound, nminus, cadence, versionsToSkip...)
    Use this to list all version between the upperBound and n-x. The candence must be set to 'minor' or 'major' depending on whether n-minus should be treated as minor or major bumps. Returns a list of strings.

    - VersionsNMinusTransition(product, edition, upperBound, nminus, current cadence, transitionVersion, priorCadence, versionsToSkip skip...)
    Like VersionsNMinus, except that it supports two different release candences by supplying both and a transition version. Useful for navigating a change in release candence. Returns a list of strings.

  - VersionsBounded(product, edition, upperBound, lowerBound, cadence, versionsToSkip...)
    Use this to list all version between two bounds. The candence must be set to 'minor' or 'major' depending on whether n-minus should be treated as minor or major bumps. Returns a list of strings.

  - VersionsBoundedTransition(product, edition, upperBound, lowerBound, current cadence, transitionVersion, priorCadence, versionsToSkip ...)
    Like VersionsBounded, except that it supports two different release candences by supplying both and a transition version. Useful for navigating a change in release candence. Returns a list of strings.

  - ParseVersion(version)
    Parse version string to semver. Returns a *semver.Version that you can access methods and fields on. See https://pkg.go.dev/github.com/Masterminds/semver

  - CompareVersions(v1, v2)
    Parse and compare two different versions. Returns an int. See https://pkg.go.dev/github.com/Masterminds/semver#Version.Compare

  - FilterVersions(listOfVersions, semverConstraint)
    Filter a list of versions by a semver constraint. Returns a list of strings.

Examples:
  # Render template with current version
  pipeline generate template my-template.tmpl output.txt --version 1.18.0 --edition enterprise

  # Read from stdin, write to stdout
  cat template.tmpl | pipeline generate template - --version 1.18.0

  # Use current version in template
  {{ .Version }} - Current version
  {{ VersionsNMinus "vault" "enterprise" .Version 3 "minor" }} - Get all n-3 versions from the .Version (passed in via --version) using a minor version cadence.`,
	}

	templateCmd.Flags().StringVarP(&generateTemplateReq.Version, "version", "v", "", "Current version to expose in template context")
	templateCmd.Flags().StringVarP(&generateTemplateReq.Edition, "edition", "e", "", "Current edition to expose in template context")

	return templateCmd
}

func runGenerateTemplateCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	generateTemplateReq.TemplatePath = args[0]
	if len(args) > 1 {
		generateTemplateReq.OutputPath = args[1]
	}

	res, err := generateTemplateReq.Run(cmd.Context())
	if err != nil {
		return fmt.Errorf("generating template: %w", err)
	}

	if res.OutputPath != "" {
		err = os.WriteFile(res.OutputPath, res.RenderedTemplate, 0o755)
	} else {
		_, err = io.Writer.Write(os.Stdout, res.RenderedTemplate)
	}

	if err != nil {
		return fmt.Errorf("generating template: %w", err)
	}

	return nil
}
