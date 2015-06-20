// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"gopkg.in/ldap.v1"
)

var (
	LdapServer string   = "localhost"
	LdapPort   uint16   = 389
	BaseDN     string   = "dc=enterprise,dc=org"
	Filter     string   = "(cn=kirkj)"
	Attributes []string = []string{"mail"}
)

func main() {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", LdapServer, LdapPort))
	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}
	defer l.Close()
	l.Debug = true

	search := ldap.NewSearchRequest(
		BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		Filter,
		Attributes,
		nil)

	// First search without tls.
	sr, err := l.Search(search)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
	}

	log.Printf("Search: %s -> num of entries = %d\n", search.Filter, len(sr.Entries))
	sr.PrettyPrint(0)

	// Then startTLS
	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}

	sr, err = l.Search(search)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
	}

	log.Printf("Search: %s -> num of entries = %d\n", search.Filter, len(sr.Entries))
	sr.PrettyPrint(0)
}
