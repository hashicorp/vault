package command

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/propagators"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
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
	format string
}

// setupEnv parses args and may replace them and sets some env vars to known
// values based on format options
func setupEnv(args []string) (retArgs []string, format string, outputCurlString, openTelemetry bool) {
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

		if len(args) == 1 && (arg == "-v" || arg == "-version" || arg == "--version") {
			args = []string{"version"}
			break
		}

		if arg == "-output-curl-string" {
			outputCurlString = true
			continue
		}

		if arg == "-open-telemetry" {
			openTelemetry = true
			continue
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

	return args, format, outputCurlString, openTelemetry
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

	var format string
	var outputCurlString bool
	var openTelemetryEnable bool
	args, format, outputCurlString, openTelemetryEnable = setupEnv(args)

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
	if outputCurlString {
		uiErrWriter = ioutil.Discard
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
		format: format,
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
		format: format,
	}

	if _, ok := Formatters[format]; !ok {
		ui.Error(fmt.Sprintf("Invalid output format: %s", format))
		return 1
	}

	initCommands(ui, serverCmdUi, runOpts)

	hiddenCommands := []string{"version"}

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

	if openTelemetryEnable {
		// Construct and register an export pipeline using the Jaeger
		// exporter and a span processor.
		flush, err := jaeger.InstallNewPipeline(
			jaeger.WithAgentEndpoint("localhost:6831"),
			jaeger.WithProcess(jaeger.Process{
				ServiceName: "Vault-CLI",
			}))
		if err != nil {
			fmt.Printf("failed to initialize stdout export pipeline: %v", err)
			return 1
		}

		defer flush()

		global.SetTextMapPropagator(otel.NewCompositeTextMapPropagator(b3.B3{}, propagators.TraceContext{}, propagators.Baggage{}))

		// Replace the default HTTP client with one which propagates tracing context
		http.DefaultClient = &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		}
		cleanhttp.EnableOpenTelemetry = true
	}

	exitCode, err := cli.Run()
	if outputCurlString {
		if exitCode == 0 {
			fmt.Fprint(runOpts.Stderr, "Could not generate cURL command")
			return 1
		} else {
			if api.LastOutputStringError == nil {
				if exitCode == 127 {
					// Usage, just pass it through
					return exitCode
				}
				fmt.Fprint(runOpts.Stderr, "cURL command not set by API operation; run without -output-curl-string to see the generated error\n")
				return exitCode
			}
			if api.LastOutputStringError.Error() != api.ErrOutputStringRequest {
				runOpts.Stdout.Write([]byte(fmt.Sprintf("Error creating request string: %s\n", api.LastOutputStringError.Error())))
				return 1
			}
			runOpts.Stdout.Write([]byte(fmt.Sprintf("%s\n", api.LastOutputStringError.CurlString())))
			return 0
		}
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
