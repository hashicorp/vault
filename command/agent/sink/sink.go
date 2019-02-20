package sink

import (
	"context"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"sync/atomic"
	"time"

	"github.com/hashicorp/errwrap"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/dhutil"
	"github.com/hashicorp/vault/helper/jsonutil"
)

type Sink interface {
	WriteToken(string) error
}

type SinkReader interface {
	Token() string
}

type SinkConfig struct {
	Sink
	Logger             hclog.Logger
	Config             map[string]interface{}
	Client             *api.Client
	WrapTTL            time.Duration
	DHType             string
	DHPath             string
	AAD                string
	cachedRemotePubKey []byte
	cachedPubKey       []byte
	cachedPriKey       []byte
}

type SinkServerConfig struct {
	Logger        hclog.Logger
	Client        *api.Client
	Context       context.Context
	ExitAfterAuth bool
}

// SinkServer is responsible for pushing tokens to sinks
type SinkServer struct {
	DoneCh        chan struct{}
	logger        hclog.Logger
	client        *api.Client
	random        *rand.Rand
	exitAfterAuth bool
	remaining     *int32
}

func NewSinkServer(conf *SinkServerConfig) *SinkServer {
	ss := &SinkServer{
		DoneCh:        make(chan struct{}),
		logger:        conf.Logger,
		client:        conf.Client,
		random:        rand.New(rand.NewSource(int64(time.Now().Nanosecond()))),
		exitAfterAuth: conf.ExitAfterAuth,
		remaining:     new(int32),
	}

	return ss
}

// Run executes the server's run loop, which is responsible for reading
// in new tokens and pushing them out to the various sinks.
func (ss *SinkServer) Run(ctx context.Context, incoming chan string, sinks []*SinkConfig) {
	if incoming == nil {
		panic("incoming or shutdown channel are nil")
	}

	ss.logger.Info("starting sink server")
	defer func() {
		ss.logger.Info("sink server stopped")
		close(ss.DoneCh)
	}()

	latestToken := new(string)
	sinkCh := make(chan func() error, len(sinks))
	for {
		select {
		case <-ctx.Done():
			return

		case token := <-incoming:
			if token != *latestToken {

				// Drain the existing funcs
			drainLoop:
				for {
					select {
					case <-sinkCh:
						atomic.AddInt32(ss.remaining, -1)
					default:
						break drainLoop
					}
				}

				*latestToken = token

				sinkFunc := func(currSink *SinkConfig, currToken string) func() error {
					return func() error {
						if currToken != *latestToken {
							return nil
						}
						var err error

						if currSink.WrapTTL != 0 {
							if currToken, err = currSink.wrapToken(ss.client, currSink.WrapTTL, currToken); err != nil {
								return err
							}
						}

						if currSink.DHType != "" {
							if currToken, err = currSink.encryptToken(currToken); err != nil {
								return err
							}
						}

						return currSink.WriteToken(currToken)
					}
				}

				for _, s := range sinks {
					atomic.AddInt32(ss.remaining, 1)
					sinkCh <- sinkFunc(s, token)
				}
			}

		case sinkFunc := <-sinkCh:
			atomic.AddInt32(ss.remaining, -1)
			select {
			case <-ctx.Done():
				return
			default:
			}

			if err := sinkFunc(); err != nil {
				backoff := 2*time.Second + time.Duration(ss.random.Int63()%int64(time.Second*2)-int64(time.Second))
				ss.logger.Error("error returned by sink function, retrying", "error", err, "backoff", backoff.String())
				select {
				case <-ctx.Done():
					return
				case <-time.After(backoff):
					atomic.AddInt32(ss.remaining, 1)
					sinkCh <- sinkFunc
				}
			} else {
				if atomic.LoadInt32(ss.remaining) == 0 && ss.exitAfterAuth {
					return
				}
			}
		}
	}
}

func (s *SinkConfig) encryptToken(token string) (string, error) {
	var aesKey []byte
	var err error
	resp := new(dhutil.Envelope)
	switch s.DHType {
	case "curve25519":
		if len(s.cachedRemotePubKey) == 0 {
			_, err = os.Lstat(s.DHPath)
			if err != nil {
				if !os.IsNotExist(err) {
					return "", errwrap.Wrapf("error stat-ing dh parameters file: {{err}}", err)
				}
				return "", errors.New("no dh parameters file found, and no cached pub key")
			}
			fileBytes, err := ioutil.ReadFile(s.DHPath)
			if err != nil {
				return "", errwrap.Wrapf("error reading file for dh parameters: {{err}}", err)
			}
			theirPubKey := new(dhutil.PublicKeyInfo)
			if err := jsonutil.DecodeJSON(fileBytes, theirPubKey); err != nil {
				return "", errwrap.Wrapf("error decoding public key: {{err}}", err)
			}
			if len(theirPubKey.Curve25519PublicKey) == 0 {
				return "", errors.New("public key is nil")
			}
			s.cachedRemotePubKey = theirPubKey.Curve25519PublicKey
		}
		if len(s.cachedPubKey) == 0 {
			s.cachedPubKey, s.cachedPriKey, err = dhutil.GeneratePublicPrivateKey()
			if err != nil {
				return "", errwrap.Wrapf("error generating pub/pri curve25519 keys: {{err}}", err)
			}
		}
		resp.Curve25519PublicKey = s.cachedPubKey
	}

	aesKey, err = dhutil.GenerateSharedKey(s.cachedPriKey, s.cachedRemotePubKey)
	if err != nil {
		return "", errwrap.Wrapf("error deriving shared key: {{err}}", err)
	}
	if len(aesKey) == 0 {
		return "", errors.New("derived AES key is empty")
	}

	resp.EncryptedPayload, resp.Nonce, err = dhutil.EncryptAES(aesKey, []byte(token), []byte(s.AAD))
	if err != nil {
		return "", errwrap.Wrapf("error encrypting with shared key: {{err}}", err)
	}
	m, err := jsonutil.EncodeJSON(resp)
	if err != nil {
		return "", errwrap.Wrapf("error encoding encrypted payload: {{err}}", err)
	}

	return string(m), nil
}

func (s *SinkConfig) wrapToken(client *api.Client, wrapTTL time.Duration, token string) (string, error) {
	wrapClient, err := client.Clone()
	if err != nil {
		return "", errwrap.Wrapf("error deriving client for wrapping, not writing out to sink: {{err}})", err)
	}
	wrapClient.SetToken(token)
	wrapClient.SetWrappingLookupFunc(func(string, string) string {
		return wrapTTL.String()
	})
	secret, err := wrapClient.Logical().Write("sys/wrapping/wrap", map[string]interface{}{
		"token": token,
	})
	if err != nil {
		return "", errwrap.Wrapf("error wrapping token, not writing out to sink: {{err}})", err)
	}
	if secret == nil {
		return "", errors.New("nil secret returned, not writing out to sink")
	}
	if secret.WrapInfo == nil {
		return "", errors.New("nil wrap info returned, not writing out to sink")
	}

	m, err := jsonutil.EncodeJSON(secret.WrapInfo)
	if err != nil {
		return "", errwrap.Wrapf("error marshaling token, not writing out to sink: {{err}})", err)
	}

	return string(m), nil
}
