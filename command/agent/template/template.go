// Package template is responsible for rendering user supplied templates to
// disk. The Server type accepts configuration to communicate to a Vault server
// and a Vault token for authentication. Internally, the Server creates a Consul
// Template Runner which manages reading secrets from Vault and rendering
// templates to disk at configured locations
package template

import (
	"context"
	"strings"

	ctconfig "github.com/hashicorp/consul-template/config"
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
	VaultConf     *config.Vault
	ExitAfterAuth bool

	Namespace string
}

// Server manages the Consul Template Runner which renders templates
type Server struct {
	// UnblockCh is used to block until a template is rendered
	UnblockCh chan struct{}

	// config holds the ServerConfig used to create it. It's passed along in other
	// methods
	config *ServerConfig

	// runner is the consul-template runner
	runner *manager.Runner

	// Templates holds the parsed Consul Templates
	Templates []*ctconfig.TemplateConfig

	// TODO: remove donech?
	DoneCh        chan struct{}
	logger        hclog.Logger
	exitAfterAuth bool
}

// NewServer returns a new configured server
func NewServer(conf *ServerConfig) *Server {
	ts := Server{
		DoneCh:        make(chan struct{}),
		logger:        conf.Logger,
		UnblockCh:     make(chan struct{}),
		config:        conf,
		exitAfterAuth: conf.ExitAfterAuth,
	}
	return &ts
}

// Run kicks off the internal Consul Template runner, and listens for changes to
// the token from the AuthHandler. If Done() is called on the context, shut down
// the Runner and return
func (ts *Server) Run(ctx context.Context, incoming chan string, templates []*ctconfig.TemplateConfig) {
	latestToken := new(string)
	ts.logger.Info("starting template server")
	defer func() {
		ts.logger.Info("template server stopped")
		close(ts.DoneCh)
	}()

	if incoming == nil {
		panic("incoming channel is nil")
	}

	// If there are no templates, close the UnblockCh
	if len(templates) == 0 {
		// nothing to do
		ts.logger.Info("no templates found")
		close(ts.UnblockCh)
		return
	}

	// construct a consul template vault config based the agents vault
	// configuration
	var runnerConfig *ctconfig.Config
	if runnerConfig = newRunnerConfig(ts.config, templates); runnerConfig == nil {
		ts.logger.Error("template server failed to generate runner config")
		close(ts.UnblockCh)
		return
	}

	var err error
	ts.runner, err = manager.NewRunner(runnerConfig, false)
	if err != nil {
		ts.logger.Error("template server failed to create", "error", err)
		close(ts.UnblockCh)
		return
	}

	for {
		select {
		case <-ctx.Done():
			ts.runner.StopImmediately()
			ts.runner = nil
			return

		case token := <-incoming:
			if token != *latestToken {
				ts.logger.Info("template server received new token")
				ts.runner.Stop()
				*latestToken = token
				ctv := ctconfig.Config{
					Vault: &ctconfig.VaultConfig{
						Token: latestToken,
					},
				}
				runnerConfig.Merge(&ctv)
				runnerConfig.Finalize()
				var runnerErr error
				ts.runner, runnerErr = manager.NewRunner(runnerConfig, false)
				if runnerErr != nil {
					ts.logger.Error("template server failed with new Vault token", "error", runnerErr)
					continue
				} else {
					go ts.runner.Start()
				}
			}
		case err := <-ts.runner.ErrCh:
			ts.logger.Error("template server error", "error", err.Error())
			close(ts.UnblockCh)
			return
		case <-ts.runner.TemplateRenderedCh():
			// A template has been rendered, unblock
			if ts.exitAfterAuth {
				// if we want to exit after auth, go ahead and shut down the runner
				ts.runner.Stop()
			}
			close(ts.UnblockCh)
		}
	}
}

// newRunnerConfig returns a consul-template runner configuration, setting the
// Vault and Consul configurations based on the clients configs.
func newRunnerConfig(sc *ServerConfig, templates ctconfig.TemplateConfigs) *ctconfig.Config {
	// TODO only use default Vault config
	conf := ctconfig.DefaultConfig()
	conf.Templates = templates.Copy()

	// Setup the Vault config
	// Always set these to ensure nothing is picked up from the environment
	conf.Vault.RenewToken = pointerutil.BoolPtr(false)
	conf.Vault.Token = pointerutil.StringPtr("")
	conf.Vault.Address = &sc.VaultConf.Address

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

	if strings.HasPrefix(sc.VaultConf.Address, "https") || sc.VaultConf.CACert != "" {
		skipVerify := sc.VaultConf.TLSSkipVerify
		verify := !skipVerify
		conf.Vault.SSL = &ctconfig.SSLConfig{
			Enabled: pointerutil.BoolPtr(true),
			Verify:  &verify,
			Cert:    &sc.VaultConf.ClientCert,
			Key:     &sc.VaultConf.ClientKey,
			CaCert:  &sc.VaultConf.CACert,
			CaPath:  &sc.VaultConf.CAPath,
		}
	}

	conf.Finalize()
	return conf
}
