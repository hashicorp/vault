package command

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	systemd "github.com/coreos/go-systemd/daemon"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/gatedwriter"
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
	"github.com/hashicorp/vault/command/agent/cache/cacheboltdb"
	"github.com/hashicorp/vault/command/agent/cache/cachememdb"
	"github.com/hashicorp/vault/command/agent/cache/keymanager"
	agentConfig "github.com/hashicorp/vault/command/agent/config"
	"github.com/hashicorp/vault/command/agent/sink"
	"github.com/hashicorp/vault/command/agent/sink/file"
	"github.com/hashicorp/vault/command/agent/sink/inmem"
	"github.com/hashicorp/vault/command/agent/template"
	"github.com/hashicorp/vault/command/agent/winsvc"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/internalshared/listenerutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/useragent"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/version"
	"github.com/kr/pretty"
	"github.com/mitchellh/cli"
	"github.com/oklog/run"
	"github.com/posener/complete"
	"google.golang.org/grpc/test/bufconn"
)

var (
	_ cli.Command             = (*AgentCommand)(nil)
	_ cli.CommandAutocomplete = (*AgentCommand)(nil)
)

type AgentCommand struct {
	*BaseCommand

	ShutdownCh chan struct{}
	SighupCh   chan struct{}

	logWriter io.Writer
	logGate   *gatedwriter.Writer
	logger    log.Logger

	// Telemetry object
	metricsHelper *metricsutil.MetricsHelper

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
		Completion: complete.PredictSet("trace", "debug", "info", "warn", "error"),
		Usage: "Log verbosity level. Supported values (in order of detail) are " +
			"\"trace\", \"debug\", \"info\", \"warn\", and \"error\".",
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
	defer cancelFunc()

	// telemetry configuration
	inmemMetrics, _, prometheusEnabled, err := configutil.SetupTelemetry(&configutil.SetupTelemetryOpts{
		Config:      config.Telemetry,
		Ui:          c.UI,
		ServiceName: "vault",
		DisplayName: "Vault",
		UserAgent:   useragent.String(),
		ClusterName: config.ClusterName,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error initializing telemetry: %s", err))
		return 1
	}
	c.metricsHelper = metricsutil.NewMetricsHelper(inmemMetrics, prometheusEnabled)

	var method auth.AuthMethod
	var sinks []*sink.SinkConfig
	var templateNamespace string
	if config.AutoAuth != nil {
		if client.Headers().Get(consts.NamespaceHeaderName) == "" && config.AutoAuth.Method.Namespace != "" {
			client.SetNamespace(config.AutoAuth.Method.Namespace)
		}
		templateNamespace = client.Headers().Get(consts.NamespaceHeaderName)

		sinkClient, err := client.CloneWithHeaders()
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error cloning client for file sink: %v", err))
			return 1
		}

		if config.DisableIdleConnsAutoAuth {
			sinkClient.SetMaxIdleConnections(-1)
		}

		if config.DisableKeepAlivesAutoAuth {
			sinkClient.SetDisableKeepAlives(true)
		}

		for _, sc := range config.AutoAuth.Sinks {
			switch sc.Type {
			case "file":
				config := &sink.SinkConfig{
					Logger:    c.logger.Named("sink.file"),
					Config:    sc.Config,
					Client:    sinkClient,
					WrapTTL:   sc.WrapTTL,
					DHType:    sc.DHType,
					DeriveKey: sc.DeriveKey,
					DHPath:    sc.DHPath,
					AAD:       sc.AAD,
				}
				s, err := file.NewFileSink(config)
				if err != nil {
					c.UI.Error(fmt.Errorf("Error creating file sink: %w", err).Error())
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
			c.UI.Error(fmt.Errorf("Error creating %s auth method: %w", config.AutoAuth.Method.Type, err).Error())
			return 1
		}
	}

	// We do this after auto-auth has been configured, because we don't want to
	// confuse the issue of retries for auth failures which have their own
	// config and are handled a bit differently.
	if os.Getenv(api.EnvVaultMaxRetries) == "" {
		client.SetMaxRetries(config.Vault.Retry.NumRetries)
	}

	enforceConsistency := cache.EnforceConsistencyNever
	whenInconsistent := cache.WhenInconsistentFail
	if config.Cache != nil {
		switch config.Cache.EnforceConsistency {
		case "always":
			enforceConsistency = cache.EnforceConsistencyAlways
		case "never", "":
		default:
			c.UI.Error(fmt.Sprintf("Unknown cache setting for enforce_consistency: %q", config.Cache.EnforceConsistency))
			return 1
		}

		switch config.Cache.WhenInconsistent {
		case "retry":
			whenInconsistent = cache.WhenInconsistentRetry
		case "forward":
			whenInconsistent = cache.WhenInconsistentForward
		case "fail", "":
		default:
			c.UI.Error(fmt.Sprintf("Unknown cache setting for when_inconsistent: %q", config.Cache.WhenInconsistent))
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

	var leaseCache *cache.LeaseCache
	var previousToken string
	// Parse agent listener configurations
	if config.Cache != nil {
		cacheLogger := c.logger.Named("cache")

		proxyClient, err := client.CloneWithHeaders()
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error cloning client for caching: %v", err))
			return 1
		}

		if config.DisableIdleConnsCaching {
			proxyClient.SetMaxIdleConnections(-1)
		}

		if config.DisableKeepAlivesCaching {
			proxyClient.SetDisableKeepAlives(true)
		}

		// Create the API proxier
		apiProxy, err := cache.NewAPIProxy(&cache.APIProxyConfig{
			Client:                 proxyClient,
			Logger:                 cacheLogger.Named("apiproxy"),
			EnforceConsistency:     enforceConsistency,
			WhenInconsistentAction: whenInconsistent,
		})
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error creating API proxy: %v", err))
			return 1
		}

		// Create the lease cache proxier and set its underlying proxier to
		// the API proxier.
		leaseCache, err = cache.NewLeaseCache(&cache.LeaseCacheConfig{
			Client:      proxyClient,
			BaseContext: ctx,
			Proxier:     apiProxy,
			Logger:      cacheLogger.Named("leasecache"),
		})
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error creating lease cache: %v", err))
			return 1
		}

		// Configure persistent storage and add to LeaseCache
		if config.Cache.Persist != nil {
			if config.Cache.Persist.Path == "" {
				c.UI.Error("must specify persistent cache path")
				return 1
			}

			// Set AAD based on key protection type
			var aad string
			switch config.Cache.Persist.Type {
			case "kubernetes":
				aad, err = getServiceAccountJWT(config.Cache.Persist.ServiceAccountTokenFile)
				if err != nil {
					c.UI.Error(fmt.Sprintf("failed to read service account token from %s: %s", config.Cache.Persist.ServiceAccountTokenFile, err))
					return 1
				}
			default:
				c.UI.Error(fmt.Sprintf("persistent key protection type %q not supported", config.Cache.Persist.Type))
				return 1
			}

			// Check if bolt file exists already
			dbFileExists, err := cacheboltdb.DBFileExists(config.Cache.Persist.Path)
			if err != nil {
				c.UI.Error(fmt.Sprintf("failed to check if bolt file exists at path %s: %s", config.Cache.Persist.Path, err))
				return 1
			}
			if dbFileExists {
				// Open the bolt file, but wait to setup Encryption
				ps, err := cacheboltdb.NewBoltStorage(&cacheboltdb.BoltStorageConfig{
					Path:   config.Cache.Persist.Path,
					Logger: cacheLogger.Named("cacheboltdb"),
				})
				if err != nil {
					c.UI.Error(fmt.Sprintf("Error opening persistent cache: %v", err))
					return 1
				}

				// Get the token from bolt for retrieving the encryption key,
				// then setup encryption so that restore is possible
				token, err := ps.GetRetrievalToken()
				if err != nil {
					c.UI.Error(fmt.Sprintf("Error getting retrieval token from persistent cache: %v", err))
				}

				if err := ps.Close(); err != nil {
					c.UI.Warn(fmt.Sprintf("Failed to close persistent cache file after getting retrieval token: %s", err))
				}

				km, err := keymanager.NewPassthroughKeyManager(token)
				if err != nil {
					c.UI.Error(fmt.Sprintf("failed to configure persistence encryption for cache: %s", err))
					return 1
				}

				// Open the bolt file with the wrapper provided
				ps, err = cacheboltdb.NewBoltStorage(&cacheboltdb.BoltStorageConfig{
					Path:    config.Cache.Persist.Path,
					Logger:  cacheLogger.Named("cacheboltdb"),
					Wrapper: km.Wrapper(),
					AAD:     aad,
				})
				if err != nil {
					c.UI.Error(fmt.Sprintf("Error opening persistent cache with wrapper: %v", err))
					return 1
				}

				// Restore anything in the persistent cache to the memory cache
				if err := leaseCache.Restore(ctx, ps); err != nil {
					c.UI.Error(fmt.Sprintf("Error restoring in-memory cache from persisted file: %v", err))
					if config.Cache.Persist.ExitOnErr {
						return 1
					}
				}
				cacheLogger.Info("loaded memcache from persistent storage")

				// Check for previous auto-auth token
				oldTokenBytes, err := ps.GetAutoAuthToken(ctx)
				if err != nil {
					c.UI.Error(fmt.Sprintf("Error in fetching previous auto-auth token: %s", err))
					if config.Cache.Persist.ExitOnErr {
						return 1
					}
				}
				if len(oldTokenBytes) > 0 {
					oldToken, err := cachememdb.Deserialize(oldTokenBytes)
					if err != nil {
						c.UI.Error(fmt.Sprintf("Error in deserializing previous auto-auth token cache entry: %s", err))
						if config.Cache.Persist.ExitOnErr {
							return 1
						}
					}
					previousToken = oldToken.Token
				}

				// If keep_after_import true, set persistent storage layer in
				// leaseCache, else remove db file
				if config.Cache.Persist.KeepAfterImport {
					defer ps.Close()
					leaseCache.SetPersistentStorage(ps)
				} else {
					if err := ps.Close(); err != nil {
						c.UI.Warn(fmt.Sprintf("failed to close persistent cache file: %s", err))
					}
					dbFile := filepath.Join(config.Cache.Persist.Path, cacheboltdb.DatabaseFileName)
					if err := os.Remove(dbFile); err != nil {
						c.UI.Error(fmt.Sprintf("failed to remove persistent storage file %s: %s", dbFile, err))
						if config.Cache.Persist.ExitOnErr {
							return 1
						}
					}
				}
			} else {
				km, err := keymanager.NewPassthroughKeyManager(nil)
				if err != nil {
					c.UI.Error(fmt.Sprintf("failed to configure persistence encryption for cache: %s", err))
					return 1
				}
				ps, err := cacheboltdb.NewBoltStorage(&cacheboltdb.BoltStorageConfig{
					Path:    config.Cache.Persist.Path,
					Logger:  cacheLogger.Named("cacheboltdb"),
					Wrapper: km.Wrapper(),
					AAD:     aad,
				})
				if err != nil {
					c.UI.Error(fmt.Sprintf("Error creating persistent cache: %v", err))
					return 1
				}
				cacheLogger.Info("configured persistent storage", "path", config.Cache.Persist.Path)

				// Stash the key material in bolt
				token, err := km.RetrievalToken()
				if err != nil {
					c.UI.Error(fmt.Sprintf("Error getting persistent key: %s", err))
					return 1
				}
				if err := ps.StoreRetrievalToken(token); err != nil {
					c.UI.Error(fmt.Sprintf("Error setting key in persistent cache: %v", err))
					return 1
				}

				defer ps.Close()
				leaseCache.SetPersistentStorage(ps)
			}
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

		proxyVaultToken := !config.Cache.ForceAutoAuthToken

		// Create the request handler
		cacheHandler := cache.Handler(ctx, cacheLogger, leaseCache, inmemSink, proxyVaultToken)

		var listeners []net.Listener

		// If there are templates, add an in-process listener
		if len(config.Templates) > 0 {
			config.Listeners = append(config.Listeners, &configutil.Listener{Type: listenerutil.BufConnType})
		}
		for i, lnConfig := range config.Listeners {
			var ln net.Listener
			var tlsConf *tls.Config

			if lnConfig.Type == listenerutil.BufConnType {
				inProcListener := bufconn.Listen(1024 * 1024)
				config.Cache.InProcDialer = listenerutil.NewBufConnWrapper(inProcListener)
				ln = inProcListener
			} else {
				ln, tlsConf, err = cache.StartListener(lnConfig)
				if err != nil {
					c.UI.Error(fmt.Sprintf("Error starting listener: %v", err))
					return 1
				}
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
			quitEnabled := lnConfig.AgentAPI != nil && lnConfig.AgentAPI.EnableQuit

			mux.Handle(consts.AgentPathCacheClear, leaseCache.HandleCacheClear(ctx))
			mux.Handle(consts.AgentPathQuit, c.handleQuit(quitEnabled))
			mux.Handle(consts.AgentPathMetrics, c.handleMetrics())
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

	// Inform any tests that the server is ready
	if c.startedCh != nil {
		close(c.startedCh)
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
				// Notify systemd that the server is shutting down
				c.notifySystemd(systemd.SdNotifyStopping)
				// Let the lease cache know this is a shutdown; no need to evict
				// everything
				if leaseCache != nil {
					leaseCache.SetShuttingDown(true)
				}
				return nil
			case <-ctx.Done():
				c.notifySystemd(systemd.SdNotifyStopping)
				return nil
			case <-winsvc.ShutdownChannel():
				return nil
			}
		}
	}, func(error) {})

	// Start auto-auth and sink servers
	if method != nil {
		enableTokenCh := len(config.Templates) > 0

		// Auth Handler is going to set its own retry values, so we want to
		// work on a copy of the client to not affect other subsystems.
		ahClient, err := c.client.CloneWithHeaders()
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error cloning client for auth handler: %v", err))
			return 1
		}

		if config.DisableIdleConnsAutoAuth {
			ahClient.SetMaxIdleConnections(-1)
		}

		if config.DisableKeepAlivesAutoAuth {
			ahClient.SetDisableKeepAlives(true)
		}

		ah := auth.NewAuthHandler(&auth.AuthHandlerConfig{
			Logger:                       c.logger.Named("auth.handler"),
			Client:                       ahClient,
			WrapTTL:                      config.AutoAuth.Method.WrapTTL,
			MinBackoff:                   config.AutoAuth.Method.MinBackoff,
			MaxBackoff:                   config.AutoAuth.Method.MaxBackoff,
			EnableReauthOnNewCredentials: config.AutoAuth.EnableReauthOnNewCredentials,
			EnableTemplateTokenCh:        enableTokenCh,
			Token:                        previousToken,
			ExitOnError:                  config.AutoAuth.Method.ExitOnError,
		})

		ss := sink.NewSinkServer(&sink.SinkServerConfig{
			Logger:        c.logger.Named("sink.server"),
			Client:        ahClient,
			ExitAfterAuth: exitAfterAuth,
		})

		ts := template.NewServer(&template.ServerConfig{
			Logger:        c.logger.Named("template.server"),
			LogLevel:      level,
			LogWriter:     c.logWriter,
			AgentConfig:   config,
			Namespace:     templateNamespace,
			ExitAfterAuth: exitAfterAuth,
		})

		g.Add(func() error {
			return ah.Run(ctx, method)
		}, func(error) {
			// Let the lease cache know this is a shutdown; no need to evict
			// everything
			if leaseCache != nil {
				leaseCache.SetShuttingDown(true)
			}
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
			// Let the lease cache know this is a shutdown; no need to evict
			// everything
			if leaseCache != nil {
				leaseCache.SetShuttingDown(true)
			}
			cancelFunc()
		})

		g.Add(func() error {
			return ts.Run(ctx, ah.TemplateTokenCh, config.Templates)
		}, func(error) {
			// Let the lease cache know this is a shutdown; no need to evict
			// everything
			if leaseCache != nil {
				leaseCache.SetShuttingDown(true)
			}
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

	// Notify systemd that the server is ready (if applicable)
	c.notifySystemd(systemd.SdNotifyReady)

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
				fmt.Errorf("missing '%s' header", consts.RequestHeaderName))
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (c *AgentCommand) notifySystemd(status string) {
	sent, err := systemd.SdNotify(false, status)
	if err != nil {
		c.logger.Error("error notifying systemd", "error", err)
	} else {
		if sent {
			c.logger.Debug("sent systemd notification", "notification", status)
		} else {
			c.logger.Debug("would have sent systemd notification (systemd not present)", "notification", status)
		}
	}
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
	case configVal:
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
	pidFile, err := os.OpenFile(pidPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("could not open pid file: %w", err)
	}
	defer pidFile.Close()

	// Write out the PID
	pid := os.Getpid()
	_, err = pidFile.WriteString(fmt.Sprintf("%d", pid))
	if err != nil {
		return fmt.Errorf("could not write to pid file: %w", err)
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

// GetServiceAccountJWT reads the service account jwt from `tokenFile`. Default is
// the default service account file path in kubernetes.
func getServiceAccountJWT(tokenFile string) (string, error) {
	if len(tokenFile) == 0 {
		tokenFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	}
	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(token)), nil
}

func (c *AgentCommand) handleMetrics() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			logical.RespondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		if err := r.ParseForm(); err != nil {
			logical.RespondError(w, http.StatusBadRequest, err)
			return
		}

		format := r.Form.Get("format")
		if format == "" {
			format = metricsutil.FormatFromRequest(&logical.Request{
				Headers: r.Header,
			})
		}

		resp := c.metricsHelper.ResponseForFormat(format)

		status := resp.Data[logical.HTTPStatusCode].(int)
		w.Header().Set("Content-Type", resp.Data[logical.HTTPContentType].(string))
		switch v := resp.Data[logical.HTTPRawBody].(type) {
		case string:
			w.WriteHeader((status))
			w.Write([]byte(v))
		case []byte:
			w.WriteHeader(status)
			w.Write(v)
		default:
			logical.RespondError(w, http.StatusInternalServerError, fmt.Errorf("wrong response returned"))
		}
	})
}

func (c *AgentCommand) handleQuit(enabled bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !enabled {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		switch r.Method {
		case http.MethodPost:
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		c.logger.Debug("received quit request")
		close(c.ShutdownCh)
	})
}
