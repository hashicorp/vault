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

	token_file "github.com/hashicorp/vault/command/agent/auth/token-file"

	ctconfig "github.com/hashicorp/consul-template/config"
	"github.com/hashicorp/go-multierror"

	"github.com/hashicorp/vault/command/agent/sink/inmem"

	systemd "github.com/coreos/go-systemd/daemon"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/gatedwriter"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-secure-stdlib/reloadutil"
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
	"github.com/hashicorp/vault/command/agent/template"
	"github.com/hashicorp/vault/command/agent/winsvc"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/useragent"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/internalshared/listenerutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/version"
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

const (
	// flagNameAgentExitAfterAuth is used as an Agent specific flag to indicate
	// that agent should exit after a single successful auth
	flagNameAgentExitAfterAuth = "exit-after-auth"
)

type AgentCommand struct {
	*BaseCommand
	logFlags logFlags

	config *agentConfig.Config

	ShutdownCh chan struct{}
	SighupCh   chan struct{}

	tlsReloadFuncsLock sync.RWMutex
	tlsReloadFuncs     []reloadutil.ReloadFunc

	logWriter io.Writer
	logGate   *gatedwriter.Writer
	logger    log.Logger

	// Telemetry object
	metricsHelper *metricsutil.MetricsHelper

	cleanupGuard sync.Once

	startedCh  chan struct{} // for tests
	reloadedCh chan struct{} // for tests

	flagConfigs        []string
	flagExitAfterAuth  bool
	flagTestVerifyOnly bool
}

func (c *AgentCommand) Synopsis() string {
	return "Start a Vault agent"
}

func (c *AgentCommand) Help() string {
	helpText := `
Usage: vault agent [options]

  This command starts a Vault Agent that can perform automatic authentication
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

	// Augment with the log flags
	f.addLogFlags(&c.logFlags)

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

	f.BoolVar(&BoolVar{
		Name:    flagNameAgentExitAfterAuth,
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

	if c.logFlags.flagCombineLogs {
		c.logWriter = os.Stdout
	}

	// Validation
	if len(c.flagConfigs) < 1 {
		c.UI.Error("Must specify exactly at least one config path using -config")
		return 1
	}

	config, err := c.loadConfig(c.flagConfigs)
	if err != nil {
		c.outputErrors(err)
		return 1
	}

	if config.AutoAuth == nil {
		c.UI.Info("No auto_auth block found in config, the automatic authentication feature will not be started")
	}

	c.applyConfigOverrides(f, config) // This only needs to happen on start-up to aggregate config from flags and env vars
	c.config = config

	l, err := c.newLogger()
	if err != nil {
		c.outputErrors(err)
		return 1
	}
	c.logger = l

	infoKeys := make([]string, 0, 10)
	info := make(map[string]string)
	info["log level"] = config.LogLevel
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
				pretty.Sprint(*c.config)))
		}
		return 0
	}

	// Ignore any setting of Agent's address. This client is used by the Agent
	// to reach out to Vault. This should never loop back to agent.
	c.flagAgentAddress = ""
	client, err := c.Client()
	if err != nil {
		c.UI.Error(fmt.Sprintf(
			"Error fetching client: %v",
			err))
		return 1
	}

	serverHealth, err := client.Sys().Health()
	if err == nil {
		// We don't exit on error here, as this is not worth stopping Agent over
		serverVersion := serverHealth.Version
		agentVersion := version.GetVersion().VersionNumber()
		if serverVersion != agentVersion {
			c.UI.Info("==> Note: Vault Agent version does not match Vault server version. " +
				fmt.Sprintf("Vault Agent version: %s, Vault server version: %s", agentVersion, serverVersion))
		}
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
					c.UI.Error(fmt.Errorf("error creating file sink: %w", err).Error())
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
		case "token_file":
			method, err = token_file.NewTokenFileAuthMethod(authConfig)
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
		client.SetMaxRetries(ctconfig.DefaultRetryAttempts)
		if config.Vault != nil {
			if config.Vault.Retry != nil {
				client.SetMaxRetries(config.Vault.Retry.NumRetries)
			}
		}
	}

	enforceConsistency := cache.EnforceConsistencyNever
	whenInconsistent := cache.WhenInconsistentFail
	if config.APIProxy != nil {
		switch config.APIProxy.EnforceConsistency {
		case "always":
			enforceConsistency = cache.EnforceConsistencyAlways
		case "never", "":
		default:
			c.UI.Error(fmt.Sprintf("Unknown api_proxy setting for enforce_consistency: %q", config.APIProxy.EnforceConsistency))
			return 1
		}

		switch config.APIProxy.WhenInconsistent {
		case "retry":
			whenInconsistent = cache.WhenInconsistentRetry
		case "forward":
			whenInconsistent = cache.WhenInconsistentForward
		case "fail", "":
		default:
			c.UI.Error(fmt.Sprintf("Unknown api_proxy setting for when_inconsistent: %q", config.APIProxy.WhenInconsistent))
			return 1
		}
	}
	// Keep Cache configuration for legacy reasons, but error if defined alongside API Proxy
	if config.Cache != nil {
		switch config.Cache.EnforceConsistency {
		case "always":
			if enforceConsistency != cache.EnforceConsistencyNever {
				c.UI.Error("enforce_consistency configured in both api_proxy and cache blocks. Please remove this configuration from the cache block.")
				return 1
			} else {
				enforceConsistency = cache.EnforceConsistencyAlways
			}
		case "never", "":
		default:
			c.UI.Error(fmt.Sprintf("Unknown cache setting for enforce_consistency: %q", config.Cache.EnforceConsistency))
			return 1
		}

		switch config.Cache.WhenInconsistent {
		case "retry":
			if whenInconsistent != cache.WhenInconsistentFail {
				c.UI.Error("when_inconsistent configured in both api_proxy and cache blocks. Please remove this configuration from the cache block.")
				return 1
			} else {
				whenInconsistent = cache.WhenInconsistentRetry
			}
		case "forward":
			if whenInconsistent != cache.WhenInconsistentFail {
				c.UI.Error("when_inconsistent configured in both api_proxy and cache blocks. Please remove this configuration from the cache block.")
				return 1
			} else {
				whenInconsistent = cache.WhenInconsistentForward
			}
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
	if !c.logFlags.flagCombineLogs {
		c.UI.Output("==> Vault Agent started! Log data will stream in below:\n")
	}

	var leaseCache *cache.LeaseCache
	var previousToken string

	proxyClient, err := client.CloneWithHeaders()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error cloning client for proxying: %v", err))
		return 1
	}

	if config.DisableIdleConnsAPIProxy {
		proxyClient.SetMaxIdleConnections(-1)
	}

	if config.DisableKeepAlivesAPIProxy {
		proxyClient.SetDisableKeepAlives(true)
	}

	apiProxyLogger := c.logger.Named("apiproxy")

	// The API proxy to be used, if listeners are configured
	apiProxy, err := cache.NewAPIProxy(&cache.APIProxyConfig{
		Client:                 proxyClient,
		Logger:                 apiProxyLogger,
		EnforceConsistency:     enforceConsistency,
		WhenInconsistentAction: whenInconsistent,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating API proxy: %v", err))
		return 1
	}

	// Parse agent cache configurations
	if config.Cache != nil {
		cacheLogger := c.logger.Named("cache")

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

				km, err := keymanager.NewPassthroughKeyManager(ctx, token)
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
				km, err := keymanager.NewPassthroughKeyManager(ctx, nil)
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
				token, err := km.RetrievalToken(ctx)
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
	}

	var listeners []net.Listener

	// If there are templates, add an in-process listener
	if len(config.Templates) > 0 {
		config.Listeners = append(config.Listeners, &configutil.Listener{Type: listenerutil.BufConnType})
	}

	// Ensure we've added all the reload funcs for TLS before anyone triggers a reload.
	c.tlsReloadFuncsLock.Lock()

	for i, lnConfig := range config.Listeners {
		var ln net.Listener
		var tlsCfg *tls.Config

		if lnConfig.Type == listenerutil.BufConnType {
			inProcListener := bufconn.Listen(1024 * 1024)
			if config.Cache != nil {
				config.Cache.InProcDialer = listenerutil.NewBufConnWrapper(inProcListener)
			}
			ln = inProcListener
		} else {
			lnBundle, err := cache.StartListener(lnConfig)
			if err != nil {
				c.UI.Error(fmt.Sprintf("Error starting listener: %v", err))
				return 1
			}

			tlsCfg = lnBundle.TLSConfig
			ln = lnBundle.Listener

			// Track the reload func, so we can reload later if needed.
			c.tlsReloadFuncs = append(c.tlsReloadFuncs, lnBundle.TLSReloadFunc)
		}

		listeners = append(listeners, ln)

		proxyVaultToken := true
		var inmemSink sink.Sink
		if config.APIProxy != nil {
			if config.APIProxy.UseAutoAuthToken {
				apiProxyLogger.Debug("auto-auth token is allowed to be used; configuring inmem sink")
				inmemSink, err = inmem.New(&sink.SinkConfig{
					Logger: apiProxyLogger,
				}, leaseCache)
				if err != nil {
					c.UI.Error(fmt.Sprintf("Error creating inmem sink for cache: %v", err))
					return 1
				}
				sinks = append(sinks, &sink.SinkConfig{
					Logger: apiProxyLogger,
					Sink:   inmemSink,
				})
			}
			proxyVaultToken = !config.APIProxy.ForceAutoAuthToken
		}

		var muxHandler http.Handler
		if leaseCache != nil {
			muxHandler = cache.ProxyHandler(ctx, apiProxyLogger, leaseCache, inmemSink, proxyVaultToken)
		} else {
			muxHandler = cache.ProxyHandler(ctx, apiProxyLogger, apiProxy, inmemSink, proxyVaultToken)
		}

		// Parse 'require_request_header' listener config option, and wrap
		// the request handler if necessary
		if lnConfig.RequireRequestHeader && ("metrics_only" != lnConfig.Role) {
			muxHandler = verifyRequestHeader(muxHandler)
		}

		// Create a muxer and add paths relevant for the lease cache layer
		mux := http.NewServeMux()
		quitEnabled := lnConfig.AgentAPI != nil && lnConfig.AgentAPI.EnableQuit

		mux.Handle(consts.AgentPathMetrics, c.handleMetrics())
		if "metrics_only" != lnConfig.Role {
			mux.Handle(consts.AgentPathCacheClear, leaseCache.HandleCacheClear(ctx))
			mux.Handle(consts.AgentPathQuit, c.handleQuit(quitEnabled))
			mux.Handle("/", muxHandler)
		}

		scheme := "https://"
		if tlsCfg == nil {
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
			TLSConfig:         tlsCfg,
			Handler:           mux,
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       30 * time.Second,
			IdleTimeout:       5 * time.Minute,
			ErrorLog:          apiProxyLogger.StandardLogger(nil),
		}

		go server.Serve(ln)
	}

	c.tlsReloadFuncsLock.Unlock()

	// Ensure that listeners are closed at all the exits
	listenerCloseFunc := func() {
		for _, ln := range listeners {
			ln.Close()
		}
	}
	defer c.cleanupGuard.Do(listenerCloseFunc)

	// Inform any tests that the server is ready
	if c.startedCh != nil {
		close(c.startedCh)
	}

	var g run.Group

	g.Add(func() error {
		for {
			select {
			case <-c.SighupCh:
				c.UI.Output("==> Vault Agent config reload triggered")
				err := c.reloadConfig(c.flagConfigs)
				if err != nil {
					c.outputErrors(err)
				}
				// Send the 'reloaded' message on the relevant channel
				select {
				case c.reloadedCh <- struct{}{}:
				default:
				}
			case <-ctx.Done():
				return nil
			}
		}
	}, func(error) {
		cancelFunc()
	})

	// This run group watches for signal termination
	g.Add(func() error {
		for {
			select {
			case <-c.ShutdownCh:
				c.UI.Output("==> Vault Agent shutdown triggered")
				// Notify systemd that the server is shutting down
				// Let the lease cache know this is a shutdown; no need to evict everything
				if leaseCache != nil {
					leaseCache.SetShuttingDown(true)
				}
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
			ExitAfterAuth: config.ExitAfterAuth,
		})

		ts := template.NewServer(&template.ServerConfig{
			Logger:        c.logger.Named("template.server"),
			LogLevel:      c.logger.GetLevel(),
			LogWriter:     c.logWriter,
			AgentConfig:   c.config,
			Namespace:     templateNamespace,
			ExitAfterAuth: config.ExitAfterAuth,
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
	c.UI.Output("==> Vault Agent configuration:\n")
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

	var exitCode int
	if err := g.Run(); err != nil {
		c.logger.Error("runtime error encountered", "error", err)
		c.UI.Error("Error encountered during run, refer to logs for more details.")
		exitCode = 1
	}
	c.notifySystemd(systemd.SdNotifyStopping)
	return exitCode
}

// applyConfigOverrides ensures that the config object accurately reflects the desired
// settings as configured by the user. It applies the relevant config setting based
// on the precedence (env var overrides file config, cli overrides env var).
// It mutates the config object supplied.
func (c *AgentCommand) applyConfigOverrides(f *FlagSets, config *agentConfig.Config) {
	if config.Vault == nil {
		config.Vault = &agentConfig.Vault{}
	}

	f.applyLogConfigOverrides(config.SharedConfig)

	f.Visit(func(fl *flag.Flag) {
		if fl.Name == flagNameAgentExitAfterAuth {
			config.ExitAfterAuth = c.flagExitAfterAuth
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
}

// verifyRequestHeader wraps an http.Handler inside a Handler that checks for
// the request header that is used for SSRF protection.
func verifyRequestHeader(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if val, ok := r.Header[consts.RequestHeaderName]; !ok || len(val) != 1 || val[0] != "true" {
			logical.RespondError(w,
				http.StatusPreconditionFailed,
				fmt.Errorf("missing %q header", consts.RequestHeaderName))
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

// newLogger creates a logger based on parsed config field on the Agent Command struct.
func (c *AgentCommand) newLogger() (log.InterceptLogger, error) {
	if c.config == nil {
		return nil, fmt.Errorf("cannot create logger, no config")
	}

	var errors error

	// Parse all the log related config
	logLevel, err := logging.ParseLogLevel(c.config.LogLevel)
	if err != nil {
		errors = multierror.Append(errors, err)
	}

	logFormat, err := logging.ParseLogFormat(c.config.LogFormat)
	if err != nil {
		errors = multierror.Append(errors, err)
	}

	logRotateDuration, err := parseutil.ParseDurationSecond(c.config.LogRotateDuration)
	if err != nil {
		errors = multierror.Append(errors, err)
	}

	if errors != nil {
		return nil, errors
	}

	logCfg := &logging.LogConfig{
		Name:              "vault-agent",
		LogLevel:          logLevel,
		LogFormat:         logFormat,
		LogFilePath:       c.config.LogFile,
		LogRotateDuration: logRotateDuration,
		LogRotateBytes:    c.config.LogRotateBytes,
		LogRotateMaxFiles: c.config.LogRotateMaxFiles,
	}

	l, err := logging.Setup(logCfg, c.logWriter)
	if err != nil {
		return nil, err
	}

	return l, nil
}

// loadConfig attempts to generate an Agent config from the file(s) specified.
func (c *AgentCommand) loadConfig(paths []string) (*agentConfig.Config, error) {
	var errors error
	cfg := agentConfig.NewConfig()

	for _, configPath := range paths {
		configFromPath, err := agentConfig.LoadConfig(configPath)
		if err != nil {
			errors = multierror.Append(errors, fmt.Errorf("error loading configuration from %s: %w", configPath, err))
		} else {
			cfg = cfg.Merge(configFromPath)
		}
	}

	if errors != nil {
		return nil, errors
	}

	if err := cfg.ValidateConfig(); err != nil {
		return nil, fmt.Errorf("error validating configuration: %w", err)
	}

	return cfg, nil
}

// reloadConfig will attempt to reload the config from file(s) and adjust certain
// config values without requiring a restart of the Vault Agent.
// If config is retrieved without error it is stored in the config field of the AgentCommand.
// This operation is not atomic and could result in updated config but partially applied config settings.
// The error returned from this func may be a multierror.
// This function will most likely be called due to Vault Agent receiving a SIGHUP signal.
// Currently only reloading the following are supported:
// * log level
// * TLS certs for listeners
func (c *AgentCommand) reloadConfig(paths []string) error {
	// Notify systemd that the server is reloading
	c.notifySystemd(systemd.SdNotifyReloading)
	defer c.notifySystemd(systemd.SdNotifyReady)

	var errors error

	// Reload the config
	cfg, err := c.loadConfig(paths)
	if err != nil {
		// Returning single error as we won't continue with bad config and won't 'commit' it.
		return err
	}
	c.config = cfg

	// Update the log level
	err = c.reloadLogLevel()
	if err != nil {
		errors = multierror.Append(errors, err)
	}

	// Update certs
	err = c.reloadCerts()
	if err != nil {
		errors = multierror.Append(errors, err)
	}

	return errors
}

// reloadLogLevel will attempt to update the log level for the logger attached
// to the AgentComment struct using the value currently set in config.
func (c *AgentCommand) reloadLogLevel() error {
	logLevel, err := logging.ParseLogLevel(c.config.LogLevel)
	if err != nil {
		return err
	}

	c.logger.SetLevel(logLevel)

	return nil
}

// reloadCerts will attempt to reload certificates using a reload func which
// was provided when the listeners were configured, only funcs that were appended
// to the AgentCommand slice will be invoked.
// This function returns a multierror type so that every func can report an error
// if it encounters one.
func (c *AgentCommand) reloadCerts() error {
	var errors error

	c.tlsReloadFuncsLock.RLock()
	defer c.tlsReloadFuncsLock.RUnlock()

	for _, reloadFunc := range c.tlsReloadFuncs {
		err := reloadFunc()
		if err != nil {
			errors = multierror.Append(errors, err)
		}
	}

	return errors
}

// outputErrors will take an error or multierror and handle outputting each to the UI
func (c *AgentCommand) outputErrors(err error) {
	if err != nil {
		if me, ok := err.(*multierror.Error); ok {
			for _, err := range me.Errors {
				c.UI.Error(err.Error())
			}
		} else {
			c.UI.Error(err.Error())
		}
	}
}
