package command

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/gatedwriter"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/version"
	"github.com/mholt/archiver"
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
	flagAddresses       []string
	flagCluster         bool

	logger hclog.Logger

	// ShutdownCh is used to capture interrupt signal and end polling capture
	ShutdownCh chan struct{}

	// skipTimingChecks bypasses timing-related checks, used primarily for tests
	skipTimingChecks bool
}

type collector struct {
	client          *api.Client
	targets         []string
	outputDir       string
	duration        time.Duration
	interval        time.Duration
	metricsInterval time.Duration

	// debugIndex is used to keep track of the index state, which gets written
	// to a file at the end.
	debugIndex *debugIndex

	logger hclog.Logger

	// shutdownCh is used to capture interrupt signal and end polling capture
	shutdownCh chan struct{}

	// Collection slices to hold data
	hostInfoCollection          []map[string]interface{}
	metricsCollection           []map[string]interface{}
	replicationStatusCollection []map[string]interface{}
	serverStatusCollection      []map[string]interface{}

	errLock sync.Mutex

	uiError func(string)
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

	f.StringSliceVar(&StringSliceVar{
		Name:   "addresses",
		Target: &c.flagAddresses,
		Usage: "Vault node addresses to query. " +
			"This can be specified multiple times to capture multiple node's debug data. " +
			"For a single node just use the usual CLI mechanisms like VAULT_ADDRESS, -address",
	})

	f.BoolVar(&BoolVar{
		Name:   "cluster",
		Target: &c.flagCluster,
		Usage:  "When true, all nodes in the cluster will be queried based on the active node's sys/ha-status response",
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

	// Initialize the logger for debug output.  We want to delay emitting anything
	// on this log that writes to stderr until after we've finished with the
	// c.UI calls announcing what we'll be doing.
	gatedWriter := gatedwriter.NewWriter(os.Stderr)
	c.logger = logging.NewVaultLoggerWithWriter(gatedWriter, hclog.Trace)

	dstOutputFile, baseDebugIndex, err := c.preflight(args)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error during validation: %s", err))
		return 1
	}

	genericClient, err := c.Client()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating client: %s", err))
		return 1
	}

	var collectors []*collector
	switch {
	case len(c.flagAddresses) > 0 && c.flagCluster:
		c.UI.Error("Cannot specify both -cluster and -addresses")
		return 1
	case c.flagCluster:
		resp, err := genericClient.Sys().HAStatus()
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error reading ha-status: %s", err))
			return 1
		}
		for _, n := range resp.Nodes {
			c.flagAddresses = append(c.flagAddresses, n.APIAddress)
		}
	case len(c.flagAddresses) == 0:
		c.flagAddresses = append(c.flagAddresses, genericClient.Address())
	}

	for _, addr := range c.flagAddresses {
		debugIndex := *baseDebugIndex
		debugIndex.VaultAddress = addr
		client, err := api.NewClient(genericClient.CloneConfig())
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error creating client for addr %s: %s", addr, err))
			return 1
		}
		err = client.SetAddress(addr)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error creating client for addr %s: %s", addr, err))
			return 1
		}
		client.SetToken(genericClient.Token())

		if _, err := client.Sys().Health(); err != nil {
			c.UI.Error(fmt.Sprintf("Unable to connect to addr %s: %s", addr, err))
			return 1
		}

		u, _ := url.Parse(addr)
		shortAddr := u.Host
		outputDir := filepath.Join(c.flagOutput, shortAddr)
		os.MkdirAll(outputDir, 0o700)
		collectors = append(collectors, &collector{
			client:          client,
			targets:         c.flagTargets,
			outputDir:       filepath.Join(c.flagOutput, shortAddr),
			duration:        c.flagDuration,
			interval:        c.flagInterval,
			metricsInterval: c.flagMetricsInterval,
			debugIndex:      &debugIndex,
			logger:          c.logger.Named(shortAddr),
			uiError:         func(s string) { c.UI.Error(s) },
			shutdownCh:      c.ShutdownCh,
		})
	}

	// Print debug information
	c.UI.Output("==> Starting debug capture...")
	c.UI.Info(fmt.Sprintf("       Vault Addresses: %s", c.flagAddresses))
	c.UI.Info(fmt.Sprintf("        Client Version: %s", version.GetVersion().VersionNumber()))
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

	var wg sync.WaitGroup
	errs := make([]error, len(collectors))
	for i := range collectors {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			errs[i] = collectors[i].captureStaticTargets()
		}(i)
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error capturing static information for address %s: %s", c.flagAddresses[i], err))
			return 2
		}
	}
	c.UI.Output("")

	// Capture polling information
	c.UI.Info("==> Capturing dynamic information...")

	for i := range collectors {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			errs[i] = collectors[i].capturePollingTargets()
		}(i)
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error capturing dynamic information for address %s: %s", c.flagAddresses[i], err))
			return 2
		}
	}
	c.UI.Output("")

	c.UI.Info("Finished capturing information, bundling files...")

	for i := range collectors {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			errs[i] = collectors[i].generateIndex()
		}(i)
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error generating index for address %s: %s", c.flagAddresses[i], err))
		}
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

func (c *collector) generateIndex() error {
	outputLayout := map[string]interface{}{
		"files": []string{},
	}
	// Walk the directory to generate the output layout
	err := filepath.Walk(c.outputDir, func(path string, info os.FileInfo, err error) error {
		// Prevent panic by handling failure accessing a path
		if err != nil {
			return err
		}

		// Skip the base dir
		if path == c.outputDir {
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

		relPath, err := filepath.Rel(c.outputDir, path)
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
	if err := ioutil.WriteFile(filepath.Join(c.outputDir, "index.json"), bytes, 0o644); err != nil {
		return fmt.Errorf("error generating index file; %s", err)
	}

	return nil
}

// preflight performs various checks against the provided flags to ensure they
// are valid/reasonable values.  It returns the output file & the base debugIndex,
// or an error.
func (c *DebugCommand) preflight(rawArgs []string) (string, *debugIndex, error) {
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
			return "", nil, fmt.Errorf("unable to stat file: %s", err)
		default:
			return "", nil, fmt.Errorf("output file already exists: %s", dstOutputFile)
		}
	}

	// Stat check the directory to ensure we don't override any existing data.
	_, err := os.Stat(c.flagOutput)
	switch {
	case os.IsNotExist(err):
		err := os.MkdirAll(c.flagOutput, 0o755)
		if err != nil {
			return "", nil, fmt.Errorf("unable to create output directory: %s", err)
		}
	case err != nil:
		return "", nil, fmt.Errorf("unable to stat directory: %s", err)
	default:
		return "", nil, fmt.Errorf("output directory already exists: %s", c.flagOutput)
	}

	// Populate initial index fields
	debugIndex := &debugIndex{
		ClientVersion:          version.GetVersion().VersionNumber(),
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

	return dstOutputFile, debugIndex, nil
}

func (c *DebugCommand) defaultTargets() []string {
	return []string{"config", "host", "metrics", "pprof", "replication-status", "server-status", "log"}
}

func (c *collector) captureStaticTargets() error {
	// Capture configuration state
	if strutil.StrListContains(c.targets, "config") {
		c.logger.Info("capturing configuration state")

		resp, err := c.client.Logical().Read("sys/config/state/sanitized")
		if err != nil {
			c.captureError("config", err)
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
				c.uiError(fmt.Sprintf("Error writing data to %s: %v", "config.json", err))
			}
		}
	}

	return nil
}

// capturePollingTargets captures all dynamic targets over the specified
// duration and interval.
func (c *collector) capturePollingTargets() error {
	var g run.Group

	ctx, cancelFunc := context.WithTimeout(context.Background(), c.duration+debugDurationGrace)
	defer cancelFunc()

	// This run group watches for interrupt or duration
	g.Add(func() error {
		for {
			select {
			case <-c.shutdownCh:
				return nil
			case <-ctx.Done():
				return nil
			}
		}
	}, func(error) {})

	// Collect host-info if target is specified
	if strutil.StrListContains(c.targets, "host") {
		g.Add(func() error {
			c.collectHostInfo(ctx)
			return nil
		}, func(error) {
			cancelFunc()
		})
	}

	// Collect metrics if target is specified
	if strutil.StrListContains(c.targets, "metrics") {
		g.Add(func() error {
			c.collectMetrics(ctx)
			return nil
		}, func(error) {
			cancelFunc()
		})
	}

	// Collect pprof data if target is specified
	if strutil.StrListContains(c.targets, "pprof") {
		g.Add(func() error {
			c.collectPprof(ctx)
			return nil
		}, func(error) {
			cancelFunc()
		})
	}

	// Collect replication status if target is specified
	if strutil.StrListContains(c.targets, "replication-status") {
		g.Add(func() error {
			c.collectReplicationStatus(ctx)
			return nil
		}, func(error) {
			cancelFunc()
		})
	}

	// Collect server status if target is specified
	if strutil.StrListContains(c.targets, "server-status") {
		g.Add(func() error {
			c.collectServerStatus(ctx)
			return nil
		}, func(error) {
			cancelFunc()
		})
	}

	if strutil.StrListContains(c.targets, "log") {
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
		c.uiError(fmt.Sprintf("Error writing data to %s: %v", "metrics.json", err))
	}
	if err := c.persistCollection(c.serverStatusCollection, "server_status.json"); err != nil {
		c.uiError(fmt.Sprintf("Error writing data to %s: %v", "server_status.json", err))
	}
	if err := c.persistCollection(c.replicationStatusCollection, "replication_status.json"); err != nil {
		c.uiError(fmt.Sprintf("Error writing data to %s: %v", "replication_status.json", err))
	}
	if err := c.persistCollection(c.hostInfoCollection, "host_info.json"); err != nil {
		c.uiError(fmt.Sprintf("Error writing data to %s: %v", "host_info.json", err))
	}

	return nil
}

func (c *collector) collectHostInfo(ctx context.Context) {
	idxCount := 0
	intervalTicker := time.Tick(c.interval)

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

		r := c.client.NewRequest("GET", "/v1/sys/host-info")
		resp, err := c.client.RawRequestWithContext(ctx, r)
		if err != nil {
			c.captureError("host", err)
		}
		if resp != nil {
			defer resp.Body.Close()

			secret, err := api.ParseSecret(resp.Body)
			if err != nil {
				c.captureError("host", err)
			}
			if secret != nil && secret.Data != nil {
				hostEntry := secret.Data
				c.hostInfoCollection = append(c.hostInfoCollection, hostEntry)
			}
		}
	}
}

func (c *collector) collectMetrics(ctx context.Context) {
	idxCount := 0
	intervalTicker := time.Tick(c.metricsInterval)

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

		healthStatus, err := c.client.Sys().Health()
		if err != nil {
			c.captureError("metrics", err)
			continue
		}

		// Check replication status. We skip on processing metrics if we're one
		// a DR node, though non-perf standbys will fail if they aren't using
		// unauthenticated_metrics_access.
		switch {
		case healthStatus.ReplicationDRMode == "secondary":
			c.logger.Info("skipping metrics capture on DR secondary node")
			continue
		}

		// Perform metrics request
		r := c.client.NewRequest("GET", "/v1/sys/metrics")
		resp, err := c.client.RawRequestWithContext(ctx, r)
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

func (c *collector) collectPprof(ctx context.Context) {
	idxCount := 0
	startTime := time.Now()
	intervalTicker := time.Tick(c.interval)

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
		dirName := filepath.Join(c.outputDir, currentDir)
		if err := os.MkdirAll(dirName, 0o755); err != nil {
			c.uiError(fmt.Sprintf("Error creating sub-directory for time interval: %s", err))
			continue
		}

		var wg sync.WaitGroup

		for _, target := range []string{"threadcreate", "allocs", "block", "mutex", "goroutine", "heap"} {
			wg.Add(1)
			go func(target string) {
				defer wg.Done()
				data, err := pprofTarget(ctx, c.client, target, nil)
				if err != nil {
					c.captureError("pprof."+target, err)
					return
				}

				filename := filepath.Join(dirName, target+".prof")
				err = ioutil.WriteFile(filename, data, 0o644)
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
			data, err := pprofTarget(ctx, c.client, "goroutine", url.Values{"debug": []string{"2"}})
			if err != nil {
				c.captureError("pprof.goroutines-text", err)
				return
			}

			err = ioutil.WriteFile(filepath.Join(dirName, "goroutines.txt"), data, 0o644)
			if err != nil {
				c.captureError("pprof.goroutines-text", err)
			}
		}()

		// If the our remaining duration is less than the interval value
		// skip profile and trace.
		runDuration := currentTimestamp.Sub(startTime)
		if (c.duration+debugDurationGrace)-runDuration < c.interval {
			wg.Wait()
			continue
		}

		// Capture profile
		wg.Add(1)
		go func() {
			defer wg.Done()
			data, err := pprofProfile(ctx, c.client, c.interval)
			if err != nil {
				c.captureError("pprof.profile", err)
				return
			}

			err = ioutil.WriteFile(filepath.Join(dirName, "profile.prof"), data, 0o644)
			if err != nil {
				c.captureError("pprof.profile", err)
			}
		}()

		// Capture trace
		wg.Add(1)
		go func() {
			defer wg.Done()
			data, err := pprofTrace(ctx, c.client, c.interval)
			if err != nil {
				c.captureError("pprof.trace", err)
				return
			}

			err = ioutil.WriteFile(filepath.Join(dirName, "trace.out"), data, 0o644)
			if err != nil {
				c.captureError("pprof.trace", err)
			}
		}()

		wg.Wait()
	}
}

func (c *collector) collectReplicationStatus(ctx context.Context) {
	idxCount := 0
	intervalTicker := time.Tick(c.interval)

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

		r := c.client.NewRequest("GET", "/v1/sys/replication/status")
		resp, err := c.client.RawRequestWithContext(ctx, r)
		if err != nil {
			c.captureError("replication-status", err)
		}
		if resp != nil {
			defer resp.Body.Close()

			secret, err := api.ParseSecret(resp.Body)
			if err != nil {
				c.captureError("replication-status", err)
			}
			if secret != nil && secret.Data != nil {
				replicationEntry := secret.Data
				replicationEntry["timestamp"] = time.Now().UTC()
				c.replicationStatusCollection = append(c.replicationStatusCollection, replicationEntry)
			}
		}
	}
}

func (c *collector) collectServerStatus(ctx context.Context) {
	idxCount := 0
	intervalTicker := time.Tick(c.interval)

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

		healthInfo, err := c.client.Sys().Health()
		if err != nil {
			c.captureError("server-status.health", err)
		}
		sealInfo, err := c.client.Sys().SealStatus()
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

// persistCollection writes the collected data for a particular target onto the
// specified file. If the collection is empty, it returns immediately.
func (c *collector) persistCollection(collection []map[string]interface{}, outFile string) error {
	if len(collection) == 0 {
		return nil
	}

	// Write server-status file and update the index
	bytes, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(c.outputDir, outFile), bytes, 0o644); err != nil {
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

func (c *collector) captureError(target string, err error) {
	c.errLock.Lock()
	c.debugIndex.Errors = append(c.debugIndex.Errors, &captureError{
		TargetError: err.Error(),
		Target:      target,
		Timestamp:   time.Now().UTC(),
	})
	c.errLock.Unlock()
}

func (c *collector) writeLogs(ctx context.Context) {
	out, err := os.Create(filepath.Join(c.outputDir, "vault.log"))
	if err != nil {
		c.captureError("log", err)
		return
	}
	defer out.Close()

	logCh, err := c.client.Sys().Monitor(ctx, "trace")
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
