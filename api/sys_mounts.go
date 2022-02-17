package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

func (c *Sys) ListMounts() (map[string]*MountOutput, error) {
	r := c.c.NewRequest("GET", "/v1/sys/mounts")

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, errors.New("data from server response is empty")
	}

	mounts := map[string]*MountOutput{}
	err = mapstructure.Decode(secret.Data, &mounts)
	if err != nil {
		return nil, err
	}

	return mounts, nil
}

func (c *Sys) Mount(path string, mountInfo *MountInput) error {
	r := c.c.NewRequest("POST", fmt.Sprintf("/v1/sys/mounts/%s", path))
	if err := r.SetJSONBody(mountInfo); err != nil {
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Sys) Unmount(path string) error {
	r := c.c.NewRequest("DELETE", fmt.Sprintf("/v1/sys/mounts/%s", path))

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

// Remount kicks off a remount operation, polls the status endpoint using
// the migration ID till either success or failure state is observed
func (c *Sys) Remount(from, to string) error {
	remountResp, err := c.StartRemount(from, to)
	if err != nil {
		return err
	}

	for {
		remountStatusResp, err := c.RemountStatus(remountResp.MigrationID)
		if err != nil {
			return err
		}
		if remountStatusResp.MigrationInfo.MigrationStatus == "success" {
			return nil
		}
		if remountStatusResp.MigrationInfo.MigrationStatus == "failure" {
			return fmt.Errorf("Failure! Error encountered moving mount %s to %s, with migration ID %s", from, to, remountResp.MigrationID)
		}
		time.Sleep(1 * time.Second)
	}
}

// StartRemount kicks off a mount migration and returns a response with the migration ID
func (c *Sys) StartRemount(from, to string) (*MountMigrationOutput, error) {
	body := map[string]interface{}{
		"from": from,
		"to":   to,
	}

	r := c.c.NewRequest("POST", "/v1/sys/remount")
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, errors.New("data from server response is empty")
	}

	var result MountMigrationOutput
	err = mapstructure.Decode(secret.Data, &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}

// RemountStatus checks the status of a mount migration operation with the provided ID
func (c *Sys) RemountStatus(migrationID string) (*MountMigrationStatusOutput, error) {
	r := c.c.NewRequest("GET", fmt.Sprintf("/v1/sys/remount/status/%s", migrationID))

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, errors.New("data from server response is empty")
	}

	var result MountMigrationStatusOutput
	err = mapstructure.Decode(secret.Data, &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}

func (c *Sys) TuneMount(path string, config MountConfigInput) error {
	r := c.c.NewRequest("POST", fmt.Sprintf("/v1/sys/mounts/%s/tune", path))
	if err := r.SetJSONBody(config); err != nil {
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) MountConfig(path string) (*MountConfigOutput, error) {
	r := c.c.NewRequest("GET", fmt.Sprintf("/v1/sys/mounts/%s/tune", path))

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}
	if secret == nil || secret.Data == nil {
		return nil, errors.New("data from server response is empty")
	}

	var result MountConfigOutput
	err = mapstructure.Decode(secret.Data, &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}

type MountInput struct {
	Type                  string            `json:"type"`
	Description           string            `json:"description"`
	Config                MountConfigInput  `json:"config"`
	Local                 bool              `json:"local"`
	SealWrap              bool              `json:"seal_wrap" mapstructure:"seal_wrap"`
	ExternalEntropyAccess bool              `json:"external_entropy_access" mapstructure:"external_entropy_access"`
	Options               map[string]string `json:"options"`

	// Deprecated: Newer server responses should be returning this information in the
	// Type field (json: "type") instead.
	PluginName string `json:"plugin_name,omitempty"`
}

type MountConfigInput struct {
	Options                   map[string]string `json:"options" mapstructure:"options"`
	DefaultLeaseTTL           string            `json:"default_lease_ttl" mapstructure:"default_lease_ttl"`
	Description               *string           `json:"description,omitempty" mapstructure:"description"`
	MaxLeaseTTL               string            `json:"max_lease_ttl" mapstructure:"max_lease_ttl"`
	ForceNoCache              bool              `json:"force_no_cache" mapstructure:"force_no_cache"`
	AuditNonHMACRequestKeys   []string          `json:"audit_non_hmac_request_keys,omitempty" mapstructure:"audit_non_hmac_request_keys"`
	AuditNonHMACResponseKeys  []string          `json:"audit_non_hmac_response_keys,omitempty" mapstructure:"audit_non_hmac_response_keys"`
	ListingVisibility         string            `json:"listing_visibility,omitempty" mapstructure:"listing_visibility"`
	PassthroughRequestHeaders []string          `json:"passthrough_request_headers,omitempty" mapstructure:"passthrough_request_headers"`
	AllowedResponseHeaders    []string          `json:"allowed_response_headers,omitempty" mapstructure:"allowed_response_headers"`
	TokenType                 string            `json:"token_type,omitempty" mapstructure:"token_type"`
	AllowedManagedKeys        []string          `json:"allowed_managed_keys,omitempty" mapstructure:"allowed_managed_keys"`

	// Deprecated: This field will always be blank for newer server responses.
	PluginName string `json:"plugin_name,omitempty" mapstructure:"plugin_name"`
}

type MountOutput struct {
	UUID                  string            `json:"uuid"`
	Type                  string            `json:"type"`
	Description           string            `json:"description"`
	Accessor              string            `json:"accessor"`
	Config                MountConfigOutput `json:"config"`
	Options               map[string]string `json:"options"`
	Local                 bool              `json:"local"`
	SealWrap              bool              `json:"seal_wrap" mapstructure:"seal_wrap"`
	ExternalEntropyAccess bool              `json:"external_entropy_access" mapstructure:"external_entropy_access"`
}

type MountConfigOutput struct {
	DefaultLeaseTTL           int      `json:"default_lease_ttl" mapstructure:"default_lease_ttl"`
	MaxLeaseTTL               int      `json:"max_lease_ttl" mapstructure:"max_lease_ttl"`
	ForceNoCache              bool     `json:"force_no_cache" mapstructure:"force_no_cache"`
	AuditNonHMACRequestKeys   []string `json:"audit_non_hmac_request_keys,omitempty" mapstructure:"audit_non_hmac_request_keys"`
	AuditNonHMACResponseKeys  []string `json:"audit_non_hmac_response_keys,omitempty" mapstructure:"audit_non_hmac_response_keys"`
	ListingVisibility         string   `json:"listing_visibility,omitempty" mapstructure:"listing_visibility"`
	PassthroughRequestHeaders []string `json:"passthrough_request_headers,omitempty" mapstructure:"passthrough_request_headers"`
	AllowedResponseHeaders    []string `json:"allowed_response_headers,omitempty" mapstructure:"allowed_response_headers"`
	TokenType                 string   `json:"token_type,omitempty" mapstructure:"token_type"`
	AllowedManagedKeys        []string `json:"allowed_managed_keys,omitempty" mapstructure:"allowed_managed_keys"`

	// Deprecated: This field will always be blank for newer server responses.
	PluginName string `json:"plugin_name,omitempty" mapstructure:"plugin_name"`
}

type MountMigrationOutput struct {
	MigrationID string `mapstructure:"migration_id"`
}

type MountMigrationStatusOutput struct {
	MigrationID   string                    `mapstructure:"migration_id"`
	MigrationInfo *MountMigrationStatusInfo `mapstructure:"migration_info"`
}

type MountMigrationStatusInfo struct {
	SourceMount     string `mapstructure:"source_mount"`
	TargetMount     string `mapstructure:"target_mount"`
	MigrationStatus string `mapstructure:"status"`
}
