// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

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
}

var rootCfg = &rootCmdCfg{}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "pipeline",
		Short: "Execute pipeline tasks",
		Long:  "Pipeline automation tasks",
	}

	rootCmd.PersistentFlags().StringVar(&rootCfg.logLevel, "log", "warn", "Set the log level. One of 'debug', 'info', 'warn', 'error'")

	rootCmd.AddCommand(newGenerateCmd())
	rootCmd.AddCommand(newGithubCmd())
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
		h := slogctx.NewHandler(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: ll}), nil)
		slog.SetDefault(slog.New(h))

		return nil
	}

	return rootCmd
}

// Execute executes the root pipeline command.
func Execute() {
	rootCmd := newRootCmd()
	rootCmd.SilenceErrors = true // We handle this below

	if err := rootCmd.Execute(); err != nil {
		slog.Default().Error(err.Error())
		os.Exit(1)
	}
}
