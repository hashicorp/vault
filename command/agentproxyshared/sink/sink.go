// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package sink

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/dhutil"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
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
	DeriveKey          bool
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
	logger        hclog.Logger
	client        *api.Client
	random        *rand.Rand
	exitAfterAuth bool
	remaining     *int32
}

func NewSinkServer(conf *SinkServerConfig) *SinkServer {
	ss := &SinkServer{
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
func (ss *SinkServer) Run(ctx context.Context, incoming chan string, sinks []*SinkConfig, tokenWriteInProgress *atomic.Bool) error {
	latestToken := new(string)
	writeSink := func(currSink *SinkConfig, currToken string) error {
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

	if incoming == nil {
		return errors.New("sink server: incoming channel is nil")
	}

	ss.logger.Info("starting sink server")
	defer func() {
		tokenWriteInProgress.Store(false)
		ss.logger.Info("sink server stopped")
	}()

	type sinkToken struct {
		sink  *SinkConfig
		token string
	}
	sinkCh := make(chan sinkToken, len(sinks))
	for {
		select {
		case <-ctx.Done():
			return nil

		case token := <-incoming:
			if len(sinks) > 0 {
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

					for _, s := range sinks {
						atomic.AddInt32(ss.remaining, 1)
						sinkCh <- sinkToken{s, token}
					}
				}
			} else {
				ss.logger.Trace("no sinks, ignoring new token")
				tokenWriteInProgress.Store(false)
				if ss.exitAfterAuth {
					ss.logger.Trace("no sinks, exitAfterAuth, bye")
					return nil
				}
			}
		case st := <-sinkCh:
			atomic.AddInt32(ss.remaining, -1)
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			if err := writeSink(st.sink, st.token); err != nil {
				backoff := 2*time.Second + time.Duration(ss.random.Int63()%int64(time.Second*2)-int64(time.Second))
				ss.logger.Error("error returned by sink function, retrying", "error", err, "backoff", backoff.String())
				timer := time.NewTimer(backoff)
				select {
				case <-ctx.Done():
					timer.Stop()
					return nil
				case <-timer.C:
					atomic.AddInt32(ss.remaining, 1)
					sinkCh <- st
				}
			} else {
				if atomic.LoadInt32(ss.remaining) == 0 {
					tokenWriteInProgress.Store(false)
					if ss.exitAfterAuth {
						return nil
					}
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
					return "", fmt.Errorf("error stat-ing dh parameters file: %w", err)
				}
				return "", errors.New("no dh parameters file found, and no cached pub key")
			}
			fileBytes, err := ioutil.ReadFile(s.DHPath)
			if err != nil {
				return "", fmt.Errorf("error reading file for dh parameters: %w", err)
			}
			theirPubKey := new(dhutil.PublicKeyInfo)
			if err := jsonutil.DecodeJSON(fileBytes, theirPubKey); err != nil {
				return "", fmt.Errorf("error decoding public key: %w", err)
			}
			if len(theirPubKey.Curve25519PublicKey) == 0 {
				return "", errors.New("public key is nil")
			}
			s.cachedRemotePubKey = theirPubKey.Curve25519PublicKey
		}
		if len(s.cachedPubKey) == 0 {
			s.cachedPubKey, s.cachedPriKey, err = dhutil.GeneratePublicPrivateKey()
			if err != nil {
				return "", fmt.Errorf("error generating pub/pri curve25519 keys: %w", err)
			}
		}
		resp.Curve25519PublicKey = s.cachedPubKey
	}

	secret, err := dhutil.GenerateSharedSecret(s.cachedPriKey, s.cachedRemotePubKey)
	if err != nil {
		return "", fmt.Errorf("error calculating shared key: %w", err)
	}
	if s.DeriveKey {
		aesKey, err = dhutil.DeriveSharedKey(secret, s.cachedPubKey, s.cachedRemotePubKey)
	} else {
		aesKey = secret
	}

	if err != nil {
		return "", fmt.Errorf("error deriving shared key: %w", err)
	}
	if len(aesKey) == 0 {
		return "", errors.New("derived AES key is empty")
	}

	resp.EncryptedPayload, resp.Nonce, err = dhutil.EncryptAES(aesKey, []byte(token), []byte(s.AAD))
	if err != nil {
		return "", fmt.Errorf("error encrypting with shared key: %w", err)
	}
	m, err := jsonutil.EncodeJSON(resp)
	if err != nil {
		return "", fmt.Errorf("error encoding encrypted payload: %w", err)
	}

	return string(m), nil
}

func (s *SinkConfig) wrapToken(client *api.Client, wrapTTL time.Duration, token string) (string, error) {
	wrapClient, err := client.CloneWithHeaders()
	if err != nil {
		return "", fmt.Errorf("error deriving client for wrapping, not writing out to sink: %w)", err)
	}

	wrapClient.SetToken(token)
	wrapClient.SetWrappingLookupFunc(func(string, string) string {
		return wrapTTL.String()
	})

	secret, err := wrapClient.Logical().Write("sys/wrapping/wrap", map[string]interface{}{
		"token": token,
	})
	if err != nil {
		return "", fmt.Errorf("error wrapping token, not writing out to sink: %w)", err)
	}
	if secret == nil {
		return "", errors.New("nil secret returned, not writing out to sink")
	}
	if secret.WrapInfo == nil {
		return "", errors.New("nil wrap info returned, not writing out to sink")
	}

	m, err := jsonutil.EncodeJSON(secret.WrapInfo)
	if err != nil {
		return "", fmt.Errorf("error marshaling token, not writing out to sink: %w)", err)
	}

	return string(m), nil
}
