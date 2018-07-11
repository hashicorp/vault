package ttl

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
)

var system = &logical.StaticSystemView{
	DefaultLeaseTTLVal: time.Hour,
	MaxLeaseTTLVal:     time.Hour,
}

func TestNegativeTTLsThrowErrors(t *testing.T) {
	h := &MountHandler{
		ConfigTTL: -1,
	}
	if err := h.Validate(system); err == nil {
		t.Fatal("should have received error due to negative TTL")
	}

	h = &MountHandler{
		ConfigMaxTTL: -1,
	}
	if err := h.Validate(system); err == nil {
		t.Fatal("should have received error due to negative TTL")
	}

	h = &MountHandler{
		RoleTTL: -1,
	}
	if err := h.Validate(system); err == nil {
		t.Fatal("should have received error due to negative TTL")
	}

	h = &MountHandler{
		RoleMaxTTL: -1,
	}
	if err := h.Validate(system); err == nil {
		t.Fatal("should have received error due to negative TTL")
	}
}

func TestZeroTTLsAreOK(t *testing.T) {
	h := &MountHandler{}
	if err := h.Validate(system); err != nil {
		t.Fatal("having zeros for everything is OK")
	}
}

func TestConfigTTLsButNoRoleTTLs(t *testing.T) {
	h := &MountHandler{
		ConfigTTL:    1,
		ConfigMaxTTL: 1,
	}
	if err := h.Validate(system); err != nil {
		t.Fatal("having zeros for the role is OK")
	}
}

func TestRoleTTLsButNoConfigTTLs(t *testing.T) {
	h := &MountHandler{
		RoleTTL:    1,
		RoleMaxTTL: 1,
	}
	if err := h.Validate(system); err != nil {
		t.Fatal("having zeros for the config is OK")
	}
}

func TestRoleTTLsHigherThanConfigTTLs(t *testing.T) {
	h := &MountHandler{
		ConfigTTL:    1,
		ConfigMaxTTL: 1,
		RoleTTL:      2,
		RoleMaxTTL:   2,
	}
	if err := h.Validate(system); err == nil {
		t.Fatal("we should error when role TTLs are higher than config TTLs")
	}
}

func TestRoleTTLsHigherThanSystemTTLs(t *testing.T) {
	h := &MountHandler{
		RoleTTL:    2 * time.Hour,
		RoleMaxTTL: 2 * time.Hour,
	}
	if err := h.Validate(system); err == nil {
		t.Fatal("we should error when role TTLs are higher than system TTLs")
	}
}

func TestConfigTTLsHigherThanSystemTTLs(t *testing.T) {
	h := &MountHandler{
		ConfigTTL:    2 * time.Hour,
		ConfigMaxTTL: 2 * time.Hour,
	}
	if err := h.Validate(system); err == nil {
		t.Fatal("we should error when config TTLs are higher than system TTLs")
	}
}
