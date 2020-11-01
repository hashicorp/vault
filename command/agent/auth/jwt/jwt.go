package jwt

import (
	"errors"

	"github.com/hashicorp/vault/command/agent/auth"
)

//NewJWTAuthMethod takes an auth configuration and returns an auth method and an error.
func NewJWTAuthMethod(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("empty config data")
	}

	_, pathExists := conf.Config["path"]
	_, envVarExists := conf.Config["env-var"]

	switch {
	case pathExists && envVarExists:
		return nil, errors.New("Either path or env-var can be used in config, not both")
	case pathExists:
		return newJWTAuthMethodFromFile(conf)
	case envVarExists:
		return newJWTAuthMethodFromEnvVar(conf)
	}

	return nil, errors.New("neither path nor env-var exist in config")
}
