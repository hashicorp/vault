// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"errors"
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/spf13/cobra"
)

var checkGithubGoModReq = github.CheckGoModDiffReq{
	DiffOpts: golang.DefaultDiffOpts(),
	Paths:    []string{},
}

func newGithubCheckGoModDiffCmd() *cobra.Command {
	checkGoModCmd := &cobra.Command{
		Use:   "go-mod-diff [FLAGS]",
		Short: "Compare Go modules for equality",
		Long:  "Compare one-or-more Go modules from different branches hosted in Github",
		RunE:  runCheckGithubGoModCmd,
	}

	// Repository, branch, and path flags
	checkGoModCmd.PersistentFlags().StringVar(&checkGithubGoModReq.AOwner, "a-owner", "hashicorp", "The Github organization hosting the a branch")
	checkGoModCmd.PersistentFlags().StringVar(&checkGithubGoModReq.ARepo, "a-repo", "vault-enterprise", "The Github repository hosting the a branch")
	checkGoModCmd.PersistentFlags().StringVar(&checkGithubGoModReq.ABranch, "a-branch", "", "The name of the a branch we want diff")
	checkGoModCmd.PersistentFlags().StringVar(&checkGithubGoModReq.BOwner, "b-owner", "hashicorp", "The Github organization hosting the b branch")
	checkGoModCmd.PersistentFlags().StringVar(&checkGithubGoModReq.BRepo, "b-repo", "vault-enterprise", "The Github repository hosting the b branch")
	checkGoModCmd.PersistentFlags().StringVar(&checkGithubGoModReq.BBranch, "b-branch", "", "The name of the b branch we want diff")
	checkGoModCmd.PersistentFlags().StringSliceVarP(&checkGithubGoModReq.Paths, "path", "p", []string{}, "The go.mod paths relative to the repository to use. e.g. -p go.mod -p api/go.mod")

	err := checkGoModCmd.MarkPersistentFlagRequired("a-branch")
	if err != nil {
		panic(err)
	}
	err = checkGoModCmd.MarkPersistentFlagRequired("b-branch")
	if err != nil {
		panic(err)
	}

	// Diff option flags
	checkGoModCmd.PersistentFlags().BoolVar(&checkGithubGoModReq.DiffOpts.ParseLax, "lax", false, "Parse the modules in lax mode to ignore newer and unknown directives")

	// Enable or disable diffing at a directive level
	checkGoModCmd.PersistentFlags().BoolVar(&checkGithubGoModReq.DiffOpts.Module, "module", true, "Compare the module directives in both files")
	checkGoModCmd.PersistentFlags().BoolVar(&checkGithubGoModReq.DiffOpts.Go, "go", true, "Compare the go directives in both files")
	checkGoModCmd.PersistentFlags().BoolVar(&checkGithubGoModReq.DiffOpts.Toolchain, "toolchain", true, "Compare the toolchain directives in both files")
	checkGoModCmd.PersistentFlags().BoolVar(&checkGithubGoModReq.DiffOpts.Godebug, "godebug", true, "Compare the godebug directives in both files")
	checkGoModCmd.PersistentFlags().BoolVar(&checkGithubGoModReq.DiffOpts.Require, "require", true, "Compare the require directives in both files")
	checkGoModCmd.PersistentFlags().BoolVar(&checkGithubGoModReq.DiffOpts.Replace, "replace", true, "Compare the replace directives in both files")
	checkGoModCmd.PersistentFlags().BoolVar(&checkGithubGoModReq.DiffOpts.Retract, "retract", true, "Compare the retract directives in both files")
	checkGoModCmd.PersistentFlags().BoolVar(&checkGithubGoModReq.DiffOpts.Tool, "tool", true, "Compare the tool directives in both files")
	checkGoModCmd.PersistentFlags().BoolVar(&checkGithubGoModReq.DiffOpts.Ignore, "ignore", true, "Compare the ignore directives in both files")

	// Enable or disable strict diffing at a directive level. When strict diffing is disabled only like directives will be compared.
	checkGoModCmd.PersistentFlags().BoolVar(&checkGithubGoModReq.DiffOpts.StrictDiffRequire, "strict-require", true, "Strictly compare the requires directives in both files. When true all requires are compared, otherwise only shared requires are compared")
	checkGoModCmd.PersistentFlags().BoolVar(&checkGithubGoModReq.DiffOpts.StrictDiffExclude, "strict-exclude", true, "Strictly compare the excludes directives in both files. When true all excludes are compared, otherwise only shared excludes are compared")
	checkGoModCmd.PersistentFlags().BoolVar(&checkGithubGoModReq.DiffOpts.StrictDiffReplace, "strict-replace", true, "Strictly compare the replace directives in both files. When true all replaces are compared, otherwise only shared replaces are compared")
	checkGoModCmd.PersistentFlags().BoolVar(&checkGithubGoModReq.DiffOpts.StrictDiffRetract, "strict-retract", true, "Strictly compare the retract directives in both files. When true all retracts are compared, otherwise only shared retract directives are compared")

	return checkGoModCmd
}

func runCheckGithubGoModCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	res, err := checkGithubGoModReq.Run(cmd.Context(), githubCmdState.GithubV3, rootCfg.git)
	if err != nil {
		return err
	}

	switch rootCfg.format {
	case "json":
		b, err1 := res.ToJSON()
		if err1 != nil {
			err = errors.Join(err, err1)
		} else {
			fmt.Println(string(b))
		}
	case "markdown":
		tbl, err1 := res.ToTable(err)
		if err1 != nil {
			err = errors.Join(err, err1)
		} else {
			tbl.SetTitle("Go Mod Diff")
			fmt.Println(tbl.RenderMarkdown())
		}
	default:
		tbl, err1 := res.ToTable(err)
		if err1 != nil {
			err = errors.Join(err, err1)
		} else {
			if text := tbl.Render(); text != "" {
				fmt.Println(text)
			}
		}
	}

	for _, check := range res.Diffs {
		if l := len(check.ModDiff); l > 0 {
			err = errors.Join(fmt.Errorf("%d differences were found", l), err)
			break
		}
	}

	return err
}
