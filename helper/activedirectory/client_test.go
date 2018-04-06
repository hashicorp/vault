package activedirectory

import (
	"crypto/tls"
	"fmt"
	"testing"

	"github.com/go-ldap/ldap"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/ldapifc"
)

func TestSearch(t *testing.T) {

	config := emptyConfig()

	conn := &fakeLDAPConnection{
		searchRequestToExpect: testSearchRequest(),
		searchResultToReturn:  testSearchResult(),
	}

	ldapClient := &fakeLDAPClient{conn}

	client := NewClientWith(hclog.NewNullLogger(), config, ldapClient)

	filters := map[*Field][]string{
		FieldRegistry.Surname: {"Jones"},
	}

	entries, err := client.Search(filters)
	if err != nil {
		t.Error(err)
	}

	if len(entries) != 1 {
		t.FailNow()
	}
	entry := entries[0]

	result, _ := entry.GetJoined(FieldRegistry.Surname)
	if result != "Jones" {
		t.FailNow()
	}

	result, _ = entry.GetJoined(FieldRegistry.BadPasswordTime)
	if result != "131653637947737037" {
		t.FailNow()
	}

	result, _ = entry.GetJoined(FieldRegistry.PasswordLastSet)
	if result != "0" {
		t.FailNow()
	}

	result, _ = entry.GetJoined(FieldRegistry.PrimaryGroupID)
	if result != "513" {
		t.Error(err)
		t.FailNow()
	}

	result, _ = entry.GetJoined(FieldRegistry.UserPrincipalName)
	if result != "jim@example.com" {
		t.FailNow()
	}

	result, _ = entry.GetJoined(FieldRegistry.ObjectClass)
	if result != "top,person,organizationalPerson,user" {
		t.FailNow()
	}
}

func TestUpdateEntry(t *testing.T) {

	config := emptyConfig()

	conn := &fakeLDAPConnection{
		searchRequestToExpect: testSearchRequest(),
		searchResultToReturn:  testSearchResult(),
	}

	conn.modifyRequestToExpect = &ldap.ModifyRequest{
		DN: "CN=Jim H.. Jones,OU=Vault,OU=Engineering,DC=example,DC=com",
		ReplaceAttributes: []ldap.PartialAttribute{
			{
				Type: "cn",
				Vals: []string{"Blue", "Red"},
			},
		},
	}
	ldapClient := &fakeLDAPClient{conn}

	client := NewClientWith(hclog.NewNullLogger(), config, ldapClient)

	filters := map[*Field][]string{
		FieldRegistry.Surname: {"Jones"},
	}

	newValues := map[*Field][]string{
		FieldRegistry.CommonName: {"Blue", "Red"},
	}

	if err := client.UpdateEntry(filters, newValues); err != nil {
		t.Error(err)
	}
}

func TestUpdatePassword(t *testing.T) {

	testPass := "hell0$catz*"

	config := emptyConfig()
	config.Username = "cats"
	config.Password = "dogs"
	config.StartTLS = true

	conn := &fakeLDAPConnection{
		searchRequestToExpect: testSearchRequest(),
		searchResultToReturn:  testSearchResult(),
	}

	expectedPass, err := formatPassword(testPass)
	if err != nil {
		t.Error(err)
	}
	conn.modifyRequestToExpect = &ldap.ModifyRequest{
		DN: "CN=Jim H.. Jones,OU=Vault,OU=Engineering,DC=example,DC=com",
		ReplaceAttributes: []ldap.PartialAttribute{
			{
				Type: "unicodePwd",
				Vals: []string{expectedPass},
			},
		},
	}
	ldapClient := &fakeLDAPClient{conn}

	client := NewClientWith(hclog.NewNullLogger(), config, ldapClient)

	filters := map[*Field][]string{
		FieldRegistry.Surname: {"Jones"},
	}

	if err := client.UpdatePassword(filters, testPass); err != nil {
		t.Error(err)
	}
}

type fakeLDAPClient struct {
	connToReturn ldapifc.Connection
}

func (f *fakeLDAPClient) Dial(network, addr string) (ldapifc.Connection, error) {
	return f.connToReturn, nil
}

func (f *fakeLDAPClient) DialTLS(network, addr string, config *tls.Config) (ldapifc.Connection, error) {
	return f.connToReturn, nil
}

type fakeLDAPConnection struct {
	modifyRequestToExpect *ldap.ModifyRequest

	searchRequestToExpect *ldap.SearchRequest
	searchResultToReturn  *ldap.SearchResult
}

func (f *fakeLDAPConnection) Bind(username, password string) error {
	return nil
}

func (f *fakeLDAPConnection) Close() {}

func (f *fakeLDAPConnection) Modify(modifyRequest *ldap.ModifyRequest) error {
	if f.modifyRequestToExpect.DN != modifyRequest.DN {
		return fmt.Errorf("expected modifyRequest of %s, but received %s", f.modifyRequestToExpect, modifyRequest)
	}
	if len(f.modifyRequestToExpect.ReplaceAttributes) != len(modifyRequest.ReplaceAttributes) {
		return fmt.Errorf("expected modifyRequest of %s, but received %s", f.modifyRequestToExpect, modifyRequest)
	}
	if f.modifyRequestToExpect.ReplaceAttributes[0].Type != modifyRequest.ReplaceAttributes[0].Type {
		return fmt.Errorf("expected modifyRequest of %s, but received %s", f.modifyRequestToExpect, modifyRequest)
	}
	if len(f.modifyRequestToExpect.ReplaceAttributes[0].Vals) != len(modifyRequest.ReplaceAttributes[0].Vals) {
		return fmt.Errorf("expected modifyRequest of %s, but received %s", f.modifyRequestToExpect, modifyRequest)
	}
	if f.modifyRequestToExpect.ReplaceAttributes[0].Vals[0] != modifyRequest.ReplaceAttributes[0].Vals[0] {
		return fmt.Errorf("expected modifyRequest of %s, but received %s", f.modifyRequestToExpect, modifyRequest)
	}
	return nil
}

func (f *fakeLDAPConnection) Search(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error) {
	if f.searchRequestToExpect.BaseDN != searchRequest.BaseDN {
		return nil, fmt.Errorf("expected searchRequest of %v, but received %v", f.searchRequestToExpect, searchRequest)
	}
	if f.searchRequestToExpect.Scope != searchRequest.Scope {
		return nil, fmt.Errorf("expected searchRequest of %v, but received %v", f.searchRequestToExpect, searchRequest)
	}
	if f.searchRequestToExpect.Filter != searchRequest.Filter {
		return nil, fmt.Errorf("expected searchRequest of %v, but received %v", f.searchRequestToExpect, searchRequest)
	}
	return f.searchResultToReturn, nil
}

func (f *fakeLDAPConnection) StartTLS(config *tls.Config) error {
	return nil
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
