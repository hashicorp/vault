package command

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/gatedwriter"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/osutil"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/version"
	"github.com/mholt/archiver/v3"
	"github.com/mitchellh/cli"
	"github.com/oklog/run"
	"github.com/posener/complete"
)

const (
	// debugIndexVersion tracks the canonical version in the index file
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
	Version                int                    `json:"version"`
	VaultAddress           string                 `json:"vault_address"`
	ClientVersion          string                 `json:"client_version"`
	ServerVersion          string                 `json:"server_version"`
	Timestamp              time.Time              `json:"timestamp"`
	DurationSeconds        int                    `json:"duration_seconds"`
	IntervalSeconds        int                    `json:"interval_seconds"`
	MetricsIntervalSeconds int                    `json:"metrics_interval_seconds"`
	Compress               bool                   `json:"compress"`
	RawArgs                []string               `json:"raw_args"`
	Targets                []string               `json:"targets"`
	Output                 map[string]interface{} `json:"output"`
	Errors                 []*captureError        `json:"errors"`
}

// captureError holds an error entry that can occur during polling capture.
// It includes the timestamp, the target, and the error itself.
type captureError struct {
	TargetError string    `json:"error"`
	Target      string    `json:"target"`
	Timestamp   time.Time `json:"timestamp"`
}

var (
	_ cli.Command             = (*DebugCommand)(nil)
	_ cli.CommandAutocomplete = (*DebugCommand)(nil)
)

type DebugCommand struct {
	*BaseCommand

	flagCompress        bool
	flagDuration        time.Duration
	flagInterval        time.Duration
	flagMetricsInterval time.Duration
	flagOutput          string
	flagTargets         []string

	// logFormat defines the output format for Monitor
	logFormat string

	// debugIndex is used to keep track of the index state, which gets written
	// to a file at the end.
	debugIndex *debugIndex

	// skipTimingChecks bypasses timing-related checks, used primarily for tests
	skipTimingChecks bool
	// logger is the logger used for outputting capture progress
	logger hclog.Logger

	// ShutdownCh is used to capture interrupt signal and end polling capture
	ShutdownCh chan struct{}

	// Collection slices to hold data
	hostInfoCollection          []map[string]interface{}
	metricsCollection           []map[string]interface{}
	replicationStatusCollection []map[string]interface{}
	serverStatusCollection      []map[string]interface{}
	inFlightReqStatusCollection []map[string]interface{}

	// cachedClient holds the client retrieved during preflight
	cachedClient *api.Client

	// errLock is used to lock error capture into the index file
	errLock sync.Mutex
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
		Usage:      "The polling interval at which to collect profiling data and server state.",
	})

	f.DurationVar(&DurationVar{
		Name:       "metrics-interval",
		Target:     &c.flagMetricsInterval,
		Completion: complete.PredictAnything,
		Default:    10 * time.Second,
		Usage:      "The polling interval at which to collect metrics data.",
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
			"replication-status, server-status, log.",
	})

	f.StringVar(&StringVar{
		Name:    "log-format",
		Target:  &c.logFormat,
		Default: "standard",
		Usage: "Log format to be captured if \"log\" target specified. " +
			"Supported values are \"standard\" and \"json\". The default is \"standard\".",
	})

	return set
}

func (c *DebugCommand) Help() string {
	helpText := `
Usage: vault debug [options]

  Probes a specific Vault server node for a specified period of time, recording
  information about the node, its cluster, and its host environment. The
  information collected is packaged and written to the specified path.

  Certain endpoints that this command uses require ACL permissions to access.
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

  $ vault debug -target=host -target=metrics

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
	gatedWriter := gatedwriter.NewWriter(os.Stderr)
	if c.logger == nil {
		c.logger = logging.NewVaultLoggerWithWriter(gatedWriter, hclog.Trace)
	}

	dstOutputFile, err := c.preflight(args)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error during validation: %s", err))
		return 1
	}

	// Print debug information
	c.UI.Output("==> Starting debug capture...")
	c.UI.Info(fmt.Sprintf("         Vault Address: %s", c.debugIndex.VaultAddress))
	c.UI.Info(fmt.Sprintf("        Client Version: %s", c.debugIndex.ClientVersion))
	c.UI.Info(fmt.Sprintf("        Server Version: %s", c.debugIndex.ServerVersion))
	c.UI.Info(fmt.Sprintf("              Duration: %s", c.flagDuration))
	c.UI.Info(fmt.Sprintf("              Interval: %s", c.flagInterval))
	c.UI.Info(fmt.Sprintf("      Metrics Interval: %s", c.flagMetricsInterval))
	c.UI.Info(fmt.Sprintf("               Targets: %s", strings.Join(c.flagTargets, ", ")))
	c.UI.Info(fmt.Sprintf("                Output: %s", dstOutputFile))
	c.UI.Output("")

	// Release the log gate.
	c.logger.(hclog.OutputResettable).ResetOutputWithFlush(&hclog.LoggerOptions{
		Output: os.Stderr,
	}, gatedWriter)

	// Capture static information
	c.UI.Info("==> Capturing static information...")
	if err := c.captureStaticTargets(); err != nil {
		c.UI.Error(fmt.Sprintf("Error capturing static information: %s", err))
		return 2
	}

	c.UI.Output("")

	// Capture polling information
	c.UI.Info("==> Capturing dynamic information...")
	if err := c.capturePollingTargets(); err != nil {
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
	if err := ioutil.WriteFile(filepath.Join(c.flagOutput, "index.json"), bytes, 0o600); err != nil {
		return fmt.Errorf("error generating index file; %s", err)
	}

	return nil
}

// preflight performs various checks against the provided flags to ensure they
// are valid/reasonable values. It also takes care of instantiating a client and
// index object for use by the command.
func (c *DebugCommand) preflight(rawArgs []string) (string, error) {
	if !c.skipTimingChecks {
		// Guard duration and interval values to acceptable values
		if c.flagDuration < debugMinInterval {
			c.UI.Info(fmt.Sprintf("Overwriting duration value %q to the minimum value of %q", c.flagDuration, debugMinInterval))
			c.flagDuration = debugMinInterval
		}
		if c.flagInterval < debugMinInterval {
			c.UI.Info(fmt.Sprintf("Overwriting interval value %q to the minimum value of %q", c.flagInterval, debugMinInterval))
			c.flagInterval = debugMinInterval
		}
		if c.flagMetricsInterval < debugMinInterval {
			c.UI.Info(fmt.Sprintf("Overwriting metrics interval value %q to the minimum value of %q", c.flagMetricsInterval, debugMinInterval))
			c.flagMetricsInterval = debugMinInterval
		}
	}

	// These timing checks are always applicable since interval shouldn't be
	// greater than the duration
	if c.flagInterval > c.flagDuration {
		c.UI.Info(fmt.Sprintf("Overwriting interval value %q to the duration value %q", c.flagInterval, c.flagDuration))
		c.flagInterval = c.flagDuration
	}
	if c.flagMetricsInterval > c.flagDuration {
		c.UI.Info(fmt.Sprintf("Overwriting metrics interval value %q to the duration value %q", c.flagMetricsInterval, c.flagDuration))
		c.flagMetricsInterval = c.flagDuration
	}

	if len(c.flagTargets) == 0 {
		c.flagTargets = c.defaultTargets()
	} else {
		// Check for any invalid targets and ignore them if found
		invalidTargets := strutil.Difference(c.flagTargets, c.defaultTargets(), true)
		if len(invalidTargets) != 0 {
			c.UI.Info(fmt.Sprintf("Ignoring invalid targets: %s", strings.Join(invalidTargets, ", ")))
			c.flagTargets = strutil.Difference(c.flagTargets, invalidTargets, true)
		}
	}

	// Make sure we can talk to the server
	client, err := c.Client()
	if err != nil {
		return "", fmt.Errorf("unable to create client to connect to Vault: %s", err)
	}
	serverHealth, err := client.Sys().Health()
	if err != nil {
		return "", fmt.Errorf("unable to connect to the server: %s", err)
	}

	// Check if server is DR Secondary and we need to further
	// ignore any targets due to endpoint restrictions
	if serverHealth.ReplicationDRMode == "secondary" {
		invalidDRTargets := strutil.Difference(c.flagTargets, c.validDRSecondaryTargets(), true)
		if len(invalidDRTargets) != 0 {
			c.UI.Info(fmt.Sprintf("Ignoring invalid targets for DR Secondary: %s", strings.Join(invalidDRTargets, ", ")))
			c.flagTargets = strutil.Difference(c.flagTargets, invalidDRTargets, true)
		}
	}
	c.cachedClient = client

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
		_, err := os.Stat(dstOutputFile)
		switch {
		case os.IsNotExist(err):
			c.flagOutput = strings.TrimSuffix(c.flagOutput, ".tar.gz")
			c.flagOutput = strings.TrimSuffix(c.flagOutput, ".tgz")
		case err != nil:
			return "", fmt.Errorf("unable to stat file: %s", err)
		default:
			return "", fmt.Errorf("output file already exists: %s", dstOutputFile)
		}
	}

	// Stat check the directory to ensure we don't override any existing data.
	_, err = os.Stat(c.flagOutput)
	switch {
	case os.IsNotExist(err):
		err := os.MkdirAll(c.flagOutput, 0o700)
		if err != nil {
			return "", fmt.Errorf("unable to create output directory: %s", err)
		}
	case err != nil:
		return "", fmt.Errorf("unable to stat directory: %s", err)
	default:
		return "", fmt.Errorf("output directory already exists: %s", c.flagOutput)
	}

	// Populate initial index fields
	c.debugIndex = &debugIndex{
		VaultAddress:           client.Address(),
		ClientVersion:          version.GetVersion().VersionNumber(),
		ServerVersion:          serverHealth.Version,
		Compress:               c.flagCompress,
		DurationSeconds:        int(c.flagDuration.Seconds()),
		IntervalSeconds:        int(c.flagInterval.Seconds()),
		MetricsIntervalSeconds: int(c.flagMetricsInterval.Seconds()),
		RawArgs:                rawArgs,
		Version:                debugIndexVersion,
		Targets:                c.flagTargets,
		Timestamp:              captureTime,
		Errors:                 []*captureError{},
	}

	return dstOutputFile, nil
}

func (c *DebugCommand) defaultTargets() []string {
	return []string{"config", "host", "requests", "metrics", "pprof", "replication-status", "server-status", "log"}
}

func (c *DebugCommand) validDRSecondaryTargets() []string {
	return []string{"metrics", "replication-status", "server-status"}
}

func (c *DebugCommand) captureStaticTargets() error {
	// Capture configuration state
	if strutil.StrListContains(c.flagTargets, "config") {
		c.logger.Info("capturing configuration state")

		resp, err := c.cachedClient.Logical().Read("sys/config/state/sanitized")
		if err != nil {
			c.captureError("config", err)
			c.logger.Error("config: error capturing config state", "error", err)
			return nil
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
func (c *DebugCommand) capturePollingTargets() error {
	var g run.Group

	ctx, cancelFunc := context.WithTimeout(context.Background(), c.flagDuration+debugDurationGrace)
	defer cancelFunc()

	// This run group watches for interrupt or duration
	g.Add(func() error {
		for {
			select {
			case <-c.ShutdownCh:
				return nil
			case <-ctx.Done():
				return nil
			}
		}
	}, func(error) {})

	// Collect host-info if target is specified
	if strutil.StrListContains(c.flagTargets, "host") {
		g.Add(func() error {
			c.collectHostInfo(ctx)
			return nil
		}, func(error) {
			cancelFunc()
		})
	}

	// Collect metrics if target is specified
	if strutil.StrListContains(c.flagTargets, "metrics") {
		g.Add(func() error {
			c.collectMetrics(ctx)
			return nil
		}, func(error) {
			cancelFunc()
		})
	}

	// Collect pprof data if target is specified
	if strutil.StrListContains(c.flagTargets, "pprof") {
		g.Add(func() error {
			c.collectPprof(ctx)
			return nil
		}, func(error) {
			cancelFunc()
		})
	}

	// Collect replication status if target is specified
	if strutil.StrListContains(c.flagTargets, "replication-status") {
		g.Add(func() error {
			c.collectReplicationStatus(ctx)
			return nil
		}, func(error) {
			cancelFunc()
		})
	}

	// Collect server status if target is specified
	if strutil.StrListContains(c.flagTargets, "server-status") {
		g.Add(func() error {
			c.collectServerStatus(ctx)
			return nil
		}, func(error) {
			cancelFunc()
		})
	}

	// Collect in-flight request status if target is specified
	if strutil.StrListContains(c.flagTargets, "requests") {
		g.Add(func() error {
			c.collectInFlightRequestStatus(ctx)
			return nil
		}, func(error) {
			cancelFunc()
		})
	}

	if strutil.StrListContains(c.flagTargets, "log") {
		g.Add(func() error {
			c.writeLogs(ctx)
			// If writeLogs returned earlier due to an error, wait for context
			// to terminate so we don't abort everything.
			<-ctx.Done()
			return nil
		}, func(error) {
			cancelFunc()
		})
	}

	// We shouldn't bump across errors since none is returned by the interrupts,
	// but we error check for sanity here.
	if err := g.Run(); err != nil {
		return err
	}

	// Write collected data to their corresponding files
	if err := c.persistCollection(c.metricsCollection, "metrics.json"); err != nil {
		c.UI.Error(fmt.Sprintf("Error writing data to %s: %v", "metrics.json", err))
	}
	if err := c.persistCollection(c.serverStatusCollection, "server_status.json"); err != nil {
		c.UI.Error(fmt.Sprintf("Error writing data to %s: %v", "server_status.json", err))
	}
	if err := c.persistCollection(c.replicationStatusCollection, "replication_status.json"); err != nil {
		c.UI.Error(fmt.Sprintf("Error writing data to %s: %v", "replication_status.json", err))
	}
	if err := c.persistCollection(c.hostInfoCollection, "host_info.json"); err != nil {
		c.UI.Error(fmt.Sprintf("Error writing data to %s: %v", "host_info.json", err))
	}
	if err := c.persistCollection(c.inFlightReqStatusCollection, "requests.json"); err != nil {
		c.UI.Error(fmt.Sprintf("Error writing data to %s: %v", "requests.json", err))
	}
	return nil
}

func (c *DebugCommand) collectHostInfo(ctx context.Context) {
	idxCount := 0
	intervalTicker := time.Tick(c.flagInterval)

	for {
		if idxCount > 0 {
			select {
			case <-ctx.Done():
				return
			case <-intervalTicker:
			}
		}

		c.logger.Info("capturing host information", "count", idxCount)
		idxCount++

		r := c.cachedClient.NewRequest("GET", "/v1/sys/host-info")
		resp, err := c.cachedClient.RawRequestWithContext(ctx, r)
		if err != nil {
			c.captureError("host", err)
			return
		}
		if resp != nil {
			defer resp.Body.Close()

			secret, err := api.ParseSecret(resp.Body)
			if err != nil {
				c.captureError("host", err)
				return
			}
			if secret != nil && secret.Data != nil {
				hostEntry := secret.Data
				c.hostInfoCollection = append(c.hostInfoCollection, hostEntry)
			}
		}
	}
}

func (c *DebugCommand) collectMetrics(ctx context.Context) {
	idxCount := 0
	intervalTicker := time.Tick(c.flagMetricsInterval)

	for {
		if idxCount > 0 {
			select {
			case <-ctx.Done():
				return
			case <-intervalTicker:
			}
		}

		c.logger.Info("capturing metrics", "count", idxCount)
		idxCount++

		// Perform metrics request
		r := c.cachedClient.NewRequest("GET", "/v1/sys/metrics")
		resp, err := c.cachedClient.RawRequestWithContext(ctx, r)
		if err != nil {
			c.captureError("metrics", err)
			continue
		}
		if resp != nil {
			defer resp.Body.Close()

			metricsEntry := make(map[string]interface{})
			err := json.NewDecoder(resp.Body).Decode(&metricsEntry)
			if err != nil {
				c.captureError("metrics", err)
				continue
			}
			c.metricsCollection = append(c.metricsCollection, metricsEntry)
		}
	}
}

func (c *DebugCommand) collectPprof(ctx context.Context) {
	idxCount := 0
	startTime := time.Now()
	intervalTicker := time.Tick(c.flagInterval)

	for {
		if idxCount > 0 {
			select {
			case <-ctx.Done():
				return
			case <-intervalTicker:
			}
		}

		currentTimestamp := time.Now().UTC()
		c.logger.Info("capturing pprof data", "count", idxCount)
		idxCount++

		// Create a sub-directory for pprof data
		currentDir := currentTimestamp.Format(fileFriendlyTimeFormat)
		dirName := filepath.Join(c.flagOutput, currentDir)
		if err := os.MkdirAll(dirName, 0o700); err != nil {
			c.UI.Error(fmt.Sprintf("Error creating sub-directory for time interval: %s", err))
			continue
		}

		var wg sync.WaitGroup

		for _, target := range []string{"threadcreate", "allocs", "block", "mutex", "goroutine", "heap"} {
			wg.Add(1)
			go func(target string) {
				defer wg.Done()
				data, err := pprofTarget(ctx, c.cachedClient, target, nil)
				if err != nil {
					c.captureError("pprof."+target, err)
					return
				}

				err = ioutil.WriteFile(filepath.Join(dirName, target+".prof"), data, 0o600)
				if err != nil {
					c.captureError("pprof."+target, err)
				}
			}(target)
		}

		// As a convenience, we'll also fetch the goroutine target using debug=2, which yields a text
		// version of the stack traces that don't require using `go tool pprof` to view.
		wg.Add(1)
		go func() {
			defer wg.Done()
			data, err := pprofTarget(ctx, c.cachedClient, "goroutine", url.Values{"debug": []string{"2"}})
			if err != nil {
				c.captureError("pprof.goroutines-text", err)
				return
			}

			err = ioutil.WriteFile(filepath.Join(dirName, "goroutines.txt"), data, 0o600)
			if err != nil {
				c.captureError("pprof.goroutines-text", err)
			}
		}()

		// If the our remaining duration is less than the interval value
		// skip profile and trace.
		runDuration := currentTimestamp.Sub(startTime)
		if (c.flagDuration+debugDurationGrace)-runDuration < c.flagInterval {
			wg.Wait()
			continue
		}

		// Capture profile
		wg.Add(1)
		go func() {
			defer wg.Done()
			data, err := pprofProfile(ctx, c.cachedClient, c.flagInterval)
			if err != nil {
				c.captureError("pprof.profile", err)
				return
			}

			err = ioutil.WriteFile(filepath.Join(dirName, "profile.prof"), data, 0o600)
			if err != nil {
				c.captureError("pprof.profile", err)
			}
		}()

		// Capture trace
		wg.Add(1)
		go func() {
			defer wg.Done()
			data, err := pprofTrace(ctx, c.cachedClient, c.flagInterval)
			if err != nil {
				c.captureError("pprof.trace", err)
				return
			}

			err = ioutil.WriteFile(filepath.Join(dirName, "trace.out"), data, 0o600)
			if err != nil {
				c.captureError("pprof.trace", err)
			}
		}()

		wg.Wait()
	}
}

func (c *DebugCommand) collectReplicationStatus(ctx context.Context) {
	idxCount := 0
	intervalTicker := time.Tick(c.flagInterval)

	for {
		if idxCount > 0 {
			select {
			case <-ctx.Done():
				return
			case <-intervalTicker:
			}
		}

		c.logger.Info("capturing replication status", "count", idxCount)
		idxCount++

		r := c.cachedClient.NewRequest("GET", "/v1/sys/replication/status")
		resp, err := c.cachedClient.RawRequestWithContext(ctx, r)
		if err != nil {
			c.captureError("replication-status", err)
			return
		}
		if resp != nil {
			defer resp.Body.Close()

			secret, err := api.ParseSecret(resp.Body)
			if err != nil {
				c.captureError("replication-status", err)
				return
			}
			if secret != nil && secret.Data != nil {
				replicationEntry := secret.Data
				replicationEntry["timestamp"] = time.Now().UTC()
				c.replicationStatusCollection = append(c.replicationStatusCollection, replicationEntry)
			}
		}
	}
}

func (c *DebugCommand) collectServerStatus(ctx context.Context) {
	idxCount := 0
	intervalTicker := time.Tick(c.flagInterval)

	for {
		if idxCount > 0 {
			select {
			case <-ctx.Done():
				return
			case <-intervalTicker:
			}
		}

		c.logger.Info("capturing server status", "count", idxCount)
		idxCount++

		healthInfo, err := c.cachedClient.Sys().Health()
		if err != nil {
			c.captureError("server-status.health", err)
		}
		sealInfo, err := c.cachedClient.Sys().SealStatus()
		if err != nil {
			c.captureError("server-status.seal", err)
		}

		statusEntry := map[string]interface{}{
			"timestamp": time.Now().UTC(),
			"health":    healthInfo,
			"seal":      sealInfo,
		}
		c.serverStatusCollection = append(c.serverStatusCollection, statusEntry)
	}
}

func (c *DebugCommand) collectInFlightRequestStatus(ctx context.Context) {
	idxCount := 0
	intervalTicker := time.Tick(c.flagInterval)

	for {
		if idxCount > 0 {
			select {
			case <-ctx.Done():
				return
			case <-intervalTicker:
			}
		}

		c.logger.Info("capturing in-flight request status", "count", idxCount)
		idxCount++

		req := c.cachedClient.NewRequest("GET", "/v1/sys/in-flight-req")
		resp, err := c.cachedClient.RawRequestWithContext(ctx, req)
		if err != nil {
			c.captureError("requests", err)
			return
		}

		var data map[string]interface{}
		if resp != nil {
			defer resp.Body.Close()
			err = jsonutil.DecodeJSONFromReader(resp.Body, &data)
			if err != nil {
				c.captureError("requests", err)
				return
			}

			statusEntry := map[string]interface{}{
				"timestamp":          time.Now().UTC(),
				"in_flight_requests": data,
			}
			c.inFlightReqStatusCollection = append(c.inFlightReqStatusCollection, statusEntry)
		}
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
	if err := ioutil.WriteFile(filepath.Join(c.flagOutput, outFile), bytes, 0o600); err != nil {
		return err
	}

	return nil
}

func (c *DebugCommand) compress(dst string) error {
	if runtime.GOOS != "windows" {
		defer osutil.Umask(osutil.Umask(0o077))
	}

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

func pprofTarget(ctx context.Context, client *api.Client, target string, params url.Values) ([]byte, error) {
	req := client.NewRequest("GET", "/v1/sys/pprof/"+target)
	if params != nil {
		req.Params = params
	}
	resp, err := client.RawRequestWithContext(ctx, req)
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

func pprofProfile(ctx context.Context, client *api.Client, duration time.Duration) ([]byte, error) {
	seconds := int(duration.Seconds())
	secStr := strconv.Itoa(seconds)

	req := client.NewRequest("GET", "/v1/sys/pprof/profile")
	req.Params.Add("seconds", secStr)
	resp, err := client.RawRequestWithContext(ctx, req)
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

func pprofTrace(ctx context.Context, client *api.Client, duration time.Duration) ([]byte, error) {
	seconds := int(duration.Seconds())
	secStr := strconv.Itoa(seconds)

	req := client.NewRequest("GET", "/v1/sys/pprof/trace")
	req.Params.Add("seconds", secStr)
	resp, err := client.RawRequestWithContext(ctx, req)
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

// newCaptureError instantiates a new captureError.
func (c *DebugCommand) captureError(target string, err error) {
	c.errLock.Lock()
	c.debugIndex.Errors = append(c.debugIndex.Errors, &captureError{
		TargetError: err.Error(),
		Target:      target,
		Timestamp:   time.Now().UTC(),
	})
	c.errLock.Unlock()
}

func (c *DebugCommand) writeLogs(ctx context.Context) {
	out, err := os.OpenFile(filepath.Join(c.flagOutput, "vault.log"), os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		c.captureError("log", err)
		return
	}
	defer out.Close()

	logCh, err := c.cachedClient.Sys().Monitor(ctx, "trace", c.logFormat)
	if err != nil {
		c.captureError("log", err)
		return
	}

	for {
		select {
		case log := <-logCh:
			if !strings.HasSuffix(log, "\n") {
				log += "\n"
			}
			_, err = out.WriteString(log)
			if err != nil {
				c.captureError("log", err)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
