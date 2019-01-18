package alicloud

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	aliCloudAuth "github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials/providers"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault-plugin-auth-alicloud/tools"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/auth"
)

/*

	Creds can be inferred from instance metadata, and those creds
	expire every 60 minutes, so we're going to need to poll for new
	creds. Since we're polling anyways, let's poll once a minute so
	all changes can be picked up rather quickly. This is configurable,
	however.

*/
const defaultCredCheckFreqSeconds = 60

func NewAliCloudAuthMethod(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("empty config data")
	}

	a := &alicloudMethod{
		logger:     conf.Logger,
		mountPath:  conf.MountPath,
		credsFound: make(chan struct{}),
		stopCh:     make(chan struct{}),
	}

	// Build the required information we'll need to create a client.
	if roleRaw, ok := conf.Config["role"]; !ok {
		return nil, errors.New("'role' is required but is not provided")
	} else {
		if a.role, ok = roleRaw.(string); !ok {
			return nil, errors.New("could not convert 'role' config value to string")
		}
	}
	if regionRaw, ok := conf.Config["region"]; !ok {
		return nil, errors.New("'region' is required but is not provided")
	} else {
		if a.region, ok = regionRaw.(string); !ok {
			return nil, errors.New("could not convert 'region' config value to string")
		}
	}

	// Check for an optional custom frequency at which we should poll for creds.
	credCheckFreqSec := defaultCredCheckFreqSeconds
	if checkFreqRaw, ok := conf.Config["credential_poll_interval"]; ok {
		if credFreq, ok := checkFreqRaw.(int); ok {
			credCheckFreqSec = credFreq
		} else {
			return nil, errors.New("could not convert 'credential_poll_interval' config value to int")
		}
	}

	// Build the optional, configuration-based piece of the credential chain.
	credConfig := &providers.Configuration{}

	if accessKeyRaw, ok := conf.Config["access_key"]; ok {
		if credConfig.AccessKeyID, ok = accessKeyRaw.(string); !ok {
			return nil, errors.New("could not convert 'access_key' config value to string")
		}
	}

	if accessSecretRaw, ok := conf.Config["access_secret"]; ok {
		if credConfig.AccessKeySecret, ok = accessSecretRaw.(string); !ok {
			return nil, errors.New("could not convert 'access_secret' config value to string")
		}
	}

	if accessTokenRaw, ok := conf.Config["access_token"]; ok {
		if credConfig.AccessKeyStsToken, ok = accessTokenRaw.(string); !ok {
			return nil, errors.New("could not convert 'access_token' config value to string")
		}
	}

	if roleArnRaw, ok := conf.Config["role_arn"]; ok {
		if credConfig.RoleArn, ok = roleArnRaw.(string); !ok {
			return nil, errors.New("could not convert 'role_arn' config value to string")
		}
	}

	if roleSessionNameRaw, ok := conf.Config["role_session_name"]; ok {
		if credConfig.RoleSessionName, ok = roleSessionNameRaw.(string); !ok {
			return nil, errors.New("could not convert 'role_session_name' config value to string")
		}
	}

	if roleSessionExpirationRaw, ok := conf.Config["role_session_expiration"]; ok {
		if roleSessionExpiration, ok := roleSessionExpirationRaw.(int); !ok {
			return nil, errors.New("could not convert 'role_session_expiration' config value to int")
		} else {
			credConfig.RoleSessionExpiration = &roleSessionExpiration
		}
	}

	if privateKeyRaw, ok := conf.Config["private_key"]; ok {
		if credConfig.PrivateKey, ok = privateKeyRaw.(string); !ok {
			return nil, errors.New("could not convert 'private_key' config value to string")
		}
	}

	if publicKeyIdRaw, ok := conf.Config["public_key_id"]; ok {
		if credConfig.PublicKeyID, ok = publicKeyIdRaw.(string); !ok {
			return nil, errors.New("could not convert 'public_key_id' config value to string")
		}
	}

	if sessionExpirationRaw, ok := conf.Config["session_expiration"]; ok {
		if sessionExpiration, ok := sessionExpirationRaw.(int); !ok {
			return nil, errors.New("could not convert 'session_expiration' config value to int")
		} else {
			credConfig.SessionExpiration = &sessionExpiration
		}
	}

	if roleNameRaw, ok := conf.Config["role_name"]; ok {
		if credConfig.RoleName, ok = roleNameRaw.(string); !ok {
			return nil, errors.New("could not convert 'role_name' config value to string")
		}
	}

	credentialChain := []providers.Provider{
		providers.NewEnvCredentialProvider(),
		providers.NewConfigurationCredentialProvider(credConfig),
		providers.NewInstanceMetadataProvider(),
	}
	credProvider := providers.NewChainProvider(credentialChain)

	// Do an initial population of the creds because we want to err right away if we can't
	// even get a first set.
	lastCreds, err := credProvider.Retrieve()
	if err != nil {
		return nil, err
	}
	a.lastCreds = lastCreds

	go a.pollForCreds(credProvider, credCheckFreqSec)

	return a, nil
}

type alicloudMethod struct {
	logger    hclog.Logger
	mountPath string

	// These parameters are fed into building login data.
	role   string
	region string

	// These are used to share the latest creds safely across goroutines.
	credLock  sync.Mutex
	lastCreds aliCloudAuth.Credential

	// Notifies the outer environment that it should call Authenticate again.
	credsFound chan struct{}

	// Detects that the outer environment is closing.
	stopCh chan struct{}
}

func (a *alicloudMethod) Authenticate(context.Context, *api.Client) (string, map[string]interface{}, error) {
	a.credLock.Lock()
	defer a.credLock.Unlock()

	a.logger.Trace("beginning authentication")
	data, err := tools.GenerateLoginData(a.role, a.lastCreds, a.region)
	if err != nil {
		return "", nil, err
	}
	return fmt.Sprintf("%s/login", a.mountPath), data, nil
}

func (a *alicloudMethod) NewCreds() chan struct{} {
	return a.credsFound
}

func (a *alicloudMethod) CredSuccess() {}

func (a *alicloudMethod) Shutdown() {
	close(a.credsFound)
	close(a.stopCh)
}

func (a *alicloudMethod) pollForCreds(credProvider providers.Provider, frequencySeconds int) {
	ticker := time.NewTicker(time.Duration(frequencySeconds) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-a.stopCh:
			a.logger.Trace("shutdown triggered, stopping alicloud auth handler")
			return
		case <-ticker.C:
			if err := a.checkCreds(credProvider); err != nil {
				a.logger.Warn("unable to retrieve current creds, retaining last creds", err)
			}
		}
	}
}

func (a *alicloudMethod) checkCreds(credProvider providers.Provider) error {
	a.credLock.Lock()
	defer a.credLock.Unlock()

	a.logger.Trace("checking for new credentials")
	currentCreds, err := credProvider.Retrieve()
	if err != nil {
		return err
	}
	// These will always have different pointers regardless of whether their
	// values are identical, hence the use of DeepEqual.
	if reflect.DeepEqual(currentCreds, a.lastCreds) {
		a.logger.Trace("credentials are unchanged")
		return nil
	}
	a.lastCreds = currentCreds
	a.logger.Trace("new credentials detected, triggering Authenticate")
	a.credsFound <- struct{}{}
	return nil
}
