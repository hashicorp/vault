package api

import (
	"fmt"
	"testing"
)

func TestSSH_CreateTLSClient(t *testing.T) {
	// load the default configuration
	config, err := LoadSSHAgentConfig("./test-fixtures/agent_config.hcl")
	if err != nil {
		panic(fmt.Sprintf("error loading agent's config file: %s", err))
	}

	client, err := config.NewClient()
	if err != nil {
		panic(fmt.Sprintf("error creating the client: %s", err))
	}

	// Provide a certificate and enforce setting of transport
	config.CACert = "./test-fixtures/vault.crt"

	client, err = config.NewClient()
	if err != nil {
		panic(fmt.Sprintf("error creating the client: %s", err))
	}
	if client.config.HttpClient.Transport == nil {
		panic(fmt.Sprintf("error creating client with TLS transport"))
	}
}
