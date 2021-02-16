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
	TemplateRetry *config.TemplateRetry
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

	// testingLimitRetry is used for tests to limit the number of retries
	// performed by the template server
	testingLimitRetry int
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
						Token: latestToken,
					},
				}

				if ts.config.TemplateRetry != nil && ts.config.TemplateRetry.Enabled {
					ctv.Vault.Retry = &ctconfig.RetryConfig{
						Attempts:   &ts.config.TemplateRetry.Attempts,
						Backoff:    &ts.config.TemplateRetry.Backoff,
						MaxBackoff: &ts.config.TemplateRetry.MaxBackoff,
						Enabled:    &ts.config.TemplateRetry.Enabled,
					}
				} else if ts.testingLimitRetry != 0 {
					// If we're testing, limit retries to 3 attempts to avoid
					// long test runs from exponential back-offs
					ctv.Vault.Retry = &ctconfig.RetryConfig{Attempts: &ts.testingLimitRetry}
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
			ts.runner.StopImmediately()
			return fmt.Errorf("template server: %w", err)

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

	if sc.AgentConfig.Cache != nil && len(sc.AgentConfig.Listeners) != 0 {
		scheme := "unix:/"
		if sc.AgentConfig.Listeners[0].Type == "tcp" {
			scheme = "https://"
			if sc.AgentConfig.Listeners[0].TLSDisable {
				scheme = "http://"
			}
		}
		address := fmt.Sprintf("%s%s", scheme, sc.AgentConfig.Listeners[0].Address)
		conf.Vault.Address = &address
	}

	if sc.Namespace != "" {
		conf.Vault.Namespace = &sc.Namespace
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

	if strings.HasPrefix(*conf.Vault.Address, "https") || sc.AgentConfig.Vault.CACert != "" {
		skipVerify := sc.AgentConfig.Vault.TLSSkipVerify
		verify := !skipVerify
		conf.Vault.SSL = &ctconfig.SSLConfig{
			Enabled: pointerutil.BoolPtr(true),
			Verify:  &verify,
			Cert:    &sc.AgentConfig.Vault.ClientCert,
			Key:     &sc.AgentConfig.Vault.ClientKey,
			CaCert:  &sc.AgentConfig.Vault.CACert,
			CaPath:  &sc.AgentConfig.Vault.CAPath,
		}

		// Only configure TLS Skip Verify if CT is not going through the cache. We can
		// skip verification if its using the cache because they're part of the same agent.
		// Agent listener doesn't support mTLS listeners.
		if sc.AgentConfig.Cache != nil {
			conf.Vault.SSL.Enabled = pointerutil.BoolPtr(true)
			conf.Vault.SSL.Verify = pointerutil.BoolPtr(false)
			conf.Vault.SSL.Cert = pointerutil.StringPtr("")
			conf.Vault.SSL.Key = pointerutil.StringPtr("")
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
		levelStr = "ERROR"
	default:
		levelStr = "INFO"
	}
	return pointerutil.StringPtr(levelStr)
}
