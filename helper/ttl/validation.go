package ttl

import (
	"errors"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/logical"
)

type MountHandler struct {
	// ConfigTTL is the TTL being set at a mount's config level.
	// For example, if your plugin were the Azure secrets engine,
	// and you had a path for an overall config like <mount>/config,
	// this would be the TTL at that level.
	// Optional, can be unset.
	ConfigTTL time.Duration

	// ConfigMaxTTL is the MaxTTL being set at a mount's config level.
	// For example, if your plugin were the Azure secrets engine,
	// and you had a path for an overall config like <mount>/config,
	// this would be the MaxTTL at that level.
	// Optional, can be unset.
	ConfigMaxTTL time.Duration

	// RoleTTL is the TTL being set at a role level, which is a lower
	// and more specific level than the config level.
	// Optional, can be unset.
	RoleTTL time.Duration

	// RoleMaxTTL is the MaxTTL being set at a role level, which is a lower
	// and more specific level than the config level.
	// Optional, can be unset.
	RoleMaxTTL time.Duration
}

func (h *MountHandler) Validate(system logical.SystemView) error {

	merr := &multierror.Error{}

	// Verify the config-level TTL's alone.
	if h.ConfigTTL < 0 {
		merr = multierror.Append(merr, errors.New("config ttl < 0"))
	}
	if h.ConfigMaxTTL < 0 {
		merr = multierror.Append(merr, errors.New("config max_ttl < 0"))
	}
	if h.ConfigTTL > system.DefaultLeaseTTL() {
		merr = multierror.Append(merr, errors.New("config ttl > system defined TTL"))
	}
	if h.ConfigMaxTTL > system.MaxLeaseTTL() {
		merr = multierror.Append(merr, errors.New("config max_ttl > system defined max TTL"))
	}
	if h.ConfigMaxTTL != 0 && h.ConfigTTL > h.ConfigMaxTTL {
		merr = multierror.Append(merr, errors.New("config ttl > config max_ttl"))
	}

	// Verify the role-level TTL's alone.
	if h.RoleTTL < 0 {
		merr = multierror.Append(merr, errors.New("role ttl < 0"))
	}
	if h.RoleMaxTTL < 0 {
		merr = multierror.Append(merr, errors.New("role max_ttl < 0"))
	}
	if h.RoleTTL > system.DefaultLeaseTTL() {
		merr = multierror.Append(merr, errors.New("role ttl > system defined TTL"))
	}
	if h.RoleMaxTTL > system.MaxLeaseTTL() {
		merr = multierror.Append(merr, errors.New("role max_ttl > system defined max TTL"))
	}
	if h.RoleMaxTTL != 0 && h.RoleTTL > h.RoleMaxTTL {
		merr = multierror.Append(merr, errors.New("role ttl > role max_ttl"))
	}

	// Verify the config and role TTL's in relation to each other.
	if h.ConfigTTL != 0 && h.RoleTTL > h.ConfigTTL {
		merr = multierror.Append(merr, errors.New("role ttl > config ttl"))
	}
	if h.ConfigMaxTTL != 0 && h.RoleMaxTTL > h.ConfigMaxTTL {
		merr = multierror.Append(merr, errors.New("role max_ttl > config max_ttl"))
	}
	return merr.ErrorOrNil()
}
