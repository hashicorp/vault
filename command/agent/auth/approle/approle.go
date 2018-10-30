package approle

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/auth"
	"github.com/hashicorp/vault/helper/parseutil"
)

type approleMethod struct {
	logger    hclog.Logger
	mountPath string

	roleIDFilePath                 string
	secretIDFilePath               string
	cachedRoleID                   string
	cachedSecretID                 string
	removeSecretIDFileAfterReading bool
}

func NewApproleAuthMethod(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("empty config data")
	}

	a := &approleMethod{
		logger:                         conf.Logger,
		mountPath:                      conf.MountPath,
		removeSecretIDFileAfterReading: true,
	}

	roleIDFilePathRaw, ok := conf.Config["role_id_file_path"]
	if !ok {
		return nil, errors.New("missing 'role_id_file_path' value")
	}
	a.roleIDFilePath, ok = roleIDFilePathRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'role_id_file_path' config value to string")
	}
	if a.roleIDFilePath == "" {
		return nil, errors.New("'role_id_file_path' value is empty")
	}

	secretIDFilePathRaw, ok := conf.Config["secret_id_file_path"]
	if !ok {
		return nil, errors.New("missing 'secret_id_file_path' value")
	}
	a.secretIDFilePath, ok = secretIDFilePathRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'secret_id_file_path' config value to string")
	}
	if a.secretIDFilePath == "" {
		return nil, errors.New("'secret_id_file_path' value is empty")
	}

	removeSecretIDFileAfterReadingRaw, ok := conf.Config["remove_secret_id_file_after_reading"]
	if ok {
		removeSecretIDFileAfterReading, err := parseutil.ParseBool(removeSecretIDFileAfterReadingRaw)
		if err != nil {
			return nil, errwrap.Wrapf("error parsing 'remove_secret_id_file_after_reading' value: {{err}}", err)
		}
		a.removeSecretIDFileAfterReading = removeSecretIDFileAfterReading
	}

	return a, nil
}

func (a *approleMethod) Authenticate(ctx context.Context, client *api.Client) (string, map[string]interface{}, error) {
	if _, err := os.Stat(a.roleIDFilePath); err == nil {
		roleID, err := ioutil.ReadFile(a.roleIDFilePath)
		if err != nil {
			if a.cachedRoleID == "" {
				return "", nil, errwrap.Wrapf("error reading role ID file and no cached role ID known: {{err}}", err)
			}
			a.logger.Warn("error reading role ID file", "error", err)
		}
		if len(roleID) == 0 {
			if a.cachedRoleID == "" {
				return "", nil, errors.New("role ID file empty and no cached role ID known")
			}
			a.logger.Warn("role ID file exists but read empty value, re-using cached value")
		} else {
			a.cachedRoleID = strings.TrimSpace(string(roleID))
		}
	}
	if _, err := os.Stat(a.secretIDFilePath); err == nil {
		secretID, err := ioutil.ReadFile(a.secretIDFilePath)
		if err != nil {
			if a.cachedSecretID == "" {
				return "", nil, errwrap.Wrapf("error reading secret ID file and no cached secret ID known: {{err}}", err)
			}
			a.logger.Warn("error reading secret ID file", "error", err)
		}
		if len(secretID) == 0 {
			if a.cachedSecretID == "" {
				return "", nil, errors.New("secret ID file empty and no cached secret ID known")
			}
			a.logger.Warn("secret ID file exists but read empty value, re-using cached value")
		} else {
			a.cachedSecretID = strings.TrimSpace(string(secretID))
			if a.removeSecretIDFileAfterReading {
				if err := os.Remove(a.secretIDFilePath); err != nil {
					a.logger.Error("error removing secret ID file after reading", "error", err)
				}
			}
		}
	}

	if a.cachedRoleID == "" {
		return "", nil, errors.New("no known role ID")
	}
	if a.cachedSecretID == "" {
		return "", nil, errors.New("no known secret ID")
	}

	return fmt.Sprintf("%s/login", a.mountPath), map[string]interface{}{
		"role_id":   a.cachedRoleID,
		"secret_id": a.cachedSecretID,
	}, nil
}

func (a *approleMethod) NewCreds() chan struct{} {
	return nil
}

func (a *approleMethod) CredSuccess() {
}

func (a *approleMethod) Shutdown() {
}
