package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
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
	logger                       hclog.Logger
	client                       *api.Client
	random                       *rand.Rand
	wrapTTL                      time.Duration
	enableReauthOnNewCredentials bool
	enableTemplateTokenCh        bool
	enableUseExistingToken       bool
}

type AuthHandlerConfig struct {
	Logger                       hclog.Logger
	Client                       *api.Client
	WrapTTL                      time.Duration
	EnableReauthOnNewCredentials bool
	EnableTemplateTokenCh        bool
	EnableUseExistingToken       bool
}

func NewAuthHandler(conf *AuthHandlerConfig) *AuthHandler {
	ah := &AuthHandler{
		// This is buffered so that if we try to output after the sink server
		// has been shut down, during agent shutdown, we won't block
		OutputCh:                     make(chan string, 1),
		TemplateTokenCh:              make(chan string, 1),
		logger:                       conf.Logger,
		client:                       conf.Client,
		random:                       rand.New(rand.NewSource(int64(time.Now().Nanosecond()))),
		wrapTTL:                      conf.WrapTTL,
		enableReauthOnNewCredentials: conf.EnableReauthOnNewCredentials,
		enableTemplateTokenCh:        conf.EnableTemplateTokenCh,
		enableUseExistingToken:       conf.EnableUseExistingToken,
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

		// Create a fresh backoff value
		backoff := 2*time.Second + time.Duration(ah.random.Int63()%int64(time.Second*2)-int64(time.Second))

		ah.logger.Info("authenticating")

		path, header, data, err := am.Authenticate(ctx, ah.client)
		if err != nil {
			ah.logger.Error("error getting path or data from method", "error", err, "backoff", backoff.Seconds())
			backoffOrQuit(ctx, backoff)
			continue
		}

		var clientToUse *api.Client

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

		var secret *api.Secret = new(api.Secret)
		if first && ah.enableUseExistingToken {
			var token string
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("error determining home directory: %s", err)
			}
			tokenPath := fmt.Sprintf("%s/.vault-token", home)

			if fileExists(tokenPath) {
				ah.logger.Debug("attempting to read preloaded token from path", "path", tokenPath)
				tokenRaw, err := ioutil.ReadFile(tokenPath)
				if err != nil {
					return fmt.Errorf("error opening token file: %s", err)
				}
				token = string(tokenRaw)
			} else {
				ah.logger.Debug("attempting to read preloaded token from env VAULT_TOKEN")
				token = os.Getenv("VAULT_TOKEN")
			}

			if token != "" {
				first = false

				ah.logger.Debug("lookup-self with preloaded token")

				clientToUse.SetToken(string(token))
				secret, err = clientToUse.Logical().Read("auth/token/lookup-self")
				if err != nil {
					ah.logger.Error("error looking up token", "error", err, "backoff", backoff.Seconds())
					backoffOrQuit(ctx, backoff)
					continue
				}

				duration, _ := secret.Data["ttl"].(json.Number).Int64()
				secret.Auth = &api.SecretAuth{
					ClientToken:   secret.Data["id"].(string),
					LeaseDuration: int(duration),
					Renewable:     secret.Data["renewable"].(bool),
				}
			}
		}

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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
