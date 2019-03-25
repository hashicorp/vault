package command

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/auth"
	"github.com/hashicorp/vault/command/agent/auth/alicloud"
	"github.com/hashicorp/vault/command/agent/auth/approle"
	"github.com/hashicorp/vault/command/agent/auth/aws"
	"github.com/hashicorp/vault/command/agent/auth/azure"
	"github.com/hashicorp/vault/command/agent/auth/gcp"
	"github.com/hashicorp/vault/command/agent/auth/jwt"
	"github.com/hashicorp/vault/command/agent/auth/kubernetes"
	"github.com/hashicorp/vault/command/agent/cache"
	"github.com/hashicorp/vault/command/agent/config"
	"github.com/hashicorp/vault/command/agent/sink"
	"github.com/hashicorp/vault/command/agent/sink/file"
	"github.com/hashicorp/vault/command/agent/sink/inmem"
	"github.com/hashicorp/vault/helper/consts"
	gatedwriter "github.com/hashicorp/vault/helper/gated-writer"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/version"
	"github.com/kr/pretty"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*AgentCommand)(nil)
var _ cli.CommandAutocomplete = (*AgentCommand)(nil)

type AgentCommand struct {
	*BaseCommand

	ShutdownCh chan struct{}
	SighupCh   chan struct{}

	logWriter io.Writer
	logGate   *gatedwriter.Writer
	logger    log.Logger

	cleanupGuard sync.Once

	startedCh chan (struct{}) // for tests

	flagConfigs  []string
	flagLogLevel string

	flagTestVerifyOnly bool
	flagCombineLogs    bool
}

func (c *AgentCommand) Synopsis() string {
	return "Start a Vault agent"
}

func (c *AgentCommand) Help() string {
	helpText := `
Usage: vault agent [options]

  This command starts a Vault agent that can perform automatic authentication
  in certain environments.

  Start an agent with a configuration file:

      $ vault agent -config=/etc/vault/config.hcl

  For a full list of examples, please see the documentation.

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *AgentCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.StringSliceVar(&StringSliceVar{
		Name:   "config",
		Target: &c.flagConfigs,
		Completion: complete.PredictOr(
			complete.PredictFiles("*.hcl"),
			complete.PredictFiles("*.json"),
		),
		Usage: "Path to a configuration file. This configuration file should " +
			"contain only agent directives.",
	})

	f.StringVar(&StringVar{
		Name:       "log-level",
		Target:     &c.flagLogLevel,
		Default:    "info",
		EnvVar:     "VAULT_LOG_LEVEL",
		Completion: complete.PredictSet("trace", "debug", "info", "warn", "err"),
		Usage: "Log verbosity level. Supported values (in order of detail) are " +
			"\"trace\", \"debug\", \"info\", \"warn\", and \"err\".",
	})

	// Internal-only flags to follow.
	//
	// Why hello there little source code reader! Welcome to the Vault source
	// code. The remaining options are intentionally undocumented and come with
	// no warranty or backwards-compatibility promise. Do not use these flags
	// in production. Do not build automation using these flags. Unless you are
	// developing against Vault, you should not need any of these flags.

	// TODO: should the below flags be public?
	f.BoolVar(&BoolVar{
		Name:    "combine-logs",
		Target:  &c.flagCombineLogs,
		Default: false,
		Hidden:  true,
	})

	f.BoolVar(&BoolVar{
		Name:    "test-verify-only",
		Target:  &c.flagTestVerifyOnly,
		Default: false,
		Hidden:  true,
	})

	// End internal-only flags.

	return set
}

func (c *AgentCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *AgentCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AgentCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// Create a logger. We wrap it in a gated writer so that it doesn't
	// start logging too early.
	c.logGate = &gatedwriter.Writer{Writer: os.Stderr}
	c.logWriter = c.logGate
	if c.flagCombineLogs {
		c.logWriter = os.Stdout
	}
	var level log.Level
	c.flagLogLevel = strings.ToLower(strings.TrimSpace(c.flagLogLevel))
	switch c.flagLogLevel {
	case "trace":
		level = log.Trace
	case "debug":
		level = log.Debug
	case "notice", "info", "":
		level = log.Info
	case "warn", "warning":
		level = log.Warn
	case "err", "error":
		level = log.Error
	default:
		c.UI.Error(fmt.Sprintf("Unknown log level: %s", c.flagLogLevel))
		return 1
	}

	if c.logger == nil {
		c.logger = logging.NewVaultLoggerWithWriter(c.logWriter, level)
	}

	// Validation
	if len(c.flagConfigs) != 1 {
		c.UI.Error("Must specify exactly one config path using -config")
		return 1
	}

	// Load the configuration
	config, err := config.LoadConfig(c.flagConfigs[0], c.logger)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error loading configuration from %s: %s", c.flagConfigs[0], err))
		return 1
	}

	// Ensure at least one config was found.
	if config == nil {
		c.UI.Output(wrapAtLength(
			"No configuration read. Please provide the configuration with the " +
				"-config flag."))
		return 1
	}
	if config.AutoAuth == nil && config.Cache == nil {
		c.UI.Error("No auto_auth or cache block found in config file")
		return 1
	}
	if config.AutoAuth == nil {
		c.UI.Info("No auto_auth block found in config file, not starting automatic authentication feature")
	}

	if config.Vault != nil {
		c.setStringFlag(f, config.Vault.Address, &StringVar{
			Name:    flagNameAddress,
			Target:  &c.flagAddress,
			Default: "https://127.0.0.1:8200",
			EnvVar:  api.EnvVaultAddress,
		})
		c.setStringFlag(f, config.Vault.CACert, &StringVar{
			Name:    flagNameCACert,
			Target:  &c.flagCACert,
			Default: "",
			EnvVar:  api.EnvVaultCACert,
		})
		c.setStringFlag(f, config.Vault.CAPath, &StringVar{
			Name:    flagNameCAPath,
			Target:  &c.flagCAPath,
			Default: "",
			EnvVar:  api.EnvVaultCAPath,
		})
		c.setStringFlag(f, config.Vault.ClientCert, &StringVar{
			Name:    flagNameClientCert,
			Target:  &c.flagClientCert,
			Default: "",
			EnvVar:  api.EnvVaultClientCert,
		})
		c.setStringFlag(f, config.Vault.ClientKey, &StringVar{
			Name:    flagNameClientKey,
			Target:  &c.flagClientKey,
			Default: "",
			EnvVar:  api.EnvVaultClientKey,
		})
		c.setBoolFlag(f, config.Vault.TLSSkipVerify, &BoolVar{
			Name:    flagNameTLSSkipVerify,
			Target:  &c.flagTLSSkipVerify,
			Default: false,
			EnvVar:  api.EnvVaultSkipVerify,
		})
	}

	infoKeys := make([]string, 0, 10)
	info := make(map[string]string)
	info["log level"] = c.flagLogLevel
	infoKeys = append(infoKeys, "log level")

	infoKeys = append(infoKeys, "version")
	verInfo := version.GetVersion()
	info["version"] = verInfo.FullVersionNumber(false)
	if verInfo.Revision != "" {
		info["version sha"] = strings.Trim(verInfo.Revision, "'")
		infoKeys = append(infoKeys, "version sha")
	}
	infoKeys = append(infoKeys, "cgo")
	info["cgo"] = "disabled"
	if version.CgoEnabled {
		info["cgo"] = "enabled"
	}

	// Tests might not want to start a vault server and just want to verify
	// the configuration.
	if c.flagTestVerifyOnly {
		if os.Getenv("VAULT_TEST_VERIFY_ONLY_DUMP_CONFIG") != "" {
			c.UI.Output(fmt.Sprintf(
				"\nConfiguration:\n%s\n",
				pretty.Sprint(*config)))
		}
		return 0
	}

	// Ignore any setting of agent's address. This client is used by the agent
	// to reach out to Vault. This should never loop back to agent.
	c.flagAgentAddress = ""
	client, err := c.Client()
	if err != nil {
		c.UI.Error(fmt.Sprintf(
			"Error fetching client: %v",
			err))
		return 1
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	var method auth.AuthMethod
	var sinks []*sink.SinkConfig
	if config.AutoAuth != nil {
		for _, sc := range config.AutoAuth.Sinks {
			switch sc.Type {
			case "file":
				config := &sink.SinkConfig{
					Logger:  c.logger.Named("sink.file"),
					Config:  sc.Config,
					Client:  client,
					WrapTTL: sc.WrapTTL,
					DHType:  sc.DHType,
					DHPath:  sc.DHPath,
					AAD:     sc.AAD,
				}
				s, err := file.NewFileSink(config)
				if err != nil {
					c.UI.Error(errwrap.Wrapf("Error creating file sink: {{err}}", err).Error())
					return 1
				}
				config.Sink = s
				sinks = append(sinks, config)
			default:
				c.UI.Error(fmt.Sprintf("Unknown sink type %q", sc.Type))
				return 1
			}
		}

		authConfig := &auth.AuthConfig{
			Logger:    c.logger.Named(fmt.Sprintf("auth.%s", config.AutoAuth.Method.Type)),
			MountPath: config.AutoAuth.Method.MountPath,
			Config:    config.AutoAuth.Method.Config,
		}
		switch config.AutoAuth.Method.Type {
		case "alicloud":
			method, err = alicloud.NewAliCloudAuthMethod(authConfig)
		case "aws":
			method, err = aws.NewAWSAuthMethod(authConfig)
		case "azure":
			method, err = azure.NewAzureAuthMethod(authConfig)
		case "gcp":
			method, err = gcp.NewGCPAuthMethod(authConfig)
		case "jwt":
			method, err = jwt.NewJWTAuthMethod(authConfig)
		case "kubernetes":
			method, err = kubernetes.NewKubernetesAuthMethod(authConfig)
		case "approle":
			method, err = approle.NewApproleAuthMethod(authConfig)
		default:
			c.UI.Error(fmt.Sprintf("Unknown auth method %q", config.AutoAuth.Method.Type))
			return 1
		}
		if err != nil {
			c.UI.Error(errwrap.Wrapf(fmt.Sprintf("Error creating %s auth method: {{err}}", config.AutoAuth.Method.Type), err).Error())
			return 1
		}
	}

	// Output the header that the server has started
	if !c.flagCombineLogs {
		c.UI.Output("==> Vault server started! Log data will stream in below:\n")
	}

	// Inform any tests that the server is ready
	select {
	case c.startedCh <- struct{}{}:
	default:
	}

	// Parse agent listener configurations
	if config.Cache != nil && len(config.Listeners) != 0 {
		cacheLogger := c.logger.Named("cache")

		// Create the API proxier
		apiProxy, err := cache.NewAPIProxy(&cache.APIProxyConfig{
			Client: client,
			Logger: cacheLogger.Named("apiproxy"),
		})
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error creating API proxy: %v", err))
			return 1
		}

		// Create the lease cache proxier and set its underlying proxier to
		// the API proxier.
		leaseCache, err := cache.NewLeaseCache(&cache.LeaseCacheConfig{
			Client:      client,
			BaseContext: ctx,
			Proxier:     apiProxy,
			Logger:      cacheLogger.Named("leasecache"),
		})
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error creating lease cache: %v", err))
			return 1
		}

		var inmemSink sink.Sink
		if config.Cache.UseAutoAuthToken {
			cacheLogger.Debug("auto-auth token is allowed to be used; configuring inmem sink")
			inmemSink, err = inmem.New(&sink.SinkConfig{
				Logger: cacheLogger,
			}, leaseCache)
			if err != nil {
				c.UI.Error(fmt.Sprintf("Error creating inmem sink for cache: %v", err))
				return 1
			}
			sinks = append(sinks, &sink.SinkConfig{
				Logger: cacheLogger,
				Sink:   inmemSink,
			})
		}

		// Create a muxer and add paths relevant for the lease cache layer
		mux := http.NewServeMux()
		mux.Handle(consts.AgentPathCacheClear, leaseCache.HandleCacheClear(ctx))

		mux.Handle("/", cache.Handler(ctx, cacheLogger, leaseCache, inmemSink))

		var listeners []net.Listener
		for i, lnConfig := range config.Listeners {
			ln, tlsConf, err := cache.StartListener(lnConfig)
			if err != nil {
				c.UI.Error(fmt.Sprintf("Error starting listener: %v", err))
				return 1
			}

			listeners = append(listeners, ln)

			scheme := "https://"
			if tlsConf == nil {
				scheme = "http://"
			}
			if ln.Addr().Network() == "unix" {
				scheme = "unix://"
			}

			infoKey := fmt.Sprintf("api address %d", i+1)
			info[infoKey] = scheme + ln.Addr().String()
			infoKeys = append(infoKeys, infoKey)

			server := &http.Server{
				Addr:              ln.Addr().String(),
				TLSConfig:         tlsConf,
				Handler:           mux,
				ReadHeaderTimeout: 10 * time.Second,
				ReadTimeout:       30 * time.Second,
				IdleTimeout:       5 * time.Minute,
				ErrorLog:          cacheLogger.StandardLogger(nil),
			}

			go server.Serve(ln)
		}

		// Ensure that listeners are closed at all the exits
		listenerCloseFunc := func() {
			for _, ln := range listeners {
				ln.Close()
			}
		}
		defer c.cleanupGuard.Do(listenerCloseFunc)
	}

	var ssDoneCh, ahDoneCh chan struct{}
	// Start auto-auth and sink servers
	if method != nil {
		ah := auth.NewAuthHandler(&auth.AuthHandlerConfig{
			Logger:                       c.logger.Named("auth.handler"),
			Client:                       c.client,
			WrapTTL:                      config.AutoAuth.Method.WrapTTL,
			EnableReauthOnNewCredentials: config.AutoAuth.EnableReauthOnNewCredentials,
		})
		ahDoneCh = ah.DoneCh

		ss := sink.NewSinkServer(&sink.SinkServerConfig{
			Logger:        c.logger.Named("sink.server"),
			Client:        client,
			ExitAfterAuth: config.ExitAfterAuth,
		})
		ssDoneCh = ss.DoneCh

		go ah.Run(ctx, method)
		go ss.Run(ctx, ah.OutputCh, sinks)
	}

	// Server configuration output
	padding := 24
	sort.Strings(infoKeys)
	c.UI.Output("==> Vault agent configuration:\n")
	for _, k := range infoKeys {
		c.UI.Output(fmt.Sprintf(
			"%s%s: %s",
			strings.Repeat(" ", padding-len(k)),
			strings.Title(k),
			info[k]))
	}
	c.UI.Output("")

	// Release the log gate.
	c.logGate.Flush()

	// Write out the PID to the file now that server has successfully started
	if err := c.storePidFile(config.PidFile); err != nil {
		c.UI.Error(fmt.Sprintf("Error storing PID: %s", err))
		return 1
	}

	defer func() {
		if err := c.removePidFile(config.PidFile); err != nil {
			c.UI.Error(fmt.Sprintf("Error deleting the PID file: %s", err))
		}
	}()

	select {
	case <-ssDoneCh:
		// This will happen if we exit-on-auth
		c.logger.Info("sinks finished, exiting")
	case <-c.ShutdownCh:
		c.UI.Output("==> Vault agent shutdown triggered")
		cancelFunc()
		if ahDoneCh != nil {
			<-ahDoneCh
		}
		if ssDoneCh != nil {
			<-ssDoneCh
		}
	}

	return 0
}

func (c *AgentCommand) setStringFlag(f *FlagSets, configVal string, fVar *StringVar) {
	var isFlagSet bool
	f.Visit(func(f *flag.Flag) {
		if f.Name == fVar.Name {
			isFlagSet = true
		}
	})

	flagEnvValue, flagEnvSet := os.LookupEnv(fVar.EnvVar)
	switch {
	case isFlagSet:
		// Don't do anything as the flag is already set from the command line
	case flagEnvSet:
		// Use value from env var
		*fVar.Target = flagEnvValue
	case configVal != "":
		// Use value from config
		*fVar.Target = configVal
	default:
		// Use the default value
		*fVar.Target = fVar.Default
	}
}

func (c *AgentCommand) setBoolFlag(f *FlagSets, configVal bool, fVar *BoolVar) {
	var isFlagSet bool
	f.Visit(func(f *flag.Flag) {
		if f.Name == fVar.Name {
			isFlagSet = true
		}
	})

	flagEnvValue, flagEnvSet := os.LookupEnv(fVar.EnvVar)
	switch {
	case isFlagSet:
		// Don't do anything as the flag is already set from the command line
	case flagEnvSet:
		// Use value from env var
		*fVar.Target = flagEnvValue != ""
	case configVal == true:
		// Use value from config
		*fVar.Target = configVal
	default:
		// Use the default value
		*fVar.Target = fVar.Default
	}
}

// storePidFile is used to write out our PID to a file if necessary
func (c *AgentCommand) storePidFile(pidPath string) error {
	// Quit fast if no pidfile
	if pidPath == "" {
		return nil
	}

	// Open the PID file
	pidFile, err := os.OpenFile(pidPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return errwrap.Wrapf("could not open pid file: {{err}}", err)
	}
	defer pidFile.Close()

	// Write out the PID
	pid := os.Getpid()
	_, err = pidFile.WriteString(fmt.Sprintf("%d", pid))
	if err != nil {
		return errwrap.Wrapf("could not write to pid file: {{err}}", err)
	}
	return nil
}

// removePidFile is used to cleanup the PID file if necessary
func (c *AgentCommand) removePidFile(pidPath string) error {
	if pidPath == "" {
		return nil
	}
	return os.Remove(pidPath)
}
