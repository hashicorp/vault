package command

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/logutils"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/helper/flag-slice"
	"github.com/hashicorp/vault/helper/gated-writer"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/vault"
)

// ServerCommand is a Command that starts the Vault server.
type ServerCommand struct {
	AuditBackends      map[string]audit.Factory
	CredentialBackends map[string]logical.Factory
	LogicalBackends    map[string]logical.Factory

	Meta
}

func (c *ServerCommand) Run(args []string) int {
	var dev bool
	var configPath []string
	var logLevel string
	flags := c.Meta.FlagSet("server", FlagSetDefault)
	flags.BoolVar(&dev, "dev", false, "")
	flags.StringVar(&logLevel, "log-level", "info", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	flags.Var((*sliceflag.StringFlag)(&configPath), "config", "config")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validation
	if !dev && len(configPath) == 0 {
		c.Ui.Error("At least one config path must be specified with -config")
		flags.Usage()
		return 1
	}

	// Load the configuration
	var config *server.Config
	if dev {
		config = server.DevConfig()
	}
	for _, path := range configPath {
		current, err := server.LoadConfig(path)
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error loading configuration from %s: %s", path, err))
			return 1
		}

		if config == nil {
			config = current
		} else {
			config = config.Merge(current)
		}
	}

	// Create a logger. We wrap it in a gated writer so that it doesn't
	// start logging too early.
	logGate := &gatedwriter.Writer{Writer: os.Stderr}
	logger := log.New(&logutils.LevelFilter{
		Levels: []logutils.LogLevel{
			"TRACE", "DEBUG", "INFO", "WARN", "ERR"},
		MinLevel: logutils.LogLevel(strings.ToUpper(logLevel)),
		Writer:   logGate,
	}, "", log.LstdFlags)

	// Initialize the backend
	backend, err := physical.NewBackend(
		config.Backend.Type, config.Backend.Config)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing backend of type %s: %s",
			config.Backend.Type, err))
		return 1
	}

	// Initialize the core
	core, err := vault.NewCore(&vault.CoreConfig{
		AdvertiseAddr:      config.Backend.AdvertiseAddr,
		Physical:           backend,
		AuditBackends:      c.AuditBackends,
		CredentialBackends: c.CredentialBackends,
		LogicalBackends:    c.LogicalBackends,
		Logger:             logger,
	})

	// If we're in dev mode, then initialize the core
	if dev {
		init, err := c.enableDev(core)
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error initializing dev mode: %s", err))
			return 1
		}

		c.Ui.Output(fmt.Sprintf(
			"WARNING: Dev mode is enabled!\n\n"+
				"In this mode, Vault is completely in-memory and unsealed.\n"+
				"Vault is configured to only have a single unseal key. The root\n"+
				"token has already been authenticated with the CLI, so you can\n"+
				"immediately begin using the Vault CLI.\n\n"+
				"The only step you need to take is to set the following\n"+
				"environment variable since Vault will be taking without TLS:\n\n"+
				"    export VAULT_ADDR='http://127.0.0.1:8200'\n\n"+
				"The unseal key and root token are reproduced below in case you\n"+
				"want to seal/unseal the Vault or play with authentication.\n\n"+
				"Unseal Key: %s\nRoot Token: %s\n",
			hex.EncodeToString(init.SecretShares[0]),
			init.RootToken,
		))
	}

	// Compile server information for output later
	infoKeys := make([]string, 0, 10)
	info := make(map[string]string)
	info["backend"] = config.Backend.Type
	info["log level"] = logLevel
	infoKeys = append(infoKeys, "log level", "backend")

	// If the backend supports HA, then note it
	if _, ok := backend.(physical.HABackend); ok {
		info["backend"] += " (HA available)"
	}

	// Initialize the telemetry
	if err := c.setupTelementry(config); err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing telemetry: %s", err))
		return 1
	}

	// Initialize the listeners
	lns := make([]net.Listener, 0, len(config.Listeners))
	for i, lnConfig := range config.Listeners {
		ln, props, err := server.NewListener(lnConfig.Type, lnConfig.Config)
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error initializing listener of type %s: %s",
				lnConfig.Type, err))
			return 1
		}

		// Store the listener props for output later
		key := fmt.Sprintf("listener %d", i+1)
		propsList := make([]string, 0, len(props))
		for k, v := range props {
			propsList = append(propsList, fmt.Sprintf(
				"%s: %q", k, v))
		}
		sort.Strings(propsList)
		infoKeys = append(infoKeys, key)
		info[key] = fmt.Sprintf(
			"%s (%s)", lnConfig.Type, strings.Join(propsList, ", "))

		lns = append(lns, ln)
	}

	// Initialize the HTTP server
	server := &http.Server{}
	server.Handler = vaulthttp.Handler(core)
	for _, ln := range lns {
		go server.Serve(ln)
	}

	// Server configuration output
	padding := 18
	c.Ui.Output("==> Vault server configuration:\n")
	for _, k := range infoKeys {
		c.Ui.Output(fmt.Sprintf(
			"%s%s: %s",
			strings.Repeat(" ", padding-len(k)),
			strings.Title(k),
			info[k]))
	}
	c.Ui.Output("")

	// Output the header that the server has started
	c.Ui.Output("==> Vault server started! Log data will stream in below:\n")

	// Release the log gate.
	logGate.Flush()

	<-make(chan struct{})
	return 0
}

func (c *ServerCommand) enableDev(core *vault.Core) (*vault.InitResult, error) {
	// Initialize it with a basic single key
	init, err := core.Initialize(&vault.SealConfig{
		SecretShares:    1,
		SecretThreshold: 1,
	})
	if err != nil {
		return nil, err
	}

	// Copy the key so that it can be zeroed
	key := make([]byte, len(init.SecretShares[0]))
	copy(key, init.SecretShares[0])

	// Unseal the core
	unsealed, err := core.Unseal(key)
	if err != nil {
		return nil, err
	}
	if !unsealed {
		return nil, fmt.Errorf("failed to unseal Vault for dev mode")
	}

	// Set the token
	tokenHelper, err := c.TokenHelper()
	if err != nil {
		return nil, err
	}
	if err := tokenHelper.Store(init.RootToken); err != nil {
		return nil, err
	}

	return init, nil
}

// setupTelementry is used ot setup the telemetry sub-systems
func (c *ServerCommand) setupTelementry(config *server.Config) error {
	/* Setup telemetry
	Aggregate on 10 second intervals for 1 minute. Expose the
	metrics over stderr when there is a SIGUSR1 received.
	*/
	inm := metrics.NewInmemSink(10*time.Second, time.Minute)
	metrics.DefaultInmemSignal(inm)
	metricsConf := metrics.DefaultConfig("vault")

	// Configure the statsite sink
	var fanout metrics.FanoutSink
	if config.StatsiteAddr != "" {
		sink, err := metrics.NewStatsiteSink(config.StatsiteAddr)
		if err != nil {
			return err
		}
		fanout = append(fanout, sink)
	}

	// Configure the statsd sink
	if config.StatsdAddr != "" {
		sink, err := metrics.NewStatsdSink(config.StatsdAddr)
		if err != nil {
			return err
		}
		fanout = append(fanout, sink)
	}

	// Initialize the global sink
	if len(fanout) > 0 {
		fanout = append(fanout, inm)
		metrics.NewGlobal(metricsConf, fanout)
	} else {
		metricsConf.EnableHostname = false
		metrics.NewGlobal(metricsConf, inm)
	}
	return nil
}

func (c *ServerCommand) Synopsis() string {
	return "Start a Vault server"
}

func (c *ServerCommand) Help() string {
	helpText := `
Usage: vault server [options]

  Start a Vault server.

  This command starts a Vault server that responds to API requests.
  Vault will start in a "sealed" state. The Vault must be unsealed
  with "vault unseal" or the API before this server can respond to requests.
  This must be done for every server.

  If the server is being started against a storage backend that has
  brand new (no existing Vault data in it), it must be initialized with
  "vault init" or the API first.


General Options:

  -config=<path>      Path to the configuration file or directory. This can be
                      specified multiple times. If it is a directory, all
                      files with a ".hcl" or ".json" suffix will be loaded.

  -log-level=info     Log verbosity. Defaults to "info", will be outputted
                      to stderr.

`
	return strings.TrimSpace(helpText)
}
