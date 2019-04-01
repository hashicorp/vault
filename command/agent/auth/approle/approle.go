package approle

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hashicorp/errwrap"
	hclog "github.com/hashicorp/go-hclog"
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
	secretIDResponseWrappingPath   string
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
	if ok {
		a.secretIDFilePath, ok = secretIDFilePathRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'secret_id_file_path' config value to string")
		}
		if a.secretIDFilePath == "" {
			return a, nil
		}

		removeSecretIDFileAfterReadingRaw, ok := conf.Config["remove_secret_id_file_after_reading"]
		if ok {
			removeSecretIDFileAfterReading, err := parseutil.ParseBool(removeSecretIDFileAfterReadingRaw)
			if err != nil {
				return nil, errwrap.Wrapf("error parsing 'remove_secret_id_file_after_reading' value: {{err}}", err)
			}
			a.removeSecretIDFileAfterReading = removeSecretIDFileAfterReading
		}

		secretIDResponseWrappingPathRaw, ok := conf.Config["secret_id_response_wrapping_path"]
		if ok {
			a.secretIDResponseWrappingPath, ok = secretIDResponseWrappingPathRaw.(string)
			if !ok {
				return nil, errors.New("could not convert 'secret_id_response_wrapping_path' config value to string")
			}
			if a.secretIDResponseWrappingPath == "" {
				return nil, errors.New("'secret_id_response_wrapping_path' value is empty")
			}
		}
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

	if a.cachedRoleID == "" {
		return "", nil, errors.New("no known role ID")
	}

	if a.secretIDFilePath == "" {
		return fmt.Sprintf("%s/login", a.mountPath), map[string]interface{}{
			"role_id": a.cachedRoleID,
		}, nil
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
			stringSecretID := strings.TrimSpace(string(secretID))
			if a.secretIDResponseWrappingPath != "" {
				clonedClient, err := client.Clone()
				if err != nil {
					return "", nil, errwrap.Wrapf("error cloning client to unwrap secret ID: {{err}}", err)
				}
				clonedClient.SetToken(stringSecretID)
				// Validate the creation path
				resp, err := clonedClient.Logical().Read("sys/wrapping/lookup")
				if err != nil {
					return "", nil, errwrap.Wrapf("error looking up wrapped secret ID: {{err}}", err)
				}
				if resp == nil {
					return "", nil, errors.New("response nil when looking up wrapped secret ID")
				}
				if resp.Data == nil {
					return "", nil, errors.New("data in response nil when looking up wrapped secret ID")
				}
				creationPathRaw, ok := resp.Data["creation_path"]
				if !ok {
					return "", nil, errors.New("creation_path in response nil when looking up wrapped secret ID")
				}
				creationPath, ok := creationPathRaw.(string)
				if !ok {
					return "", nil, errors.New("creation_path in response could not be parsed as string when looking up wrapped secret ID")
				}
				if creationPath != a.secretIDResponseWrappingPath {
					a.logger.Error("SECURITY: unable to validate wrapping token creation path", "expected", a.secretIDResponseWrappingPath, "found", creationPath)
					return "", nil, errors.New("unable to validate wrapping token creation path")
				}
				// Now get the secret ID
				resp, err = clonedClient.Logical().Unwrap("")
				if err != nil {
					return "", nil, errwrap.Wrapf("error unwrapping secret ID: {{err}}", err)
				}
				if resp == nil {
					return "", nil, errors.New("response nil when unwrapping secret ID")
				}
				if resp.Data == nil {
					return "", nil, errors.New("data in response nil when unwrapping secret ID")
				}
				secretIDRaw, ok := resp.Data["secret_id"]
				if !ok {
					return "", nil, errors.New("secret_id in response nil when unwrapping secret ID")
				}
				secretID, ok := secretIDRaw.(string)
				if !ok {
					return "", nil, errors.New("secret_id in response could not be parsed as string when unwrapping secret ID")
				}
				stringSecretID = secretID
			}
			a.cachedSecretID = stringSecretID
			if a.removeSecretIDFileAfterReading {
				if err := os.Remove(a.secretIDFilePath); err != nil {
					a.logger.Error("error removing secret ID file after reading", "error", err)
				}
			}
		}
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
