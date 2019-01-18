package aws

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/hashicorp/errwrap"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	hclog "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	awsauth "github.com/hashicorp/vault/builtin/credential/aws"
	"github.com/hashicorp/vault/command/agent/auth"
)

const (
	typeEC2          = "ec2"
	typeIAM          = "iam"
	identityEndpoint = "http://169.254.169.254/latest/dynamic/instance-identity"

	/*

		IAM creds can be inferred from instance metadata or the container
		identity service, and those creds expire at varying intervals with
		new creds becoming available at likewise varying intervals. Let's
		default to polling once a minute so all changes can be picked up
		rather quickly. This is configurable, however.

	*/
	defaultCredentialPollInterval = 60
)

type awsMethod struct {
	logger      hclog.Logger
	authType    string
	nonce       string
	mountPath   string
	role        string
	headerValue string

	// These are used to share the latest creds safely across goroutines.
	credLock  sync.Mutex
	lastCreds *credentials.Credentials

	// Notifies the outer environment that it should call Authenticate again.
	credsFound chan struct{}

	// Detects that the outer environment is closing.
	stopCh chan struct{}
}

func NewAWSAuthMethod(conf *auth.AuthConfig) (auth.AuthMethod, error) {
	if conf == nil {
		return nil, errors.New("empty config")
	}
	if conf.Config == nil {
		return nil, errors.New("empty config data")
	}

	a := &awsMethod{
		logger:     conf.Logger,
		mountPath:  conf.MountPath,
		credsFound: make(chan struct{}),
		stopCh:     make(chan struct{}),
	}

	typeRaw, ok := conf.Config["type"]
	if !ok {
		return nil, errors.New("missing 'type' value")
	}
	a.authType, ok = typeRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'type' config value to string")
	}

	roleRaw, ok := conf.Config["role"]
	if !ok {
		return nil, errors.New("missing 'role' value")
	}
	a.role, ok = roleRaw.(string)
	if !ok {
		return nil, errors.New("could not convert 'role' config value to string")
	}

	switch {
	case a.role == "":
		return nil, errors.New("'role' value is empty")
	case a.authType == "":
		return nil, errors.New("'type' value is empty")
	case a.authType != typeEC2 && a.authType != typeIAM:
		return nil, errors.New("'type' value is invalid")
	}

	accessKey := ""
	accessKeyRaw, ok := conf.Config["access_key"]
	if ok {
		accessKey, ok = accessKeyRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'access_key' value into string")
		}
	}

	secretKey := ""
	secretKeyRaw, ok := conf.Config["secret_key"]
	if ok {
		secretKey, ok = secretKeyRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'secret_key' value into string")
		}
	}

	sessionToken := ""
	sessionTokenRaw, ok := conf.Config["session_token"]
	if ok {
		sessionToken, ok = sessionTokenRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'session_token' value into string")
		}
	}

	headerValueRaw, ok := conf.Config["header_value"]
	if ok {
		a.headerValue, ok = headerValueRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'header_value' value into string")
		}
	}

	if a.authType == typeIAM {

		// Check for an optional custom frequency at which we should poll for creds.
		credentialPollIntervalSec := defaultCredentialPollInterval
		if credentialPollIntervalRaw, ok := conf.Config["credential_poll_interval"]; ok {
			if credentialPollInterval, ok := credentialPollIntervalRaw.(int); ok {
				credentialPollIntervalSec = credentialPollInterval
			} else {
				return nil, errors.New("could not convert 'credential_poll_interval' into int")
			}
		}

		// Do an initial population of the creds because we want to err right away if we can't
		// even get a first set.
		creds, err := awsauth.RetrieveCreds(accessKey, secretKey, sessionToken)
		if err != nil {
			return nil, err
		}
		a.lastCreds = creds

		go a.pollForCreds(accessKey, secretKey, sessionToken, credentialPollIntervalSec)
	}

	return a, nil
}

func (a *awsMethod) Authenticate(ctx context.Context, client *api.Client) (retToken string, retData map[string]interface{}, retErr error) {
	a.logger.Trace("beginning authentication")

	data := make(map[string]interface{})

	switch a.authType {
	case typeEC2:
		client := cleanhttp.DefaultClient()

		// Fetch document
		{
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/document", identityEndpoint), nil)
			if err != nil {
				retErr = errwrap.Wrapf("error creating request: {{err}}", err)
				return
			}
			req = req.WithContext(ctx)
			resp, err := client.Do(req)
			if err != nil {
				retErr = errwrap.Wrapf("error fetching instance document: {{err}}", err)
				return
			}
			if resp == nil {
				retErr = errors.New("empty response fetching instance document")
				return
			}
			defer resp.Body.Close()
			doc, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				retErr = errwrap.Wrapf("error reading instance document response body: {{err}}", err)
				return
			}
			data["identity"] = base64.StdEncoding.EncodeToString(doc)
		}

		// Fetch signature
		{
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/signature", identityEndpoint), nil)
			if err != nil {
				retErr = errwrap.Wrapf("error creating request: {{err}}", err)
				return
			}
			req = req.WithContext(ctx)
			resp, err := client.Do(req)
			if err != nil {
				retErr = errwrap.Wrapf("error fetching instance document signature: {{err}}", err)
				return
			}
			if resp == nil {
				retErr = errors.New("empty response fetching instance document signature")
				return
			}
			defer resp.Body.Close()
			sig, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				retErr = errwrap.Wrapf("error reading instance document signature response body: {{err}}", err)
				return
			}
			data["signature"] = string(sig)
		}

		// Add the reauthentication value, if we have one
		if a.nonce == "" {
			uid, err := uuid.GenerateUUID()
			if err != nil {
				retErr = errwrap.Wrapf("error generating uuid for reauthentication value: {{err}}", err)
				return
			}
			a.nonce = uid
		}
		data["nonce"] = a.nonce

	default:
		// This is typeIAM.
		a.credLock.Lock()
		defer a.credLock.Unlock()

		var err error
		data, err = awsauth.GenerateLoginData(a.lastCreds, a.headerValue)
		if err != nil {
			retErr = errwrap.Wrapf("error creating login value: {{err}}", err)
			return
		}
	}

	data["role"] = a.role

	return fmt.Sprintf("%s/login", a.mountPath), data, nil
}

func (a *awsMethod) NewCreds() chan struct{} {
	return a.credsFound
}

func (a *awsMethod) CredSuccess() {}

func (a *awsMethod) Shutdown() {
	close(a.credsFound)
	close(a.stopCh)
}

func (a *awsMethod) pollForCreds(accessKey, secretKey, sessionToken string, frequencySeconds int) {
	ticker := time.NewTicker(time.Duration(frequencySeconds) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-a.stopCh:
			a.logger.Trace("shutdown triggered, stopping aws auth handler")
			return
		case <-ticker.C:
			if err := a.checkCreds(accessKey, secretKey, sessionToken); err != nil {
				a.logger.Warn("unable to retrieve current creds, retaining last creds", err)
			}
		}
	}
}

func (a *awsMethod) checkCreds(accessKey, secretKey, sessionToken string) error {
	a.credLock.Lock()
	defer a.credLock.Unlock()

	a.logger.Trace("checking for new credentials")
	currentCreds, err := awsauth.RetrieveCreds(accessKey, secretKey, sessionToken)
	if err != nil {
		return err
	}

	currentVal, err := currentCreds.Get()
	if err != nil {
		return err
	}
	lastVal, err := a.lastCreds.Get()
	if err != nil {
		return err
	}

	// These will always have different pointers regardless of whether their
	// values are identical, hence the use of DeepEqual.
	if !a.lastCreds.IsExpired() && reflect.DeepEqual(currentVal, lastVal) {
		a.logger.Trace("credentials are unchanged and still valid")
		return nil
	}

	a.lastCreds = currentCreds
	a.logger.Trace("new credentials detected, triggering Authenticate")
	a.credsFound <- struct{}{}
	return nil
}
