// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap"
)

var (
	ldapServer string = "localhost"
	ldapPort   uint16 = 389
	user       string = "*"
	passwd     string = "*"
	dn         string = "uid=*,cn=*,dc=*,dc=*"
	attribute  string = "uid"
	value      string = "username"
)

func main() {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapServer, ldapPort))
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	defer l.Close()
	// l.Debug = true

	err = l.Bind(user, passwd)
	if err != nil {
		log.Printf("ERROR: Cannot bind: %s\n", err.Error())
		return
	}

	fmt.Println(l.Compare(dn, attribute, value))
}
