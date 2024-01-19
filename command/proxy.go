// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"context"
	"crypto/tls"
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

	systemd "github.com/coreos/go-systemd/daemon"
	"github.com/hashicorp/cli"
	ctconfig "github.com/hashicorp/consul-template/config"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/gatedwriter"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-secure-stdlib/reloadutil"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agentproxyshared"
	"github.com/hashicorp/vault/command/agentproxyshared/auth"
	"github.com/hashicorp/vault/command/agentproxyshared/cache"
	"github.com/hashicorp/vault/command/agentproxyshared/sink"
	"github.com/hashicorp/vault/command/agentproxyshared/sink/file"
	"github.com/hashicorp/vault/command/agentproxyshared/sink/inmem"
	"github.com/hashicorp/vault/command/agentproxyshared/winsvc"
	proxyConfig "github.com/hashicorp/vault/command/proxy/config"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/useragent"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/internalshared/listenerutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/version"
	"github.com/kr/pretty"
	"github.com/oklog/run"
	"github.com/posener/complete"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/grpc/test/bufconn"
)

var (
	_ cli.Command             = (*ProxyCommand)(nil)
	_ cli.CommandAutocomplete = (*ProxyCommand)(nil)
)

const (
	// flagNameProxyExitAfterAuth is used as a Proxy specific flag to indicate
	// that proxy should exit after a single successful auth
	flagNameProxyExitAfterAuth = "exit-after-auth"
	nameProxy                  = "proxy"
)

type ProxyCommand struct {
	*BaseCommand
	logFlags logFlags

	config *proxyConfig.Config

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

func (c *ProxyCommand) Synopsis() string {
	return "Start a Vault Proxy"
}

func (c *ProxyCommand) Help() string {
	helpText := `
Usage: vault proxy [options]

  This command starts a Vault Proxy that can perform automatic authentication
  in certain environments.

  Start a proxy with a configuration file:

      $ vault proxy -config=/etc/vault/config.hcl

  For a full list of examples, please see the documentation.

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *ProxyCommand) Flags() *FlagSets {
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
			"contain only proxy directives.",
	})

	f.BoolVar(&BoolVar{
		Name:    flagNameProxyExitAfterAuth,
		Target:  &c.flagExitAfterAuth,
		Default: false,
		Usage: "If set to true, the proxy will exit with code 0 after a single " +
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

func (c *ProxyCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ProxyCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ProxyCommand) Run(args []string) int {
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

	// release log gate if the disable-gated-logs flag is set
	if c.logFlags.flagDisableGatedLogs {
		c.logGate.Flush()
	}

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

	// Ignore any setting of Agent/Proxy's address. This client is used by the Proxy
	// to reach out to Vault. This should never loop back to the proxy.
	c.flagAgentProxyAddress = ""
	client, err := c.Client()
	if err != nil {
		c.UI.Error(fmt.Sprintf(
			"Error fetching client: %v",
			err))
		return 1
	}

	serverHealth, err := client.Sys().Health()
	// We don't have any special behaviour if the error != nil, as this
	// is not worth stopping the Proxy process over.
	if err == nil {
		// Note that we don't exit if the versions don't match, as this is a valid
		// configuration, but we should still let the user know.
		serverVersion := serverHealth.Version
		proxyVersion := version.GetVersion().VersionNumber()
		if serverVersion != proxyVersion {
			c.UI.Info("==> Note: Vault Proxy version does not match Vault server version. " +
				fmt.Sprintf("Vault Proxy version: %s, Vault server version: %s", proxyVersion, serverVersion))
		}
	}

	// telemetry configuration
	inmemMetrics, _, prometheusEnabled, err := configutil.SetupTelemetry(&configutil.SetupTelemetryOpts{
		Config:      config.Telemetry,
		Ui:          c.UI,
		ServiceName: "vault",
		DisplayName: "Vault",
		UserAgent:   useragent.ProxyString(),
		ClusterName: config.ClusterName,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error initializing telemetry: %s", err))
		return 1
	}
	c.metricsHelper = metricsutil.NewMetricsHelper(inmemMetrics, prometheusEnabled)

	// This indicates whether the namespace for the client has been set by environment variable.
	// If it has, we don't touch it
	namespaceSetByEnvironmentVariable := client.Namespace() != ""

	if !namespaceSetByEnvironmentVariable && config.Vault != nil && config.Vault.Namespace != "" {
		client.SetNamespace(config.Vault.Namespace)
	}

	var method auth.AuthMethod
	var sinks []*sink.SinkConfig
	if config.AutoAuth != nil {
		// Note: This will only set namespace header to the value in config.AutoAuth.Method.Namespace
		// only if it hasn't been set by config.Vault.Namespace above. In that case, the config value
		// present at config.AutoAuth.Method.Namespace will still be used for auto-auth.
		if !namespaceSetByEnvironmentVariable && config.AutoAuth.Method.Namespace != "" {
			client.SetNamespace(config.AutoAuth.Method.Namespace)
		}

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
		method, err = agentproxyshared.GetAutoAuthMethodFromConfig(config.AutoAuth.Method.Type, authConfig, config.Vault.Address)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error creating %s auth method: %v", config.AutoAuth.Method.Type, err))
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

	// Warn if cache _and_ cert auto-auth is enabled but certificates were not
	// provided in the auto_auth.method["cert"].config stanza.
	if config.Cache != nil && (config.AutoAuth != nil && config.AutoAuth.Method != nil && config.AutoAuth.Method.Type == "cert") {
		_, okCertFile := config.AutoAuth.Method.Config["client_cert"]
		_, okCertKey := config.AutoAuth.Method.Config["client_key"]

		// If neither of these exists in the cert stanza, proxy will use the
		// certs from the vault stanza.
		if !okCertFile && !okCertKey {
			c.UI.Warn(wrapAtLength("WARNING! Cache is enabled and using the same certificates " +
				"from the 'cert' auto-auth method specified in the 'vault' stanza. Consider " +
				"specifying certificate information in the 'cert' auto-auth's config stanza."))
		}

	}

	// Output the header that the proxy has started
	if !c.logFlags.flagCombineLogs {
		c.UI.Output("==> Vault Proxy started! Log data will stream in below:\n")
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
		Client:                     proxyClient,
		Logger:                     apiProxyLogger,
		EnforceConsistency:         enforceConsistency,
		WhenInconsistentAction:     whenInconsistent,
		UserAgentStringFunction:    useragent.ProxyStringWithProxiedUserAgent,
		UserAgentString:            useragent.ProxyAPIProxyString(),
		PrependConfiguredNamespace: config.APIProxy != nil && config.APIProxy.PrependConfiguredNamespace,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating API proxy: %v", err))
		return 1
	}

	// ctx and cancelFunc are passed to the AuthHandler, SinkServer,
	// and other subsystems, so that they can listen for ctx.Done() to
	// fire and shut down accordingly.
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	var updater *cache.StaticSecretCacheUpdater

	// Parse proxy cache configurations
	if config.Cache != nil {
		cacheLogger := c.logger.Named("cache")

		// Create the lease cache proxier and set its underlying proxier to
		// the API proxier.
		leaseCache, err = cache.NewLeaseCache(&cache.LeaseCacheConfig{
			Client:             proxyClient,
			BaseContext:        ctx,
			Proxier:            apiProxy,
			Logger:             cacheLogger.Named("leasecache"),
			CacheStaticSecrets: config.Cache.CacheStaticSecrets,
			// dynamic secrets are configured as default-on to preserve backwards compatibility
			CacheDynamicSecrets: !config.Cache.DisableCachingDynamicSecrets,
			UserAgentToUse:      useragent.AgentProxyString(),
		})
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error creating lease cache: %v", err))
			return 1
		}

		cacheLogger.Info("cache configured", "cache_static_secrets", config.Cache.CacheStaticSecrets, "disable_caching_dynamic_secrets", config.Cache.DisableCachingDynamicSecrets)

		// Configure persistent storage and add to LeaseCache
		if config.Cache.Persist != nil {
			deferFunc, oldToken, err := agentproxyshared.AddPersistentStorageToLeaseCache(ctx, leaseCache, config.Cache.Persist, cacheLogger)
			if err != nil {
				c.UI.Error(fmt.Sprintf("Error creating persistent cache: %v", err))
				return 1
			}
			previousToken = oldToken
			if deferFunc != nil {
				defer deferFunc()
			}
		}

		// If we're caching static secrets, we need to start the updater, too
		if config.Cache.CacheStaticSecrets {
			staticSecretCacheUpdaterLogger := c.logger.Named("cache.staticsecretcacheupdater")
			inmemSink, err := inmem.New(&sink.SinkConfig{
				Logger: staticSecretCacheUpdaterLogger,
			}, leaseCache)
			if err != nil {
				c.UI.Error(fmt.Sprintf("Error creating inmem sink for static secret updater susbsystem: %v", err))
				return 1
			}
			sinks = append(sinks, &sink.SinkConfig{
				Logger: staticSecretCacheUpdaterLogger,
				Sink:   inmemSink,
			})

			updater, err = cache.NewStaticSecretCacheUpdater(&cache.StaticSecretCacheUpdaterConfig{
				Client:     client,
				LeaseCache: leaseCache,
				Logger:     staticSecretCacheUpdaterLogger,
				TokenSink:  inmemSink,
			})
			if err != nil {
				c.UI.Error(fmt.Sprintf("Error creating static secret cache updater: %v", err))
				return 1
			}

			capabilityManager, err := cache.NewStaticSecretCapabilityManager(&cache.StaticSecretCapabilityManagerConfig{
				LeaseCache: leaseCache,
				Logger:     c.logger.Named("cache.staticsecretcapabilitymanager"),
				Client:     client,
				StaticSecretTokenCapabilityRefreshInterval: config.Cache.StaticSecretTokenCapabilityRefreshInterval,
			})
			if err != nil {
				c.UI.Error(fmt.Sprintf("Error creating static secret capability manager: %v", err))
				return 1
			}
			leaseCache.SetCapabilityManager(capabilityManager)
		}
	}

	var listeners []net.Listener

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
				apiProxyLogger.Debug("configuring inmem auto-auth sink")
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
		quitEnabled := lnConfig.ProxyAPI != nil && lnConfig.ProxyAPI.EnableQuit

		mux.Handle(consts.ProxyPathMetrics, c.handleMetrics())
		if "metrics_only" != lnConfig.Role {
			mux.Handle(consts.ProxyPathCacheClear, leaseCache.HandleCacheClear(ctx))
			mux.Handle(consts.ProxyPathQuit, c.handleQuit(quitEnabled))
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
				c.UI.Output("==> Vault Proxy config reload triggered")
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
				c.UI.Output("==> Vault Proxy shutdown triggered")
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
		// Auth Handler is going to set its own retry values, so we want to
		// work on a copy of the client to not affect other subsystems.
		ahClient, err := c.client.CloneWithHeaders()
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error cloning client for auth handler: %v", err))
			return 1
		}

		// Override the set namespace with the auto-auth specific namespace
		if !namespaceSetByEnvironmentVariable && config.AutoAuth.Method.Namespace != "" {
			ahClient.SetNamespace(config.AutoAuth.Method.Namespace)
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
			Token:                        previousToken,
			ExitOnError:                  config.AutoAuth.Method.ExitOnError,
			UserAgent:                    useragent.ProxyAutoAuthString(),
			MetricsSignifier:             "proxy",
		})

		ss := sink.NewSinkServer(&sink.SinkServerConfig{
			Logger:        c.logger.Named("sink.server"),
			Client:        ahClient,
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

			return err
		}, func(error) {
			// Let the lease cache know this is a shutdown; no need to evict
			// everything
			if leaseCache != nil {
				leaseCache.SetShuttingDown(true)
			}
			cancelFunc()
		})
	}

	// Add the static secret cache updater, if appropriate
	if updater != nil {
		g.Add(func() error {
			err := updater.Run(ctx)
			return err
		}, func(error) {
			cancelFunc()
		})
	}

	// Server configuration output
	padding := 24
	sort.Strings(infoKeys)
	caser := cases.Title(language.English)
	c.UI.Output("==> Vault Proxy configuration:\n")
	for _, k := range infoKeys {
		c.UI.Output(fmt.Sprintf(
			"%s%s: %s",
			strings.Repeat(" ", padding-len(k)),
			caser.String(k),
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
func (c *ProxyCommand) applyConfigOverrides(f *FlagSets, config *proxyConfig.Config) {
	if config.Vault == nil {
		config.Vault = &proxyConfig.Vault{}
	}

	f.applyLogConfigOverrides(config.SharedConfig)

	f.Visit(func(fl *flag.Flag) {
		if fl.Name == flagNameProxyExitAfterAuth {
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

func (c *ProxyCommand) notifySystemd(status string) {
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

func (c *ProxyCommand) setStringFlag(f *FlagSets, configVal string, fVar *StringVar) {
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

func (c *ProxyCommand) setBoolFlag(f *FlagSets, configVal bool, fVar *BoolVar) {
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
func (c *ProxyCommand) storePidFile(pidPath string) error {
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
func (c *ProxyCommand) removePidFile(pidPath string) error {
	if pidPath == "" {
		return nil
	}
	return os.Remove(pidPath)
}

func (c *ProxyCommand) handleMetrics() http.Handler {
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
			w.WriteHeader(status)
			w.Write([]byte(v))
		case []byte:
			w.WriteHeader(status)
			w.Write(v)
		default:
			logical.RespondError(w, http.StatusInternalServerError, fmt.Errorf("wrong response returned"))
		}
	})
}

func (c *ProxyCommand) handleQuit(enabled bool) http.Handler {
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

// newLogger creates a logger based on parsed config field on the Proxy Command struct.
func (c *ProxyCommand) newLogger() (log.InterceptLogger, error) {
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

	logCfg, err := logging.NewLogConfig(nameProxy)
	if err != nil {
		return nil, err
	}
	logCfg.Name = nameProxy
	logCfg.LogLevel = logLevel
	logCfg.LogFormat = logFormat
	logCfg.LogFilePath = c.config.LogFile
	logCfg.LogRotateDuration = logRotateDuration
	logCfg.LogRotateBytes = c.config.LogRotateBytes
	logCfg.LogRotateMaxFiles = c.config.LogRotateMaxFiles

	l, err := logging.Setup(logCfg, c.logWriter)
	if err != nil {
		return nil, err
	}

	return l, nil
}

// loadConfig attempts to generate a Proxy config from the file(s) specified.
func (c *ProxyCommand) loadConfig(paths []string) (*proxyConfig.Config, error) {
	var errors error
	cfg := proxyConfig.NewConfig()

	for _, configPath := range paths {
		configFromPath, err := proxyConfig.LoadConfig(configPath)
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
// config values without requiring a restart of the Vault Proxy.
// If config is retrieved without error it is stored in the config field of the ProxyCommand.
// This operation is not atomic and could result in updated config but partially applied config settings.
// The error returned from this func may be a multierror.
// This function will most likely be called due to Vault Proxy receiving a SIGHUP signal.
// Currently only reloading the following are supported:
// * log level
// * TLS certs for listeners
func (c *ProxyCommand) reloadConfig(paths []string) error {
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
// to the ProxyCommand struct using the value currently set in config.
func (c *ProxyCommand) reloadLogLevel() error {
	logLevel, err := logging.ParseLogLevel(c.config.LogLevel)
	if err != nil {
		return err
	}

	c.logger.SetLevel(logLevel)

	return nil
}

// reloadCerts will attempt to reload certificates using a reload func which
// was provided when the listeners were configured, only funcs that were appended
// to the ProxyCommand slice will be invoked.
// This function returns a multierror type so that every func can report an error
// if it encounters one.
func (c *ProxyCommand) reloadCerts() error {
	var errors error

	c.tlsReloadFuncsLock.RLock()
	defer c.tlsReloadFuncsLock.RUnlock()

	for _, reloadFunc := range c.tlsReloadFuncs {
		// Non-TLS listeners will have a nil reload func.
		if reloadFunc != nil {
			err := reloadFunc()
			if err != nil {
				errors = multierror.Append(errors, err)
			}
		}
	}

	return errors
}

// outputErrors will take an error or multierror and handle outputting each to the UI
func (c *ProxyCommand) outputErrors(err error) {
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
