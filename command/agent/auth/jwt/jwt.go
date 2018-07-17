package jwt

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-hclog"
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
	watcher         *fsnotify.Watcher
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

	var err error
	j.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, errwrap.Wrapf("error creating watcher: {{err}}", err)
	}
	err = j.watcher.Add(j.path)
	if err != nil {
		return nil, errwrap.Wrapf("error adding path to watcher: {{err}}", err)
	}

	go j.runWatcher()

	return j, nil
}

func (j *jwtMethod) Authenticate(client *api.Client) (*api.Secret, error) {
	j.ingressToken()

	latestToken := j.latestToken.Load().(string)
	if latestToken == "" {
		return nil, errors.New("latest known jwt is empty, cannot authenticate")
	}

	secret, err := client.Logical().Write(fmt.Sprintf("%s/login", j.mountPath), map[string]interface{}{
		"role": j.role,
		"jwt":  latestToken,
	})

	if err != nil {
		return nil, errwrap.Wrapf("error logging in: {{err}}", err)
	}

	return secret, nil
}

func (j *jwtMethod) NewCreds() chan struct{} {
	return j.credsFound
}

func (j *jwtMethod) Shutdown() {
	j.watcher.Close()
	close(j.stopCh)
	<-j.doneCh
}

func (j *jwtMethod) runWatcher() {
	defer close(j.doneCh)

	// Drain the watcher in case events have been queueing up
drainloop:
	for {
		select {
		case <-j.watcher.Errors:
		case <-j.watcher.Events:
		case <-j.stopCh:
			return

		case <-j.credSuccessGate:
			// We only start the next loop once we're initially successful,
			// since at startup Authenticate will be called and we don't want
			// to end up immediately reauthenticating by having found a new
			// value
			break drainloop
		}
	}

	for {
		select {
		case <-j.stopCh:
			return

		case err := <-j.watcher.Errors:
			j.logger.Error("error from watcher", "error", err)

		case event := <-j.watcher.Events:
			switch event.Op {
			case fsnotify.Create, fsnotify.Write:
				latestToken := j.latestToken.Load().(string)
				j.ingressToken()
				newToken := j.latestToken.Load().(string)
				if newToken != latestToken {
					j.credsFound <- struct{}{}
				}
			}
		}
	}
}

func (j *jwtMethod) ingressToken() {
	j.logger.Info("ingressing jwt token file")

	fi, err := os.Lstat(j.path)
	if err != nil {
		if os.IsNotExist(err) {
			j.logger.Info("no current jwt file found, not updating")
			return
		}
		j.logger.Error("error encountered stat'ing jwt file", "error", err)
		return
	}
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
