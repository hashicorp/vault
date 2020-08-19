package agent

import (
	"context"
	"github.com/hashicorp/consul-template/config"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/auth"
	config2 "github.com/hashicorp/vault/command/agent/config"
	"github.com/hashicorp/vault/command/agent/sink"
	"github.com/hashicorp/vault/command/agent/template"
	"github.com/oklog/run"
	"io"
	"time"
)

type AgentServerConfig struct {
	Logger                               hclog.Logger
	Level                                hclog.Level
	Writer                               io.Writer
	Sinks                                []*sink.SinkConfig
	Templates                            []*config.TemplateConfig
	Namespace                            string
	ExitAfterAuth                        bool
	AutoAuthWrapTTL                      time.Duration
	AutoAuthEnableReauthOnNewCredentials bool
	VaultConf                            *config2.Vault
}

func RunAgent(ctx context.Context, cancelFunc context.CancelFunc, group *run.Group, cfg *AgentServerConfig, client *api.Client, method auth.AuthMethod) {
	ah := auth.NewAuthHandler(&auth.AuthHandlerConfig{
		Logger:                       cfg.Logger.Named("auth.handler"),
		Client:                       client,
		WrapTTL:                      cfg.AutoAuthWrapTTL,
		EnableReauthOnNewCredentials: cfg.AutoAuthEnableReauthOnNewCredentials,
		EnableTemplateTokenCh:        len(cfg.Templates) > 0,
	})

	ss := sink.NewSinkServer(&sink.SinkServerConfig{
		Logger:        cfg.Logger.Named("sink.server"),
		Client:        client,
		ExitAfterAuth: cfg.ExitAfterAuth,
	})

	ts := template.NewServer(&template.ServerConfig{
		Logger:        cfg.Logger.Named("template.server"),
		LogLevel:      cfg.Level,
		LogWriter:     cfg.Writer,
		VaultConf:     cfg.VaultConf,
		Namespace:     cfg.Namespace,
		ExitAfterAuth: cfg.ExitAfterAuth,
	})

	group.Add(func() error {
		return ah.Run(ctx, method)
	}, func(error) {
		cancelFunc()
	})

	group.Add(func() error {
		err := ss.Run(ctx, ah.OutputCh, cfg.Sinks)
		cfg.Logger.Info("sinks finished, exiting")

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
		if len(cfg.Templates) > 0 {
			<-ts.DoneCh
		}

		return err
	}, func(error) {
		cancelFunc()
	})

	group.Add(func() error {
		return ts.Run(ctx, ah.TemplateTokenCh, cfg.Templates)
	}, func(error) {
		cancelFunc()
		ts.Stop()
	})
}
