// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

// Package cmd defines the pipeline CLI commands.
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/config"
	git "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git/client"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/spf13/cobra"
	slogctx "github.com/veqryn/slog-context"
)

type rootCmdCfg struct {
	logLevel          string
	format            string
	git               *git.Client
	configDecodeRes   *config.DecodeRes
	versionsDecodeRes *releases.DecodeRes
}

var rootCfg = &rootCmdCfg{
	git: git.NewClient(git.WithLoadTokenFromEnv()),
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "pipeline",
		Short: "Execute pipeline tasks",
		Long:  "Pipeline automation tasks",
	}

	var pipelineCfgPath string
	var versionsConfigPath string

	rootCmd.PersistentFlags().StringVar(&rootCfg.logLevel, "log", "warn", "Set the log level. One of 'debug', 'info', 'warn', 'error'")
	rootCmd.PersistentFlags().StringVarP(&rootCfg.format, "format", "f", "table", "The output format. Can be 'json', 'table', and sometimes 'markdown'")
	rootCmd.PersistentFlags().StringVar(&pipelineCfgPath, "pipeline-config", "", "Specify the path to pipeline.hcl configuration file (default: <git repo root>/.release/pipeline.hcl)")
	rootCmd.PersistentFlags().StringVar(&versionsConfigPath, "versions-config", "", "Specify the path to versions.hcl configuration file (default: <git repo root>/.release/versions.hcl)")

	rootCmd.AddCommand(newConfigCmd())
	rootCmd.AddCommand(newGenerateCmd())
	rootCmd.AddCommand(newGitCmd())
	rootCmd.AddCommand(newGithubCmd())
	rootCmd.AddCommand(newGoCmd())
	rootCmd.AddCommand(newHCPCmd())
	rootCmd.AddCommand(newReleasesCmd())

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		// Setup a default logger before we process anything
		var ll slog.Level
		switch rootCfg.logLevel {
		case "debug":
			ll = slog.LevelDebug
		case "info":
			ll = slog.LevelInfo
		case "warn":
			ll = slog.LevelWarn
		case "error":
			ll = slog.LevelError
		default:
			return fmt.Errorf("unsupported log level: %s", rootCfg.logLevel)
		}
		h := slogctx.NewHandler(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: ll}), nil)
		slog.SetDefault(slog.New(h))

		switch rootCfg.format {
		case "json", "table", "markdown":
		default:
			return fmt.Errorf("unsupported format: %s", rootCfg.format)
		}

		getRepoRoot := sync.OnceValues(func() (string, error) {
			slog.DebugContext(ctx, "determining repository root to load configuration")
			revParse, err := rootCfg.git.RevParse(ctx, &git.RevParseOpts{
				ShowTopLevel: true,
			})
			if err != nil {
				return "", fmt.Errorf("getting the repository root %s: %w", revParse.String(), err)
			}

			return filepath.Join(strings.TrimSpace(string(revParse.Stdout))), nil
		})

		// Get repo root if needed
		var repoRoot string
		var err error
		if pipelineCfgPath == "" || versionsConfigPath == "" {
			repoRoot, err = getRepoRoot()
			if err != nil {
				return err
			}
		}

		// Decode the pipeline config. Store the result (including any errors)
		// for commands to handle as needed.
		if pipelineCfgPath == "" {
			pipelineCfgPath = filepath.Join(repoRoot, ".release", "pipeline.hcl")
		}
		rootCfg.configDecodeRes = config.Decode(ctx, &config.DecodeReq{
			Path: pipelineCfgPath,
		})

		// Decode the versions config. Store the result (including any errors)
		// for commands to handle as needed.
		if versionsConfigPath == "" {
			versionsConfigPath = filepath.Join(repoRoot, ".release", "versions.hcl")
		}
		rootCfg.versionsDecodeRes = releases.Decode(ctx, &releases.DecodeReq{
			Path: versionsConfigPath,
		})

		return nil
	}

	return rootCmd
}

// Execute executes the root pipeline command.
func Execute() {
	cobra.EnableTraverseRunHooks = true // Automatically chain run hooks
	rootCmd := newRootCmd()
	rootCmd.SilenceErrors = true // We handle this below

	if err := rootCmd.Execute(); err != nil {
		slog.Default().Error(err.Error())
		os.Exit(1)
	}
}
