package cf

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	cf "github.com/hashicorp/vault-plugin-auth-cf"
	"github.com/hashicorp/vault-plugin-auth-cf/signatures"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/auth"
)

type cfMethod struct {
	mountPath string
	roleName  string
}

func NewCFAuthMethod(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("empty config data")
	}
	a := &cfMethod{
		mountPath: conf.MountPath,
	}
	if raw, ok := conf.Config["role"]; ok {
		if roleName, ok := raw.(string); ok {
			a.roleName = roleName
		} else {
			return nil, errors.New("could not convert 'role' config value to string")
		}
	} else {
		return nil, errors.New("missing 'role' value")
	}
	return a, nil
}

func (p *cfMethod) Authenticate(ctx context.Context, client *api.Client) (string, http.Header, map[string]interface{}, error) {
	pathToClientCert := os.Getenv(cf.EnvVarInstanceCertificate)
	if pathToClientCert == "" {
		return "", nil, nil, fmt.Errorf("missing %q value", cf.EnvVarInstanceCertificate)
	}
	certBytes, err := ioutil.ReadFile(pathToClientCert)
	if err != nil {
		return "", nil, nil, err
	}
	pathToClientKey := os.Getenv(cf.EnvVarInstanceKey)
	if pathToClientKey == "" {
		return "", nil, nil, fmt.Errorf("missing %q value", cf.EnvVarInstanceKey)
	}
	signingTime := time.Now().UTC()
	signatureData := &signatures.SignatureData{
		SigningTime:            signingTime,
		Role:                   p.roleName,
		CFInstanceCertContents: string(certBytes),
	}
	signature, err := signatures.Sign(pathToClientKey, signatureData)
	if err != nil {
		return "", nil, nil, err
	}
	data := map[string]interface{}{
		"role":             p.roleName,
		"cf_instance_cert": string(certBytes),
		"signing_time":     signingTime.Format(signatures.TimeFormat),
		"signature":        signature,
	}
	return fmt.Sprintf("%s/login", p.mountPath), nil, data, nil
}

func (p *cfMethod) NewCreds() chan struct{} {
	return nil
}

func (p *cfMethod) CredSuccess() {}

func (p *cfMethod) Shutdown() {}
