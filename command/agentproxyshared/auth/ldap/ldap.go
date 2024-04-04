// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ldap

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agentproxyshared/auth"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
)

type ldapMethod struct {
	logger    hclog.Logger
	mountPath string

	username                      string
	passwordFilePath              string
	removePasswordAfterReading    bool
	removePasswordFollowsSymlinks bool
	credsFound                    chan struct{}
	watchCh                       chan string
	stopCh                        chan struct{}
	doneCh                        chan struct{}
	credSuccessGate               chan struct{}
	ticker                        *time.Ticker
	once                          *sync.Once
	latestPass                    *atomic.Value
}

// NewLdapMethod reads the user configuration and returns a configured
// LdapAuthMethod
func NewLdapAuthMethod(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("empty config data")
	}

	k := &ldapMethod{
		logger:                     conf.Logger,
		mountPath:                  conf.MountPath,
		removePasswordAfterReading: true,
		credsFound:                 make(chan struct{}),
		watchCh:                    make(chan string),
		stopCh:                     make(chan struct{}),
		doneCh:                     make(chan struct{}),
		credSuccessGate:            make(chan struct{}),
		once:                       new(sync.Once),
		latestPass:                 new(atomic.Value),
	}

	k.latestPass.Store("")
	usernameRaw, ok := conf.Config["username"]
	if !ok {
		return nil, errors.New("missing 'username' value")
	}
	k.username, ok = usernameRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'username' config value to string")
	}

	passFilePathRaw, ok := conf.Config["password_file_path"]
	if !ok {
		return nil, errors.New("missing 'password_file_path' value")
	}
	k.passwordFilePath, ok = passFilePathRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'password_file_path' config value to string")
	}
	if removePassAfterReadingRaw, ok := conf.Config["remove_password_after_reading"]; ok {
		removePassAfterReading, err := parseutil.ParseBool(removePassAfterReadingRaw)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'remove_password_after_reading' value: %w", err)
		}
		k.removePasswordAfterReading = removePassAfterReading
	}

	if removePassFollowsSymlinksRaw, ok := conf.Config["remove_password_follows_symlinks"]; ok {
		removePassFollowsSymlinks, err := parseutil.ParseBool(removePassFollowsSymlinksRaw)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'remove_password_follows_symlinks' value: %w", err)
		}
		k.removePasswordFollowsSymlinks = removePassFollowsSymlinks
	}
	switch {
	case k.passwordFilePath == "":
		return nil, errors.New("'password_file_path' value is empty")
	case k.username == "":
		return nil, errors.New("'username' value is empty")
	}

	// Default readPeriod
	readPeriod := 1 * time.Minute

	if passReadPeriodRaw, ok := conf.Config["password_read_period"]; ok {
		passReadPeriod, err := parseutil.ParseDurationSecond(passReadPeriodRaw)
		if err != nil || passReadPeriod <= 0 {
			return nil, fmt.Errorf("error parsing 'password_read_period' value into a positive value: %w", err)
		}
		readPeriod = passReadPeriod
	} else {
		// If we don't delete the password after reading, use a slower reload period,
		// otherwise we would re-read the whole file every 500ms, instead of just
		// doing a stat on the file every 500ms.
		if k.removePasswordAfterReading {
			readPeriod = 500 * time.Millisecond
		}
	}

	k.ticker = time.NewTicker(readPeriod)

	go k.runWatcher()

	k.logger.Info("ldap auth method created", "password_file_path", k.passwordFilePath)

	return k, nil
}

func (k *ldapMethod) Authenticate(ctx context.Context, client *api.Client) (string, http.Header, map[string]interface{}, error) {
	k.logger.Trace("beginning authentication")

	k.ingressPass()

	latestPass := k.latestPass.Load().(string)

	if latestPass == "" {
		return "", nil, nil, errors.New("latest known password is empty, cannot authenticate")
	}
	k.logger.Info("last known password in Authentication setup is")
	return fmt.Sprintf("%s/login/%s", k.mountPath, k.username), nil, map[string]interface{}{
		"password": latestPass,
	}, nil
}

func (k *ldapMethod) NewCreds() chan struct{} {
	return k.credsFound
}

func (k *ldapMethod) CredSuccess() {
	k.once.Do(func() {
		close(k.credSuccessGate)
	})
}

func (k *ldapMethod) Shutdown() {
	k.ticker.Stop()
	close(k.stopCh)
	<-k.doneCh
}

func (k *ldapMethod) runWatcher() {
	defer close(k.doneCh)

	select {
	case <-k.stopCh:
		return

	case <-k.credSuccessGate:
		// We only start the next loop once we're initially successful,
		// since at startup Authenticate will be called, and we don't want
		// to end up immediately re-authenticating by having found a new
		// value
	}

	for {
		select {
		case <-k.stopCh:
			return

		case <-k.ticker.C:
			latestPass := k.latestPass.Load().(string)
			k.ingressPass()
			newPass := k.latestPass.Load().(string)
			if newPass != latestPass {
				k.logger.Debug("new password file found")
				k.credsFound <- struct{}{}
			}
		}
	}
}

func (k *ldapMethod) ingressPass() {
	fi, err := os.Lstat(k.passwordFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		k.logger.Error("error encountered stat'ing password file", "error", err)
		return
	}

	// Check that the path refers to a file.
	// If it's a symlink, it could still be a symlink to a directory,
	// but os.ReadFile below will return a descriptive error.
	evalSymlinkPath := k.passwordFilePath
	switch mode := fi.Mode(); {
	case mode.IsRegular():
		// regular file
	case mode&fs.ModeSymlink != 0:
		// If our file path is a symlink, we should also return early (like above) without error
		// if the file that is linked to is not present, otherwise we will error when trying
		// to read that file by following the link in the os.ReadFile call.
		evalSymlinkPath, err = filepath.EvalSymlinks(k.passwordFilePath)
		if err != nil {
			k.logger.Error("error encountered evaluating symlinks", "error", err)
			return
		}
		_, err := os.Stat(evalSymlinkPath)
		if err != nil {
			if os.IsNotExist(err) {
				return
			}
			k.logger.Error("error encountered stat'ing password file after evaluating symlinks", "error", err)
			return
		}
	default:
		k.logger.Error("password file is not a regular file or symlink")
		return
	}

	pass, err := os.ReadFile(k.passwordFilePath)
	if err != nil {
		k.logger.Error("failed to read password file", "error", err)
		return
	}

	switch len(pass) {
	case 0:
		k.logger.Warn("empty password file read")

	default:
		k.latestPass.Store(string(pass))
	}

	if k.removePasswordAfterReading {
		pathToRemove := k.passwordFilePath
		if k.removePasswordFollowsSymlinks {
			// If removePassFollowsSymlinks is set, we follow the symlink and delete the password,
			// not just the symlink that links to the password file
			pathToRemove = evalSymlinkPath
		}
		if err := os.Remove(pathToRemove); err != nil {
			k.logger.Error("error removing password file", "error", err)
		}
	}
}
