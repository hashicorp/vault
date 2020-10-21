package command

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/hashicorp/vault/command/agent/winsvc"
	"io"
	"net"
	"net/http"
	"os"
	"path"
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
	"github.com/hashicorp/vault/command/agent/auth/cert"
	"github.com/hashicorp/vault/command/agent/auth/cf"
	"github.com/hashicorp/vault/command/agent/auth/gcp"
	"github.com/hashicorp/vault/command/agent/auth/jwt"
	"github.com/hashicorp/vault/command/agent/auth/kerberos"
	"github.com/hashicorp/vault/command/agent/auth/kubernetes"
	"github.com/hashicorp/vault/command/agent/cache"
	agentConfig "github.com/hashicorp/vault/command/agent/config"
	"github.com/hashicorp/vault/command/agent/sink"
	"github.com/hashicorp/vault/command/agent/sink/file"
	"github.com/hashicorp/vault/command/agent/sink/inmem"
	"github.com/hashicorp/vault/command/agent/template"
	_ "github.com/hashicorp/vault/command/agent/winsvc"
	"github.com/hashicorp/vault/internalshared/gatedwriter"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/version"
	"github.com/kr/pretty"
	"github.com/mitchellh/cli"
	"github.com/oklog/run"
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

	flagConfigs       []string
	flagLogLevel      string
	flagExitAfterAuth bool

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

	f.BoolVar(&BoolVar{
		Name:    "exit-after-auth",
		Target:  &c.flagExitAfterAuth,
		Default: false,
		Usage: "If set to true, the agent will exit with code 0 after a single " +
			"successful auth, where success means that a token was retrieved and " +
			"all sinks successfully wrote it",
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
	c.logGate = gatedwriter.NewWriter(os.Stderr)
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
	config, err := agentConfig.LoadConfig(c.flagConfigs[0])
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

	// create an empty Vault configuration if none was loaded from file. The
	// follow-up setStringFlag calls will populate with defaults if otherwise
	// omitted
	if config.Vault == nil {
		config.Vault = new(agentConfig.Vault)
	}

	exitAfterAuth := config.ExitAfterAuth
	f.Visit(func(fl *flag.Flag) {
		if fl.Name == "exit-after-auth" {
			exitAfterAuth = c.flagExitAfterAuth
		}
	})

	c.setStringFlag(f, config.Vault.Address, &StringVar{
		Name:    flagNameAddress,
		Target:  &c.flagAddress,
		Default: "https://127.0.0.1:8200",
		EnvVar:  api.EnvVaultAddress,
	})
	config.Vault.Address = c.flagAddress
	c.setStringFlag(f, config.Vault.CACert, &StringVar{
		Name:    flagNameCACert,
		Target:  &c.flagCACert,
		Default: "",
		EnvVar:  api.EnvVaultCACert,
	})
	config.Vault.CACert = c.flagCACert
	c.setStringFlag(f, config.Vault.CAPath, &StringVar{
		Name:    flagNameCAPath,
		Target:  &c.flagCAPath,
		Default: "",
		EnvVar:  api.EnvVaultCAPath,
	})
	config.Vault.CAPath = c.flagCAPath
	c.setStringFlag(f, config.Vault.ClientCert, &StringVar{
		Name:    flagNameClientCert,
		Target:  &c.flagClientCert,
		Default: "",
		EnvVar:  api.EnvVaultClientCert,
	})
	config.Vault.ClientCert = c.flagClientCert
	c.setStringFlag(f, config.Vault.ClientKey, &StringVar{
		Name:    flagNameClientKey,
		Target:  &c.flagClientKey,
		Default: "",
		EnvVar:  api.EnvVaultClientKey,
	})
	config.Vault.ClientKey = c.flagClientKey
	c.setBoolFlag(f, config.Vault.TLSSkipVerify, &BoolVar{
		Name:    flagNameTLSSkipVerify,
		Target:  &c.flagTLSSkipVerify,
		Default: false,
		EnvVar:  api.EnvVaultSkipVerify,
	})
	config.Vault.TLSSkipVerify = c.flagTLSSkipVerify
	c.setStringFlag(f, config.Vault.TLSServerName, &StringVar{
		Name:    flagTLSServerName,
		Target:  &c.flagTLSServerName,
		Default: "",
		EnvVar:  api.EnvVaultTLSServerName,
	})
	config.Vault.TLSServerName = c.flagTLSServerName

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

	// ctx and cancelFunc are passed to the AuthHandler, SinkServer, and
	// TemplateServer that periodically listen for ctx.Done() to fire and shut
	// down accordingly.
	ctx, cancelFunc := context.WithCancel(context.Background())

	var method auth.AuthMethod
	var sinks []*sink.SinkConfig
	var namespace string
	if config.AutoAuth != nil {
		for _, sc := range config.AutoAuth.Sinks {
			switch sc.Type {
			case "file":
				config := &sink.SinkConfig{
					Logger:    c.logger.Named("sink.file"),
					Config:    sc.Config,
					Client:    client,
					WrapTTL:   sc.WrapTTL,
					DHType:    sc.DHType,
					DeriveKey: sc.DeriveKey,
					DHPath:    sc.DHPath,
					AAD:       sc.AAD,
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

		// Check if a default namespace has been set
		mountPath := config.AutoAuth.Method.MountPath
		if config.AutoAuth.Method.Namespace != "" {
			namespace = config.AutoAuth.Method.Namespace
			mountPath = path.Join(namespace, mountPath)
		}

		authConfig := &auth.AuthConfig{
			Logger:    c.logger.Named(fmt.Sprintf("auth.%s", config.AutoAuth.Method.Type)),
			MountPath: mountPath,
			Config:    config.AutoAuth.Method.Config,
		}
		switch config.AutoAuth.Method.Type {
		case "alicloud":
			method, err = alicloud.NewAliCloudAuthMethod(authConfig)
		case "aws":
			method, err = aws.NewAWSAuthMethod(authConfig)
		case "azure":
			method, err = azure.NewAzureAuthMethod(authConfig)
		case "cert":
			method, err = cert.NewCertAuthMethod(authConfig)
		case "cf":
			method, err = cf.NewCFAuthMethod(authConfig)
		case "gcp":
			method, err = gcp.NewGCPAuthMethod(authConfig)
		case "jwt":
			method, err = jwt.NewJWTAuthMethod(authConfig)
		case "kerberos":
			method, err = kerberos.NewKerberosAuthMethod(authConfig)
		case "kubernetes":
			method, err = kubernetes.NewKubernetesAuthMethod(authConfig)
		case "approle":
			method, err = approle.NewApproleAuthMethod(authConfig)
		case "pcf": // Deprecated.
			method, err = cf.NewCFAuthMethod(authConfig)
		default:
			c.UI.Error(fmt.Sprintf("Unknown auth method %q", config.AutoAuth.Method.Type))
			return 1
		}
		if err != nil {
			c.UI.Error(errwrap.Wrapf(fmt.Sprintf("Error creating %s auth method: {{err}}", config.AutoAuth.Method.Type), err).Error())
			return 1
		}
	}

	// Warn if cache _and_ cert auto-auth is enabled but certificates were not
	// provided in the auto_auth.method["cert"].config stanza.
	if config.Cache != nil && (config.AutoAuth != nil && config.AutoAuth.Method != nil && config.AutoAuth.Method.Type == "cert") {
		_, okCertFile := config.AutoAuth.Method.Config["client_cert"]
		_, okCertKey := config.AutoAuth.Method.Config["client_key"]

		// If neither of these exists in the cert stanza, agent will use the
		// certs from the vault stanza.
		if !okCertFile && !okCertKey {
			c.UI.Warn(wrapAtLength("WARNING! Cache is enabled and using the same certificates " +
				"from the 'cert' auto-auth method specified in the 'vault' stanza. Consider " +
				"specifying certificate information in the 'cert' auto-auth's config stanza."))
		}

	}

	// Output the header that the agent has started
	if !c.flagCombineLogs {
		c.UI.Output("==> Vault agent started! Log data will stream in below:\n")
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

		var proxyVaultToken = !config.Cache.ForceAutoAuthToken

		// Create the request handler
		cacheHandler := cache.Handler(ctx, cacheLogger, leaseCache, inmemSink, proxyVaultToken)

		var listeners []net.Listener
		for i, lnConfig := range config.Listeners {
			ln, tlsConf, err := cache.StartListener(lnConfig)
			if err != nil {
				c.UI.Error(fmt.Sprintf("Error starting listener: %v", err))
				return 1
			}

			listeners = append(listeners, ln)

			// Parse 'require_request_header' listener config option, and wrap
			// the request handler if necessary
			muxHandler := cacheHandler
			if lnConfig.RequireRequestHeader {
				muxHandler = verifyRequestHeader(muxHandler)
			}

			// Create a muxer and add paths relevant for the lease cache layer
			mux := http.NewServeMux()
			mux.Handle(consts.AgentPathCacheClear, leaseCache.HandleCacheClear(ctx))
			mux.Handle("/", muxHandler)

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

	// Listen for signals
	// TODO: implement support for SIGHUP reloading of configuration
	// signal.Notify(c.signalCh)

	var g run.Group

	// This run group watches for signal termination
	g.Add(func() error {
		for {
			select {
			case <-c.ShutdownCh:
				c.UI.Output("==> Vault agent shutdown triggered")
				return nil
			case <-ctx.Done():
				return nil
			case <-winsvc.ShutdownChannel():
				return nil
			}
		}
	}, func(error) {})

	// Start auto-auth and sink servers
	if method != nil {
		enableTokenCh := len(config.Templates) > 0
		ah := auth.NewAuthHandler(&auth.AuthHandlerConfig{
			Logger:                       c.logger.Named("auth.handler"),
			Client:                       c.client,
			WrapTTL:                      config.AutoAuth.Method.WrapTTL,
			EnableReauthOnNewCredentials: config.AutoAuth.EnableReauthOnNewCredentials,
			EnableTemplateTokenCh:        enableTokenCh,
		})

		ss := sink.NewSinkServer(&sink.SinkServerConfig{
			Logger:        c.logger.Named("sink.server"),
			Client:        client,
			ExitAfterAuth: exitAfterAuth,
		})

		ts := template.NewServer(&template.ServerConfig{
			Logger:        c.logger.Named("template.server"),
			LogLevel:      level,
			LogWriter:     c.logWriter,
			VaultConf:     config.Vault,
			Namespace:     namespace,
			ExitAfterAuth: exitAfterAuth,
		})

		g.Add(func() error {
			return ah.Run(ctx, method)
		}, func(error) {
			cancelFunc()
		})

		g.Add(func() error {
			err := ss.Run(ctx, ah.OutputCh, sinks)
			c.logger.Info("sinks finished, exiting")

			// Start goroutine to drain from ah.OutputCh from this point onward
			// to prevent ah.Run from being blocked.
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					case <-ah.OutputCh:
					}
				}
			}()

			// Wait until templates are rendered
			if len(config.Templates) > 0 {
				<-ts.DoneCh
			}

			return err
		}, func(error) {
			cancelFunc()
		})

		g.Add(func() error {
			return ts.Run(ctx, ah.TemplateTokenCh, config.Templates)
		}, func(error) {
			cancelFunc()
			ts.Stop()
		})

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

	if err := g.Run(); err != nil {
		c.logger.Error("runtime error encountered", "error", err)
		c.UI.Error("Error encountered during run, refer to logs for more details.")
		return 1
	}

	return 0
}

// verifyRequestHeader wraps an http.Handler inside a Handler that checks for
// the request header that is used for SSRF protection.
func verifyRequestHeader(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if val, ok := r.Header[consts.RequestHeaderName]; !ok || len(val) != 1 || val[0] != "true" {
			logical.RespondError(w,
				http.StatusPreconditionFailed,
				errors.New(fmt.Sprintf("missing '%s' header", consts.RequestHeaderName)))
			return
		}

		handler.ServeHTTP(w, r)
	})
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
