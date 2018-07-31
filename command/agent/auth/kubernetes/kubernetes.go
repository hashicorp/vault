package kubernetes

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

const (
	serviceAccountFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
)

type kubernetesMethod struct {
	logger    hclog.Logger
	mountPath string

	role string
}

func NewKubernetesAuthMethod(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("empty config data")
	}

	k := &kubernetesMethod{
		logger:    conf.Logger,
		mountPath: conf.MountPath,
	}

	roleRaw, ok := conf.Config["role"]
	if !ok {
		return nil, errors.New("missing 'role' value")
	}
	k.role, ok = roleRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'role' config value to string")
	}
	if k.role == "" {
		return nil, errors.New("'role' value is empty")
	}

	return k, nil
}

func (k *kubernetesMethod) Authenticate(ctx context.Context, client *api.Client) (string, map[string]interface{}, error) {
	k.logger.Trace("beginning authentication")
	content, err := ioutil.ReadFile(serviceAccountFile)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s/login", k.mountPath), map[string]interface{}{
		"role": k.role,
		"jwt":  strings.TrimSpace(string(content)),
	}, nil
}

func (k *kubernetesMethod) NewCreds() chan struct{} {
	return nil
}

func (k *kubernetesMethod) CredSuccess() {
}

func (k *kubernetesMethod) Shutdown() {
}
