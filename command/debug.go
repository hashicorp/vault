package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	gatedwriter "github.com/hashicorp/vault/helper/gated-writer"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/version"
	"github.com/mholt/archiver"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

const (
	// debugIndexVersion is tracks the canonical version in the index file
	// for compatibility with future format/layout changes on the bundle.
	debugIndexVersion = 1

	// debugMinInterval is the minimum acceptable interval capture value. This
	// value applies to duration and all interval-related flags.
	debugMinInterval = 5 * time.Second

	// debugDurationGrace is the grace period added to duration to allow for
	// "last frame" capture if the interval falls into the last duration time
	// value. For instance, using default values, adding a grace duration lets
	// the command capture 5 intervals (0, 30, 60, 90, and 120th second) before
	// exiting.
	debugDurationGrace = 1 * time.Second

	// debugCompressionExt is the default compression extension used if
	// compression is enabled.
	debugCompressionExt = ".tar.gz"

	// fileFriendlyTimeFormat is the time format used for file and directory
	// naming.
	fileFriendlyTimeFormat = "2006-01-02T15-04-05Z"
)

// debugIndex represents the data structure in the index file
type debugIndex struct {
	Version         int                    `json:"version"`
	VaultAddress    string                 `json:"vault_address"`
	ClientVersion   string                 `json:"client_version"`
	Timestamp       time.Time              `json:"timestamp"`
	Duration        int                    `json:"duration_seconds"`
	Interval        int                    `json:"interval_seconds"`
	MetricsInterval int                    `json:"metrics_interval_seconds"`
	Compress        bool                   `json:"compress"`
	RawArgs         []string               `json:"raw_args"`
	Targets         []string               `json:"targets"`
	Output          map[string]interface{} `json:"output"`
	Errors          []*captureError        `json:"errors"`
}

// captureError holds an error entry that can occur during polling capture.
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

var _ cli.Command = (*DebugCommand)(nil)
var _ cli.CommandAutocomplete = (*DebugCommand)(nil)

type DebugCommand struct {
	*BaseCommand

	flagCompress        bool
	flagDuration        time.Duration
	flagInterval        time.Duration
	flagMetricsInterval time.Duration
	flagOutput          string
	flagTargets         []string

	// debugIndex is used to keep track of the index state, which gets written
	// to a file at the end.
	debugIndex *debugIndex

	// pprofWg ensures that we are only sending a single profile/trace request
	// at a time since the server pprof handler cannot handle concurrent
	// requests of these types.
	pprofWg *sync.WaitGroup

	// pollingWg ensures that polling capture goroutines are done before the
	// bundle gets generated
	pollingWg *sync.WaitGroup

	// skipTimingChecks bypasses timing-related checks, used primarily for tests
	skipTimingChecks bool
	// logger is the logger used for outputting capture progress
	logger hclog.Logger

	// ShutdownCh is used to capture interrupt signal and end polling capture
	ShutdownCh chan struct{}
	// doneCh is used to signal terminal
	doneCh chan struct{}

	// Various channels for receiving state during polling capture
	metricsCh           chan map[string]interface{}
	serverStatusCh      chan map[string]interface{}
	replicationStatusCh chan map[string]interface{}
	hostInfoCh          chan map[string]interface{}
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
		Usage: "The interval in which to perform profiling and server state " +
			"capture, excluding metrics.",
	})

	f.DurationVar(&DurationVar{
		Name:       "metrics-interval",
		Target:     &c.flagMetricsInterval,
		Completion: complete.PredictAnything,
		Default:    10 * time.Second,
		Usage:      "The interval in which to perform metrics capture.",
	})

	f.StringVar(&StringVar{
		Name:       "output",
		Target:     &c.flagOutput,
		Completion: complete.PredictAnything,
		Usage:      "Specifies the output path for the debug package.",
	})

	f.StringSliceVar(&StringSliceVar{
		Name:   "target",
		Target: &c.flagTargets,
		Usage: "Target to capture, defaulting to all if none specified. " +
			"This can be specified multiple times to capture multiple targets. " +
			"Available targets are: config, host, metrics, pprof, " +
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
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	parsedArgs := f.Args()
	if len(parsedArgs) > 0 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0, got %d)", len(parsedArgs)))
		return 1
	}

	// Initialize the logger for debug output
	logWriter := &gatedwriter.Writer{Writer: os.Stderr}
	if c.logger == nil {
		c.logger = logging.NewVaultLoggerWithWriter(logWriter, hclog.Trace)
	}

	client, dstOutputFile, err := c.preflight(args)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error during validation: %s", err))
		return 1
	}

	// Print debug information
	c.UI.Output("==> Starting debug capture...")
	c.UI.Info(fmt.Sprintf("         Vault Address: %s", c.debugIndex.VaultAddress))
	c.UI.Info(fmt.Sprintf("        Client Version: %s", c.debugIndex.ClientVersion))
	c.UI.Info(fmt.Sprintf("              Duration: %s", c.flagDuration))
	c.UI.Info(fmt.Sprintf("              Interval: %s", c.flagInterval))
	c.UI.Info(fmt.Sprintf("      Metrics Interval: %s", c.flagMetricsInterval))
	c.UI.Info(fmt.Sprintf("               Targets: %s", strings.Join(c.flagTargets, ", ")))
	c.UI.Info(fmt.Sprintf("                Output: %s", dstOutputFile))
	c.UI.Output("")

	// Release the log gate.
	logWriter.Flush()

	// Setup variables
	c.doneCh = make(chan struct{})
	c.pprofWg = &sync.WaitGroup{}
	c.pollingWg = &sync.WaitGroup{}

	// Capture static information
	c.UI.Info("==> Capturing static information...")
	if err := c.captureStaticTargets(client); err != nil {
		c.UI.Error(fmt.Sprintf("Error capturing static information: %s", err))
		return 2
	}

	c.UI.Output("")

	// Capture polling information
	c.UI.Info("==> Capturing dynamic information...")
	if err := c.capturePollingTargets(client); err != nil {
		c.UI.Error(fmt.Sprintf("Error capturing dynamic information: %s", err))
		return 2
	}

	c.UI.Output("Finished capturing information, bundling files...")

	// Generate index file
	if err := c.generateIndex(); err != nil {
		c.UI.Error(fmt.Sprintf("Error generating index: %s", err))
		return 1
	}

	if c.flagCompress {
		if err := c.compress(dstOutputFile); err != nil {
			c.UI.Error(fmt.Sprintf("Error encountered during bundle compression: %s", err))
			// We want to inform that data collection was captured and stored in
			// a directory even if compression fails
			c.UI.Info(fmt.Sprintf("Data written to: %s", c.flagOutput))
			return 1
		}
	}

	c.UI.Info(fmt.Sprintf("Success! Bundle written to: %s", dstOutputFile))
	return 0
}

func (c *DebugCommand) Synopsis() string {
	return "Runs the debug command"
}

func (c *DebugCommand) generateIndex() error {
	outputLayout := map[string]interface{}{
		"files": []string{},
	}
	// Walk the directory to generate the output layout
	err := filepath.Walk(c.flagOutput, func(path string, info os.FileInfo, err error) error {
		// Prevent panic by handling failure accessing a path
		if err != nil {
			return err
		}

		// Skip the base dir
		if path == c.flagOutput {
			return nil
		}

		// If we're a directory, simply add a corresponding map
		if info.IsDir() {
			parsedTime, err := time.Parse(fileFriendlyTimeFormat, info.Name())
			if err != nil {
				return err
			}

			outputLayout[info.Name()] = map[string]interface{}{
				"timestamp": parsedTime,
				"files":     []string{},
			}
			return nil
		}

		relPath, err := filepath.Rel(c.flagOutput, path)
		if err != nil {
			return err
		}

		dir, file := filepath.Split(relPath)
		if len(dir) != 0 {
			dir = strings.TrimSuffix(dir, "/")
			filesArr := outputLayout[dir].(map[string]interface{})["files"]
			outputLayout[dir].(map[string]interface{})["files"] = append(filesArr.([]string), file)
		} else {
			outputLayout["files"] = append(outputLayout["files"].([]string), file)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error generating directory output layout: %s", err)
	}

	c.debugIndex.Output = outputLayout

	// Marshal into json
	bytes, err := json.MarshalIndent(c.debugIndex, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling index file: %s", err)
	}

	// Write out file
	if err := ioutil.WriteFile(filepath.Join(c.flagOutput, "index.json"), bytes, 0644); err != nil {
		return fmt.Errorf("error generating index file; %s", err)
	}

	return nil
}

// preflight performs various checks against the provided flags to ensure they
// are valid/reasonable values. It also takes care of instantiating a client and
// index object for use by the command.
func (c *DebugCommand) preflight(rawArgs []string) (*api.Client, string, error) {
	if !c.skipTimingChecks {
		// Guard duration and interval values to acceptable values
		if c.flagDuration < debugMinInterval {
			c.UI.Info(fmt.Sprintf("Overwriting duration value %q to the minimum value of %q", c.flagDuration, debugMinInterval))
			c.flagDuration = debugMinInterval
		}
		if c.flagInterval < debugMinInterval {
			c.UI.Info(fmt.Sprintf("Overwriting inteval value %q to the minimum value of %q", c.flagInterval, debugMinInterval))
			c.flagInterval = debugMinInterval
		}
		if c.flagMetricsInterval < debugMinInterval {
			c.UI.Info(fmt.Sprintf("Overwriting metrics inteval value %q to the minimum value of %q", c.flagMetricsInterval, debugMinInterval))
			c.flagMetricsInterval = debugMinInterval
		}
	}

	// These timing checks are always applicable since interval shouldn't be
	// greater than the duration
	if c.flagInterval > c.flagDuration {
		c.UI.Info(fmt.Sprintf("Overwriting inteval value %q to the duration value %q", c.flagInterval, c.flagDuration))
		c.flagInterval = c.flagDuration
	}
	if c.flagMetricsInterval > c.flagDuration {
		c.UI.Info(fmt.Sprintf("Overwriting metrics inteval value %q to the duration value %q", c.flagMetricsInterval, c.flagDuration))
		c.flagMetricsInterval = c.flagDuration
	}

	if len(c.flagTargets) == 0 {
		c.flagTargets = c.defaultTargets()
	}

	// Make sure we can talk to the server
	client, err := c.Client()
	if err != nil {
		return nil, "", fmt.Errorf("unable to create client to connect to Vault: %s", err)
	}
	if _, err := client.Sys().Health(); err != nil {
		return nil, "", fmt.Errorf("unable to connect to the server: %s", err)
	}

	captureTime := time.Now().UTC()
	if len(c.flagOutput) == 0 {
		formattedTime := captureTime.Format(fileFriendlyTimeFormat)
		c.flagOutput = fmt.Sprintf("vault-debug-%s", formattedTime)
	}

	// Strip trailing slash before proceeding
	c.flagOutput = strings.TrimSuffix(c.flagOutput, "/")

	// If compression is enabled, trim the extension so that the files are
	// written to a directory even if compression somehow fails. We ensure the
	// extension during compression. We also prevent overwriting if the file
	// already exists.
	dstOutputFile := c.flagOutput
	if c.flagCompress {
		if !strings.HasSuffix(dstOutputFile, ".tar.gz") && !strings.HasSuffix(dstOutputFile, ".tgz") {
			dstOutputFile = dstOutputFile + debugCompressionExt
		}

		// Ensure that the file doesn't already exist, and ensure that we always
		// trim the extension from flagOutput since we'll be progressively
		// writing to that.
		if _, err := os.Stat(dstOutputFile); os.IsNotExist(err) {
			c.flagOutput = strings.TrimSuffix(c.flagOutput, ".tar.gz")
			c.flagOutput = strings.TrimSuffix(c.flagOutput, ".tgz")
		} else {
			return nil, "", fmt.Errorf("output file already exists: %s", dstOutputFile)
		}
	}

	// Stat check the directory to ensure we don't override any existing data.
	if _, err := os.Stat(c.flagOutput); os.IsNotExist(err) {
		err := os.MkdirAll(c.flagOutput, 0755)
		if err != nil {
			return nil, "", fmt.Errorf("unable to create output directory: %s", err)
		}
	} else {
		return nil, "", fmt.Errorf("output directory already exists: %s", c.flagOutput)
	}

	// Instantiate the channels for polling
	c.metricsCh = make(chan map[string]interface{})
	c.serverStatusCh = make(chan map[string]interface{})
	c.replicationStatusCh = make(chan map[string]interface{})
	c.hostInfoCh = make(chan map[string]interface{})

	// Populate initial index fields
	c.debugIndex = &debugIndex{
		VaultAddress:    client.Address(),
		ClientVersion:   version.GetVersion().VersionNumber(),
		Compress:        c.flagCompress,
		Duration:        int(c.flagDuration.Seconds()),
		Interval:        int(c.flagInterval.Seconds()),
		MetricsInterval: int(c.flagMetricsInterval.Seconds()),
		RawArgs:         rawArgs,
		Version:         debugIndexVersion,
		Targets:         c.flagTargets,
		Timestamp:       captureTime,
		Errors:          []*captureError{},
	}

	return client, dstOutputFile, nil
}

func (c *DebugCommand) defaultTargets() []string {
	return []string{"config", "host", "metrics", "pprof", "replication-status", "server-status"}
}

func (c *DebugCommand) captureStaticTargets(client *api.Client) error {
	// Capture configuration state
	if strutil.StrListContains(c.flagTargets, "config") {
		c.logger.Info("capturing configuration state")

		resp, err := client.Logical().Read("sys/config/state/sanitized")
		if err != nil {
			captErr := newCaptureError("config", err)
			c.debugIndex.Errors = append(c.debugIndex.Errors, captErr)

			c.logger.Error("config: error capturing config state", "error", err)
		}

		if resp != nil && resp.Data != nil {
			collection := []map[string]interface{}{
				{
					"timestamp": time.Now().UTC(),
					"config":    resp.Data,
				},
			}
			if err := c.persistCollection(collection, "config.json"); err != nil {
				c.UI.Error(fmt.Sprintf("Error writing data to %s: %v", "config.json", err))
			}
		}
	}

	return nil
}

// capturePollingTargets captures all dynamic targets over the specified
// duration and interval.
func (c *DebugCommand) capturePollingTargets(client *api.Client) error {
	startTime := time.Now()
	durationCh := time.After(c.flagDuration + debugDurationGrace)

	// Track current interval count to show progress.
	// totalCount := int(c.flagDuration.Seconds()/c.flagInterval.Seconds()) + 1
	var idxCount, mIdxCount int

	// errCh is used to capture any polling errors and recorded in the index
	errCh := make(chan *captureError)

	var serverStatusCollection []map[string]interface{}
	var replicationStatusCollection []map[string]interface{}
	var metricsCollection []map[string]interface{}
	var hostInfoCollection []map[string]interface{}

	// Start a goroutine to collect data as soon as it's received
	collectCh := make(chan struct{})
	go func() {
		for {
			select {
			case metricsEntry := <-c.metricsCh:
				metricsCollection = append(metricsCollection, metricsEntry)
			case serverStatusEntry := <-c.serverStatusCh:
				serverStatusCollection = append(serverStatusCollection, serverStatusEntry)
			case replicationStatusEntry := <-c.replicationStatusCh:
				replicationStatusCollection = append(replicationStatusCollection, replicationStatusEntry)
			case hostInfoEntry := <-c.hostInfoCh:
				hostInfoCollection = append(hostInfoCollection, hostInfoEntry)
			case <-c.doneCh:
				close(collectCh)
				return
			}
		}
	}()

	// Start capture by capturing the first interval before we hit the first
	// ticker.
	c.pollingWg.Add(1)
	go c.intervalCapture(client, idxCount, startTime, errCh)

	c.pollingWg.Add(1)
	go c.metricsIntervalCapture(client, mIdxCount, errCh)

	// Capture at each interval, until end of duration or interrupt.
	intervalTicker := time.Tick(c.flagInterval)
	metricsIntervalTicker := time.Tick(c.flagMetricsInterval)
POLLING:
	for {
		select {
		case err := <-errCh:
			c.debugIndex.Errors = append(c.debugIndex.Errors, err)
		case <-intervalTicker:
			idxCount++
			c.pollingWg.Add(1)
			go c.intervalCapture(client, idxCount, startTime, errCh)
		case <-metricsIntervalTicker:
			mIdxCount++
			c.pollingWg.Add(1)
			go c.metricsIntervalCapture(client, mIdxCount, errCh)
		case <-durationCh:
			// Wait for polling requests to finish before breaking
			c.pollingWg.Wait()

			break POLLING
		case <-c.ShutdownCh:
			c.UI.Info("Caught interrupt signal, exiting...")
			break POLLING
		}
	}

	// Close the done channel to signal termination, make sure collection
	// goroutine is terminated before proceeding to persisting the info.
	close(c.doneCh)
	<-collectCh

	// Write collected data to their corresponding files
	if err := c.persistCollection(metricsCollection, "metrics.json"); err != nil {
		c.UI.Error(fmt.Sprintf("Error writing data to %s: %v", "metrics.json", err))
	}
	if err := c.persistCollection(serverStatusCollection, "server_status.json"); err != nil {
		c.UI.Error(fmt.Sprintf("Error writing data to %s: %v", "server_status.json", err))
	}
	if err := c.persistCollection(replicationStatusCollection, "replication_status.json"); err != nil {
		c.UI.Error(fmt.Sprintf("Error writing data to %s: %v", "replication_status.json", err))
	}
	if err := c.persistCollection(hostInfoCollection, "host_info.json"); err != nil {
		c.UI.Error(fmt.Sprintf("Error writing data to %s: %v", "host_info.json", err))
	}

	return nil
}

func (c *DebugCommand) intervalCapture(client *api.Client, idxCount int, startTime time.Time, errCh chan<- *captureError) {
	defer c.pollingWg.Done()

	var wg sync.WaitGroup
	currentTimestamp := time.Now().UTC()

	if strutil.StrListContains(c.flagTargets, "host") {
		c.logger.Info("capturing host information", "count", idxCount)

		wg.Add(1)
		go func() {
			r := client.NewRequest("GET", "/v1/sys/host-info")
			resp, err := client.RawRequest(r)
			if err != nil {
				errCh <- newCaptureError("host", err)
			}
			if resp != nil {
				defer resp.Body.Close()

				secret, err := api.ParseSecret(resp.Body)
				if err != nil {
					errCh <- newCaptureError("host", err)
				}
				if hostEntry := secret.Data; hostEntry != nil {
					c.hostInfoCh <- hostEntry
				}
			}
			wg.Done()
		}()
	}

	if strutil.StrListContains(c.flagTargets, "pprof") {
		c.logger.Info("capturing pprof data", "count", idxCount)

		wg.Add(1)
		go func() {
			defer wg.Done()

			// Create a sub-directory for pprof data
			currentDir := currentTimestamp.Format(fileFriendlyTimeFormat)
			dirName := filepath.Join(c.flagOutput, currentDir)
			if err := os.MkdirAll(dirName, 0755); err != nil {
				c.UI.Error(fmt.Sprintf("Error creating sub-directory for time interval: %s", err))
				return
			}

			// Capture goroutines
			data, err := pprofGoroutine(client)
			if err != nil {
				errCh <- newCaptureError("pprof.goroutine", err)
			} else {
				err = ioutil.WriteFile(filepath.Join(dirName, "goroutine.prof"), data, 0644)
				if err != nil {
					errCh <- newCaptureError("pprof.goroutine", err)
				}
			}

			// Capture heap
			data, err = pprofHeap(client)
			if err != nil {
				errCh <- newCaptureError("pprof.heap", err)
			} else {
				err = ioutil.WriteFile(filepath.Join(dirName, "heap.prof"), data, 0644)
				if err != nil {
					errCh <- newCaptureError("pprof.heap", err)
				}
			}

			// If the our remaining duration is less than the interval value
			// skip profile and trace.
			runDuration := currentTimestamp.Sub(startTime)
			if (c.flagDuration+debugDurationGrace)-runDuration < c.flagInterval {
				return
			}

			// Wait until all other profile/trace requests are finished or until
			// we're terminated.
			contCh := make(chan struct{})
			go func() {
				c.pprofWg.Wait()
				close(contCh)
			}()
			select {
			case <-contCh:
			case <-c.doneCh:
				return
			}

			// We want to add 2 at a time to ensure that we capture both at the
			// same interval slice.
			c.pprofWg.Add(2)

			// Capture profile
			go func() {
				defer c.pprofWg.Done()
				data, err := pprofProfile(client, c.flagInterval)
				if err != nil {
					errCh <- newCaptureError("pprof.profile", err)
					return
				}

				err = ioutil.WriteFile(filepath.Join(dirName, "profile.prof"), data, 0644)
				if err != nil {
					errCh <- newCaptureError("pprof.profile", err)
				}
			}()

			// Capture trace
			go func() {
				defer c.pprofWg.Done()
				data, err := pprofTrace(client, c.flagInterval)
				if err != nil {
					errCh <- newCaptureError("pprof.trace", err)
					return
				}

				err = ioutil.WriteFile(filepath.Join(dirName, "trace.out"), data, 0644)
				if err != nil {
					errCh <- newCaptureError("pprof.trace", err)
				}
			}()

			// Wait until profile/trace is done
			c.pprofWg.Wait()
		}()
	}

	if strutil.StrListContains(c.flagTargets, "replication-status") {
		c.logger.Info("capturing replication status", "count", idxCount)

		wg.Add(1)
		go func() {
			r := client.NewRequest("GET", "/v1/sys/replication/status")
			resp, err := client.RawRequest(r)
			if err != nil {
				errCh <- newCaptureError("replication-status", err)
			}
			if resp != nil {
				defer resp.Body.Close()

				secret, err := api.ParseSecret(resp.Body)
				if err != nil {
					errCh <- newCaptureError("replication-status", err)
				}
				if replicationEntry := secret.Data; replicationEntry != nil {
					replicationEntry["timestamp"] = currentTimestamp
					c.replicationStatusCh <- replicationEntry
				}
			}
			wg.Done()
		}()
	}

	if strutil.StrListContains(c.flagTargets, "server-status") {
		c.logger.Info("capturing server status", "count", idxCount)

		wg.Add(1)
		go func() {
			healthInfo, err := client.Sys().Health()
			if err != nil {
				errCh <- newCaptureError("server-status.health", err)
			}
			sealInfo, err := client.Sys().SealStatus()
			if err != nil {
				errCh <- newCaptureError("server-status.seal", err)
			}

			statusEntry := map[string]interface{}{
				"timestamp": currentTimestamp,
				"health":    healthInfo,
				"seal":      sealInfo,
			}
			c.serverStatusCh <- statusEntry

			wg.Done()
		}()
	}

	// Wait for all dynamic information to be captured before returning
	wg.Wait()
}

func (c *DebugCommand) metricsIntervalCapture(client *api.Client, mIdxCount int, errCh chan<- *captureError) {
	defer c.pollingWg.Done()

	if !strutil.StrListContains(c.flagTargets, "metrics") {
		return
	}

	c.logger.Info("capturing metrics", "count", mIdxCount)

	healthStatus, err := client.Sys().Health()
	if err != nil {
		errCh <- newCaptureError("metrics", err)
		return
	}

	// Check replication status. We skip on processing metrics if we're one
	// of the following (since the request will be forwarded):
	// 1. Any type of DR Node
	// 2. Non-DR, non-performance standby nodes
	switch {
	case healthStatus.ReplicationDRMode == "secondary":
		c.logger.Info("skipping metrics capture on DR secondary node")
		return
	case healthStatus.Standby && !healthStatus.PerformanceStandby:
		c.logger.Info("skipping metrics on standby node")
		return
	}

	// Perform metrics request
	r := client.NewRequest("GET", "/v1/sys/metrics")
	resp, err := client.RawRequest(r)
	if err != nil {
		errCh <- newCaptureError("metrics", err)
		return
	}
	if resp != nil {
		defer resp.Body.Close()

		metricsEntry := make(map[string]interface{})
		err := json.NewDecoder(resp.Body).Decode(&metricsEntry)
		if err != nil {
			errCh <- newCaptureError("metrics", err)
		}
		c.metricsCh <- metricsEntry
	}
}

// persistCollection writes the collected data for a particular target onto the
// specified file. If the collection is empty, it returns immediately.
func (c *DebugCommand) persistCollection(collection []map[string]interface{}, outFile string) error {
	if len(collection) == 0 {
		return nil
	}

	// Write server-status file and update the index
	bytes, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(c.flagOutput, outFile), bytes, 0644); err != nil {
		return err
	}

	return nil
}

func (c *DebugCommand) compress(dst string) error {
	tgz := archiver.NewTarGz()
	if err := tgz.Archive([]string{c.flagOutput}, dst); err != nil {
		return fmt.Errorf("failed to compress data: %s", err)
	}

	// If everything is fine up to this point, remove original directory
	if err := os.RemoveAll(c.flagOutput); err != nil {
		return fmt.Errorf("failed to remove data directory: %s", err)
	}

	return nil
}

func pprofGoroutine(client *api.Client) ([]byte, error) {
	req := client.NewRequest("GET", "/v1/sys/pprof/goroutine")
	resp, err := client.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func pprofHeap(client *api.Client) ([]byte, error) {
	req := client.NewRequest("GET", "/v1/sys/pprof/heap")
	resp, err := client.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func pprofProfile(client *api.Client, duration time.Duration) ([]byte, error) {
	seconds := int(duration.Seconds())
	secStr := strconv.Itoa(seconds)

	req := client.NewRequest("GET", "/v1/sys/pprof/profile")
	req.Params.Add("seconds", secStr)
	resp, err := client.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func pprofTrace(client *api.Client, duration time.Duration) ([]byte, error) {
	seconds := int(duration.Seconds())
	secStr := strconv.Itoa(seconds)

	req := client.NewRequest("GET", "/v1/sys/pprof/trace")
	req.Params.Add("seconds", secStr)
	resp, err := client.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
