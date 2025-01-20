// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package awsauth

import (
	"context"
	"errors"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathConfigClient() *framework.Path {
	return &framework.Path{
		Pattern: "config/client$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAWS,
		},

		Fields: map[string]*framework.FieldSchema{
			"access_key": {
				Type:        framework.TypeString,
				Default:     "",
				Description: "AWS Access Key ID for the account used to make AWS API requests.",
			},

			"secret_key": {
				Type:        framework.TypeString,
				Default:     "",
				Description: "AWS Secret Access Key for the account used to make AWS API requests.",
			},

			"endpoint": {
				Type:        framework.TypeString,
				Default:     "",
				Description: "URL to override the default generated endpoint for making AWS EC2 API calls.",
			},

			"iam_endpoint": {
				Type:        framework.TypeString,
				Default:     "",
				Description: "URL to override the default generated endpoint for making AWS IAM API calls.",
			},

			"sts_endpoint": {
				Type:        framework.TypeString,
				Default:     "",
				Description: "URL to override the default generated endpoint for making AWS STS API calls.",
			},

			"sts_region": {
				Type:        framework.TypeString,
				Default:     "",
				Description: "The region ID for the sts_endpoint, if set.",
			},

			"use_sts_region_from_client": {
				Type:        framework.TypeBool,
				Default:     false,
				Description: "Uses the STS region from client requests for making AWS STS API calls.",
			},

			"iam_server_id_header_value": {
				Type:        framework.TypeString,
				Default:     "",
				Description: "Value to require in the X-Vault-AWS-IAM-Server-ID request header",
			},

			"allowed_sts_header_values": {
				Type:        framework.TypeCommaStringSlice,
				Default:     nil,
				Description: "List of additional headers that are allowed to be in AWS STS request headers",
			},

			"max_retries": {
				Type:        framework.TypeInt,
				Default:     aws.UseServiceDefaultRetries,
				Description: "Maximum number of retries for recoverable exceptions of AWS APIs",
			},
		},

		ExistenceCheck: b.pathConfigClientExistenceCheck,

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathConfigClientCreateUpdate,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "client",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigClientCreateUpdate,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "client",
				},
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathConfigClientDelete,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "client-configuration",
				},
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigClientRead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "client-configuration",
				},
			},
		},

		HelpSynopsis:    pathConfigClientHelpSyn,
		HelpDescription: pathConfigClientHelpDesc,
	}
}

// Establishes dichotomy of request operation between CreateOperation and UpdateOperation.
// Returning 'true' forces an UpdateOperation, CreateOperation otherwise.
func (b *backend) pathConfigClientExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := b.lockedClientConfigEntry(ctx, req.Storage)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

// Fetch the client configuration required to access the AWS API, after acquiring an exclusive lock.
func (b *backend) lockedClientConfigEntry(ctx context.Context, s logical.Storage) (*clientConfig, error) {
	b.configMutex.RLock()
	defer b.configMutex.RUnlock()

	return b.nonLockedClientConfigEntry(ctx, s)
}

// Fetch the client configuration required to access the AWS API.
func (b *backend) nonLockedClientConfigEntry(ctx context.Context, s logical.Storage) (*clientConfig, error) {
	entry, err := s.Get(ctx, "config/client")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result clientConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (b *backend) pathConfigClientRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	clientConfig, err := b.lockedClientConfigEntry(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if clientConfig == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"access_key":                 clientConfig.AccessKey,
			"endpoint":                   clientConfig.Endpoint,
			"iam_endpoint":               clientConfig.IAMEndpoint,
			"sts_endpoint":               clientConfig.STSEndpoint,
			"sts_region":                 clientConfig.STSRegion,
			"use_sts_region_from_client": clientConfig.UseSTSRegionFromClient,
			"iam_server_id_header_value": clientConfig.IAMServerIdHeaderValue,
			"max_retries":                clientConfig.MaxRetries,
			"allowed_sts_header_values":  clientConfig.AllowedSTSHeaderValues,
		},
	}, nil
}

func (b *backend) pathConfigClientDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	if err := req.Storage.Delete(ctx, "config/client"); err != nil {
		return nil, err
	}

	// Remove all the cached EC2 client objects in the backend.
	b.flushCachedEC2Clients()

	// Remove all the cached EC2 client objects in the backend.
	b.flushCachedIAMClients()

	// unset the cached default AWS account ID
	b.defaultAWSAccountID = ""

	return nil, nil
}

// pathConfigClientCreateUpdate is used to register the 'aws_secret_key' and 'aws_access_key'
// that can be used to interact with AWS EC2 API.
func (b *backend) pathConfigClientCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.configMutex.Lock()
	defer b.configMutex.Unlock()

	configEntry, err := b.nonLockedClientConfigEntry(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if configEntry == nil {
		configEntry = &clientConfig{}
	}

	// changedCreds is whether we need to flush the cached AWS clients and store in the backend
	changedCreds := false
	// changedOtherConfig is whether other config has changed that requires storing in the backend
	// but does not require flushing the cached clients
	changedOtherConfig := false

	accessKeyStr, ok := data.GetOk("access_key")
	if ok {
		if configEntry.AccessKey != accessKeyStr.(string) {
			changedCreds = true
			configEntry.AccessKey = accessKeyStr.(string)
		}
	} else if req.Operation == logical.CreateOperation {
		// Use the default
		configEntry.AccessKey = data.Get("access_key").(string)
	}

	secretKeyStr, ok := data.GetOk("secret_key")
	if ok {
		if configEntry.SecretKey != secretKeyStr.(string) {
			changedCreds = true
			configEntry.SecretKey = secretKeyStr.(string)
		}
	} else if req.Operation == logical.CreateOperation {
		configEntry.SecretKey = data.Get("secret_key").(string)
	}

	endpointStr, ok := data.GetOk("endpoint")
	if ok {
		if configEntry.Endpoint != endpointStr.(string) {
			changedCreds = true
			configEntry.Endpoint = endpointStr.(string)
		}
	} else if req.Operation == logical.CreateOperation {
		configEntry.Endpoint = data.Get("endpoint").(string)
	}

	iamEndpointStr, ok := data.GetOk("iam_endpoint")
	if ok {
		if configEntry.IAMEndpoint != iamEndpointStr.(string) {
			changedCreds = true
			configEntry.IAMEndpoint = iamEndpointStr.(string)
		}
	} else if req.Operation == logical.CreateOperation {
		configEntry.IAMEndpoint = data.Get("iam_endpoint").(string)
	}

	stsEndpointStr, ok := data.GetOk("sts_endpoint")
	if ok {
		if configEntry.STSEndpoint != stsEndpointStr.(string) {
			// We don't directly cache STS clients as they are never directly used.
			// However, they are potentially indirectly used as credential providers
			// for the EC2 and IAM clients, and thus we would be indirectly caching
			// them there. So, if we change the STS endpoint, we should flush those
			// cached clients.
			changedCreds = true
			configEntry.STSEndpoint = stsEndpointStr.(string)
		}
	} else if req.Operation == logical.CreateOperation {
		configEntry.STSEndpoint = data.Get("sts_endpoint").(string)
	}

	stsRegionStr, ok := data.GetOk("sts_region")
	if ok {
		if configEntry.STSRegion != stsRegionStr.(string) {
			// Region is used when building STS clients. As such, all the comments
			// regarding the sts_endpoint changing apply here as well.
			changedCreds = true
			configEntry.STSRegion = stsRegionStr.(string)
		}
	}

	useSTSRegionFromClientRaw, ok := data.GetOk("use_sts_region_from_client")
	if ok {
		if configEntry.UseSTSRegionFromClient != useSTSRegionFromClientRaw.(bool) {
			changedCreds = true
			configEntry.UseSTSRegionFromClient = useSTSRegionFromClientRaw.(bool)
		}
	}

	headerValStr, ok := data.GetOk("iam_server_id_header_value")
	if ok {
		if configEntry.IAMServerIdHeaderValue != headerValStr.(string) {
			// NOT setting changedCreds here, since this isn't really cached
			configEntry.IAMServerIdHeaderValue = headerValStr.(string)
			changedOtherConfig = true
		}
	} else if req.Operation == logical.CreateOperation {
		configEntry.IAMServerIdHeaderValue = data.Get("iam_server_id_header_value").(string)
	}

	aHeadersValStr, ok := data.GetOk("allowed_sts_header_values")
	if ok {
		aHeadersValSl := aHeadersValStr.([]string)
		for i, v := range aHeadersValSl {
			aHeadersValSl[i] = textproto.CanonicalMIMEHeaderKey(v)
		}
		if !strutil.EquivalentSlices(configEntry.AllowedSTSHeaderValues, aHeadersValSl) {
			// NOT setting changedCreds here, since this isn't really cached
			configEntry.AllowedSTSHeaderValues = aHeadersValSl
			changedOtherConfig = true
		}
	} else if req.Operation == logical.CreateOperation {
		ah, ok := data.GetOk("allowed_sts_header_values")
		if ok {
			configEntry.AllowedSTSHeaderValues = ah.([]string)
		}
	}

	maxRetriesInt, ok := data.GetOk("max_retries")
	if ok {
		configEntry.MaxRetries = maxRetriesInt.(int)
		changedOtherConfig = true
	} else if req.Operation == logical.CreateOperation {
		configEntry.MaxRetries = data.Get("max_retries").(int)
	}

	// Since this endpoint supports both create operation and update operation,
	// the error checks for access_key and secret_key not being set are not present.
	// This allows calling this endpoint multiple times to provide the values.
	// Hence, the readers of this endpoint should do the validation on
	// the validation of keys before using them.
	entry, err := b.configClientToEntry(configEntry)
	if err != nil {
		return nil, err
	}

	if changedCreds || changedOtherConfig || req.Operation == logical.CreateOperation {
		if err := req.Storage.Put(ctx, entry); err != nil {
			return nil, err
		}
	}

	if changedCreds {
		b.flushCachedEC2Clients()
		b.flushCachedIAMClients()
		b.defaultAWSAccountID = ""
	}

	return nil, nil
}

// configClientToEntry allows the client config code to encapsulate its
// knowledge about where its config is stored. It also provides a way
// for other endpoints to update the config properly.
func (b *backend) configClientToEntry(conf *clientConfig) (*logical.StorageEntry, error) {
	entry, err := logical.StorageEntryJSON("config/client", conf)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

// Struct to hold 'aws_access_key' and 'aws_secret_key' that are required to
// interact with the AWS EC2 API.
type clientConfig struct {
	AccessKey              string   `json:"access_key"`
	SecretKey              string   `json:"secret_key"`
	Endpoint               string   `json:"endpoint"`
	IAMEndpoint            string   `json:"iam_endpoint"`
	STSEndpoint            string   `json:"sts_endpoint"`
	STSRegion              string   `json:"sts_region"`
	UseSTSRegionFromClient bool     `json:"use_sts_region_from_client"`
	IAMServerIdHeaderValue string   `json:"iam_server_id_header_value"`
	AllowedSTSHeaderValues []string `json:"allowed_sts_header_values"`
	MaxRetries             int      `json:"max_retries"`
}

func (c *clientConfig) validateAllowedSTSHeaderValues(headers http.Header) error {
	for k := range headers {
		h := textproto.CanonicalMIMEHeaderKey(k)
		if h == "X-Amz-Signedheaders" {
			h = amzSignedHeaders
		}
		if strings.HasPrefix(h, amzHeaderPrefix) &&
			!strutil.StrListContains(defaultAllowedSTSRequestHeaders, h) &&
			!strutil.StrListContains(c.AllowedSTSHeaderValues, h) {
			return errors.New("invalid request header: " + k)
		}
	}
	return nil
}

func (c *clientConfig) validateAllowedSTSQueryValues(params url.Values) error {
	for k := range params {
		h := textproto.CanonicalMIMEHeaderKey(k)
		if h == "X-Amz-Signedheaders" {
			h = amzSignedHeaders
		}
		if strings.HasPrefix(h, amzHeaderPrefix) &&
			!strutil.StrListContains(defaultAllowedSTSRequestHeaders, k) &&
			!strutil.StrListContains(c.AllowedSTSHeaderValues, k) {
			return errors.New("invalid request query param: " + k)
		}
	}
	return nil
}

const pathConfigClientHelpSyn = `
Configure AWS IAM credentials that are used to query instance and role details from the AWS API.
`

const pathConfigClientHelpDesc = `
The aws-ec2 auth method makes AWS API queries to retrieve information
regarding EC2 instances that perform login operations. The 'aws_secret_key' and
'aws_access_key' parameters configured here should map to an AWS IAM user that
has permission to make the following API queries:

* ec2:DescribeInstances
* iam:GetInstanceProfile (if IAM Role binding is used)
`
