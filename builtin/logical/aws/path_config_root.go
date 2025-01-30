// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/automatedrotationutil"
	"github.com/hashicorp/vault/sdk/helper/pluginidentityutil"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/rotation"
)

// A single default template that supports both the different credential types (IAM/STS) that are capped at differing length limits (64 chars/32 chars respectively)
const (
	defaultUserNameTemplate = `{{ if (eq .Type "STS") }}{{ printf "vault-%s-%s"  (unix_time) (random 20) | truncate 32 }}{{ else }}{{ printf "vault-%s-%s-%s" (printf "%s-%s" (.DisplayName) (.PolicyName) | truncate 42) (unix_time) (random 20) | truncate 64 }}{{ end }}`
	rootRotationJobName     = "aws-root-creds"
)

func pathConfigRoot(b *backend) *framework.Path {
	p := &framework.Path{
		Pattern: "config/root",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAWS,
		},

		Fields: map[string]*framework.FieldSchema{
			"access_key": {
				Type:        framework.TypeString,
				Description: "Access key with permission to create new keys.",
			},

			"secret_key": {
				Type:        framework.TypeString,
				Description: "Secret key with permission to create new keys.",
			},

			"region": {
				Type:        framework.TypeString,
				Description: "Region for API calls.",
			},
			"iam_endpoint": {
				Type:        framework.TypeString,
				Description: "Endpoint to custom IAM server URL",
			},
			"sts_endpoint": {
				Type:        framework.TypeString,
				Description: "Endpoint to custom STS server URL",
			},
			"sts_region": {
				Type:        framework.TypeString,
				Description: "Specific region for STS API calls.",
			},
			"sts_fallback_endpoints": {
				Type:        framework.TypeCommaStringSlice,
				Description: "Fallback endpoints if sts_endpoint is unreachable",
			},
			"sts_fallback_regions": {
				Type:        framework.TypeCommaStringSlice,
				Description: "Fallback regions if sts_region is unreachable",
			},
			"max_retries": {
				Type:        framework.TypeInt,
				Default:     aws.UseServiceDefaultRetries,
				Description: "Maximum number of retries for recoverable exceptions of AWS APIs",
			},
			"username_template": {
				Type:        framework.TypeString,
				Description: "Template to generate custom IAM usernames",
			},
			"role_arn": {
				Type:        framework.TypeString,
				Description: "Role ARN to assume for plugin identity token federation",
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigRootRead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "root-iam-credentials-configuration",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigRootWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "root-iam-credentials",
				},
			},
		},

		HelpSynopsis:    pathConfigRootHelpSyn,
		HelpDescription: pathConfigRootHelpDesc,
	}
	pluginidentityutil.AddPluginIdentityTokenFields(p.Fields)
	automatedrotationutil.AddAutomatedRotationFields(p.Fields)

	return p
}

func (b *backend) pathConfigRootRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.clientMutex.RLock()
	defer b.clientMutex.RUnlock()

	config, exists, err := getConfigFromStorage(ctx, req)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	configData := map[string]interface{}{
		"access_key":             config.AccessKey,
		"region":                 config.Region,
		"iam_endpoint":           config.IAMEndpoint,
		"sts_endpoint":           config.STSEndpoint,
		"sts_region":             config.STSRegion,
		"sts_fallback_endpoints": config.STSFallbackEndpoints,
		"sts_fallback_regions":   config.STSFallbackRegions,
		"max_retries":            config.MaxRetries,
		"username_template":      config.UsernameTemplate,
		"role_arn":               config.RoleARN,
	}

	config.PopulatePluginIdentityTokenData(configData)
	config.PopulateAutomatedRotationData(configData)

	return &logical.Response{
		Data: configData,
	}, nil
}

func (b *backend) pathConfigRootWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	region := data.Get("region").(string)
	iamendpoint := data.Get("iam_endpoint").(string)
	stsendpoint := data.Get("sts_endpoint").(string)
	stsregion := data.Get("sts_region").(string)
	maxretries := data.Get("max_retries").(int)
	roleARN := data.Get("role_arn").(string)
	usernameTemplate := data.Get("username_template").(string)
	if usernameTemplate == "" {
		usernameTemplate = defaultUserNameTemplate
	}

	stsFallbackEndpoints := data.Get("sts_fallback_endpoints").([]string)
	stsFallbackRegions := data.Get("sts_fallback_regions").([]string)

	if len(stsFallbackEndpoints) != len(stsFallbackRegions) {
		return logical.ErrorResponse("fallback endpoints and fallback regions must be the same length"), nil
	}

	b.clientMutex.Lock()
	defer b.clientMutex.Unlock()

	// check for existing config
	previousCfg, previousCfgExists, err := getConfigFromStorage(ctx, req)
	if err != nil {
		return nil, err
	}

	rc := rootConfig{
		AccessKey:            data.Get("access_key").(string),
		SecretKey:            data.Get("secret_key").(string),
		IAMEndpoint:          iamendpoint,
		STSEndpoint:          stsendpoint,
		STSRegion:            stsregion,
		STSFallbackEndpoints: stsFallbackEndpoints,
		STSFallbackRegions:   stsFallbackRegions,
		Region:               region,
		MaxRetries:           maxretries,
		UsernameTemplate:     usernameTemplate,
		RoleARN:              roleARN,
	}
	if err := rc.ParsePluginIdentityTokenFields(data); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	if err := rc.ParseAutomatedRotationFields(data); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	if rc.IdentityTokenAudience != "" && rc.AccessKey != "" {
		return logical.ErrorResponse("only one of 'access_key' or 'identity_token_audience' can be set"), nil
	}

	if rc.IdentityTokenAudience != "" && rc.RoleARN == "" {
		return logical.ErrorResponse("missing required 'role_arn' when 'identity_token_audience' is set"), nil
	}

	if rc.IdentityTokenAudience != "" {
		_, err := b.System().GenerateIdentityToken(ctx, &pluginutil.IdentityTokenRequest{
			Audience: rc.IdentityTokenAudience,
		})
		if err != nil {
			if errors.Is(err, pluginidentityutil.ErrPluginWorkloadIdentityUnsupported) {
				return logical.ErrorResponse(err.Error()), nil
			}
			return nil, err
		}
	}

	// Save the initial config only if it does not already exist
	if !previousCfgExists {
		if err := putConfigToStorage(ctx, req, &rc); err != nil {
			return nil, err
		}
	}

	// Now that the root config is set up, register the rotation job if it required
	if rc.ShouldRegisterRotationJob() {
		cfgReq := &rotation.RotationJobConfigureRequest{
			Name:             rootRotationJobName,
			MountType:        req.MountType,
			ReqPath:          req.Path,
			RotationSchedule: rc.RotationSchedule,
			RotationWindow:   rc.RotationWindow,
			RotationPeriod:   rc.RotationPeriod,
		}

		_, err = b.System().RegisterRotationJob(ctx, cfgReq)
		if err != nil {
			return logical.ErrorResponse("error registering rotation job: %s", err), nil
		}
	}

	// Disable Automated Rotation and Deregister credentials if required
	if rc.DisableAutomatedRotation {
		// Ensure de-registering only occurs on updates and if
		// a credential has actually been registered (rotation_period or rotation_schedule is set)
		deregisterReq := &rotation.RotationJobDeregisterRequest{
			MountType: req.MountType,
			ReqPath:   req.Path,
		}
		if previousCfgExists && previousCfg.ShouldRegisterRotationJob() {
			err := b.System().DeregisterRotationJob(ctx, deregisterReq)
			if err != nil {
				return logical.ErrorResponse("error de-registering rotation job: %s", err), nil
			}
		}
	}

	// update config entry with rotation ID
	if err := putConfigToStorage(ctx, req, &rc); err != nil {
		return nil, err
	}

	// clear possible cached IAM / STS clients after successfully updating
	// config/root
	b.iamClient = nil
	b.stsClient = nil

	return nil, nil
}

func getConfigFromStorage(ctx context.Context, req *logical.Request) (*rootConfig, bool, error) {
	entry, err := req.Storage.Get(ctx, "config/root")
	if err != nil {
		return nil, false, err
	}
	if entry == nil {
		return nil, false, nil
	}

	var config rootConfig

	if err := entry.DecodeJSON(&config); err != nil {
		return nil, false, err
	}

	return &config, true, nil
}

func putConfigToStorage(ctx context.Context, req *logical.Request, rc *rootConfig) error {
	entry, err := logical.StorageEntryJSON("config/root", rc)
	if err != nil {
		return err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return err
	}

	return nil
}

type rootConfig struct {
	pluginidentityutil.PluginIdentityTokenParams
	automatedrotationutil.AutomatedRotationParams

	AccessKey            string   `json:"access_key"`
	SecretKey            string   `json:"secret_key"`
	IAMEndpoint          string   `json:"iam_endpoint"`
	STSEndpoint          string   `json:"sts_endpoint"`
	STSRegion            string   `json:"sts_region"`
	STSFallbackEndpoints []string `json:"sts_fallback_endpoints"`
	STSFallbackRegions   []string `json:"sts_fallback_regions"`
	Region               string   `json:"region"`
	MaxRetries           int      `json:"max_retries"`
	UsernameTemplate     string   `json:"username_template"`
	RoleARN              string   `json:"role_arn"`
}

const pathConfigRootHelpSyn = `
Configure the root credentials that are used to manage IAM.
`

const pathConfigRootHelpDesc = `
Before doing anything, the AWS backend needs credentials that are able
to manage IAM policies, users, access keys, etc. This endpoint is used
to configure those credentials. They don't necessarily need to be root
keys as long as they have permission to manage IAM.
`
