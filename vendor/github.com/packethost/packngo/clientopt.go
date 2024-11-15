package packngo

import (
	"net/http"
	"net/url"
)

// ClientOpt is an option usable as an argument to NewClient constructor.
type ClientOpt func(*Client) error

// WithAuth configures Client with a specific consumerToken and apiKey for subsequent HTTP requests.
func WithAuth(consumerToken string, apiKey string) ClientOpt {
	return func(c *Client) error {
		c.ConsumerToken = consumerToken
		c.APIKey = apiKey
		c.apiKeySet = true

		return nil
	}
}

// WithHTTPClient configures Client to use a specific httpClient for subsequent HTTP requests.
func WithHTTPClient(httpClient *http.Client) ClientOpt {
	return func(c *Client) error {
		c.client = httpClient

		return nil
	}
}

// WithBaseURL configures Client to use a nonstandard API URL, e.g. for mocking the remote API.
func WithBaseURL(apiBaseURL string) ClientOpt {
	return func(c *Client) error {
		u, err := url.Parse(apiBaseURL)
		if err != nil {
			return err
		}

		c.BaseURL = u

		return nil
	}
}

// WithHeader configures Client to use the given HTTP header set.
// The headers X-Auth-Token, X-Consumer-Token, User-Agent will be ignored even if provided in the set.
func WithHeader(header http.Header) ClientOpt {
	return func(c *Client) error {
		for k, v := range header {
			c.header[k] = v
		}

		return nil
	}
}
