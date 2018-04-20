package roles

import (
	"testing"

	"github.com/hashicorp/vault/builtin/logical/ad/config"
	"github.com/hashicorp/vault/logical/framework"
)

var (
	mgr    = &handler{}
	schema = mgr.Path().Fields
)

func TestOnlyDefaultTTLs(t *testing.T) {

	passwordConf := &config.PasswordConf{
		TTL:    config.DefaultPasswordTTLs,
		MaxTTL: config.DefaultPasswordTTLs,
		Length: config.DefaultPasswordLength,
	}

	fieldData := &framework.FieldData{
		Raw: map[string]interface{}{
			"service_account_name": "kibana@example.com",
		},
		Schema: schema,
	}

	ttl, err := getValidatedTTL(passwordConf, fieldData)
	if err != nil {
		t.Fatal(err)
	}

	if ttl != config.DefaultPasswordTTLs {
		t.Fatal("ttl is not defaulting properly")
	}
}

func TestCustomOperatorTTLButDefaultRoleTTL(t *testing.T) {

	passwordConf := &config.PasswordConf{
		TTL:    10,
		MaxTTL: config.DefaultPasswordTTLs,
		Length: config.DefaultPasswordLength,
	}

	fieldData := &framework.FieldData{
		Raw: map[string]interface{}{
			"service_account_name": "kibana@example.com",
		},
		Schema: schema,
	}

	ttl, err := getValidatedTTL(passwordConf, fieldData)
	if err != nil {
		t.Fatal(err)
	}

	if ttl != 10 {
		t.Fatal("ttl is not defaulting properly")
	}
}

func TestTTLTooHigh(t *testing.T) {

	passwordConf := &config.PasswordConf{
		TTL:    10,
		MaxTTL: 10,
		Length: config.DefaultPasswordLength,
	}

	fieldData := &framework.FieldData{
		Raw: map[string]interface{}{
			"service_account_name": "kibana@example.com",
			"ttl": 100,
		},
		Schema: schema,
	}

	_, err := getValidatedTTL(passwordConf, fieldData)
	if err == nil {
		t.Fatal("should error when ttl is too high")
	}
}

func TestNegativeTTL(t *testing.T) {

	passwordConf := &config.PasswordConf{
		TTL:    10,
		MaxTTL: config.DefaultPasswordTTLs,
		Length: config.DefaultPasswordLength,
	}

	fieldData := &framework.FieldData{
		Raw: map[string]interface{}{
			"service_account_name": "kibana@example.com",
			"ttl": -100,
		},
		Schema: schema,
	}

	_, err := getValidatedTTL(passwordConf, fieldData)
	if err == nil {
		t.Fatal("should error then ttl is negative")
	}
}

func TestZeroTTL(t *testing.T) {

	passwordConf := &config.PasswordConf{
		TTL:    10,
		MaxTTL: config.DefaultPasswordTTLs,
		Length: config.DefaultPasswordLength,
	}

	fieldData := &framework.FieldData{
		Raw: map[string]interface{}{
			"service_account_name": "kibana@example.com",
			"ttl": 0,
		},
		Schema: schema,
	}

	_, err := getValidatedTTL(passwordConf, fieldData)
	if err == nil {
		t.Fatal("should error then ttl is zero")
	}
}
