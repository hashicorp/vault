// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/awsutil/v2"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agentproxyshared/auth"
	internal_awsutilv2 "github.com/hashicorp/vault/internal/awsutil/v2"
)

const (
	typeEC2 = "ec2"
	typeIAM = "iam"

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
	region      string

	// These are used to share the latest AWS config (with credentials) safely
	// across goroutines.
	credLock sync.Mutex
	lastCfg  *aws.Config

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
		region:     awsutil.DefaultRegion,
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

	nonceRaw, ok := conf.Config["nonce"]
	if ok {
		a.nonce, ok = nonceRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'nonce' value into string")
		}
	}

	regionRaw, ok := conf.Config["region"]
	if ok {
		a.region, ok = regionRaw.(string)
		if !ok {
			return nil, errors.New("could not convert 'region' value into string")
		}
	}

	if a.authType == typeIAM {

		// Check for an optional custom frequency at which we should poll for creds.
		credentialPollIntervalSec := defaultCredentialPollInterval
		if credentialPollIntervalRaw, ok := conf.Config["credential_poll_interval"]; ok {
			if credentialPollInterval, ok := credentialPollIntervalRaw.(int); ok && credentialPollInterval > 0 {
				credentialPollIntervalSec = credentialPollInterval
			} else {
				return nil, errors.New("could not convert 'credential_poll_interval' into positive int")
			}
		}

		// Load the AWS config up front. For static credentials, values are set directly on the config.
		// For dynamic credentials (IMDS/ECS/env), the chain is probed immediately to preserve the
		// v1 fail-fast behavior; credentials are then refreshed by the poller.

		cfgOpts := []func(*config.LoadOptions) error{}
		if a.region != "" {
			cfgOpts = append(cfgOpts, config.WithRegion(a.region))
		}
		cfg, err := config.LoadDefaultConfig(context.Background(), cfgOpts...)
		if err != nil {
			return nil, err
		}

		if accessKey != "" && secretKey != "" {
			cfg.Credentials = aws.NewCredentialsCache(
				credentials.NewStaticCredentialsProvider(
					accessKey,
					secretKey,
					sessionToken,
				),
			)
		} else {
			// fail-fast: probe the dynamic credential chain (env vars,
			// shared config, IMDS/ECS) at startup so a misconfigured instance fails
			// immediately rather than on the first Authenticate call.
			initCtx, initCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer initCancel()
			if _, err := cfg.Credentials.Retrieve(initCtx); err != nil {
				return nil, fmt.Errorf("failed to retrieve initial credentials: %w", err)
			}
		}

		a.lastCfg = &cfg

		go a.pollForCreds(accessKey, secretKey, sessionToken, credentialPollIntervalSec)
	}

	return a, nil
}

func (a *awsMethod) Authenticate(ctx context.Context, client *api.Client) (retToken string, header http.Header, retData map[string]interface{}, retErr error) {
	a.logger.Trace("beginning authentication")

	data := make(map[string]interface{})

	switch a.authType {
	case typeEC2:
		metadataSvc := imds.New(imds.Options{})

		// Fetch document
		{
			docOutput, err := metadataSvc.GetDynamicData(ctx, &imds.GetDynamicDataInput{
				Path: "instance-identity/document",
			})
			if err != nil {
				retErr = fmt.Errorf("error requesting doc: %w", err)
				return
			}
			defer docOutput.Content.Close()
			doc, err := io.ReadAll(docOutput.Content)
			if err != nil {
				retErr = fmt.Errorf("error reading doc response: %w", err)
				return
			}
			data["identity"] = base64.StdEncoding.EncodeToString(doc)
		}

		// Fetch signature
		{
			signatureOutput, err := metadataSvc.GetDynamicData(ctx, &imds.GetDynamicDataInput{
				Path: "instance-identity/signature",
			})
			if err != nil {
				retErr = fmt.Errorf("error requesting signature: %w", err)
				return
			}

			defer signatureOutput.Content.Close()
			signature, err := io.ReadAll(signatureOutput.Content)
			if err != nil {
				retErr = fmt.Errorf("error reading signature response: %w", err)
				return
			}
			data["signature"] = string(signature)
		}

		// Add the reauthentication value, if we have one
		if a.nonce == "" {
			uid, err := uuid.GenerateUUID()
			if err != nil {
				retErr = fmt.Errorf("error generating uuid for reauthentication value: %w", err)
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
		data, err = internal_awsutilv2.GenerateLoginDataV2(ctx, a.lastCfg, a.headerValue, a.region, a.logger)
		if err != nil {
			retErr = fmt.Errorf("error creating login value: %w", err)
			return
		}
	}

	data["role"] = a.role

	return fmt.Sprintf("%s/login", a.mountPath), nil, data, nil
}

func (a *awsMethod) NewCreds() chan struct{} {
	return a.credsFound
}

func (a *awsMethod) CredSuccess() {}

func (a *awsMethod) Shutdown() {
	close(a.stopCh)
	// Note: credsFound is not closed to avoid panic if poller tries to send during shutdown.
	// The select statement in checkCreds handles shutdown coordination via stopCh.
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
				a.logger.Warn("unable to retrieve current creds, retaining last creds", "error", err)
			}
		}
	}
}

func (a *awsMethod) checkCreds(accessKey, secretKey, sessionToken string) error {
	a.credLock.Lock()
	defer a.credLock.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	a.logger.Trace("checking for new credentials")

	// When no static credentials are provided, don't force the shared "default"
	// profile. In SDK v2, setting the shared profile makes credential resolution
	// short-circuit to that profile and skip the environment and IMDS/ECS
	// providers; leaving it unset lets the default chain (env vars, shared config,
	// IMDS/ECS) resolve naturally, matching the v1 SDK behavior.
	var credOpts []awsutil.Option
	if accessKey == "" && secretKey == "" {
		credOpts = append(credOpts, awsutil.WithSharedCredentials(false))
	}

	currentConfig, err := awsutil.RetrieveCreds(ctx, accessKey, secretKey, sessionToken, a.logger, credOpts...)
	if err != nil {
		return err
	}

	currentVal, err := currentConfig.Credentials.Retrieve(ctx)
	if err != nil {
		return err
	}
	lastVal, err := a.lastCfg.Credentials.Retrieve(ctx)
	if err != nil {
		return err
	}

	// These will always have different pointers regardless of whether their
	// values are identical, hence the use of DeepEqual.
	if !lastVal.Expired() && reflect.DeepEqual(currentVal, lastVal) {
		a.logger.Trace("credentials are unchanged and still valid")
		return nil
	}

	// Only update the credentials, preserving region and other options set during initialization.
	a.lastCfg.Credentials = currentConfig.Credentials
	a.logger.Trace("new credentials detected, triggering Authenticate")

	// Use select to avoid panic if stopCh is closed during shutdown
	select {
	case a.credsFound <- struct{}{}:
	case <-a.stopCh:
		a.logger.Trace("shutdown detected, skipping credential notification")
	default:
	}
	return nil
}
