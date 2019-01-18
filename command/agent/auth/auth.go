package auth

import (
	"context"
	"math/rand"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/jsonutil"
)

type AuthMethod interface {
	Authenticate(context.Context, *api.Client) (string, map[string]interface{}, error)
	NewCreds() chan struct{}
	CredSuccess()
	Shutdown()
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
	DoneCh                       chan struct{}
	OutputCh                     chan string
	logger                       hclog.Logger
	client                       *api.Client
	random                       *rand.Rand
	wrapTTL                      time.Duration
	enableReauthOnNewCredentials bool
}

type AuthHandlerConfig struct {
	Logger                       hclog.Logger
	Client                       *api.Client
	WrapTTL                      time.Duration
	EnableReauthOnNewCredentials bool
}

func NewAuthHandler(conf *AuthHandlerConfig) *AuthHandler {
	ah := &AuthHandler{
		DoneCh: make(chan struct{}),
		// This is buffered so that if we try to output after the sink server
		// has been shut down, during agent shutdown, we won't block
		OutputCh:                     make(chan string, 1),
		logger:                       conf.Logger,
		client:                       conf.Client,
		random:                       rand.New(rand.NewSource(int64(time.Now().Nanosecond()))),
		wrapTTL:                      conf.WrapTTL,
		enableReauthOnNewCredentials: conf.EnableReauthOnNewCredentials,
	}

	return ah
}

func backoffOrQuit(ctx context.Context, backoff time.Duration) {
	select {
	case <-time.After(backoff):
	case <-ctx.Done():
	}
}

func (ah *AuthHandler) Run(ctx context.Context, am AuthMethod) {
	if am == nil {
		panic("nil auth method")
	}

	ah.logger.Info("starting auth handler")
	defer func() {
		am.Shutdown()
		close(ah.OutputCh)
		close(ah.DoneCh)
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

	var renewer *api.Renewer

	for {
		select {
		case <-ctx.Done():
			return

		default:
		}

		// Create a fresh backoff value
		backoff := 2*time.Second + time.Duration(ah.random.Int63()%int64(time.Second*2)-int64(time.Second))

		ah.logger.Info("authenticating")
		path, data, err := am.Authenticate(ctx, ah.client)
		if err != nil {
			ah.logger.Error("error getting path or data from method", "error", err, "backoff", backoff.Seconds())
			backoffOrQuit(ctx, backoff)
			continue
		}

		clientToUse := ah.client
		if ah.wrapTTL > 0 {
			wrapClient, err := ah.client.Clone()
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

		secret, err := clientToUse.Logical().Write(path, data)
		// Check errors/sanity
		if err != nil {
			ah.logger.Error("error authenticating", "error", err, "backoff", backoff.Seconds())
			backoffOrQuit(ctx, backoff)
			continue
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

			am.CredSuccess()
		}

		if renewer != nil {
			renewer.Stop()
		}

		renewer, err = ah.client.NewRenewer(&api.RenewerInput{
			Secret: secret,
		})
		if err != nil {
			ah.logger.Error("error creating renewer, backing off and retrying", "error", err, "backoff", backoff.Seconds())
			backoffOrQuit(ctx, backoff)
			continue
		}

		// Start the renewal process
		ah.logger.Info("starting renewal process")
		go renewer.Renew()

	RenewerLoop:
		for {
			select {
			case <-ctx.Done():
				ah.logger.Info("shutdown triggered, stopping renewer")
				renewer.Stop()
				break RenewerLoop

			case err := <-renewer.DoneCh():
				ah.logger.Info("renewer done channel triggered")
				if err != nil {
					ah.logger.Error("error renewing token", "error", err)
				}
				break RenewerLoop

			case <-renewer.RenewCh():
				ah.logger.Info("renewed auth token")

			case <-credCh:
				ah.logger.Info("auth method found new credentials, re-authenticating")
				break RenewerLoop
			}
		}
	}
}
