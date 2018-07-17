package auth

import (
	"math/rand"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
)

type AuthMethod interface {
	Authenticate(*api.Client) (*api.Secret, error)
}

type AuthConfig struct {
	Logger hclog.Logger
	Config map[string]interface{}
}

// AuthHandler is responsible for keeping a token alive and renewed and passing
// new tokens to the sink server
type AuthHandler struct {
	DoneCh     chan struct{}
	ShutdownCh chan struct{}
	OutputCh   chan string
	logger     hclog.Logger
	client     *api.Client
	random     *rand.Rand
}

type AuthHandlerConfig struct {
	Logger hclog.Logger
	Client *api.Client
}

func NewAuthHandler(conf *AuthHandlerConfig) *AuthHandler {
	ah := &AuthHandler{
		ShutdownCh: make(chan struct{}),
		DoneCh:     make(chan struct{}),
		OutputCh:   make(chan string),
		logger:     conf.Logger,
		client:     conf.Client,
		random:     rand.New(rand.NewSource(int64(time.Now().Nanosecond()))),
	}

	return ah
}

func (ah *AuthHandler) Run(am AuthMethod) {
	if am == nil {
		panic("nil auth method")
	}

	ah.logger.Info("starting auth handler")
	defer func() {
		ah.logger.Info("auth handler stopped")
		close(ah.DoneCh)
	}()

	for {
		select {
		case <-ah.ShutdownCh:
			return
		default:
		}

		// Create a fresh backoff value
		backoff := 2*time.Second + time.Duration(ah.random.Int63()%int64(time.Second*2)-int64(time.Second))

		ah.logger.Info("authenticating")
		secret, err := am.Authenticate(ah.client)
		// Check errors/sanity
		if err != nil {
			ah.logger.Error("error authenticating, backing off and retrying", "error", err, "backoff", backoff.Seconds())
			time.Sleep(backoff)
			continue
		}
		if secret.Auth == nil {
			ah.logger.Error("authentication returned nil auth info, backing off and retrying", "backoff", backoff.Seconds())
			time.Sleep(backoff)
			continue
		}
		if secret.Auth.ClientToken == "" {
			ah.logger.Error("authentication returned empty client token, backing off and retrying", "backoff", backoff.Seconds())
			time.Sleep(backoff)
			continue
		}

		// Output to the sinks
		ah.logger.Info("authentication successful, sending token to sinks")
		ah.OutputCh <- secret.Auth.ClientToken

		renewer, err := ah.client.NewRenewer(&api.RenewerInput{
			Secret: secret,
		})
		if err != nil {
			ah.logger.Error("error creating renewer, backing off and retrying", "error", err, "backoff", backoff.Seconds())
			time.Sleep(backoff)
			continue
		}

		// Start the renewal process
		ah.logger.Info("starting renewal process")
		go renewer.Renew()

	RenewerLoop:
		for {
			select {
			case <-ah.ShutdownCh:
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
			}
		}
	}
}
