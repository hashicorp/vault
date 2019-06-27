package pcf

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	pcf "github.com/hashicorp/vault-plugin-auth-pcf"
	"github.com/hashicorp/vault-plugin-auth-pcf/signatures"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/auth"
)

type pcfMethod struct {
	mountPath string
	roleName  string
}

func NewPCFAuthMethod(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("empty config data")
	}
	a := &pcfMethod{
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

func (p *pcfMethod) Authenticate(ctx context.Context, client *api.Client) (string, map[string]interface{}, error) {
	pathToClientCert := os.Getenv(pcf.EnvVarInstanceCertificate)
	if pathToClientCert == "" {
		return "", nil, fmt.Errorf("missing %q value", pcf.EnvVarInstanceCertificate)
	}
	certBytes, err := ioutil.ReadFile(pathToClientCert)
	if err != nil {
		return "", nil, err
	}
	pathToClientKey := os.Getenv(pcf.EnvVarInstanceKey)
	if pathToClientKey == "" {
		return "", nil, fmt.Errorf("missing %q value", pcf.EnvVarInstanceKey)
	}
	signingTime := time.Now().UTC()
	signatureData := &signatures.SignatureData{
		SigningTime:            signingTime,
		Role:                   p.roleName,
		CFInstanceCertContents: string(certBytes),
	}
	signature, err := signatures.Sign(pathToClientKey, signatureData)
	if err != nil {
		return "", nil, err
	}
	data := map[string]interface{}{
		"role":             p.roleName,
		"cf_instance_cert": string(certBytes),
		"signing_time":     signingTime.Format(signatures.TimeFormat),
		"signature":        signature,
	}
	return fmt.Sprintf("%s/login", p.mountPath), data, nil
}

func (p *pcfMethod) NewCreds() chan struct{} {
	return nil
}

func (p *pcfMethod) CredSuccess() {}

func (p *pcfMethod) Shutdown() {}
