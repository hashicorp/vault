// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package hcp

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"sync"

	slogctx "github.com/veqryn/slog-context"
)

// Client is the HCP client.
type Client struct {
	Environment Environment
	HTTPClient  *http.Client
	// Basic auth
	Username string
	Password string
	once     *sync.Once
}

// ClientOpt is an option to NewClient.
type ClientOpt func(*Client)

// Requester is an interface that defines a request that can be configured
// for an environment.
type Requester interface {
	Request(Environment) (*http.Request, error)
}

// Environment is an HCP portal environment
type Environment string

const (
	EnvironmentUnknown Environment = ""
	EnvironmentDev     Environment = "dev"
	EnvironmentInt     Environment = "int"
	EnvironmentProd    Environment = "prod"
)

// Addr is the URL for each environment
func (g Environment) Addr() string {
	switch g {
	case EnvironmentDev:
		return "https://api.hcp.dev"
	case EnvironmentInt:
		return "https://api.hcp.to"
	case EnvironmentProd:
		return "https://api.hashicorp.cloud"
	default:
		return ""
	}
}

// NewClient takes none-or-more options and returns a new Client.
func NewClient(opts ...ClientOpt) *Client {
	c := &Client{
		HTTPClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithEnvironment sets the client environment.
func WithEnvironment(env Environment) ClientOpt {
	return func(c *Client) {
		c.Environment = env
	}
}

// WithHTTPClient sets the client HTTP Client.
func WithHTTPClient(httpClient *http.Client) ClientOpt {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

// WithUsername sets the basic auth username for internal APIs.
func WithUsername(username string) ClientOpt {
	return func(c *Client) {
		c.Username = username
	}
}

// WithUsername sets the base auth password for internal APIs>
func WithPassword(password string) ClientOpt {
	return func(c *Client) {
		c.Password = password
	}
}

// WithLoadTokenFromEnv sets the basic auth username and token from known env
// vars.
func WithLoadAuthFromEnv() ClientOpt {
	return func(client *Client) {
		if username, ok := os.LookupEnv("HCP_USERNAME"); ok {
			client.Username = username
		}
		if password, ok := os.LookupEnv("HCP_PASSWORD"); ok {
			client.Password = password
		}
	}
}

// Do takes in a Requester and performs the request. It returns the raw http
// Response.
func (c *Client) Do(ctx context.Context, req Requester) (*http.Response, error) {
	logArgs := []any{
		slog.String("env", string(c.Environment)),
		slog.String("api-addr", string(c.Environment.Addr())),
	}
	httpReq, err := req.Request(c.Environment)
	if err != nil {
		slog.Default().ErrorContext(slogctx.Append(ctx,
			append(logArgs, slog.String("error", err.Error()))),
			"performing request",
		)
		return nil, err
	}

	logArgs = append(logArgs,
		slog.String("method", httpReq.Method),
		slog.String("url", httpReq.URL.String()),
	)

	ctx = slogctx.Append(ctx, logArgs...)
	slog.Default().DebugContext(ctx, "performing request")
	httpReq.SetBasicAuth(c.Username, c.Password)
	httpReq.WithContext(ctx)

	return c.HTTPClient.Do(httpReq)
}
