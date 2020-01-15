package auth

import (
	"context"
	"math/rand"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
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
	TemplateTokenCh              chan string
	logger                       hclog.Logger
	client                       *api.Client
	random                       *rand.Rand
	wrapTTL                      time.Duration
	enableReauthOnNewCredentials bool
	enableTemplateTokenCh        bool
}

type AuthHandlerConfig struct {
	Logger                       hclog.Logger
	Client                       *api.Client
	WrapTTL                      time.Duration
	EnableReauthOnNewCredentials bool
	EnableTemplateTokenCh        bool
}

func NewAuthHandler(conf *AuthHandlerConfig) *AuthHandler {
	ah := &AuthHandler{
		DoneCh: make(chan struct{}),
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
	}

	return ah
}

func backoffOrQuit(ctx context.Context, backoff time.Duration) {
	select {
	case <-time.After(backoff):
	case <-ctx.Done():
	}
}

// Run is a convienence method to call RunWithConnectionTimeout with the default value
// of 0 connTimeout, which indicates no maximum attempts.
func (ah *AuthHandler) Run(ctx context.Context, am AuthMethod) {
	ah.RunWithConnectionTimeout(ctx, am, 0)
}

// Run is a long running process that authenticates with Vault, retrieves an
// access token, and passes that token to it's output channels to be consumed
// down the pipeline. The MaxAttempts variable determines how many attempts can
// be made to a non-responsive Vault instance. A value of 0 for
// connTimeout means there is no limit to the connection attemts
// (Vault Agent will not abort due to no connection being established).
func (ah *AuthHandler) RunWithConnectionTimeout(ctx context.Context, am AuthMethod, connTimeout time.Duration) {
	if am == nil {
		panic("nil auth method")
	}

	var connTimer *time.Timer
	ah.logger.Info("starting auth handler")
	defer func() {
		am.Shutdown()
		close(ah.OutputCh)
		close(ah.DoneCh)
		close(ah.TemplateTokenCh)
		ah.logger.Info("auth handler stopped")
		if connTimer != nil {
			if !connTimer.Stop() {
				<-connTimer.C
			}
		}
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

	timeLimit := connTimeout

	var connErrCount int
	for {
		// If the timer is not nil and connection timeout is set, create and start
		// timer with the configured duration
		//
		// If the timer is not nil, then we're in the connection loop and we simply
		// select. If it fires, we've exceeded the timeout. If we are successful in
		// connecting, then the timer will be stoped

		if connTimer == nil {
			connTimer = time.NewTimer(timeLimit)
		}

		select {
		case <-ctx.Done():
			return

		default:
		}

		if timeLimit != 0 {
			select {
			case <-connTimer.C:
				ah.logger.Error("connection attempt timeout hit, aborting", "timeout", timeLimit.String())
				return
			default:
			}
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
			connErrCount++
			continue
		}

		// reset the connection error count
		connErrCount = 0
		if !connTimer.Stop() {
			<-connTimer.C
		}
		connTimer = nil

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

		watcher, err = ah.client.NewLifetimeWatcher(&api.LifetimeWatcherInput{
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
