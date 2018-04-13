package roles

import (
	"testing"

	"github.com/go-ldap/ldap"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/builtin/logical/ad/config"
	"github.com/hashicorp/vault/helper/activedirectory"
	"github.com/hashicorp/vault/helper/ldapifc"
	"github.com/hashicorp/vault/logical/framework"
)

var (
	manager  = &Manager{}
	schema   = manager.Path().Fields
	adClient = activedirectory.NewClientWith(hclog.NewNullLogger(), emptyConfig(), validLDAPClient())
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

	role, err := newRole(adClient, passwordConf, "kibana", fieldData)
	if err != nil {
		t.Fatal(err)
	}

	if role.TTL != config.DefaultPasswordTTLs {
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

	role, err := newRole(adClient, passwordConf, "kibana", fieldData)
	if err != nil {
		t.Fatal(err)
	}

	if role.TTL != 10 {
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

	_, err := newRole(adClient, passwordConf, "kibana", fieldData)
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

	_, err := newRole(adClient, passwordConf, "kibana", fieldData)
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

	_, err := newRole(adClient, passwordConf, "kibana", fieldData)
	if err == nil {
		t.Fatal("should error then ttl is zero")
	}
}

func validLDAPClient() ldapifc.Client {
	return &ldapifc.FakeLDAPClient{
		ConnToReturn: &ldapifc.FakeLDAPConnection{
			SearchRequestToExpect: &ldap.SearchRequest{
				BaseDN: "dc=example,dc=com",
				Filter: "(userPrincipalName=kibana@example.com)",
				Scope:  2,
			},
			SearchResultToReturn: &ldap.SearchResult{
				Entries: []*ldap.Entry{
					{
						DN: "CN=Jim H.. Jones,OU=Vault,OU=Engineering,DC=example,DC=com",
						Attributes: []*ldap.EntryAttribute{
							{
								Name:   activedirectory.FieldRegistry.LastLogon.String(),
								Values: []string{"131680504285591921"},
							},
						},
					},
				},
			},
		},
	}
}

func emptyConfig() *activedirectory.Configuration {
	return &activedirectory.Configuration{
		RootDomainName: "example,com",
		URLs:           []string{"ldap://127.0.0.1"},
	}
}
