package token_file

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/parseutil"

	"github.com/hashicorp/vault/api"

	"github.com/hashicorp/vault/command/agent/auth"

	"github.com/hashicorp/go-hclog"
)

type TokenFileMethod struct {
	logger    hclog.Logger
	mountPath string

	cachedToken                 string
	tokenFilePath               string
	followSymlinks              bool
	removeTokenFileAfterReading bool
}

func NewTokenFileAuthMethod(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("empty config data")
	}

	a := &TokenFileMethod{
		logger:    conf.Logger,
		mountPath: "auth/token",
	}

	tokenFilePathRaw, ok := conf.Config["token_file_path"]
	if !ok {
		return nil, errors.New("missing 'token_file_path' value")
	}
	a.tokenFilePath, ok = tokenFilePathRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'token_file_path' config value to string")
	}
	if a.tokenFilePath == "" {
		return nil, errors.New("'token_file_path' value is empty")
	}

	removeTokenFileAfterReadingRaw, ok := conf.Config["remove_token_file_after_reading"]
	if ok {
		removeTokenFileAfterReading, err := parseutil.ParseBool(removeTokenFileAfterReadingRaw)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'remove_token_file_after_reading' value: %w", err)
		}
		a.removeTokenFileAfterReading = removeTokenFileAfterReading
	}

	followSymlinksRaw, ok := conf.Config["follow_symlinks"]
	if ok {
		followSymlinks, err := parseutil.ParseBool(followSymlinksRaw)
		if err != nil {
			return nil, fmt.Errorf("error parsing 'follow_symlinks' value: %w", err)
		}
		a.followSymlinks = followSymlinks
	}

	return a, nil
}

func (a *TokenFileMethod) Authenticate(ctx context.Context, client *api.Client) (string, http.Header, map[string]interface{}, error) {
	_, err := os.Stat(a.tokenFilePath)
	if err != nil {
		if a.cachedToken == "" {
			return "", nil, nil, fmt.Errorf("token file not found and no cached token known: %w", err)
		}
		a.logger.Warn("token file not found", "error", err)
	}

	token, err := os.ReadFile(a.tokenFilePath)
	if err != nil {
		if a.cachedToken == "" {
			return "", nil, nil, fmt.Errorf("error reading token file and no cached token known: %w", err)
		}
		a.logger.Warn("error reading token file", "error", err)
	}
	if len(token) == 0 {
		if a.cachedToken == "" {
			return "", nil, nil, errors.New("token file empty and no cached token known")
		}
		a.logger.Warn("role ID file exists but read empty value, re-using cached value")
	} else {
		a.cachedToken = strings.TrimSpace(string(token))
	}

	if a.removeTokenFileAfterReading {
		a.logger.Info("removing token file after reading, because 'remove_token_file_after_reading' is true")
		if err := os.Remove(a.tokenFilePath); err != nil {
			a.logger.Error("error removing token file after reading", "error", err)
		}
	}

	// i.e. auth/token/lookup
	return fmt.Sprintf("%s/lookup", a.mountPath), nil, map[string]interface{}{
		"token": a.cachedToken,
	}, nil
}

func (a *TokenFileMethod) NewCreds() chan struct{} {
	return nil
}

func (a *TokenFileMethod) CredSuccess() {
}

func (a *TokenFileMethod) Shutdown() {
}
