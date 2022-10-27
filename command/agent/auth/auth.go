package auth

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
)

const (
	defaultMinBackoff = 1 * time.Second
	defaultMaxBackoff = 5 * time.Minute
)

// AuthMethod is the interface that auto-auth methods implement for the agent
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
	token                        string
	logger                       hclog.Logger
	client                       *api.Client
	random                       *rand.Rand
	wrapTTL                      time.Duration
	maxBackoff                   time.Duration
	minBackoff                   time.Duration
	enableReauthOnNewCredentials bool
	enableTemplateTokenCh        bool
	exitOnError                  bool
}

type AuthHandlerConfig struct {
	Logger                       hclog.Logger
	Client                       *api.Client
	WrapTTL                      time.Duration
	MaxBackoff                   time.Duration
	MinBackoff                   time.Duration
	Token                        string
	EnableReauthOnNewCredentials bool
	EnableTemplateTokenCh        bool
	ExitOnError                  bool
}

func NewAuthHandler(conf *AuthHandlerConfig) *AuthHandler {
	ah := &AuthHandler{
		// This is buffered so that if we try to output after the sink server
		// has been shut down, during agent shutdown, we won't block
		OutputCh:                     make(chan string, 1),
		TemplateTokenCh:              make(chan string, 1),
		token:                        conf.Token,
		logger:                       conf.Logger,
		client:                       conf.Client,
		random:                       rand.New(rand.NewSource(int64(time.Now().Nanosecond()))),
		wrapTTL:                      conf.WrapTTL,
		minBackoff:                   conf.MinBackoff,
		maxBackoff:                   conf.MaxBackoff,
		enableReauthOnNewCredentials: conf.EnableReauthOnNewCredentials,
		enableTemplateTokenCh:        conf.EnableTemplateTokenCh,
		exitOnError:                  conf.ExitOnError,
	}

	return ah
}

func backoff(ctx context.Context, backoff *agentBackoff) bool {
	if backoff.exitOnErr {
		return false
	}

	select {
	case <-time.After(backoff.current):
	case <-ctx.Done():
	}

	// Increase exponential backoff for the next time if we don't
	// successfully auth/renew/etc.
	backoff.next()
	return true
}

func (ah *AuthHandler) Run(ctx context.Context, am AuthMethod) error {
	if am == nil {
		return errors.New("auth handler: nil auth method")
	}

	if ah.minBackoff <= 0 {
		ah.minBackoff = defaultMinBackoff
	}

	backoffCfg := newAgentBackoff(ah.minBackoff, ah.maxBackoff, ah.exitOnError)

	if backoffCfg.min >= backoffCfg.max {
		return errors.New("auth handler: min_backoff cannot be greater than max_backoff")
	}

	ah.logger.Info("starting auth handler")
	defer func() {
		am.Shutdown()
		close(ah.OutputCh)
		close(ah.TemplateTokenCh)
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

		switch am.(type) {
		case AuthMethodWithClient:
			clientToUse, err = am.(AuthMethodWithClient).AuthClient(ah.client)
			if err != nil {
				ah.logger.Error("error creating client for authentication call", "error", err, "backoff", backoff)
				metrics.IncrCounter([]string{"agent", "auth", "failure"}, 1)

				if backoff(ctx, backoffCfg) {
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
				metrics.IncrCounter([]string{"agent", "auth", "failure"}, 1)

				if backoff(ctx, backoffCfg) {
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
				metrics.IncrCounter([]string{"agent", "auth", "failure"}, 1)

				if backoff(ctx, backoffCfg) {
					continue
				}
				return err
			}
		}

		if ah.wrapTTL > 0 {
			wrapClient, err := clientToUse.Clone()
			if err != nil {
				ah.logger.Error("error creating client for wrapped call", "error", err, "backoff", backoffCfg)
				metrics.IncrCounter([]string{"agent", "auth", "failure"}, 1)

				if backoff(ctx, backoffCfg) {
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
		//  or if a preloaded token has expired and is now switching to auto-auth.
		if secret.Auth == nil {
			secret, err = clientToUse.Logical().WriteWithContext(ctx, path, data)
			// Check errors/sanity
			if err != nil {
				ah.logger.Error("error authenticating", "error", err, "backoff", backoffCfg)
				metrics.IncrCounter([]string{"agent", "auth", "failure"}, 1)

				if backoff(ctx, backoffCfg) {
					continue
				}
				return err
			}
		}

		switch {
		case ah.wrapTTL > 0:
			if secret.WrapInfo == nil {
				ah.logger.Error("authentication returned nil wrap info", "backoff", backoffCfg)
				metrics.IncrCounter([]string{"agent", "auth", "failure"}, 1)

				if backoff(ctx, backoffCfg) {
					continue
				}
				return err
			}
			if secret.WrapInfo.Token == "" {
				ah.logger.Error("authentication returned empty wrapped client token", "backoff", backoffCfg)
				metrics.IncrCounter([]string{"agent", "auth", "failure"}, 1)

				if backoff(ctx, backoffCfg) {
					continue
				}
				return err
			}
			wrappedResp, err := jsonutil.EncodeJSON(secret.WrapInfo)
			if err != nil {
				ah.logger.Error("failed to encode wrapinfo", "error", err, "backoff", backoffCfg)
				metrics.IncrCounter([]string{"agent", "auth", "failure"}, 1)

				if backoff(ctx, backoffCfg) {
					continue
				}
				return err
			}
			ah.logger.Info("authentication successful, sending wrapped token to sinks and pausing")
			ah.OutputCh <- string(wrappedResp)
			if ah.enableTemplateTokenCh {
				ah.TemplateTokenCh <- string(wrappedResp)
			}

			am.CredSuccess()
			backoffCfg.reset()

			select {
			case <-ctx.Done():
				ah.logger.Info("shutdown triggered")
				continue

			case <-credCh:
				ah.logger.Info("auth method found new credentials, re-authenticating")
				continue
			}

		default:
			if secret == nil || secret.Auth == nil {
				ah.logger.Error("authentication returned nil auth info", "backoff", backoffCfg)
				metrics.IncrCounter([]string{"agent", "auth", "failure"}, 1)

				if backoff(ctx, backoffCfg) {
					continue
				}
				return err
			}
			if secret.Auth.ClientToken == "" {
				ah.logger.Error("authentication returned empty client token", "backoff", backoffCfg)
				metrics.IncrCounter([]string{"agent", "auth", "failure"}, 1)

				if backoff(ctx, backoffCfg) {
					continue
				}
				return err
			}
			ah.logger.Info("authentication successful, sending token to sinks")
			ah.OutputCh <- secret.Auth.ClientToken
			if ah.enableTemplateTokenCh {
				ah.TemplateTokenCh <- secret.Auth.ClientToken
			}

			am.CredSuccess()
			backoffCfg.reset()
		}

		if watcher != nil {
			watcher.Stop()
		}

		watcher, err = clientToUse.NewLifetimeWatcher(&api.LifetimeWatcherInput{
			Secret: secret,
		})
		if err != nil {
			ah.logger.Error("error creating lifetime watcher", "error", err, "backoff", backoffCfg)
			metrics.IncrCounter([]string{"agent", "auth", "failure"}, 1)

			if backoff(ctx, backoffCfg) {
				continue
			}
			return err
		}

		// Start the renewal process
		ah.logger.Info("starting renewal process")
		metrics.IncrCounter([]string{"agent", "auth", "success"}, 1)
		go watcher.Renew()

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
					metrics.IncrCounter([]string{"agent", "auth", "failure"}, 1)
					ah.logger.Error("error renewing token", "error", err)
				}
				break LifetimeWatcherLoop

			case <-watcher.RenewCh():
				metrics.IncrCounter([]string{"agent", "auth", "success"}, 1)
				ah.logger.Info("renewed auth token")

			case <-credCh:
				ah.logger.Info("auth method found new credentials, re-authenticating")
				break LifetimeWatcherLoop
			}
		}
	}
}

// agentBackoff tracks exponential backoff state.
type agentBackoff struct {
	min       time.Duration
	max       time.Duration
	current   time.Duration
	exitOnErr bool
}

func newAgentBackoff(min, max time.Duration, exitErr bool) *agentBackoff {
	if max <= 0 {
		max = defaultMaxBackoff
	}

	if min <= 0 {
		min = defaultMinBackoff
	}

	return &agentBackoff{
		current:   min,
		max:       max,
		min:       min,
		exitOnErr: exitErr,
	}
}

// next determines the next backoff duration that is roughly twice
// the current value, capped to a max value, with a measure of randomness.
func (b *agentBackoff) next() {
	maxBackoff := 2 * b.current

	if maxBackoff > b.max {
		maxBackoff = b.max
	}

	// Trim a random amount (0-25%) off the doubled duration
	trim := rand.Int63n(int64(maxBackoff) / 4)
	b.current = maxBackoff - time.Duration(trim)
}

func (b *agentBackoff) reset() {
	b.current = b.min
}

func (b agentBackoff) String() string {
	return b.current.Truncate(10 * time.Millisecond).String()
}
