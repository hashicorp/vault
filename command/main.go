package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/hashicorp/vault/command/token"
	"github.com/mitchellh/cli"
	"golang.org/x/crypto/ssh/terminal"
)

type VaultUI struct {
	cli.Ui
	isTerminal bool
	format     string
}

func (u *VaultUI) Output(m string) {
	if u.isTerminal {
		u.Ui.Output(m)
	} else {
		writer := getWriterFromUI(u.Ui)
		writer.Write([]byte(m))
		writer.Write([]byte("\n"))
	}
}

// setupEnv parses args and may replace them and sets some env vars to known
// values based on format options
func setupEnv(args []string) []string {
	var format string
	var nextArgFormat bool

	for _, arg := range args {
		if nextArgFormat {
			nextArgFormat = false
			format = arg
			continue
		}

		if arg == "--" {
			break
		}

		if arg == "-v" || arg == "-version" || arg == "--version" {
			args = []string{"version"}
			break
		}

		// Parse a given flag here, which overrides the env var
		if strings.HasPrefix(arg, "--format=") {
			format = strings.TrimPrefix(arg, "--format=")
		}
		if strings.HasPrefix(arg, "-format=") {
			format = strings.TrimPrefix(arg, "-format=")
		}
		// For backwards compat, it could be specified without an equal sign
		if arg == "-format" || arg == "--format" {
			nextArgFormat = true
		}
	}

	envVaultFormat := os.Getenv(EnvVaultFormat)
	// If we did not parse a value, fetch the env var
	if format == "" && envVaultFormat != "" {
		format = envVaultFormat
	}
	// Lowercase for consistency
	format = strings.ToLower(format)
	if format == "" {
		format = "table"
	}
	// Put back into the env for later
	os.Setenv(EnvVaultFormat, format)

	return args
}

type RunOptions struct {
	TokenHelper token.TokenHelper
	Stdout      io.Writer
	Stderr      io.Writer
	Address     string
}

func Run(args []string) int {
	return RunCustom(args, nil)
}

// RunCustom allows passing in a base command template to pass to other
// commands. Currenty, this is only used for setting a custom token helper.
func RunCustom(args []string, runOpts *RunOptions) int {
	if runOpts == nil {
		runOpts = &RunOptions{}
	}
	if runOpts.Stdout == nil {
		runOpts.Stdout = os.Stdout
	}
	if runOpts.Stderr == nil {
		runOpts.Stderr = os.Stderr
	}

	args = setupEnv(args)

	// Don't use color if disabled
	color := true
	if os.Getenv(EnvVaultCLINoColor) != "" {
		color = false
	}

	format := format()

	isTerminal := terminal.IsTerminal(int(os.Stdout.Fd()))

	ui := &VaultUI{
		Ui: &cli.BasicUi{
			Writer:      runOpts.Stdout,
			ErrorWriter: runOpts.Stderr,
		},
		isTerminal: isTerminal,
		format:     format,
	}
	serverCmdUi := &VaultUI{
		Ui: &cli.BasicUi{
			Writer: runOpts.Stdout,
		},
		isTerminal: isTerminal,
		format:     format,
	}

	if _, ok := Formatters[format]; !ok {
		ui.Error(fmt.Sprintf("Invalid output format: %s", format))
		return 1
	}

	// Only use colored UI if stdout is a tty, and not disabled
	if isTerminal && color && format == "table" {
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

	initCommands(ui, serverCmdUi, runOpts)

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
		HelpWriter:                 runOpts.Stderr,
		HiddenCommands:             hiddenCommands,
		Autocomplete:               true,
		AutocompleteNoDefaultFlags: true,
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(runOpts.Stderr, "Error executing CLI: %s\n", err.Error())
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
