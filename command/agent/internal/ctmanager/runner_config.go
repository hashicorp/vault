package ctmanager

import (
	"io"
	"fmt"
	"strings"

	ctconfig "github.com/hashicorp/consul-template/config"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/pointerutil"
	ctlogging "github.com/hashicorp/consul-template/logging"
	"github.com/hashicorp/vault/command/agent/config"
)

type ManagerConfig struct {
	AgentConfig *config.Config
	Namespace   string
	LogLevel    hclog.Level
	LogWriter   io.Writer
}

// NewManagerConfig returns a consul-template runner configuration, setting the
// Vault and Consul configurations based on the clients configs.
func NewManagerConfig(sc ManagerConfig, templates ctconfig.TemplateConfigs) (*ctconfig.Config, error) {
	conf := ctconfig.DefaultConfig()
	conf.Templates = templates.Copy()

	// Setup the Vault config
	// Always set these to ensure nothing is picked up from the environment
	conf.Vault.RenewToken = pointerutil.BoolPtr(false)
	conf.Vault.Token = pointerutil.StringPtr("")
	conf.Vault.Address = &sc.AgentConfig.Vault.Address

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
