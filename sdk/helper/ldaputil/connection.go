// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ldaputil

import (
	"crypto/tls"
	"time"

	"github.com/go-ldap/ldap/v3"
)

// Connection provides the functionality of an LDAP connection,
// but through an interface.
type Connection interface {
	Bind(username, password string) error
	Close() error
	Add(addRequest *ldap.AddRequest) error
	Modify(modifyRequest *ldap.ModifyRequest) error
	Del(delRequest *ldap.DelRequest) error
	Search(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error)
	StartTLS(config *tls.Config) error
	SetTimeout(timeout time.Duration)
	UnauthenticatedBind(username string) error
}

type PagingConnection interface {
	Connection
	SearchWithPaging(searchRequest *ldap.SearchRequest, pagingSize uint32) (*ldap.SearchResult, error)
}
