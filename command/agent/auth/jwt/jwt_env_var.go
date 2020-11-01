package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/auth"
)

type jwtMethodFromEnvVar struct {
	logger          hclog.Logger
	envVar          string
	mountPath       string
	role            string
	credsFound      chan struct{}
	watchCh         chan string
	stopCh          chan struct{}
	doneCh          chan struct{}
	credSuccessGate chan struct{}
	ticker          *time.Ticker
	once            *sync.Once
	latestToken     *atomic.Value
}

func newJWTAuthMethodFromEnvVar(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	j := &jwtMethodFromEnvVar{
		logger:          conf.Logger,
		mountPath:       conf.MountPath,
		credsFound:      make(chan struct{}),
		watchCh:         make(chan string),
		stopCh:          make(chan struct{}),
		doneCh:          make(chan struct{}),
		credSuccessGate: make(chan struct{}),
		once:            new(sync.Once),
		latestToken:     new(atomic.Value),
	}
	j.latestToken.Store("")

	envVarRaw, ok := conf.Config["env-var"]
	if !ok {
		return nil, errors.New("missing 'env-var' value")
	}
	j.envVar, ok = envVarRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'path' config value to string")
	}

	roleRaw, ok := conf.Config["role"]
	if !ok {
		return nil, errors.New("missing 'role' value")
	}
	j.role, ok = roleRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'role' config value to string")
	}

	switch {
	case j.envVar == "":
		return nil, errors.New("'env-var' value is empty")
	case j.role == "":
		return nil, errors.New("'role' value is empty")
	}

	j.ticker = time.NewTicker(500 * time.Millisecond)

	go j.runWatcher()

	return j, nil
}

func (j *jwtMethodFromEnvVar) Authenticate(_ context.Context, client *api.Client) (string, http.Header, map[string]interface{}, error) {
	j.logger.Trace("beginning authentication")

	j.ingressToken()

	latestToken := j.latestToken.Load().(string)
	if latestToken == "" {
		return "", nil, nil, errors.New("latest known jwt is empty, cannot authenticate")
	}

	return fmt.Sprintf("%s/login", j.mountPath), nil, map[string]interface{}{
		"role": j.role,
		"jwt":  latestToken,
	}, nil
}

func (j *jwtMethodFromEnvVar) NewCreds() chan struct{} {
	return j.credsFound
}

func (j *jwtMethodFromEnvVar) CredSuccess() {
	j.once.Do(func() {
		close(j.credSuccessGate)
	})
}

func (j *jwtMethodFromEnvVar) Shutdown() {
	j.ticker.Stop()
	close(j.stopCh)
	<-j.doneCh
}

func (j *jwtMethodFromEnvVar) runWatcher() {
	defer close(j.doneCh)

	select {
	case <-j.stopCh:
		return

	case <-j.credSuccessGate:
		// We only start the next loop once we're initially successful,
		// since at startup Authenticate will be called and we don't want
		// to end up immediately reauthenticating by having found a new
		// value
	}

	for {
		select {
		case <-j.stopCh:
			return

		case <-j.ticker.C:
			latestToken := j.latestToken.Load().(string)
			j.ingressToken()
			newToken := j.latestToken.Load().(string)
			if newToken != latestToken {
				j.credsFound <- struct{}{}
			}
		}
	}
}

func (j *jwtMethodFromEnvVar) ingressToken() {
	token, ok := os.LookupEnv(j.envVar)
	if !ok {
		return
	}

	debugMsg := fmt.Sprintf("new jwt value from environment variable %s is found", j.envVar)
	j.logger.Debug(debugMsg)

	switch len(token) {
	case 0:
		j.logger.Warn("empty jwt value read")
	default:
		j.latestToken.Store(token)
	}

	err := os.Unsetenv(j.envVar)
	if err != nil {
		j.logger.Error("error unsetting env-var", "error", err)
	}
}
