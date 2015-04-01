package main

import (
	"os"

	"github.com/hashicorp/vault/builtin/credential/github"
	"github.com/hashicorp/vault/builtin/logical/aws"
	"github.com/hashicorp/vault/builtin/logical/consul"
	tokenDisk "github.com/hashicorp/vault/builtin/token/disk"
	"github.com/hashicorp/vault/command"
	"github.com/hashicorp/vault/logical"
	"github.com/mitchellh/cli"
)

// Commands is the mapping of all the available Vault commands. CommandsInclude
// are the commands to include for help.
var Commands map[string]cli.CommandFactory
var CommandsInclude []string

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

		"revoke": func() (cli.Command, error) {
			return &command.RevokeCommand{
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
				CredentialBackends: map[string]logical.Factory{
					"github": github.Factory,
				},
				LogicalBackends: map[string]logical.Factory{
					"aws":    aws.Factory,
					"consul": consul.Factory,
				},
			}, nil
		},

		"mount": func() (cli.Command, error) {
			return &command.MountCommand{
				Meta: meta,
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

	// Build the commands to include in the help now
	CommandsInclude = make([]string, 0, len(Commands))
	for k, _ := range Commands {
		CommandsInclude = append(CommandsInclude, k)
	}

	// The commands below are hidden from the help output
	Commands["token-disk"] = func() (cli.Command, error) {
		return &tokenDisk.Command{}, nil
	}
}
