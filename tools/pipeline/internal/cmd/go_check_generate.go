// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/spf13/cobra"
)

var goCheckGenerateReq = &golang.CheckGenerateReq{}

func newGoCheckGenerateCmd() *cobra.Command {
	goCheckGenerateCmd := &cobra.Command{
		Use:   "generate [--go-work PATH] [--tags TAGS] [--module DIR]...",
		Short: "Check that committed generated Go code is up to date",
		Long: `Check that committed generated Go code is up to date.

Runs 'go generate ./...' with GOWORK=off in each module's directory and reports
every file it modified, added or deleted, including new git-ignored files, which
builds that skip generation would never see. Tags are passed in GOFLAGS so
generators that load packages see them too. GOOS, GOARCH, CC and CC_FOR_TARGET
are cleared. Changes are left in the working tree so they can be committed.

Generators that directives call must be installed first.

tools/tools.go is skipped because of its tools build tag. Never pass 'tools' in
--tags: its directive installs every tool.

Known limitations:
  - git reports an ignored directory (e.g. bin/) as a single entry, so files
    written inside a directory that was already ignored aren't detected.
  - Ignored files that existed before the run aren't hashed, so a generator
    that rewrites one isn't detected.
  - On a dirty tree, the diff of a tracked file includes changes made before the
    run. Detection still compares file contents, so it's unaffected.`,
		RunE: runGoCheckGenerateCmd,
		Args: cobra.NoArgs,
	}

	goCheckGenerateCmd.PersistentFlags().StringVar(&goCheckGenerateReq.GoWork, "go-work", "", "Path to go.work. Defaults to the first go.work in the current directory or its parents")
	goCheckGenerateCmd.PersistentFlags().StringVar(&goCheckGenerateReq.Tags, "tags", "", "Comma- or space-separated build tags")
	goCheckGenerateCmd.PersistentFlags().StringArrayVar(&goCheckGenerateReq.Modules, "module", nil, "Repeatable module directory relative to the workspace root. Defaults to every module")

	return goCheckGenerateCmd
}

func runGoCheckGenerateCmd(cmd *cobra.Command, _ []string) error {
	cmd.SilenceUsage = true

	res, err := goCheckGenerateReq.Run(cmd.Context())
	if err != nil {
		return fmt.Errorf("checking go generate: %w", err)
	}

	if err := printGoCmdResult(cmd, res); err != nil {
		return err
	}

	if !res.Success() {
		return fmt.Errorf(
			"check failed: %s; to regenerate, run %s and commit the result",
			res.String(), goCheckGenerateRerun(goCheckGenerateReq),
		)
	}

	return nil
}

// goCheckGenerateRerun returns the command that reruns the check with the same
// flags. The check leaves regenerated files in place, so it doubles as the fix.
func goCheckGenerateRerun(req *golang.CheckGenerateReq) string {
	args := []string{"pipeline", "go", "check", "generate"}
	if req.GoWork != "" {
		args = append(args, "--go-work", shellQuote(req.GoWork))
	}
	if tags := golang.NormalizeTags(req.Tags); len(tags) > 0 {
		args = append(args, "--tags", shellQuote(strings.Join(tags, ",")))
	}
	for _, m := range req.Modules {
		args = append(args, "--module", shellQuote(m))
	}

	return strings.Join(args, " ")
}

// shellQuote single-quotes s when a POSIX shell would otherwise split or
// expand it.
func shellQuote(s string) string {
	if s != "" && !strings.ContainsAny(s, " \t\n'\"\\$`*?[]{}()<>|&;#~!") {
		return s
	}

	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
