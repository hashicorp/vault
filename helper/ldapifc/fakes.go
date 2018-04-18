package ldapifc

import (
	"crypto/tls"
	"fmt"

	"github.com/go-ldap/ldap"
)

// FakeLDAPClient can be used to inspect the LDAP requests that have been constructed,
// and to inject responses.
type FakeLDAPClient struct {
	ConnToReturn Connection
}

func (f *FakeLDAPClient) Dial(network, addr string) (Connection, error) {
	return f.ConnToReturn, nil
}

func (f *FakeLDAPClient) DialTLS(network, addr string, config *tls.Config) (Connection, error) {
	return f.ConnToReturn, nil
}

type FakeLDAPConnection struct {
	ModifyRequestToExpect *ldap.ModifyRequest
	SearchRequestToExpect *ldap.SearchRequest
	SearchResultToReturn  *ldap.SearchResult
}

func (f *FakeLDAPConnection) Bind(username, password string) error {
	return nil
}

func (f *FakeLDAPConnection) Close() {}

func (f *FakeLDAPConnection) Modify(modifyRequest *ldap.ModifyRequest) error {

	if f.ModifyRequestToExpect.DN != modifyRequest.DN {
		return fmt.Errorf("expected modifyRequest of %s, but received %s", f.ModifyRequestToExpect, modifyRequest)
	}

	if len(f.ModifyRequestToExpect.ReplaceAttributes) != len(modifyRequest.ReplaceAttributes) {
		return fmt.Errorf("expected modifyRequest of %s, but received %s", f.ModifyRequestToExpect, modifyRequest)
	}

	if f.ModifyRequestToExpect.ReplaceAttributes[0].Type != modifyRequest.ReplaceAttributes[0].Type {
		return fmt.Errorf("expected modifyRequest of %s, but received %s", f.ModifyRequestToExpect, modifyRequest)
	}

	if len(f.ModifyRequestToExpect.ReplaceAttributes[0].Vals) != len(modifyRequest.ReplaceAttributes[0].Vals) {
		return fmt.Errorf("expected modifyRequest of %s, but received %s", f.ModifyRequestToExpect, modifyRequest)
	}

	if f.ModifyRequestToExpect.ReplaceAttributes[0].Vals[0] != modifyRequest.ReplaceAttributes[0].Vals[0] {
		return fmt.Errorf("expected modifyRequest of %s, but received %s", f.ModifyRequestToExpect, modifyRequest)
	}

	return nil
}

func (f *FakeLDAPConnection) Search(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error) {

	if f.SearchRequestToExpect.BaseDN != searchRequest.BaseDN {
		return nil, fmt.Errorf("expected searchRequest of %v, but received %v", f.SearchRequestToExpect, searchRequest)
	}

	if f.SearchRequestToExpect.Scope != searchRequest.Scope {
		return nil, fmt.Errorf("expected searchRequest of %v, but received %v", f.SearchRequestToExpect, searchRequest)
	}

	if f.SearchRequestToExpect.Filter != searchRequest.Filter {
		return nil, fmt.Errorf("expected searchRequest of %v, but received %v", f.SearchRequestToExpect, searchRequest)
	}

	return f.SearchResultToReturn, nil
}

func (f *FakeLDAPConnection) StartTLS(config *tls.Config) error {
	return nil
}
