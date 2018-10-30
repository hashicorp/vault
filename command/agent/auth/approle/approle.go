package approle

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/auth"
)

type approleMethod struct {
	logger    hclog.Logger
	mountPath string

	roleIDPath   string
	secretIDPath string
}

func NewApproleAuthMethod(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("empty config data")
	}

	a := &approleMethod{
		logger:    conf.Logger,
		mountPath: conf.MountPath,
	}

	roleIDPathRaw, ok := conf.Config["role_id_path"]
	if !ok {
		return nil, errors.New("missing 'role_id_path' value")
	}
	a.roleIDPath, ok = roleIDPathRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'role_id_path' config value to string")
	}
	if a.roleIDPath == "" {
		return nil, errors.New("'role_id_path' value is empty")
	}

	secretIDPathRaw, ok := conf.Config["secret_id_path"]
	if !ok {
		return nil, errors.New("missing 'secret_id_path' value")
	}
	a.secretIDPath, ok = secretIDPathRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'secret_id_path' config value to string")
	}
	if a.secretIDPath == "" {
		return nil, errors.New("'secret_id_path' value is empty")
	}

	return a, nil
}

func (a *approleMethod) Authenticate(ctx context.Context, client *api.Client) (string, map[string]interface{}, error) {
	a.logger.Trace("beginning authentication")
	roleID, err := ioutil.ReadFile(a.roleIDPath)
	if err != nil {
		log.Fatal(err)
	}
	secretID, err := ioutil.ReadFile(a.secretIDPath)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s/login", a.mountPath), map[string]interface{}{
		"role_id":   strings.TrimSpace(string(roleID)),
		"secret_id": strings.TrimSpace(string(secretID)),
	}, nil
}

func (a *approleMethod) NewCreds() chan struct{} {
	return nil
}

func (a *approleMethod) CredSuccess() {
}

func (a *approleMethod) Shutdown() {
}
