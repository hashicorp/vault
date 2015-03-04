package api

import (
	"net/http"
)

// Config is used to configure the creation of the client.
type Config struct {
	// Address is the address of the Vault server. This should be a complete
	// URL such as "http://vault.example.com". If you need a custom SSL
	// cert or want to enable insecure mode, you need to specify a custom
	// HttpClient.
	Address string

	// HttpClient is the HTTP client to use. http.DefaultClient will be
	// used if not specified.
	HttpClient *http.Client
}

// DefaultConfig returns a default configuration for the client. It is
// safe to modify the return value of this function.
func DefaultConfig() *Config {
	config := &Config{
		Address:    "https://127.0.0.1:8200",
		HttpClient: http.DefaultClient,
	}

	return config
}

// Client is the client to the Vault API. Create a client with
// NewClient.
type Client struct {
	config Config
}

// NewClient returns a new client for the given configuration.
func NewClient(c Config) (*Client, error) {
	return &Client{
		config: c,
	}, nil
}
