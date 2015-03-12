package api

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// Secret is the structure returned for every secret within Vault.
type Secret struct {
	VaultId          string `mapstructure:"vault_id"`
	Renewable        bool
	LeaseDuration    int                    `mapstructure:"lease_duration"`
	LeaseDurationMax int                    `mapstructure:"lease_duration_max"`
	Data             map[string]interface{} `mapstructure:"-"`
}

// ParseSecret is used to parse a secret value from JSON from an io.Reader.
func ParseSecret(r io.Reader) (*Secret, error) {
	// First decode the JSON into a map[string]interface{}
	var raw map[string]interface{}
	dec := json.NewDecoder(r)
	if err := dec.Decode(&raw); err != nil {
		return nil, err
	}

	// Use mapstructure to get as much as possible
	var result Secret
	var metadata mapstructure.Metadata
	mdec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: &metadata,
		Result:   &result,
	})
	if err != nil {
		return nil, err
	}
	if err := mdec.Decode(raw); err != nil {
		return nil, err
	}

	// Delete the keys we decoded from the raw value, then set that
	// raw value as the resulting data.
	for _, k := range metadata.Keys {
		delete(raw, strings.ToLower(k))
	}

	result.Data = raw
	return &result, nil
}
