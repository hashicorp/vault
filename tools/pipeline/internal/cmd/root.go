// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

// Package cmd defines the pipeline CLI commands.
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	slogctx "github.com/veqryn/slog-context"
)

type rootCmdCfg struct {
	logLevel string
	format   string
}

var rootCfg = &rootCmdCfg{}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "pipeline",
		Short: "Execute pipeline tasks",
		Long:  "Pipeline automation tasks",
	}

	rootCmd.PersistentFlags().StringVar(&rootCfg.logLevel, "log", "warn", "Set the log level. One of 'debug', 'info', 'warn', 'error'")
	rootCmd.PersistentFlags().StringVarP(&rootCfg.format, "format", "f", "table", "The output format. Can be 'json', 'table', and sometimes 'markdown'")

	rootCmd.AddCommand(newGenerateCmd())
	rootCmd.AddCommand(newGithubCmd())
	rootCmd.AddCommand(newGoCmd())
	rootCmd.AddCommand(newHCPCmd())
	rootCmd.AddCommand(newReleasesCmd())

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
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
