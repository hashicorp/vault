// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/spf13/cobra"
)

var goModDiffReq = &golang.DiffModReq{
	A:    &golang.ModSource{},
	B:    &golang.ModSource{},
	Opts: golang.DefaultDiffOpts(),
}

func newGoDiffModCmd() *cobra.Command {
	goModDiffCmd := &cobra.Command{
		Use:   "mod </path/a/go.mod> </path/b/go.mod> [ARGS]",
		Short: "Diff two local go.mod files",
		Long:  "Diff two local go.mod files",
		RunE:  runGoModDiffCmd,
		Args: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 2:
				err := golang.SetUpGoModSourceFromPath(args[0], goModDiffReq.A)
				if err != nil {
					return err
				}
				return golang.SetUpGoModSourceFromPath(args[1], goModDiffReq.B)
			case 0, 1:
				return errors.New("invalid arguments: you must provide two local file paths")
			default:
				return fmt.Errorf("invalid arguments: expected two paths as arguments, received %d arguments", len(args))
			}
		},
	}

	goModDiffCmd.PersistentFlags().BoolVar(&goModDiffReq.Opts.ParseLax, "lax", false, "Parse the modules in lax mode to ignore newer and unknown directives")

	goModDiffCmd.PersistentFlags().BoolVar(&goModDiffReq.Opts.Module, "module", true, "Compare the module directives in both files")
	goModDiffCmd.PersistentFlags().BoolVar(&goModDiffReq.Opts.Go, "go", true, "Compare the go directives in both files")
	goModDiffCmd.PersistentFlags().BoolVar(&goModDiffReq.Opts.Toolchain, "toolchain", true, "Compare the toolchain directives in both files")
	goModDiffCmd.PersistentFlags().BoolVar(&goModDiffReq.Opts.Godebug, "godebug", true, "Compare the godebug directives in both files")
	goModDiffCmd.PersistentFlags().BoolVar(&goModDiffReq.Opts.Require, "require", true, "Compare the require directives in both files")
	goModDiffCmd.PersistentFlags().BoolVar(&goModDiffReq.Opts.Replace, "replace", true, "Compare the replace directives in both files")
	goModDiffCmd.PersistentFlags().BoolVar(&goModDiffReq.Opts.Retract, "retract", true, "Compare the retract directives in both files")
	goModDiffCmd.PersistentFlags().BoolVar(&goModDiffReq.Opts.Tool, "tool", true, "Compare the tool directives in both files")
	goModDiffCmd.PersistentFlags().BoolVar(&goModDiffReq.Opts.Ignore, "ignore", true, "Compare the ignore directives in both files")

	// Enable or disable strict diffing at a directive level. When strict diffing is disabled only like directives will be compared.
	goModDiffCmd.PersistentFlags().BoolVar(&goModDiffReq.Opts.StrictDiffRequire, "strict-require", true, "Strictly compare the requires directives in both files. When true all requires are compared, otherwise only shared requires are compared")
	goModDiffCmd.PersistentFlags().BoolVar(&goModDiffReq.Opts.StrictDiffExclude, "strict-exclude", true, "Strictly compare the excludes directives in both files. When true all excludes are compared, otherwise only shared excludes are compared")
	goModDiffCmd.PersistentFlags().BoolVar(&goModDiffReq.Opts.StrictDiffReplace, "strict-replace", true, "Strictly compare the replace directives in both files. When true all replaces are compared, otherwise only shared replaces are compared")
	goModDiffCmd.PersistentFlags().BoolVar(&goModDiffReq.Opts.StrictDiffRetract, "strict-retract", true, "Strictly compare the retract directives in both files. When true all retracts are compared, otherwise only shared retract directives are compared")

	// Explicitly exclude individual directives. Can use literal directive values or path style globs.
	goModDiffCmd.PersistentFlags().StringArrayVar(&goModDiffReq.Opts.ExcludeRequire, "exclude-require", nil, "repeatable glob/literal pattern; skip matching module paths in require directives")
	goModDiffCmd.PersistentFlags().StringArrayVar(&goModDiffReq.Opts.ExcludeReplace, "exclude-replace", nil, "repeatable glob/literal pattern; skip matching module paths in replace directives (old or new path)")

	return goModDiffCmd
}

func runGoModDiffCmd(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true // Don't spam the usage on failure

	res, err := goModDiffReq.Run(context.TODO())
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

	if l := len(res.ModDiff); l > 0 {
		err = errors.Join(fmt.Errorf("%d differences were found", l), err)
	}

	return err
}
