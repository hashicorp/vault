// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package auth

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/backoff"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
)

const (
	defaultMinBackoff = 1 * time.Second
	defaultMaxBackoff = 5 * time.Minute
)

// AuthMethod is the interface that auto-auth methods implement for the agent/proxy
// to use.
type AuthMethod interface {
	// Authenticate returns a mount path, header, request body, and error.
	// The header may be nil if no special header is needed.
	Authenticate(context.Context, *api.Client) (string, http.Header, map[string]interface{}, error)
	NewCreds() chan struct{}
	CredSuccess()
	Shutdown()
}

// AuthMethodWithClient is an extended interface that can return an API client
// for use during the authentication call.
type AuthMethodWithClient interface {
	AuthMethod
	AuthClient(client *api.Client) (*api.Client, error)
}

type AuthConfig struct {
	Logger    hclog.Logger
	MountPath string
	WrapTTL   time.Duration
	Config    map[string]interface{}
}

// AuthHandler is responsible for keeping a token alive and renewed and passing
// new tokens to the sink server
type AuthHandler struct {
	OutputCh                     chan string
	TemplateTokenCh              chan string
	ExecTokenCh                  chan string
	token                        string
	userAgent                    string
	metricsSignifier             string
	logger                       hclog.Logger
	client                       *api.Client
	random                       *rand.Rand
	wrapTTL                      time.Duration
	maxBackoff                   time.Duration
	minBackoff                   time.Duration
	enableReauthOnNewCredentials bool
	enableTemplateTokenCh        bool
	enableExecTokenCh            bool
	exitOnError                  bool
}

type AuthHandlerConfig struct {
	Logger     hclog.Logger
	Client     *api.Client
	WrapTTL    time.Duration
	MaxBackoff time.Duration
	MinBackoff time.Duration
	Token      string
	// UserAgent is the HTTP UserAgent header auto-auth will use when
	// communicating with Vault.
	UserAgent string
	// MetricsSignifier is the first argument we will give to
	// metrics.IncrCounter, signifying what the name of the application is
	MetricsSignifier             string
	EnableReauthOnNewCredentials bool
	EnableTemplateTokenCh        bool
	EnableExecTokenCh            bool
	ExitOnError                  bool
}

func NewAuthHandler(conf *AuthHandlerConfig) *AuthHandler {
	ah := &AuthHandler{
		// This is buffered so that if we try to output after the sink server
		// has been shut down, during agent/proxy shutdown, we won't block
		OutputCh:                     make(chan string, 1),
		TemplateTokenCh:              make(chan string, 1),
		ExecTokenCh:                  make(chan string, 1),
		token:                        conf.Token,
		logger:                       conf.Logger,
		client:                       conf.Client,
		random:                       rand.New(rand.NewSource(int64(time.Now().Nanosecond()))),
		wrapTTL:                      conf.WrapTTL,
		minBackoff:                   conf.MinBackoff,
		maxBackoff:                   conf.MaxBackoff,
		enableReauthOnNewCredentials: conf.EnableReauthOnNewCredentials,
		enableTemplateTokenCh:        conf.EnableTemplateTokenCh,
		enableExecTokenCh:            conf.EnableExecTokenCh,
		exitOnError:                  conf.ExitOnError,
		userAgent:                    conf.UserAgent,
		metricsSignifier:             conf.MetricsSignifier,
	}

	return ah
}

func backoffSleep(ctx context.Context, backoff *autoAuthBackoff) bool {
	nextSleep, err := backoff.backoff.Next()
	if err != nil {
		return false
	}
	select {
	case <-time.After(nextSleep):
	case <-ctx.Done():
	}
	return true
}

func (ah *AuthHandler) Run(ctx context.Context, am AuthMethod) error {
	if am == nil {
		return errors.New("auth handler: nil auth method")
	}

	if ah.minBackoff <= 0 {
		ah.minBackoff = defaultMinBackoff
	}
	if ah.maxBackoff <= 0 {
		ah.maxBackoff = defaultMaxBackoff
	}
	if ah.minBackoff > ah.maxBackoff {
		return errors.New("auth handler: min_backoff cannot be greater than max_backoff")
	}
	backoffCfg := newAutoAuthBackoff(ah.minBackoff, ah.maxBackoff, ah.exitOnError)

	ah.logger.Info("starting auth handler")
	defer func() {
		am.Shutdown()
		close(ah.OutputCh)
		close(ah.TemplateTokenCh)
		close(ah.ExecTokenCh)
		ah.logger.Info("auth handler stopped")
	}()

	credCh := am.NewCreds()
	if !ah.enableReauthOnNewCredentials {
		realCredCh := credCh
		credCh = nil
		if realCredCh != nil {
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					case <-realCredCh:
					}
				}
			}()
		}
	}
	if credCh == nil {
		credCh = make(chan struct{})
	}

	if ah.client != nil {
		headers := ah.client.Headers()
		if headers == nil {
			headers = make(http.Header)
		}
		headers.Set("User-Agent", ah.userAgent)
		ah.client.SetHeaders(headers)
	}

	var watcher *api.LifetimeWatcher
	first := true

	for {
		select {
		case <-ctx.Done():
			return nil

		default:
		}

		var clientToUse *api.Client
		var err error
		var path string
		var data map[string]interface{}
		var header http.Header
		var isTokenFileMethod bool

		switch am.(type) {
		case AuthMethodWithClient:
			clientToUse, err = am.(AuthMethodWithClient).AuthClient(ah.client)
			if err != nil {
				ah.logger.Error("error creating client for authentication call", "error", err, "backoff", backoffCfg)
				metrics.IncrCounter([]string{ah.metricsSignifier, "auth", "failure"}, 1)

				if backoffSleep(ctx, backoffCfg) {
					continue
				}

				return err
			}
		default:
			clientToUse = ah.client
		}

		// Disable retry on the client to ensure our backoffOrQuit function is
		// the only source of retry/backoff.
		clientToUse.SetMaxRetries(0)

		var secret *api.Secret = new(api.Secret)
		if first && ah.token != "" {
			ah.logger.Debug("using preloaded token")

			first = false
			ah.logger.Debug("lookup-self with preloaded token")
			clientToUse.SetToken(ah.token)

			secret, err = clientToUse.Auth().Token().LookupSelfWithContext(ctx)
			if err != nil {
				ah.logger.Error("could not look up token", "err", err, "backoff", backoffCfg)
				metrics.IncrCounter([]string{ah.metricsSignifier, "auth", "failure"}, 1)

				if backoffSleep(ctx, backoffCfg) {
					continue
				}
				return err
			}

			duration, _ := secret.Data["ttl"].(json.Number).Int64()
			secret.Auth = &api.SecretAuth{
				ClientToken:   secret.Data["id"].(string),
				LeaseDuration: int(duration),
				Renewable:     secret.Data["renewable"].(bool),
			}
		} else {
			ah.logger.Info("authenticating")

			path, header, data, err = am.Authenticate(ctx, ah.client)
			if err != nil {
				ah.logger.Error("error getting path or data from method", "error", err, "backoff", backoffCfg)
				metrics.IncrCounter([]string{ah.metricsSignifier, "auth", "failure"}, 1)

				if backoffSleep(ctx, backoffCfg) {
					continue
				}
				return err
			}
		}

		if ah.wrapTTL > 0 {
			wrapClient, err := clientToUse.CloneWithHeaders()
			if err != nil {
				ah.logger.Error("error creating client for wrapped call", "error", err, "backoff", backoffCfg)
				metrics.IncrCounter([]string{ah.metricsSignifier, "auth", "failure"}, 1)

				if backoffSleep(ctx, backoffCfg) {
					continue
				}
				return err
			}
			wrapClient.SetWrappingLookupFunc(func(string, string) string {
				return ah.wrapTTL.String()
			})
			clientToUse = wrapClient
		}
		for key, values := range header {
			for _, value := range values {
				clientToUse.AddHeader(key, value)
			}
		}

		// This should only happen if there's no preloaded token (regular auto-auth login)
		// or if a preloaded token has expired and is now switching to auto-auth.
		if secret.Auth == nil {
			isTokenFileMethod = path == "auth/token/lookup-self"
			if isTokenFileMethod {
				token, _ := data["token"].(string)
				lookupSelfClient, err := clientToUse.CloneWithHeaders()
				if err != nil {
					ah.logger.Error("failed to clone client to perform token lookup")
					return err
				}
				lookupSelfClient.SetToken(token)
				secret, err = lookupSelfClient.Auth().Token().LookupSelf()
			} else {
				secret, err = clientToUse.Logical().WriteWithContext(ctx, path, data)
			}

			// Check errors/sanity
			if err != nil {
				ah.logger.Error("error authenticating", "error", err, "backoff", backoffCfg)
				metrics.IncrCounter([]string{ah.metricsSignifier, "auth", "failure"}, 1)

				if backoffSleep(ctx, backoffCfg) {
					continue
				}
				return err
			}
		}

		var leaseDuration int

		switch {
		case ah.wrapTTL > 0:
			if secret.WrapInfo == nil {
				ah.logger.Error("authentication returned nil wrap info", "backoff", backoffCfg)
				metrics.IncrCounter([]string{ah.metricsSignifier, "auth", "failure"}, 1)

				if backoffSleep(ctx, backoffCfg) {
					continue
				}
				return err
			}
			if secret.WrapInfo.Token == "" {
				ah.logger.Error("authentication returned empty wrapped client token", "backoff", backoffCfg)
				metrics.IncrCounter([]string{ah.metricsSignifier, "auth", "failure"}, 1)

				if backoffSleep(ctx, backoffCfg) {
					continue
				}
				return err
			}
			wrappedResp, err := jsonutil.EncodeJSON(secret.WrapInfo)
			if err != nil {
				ah.logger.Error("failed to encode wrapinfo", "error", err, "backoff", backoffCfg)
				metrics.IncrCounter([]string{ah.metricsSignifier, "auth", "failure"}, 1)

				if backoffSleep(ctx, backoffCfg) {
					continue
				}
				return err
			}
			ah.logger.Info("authentication successful, sending wrapped token to sinks and pausing")
			ah.OutputCh <- string(wrappedResp)
			if ah.enableTemplateTokenCh {
				ah.TemplateTokenCh <- string(wrappedResp)
			}
			if ah.enableExecTokenCh {
				ah.ExecTokenCh <- string(wrappedResp)
			}

			am.CredSuccess()
			backoffCfg.backoff.Reset()

			select {
			case <-ctx.Done():
				ah.logger.Info("shutdown triggered")
				continue

			case <-credCh:
				ah.logger.Info("auth method found new credentials, re-authenticating")
				continue
			}

		default:
			// We handle the token_file method specially, as it's the only
			// auth method that isn't actually authenticating, i.e. the secret
			// returned does not have an Auth struct attached
			isTokenFileMethod := path == "auth/token/lookup-self"
			if isTokenFileMethod {
				// We still check the response of the request to ensure the token is valid
				// i.e. if the token is invalid, we will fail in the authentication step
				if secret == nil || secret.Data == nil {
					ah.logger.Error("token file validation failed, token may be invalid", "backoff", backoffCfg)
					metrics.IncrCounter([]string{ah.metricsSignifier, "auth", "failure"}, 1)

					if backoffSleep(ctx, backoffCfg) {
						continue
					}
					return err
				}
				token, ok := secret.Data["id"].(string)
				if !ok || token == "" {
					ah.logger.Error("token file validation returned empty client token", "backoff", backoffCfg)
					metrics.IncrCounter([]string{ah.metricsSignifier, "auth", "failure"}, 1)

					if backoffSleep(ctx, backoffCfg) {
						continue
					}
					return err
				}

				duration, _ := secret.Data["ttl"].(json.Number).Int64()
				leaseDuration = int(duration)
				renewable, _ := secret.Data["renewable"].(bool)
				secret.Auth = &api.SecretAuth{
					ClientToken:   token,
					LeaseDuration: int(duration),
					Renewable:     renewable,
				}
				ah.logger.Info("authentication successful, sending token to sinks")
				ah.OutputCh <- token
				if ah.enableTemplateTokenCh {
					ah.TemplateTokenCh <- token
				}
				if ah.enableExecTokenCh {
					ah.ExecTokenCh <- token
				}

				tokenType := secret.Data["type"].(string)
				if tokenType == "batch" {
					ah.logger.Info("note that this token type is batch, and batch tokens cannot be renewed", "ttl", leaseDuration)
				}
			} else {
				if secret == nil || secret.Auth == nil {
					ah.logger.Error("authentication returned nil auth info", "backoff", backoffCfg)
					metrics.IncrCounter([]string{ah.metricsSignifier, "auth", "failure"}, 1)

					if backoffSleep(ctx, backoffCfg) {
						continue
					}
					return err
				}
				if secret.Auth.ClientToken == "" {
					ah.logger.Error("authentication returned empty client token", "backoff", backoffCfg)
					metrics.IncrCounter([]string{ah.metricsSignifier, "auth", "failure"}, 1)

					if backoffSleep(ctx, backoffCfg) {
						continue
					}
					return err
				}

				leaseDuration = secret.LeaseDuration
				ah.logger.Info("authentication successful, sending token to sinks")
				ah.OutputCh <- secret.Auth.ClientToken
				if ah.enableTemplateTokenCh {
					ah.TemplateTokenCh <- secret.Auth.ClientToken
				}
				if ah.enableExecTokenCh {
					ah.ExecTokenCh <- secret.Auth.ClientToken
				}
			}

			am.CredSuccess()
			backoffCfg.backoff.Reset()
		}

		if watcher != nil {
			watcher.Stop()
		}

		watcher, err = clientToUse.NewLifetimeWatcher(&api.LifetimeWatcherInput{
			Secret: secret,
		})
		if err != nil {
			ah.logger.Error("error creating lifetime watcher", "error", err, "backoff", backoffCfg)
			metrics.IncrCounter([]string{ah.metricsSignifier, "auth", "failure"}, 1)

			if backoffSleep(ctx, backoffCfg) {
				continue
			}
			return err
		}

		metrics.IncrCounter([]string{ah.metricsSignifier, "auth", "success"}, 1)
		// We don't want to trigger the renewal process for tokens with
		// unlimited TTL, such as the root token.
		if leaseDuration == 0 && isTokenFileMethod {
			ah.logger.Info("not starting token renewal process, as token has unlimited TTL")
		} else {
			ah.logger.Info("starting renewal process")
			go watcher.Renew()
		}

	LifetimeWatcherLoop:
		for {
			select {
			case <-ctx.Done():
				ah.logger.Info("shutdown triggered, stopping lifetime watcher")
				watcher.Stop()
				break LifetimeWatcherLoop

			case err := <-watcher.DoneCh():
				ah.logger.Info("lifetime watcher done channel triggered")
				if err != nil {
					metrics.IncrCounter([]string{ah.metricsSignifier, "auth", "failure"}, 1)
					ah.logger.Error("error renewing token", "error", err)
				}
				break LifetimeWatcherLoop

			case <-watcher.RenewCh():
				metrics.IncrCounter([]string{ah.metricsSignifier, "auth", "success"}, 1)
				ah.logger.Info("renewed auth token")

			case <-credCh:
				ah.logger.Info("auth method found new credentials, re-authenticating")
				break LifetimeWatcherLoop
			}
		}
	}
}

// autoAuthBackoff tracks exponential backoff state.
type autoAuthBackoff struct {
	backoff *backoff.Backoff
}

func newAutoAuthBackoff(min, max time.Duration, exitErr bool) *autoAuthBackoff {
	if max <= 0 {
		max = defaultMaxBackoff
	}

	if min <= 0 {
		min = defaultMinBackoff
	}

	retries := math.MaxInt
	if exitErr {
		retries = 0
	}

	b := backoff.NewBackoff(retries, min, max)

	return &autoAuthBackoff{
		backoff: b,
	}
}

func (b autoAuthBackoff) String() string {
	return b.backoff.Current().Truncate(10 * time.Millisecond).String()
}
