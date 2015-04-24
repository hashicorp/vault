package api

import (
	"encoding/json"
	"io"
)

// Secret is the structure returned for every secret within Vault.
type Secret struct {
	LeaseID       string `json:"lease_id"`
	LeaseDuration int    `json:"lease_duration"`
	Renewable     bool   `json:"renewable"`

	// Data is the actual contents of the secret. The format of the data
	// is arbitrary and up to the secret backend.
	Data map[string]interface{} `json:"data"`

	// Auth, if non-nil, means that there was authentication information
	// attached to this response.
	Auth *SecretAuth `json:"auth,omitempty"`
}

// Auth is the structure containing auth information if we have it.
type SecretAuth struct {
	ClientToken string            `json:"client_token"`
	Policies    []string          `json:"policies"`
	Metadata    map[string]string `json:"metadata"`

	LeaseDuration int  `json:"lease_duration"`
	Renewable     bool `json:"renewable"`
}

// ParseSecret is used to parse a secret value from JSON from an io.Reader.
func ParseSecret(r io.Reader) (*Secret, error) {
	// First decode the JSON into a map[string]interface{}
	var secret Secret
	dec := json.NewDecoder(r)
	if err := dec.Decode(&secret); err != nil {
		return nil, err
	}

	return &secret, nil
}
