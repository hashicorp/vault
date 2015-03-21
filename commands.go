package main

import (
	"os"

	"github.com/hashicorp/vault/builtin/logical/aws"
	"github.com/hashicorp/vault/builtin/logical/consul"
	"github.com/hashicorp/vault/command"
	"github.com/hashicorp/vault/logical"
	"github.com/mitchellh/cli"
)

// Commands is the mapping of all the available Consul commands.
var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}
	meta := command.Meta{Ui: ui}

	Commands = map[string]cli.CommandFactory{
		"auth": func() (cli.Command, error) {
			return &command.AuthCommand{
				Meta: meta,
			}, nil
		},

		"read": func() (cli.Command, error) {
			return &command.ReadCommand{
				Meta: meta,
			}, nil
		},

		"write": func() (cli.Command, error) {
			return &command.WriteCommand{
				Meta: meta,
			}, nil
		},

		"seal": func() (cli.Command, error) {
			return &command.SealCommand{
				Meta: meta,
			}, nil
		},

		"seal-status": func() (cli.Command, error) {
			return &command.SealStatusCommand{
				Meta: meta,
			}, nil
		},

		"unseal": func() (cli.Command, error) {
			return &command.UnsealCommand{
				Meta: meta,
			}, nil
		},

		"init": func() (cli.Command, error) {
			return &command.InitCommand{
				Meta: meta,
			}, nil
		},

		"server": func() (cli.Command, error) {
			return &command.ServerCommand{
				Meta: meta,
				LogicalBackends: map[string]logical.Factory{
					"aws":    aws.Factory,
					"consul": consul.Factory,
				},
			}, nil
		},

		"mounts": func() (cli.Command, error) {
			return &command.MountsCommand{
				Meta: meta,
			}, nil
		},

		"version": func() (cli.Command, error) {
			ver := Version
			rel := VersionPrerelease
			if GitDescribe != "" {
				ver = GitDescribe
			}
			if GitDescribe == "" && rel == "" {
				rel = "dev"
			}

			return &command.VersionCommand{
				Revision:          GitCommit,
				Version:           ver,
				VersionPrerelease: rel,
				Ui:                ui,
			}, nil
		},
	}
}
