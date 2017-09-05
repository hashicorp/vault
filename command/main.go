package command

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

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
		HelpFunc: helpFunc,

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

// helpFunc is a cli.HelpFunc that can is used to output the help for Vault.
func helpFunc(commands map[string]cli.CommandFactory) string {
	commonNames := map[string]struct{}{
		"delete": struct{}{},
		"read":   struct{}{},
		"renew":  struct{}{},
		"revoke": struct{}{},
		"server": struct{}{},
		"status": struct{}{},
		"unwrap": struct{}{},
		"write":  struct{}{},
	}

	// Determine the maximum key length, and classify based on type
	commonCommands := make(map[string]cli.CommandFactory)
	otherCommands := make(map[string]cli.CommandFactory)

	commonKeyLen, otherKeyLen := 0, 0
	for key, f := range commands {
		if _, ok := commonNames[key]; ok {
			if len(key) > commonKeyLen {
				commonKeyLen = len(key)
			}
			commonCommands[key] = f
		} else {
			if len(key) > otherKeyLen {
				otherKeyLen = len(key)
			}
			otherCommands[key] = f
		}
	}

	var buf bytes.Buffer
	buf.WriteString("Usage: vault <command> [args]\n\n")
	buf.WriteString("Common commands:\n\n")
	buf.WriteString(listCommands(commonCommands, commonKeyLen))
	buf.WriteString("\n")
	buf.WriteString("Other commands:\n\n")
	buf.WriteString(listCommands(otherCommands, otherKeyLen))
	return strings.TrimSpace(buf.String())
}

// listCommands just lists the commands in the map with the
// given maximum key length.
func listCommands(commands map[string]cli.CommandFactory, maxKeyLen int) string {
	var buf bytes.Buffer

	// Get the list of keys so we can sort them, and also get the maximum
	// key length so they can be aligned properly.
	keys := make([]string, 0, len(commands))
	for key, _ := range commands {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		commandFunc, ok := commands[key]
		if !ok {
			// This should never happen since we JUST built the list of
			// keys.
			panic("command not found: " + key)
		}

		command, err := commandFunc()
		if err != nil {
			panic(fmt.Sprintf("command '%s' failed to load: %s", key, err))
		}

		key = fmt.Sprintf("%s%s", key, strings.Repeat(" ", maxKeyLen-len(key)))
		buf.WriteString(fmt.Sprintf("    %s    %s\n", key, command.Synopsis()))
	}

	return buf.String()
}
