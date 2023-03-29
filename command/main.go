// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/token"
	colorable "github.com/mattn/go-colorable"
	"github.com/mitchellh/cli"
)

type VaultUI struct {
	cli.Ui
	format   string
	detailed bool
}

const (
	globalFlagOutputCurlString = "output-curl-string"
	globalFlagOutputPolicy     = "output-policy"
	globalFlagFormat           = "format"
	globalFlagDetailed         = "detailed"
)

var globalFlags = []string{
	globalFlagOutputCurlString, globalFlagOutputPolicy, globalFlagFormat, globalFlagDetailed,
}

// setupEnv parses args and may replace them and sets some env vars to known
// values based on format options
func setupEnv(args []string) (retArgs []string, gf parsedGlobalFlags) {
	var err error
	var nextArgFormat bool
	var haveDetailed bool

	for _, arg := range args {
		retArgs = append(retArgs, arg)
		if nextArgFormat {
			nextArgFormat = false
			gf.format = arg
			retArgs = retArgs[:len(retArgs)-1]
			continue
		}

		if arg == "--" {
			break
		}

		if len(args) == 1 && (arg == "-v" || arg == "-version" || arg == "--version") {
			args = []string{"version"}
			break
		}

		if isGlobalFlag(arg, globalFlagOutputCurlString) {
			gf.outputCurlString = true
			retArgs = retArgs[:len(retArgs)-1]
			continue
		}

		if isGlobalFlag(arg, globalFlagOutputPolicy) {
			gf.outputPolicy = true
			retArgs = retArgs[:len(retArgs)-1]
			continue
		}

		// Parse the 'format' flag, which overrides the env var
		if isGlobalFlagWithValue(arg, globalFlagFormat) {
			gf.format = getGlobalFlagValue(arg)
			retArgs = retArgs[:len(retArgs)-1]
			continue
		}
		// For backwards compat, it could be specified without an equal sign
		if isGlobalFlag(arg, globalFlagFormat) {
			nextArgFormat = true
			retArgs = retArgs[:len(retArgs)-1]
			continue
		}

		// Parse the 'detailed' flag, which overrides the env var
		if isGlobalFlagWithValue(arg, globalFlagDetailed) {
			gf.detailed, err = strconv.ParseBool(getGlobalFlagValue(globalFlagDetailed))
			if err != nil {
				gf.detailed = false
			}
			haveDetailed = true
			retArgs = retArgs[:len(retArgs)-1]
			continue
		}
		// For backwards compat, it could be specified without an equal sign to enable
		// detailed output.
		if isGlobalFlag(arg, globalFlagDetailed) {
			gf.detailed = true
			haveDetailed = true
			retArgs = retArgs[:len(retArgs)-1]
			continue
		}
	}

	envVaultFormat := os.Getenv(EnvVaultFormat)
	// If we did not parse a value, fetch the env var
	if gf.format == "" && envVaultFormat != "" {
		gf.format = envVaultFormat
	}
	// Lowercase for consistency
	gf.format = strings.ToLower(gf.format)
	if gf.format == "" {
		gf.format = "table"
	}

	envVaultDetailed := os.Getenv(EnvVaultDetailed)
	// If we did not parse a value, fetch the env var
	if !haveDetailed && envVaultDetailed != "" {
		gf.detailed, err = strconv.ParseBool(envVaultDetailed)
		if err != nil {
			gf.detailed = false
		}
	}

	return retArgs, gf
}

func isGlobalFlag(arg string, flag string) bool {
	return arg == "-"+flag || arg == "--"+flag
}

func isGlobalFlagWithValue(arg string, flag string) bool {
	return strings.HasPrefix(arg, "--"+flag+"=") || strings.HasPrefix(arg, "-"+flag+"=")
}

func getGlobalFlagValue(arg string) string {
	_, value, _ := strings.Cut(arg, "=")

	return value
}

type RunOptions struct {
	TokenHelper token.TokenHelper
	Stdout      io.Writer
	Stderr      io.Writer
	Address     string
	Client      *api.Client
}

func Run(args []string) int {
	return RunCustom(args, nil)
}

// RunCustom allows passing in a base command template to pass to other
// commands. Currently, this is only used for setting a custom token helper.
func RunCustom(args []string, runOpts *RunOptions) int {
	if runOpts == nil {
		runOpts = &RunOptions{}
	}

	args, gf := setupEnv(args)

	// Don't use color if disabled
	useColor := true
	if os.Getenv(EnvVaultCLINoColor) != "" || color.NoColor {
		useColor = false
	}

	if runOpts.Stdout == nil {
		runOpts.Stdout = os.Stdout
	}
	if runOpts.Stderr == nil {
		runOpts.Stderr = os.Stderr
	}

	// Only use colored UI if stdout is a tty, and not disabled
	if useColor && gf.format == "table" {
		if f, ok := runOpts.Stdout.(*os.File); ok {
			runOpts.Stdout = colorable.NewColorable(f)
		}
		if f, ok := runOpts.Stderr.(*os.File); ok {
			runOpts.Stderr = colorable.NewColorable(f)
		}
	} else {
		runOpts.Stdout = colorable.NewNonColorable(runOpts.Stdout)
		runOpts.Stderr = colorable.NewNonColorable(runOpts.Stderr)
	}

	uiErrWriter := runOpts.Stderr
	if gf.outputCurlString || gf.outputPolicy {
		uiErrWriter = &bytes.Buffer{}
	}

	ui := &VaultUI{
		Ui: &cli.ColoredUi{
			ErrorColor: cli.UiColorRed,
			WarnColor:  cli.UiColorYellow,
			Ui: &cli.BasicUi{
				Reader:      bufio.NewReader(os.Stdin),
				Writer:      runOpts.Stdout,
				ErrorWriter: uiErrWriter,
			},
		},
		format:   gf.format,
		detailed: gf.detailed,
	}

	serverCmdUi := &VaultUI{
		Ui: &cli.ColoredUi{
			ErrorColor: cli.UiColorRed,
			WarnColor:  cli.UiColorYellow,
			Ui: &cli.BasicUi{
				Reader: bufio.NewReader(os.Stdin),
				Writer: runOpts.Stdout,
			},
		},
		format: gf.format,
	}

	if _, ok := Formatters[gf.format]; !ok {
		ui.Error(fmt.Sprintf("Invalid output format: %s", gf.format))
		return 1
	}

	commands := initCommands(ui, serverCmdUi, runOpts, gf)

	hiddenCommands := []string{"version"}

	cli := &cli.CLI{
		Name:     "vault",
		Args:     args,
		Commands: commands,
		HelpFunc: groupedHelpFunc(
			cli.BasicHelpFunc("vault"),
		),
		HelpWriter:                 runOpts.Stdout,
		ErrorWriter:                runOpts.Stderr,
		HiddenCommands:             hiddenCommands,
		Autocomplete:               true,
		AutocompleteNoDefaultFlags: true,
	}

	exitCode, err := cli.Run()
	if gf.outputCurlString {
		return generateCurlString(exitCode, runOpts, uiErrWriter.(*bytes.Buffer))
	} else if gf.outputPolicy {
		return generatePolicy(exitCode, runOpts, uiErrWriter.(*bytes.Buffer))
	} else if err != nil {
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
	"agent",
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

func generateCurlString(exitCode int, runOpts *RunOptions, preParsingErrBuf *bytes.Buffer) int {
	if exitCode == 0 {
		fmt.Fprint(runOpts.Stderr, "Could not generate cURL command")
		return 1
	}

	if api.LastOutputStringError == nil {
		if exitCode == 127 {
			// Usage, just pass it through
			return exitCode
		}
		runOpts.Stderr.Write(preParsingErrBuf.Bytes())
		runOpts.Stderr.Write([]byte("Unable to generate cURL string from command\n"))
		return exitCode
	}

	cs, err := api.LastOutputStringError.CurlString()
	if err != nil {
		runOpts.Stderr.Write([]byte(fmt.Sprintf("Error creating request string: %s\n", err)))
		return 1
	}

	runOpts.Stdout.Write([]byte(fmt.Sprintf("%s\n", cs)))
	return 0
}

func generatePolicy(exitCode int, runOpts *RunOptions, preParsingErrBuf *bytes.Buffer) int {
	if exitCode == 0 {
		fmt.Fprint(runOpts.Stderr, "Could not generate policy")
		return 1
	}

	if api.LastOutputPolicyError == nil {
		if exitCode == 127 {
			// Usage, just pass it through
			return exitCode
		}
		runOpts.Stderr.Write(preParsingErrBuf.Bytes())
		runOpts.Stderr.Write([]byte("Unable to generate policy from command\n"))
		return exitCode
	}

	hcl, err := api.LastOutputPolicyError.HCLString()
	if err != nil {
		runOpts.Stderr.Write([]byte(fmt.Sprintf("Error assembling policy HCL: %s\n", err)))
		return 1
	}

	runOpts.Stdout.Write([]byte(fmt.Sprintf("%s\n", hcl)))
	return 0
}
