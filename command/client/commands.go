// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/cli"
	hcpvlib "github.com/hashicorp/vault-hcp-lib"
	base "github.com/hashicorp/vault/command/base"
	"github.com/hashicorp/vault/version"
)

// Run is the entry point for the client binary.
func Run(args []string) int {
	return RunCustom(args, nil)
}

// RunCustom allows passing in base run options (e.g. a custom token helper).
func RunCustom(args []string, runOpts *base.RunOptions) int {
	return base.RunCustomWithCommandsFn(args, runOpts, InitClientBinaryCommands, unavailableClientCommands...)
}

// serverInstallURL is where users can obtain the vault binary, which includes
// the server that vault-client leaves out.
const serverInstallURL = "https://developer.hashicorp.com/vault/install"

// UnavailableCommand points users to the full Vault binary for commands we
// don't include in vault-client.
type UnavailableCommand struct {
	ui   cli.Ui
	name string
}

var _ cli.Command = (*UnavailableCommand)(nil)

func (c *UnavailableCommand) Synopsis() string {
	return "Not included in this binary"
}

func (c *UnavailableCommand) Help() string {
	return c.message()
}

func (c *UnavailableCommand) Run(_ []string) int {
	c.ui.Error(c.message())
	return 1
}

func (c *UnavailableCommand) message() string {
	return fmt.Sprintf("%q is not available in vault-client v%s.\n"+
		"This command requires the full Vault Server image.\n"+
		"Obtain the full Vault binary: %s",
		c.name, version.GetVersion().VersionNumber(), serverInstallURL)
}

// Hide these commands from help and completion. Keep them out of
// InitClientCommands so we don't replace real commands in the vault binary.
// The operator key-status, rotate, and unseal commands are supported client
// commands (PRD VLT-2015 AC-1.3), so they must not be listed here.
var unavailableClientCommands = []string{
	"audit",
	"audit disable",
	"audit enable",
	"audit list",
	"debug",
	"events",
	"events subscribe",
	"events subscriptions",
	"events subscriptions create",
	"events subscriptions delete",
	"events subscriptions get",
	"events subscriptions list",
	"license",
	"license get",
	"license inspect",
	"monitor",
	"operator diagnose",
	"operator generate-root",
	"operator import",
	"operator init",
	"operator members",
	"operator migrate",
	"operator raft",
	"operator raft autopilot",
	"operator raft autopilot get-config",
	"operator raft autopilot set-config",
	"operator raft autopilot state",
	"operator raft join",
	"operator raft list-peers",
	"operator raft remove-peer",
	"operator raft snapshot",
	"operator raft snapshot inspect",
	"operator raft snapshot load",
	"operator raft snapshot restore",
	"operator raft snapshot save",
	"operator raft snapshot unload",
	"operator rekey",
	"operator seal",
	"operator step-down",
	"operator usage",
	"operator utilization",
	"pki",
	"pki health-check",
	"pki issue",
	"pki list-intermediates",
	"pki reissue",
	"pki verify-sign",
	"plugin",
	"plugin deregister",
	"plugin info",
	"plugin list",
	"plugin register",
	"plugin reload",
	"plugin reload-status",
	"plugin runtime",
	"plugin runtime deregister",
	"plugin runtime info",
	"plugin runtime list",
	"plugin runtime register",
	"policy",
	"policy delete",
	"policy fmt",
	"policy list",
	"policy read",
	"policy write",
	"server",
	"transform",
	"transform import",
	"transform import-version",
	"transit",
	"transit envelope",
	"transit envelope decrypt",
	"transit envelope encrypt",
	"transit envelope header",
	"transit import",
	"transit import-version",
	"version-history",
}

// InitClientBinaryCommands is like InitClientCommands only it adds overrides
// and placeholder shims commands that ought to only exist in the vault-client
// binary. Use InitClientCommands for all of the real Client commands.
func InitClientBinaryCommands(ui, serverCmdUi cli.Ui, runOpts *base.RunOptions) map[string]cli.CommandFactory {
	return initClientBinaryCommands(ui, serverCmdUi, runOpts)
}

func initClientBinaryCommands(ui, serverCmdUi cli.Ui, runOpts *base.RunOptions) map[string]cli.CommandFactory {
	commands := InitClientCommands(ui, serverCmdUi, runOpts)

	// Override version so that it shows Vault Client instead of Vault.
	commands["version"] = func() (cli.Command, error) {
		versionInfo := version.GetVersion()
		versionInfo.BinaryName = "Vault Client"
		return &VersionCommand{
			VersionInfo: versionInfo,
			BaseCommand: base.NewBaseCommand(ui, runOpts),
		}, nil
	}

	// Override operator so its examples don't advertise the unavailable
	// operator subcommands.
	commands["operator"] = func() (cli.Command, error) {
		return &OperatorCommand{
			BaseCommand: base.NewBaseCommand(ui, runOpts),
			clientHelp:  true,
		}, nil
	}

	// Add unavailable command shims for commands that exist in server but
	// do not exist in client.
	for _, name := range unavailableClientCommands {
		commands[name] = func() (cli.Command, error) {
			return &UnavailableCommand{ui: ui, name: name}, nil
		}
	}

	return commands
}

// InitClientCommands returns the set of CLI commands for the client binary,
// as defined by PRD VLT-2015. The vault server binary registers these and
// adds server-only commands on top of them.
func InitClientCommands(ui, serverCmdUi cli.Ui, runOpts *base.RunOptions) map[string]cli.CommandFactory {
	getBaseCommand := func() *base.BaseCommand {
		return base.NewBaseCommand(ui, runOpts)
	}

	binaryName := getBaseCommand().BinaryName
	loginHandlers := LoginHandlers(binaryName)

	commands := map[string]cli.CommandFactory{
		"agent": func() (cli.Command, error) {
			return &AgentCommand{
				BaseCommand: &base.BaseCommand{
					UI:         serverCmdUi,
					BinaryName: binaryName,
				},
				ShutdownCh: MakeShutdownCh(),
				SighupCh:   MakeSighupCh(),
				SigUSR2Ch:  MakeSigUSR2Ch(),
			}, nil
		},
		"agent generate-config": func() (cli.Command, error) {
			return &AgentGenerateConfigCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"auth tune": func() (cli.Command, error) {
			return &AuthTuneCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"auth": func() (cli.Command, error) {
			return &AuthCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"auth disable": func() (cli.Command, error) {
			return &AuthDisableCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"auth enable": func() (cli.Command, error) {
			return &AuthEnableCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"auth help": func() (cli.Command, error) {
			return &AuthHelpCommand{
				BaseCommand: getBaseCommand(),
				Handlers:    loginHandlers,
			}, nil
		},
		"auth list": func() (cli.Command, error) {
			return &AuthListCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"auth move": func() (cli.Command, error) {
			return &AuthMoveCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"delete": func() (cli.Command, error) {
			return &DeleteCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"lease": func() (cli.Command, error) {
			return &LeaseCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"lease renew": func() (cli.Command, error) {
			return &LeaseRenewCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"lease lookup": func() (cli.Command, error) {
			return &LeaseLookupCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"lease revoke": func() (cli.Command, error) {
			return &LeaseRevokeCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"list": func() (cli.Command, error) {
			return &ListCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"login": func() (cli.Command, error) {
			return &LoginCommand{
				BaseCommand: getBaseCommand(),
				Handlers:    loginHandlers,
			}, nil
		},
		"namespace": func() (cli.Command, error) {
			return &NamespaceCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"namespace list": func() (cli.Command, error) {
			return &NamespaceListCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"namespace lookup": func() (cli.Command, error) {
			return &NamespaceLookupCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"namespace create": func() (cli.Command, error) {
			return &NamespaceCreateCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"namespace patch": func() (cli.Command, error) {
			return &NamespacePatchCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"namespace delete": func() (cli.Command, error) {
			return &NamespaceDeleteCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"namespace lock": func() (cli.Command, error) {
			return &NamespaceAPILockCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"namespace unlock": func() (cli.Command, error) {
			return &NamespaceAPIUnlockCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator": func() (cli.Command, error) {
			return &OperatorCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator key-status": func() (cli.Command, error) {
			return &OperatorKeyStatusCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator rotate": func() (cli.Command, error) {
			return &OperatorRotateCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator unseal": func() (cli.Command, error) {
			return &OperatorUnsealCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"patch": func() (cli.Command, error) {
			return &PatchCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"path-help": func() (cli.Command, error) {
			return &PathHelpCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"proxy": func() (cli.Command, error) {
			return &ProxyCommand{
				BaseCommand: &base.BaseCommand{
					UI:         serverCmdUi,
					BinaryName: binaryName,
				},
				ShutdownCh: MakeShutdownCh(),
				SighupCh:   MakeSighupCh(),
				SigUSR2Ch:  MakeSigUSR2Ch(),
			}, nil
		},
		"print": func() (cli.Command, error) {
			return &PrintCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"print token": func() (cli.Command, error) {
			return &PrintTokenCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"read": func() (cli.Command, error) {
			return &ReadCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"secrets": func() (cli.Command, error) {
			return &SecretsCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"secrets disable": func() (cli.Command, error) {
			return &SecretsDisableCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"secrets enable": func() (cli.Command, error) {
			return &SecretsEnableCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"secrets list": func() (cli.Command, error) {
			return &SecretsListCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"secrets move": func() (cli.Command, error) {
			return &SecretsMoveCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"secrets tune": func() (cli.Command, error) {
			return &SecretsTuneCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"ssh": func() (cli.Command, error) {
			return &SSHCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"status": func() (cli.Command, error) {
			return &StatusCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"token": func() (cli.Command, error) {
			return &TokenCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"token create": func() (cli.Command, error) {
			return &TokenCreateCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"token capabilities": func() (cli.Command, error) {
			return &TokenCapabilitiesCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"token lookup": func() (cli.Command, error) {
			return &TokenLookupCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"token renew": func() (cli.Command, error) {
			return &TokenRenewCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"token revoke": func() (cli.Command, error) {
			return &TokenRevokeCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"unwrap": func() (cli.Command, error) {
			return &UnwrapCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &VersionCommand{
				VersionInfo: version.GetVersion(),
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"write": func() (cli.Command, error) {
			return &WriteCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv": func() (cli.Command, error) {
			return &KVCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv put": func() (cli.Command, error) {
			return &KVPutCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv patch": func() (cli.Command, error) {
			return &KVPatchCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv rollback": func() (cli.Command, error) {
			return &KVRollbackCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv get": func() (cli.Command, error) {
			return &KVGetCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv delete": func() (cli.Command, error) {
			return &KVDeleteCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv list": func() (cli.Command, error) {
			return &KVListCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv destroy": func() (cli.Command, error) {
			return &KVDestroyCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv undelete": func() (cli.Command, error) {
			return &KVUndeleteCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv enable-versioning": func() (cli.Command, error) {
			return &KVEnableVersioningCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv metadata": func() (cli.Command, error) {
			return &KVMetadataCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv metadata put": func() (cli.Command, error) {
			return &KVMetadataPutCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv metadata patch": func() (cli.Command, error) {
			return &KVMetadataPatchCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv metadata get": func() (cli.Command, error) {
			return &KVMetadataGetCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv metadata delete": func() (cli.Command, error) {
			return &KVMetadataDeleteCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
	}

	commands["tpm"] = func() (cli.Command, error) {
		return &TPMCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}

	commands["tpm attest"] = func() (cli.Command, error) {
		return &TPMAttestCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}

	commands["tpm inspect"] = func() (cli.Command, error) {
		return &TPMInspectCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}

	commands["tpm reattest"] = func() (cli.Command, error) {
		return &TPMReattestCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}

	commands["tpm cleanup"] = func() (cli.Command, error) {
		return &TPMCleanupCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}

	commands["tpm ek"] = func() (cli.Command, error) {
		return &TPMEkCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}

	commands["tpm enroll"] = func() (cli.Command, error) {
		return &TPMEnrollCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}

	entInitCommands(ui, serverCmdUi, runOpts, commands)
	InitHCPCommands(ui, commands, binaryName)

	return commands
}

// InitHCPCommands adds HCP-specific commands to the given commands map.
// It checks for conflicts and only adds a command if it doesn't already exist.
// Their help and hcp connect's output name binaryName instead of "vault".
func InitHCPCommands(ui cli.Ui, commands map[string]cli.CommandFactory, binaryName string) {
	for cmd, cmdFactory := range hcpvlib.InitHCPCommand(ui) {
		// check for conflicts and only put command in the map in case it doesn't conflict with existing one
		_, ok := commands[cmd]
		if !ok {
			commands[cmd] = func() (cli.Command, error) {
				command, err := cmdFactory()
				if err != nil || binaryName == "" || binaryName == "vault" {
					return command, err
				}
				switch c := command.(type) {
				case *hcpvlib.HCPConnectCommand:
					c.Ui = &hcpConnectUi{Ui: c.Ui, BinaryName: binaryName}
					return &hcpConnectCommand{HCPConnectCommand: c, BinaryName: binaryName}, nil
				case *hcpvlib.HCPDisconnectCommand:
					return &hcpDisconnectCommand{HCPDisconnectCommand: c, BinaryName: binaryName}, nil
				default:
					return command, nil
				}
			}
		} else {
			ui.Error("Failed to initialize HCP commands.")
			break
		}
	}
}

// MakeShutdownCh returns a channel that can be used for shutdown
// notifications for commands. This channel will send a message for every
// SIGINT or SIGTERM received.
func MakeShutdownCh() chan struct{} {
	resultCh := make(chan struct{})

	shutdownCh := make(chan os.Signal, 4)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-shutdownCh
		close(resultCh)
	}()
	return resultCh
}

// MakeSighupCh returns a channel that can be used for SIGHUP
// reloading. This channel will send a message for every
// SIGHUP received.
func MakeSighupCh() chan struct{} {
	resultCh := make(chan struct{})

	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, syscall.SIGHUP)
	go func() {
		for {
			<-signalCh
			resultCh <- struct{}{}
		}
	}()
	return resultCh
}
