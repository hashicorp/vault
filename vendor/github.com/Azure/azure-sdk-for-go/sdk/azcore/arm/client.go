//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package arm

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	armpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
)

// ClientOptions contains configuration settings for a client's pipeline.
type ClientOptions = armpolicy.ClientOptions

// Client is a HTTP client for use with ARM endpoints.  It consists of an endpoint, pipeline, and tracing provider.
type Client struct {
	ep string
	pl runtime.Pipeline
	tr tracing.Tracer
}

// NewClient creates a new Client instance with the provided values.
// This client is intended to be used with Azure Resource Manager endpoints.
//   - moduleName - the fully qualified name of the module where the client is defined; used by the telemetry policy and tracing provider.
//   - moduleVersion - the semantic version of the module; used by the telemetry policy and tracing provider.
//   - cred - the TokenCredential used to authenticate the request
//   - options - optional client configurations; pass nil to accept the default values
func NewClient(moduleName, moduleVersion string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	if !options.Telemetry.Disabled {
		if err := shared.ValidateModVer(moduleVersion); err != nil {
			return nil, err
		}
	}

	ep := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, cred, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}

	tr := options.TracingProvider.NewTracer(moduleName, moduleVersion)
	return &Client{ep: ep, pl: pl, tr: tr}, nil
}

// Endpoint returns the service's base URL for this client.
func (c *Client) Endpoint() string {
	return c.ep
}

// Pipeline returns the pipeline for this client.
func (c *Client) Pipeline() runtime.Pipeline {
	return c.pl
}

// Tracer returns the tracer for this client.
func (c *Client) Tracer() tracing.Tracer {
	return c.tr
}
