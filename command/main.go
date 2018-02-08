package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/mitchellh/cli"
	"golang.org/x/crypto/ssh/terminal"
)

type VaultUI struct {
	cli.Ui
	isTerminal bool
}

func (u *VaultUI) Output(m string) {
	if u.isTerminal {
		u.Ui.Output(m)
	} else {
		getWriterFromUI(u.Ui).Write([]byte(m))
	}
}

func Run(args []string) int {
	color := true

	// Handle -v shorthand
	for _, arg := range args {
		if arg == "--" {
			break
		}

		if arg == "-v" || arg == "-version" || arg == "--version" {
			args = []string{"version"}
			break
		}

		if arg == "-no-color" {
			color = false
		}
	}

	if os.Getenv(EnvVaultCLINoColor) != "" {
		color = false
	}

	isTerminal := terminal.IsTerminal(int(os.Stdout.Fd()))

	ui := &VaultUI{
		Ui: &cli.BasicUi{
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
		isTerminal: isTerminal,
	}
	serverCmdUi := &VaultUI{
		Ui: &cli.BasicUi{
			Writer: os.Stdout,
		},
		isTerminal: isTerminal,
	}

	// Only use colored UI if stdoout is a tty, and not disabled
	if isTerminal && color {
		ui.Ui = &cli.ColoredUi{
			ErrorColor: cli.UiColorRed,
			WarnColor:  cli.UiColorYellow,
			Ui:         ui.Ui,
		}

		serverCmdUi.Ui = &cli.ColoredUi{
			ErrorColor: cli.UiColorRed,
			WarnColor:  cli.UiColorYellow,
			Ui:         serverCmdUi.Ui,
		}
	}

	initCommands(ui, serverCmdUi)

	// Calculate hidden commands from the deprecated ones
	hiddenCommands := make([]string, 0, len(DeprecatedCommands)+1)
	for k := range DeprecatedCommands {
		hiddenCommands = append(hiddenCommands, k)
	}
	hiddenCommands = append(hiddenCommands, "version")

	cli := &cli.CLI{
		Name:     "vault",
		Args:     args,
		Commands: Commands,
		HelpFunc: groupedHelpFunc(
			cli.BasicHelpFunc("vault"),
		),
		HiddenCommands:             hiddenCommands,
		Autocomplete:               true,
		AutocompleteNoDefaultFlags: true,
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}

var commonCommands = []string{
	"read",
	"write",
	"delete",
	"list",
	"login",
	"server",
	"status",
	"unwrap",
}

func groupedHelpFunc(f cli.HelpFunc) cli.HelpFunc {
	return func(commands map[string]cli.CommandFactory) string {
		var b bytes.Buffer
		tw := tabwriter.NewWriter(&b, 0, 2, 6, ' ', 0)

		fmt.Fprintf(tw, "Usage: vault <command> [args]\n\n")
		fmt.Fprintf(tw, "Common commands:\n")
		for _, v := range commonCommands {
			printCommand(tw, v, commands[v])
		}

		otherCommands := make([]string, 0, len(commands))
		for k := range commands {
			found := false
			for _, v := range commonCommands {
				if k == v {
					found = true
					break
				}
			}

			if !found {
				otherCommands = append(otherCommands, k)
			}
		}
		sort.Strings(otherCommands)

		fmt.Fprintf(tw, "\n")
		fmt.Fprintf(tw, "Other commands:\n")
		for _, v := range otherCommands {
			printCommand(tw, v, commands[v])
		}

		tw.Flush()

		return strings.TrimSpace(b.String())
	}
}

func printCommand(w io.Writer, name string, cmdFn cli.CommandFactory) {
	cmd, err := cmdFn()
	if err != nil {
		panic(fmt.Sprintf("failed to load %q command: %s", name, err))
	}
	fmt.Fprintf(w, "    %s\t%s\n", name, cmd.Synopsis())
}
