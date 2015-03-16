package api

import (
	"encoding/json"
	"io"
)

// Secret is the structure returned for every secret within Vault.
type Secret struct {
	VaultId          string `json:"vault_id"`
	Renewable        bool
	LeaseDuration    int                    `json:"lease_duration"`
	LeaseDurationMax int                    `json:"lease_duration_max"`
	Data             map[string]interface{} `json:"data"`
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
