// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package spiffe

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/client/auth"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
)

const (
	authTypeJWT  = "jwt"
	authTypeCert = "cert"
)

type spiffeMethod struct {
	logger        hclog.Logger
	mountPath     string
	authType      string
	role          string
	svidTokenFile string

	tlsCfg *api.TLSConfig

	stopCh          chan struct{}
	doneCh          chan struct{}
	credSuccessGate chan struct{}
	once            *sync.Once
	credsFound      chan struct{}

	// JWT state
	latestToken atomic.Pointer[string] // last successfully authenticated token

	// Cert state
	latestHash atomic.Pointer[[]byte] // last successfully authenticated hash
}

var (
	_ auth.AuthMethodWithClient = &spiffeMethod{}
	_ auth.AuthMethodReady      = &spiffeMethod{}
)

// NewSpiffeAuthMethod creates a new SPIFFE auto-auth method handler.
func NewSpiffeAuthMethod(conf *auth.AuthConfig, tlsCfg *api.TLSConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("missing config")
	}

	s := &spiffeMethod{
		logger:          conf.Logger,
		mountPath:       conf.MountPath,
		tlsCfg:          tlsCfg,
		stopCh:          make(chan struct{}),
		doneCh:          make(chan struct{}),
		credSuccessGate: make(chan struct{}),
		once:            new(sync.Once),
		credsFound:      make(chan struct{}),
	}

	typeRaw, ok := conf.Config["type"]
	if !ok {
		return nil, errors.New("missing required 'type' config value")
	}
	authType, ok := typeRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'type' config value to string")
	}
	s.authType = authType

	if roleRaw, ok := conf.Config["role"]; ok {
		roleStr, ok := roleRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'role' config value to string")
		}
		s.role = roleStr
	}

	// client_cert and client_key in the method stanza override the values
	// inherited from the vault stanza. Setting them here keeps the auth client
	// certificate separate from the TLS settings the agent cache uses when
	// proxying requests. Copy tlsCfg so the shared config is not mutated.
	var clientCert, clientKey string
	if raw, ok := conf.Config["client_cert"]; ok {
		clientCert, ok = raw.(string)
		if !ok {
			return nil, errors.New("could not convert 'client_cert' config value to string")
		}
	}
	if raw, ok := conf.Config["client_key"]; ok {
		clientKey, ok = raw.(string)
		if !ok {
			return nil, errors.New("could not convert 'client_key' config value to string")
		}
	}
	if clientCert != "" || clientKey != "" {
		merged := &api.TLSConfig{}
		if tlsCfg != nil {
			*merged = *tlsCfg
		}
		if clientCert != "" {
			merged.ClientCert = clientCert
		}
		if clientKey != "" {
			merged.ClientKey = clientKey
		}
		s.tlsCfg = merged
	}

	switch authType {
	case authTypeJWT:
		svidTokenFileRaw, ok := conf.Config["svid_token_file"]
		if !ok {
			return nil, errors.New("missing required 'svid_token_file' config value for jwt auth type")
		}
		svidTokenFile, ok := svidTokenFileRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'svid_token_file' config value to string")
		}
		if svidTokenFile == "" {
			return nil, errors.New("'svid_token_file' config value cannot be empty")
		}
		s.svidTokenFile = svidTokenFile

	case authTypeCert:
		if s.tlsCfg == nil || s.tlsCfg.ClientCert == "" || s.tlsCfg.ClientKey == "" {
			return nil, errors.New("missing required client certificate and key in tls configuration for cert auth type")
		}
	default:
		return nil, fmt.Errorf("invalid auth type %q: must be 'jwt' or 'cert'", authType)
	}

	readPeriod := 1 * time.Minute
	if readPeriodRaw, ok := conf.Config["reload_period"]; ok {
		period, err := parseutil.ParseDurationSecond(readPeriodRaw)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'reload_period' value: %w", err)
		}
		readPeriod = period
	}

	switch authType {
	case authTypeJWT:
		go s.runWatcher(readPeriod, func() bool {
			s.logger.Debug("checking if svid token file changed")
			currentToken, err := s.readToken()
			if err != nil {
				s.logger.Error("failed to read svid token file during credential refresh", "error", err)
				return false
			}

			if currentToken != "" && currentToken != *s.latestToken.Load() {
				s.logger.Info("svid token file changed, re-authenticating")
				return true
			}
			return false
		})
	case authTypeCert:
		go s.runWatcher(readPeriod, func() bool {
			s.logger.Debug("checking if cert/key files changed")
			newHash, err := s.hashCert(s.tlsCfg.ClientCert, s.tlsCfg.ClientKey)
			if err != nil {
				s.logger.Error("failed to read cert/key files during credential refresh", "error", err)
				return false
			}

			if !bytes.Equal(*s.latestHash.Load(), newHash) {
				s.logger.Info("cert/key files changed, re-authenticating")
				return true
			}
			return false
		})
	}

	s.logger.Info("spiffe auth method created", "type", authType, "mount_path", s.mountPath)

	return s, nil
}

func (s *spiffeMethod) Authenticate(_ context.Context, _ *api.Client) (string, http.Header, map[string]interface{}, error) {
	s.logger.Trace("beginning authentication")

	authMap := map[string]interface{}{
		"type": s.authType,
	}
	if s.role != "" {
		authMap["role"] = s.role
	}

	var header http.Header
	switch s.authType {
	case authTypeJWT:
		token, err := s.readToken()
		if err != nil {
			return "", nil, nil, fmt.Errorf("cannot authenticate, failed to read svid token file: %w", err)
		}
		if token == "" {
			return "", nil, nil, errors.New("cannot authenticate, svid token file is empty")
		}
		s.latestToken.Store(&token)
		header = make(http.Header)
		header.Set("Authorization", "Bearer "+token)
	}

	s.logger.Debug("SPIFFE authenticate", "mount_path", s.mountPath, "type", s.authType, "auth_map", authMap)

	return fmt.Sprintf("%s/login", s.mountPath), header, authMap, nil
}

func (s *spiffeMethod) NewCreds() chan struct{} {
	return s.credsFound
}

func (s *spiffeMethod) CredSuccess() {
	s.once.Do(func() {
		close(s.credSuccessGate)
	})
}

func (s *spiffeMethod) Shutdown() {
	close(s.stopCh)
	<-s.doneCh
}

// WaitForCreds implements auth.AuthMethodReady. It blocks until the credential
// file(s) are present and valid, then returns nil.
//
// For jwt type, it waits until svid_token_file contains a non-empty token.
// For cert type, it waits until the client cert and key files both exist and
// form a valid TLS key pair.
//
// It returns an error if the context or the method's own stop channel is
// closed before a valid credential is found.
func (s *spiffeMethod) WaitForCreds(ctx context.Context) error {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for {
		var ready bool
		switch s.authType {
		case authTypeJWT:
			token, err := s.readToken()
			ready = err == nil && token != ""
		case authTypeCert:
			_, err := s.hashCert(s.tlsCfg.ClientCert, s.tlsCfg.ClientKey)
			ready = err == nil
		}

		if ready {
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("spiffe credentials not ready: %w", ctx.Err())
		case <-s.stopCh:
			return errors.New("spiffe auth method stopped while waiting for credentials")
		case <-ticker.C:
		}
	}
}

// AuthClient returns a client configured with TLS client certificates when type == "cert".
func (s *spiffeMethod) AuthClient(client *api.Client) (*api.Client, error) {
	s.logger.Trace("deriving auth client to use")

	if s.authType != authTypeCert {
		return client, nil
	}

	config := api.DefaultConfig()
	if config.Error != nil {
		return nil, config.Error
	}
	config.Address = client.Address()

	if err := config.ConfigureTLS(s.tlsCfg); err != nil {
		return nil, fmt.Errorf("error configuring TLS: %w", err)
	}

	hash, err := s.hashCert(s.tlsCfg.ClientCert, s.tlsCfg.ClientKey)
	if err != nil {
		return nil, fmt.Errorf("error reading certificates: %w", err)
	}
	s.latestHash.Store(&hash)

	clientToAuth, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	if ns := client.Headers().Get(consts.NamespaceHeaderName); ns != "" {
		clientToAuth.SetNamespace(ns)
	}

	return clientToAuth, nil
}

func (s *spiffeMethod) readToken() (string, error) {
	raw, err := os.ReadFile(s.svidTokenFile)
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(raw)), nil
}

// runWatcher is the shared polling loop for both credential types. checkFn is
// called on each tick; it reads the current credentials, compares them to the
// last known value, and returns true when a change is detected.
func (s *spiffeMethod) runWatcher(readPeriod time.Duration, checkFn func() bool) {
	defer close(s.doneCh)

	select {
	case <-s.stopCh:
		return
	case <-s.credSuccessGate:
		// We only start the next loop once we're initially successful,
		// since at startup Authenticate will be called, and we don't want
		// to end up immediately re-authenticating by having found a new
		// value
	}

	ticker := time.NewTicker(readPeriod)
	defer ticker.Stop()

	for {
		changed := false
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			changed = checkFn()
		}

		if changed {
			select {
			case s.credsFound <- struct{}{}:
			case <-s.stopCh:
				return
			}
		}
	}
}

func (s *spiffeMethod) hashCert(certFile, keyFile string) ([]byte, error) {
	h := sha256.New()

	var (
		certPEM []byte
		keyPEM  []byte
		err     error
	)

	if certFile != "" {
		certPEM, err = os.ReadFile(certFile)
		if err != nil {
			return nil, err
		}
		h.Write(certPEM)
	}
	if keyFile != "" {
		keyPEM, err = os.ReadFile(keyFile)
		if err != nil {
			return nil, err
		}
		h.Write(keyPEM)
	}

	// Verify the cert/key pair
	if certFile != "" && keyFile != "" {
		if _, err := tls.X509KeyPair(certPEM, keyPEM); err != nil {
			return nil, err
		}
	}

	return h.Sum(nil), nil
}
