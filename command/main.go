package command

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/mitchellh/cli"
)

func Run(args []string) int {
	// Handle -v shorthand
	for _, arg := range args {
		if arg == "--" {
			break
		}

		if arg == "-v" || arg == "-version" || arg == "--version" {
			args = []string{"version"}
			break
		}
	}

	cli := &cli.CLI{
		Name:     "vault",
		Args:     args,
		Commands: Commands,

		HelpFunc: FilterDeprecatedFunc(
			FilterCommandFunc("version",
				groupedHelpFunc(
					cli.BasicHelpFunc("vault"),
				),
			),
		),

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

func FilterCommandFunc(name string, f cli.HelpFunc) cli.HelpFunc {
	return func(commands map[string]cli.CommandFactory) string {
		newCommands := make(map[string]cli.CommandFactory, len(commands))
		for k, v := range commands {
			if k != name {
				newCommands[k] = v
			}
		}
		return f(newCommands)
	}
}

// FilterDeprecatedFunc filters deprecated
func FilterDeprecatedFunc(f cli.HelpFunc) cli.HelpFunc {
	return func(commands map[string]cli.CommandFactory) string {
		newCommands := make(map[string]cli.CommandFactory)

		for k, cmdFn := range commands {
			command, err := cmdFn()
			if err != nil {
				log.Printf("[ERR] cli: Command %q failed to load: %s", k, err)
			}

			if _, ok := command.(*DeprecatedCommand); ok {
				continue
			}

			newCommands[k] = cmdFn
		}

		return f(newCommands)
	}
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
