// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"errors"
	"fmt"

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

		ExistenceCheck: b.pathConfigRootExistenceCheck,

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
			logical.CreateOperation: &framework.PathOperation{
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

// Establishes dichotomy of request operation between CreateOperation and UpdateOperation.
// Returning 'true' forces an UpdateOperation, CreateOperation otherwise.
func (b *backend) pathConfigRootExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := getConfigFromStorage(ctx, req)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

func (b *backend) pathConfigRootRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.clientMutex.RLock()
	defer b.clientMutex.RUnlock()

	config, err := getConfigFromStorage(ctx, req)
	if err != nil {
		return nil, err
	}
	if config == nil {
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
	b.clientMutex.Lock()
	defer b.clientMutex.Unlock()

	// check for existing config
	rc, err := getConfigFromStorage(ctx, req)
	if err != nil {
		return nil, err
	}

	if rc == nil {
		// Baseline
		rc = &rootConfig{}
	}

	if accessKey, ok := data.GetOk("access_key"); ok {
		rc.AccessKey = accessKey.(string)
	} else if req.Operation == logical.CreateOperation {
		rc.AccessKey = data.Get("access_key").(string)
	}

	if secretKey, ok := data.GetOk("secret_key"); ok {
		rc.SecretKey = secretKey.(string)
	} else if req.Operation == logical.CreateOperation {
		rc.SecretKey = data.Get("secret_key").(string)
	}

	if region, ok := data.GetOk("region"); ok {
		rc.Region = region.(string)
	} else if req.Operation == logical.CreateOperation {
		rc.Region = data.Get("region").(string)
	}

	if iamEndpoint, ok := data.GetOk("iam_endpoint"); ok {
		rc.IAMEndpoint = iamEndpoint.(string)
	} else if req.Operation == logical.CreateOperation {
		rc.IAMEndpoint = data.Get("iam_endpoint").(string)
	}

	if stsEndpoint, ok := data.GetOk("sts_endpoint"); ok {
		rc.STSEndpoint = stsEndpoint.(string)
	} else if req.Operation == logical.CreateOperation {
		rc.STSEndpoint = data.Get("sts_endpoint").(string)
	}

	if stsRegion, ok := data.GetOk("sts_region"); ok {
		rc.STSRegion = stsRegion.(string)
	} else if req.Operation == logical.CreateOperation {
		rc.STSRegion = data.Get("sts_region").(string)
	}

	if maxRetries, ok := data.GetOk("max_retries"); ok {
		rc.MaxRetries = maxRetries.(int)
	} else if req.Operation == logical.CreateOperation {
		rc.MaxRetries = data.Get("max_retries").(int)
	}

	if roleARN, ok := data.GetOk("role_arn"); ok {
		rc.RoleARN = roleARN.(string)
	} else if req.Operation == logical.CreateOperation {
		rc.RoleARN = data.Get("role_arn").(string)
	}

	if stsFallbackEndpoints, ok := data.GetOk("sts_fallback_endpoints"); ok {
		rc.STSFallbackEndpoints = stsFallbackEndpoints.([]string)
	} else if req.Operation == logical.CreateOperation {
		rc.STSFallbackEndpoints = data.Get("sts_fallback_endpoints").([]string)
	}

	if stsFallbackRegions, ok := data.GetOk("sts_fallback_regions"); ok {
		rc.STSFallbackRegions = stsFallbackRegions.([]string)
	} else if req.Operation == logical.CreateOperation {
		rc.STSFallbackRegions = data.Get("sts_fallback_regions").([]string)
	}

	usernameTemplate := data.Get("username_template").(string)
	if usernameTemplate == "" {
		usernameTemplate = defaultUserNameTemplate
	}
	rc.UsernameTemplate = usernameTemplate

	if len(rc.STSFallbackEndpoints) != len(rc.STSFallbackRegions) {
		return logical.ErrorResponse("fallback endpoints and fallback regions must be the same length"), nil
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

	var performedRotationManagerOpern string
	if rc.ShouldDeregisterRotationJob() {
		performedRotationManagerOpern = rotation.PerformedDeregistration
		// Disable Automated Rotation and Deregister credentials if required
		deregisterReq := &rotation.RotationJobDeregisterRequest{
			MountPoint: req.MountPoint,
			ReqPath:    req.Path,
		}

		b.Logger().Debug("Deregistering rotation job", "mount", req.MountPoint+req.Path)
		if err := b.System().DeregisterRotationJob(ctx, deregisterReq); err != nil {
			return logical.ErrorResponse("error deregistering rotation job: %s", err), nil
		}
	} else if rc.ShouldRegisterRotationJob() {
		performedRotationManagerOpern = rotation.PerformedRegistration
		// Register the rotation job if it's required.
		cfgReq := &rotation.RotationJobConfigureRequest{
			MountPoint:       req.MountPoint,
			ReqPath:          req.Path,
			RotationSchedule: rc.RotationSchedule,
			RotationWindow:   rc.RotationWindow,
			RotationPeriod:   rc.RotationPeriod,
		}

		b.Logger().Debug("Registering rotation job", "mount", req.MountPoint+req.Path)
		if _, err = b.System().RegisterRotationJob(ctx, cfgReq); err != nil {
			return logical.ErrorResponse("error registering rotation job: %s", err), nil
		}
	}

	// Save the config
	if err := putConfigToStorage(ctx, req, rc); err != nil {
		wrappedError := err
		if performedRotationManagerOpern != "" {
			b.Logger().Error("write to storage failed but the rotation manager still succeeded.",
				"operation", performedRotationManagerOpern, "mount", req.MountPoint, "path", req.Path)
			wrappedError = fmt.Errorf("write to storage failed but the rotation manager still succeeded; "+
				"operation=%s, mount=%s, path=%s, storageError=%s", performedRotationManagerOpern, req.MountPoint, req.Path, err)
		}
		return nil, wrappedError
	}

	// clear possible cached IAM / STS clients after successfully updating
	// config/root
	b.iamClient = nil
	b.stsClient = nil

	return nil, nil
}

func getConfigFromStorage(ctx context.Context, req *logical.Request) (*rootConfig, error) {
	entry, err := req.Storage.Get(ctx, "config/root")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var config rootConfig

	if err := entry.DecodeJSON(&config); err != nil {
		return nil, err
	}

	return &config, nil
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
