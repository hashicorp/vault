// Package template is responsible for rendering user supplied templates to
// disk. The Server type managing the lifecycle of an internal Consul Template
// Runner
package template

import (
	"context"
	"strings"

	ctconfig "github.com/hashicorp/consul-template/config"
	"github.com/hashicorp/consul-template/manager"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agent/config"
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
	// config holds the ServerConfig used to create it. It's passed along in other
	// methods
	config *ServerConfig

	// // lookup allows looking up the set of Nomad templates by their consul-template ID
	// lookup map[string][]*structs.Template

	// runner is the consul-template runner
	runner *manager.Runner

	// shutdownCh is used to signal the started goroutine to shutdown
	// shutdownCh chan struct{}

	// shutdown marks whether the manager has been shutdown
	// shutdown     bool
	// shutdownLock sync.Mutex
	Templates []*ctconfig.TemplateConfig

	DoneCh chan struct{}
	logger hclog.Logger
	// client        *api.Client
	exitAfterAuth bool
}

// NewServer returns a new configured server
func NewServer(conf *ServerConfig) *Server {
	ts := Server{
		DoneCh:        make(chan struct{}),
		logger:        conf.Logger,
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
	ts.logger.Info("starting sink server")
	defer func() {
		ts.logger.Info("template server stopped")
		close(ts.DoneCh)
	}()

	if incoming == nil {
		panic("incoming channel is nil")
	}

	// construct a consul template vault config based the agents vault
	// configuration
	var runnerConfig *ctconfig.Config
	if runnerConfig = newRunnerConfig(ts.config, ts.Templates); runnerConfig == nil {
		ts.logger.Info("template server failed to generate runner config")
		close(ts.DoneCh)
		return
	}

	runner, err := manager.NewRunner(runnerConfig, false)
	if err != nil {
		ts.logger.Info("template server failed to create")
		close(ts.DoneCh)
		return
	}

	for {
		select {
		case <-ctx.Done():
			ts.runner.StopImmediately()
			close(ts.DoneCh)
			return

		case token := <-incoming:
			// q.Q(">> incoming token")
			if token != *latestToken {
				// q.Q(">>:: new token")
				ts.runner.Stop()
				*latestToken = token
				ctv := ctconfig.Config{
					Vault: &ctconfig.VaultConfig{
						Token: latestToken,
					},
				}
				runnerConfig.Merge(&ctv)
				var runnerErr error
				runner, runnerErr = manager.NewRunner(runnerConfig, false)
				if runnerErr != nil {
					ts.logger.Info("template server failed with new Vault token")
					close(ts.DoneCh)
					return
				}
				go ts.runner.Start()
			}
		case err := <-runner.ErrCh:
			ts.logger.Info("template server error:", err)
			close(ts.DoneCh)
			return
		}
	}
}

// newRunnerConfig returns a consul-template runner configuration, setting the
// Vault and Consul configurations based on the clients configs.
func newRunnerConfig(sc *ServerConfig, templates ctconfig.TemplateConfigs) *ctconfig.Config {
	conf := ctconfig.DefaultConfig()
	conf.Templates = templates.Copy()

	// Setup the Vault config
	// Always set these to ensure nothing is picked up from the environment
	emptyStr := ""
	conf.Vault.RenewToken = boolPtr(false)
	// TODO: need token here
	conf.Vault.Token = &emptyStr
	conf.Vault.Address = &sc.VaultConf.Address
	// conf.Vault.Token = &config.VaultToken
	if sc.Namespace != "" {
		conf.Vault.Namespace = &sc.Namespace
	}

	if strings.HasPrefix(sc.VaultConf.Address, "https") || sc.VaultConf.CACert != "" {
		skipVerify := sc.VaultConf.TLSSkipVerify
		verify := !skipVerify
		conf.Vault.SSL = &ctconfig.SSLConfig{
			Enabled: boolPtr(true),
			Verify:  &verify,
			Cert:    &sc.VaultConf.ClientCert,
			Key:     &sc.VaultConf.ClientKey,
			CaCert:  &sc.VaultConf.CACert,
			CaPath:  &sc.VaultConf.CAPath,
			// ServerName: &sc.VaultConf.TLSServerName,
		}
	} else {
		conf.Vault.SSL = &ctconfig.SSLConfig{
			Enabled:    boolPtr(false),
			Verify:     boolPtr(false),
			Cert:       &emptyStr,
			Key:        &emptyStr,
			CaCert:     &emptyStr,
			CaPath:     &emptyStr,
			ServerName: &emptyStr,
		}
	}

	conf.Finalize()
	return conf
}

func boolPtr(b bool) *bool {
	return &b
}
