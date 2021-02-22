package auth

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
)

const (
	initialBackoff    = 1 * time.Second
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
	enableReauthOnNewCredentials bool
	enableTemplateTokenCh        bool
}

type AuthHandlerConfig struct {
	Logger                       hclog.Logger
	Client                       *api.Client
	WrapTTL                      time.Duration
	MaxBackoff                   time.Duration
	Token                        string
	EnableReauthOnNewCredentials bool
	EnableTemplateTokenCh        bool
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
		maxBackoff:                   conf.MaxBackoff,
		enableReauthOnNewCredentials: conf.EnableReauthOnNewCredentials,
		enableTemplateTokenCh:        conf.EnableTemplateTokenCh,
	}

	return ah
}

func backoffOrQuit(ctx context.Context, backoff time.Duration) {
	select {
	case <-time.After(backoff):
	case <-ctx.Done():
	}
}

func (ah *AuthHandler) Run(ctx context.Context, am AuthMethod) error {
	if am == nil {
		return errors.New("auth handler: nil auth method")
	}

	backoff := initialBackoff
	maxBackoff := defaultMaxBackoff

	if ah.maxBackoff > 0 {
		maxBackoff = ah.maxBackoff
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

		backoff = calculateBackoff(backoff, maxBackoff)

		var clientToUse *api.Client
		var err error
		var path string
		var data map[string]interface{}
		var header http.Header

		switch am.(type) {
		case AuthMethodWithClient:
			clientToUse, err = am.(AuthMethodWithClient).AuthClient(ah.client)
			if err != nil {
				ah.logger.Error("error creating client for authentication call", "error", err, "backoff", backoff.Seconds())
				backoffOrQuit(ctx, backoff)
				continue
			}
		default:
			clientToUse = ah.client
		}

		var secret *api.Secret = new(api.Secret)
		if first && ah.token != "" {
			ah.logger.Debug("using preloaded token")

			first = false
			ah.logger.Debug("lookup-self with preloaded token")
			clientToUse.SetToken(ah.token)

			secret, err = clientToUse.Logical().Read("auth/token/lookup-self")
			if err != nil {
				ah.logger.Error("could not look up token", "err", err, "backoff", backoff.Seconds())
				backoffOrQuit(ctx, backoff)
				continue
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
				ah.logger.Error("error getting path or data from method", "error", err, "backoff", backoff.Seconds())
				backoffOrQuit(ctx, backoff)
				continue
			}
		}

		if ah.wrapTTL > 0 {
			wrapClient, err := clientToUse.Clone()
			if err != nil {
				ah.logger.Error("error creating client for wrapped call", "error", err, "backoff", backoff.Seconds())
				backoffOrQuit(ctx, backoff)
				continue
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
			secret, err = clientToUse.Logical().Write(path, data)
			// Check errors/sanity
			if err != nil {
				ah.logger.Error("error authenticating", "error", err, "backoff", backoff.Seconds())
				backoffOrQuit(ctx, backoff)
				continue
			}
		}

		switch {
		case ah.wrapTTL > 0:
			if secret.WrapInfo == nil {
				ah.logger.Error("authentication returned nil wrap info", "backoff", backoff.Seconds())
				backoffOrQuit(ctx, backoff)
				continue
			}
			if secret.WrapInfo.Token == "" {
				ah.logger.Error("authentication returned empty wrapped client token", "backoff", backoff.Seconds())
				backoffOrQuit(ctx, backoff)
				continue
			}
			wrappedResp, err := jsonutil.EncodeJSON(secret.WrapInfo)
			if err != nil {
				ah.logger.Error("failed to encode wrapinfo", "error", err, "backoff", backoff.Seconds())
				backoffOrQuit(ctx, backoff)
				continue
			}
			ah.logger.Info("authentication successful, sending wrapped token to sinks and pausing")
			ah.OutputCh <- string(wrappedResp)
			if ah.enableTemplateTokenCh {
				ah.TemplateTokenCh <- string(wrappedResp)
			}

			am.CredSuccess()

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
				ah.logger.Error("authentication returned nil auth info", "backoff", backoff.Seconds())
				backoffOrQuit(ctx, backoff)
				continue
			}
			if secret.Auth.ClientToken == "" {
				ah.logger.Error("authentication returned empty client token", "backoff", backoff.Seconds())
				backoffOrQuit(ctx, backoff)
				continue
			}
			ah.logger.Info("authentication successful, sending token to sinks")
			ah.OutputCh <- secret.Auth.ClientToken
			if ah.enableTemplateTokenCh {
				ah.TemplateTokenCh <- secret.Auth.ClientToken
			}

			am.CredSuccess()
		}

		if watcher != nil {
			watcher.Stop()
		}

		watcher, err = clientToUse.NewLifetimeWatcher(&api.LifetimeWatcherInput{
			Secret: secret,
		})
		if err != nil {
			ah.logger.Error("error creating lifetime watcher, backing off and retrying", "error", err, "backoff", backoff.Seconds())
			backoffOrQuit(ctx, backoff)
			continue
		}

		// Start the renewal process
		ah.logger.Info("starting renewal process")
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
					ah.logger.Error("error renewing token", "error", err)
				}
				break LifetimeWatcherLoop

			case <-watcher.RenewCh():
				ah.logger.Info("renewed auth token")

			case <-credCh:
				ah.logger.Info("auth method found new credentials, re-authenticating")
				break LifetimeWatcherLoop
			}
		}
	}
}

// calculateBackoff determines a new backoff duration that is roughly twice
// the previous value, capped to a max value, with a measure of randomness.
func calculateBackoff(previous, max time.Duration) time.Duration {
	maxBackoff := math.Min(2*previous.Seconds(), max.Seconds())
	backoff := 0.75*maxBackoff + 0.25*rand.Float64()*maxBackoff

	// Scale to ms to avoid truncation errors for < 2s values, which can
	// leave the backoff stuck at 1s on successive calls.
	return time.Duration(backoff*1000) * time.Millisecond
}
