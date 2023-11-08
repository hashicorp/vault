package command_server

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command"
	"github.com/mattn/go-colorable"
	"github.com/mitchellh/cli"
	"os"
)

func Run(args []string) int {
	return RunCustom(args, nil)
}

// RunCustom allows passing in a base command template to pass to other
// commands. Currently, this is only used for setting a custom token helper.
func RunCustom(args []string, runOpts *command.RunOptions) int {
	if runOpts == nil {
		runOpts = &command.RunOptions{}
	}

	var format string
	var detailed bool
	var outputCurlString bool
	var outputPolicy bool
	args, format, detailed, outputCurlString, outputPolicy = command.SetupEnv(args)

	// Don't use color if disabled
	useColor := true
	if os.Getenv(command.EnvVaultCLINoColor) != "" || color.NoColor {
		useColor = false
	}

	if runOpts.Stdout == nil {
		runOpts.Stdout = os.Stdout
	}
	if runOpts.Stderr == nil {
		runOpts.Stderr = os.Stderr
	}

	// Only use colored UI if stdout is a tty, and not disabled
	if useColor && format == "table" {
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
	if outputCurlString || outputPolicy {
		uiErrWriter = &bytes.Buffer{}
	}

	ui := &command.VaultUI{
		Ui: &cli.ColoredUi{
			ErrorColor: cli.UiColorRed,
			WarnColor:  cli.UiColorYellow,
			Ui: &cli.BasicUi{
				Reader:      bufio.NewReader(os.Stdin),
				Writer:      runOpts.Stdout,
				ErrorWriter: uiErrWriter,
			},
		},
		Format:   format,
		Detailed: detailed,
	}

	serverCmdUi := &command.VaultUI{
		Ui: &cli.ColoredUi{
			ErrorColor: cli.UiColorRed,
			WarnColor:  cli.UiColorYellow,
			Ui: &cli.BasicUi{
				Reader: bufio.NewReader(os.Stdin),
				Writer: runOpts.Stdout,
			},
		},
		Format: format,
	}

	if _, ok := command.Formatters[format]; !ok {
		ui.Error(fmt.Sprintf("Invalid output format: %s", format))
		return 1
	}

	commands := command.InitCommands(ui, runOpts)
	for s, c := range initServerCommands(serverCmdUi, runOpts) {
		commands[s] = c
	}

	hiddenCommands := []string{"version"}

	cli := &cli.CLI{
		Name:     "vault",
		Args:     args,
		Commands: commands,
		HelpFunc: command.GroupedHelpFunc(
			cli.BasicHelpFunc("vault"),
		),
		HelpWriter:                 runOpts.Stdout,
		ErrorWriter:                runOpts.Stderr,
		HiddenCommands:             hiddenCommands,
		Autocomplete:               true,
		AutocompleteNoDefaultFlags: true,
	}

	exitCode, err := cli.Run()
	if outputCurlString {
		return generateCurlString(exitCode, runOpts, uiErrWriter.(*bytes.Buffer))
	} else if outputPolicy {
		return generatePolicy(exitCode, runOpts, uiErrWriter.(*bytes.Buffer))
	} else if err != nil {
		fmt.Fprintf(runOpts.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}

func generateCurlString(exitCode int, runOpts *command.RunOptions, preParsingErrBuf *bytes.Buffer) int {
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

func generatePolicy(exitCode int, runOpts *command.RunOptions, preParsingErrBuf *bytes.Buffer) int {
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
