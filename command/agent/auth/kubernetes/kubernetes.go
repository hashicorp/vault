package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	hclog "github.com/hashicorp/go-hclog"
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

	// tokenPath is an optional path to a projected service account token inside
	// the pod, for use instead of the default service account token.
	tokenPath string

	// jwtData is a ReadCloser used to inject a ReadCloser for mocking tests.
	jwtData io.ReadCloser
}

// NewKubernetesAuthMethod reads the user configuration and returns a configured
// AuthMethod
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

	tokenPathRaw, ok := conf.Config["token_path"]
	if ok {
		k.tokenPath, ok = tokenPathRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'token_path' config value to string")
		}
	}

	if k.role == "" {
		return nil, errors.New("'role' value is empty")
	}

	return k, nil
}

func (k *kubernetesMethod) Authenticate(ctx context.Context, client *api.Client) (string, http.Header, map[string]interface{}, error) {
	k.logger.Trace("beginning authentication")

	jwtString, err := k.readJWT()
	if err != nil {
		return "", nil, nil, fmt.Errorf("error reading JWT with Kubernetes Auth: %w", err)
	}

	return fmt.Sprintf("%s/login", k.mountPath), nil, map[string]interface{}{
		"role": k.role,
		"jwt":  jwtString,
	}, nil
}

func (k *kubernetesMethod) NewCreds() chan struct{} {
	return nil
}

func (k *kubernetesMethod) CredSuccess() {
}

func (k *kubernetesMethod) Shutdown() {
}

// readJWT reads the JWT data for the Agent to submit to Vault. The default is
// to read the JWT from the default service account location, defined by the
// constant serviceAccountFile. In normal use k.jwtData is nil at invocation and
// the method falls back to reading the token path with os.Open, opening a file
// from either the default location or from the token_path path specified in
// configuration.
func (k *kubernetesMethod) readJWT() (string, error) {
	// load configured token path if set, default to serviceAccountFile
	tokenFilePath := serviceAccountFile
	if k.tokenPath != "" {
		tokenFilePath = k.tokenPath
	}

	data := k.jwtData
	// k.jwtData should only be non-nil in tests
	if data == nil {
		f, err := os.Open(tokenFilePath)
		if err != nil {
			return "", err
		}
		data = f
	}
	defer data.Close()

	contentBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(contentBytes)), nil
}
