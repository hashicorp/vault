package api

import (
	"encoding/json"
	"io"
	"strconv"
	"time"

	"github.com/hashicorp/vault/helper/jsonutil"
)

// Secret is the structure returned for every secret within Vault.
type Secret struct {
	// The request ID that generated this response
	RequestID string `json:"request_id"`

	LeaseID       string `json:"lease_id"`
	LeaseDuration int    `json:"lease_duration"`
	Renewable     bool   `json:"renewable"`

	// Data is the actual contents of the secret. The format of the data
	// is arbitrary and up to the secret backend.
	Data map[string]interface{} `json:"data"`

	// Warnings contains any warnings related to the operation. These
	// are not issues that caused the command to fail, but that the
	// client should be aware of.
	Warnings []string `json:"warnings"`

	// Auth, if non-nil, means that there was authentication information
	// attached to this response.
	Auth *SecretAuth `json:"auth,omitempty"`

	// WrapInfo, if non-nil, means that the initial response was wrapped in the
	// cubbyhole of the given token (which has a TTL of the given number of
	// seconds)
	WrapInfo *SecretWrapInfo `json:"wrap_info,omitempty"`
}

// TokenID returns the standardized token ID (token) for the given secret.
func (s *Secret) TokenID() string {
	if s == nil {
		return ""
	}

	if s.Auth != nil && len(s.Auth.ClientToken) > 0 {
		return s.Auth.ClientToken
	}

	if s.Data == nil || s.Data["id"] == nil {
		return ""
	}

	id, ok := s.Data["id"].(string)
	if !ok {
		return ""
	}

	return id
}

// TokenAccessor returns the standardized token accessor for the given secret.
// If the secret is nil or does not contain an accessor, this returns the empty
// string.
func (s *Secret) TokenAccessor() string {
	if s == nil {
		return ""
	}

	if s.Auth != nil && len(s.Auth.Accessor) > 0 {
		return s.Auth.Accessor
	}

	if s.Data == nil || s.Data["accessor"] == nil {
		return ""
	}

	accessor, ok := s.Data["accessor"].(string)
	if !ok {
		return ""
	}

	return accessor
}

// TokenMeta returns the standardized token metadata for the given secret.
// If the secret is nil or does not contain an accessor, this returns the empty
// string. Metadata is usually modeled as an map[string]interface{}, but token
// metdata is always a map[string]string. This function handles the coercion.
func (s *Secret) TokenMeta() map[string]string {
	if s == nil {
		return nil
	}

	if s.Auth != nil && len(s.Auth.Metadata) > 0 {
		return s.Auth.Metadata
	}

	if s.Data == nil || s.Data["meta"] == nil {
		return nil
	}

	metaRaw, ok := s.Data["meta"].(map[string]interface{})
	if !ok {
		return nil
	}

	meta := make(map[string]string, len(metaRaw))
	for k, v := range metaRaw {
		m, ok := v.(string)
		if !ok {
			return nil
		}
		meta[k] = m
	}

	return meta
}

// TokenRemainingUses returns the standardized remaining uses for the given
// secret. If the secret is nil or does not contain the "num_uses", this returns
// 0..
func (s *Secret) TokenRemainingUses() int {
	if s == nil || s.Data == nil || s.Data["num_uses"] == nil {
		return 0
	}

	usesStr, ok := s.Data["num_uses"].(json.Number)
	if !ok {
		return 0
	}

	if string(usesStr) == "" {
		return 0
	}

	uses, err := strconv.ParseInt(string(usesStr), 10, 64)
	if err != nil {
		return 0
	}

	return int(uses)
}

// TokenPolicies returns the standardized list of policies for the given secret.
// If the secret is nil or does not contain any policies, this returns nil.
// Policies are usually returned as []interface{}, but this function ensures
// they are []string.
func (s *Secret) TokenPolicies() []string {
	if s == nil {
		return nil
	}

	if s.Auth != nil && len(s.Auth.Policies) > 0 {
		return s.Auth.Policies
	}

	if s.Data == nil || s.Data["policies"] == nil {
		return nil
	}

	list, ok := s.Data["policies"].([]interface{})
	if !ok {
		return nil
	}

	policies := make([]string, len(list))
	for i := range list {
		p, ok := list[i].(string)
		if !ok {
			return nil
		}
		policies[i] = p
	}

	return policies
}

// TokenIsRenewable returns the standardized token renewability for the given
// secret. If the secret is nil or does not contain the "renewable" key, this
// returns false.
func (s *Secret) TokenIsRenewable() bool {
	if s == nil {
		return false
	}

	if s.Auth != nil && s.Auth.Renewable {
		return s.Auth.Renewable
	}

	if s.Data == nil || s.Data["renewable"] == nil {
		return false
	}

	renewable, ok := s.Data["renewable"].(bool)
	if !ok {
		return false
	}

	return renewable
}

// TokenTTL returns the standardized remaining token TTL for the given secret.
// If the secret is nil or does not contain a TTL, this returns the 0.
func (s *Secret) TokenTTL() time.Duration {
	if s == nil {
		return 0
	}

	if s.Auth != nil && s.Auth.LeaseDuration > 0 {
		return time.Duration(s.Auth.LeaseDuration) * time.Second
	}

	if s.Data == nil || s.Data["ttl"] == nil {
		return 0
	}

	ttlStr, ok := s.Data["ttl"].(json.Number)
	if !ok {
		return 0
	}

	if string(ttlStr) == "" {
		return 0
	}

	ttl, err := time.ParseDuration(string(ttlStr) + "s")
	if err != nil {
		return 0
	}

	return ttl
}

// SecretWrapInfo contains wrapping information if we have it. If what is
// contained is an authentication token, the accessor for the token will be
// available in WrappedAccessor.
type SecretWrapInfo struct {
	Token           string    `json:"token"`
	TTL             int       `json:"ttl"`
	CreationTime    time.Time `json:"creation_time"`
	CreationPath    string    `json:"creation_path"`
	WrappedAccessor string    `json:"wrapped_accessor"`
}

// SecretAuth is the structure containing auth information if we have it.
type SecretAuth struct {
	ClientToken string            `json:"client_token"`
	Accessor    string            `json:"accessor"`
	Policies    []string          `json:"policies"`
	Metadata    map[string]string `json:"metadata"`

	LeaseDuration int  `json:"lease_duration"`
	Renewable     bool `json:"renewable"`
}

// ParseSecret is used to parse a secret value from JSON from an io.Reader.
func ParseSecret(r io.Reader) (*Secret, error) {
	// First decode the JSON into a map[string]interface{}
	var secret Secret
	if err := jsonutil.DecodeJSONFromReader(r, &secret); err != nil {
		return nil, err
	}

	return &secret, nil
}
