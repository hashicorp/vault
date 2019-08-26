package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/version"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

const (
	debugIndexVersion = 1
	debugMinDuration  = 1 * time.Minute
	debugMinInterval  = 5 * time.Second
)

var _ cli.Command = (*DebugCommand)(nil)
var _ cli.CommandAutocomplete = (*DebugCommand)(nil)

type DebugCommand struct {
	*BaseCommand

	flagCompress bool
	flagDuration time.Duration
	flagInterval time.Duration
	flagOutput   string
	flagTargets  []string

	ShutdownCh chan struct{}
}

func (c *DebugCommand) AutocompleteArgs() complete.Predictor {
	// Predict targets
	return c.PredictVaultDebugTargets()
}

func (c *DebugCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *DebugCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:    "compress",
		Target:  &c.flagCompress,
		Default: true,
		Usage:   "Toggles whether to compress output package",
	})

	f.DurationVar(&DurationVar{
		Name:       "duration",
		Target:     &c.flagDuration,
		Completion: complete.PredictAnything,
		Default:    2 * time.Minute,
		Usage:      "Duration to run the command.",
	})

	f.DurationVar(&DurationVar{
		Name:       "interval",
		Target:     &c.flagInterval,
		Completion: complete.PredictAnything,
		Default:    30 * time.Second,
		Usage:      "The interval to run the command.",
	})

	f.StringVar(&StringVar{
		Name:       "output",
		Target:     &c.flagOutput,
		Completion: complete.PredictAnything,
		Usage:      "Specifies the output path for the debug package.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:   "targets",
		Target: &c.flagTargets,
		Usage: "Comma-separated string or list of targets to capture. Available " +
			"targets are: config, host, metrics, pprof, " +
			"replication-status, server-status.",
	})

	return set
}

func (c *DebugCommand) Help() string {
	helpText := `
Usage: vault debug [options]

  Probes a specific Vault server node for a specified period of time, recording
  information about the node, its cluster, and its host environment. The
  information collected is packaged and written to the specified path.

  Certain endpoints that this command issues require ACL permissions to access.
  If not permitted, the information from these endpoints will not be part of the
  output. The command uses the Vault address and token as specified via
  the login command, environment variables, or CLI flags.

  To create a debug package using default duration and interval values in the 
  current directory that captures all applicable targets:

  $ vault debug

  To create a debug package with a specific duration and interval in the current
  directory that capture all applicable targets:

  $ vault debug -duration=10m -interval=1m

  To create a debug package in the current directory with a specific sub-set of
  targets:

  $ vault debug -targets=host,metrics

` + c.Flags().Help()

	return helpText
}

func (c *DebugCommand) Run(args []string) int {
	// Copy the raw args to output in the index
	rawArgs := append([]string(nil), args...)

	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) > 0 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0, got %d)", len(args)))
		return 1
	}

	// Guard duration and interval values to acceptable values
	if c.flagDuration < debugMinDuration {
		c.flagDuration = debugMinDuration
	}
	if c.flagInterval < debugMinInterval {
		c.flagInterval = debugMinInterval
	}
	if c.flagInterval > c.flagDuration {
		c.flagInterval = c.flagDuration
	}

	if len(c.flagTargets) == 0 {
		c.flagTargets = c.defaultTargets()
	}

	captureTime := time.Now().UTC()
	if len(c.flagOutput) == 0 {
		formattedTime := captureTime.Format("2006-01-02T15-04-05Z")
		// TODO: Remove /tmp prefix
		c.flagOutput = fmt.Sprintf("/tmp/vault-debug-%s", formattedTime)
	}

	if _, err := os.Stat(c.flagOutput); os.IsNotExist(err) {
		err := os.MkdirAll(c.flagOutput, 0755)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Unable to create output directory: %s", err))
			return 1
		}
	} else {
		c.UI.Error(fmt.Sprintf("Output directory already exists: %s", c.flagOutput))
		return 1
	}

	// Populate initial index fields
	debugIndex := &debugIndex{
		ClientVersion: version.GetVersion().VersionNumber(),
		Compress:      c.flagCompress,
		Duration:      int(c.flagDuration.Seconds()),
		Interval:      int(c.flagInterval.Seconds()),
		RawArgs:       rawArgs,
		Version:       debugIndexVersion,
		Targets:       c.flagTargets,
		Timestamp:     captureTime,
		Output:        make(map[string]interface{}),
		Errors:        []*captureError{},
	}

	// Print debug information
	c.UI.Output("==> Starting debug capture...")
	c.UI.Info(fmt.Sprintf("     Client Version: %s", debugIndex.ClientVersion))
	c.UI.Info(fmt.Sprintf("           Duration: %s", c.flagDuration))
	c.UI.Info(fmt.Sprintf("           Interval: %s", c.flagInterval))
	c.UI.Info(fmt.Sprintf("            Targets: %s", strings.Join(c.flagTargets, ", ")))
	c.UI.Info(fmt.Sprintf("              Ouput: %s", c.flagOutput))

	// Capture static information
	if err := c.captureStaticTargets(debugIndex); err != nil {
		c.UI.Error(fmt.Sprintf("Error capturing static information: %s", err))
		return 2
	}

	// Capture polling information
	if err := c.capturePollingTargets(debugIndex); err != nil {
		c.UI.Error(fmt.Sprintf("Error capturing dynamic information: %s", err))
		return 2
	}

	// Marshal and write index.js
	bytes, err := json.MarshalIndent(debugIndex, "", "  ")
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error marshalling index: %s", err))
		return 1
	}
	if err := ioutil.WriteFile(filepath.Join(c.flagOutput, "index.js"), bytes, 0644); err != nil {
		c.UI.Error(fmt.Sprintf("Unable to write index.js file: %s", err))
		return 1
	}

	// TODO: Perform compression
	if c.flagCompress {

	}

	c.UI.Info(fmt.Sprintf("Success! Bundle written to: %s", c.flagOutput))
	return 0
}

func (c *DebugCommand) Synopsis() string {
	return "Runs the debug command"
}

func (c *DebugCommand) defaultTargets() []string {
	return []string{"config", "metrics", "pprof", "replication-status", "server-status"}
}

func (c *DebugCommand) captureStaticTargets(index *debugIndex) error {
	c.UI.Info("==> Capturing static information...")
	return nil
}

func (c *DebugCommand) capturePollingTargets(index *debugIndex) error {
	durationCh := time.After(c.flagDuration)
	errCh := make(chan *captureError)
	defer close(errCh)

	var wg sync.WaitGroup
	var serverStatusCollection []*serverStatus

	client, err := c.Client()
	if err != nil {
		return err
	}

	captureInterval := func() {
		currentTimestamp := time.Now().UTC()

		if strutil.StrListContains(c.flagTargets, "config") {

		}
		if strutil.StrListContains(c.flagTargets, " metrics") {

		}
		if strutil.StrListContains(c.flagTargets, " pprof") {

		}
		if strutil.StrListContains(c.flagTargets, "replication-status") {

		}
		if strutil.StrListContains(c.flagTargets, "server-status") {
			c.UI.Info(fmt.Sprintf("     %s [INFO]: Capturing server-status", currentTimestamp.Format(time.RFC3339)))

			wg.Add(1)
			go func() {
				// Naive approach for now, but we shouldn't have to hold things
				// inmem until the end since we're appending to a file. The
				// challenge is figuring out how to return as a single
				// array of objects so that it's valid JSON.
				healthInfo, err := client.Sys().Health()
				if err != nil {
					errCh <- newCaptureError("server-status.health", err)
				}
				sealInfo, err := client.Sys().SealStatus()
				if err != nil {
					errCh <- newCaptureError("server-status.seal", err)
				}

				entry := &serverStatus{
					Timestamp: currentTimestamp,
					Health:    healthInfo,
					Seal:      sealInfo,
				}
				serverStatusCollection = append(serverStatusCollection, entry)

				wg.Done()
			}()
		}
		wg.Wait()
	}

	// Upon exit write the targets that we've collection its respective files
	// and update the index.
	defer func() {
		output, err := json.MarshalIndent(serverStatusCollection, "", "  ")
		if err != nil {
			c.UI.Error("Error marshaling server-status.json data")
			return
		}
		if err := ioutil.WriteFile(filepath.Join(c.flagOutput, "server-status.json"), output, 0644); err != nil {
			c.UI.Error("Error writing data to server-status.json")
			return
		}

		if index.Output["files"] == nil {
			index.Output["files"] = make([]string, 0)
		}
		index.Output["files"] = append(index.Output["files"].([]string), "server-status.json")
	}()

	// Start capture by capturing the first interval before we hit the first
	// ticker.
	c.UI.Info("==> Capturing dynamic information...")
	go captureInterval()

	// Capture at each interval, until end of duration or interrupt.
	intervalTicker := time.Tick(c.flagInterval)
	for {
		select {
		case err := <-errCh:
			index.Errors = append(index.Errors, err)
		case <-intervalTicker:
			go captureInterval()
		case <-durationCh:
			return nil
		case <-c.ShutdownCh:
			return nil
		}
	}
}

// debugIndex represents the data in the index file
type debugIndex struct {
	Version       int                    `json:"version"`
	ClientVersion string                 `json:"client_version"`
	Timestamp     time.Time              `json:"timestamp"`
	Duration      int                    `json:"duration_seconds"`
	Interval      int                    `json:"interval_seconds"`
	Compress      bool                   `json:"compress"`
	RawArgs       []string               `json:"raw_args"`
	Targets       []string               `json:"targets"`
	Output        map[string]interface{} `json:"output"`
	Errors        []*captureError        `json:"errors"`
}

// captureError hold an error entry that can occur during polling capture.
// It includes the timestamp, the target, and the error itself.
type captureError struct {
	TargetError string    `json:"error"`
	Target      string    `json:"target"`
	Timestamp   time.Time `json:"timestamp"`
}

// newCaptureError instantiates a new captureError.
func newCaptureError(target string, err error) *captureError {
	return &captureError{
		TargetError: err.Error(),
		Target:      target,
		Timestamp:   time.Now().UTC(),
	}
}

// serverStatus holds a single interval entry for the server-status target
type serverStatus struct {
	Timestamp time.Time               `json:"timestamp"`
	Health    *api.HealthResponse     `json:"health"`
	Seal      *api.SealStatusResponse `json:"seal"`
}
