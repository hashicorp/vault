package jwt

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"sync/atomic"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/auth"
)

type jwtMethod struct {
	logger          hclog.Logger
	path            string
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

func NewJWTAuthMethod(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("empty config data")
	}

	j := &jwtMethod{
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

	pathRaw, ok := conf.Config["path"]
	if !ok {
		return nil, errors.New("missing 'path' value")
	}
	j.path, ok = pathRaw.(string)
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
	case j.path == "":
		return nil, errors.New("'path' value is empty")
	case j.role == "":
		return nil, errors.New("'role' value is empty")
	}

	j.ticker = time.NewTicker(500 * time.Millisecond)

	go j.runWatcher()

	j.logger.Info("jwt auth method created", "path", j.path)

	return j, nil
}

func (j *jwtMethod) Authenticate(_ context.Context, client *api.Client) (string, map[string]interface{}, error) {
	j.logger.Trace("beginning authentication")

	j.ingressToken()

	latestToken := j.latestToken.Load().(string)
	if latestToken == "" {
		return "", nil, errors.New("latest known jwt is empty, cannot authenticate")
	}

	return fmt.Sprintf("%s/login", j.mountPath), map[string]interface{}{
		"role": j.role,
		"jwt":  latestToken,
	}, nil
}

func (j *jwtMethod) NewCreds() chan struct{} {
	return j.credsFound
}

func (j *jwtMethod) CredSuccess() {
	j.once.Do(func() {
		close(j.credSuccessGate)
	})
}

func (j *jwtMethod) Shutdown() {
	j.ticker.Stop()
	close(j.stopCh)
	<-j.doneCh
}

func (j *jwtMethod) runWatcher() {
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

func (j *jwtMethod) ingressToken() {
	fi, err := os.Lstat(j.path)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		j.logger.Error("error encountered stat'ing jwt file", "error", err)
		return
	}

	j.logger.Debug("new jwt file found")

	if !fi.Mode().IsRegular() {
		j.logger.Error("jwt file is not a regular file")
		return
	}

	token, err := ioutil.ReadFile(j.path)
	if err != nil {
		j.logger.Error("failed to read jwt file", "error", err)
		return
	}

	switch len(token) {
	case 0:
		j.logger.Warn("empty jwt file read")

	default:
		j.latestToken.Store(string(token))
	}

	if err := os.Remove(j.path); err != nil {
		j.logger.Error("error removing jwt file", "error", err)
	}
}
