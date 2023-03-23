// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package template is responsible for rendering user supplied templates to
// disk. The Server type accepts configuration to communicate to a Vault server
// and a Vault token for authentication. Internally, the Server creates a Consul
// Template Runner which manages reading secrets from Vault and rendering
// templates to disk at configured locations
package template

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/vault/helper/useragent"

	"go.uber.org/atomic"

	ctconfig "github.com/hashicorp/consul-template/config"
	ctlogging "github.com/hashicorp/consul-template/logging"
	"github.com/hashicorp/consul-template/manager"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agent/config"
	"github.com/hashicorp/vault/sdk/helper/pointerutil"
)

// ServerConfig is a config struct for setting up the basic parts of the
// Server
type ServerConfig struct {
	Logger hclog.Logger
	// Client        *api.Client
	AgentConfig *config.Config

	ExitAfterAuth bool
	Namespace     string

	// LogLevel is needed to set the internal Consul Template Runner's log level
	// to match the log level of Vault Agent. The internal Runner creates it's own
	// logger and can't be set externally or copied from the Template Server.
	//
	// LogWriter is needed to initialize Consul Template's internal logger to use
	// the same io.Writer that Vault Agent itself is using.
	LogLevel  hclog.Level
	LogWriter io.Writer
}

// Server manages the Consul Template Runner which renders templates
type Server struct {
	// config holds the ServerConfig used to create it. It's passed along in other
	// methods
	config *ServerConfig

	// runner is the consul-template runner
	runner        *manager.Runner
	runnerStarted *atomic.Bool

	// Templates holds the parsed Consul Templates
	Templates []*ctconfig.TemplateConfig

	// lookupMap is a list of templates indexed by their consul-template ID. This
	// is used to ensure all Vault templates have been rendered before returning
	// from the runner in the event we're using exit after auth.
	lookupMap map[string][]*ctconfig.TemplateConfig

	DoneCh  chan struct{}
	stopped *atomic.Bool

	logger        hclog.Logger
	exitAfterAuth bool
}

// NewServer returns a new configured server
func NewServer(conf *ServerConfig) *Server {
	ts := Server{
		DoneCh:        make(chan struct{}),
		stopped:       atomic.NewBool(false),
		runnerStarted: atomic.NewBool(false),

		logger:        conf.Logger,
		config:        conf,
		exitAfterAuth: conf.ExitAfterAuth,
	}
	return &ts
}

// Run kicks off the internal Consul Template runner, and listens for changes to
// the token from the AuthHandler. If Done() is called on the context, shut down
// the Runner and return
func (ts *Server) Run(ctx context.Context, incoming chan string, templates []*ctconfig.TemplateConfig) error {
	if incoming == nil {
		return errors.New("template server: incoming channel is nil")
	}

	latestToken := new(string)
	ts.logger.Info("starting template server")

	defer func() {
		ts.logger.Info("template server stopped")
	}()

	// If there are no templates, we wait for context cancellation and then return
	if len(templates) == 0 {
		ts.logger.Info("no templates found")
		<-ctx.Done()
		return nil
	}

	// construct a consul template vault config based the agents vault
	// configuration
	var runnerConfig *ctconfig.Config
	var runnerConfigErr error

	if runnerConfig, runnerConfigErr = newRunnerConfig(ts.config, templates); runnerConfigErr != nil {
		return fmt.Errorf("template server failed to runner generate config: %w", runnerConfigErr)
	}

	var err error
	ts.runner, err = manager.NewRunner(runnerConfig, false)
	if err != nil {
		return fmt.Errorf("template server failed to create: %w", err)
	}

	// Build the lookup map using the id mapping from the Template runner. This is
	// used to check the template rendering against the expected templates. This
	// returns a map with a generated ID and a slice of templates for that id. The
	// slice is determined by the source or contents of the template, so if a
	// configuration has multiple templates specified, but are the same source /
	// contents, they will be identified by the same key.
	idMap := ts.runner.TemplateConfigMapping()
	lookupMap := make(map[string][]*ctconfig.TemplateConfig, len(idMap))
	for id, ctmpls := range idMap {
		for _, ctmpl := range ctmpls {
			tl := lookupMap[id]
			tl = append(tl, ctmpl)
			lookupMap[id] = tl
		}
	}
	ts.lookupMap = lookupMap

	for {
		select {
		case <-ctx.Done():
			ts.runner.Stop()
			return nil

		case token := <-incoming:
			if token != *latestToken {
				ts.logger.Info("template server received new token")

				// If the runner was previously started and we intend to exit
				// after auth, do not restart the runner if a new token is
				// received.
				if ts.exitAfterAuth && ts.runnerStarted.Load() {
					ts.logger.Info("template server not restarting with new token with exit_after_auth set to true")
					continue
				}

				ts.runner.Stop()
				*latestToken = token
				ctv := ctconfig.Config{
					Vault: &ctconfig.VaultConfig{
						Token:           latestToken,
						ClientUserAgent: pointerutil.StringPtr(useragent.AgentTemplatingString()),
					},
				}

				runnerConfig = runnerConfig.Merge(&ctv)
				var runnerErr error
				ts.runner, runnerErr = manager.NewRunner(runnerConfig, false)
				if runnerErr != nil {
					ts.logger.Error("template server failed with new Vault token", "error", runnerErr)
					continue
				}
				ts.runnerStarted.CAS(false, true)
				go ts.runner.Start()
			}

		case err := <-ts.runner.ErrCh:
			ts.logger.Error("template server error", "error", err.Error())
			ts.runner.StopImmediately()

			// Return after stopping the runner if exit on retry failure was
			// specified
			if ts.config.AgentConfig.TemplateConfig != nil && ts.config.AgentConfig.TemplateConfig.ExitOnRetryFailure {
				return fmt.Errorf("template server: %w", err)
			}

			ts.runner, err = manager.NewRunner(runnerConfig, false)
			if err != nil {
				return fmt.Errorf("template server failed to create: %w", err)
			}
			go ts.runner.Start()

		case <-ts.runner.TemplateRenderedCh():
			// A template has been rendered, figure out what to do
			events := ts.runner.RenderEvents()

			// events are keyed by template ID, and can be matched up to the id's from
			// the lookupMap
			if len(events) < len(ts.lookupMap) {
				// Not all templates have been rendered yet
				continue
			}

			// assume the renders are finished, until we find otherwise
			doneRendering := true
			for _, event := range events {
				// This template hasn't been rendered
				if event.LastWouldRender.IsZero() {
					doneRendering = false
				}
			}

			if doneRendering && ts.exitAfterAuth {
				// if we want to exit after auth, go ahead and shut down the runner and
				// return. The deferred closing of the DoneCh will allow agent to
				// continue with closing down
				ts.runner.Stop()
				return nil
			}
		}
	}
}

func (ts *Server) Stop() {
	if ts.stopped.CAS(false, true) {
		close(ts.DoneCh)
	}
}

// newRunnerConfig returns a consul-template runner configuration, setting the
// Vault and Consul configurations based on the clients configs.
func newRunnerConfig(sc *ServerConfig, templates ctconfig.TemplateConfigs) (*ctconfig.Config, error) {
	conf := ctconfig.DefaultConfig()
	conf.Templates = templates.Copy()

	// Setup the Vault config
	// Always set these to ensure nothing is picked up from the environment
	conf.Vault.RenewToken = pointerutil.BoolPtr(false)
	conf.Vault.Token = pointerutil.StringPtr("")
	conf.Vault.Address = &sc.AgentConfig.Vault.Address
	conf.Vault.ClientUserAgent = pointerutil.StringPtr(useragent.AgentTemplatingString())

	if sc.Namespace != "" {
		conf.Vault.Namespace = &sc.Namespace
	}

	if sc.AgentConfig.TemplateConfig != nil && sc.AgentConfig.TemplateConfig.StaticSecretRenderInt != 0 {
		conf.Vault.DefaultLeaseDuration = &sc.AgentConfig.TemplateConfig.StaticSecretRenderInt
	}

	if sc.AgentConfig.DisableIdleConnsTemplating {
		idleConns := -1
		conf.Vault.Transport.MaxIdleConns = &idleConns
	}

	if sc.AgentConfig.DisableKeepAlivesTemplating {
		conf.Vault.Transport.DisableKeepAlives = pointerutil.BoolPtr(true)
	}

	conf.Vault.SSL = &ctconfig.SSLConfig{
		Enabled:    pointerutil.BoolPtr(false),
		Verify:     pointerutil.BoolPtr(false),
		Cert:       pointerutil.StringPtr(""),
		Key:        pointerutil.StringPtr(""),
		CaCert:     pointerutil.StringPtr(""),
		CaPath:     pointerutil.StringPtr(""),
		ServerName: pointerutil.StringPtr(""),
	}

	// If Vault.Retry isn't specified, use the default of 12 retries.
	// This retry value will be respected regardless of if we use the cache.
	attempts := ctconfig.DefaultRetryAttempts
	if sc.AgentConfig.Vault != nil && sc.AgentConfig.Vault.Retry != nil {
		attempts = sc.AgentConfig.Vault.Retry.NumRetries
	}

	// Use the cache if available or fallback to the Vault server values.
	if sc.AgentConfig.Cache != nil {
		if sc.AgentConfig.Cache.InProcDialer == nil {
			return nil, fmt.Errorf("missing in-process dialer configuration")
		}
		if conf.Vault.Transport == nil {
			conf.Vault.Transport = &ctconfig.TransportConfig{}
		}
		conf.Vault.Transport.CustomDialer = sc.AgentConfig.Cache.InProcDialer
		// The in-process dialer ignores the address passed in, but we're still
		// setting it here to override the setting at the top of this function,
		// and to prevent the vault/http client from defaulting to https.
		conf.Vault.Address = pointerutil.StringPtr("http://127.0.0.1:8200")
	} else if strings.HasPrefix(sc.AgentConfig.Vault.Address, "https") || sc.AgentConfig.Vault.CACert != "" {
		skipVerify := sc.AgentConfig.Vault.TLSSkipVerify
		verify := !skipVerify
		conf.Vault.SSL = &ctconfig.SSLConfig{
			Enabled:    pointerutil.BoolPtr(true),
			Verify:     &verify,
			Cert:       &sc.AgentConfig.Vault.ClientCert,
			Key:        &sc.AgentConfig.Vault.ClientKey,
			CaCert:     &sc.AgentConfig.Vault.CACert,
			CaPath:     &sc.AgentConfig.Vault.CAPath,
			ServerName: &sc.AgentConfig.Vault.TLSServerName,
		}
	}
	enabled := attempts > 0
	conf.Vault.Retry = &ctconfig.RetryConfig{
		Attempts: &attempts,
		Enabled:  &enabled,
	}

	// Sync Consul Template's retry with user set auto-auth initial backoff value.
	// This is helpful if Auto Auth cannot get a new token and CT is trying to fetch
	// secrets.
	if sc.AgentConfig.AutoAuth != nil && sc.AgentConfig.AutoAuth.Method != nil {
		if sc.AgentConfig.AutoAuth.Method.MinBackoff > 0 {
			conf.Vault.Retry.Backoff = &sc.AgentConfig.AutoAuth.Method.MinBackoff
		}

		if sc.AgentConfig.AutoAuth.Method.MaxBackoff > 0 {
			conf.Vault.Retry.MaxBackoff = &sc.AgentConfig.AutoAuth.Method.MaxBackoff
		}
	}

	conf.Finalize()

	// setup log level from TemplateServer config
	conf.LogLevel = logLevelToStringPtr(sc.LogLevel)

	if err := ctlogging.Setup(&ctlogging.Config{
		Level:  *conf.LogLevel,
		Writer: sc.LogWriter,
	}); err != nil {
		return nil, err
	}
	return conf, nil
}

// logLevelToString converts a go-hclog level to a matching, uppercase string
// value. It's used to convert Vault Agent's hclog level to a string version
// suitable for use in Consul Template's runner configuration input.
func logLevelToStringPtr(level hclog.Level) *string {
	// consul template's default level is WARN, but Vault Agent's default is INFO,
	// so we use that for the Runner's default.
	var levelStr string

	switch level {
	case hclog.Trace:
		levelStr = "TRACE"
	case hclog.Debug:
		levelStr = "DEBUG"
	case hclog.Warn:
		levelStr = "WARN"
	case hclog.Error:
		levelStr = "ERR"
	default:
		levelStr = "INFO"
	}
	return pointerutil.StringPtr(levelStr)
}
