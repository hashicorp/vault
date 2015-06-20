// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"gopkg.in/ldap.v1"
)

var (
	ldapServer string = "localhost"
	ldapPort   uint16 = 389
	baseDN     string = "dc=enterprise,dc=org"
	user       string = "uid=kirkj,ou=crew,dc=enterprise,dc=org"
	passwd     string = "*"
)

func main() {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapServer, ldapPort))
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	defer l.Close()
	l.Debug = true

	bindRequest := ldap.NewSimpleBindRequest(user, passwd, nil)

	r, err := l.SimpleBind(bindRequest)

	passwordMustChangeControl := ldap.FindControl(r.Controls, ldap.ControlTypeVChuPasswordMustChange)
	var passwordMustChange *ldap.ControlVChuPasswordMustChange
	if passwordMustChangeControl != nil {
		passwordMustChange = passwordMustChangeControl.(*ldap.ControlVChuPasswordMustChange)
	}

	if passwordMustChange != nil && passwordMustChange.MustChange {
		log.Printf("Password Must be changed.\n")
	}

	passwordWarningControl := ldap.FindControl(r.Controls, ldap.ControlTypeVChuPasswordWarning)

	var passwordWarning *ldap.ControlVChuPasswordWarning
	if passwordWarningControl != nil {
		passwordWarning = passwordWarningControl.(*ldap.ControlVChuPasswordWarning)
	} else {
		log.Printf("ppolicyControl response not available.\n")
	}
	if err != nil {
		log.Print("ERROR: Cannot bind: " + err.Error())
	} else {
		logStr := "Login Ok"
		if passwordWarning != nil {
			if passwordWarning.Expire >= 0 {
				logStr += fmt.Sprintf(". Password expires in %d seconds\n", passwordWarning.Expire)
			}
		}
		log.Print(logStr)
	}
}
