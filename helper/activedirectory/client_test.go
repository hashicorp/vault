package activedirectory

import (
	"testing"

	"github.com/go-ldap/ldap"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/ldapifc"
)

func TestSearch(t *testing.T) {

	config := emptyConfig()

	conn := &ldapifc.FakeLDAPConnection{
		SearchRequestToExpect: testSearchRequest(),
		SearchResultToReturn:  testSearchResult(),
	}

	ldapClient := &ldapifc.FakeLDAPClient{conn}

	client := NewClientWith(hclog.NewNullLogger(), config, ldapClient)

	filters := map[*Field][]string{
		FieldRegistry.Surname: {"Jones"},
	}

	entries, err := client.Search(filters)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 1 {
		t.Fatalf("only one entry was provided, but multiple were found: %s", entries)
	}
	entry := entries[0]

	result, _ := entry.GetJoined(FieldRegistry.Surname)
	if result != "Jones" {
		t.Fatalf("expected Surname of \"Jones\" but received \"%s\"", result)
	}

	result, _ = entry.GetJoined(FieldRegistry.BadPasswordTime)
	if result != "131653637947737037" {
		t.Fatalf("expected BadPasswordTime of \"131653637947737037\" but received \"%s\"", result)
	}

	result, _ = entry.GetJoined(FieldRegistry.PasswordLastSet)
	if result != "0" {
		t.Fatalf("expected PasswordLastSet of \"0\" but received \"%s\"", result)
	}

	result, _ = entry.GetJoined(FieldRegistry.PrimaryGroupID)
	if result != "513" {
		t.Fatalf("expected PrimaryGroupID of \"513\" but received \"%s\"", result)
	}

	result, _ = entry.GetJoined(FieldRegistry.UserPrincipalName)
	if result != "jim@example.com" {
		t.Fatalf("expected UserPrincipalName of \"jim@example.com\" but received \"%s\"", result)
	}

	result, _ = entry.GetJoined(FieldRegistry.ObjectClass)
	if result != "top,person,organizationalPerson,user" {
		t.Fatalf("expected ObjectClass of \"top,person,organizationalPerson,user\" but received \"%s\"", result)
	}
}

func TestUpdateEntry(t *testing.T) {

	config := emptyConfig()

	conn := &ldapifc.FakeLDAPConnection{
		SearchRequestToExpect: testSearchRequest(),
		SearchResultToReturn:  testSearchResult(),
	}

	conn.ModifyRequestToExpect = &ldap.ModifyRequest{
		DN: "CN=Jim H.. Jones,OU=Vault,OU=Engineering,DC=example,DC=com",
		ReplaceAttributes: []ldap.PartialAttribute{
			{
				Type: "cn",
				Vals: []string{"Blue", "Red"},
			},
		},
	}
	ldapClient := &ldapifc.FakeLDAPClient{conn}

	client := NewClientWith(hclog.NewNullLogger(), config, ldapClient)

	filters := map[*Field][]string{
		FieldRegistry.Surname: {"Jones"},
	}

	newValues := map[*Field][]string{
		FieldRegistry.CommonName: {"Blue", "Red"},
	}

	if err := client.UpdateEntry(filters, newValues); err != nil {
		t.Fatal(err)
	}
}

func TestUpdatePassword(t *testing.T) {

	testPass := "hell0$catz*"

	config := emptyConfig()
	config.Username = "cats"
	config.Password = "dogs"
	config.StartTLS = true

	conn := &ldapifc.FakeLDAPConnection{
		SearchRequestToExpect: testSearchRequest(),
		SearchResultToReturn:  testSearchResult(),
	}

	expectedPass, err := formatPassword(testPass)
	if err != nil {
		t.Fatal(err)
	}
	conn.ModifyRequestToExpect = &ldap.ModifyRequest{
		DN: "CN=Jim H.. Jones,OU=Vault,OU=Engineering,DC=example,DC=com",
		ReplaceAttributes: []ldap.PartialAttribute{
			{
				Type: "unicodePwd",
				Vals: []string{expectedPass},
			},
		},
	}
	ldapClient := &ldapifc.FakeLDAPClient{conn}

	client := NewClientWith(hclog.NewNullLogger(), config, ldapClient)

	filters := map[*Field][]string{
		FieldRegistry.Surname: {"Jones"},
	}

	if err := client.UpdatePassword(filters, testPass); err != nil {
		t.Fatal(err)
	}
}

func emptyConfig() *Configuration {
	return &Configuration{
		RootDomainName: "example,com",
		URLs:           []string{"ldap://127.0.0.1"},
	}
}

func testSearchRequest() *ldap.SearchRequest {
	return &ldap.SearchRequest{
		BaseDN: "dc=example,dc=com",
		Scope:  ldap.ScopeWholeSubtree,
		Filter: "(sn=Jones)",
	}
}

func testSearchResult() *ldap.SearchResult {
	return &ldap.SearchResult{
		Entries: []*ldap.Entry{
			{
				DN: "CN=Jim H.. Jones,OU=Vault,OU=Engineering,DC=example,DC=com",
				Attributes: []*ldap.EntryAttribute{
					{
						Name:   FieldRegistry.Surname.String(),
						Values: []string{"Jones"},
					},
					{
						Name:   FieldRegistry.BadPasswordTime.String(),
						Values: []string{"131653637947737037"},
					},
					{
						Name:   FieldRegistry.PasswordLastSet.String(),
						Values: []string{"0"},
					},
					{
						Name:   FieldRegistry.PrimaryGroupID.String(),
						Values: []string{"513"},
					},
					{
						Name:   FieldRegistry.UserPrincipalName.String(),
						Values: []string{"jim@example.com"},
					},
					{
						Name:   FieldRegistry.ObjectClass.String(),
						Values: []string{"top", "person", "organizationalPerson", "user"},
					},
				},
			},
		},
	}
}
