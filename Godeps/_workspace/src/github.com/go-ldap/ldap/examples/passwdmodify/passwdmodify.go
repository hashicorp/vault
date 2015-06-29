package main

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap"
)

// Example password policy. For this test pwdMinAge is 0 or subsequent password
// changes will fail.
//
// dn: cn=default,ou=policies,dc=enterprise,dc=org
// objectClass: pwdPolicy
// objectClass: person
// objectClass: top
// cn: default
// pwdAllowUserChange: TRUE
// pwdAttribute: userPassword
// pwdCheckQuality: 2
// pwdExpireWarning: 300
// pwdFailureCountInterval: 30
// pwdGraceAuthNLimit: 5
// pwdInHistory: 0
// pwdLockout: TRUE
// pwdLockoutDuration: 0
// pwdMaxAge: 300
// pwdMaxFailure: 0
// pwdMinAge: 0
// pwdMinLength: 5
// pwdMustChange: TRUE
// pwdSafeModify: TRUE
// sn: dummy value

var (
	ldapServer    string = "localhost"
	ldapPort      uint16 = 389
	baseDN        string = "dc=enterprise,dc=org"
	adminUser     string = "cn=admin,dc=enterprise,dc=org"
	adminPassword string = "*"
	user          string = "cn=kirkj,ou=crew,dc=enterprise,dc=org"
	oldPassword   string = "*"
	password1     string = "password123"
	password2     string = "password1234"
)

const (
	debug = false
)

func login(user string, password string) (*ldap.Conn, error) {

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapServer, ldapPort))
	if err != nil {
		return nil, err
	}
	l.Debug = debug

	bindRequest := ldap.NewSimpleBindRequest(user, password, nil)
	_, err = l.SimpleBind(bindRequest)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func main() {

	// Login as the admin and change the password of an user (without providing the old password)
	log.Printf("Logging in as the admin and changing the password of user (without providing the old password")
	l, err := login(adminUser, adminPassword)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}

	passwordModifyRequest := ldap.NewPasswordModifyRequest(user, "", password1)
	_, err = l.PasswordModify(passwordModifyRequest)

	if err != nil {
		l.Close()
		log.Fatalf("ERROR: Cannot change password: %s\n", err)
	}

	log.Printf("Done")
	l.Close()

	// Login as the user and change the password without providing a new password.
	log.Printf("Logging in as the user and changing the password without providing a new one")
	l, err = login(user, password1)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}

	passwordModifyRequest = ldap.NewPasswordModifyRequest("", password1, "")
	passwordModifyResponse, err := l.PasswordModify(passwordModifyRequest)

	if err != nil {
		l.Close()
		log.Fatalf("ERROR: Cannot change password: %s\n", err)
	}

	generatedPassword := passwordModifyResponse.GeneratedPassword
	log.Printf("Done. Generated password: %s\n", generatedPassword)

	l.Close()

	// Login as the user with the generated password and change it to another one
	log.Printf("Logging in as the user and changing the password")
	l, err = login(user, generatedPassword)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}

	passwordModifyRequest = ldap.NewPasswordModifyRequest("", generatedPassword, password2)
	_, err = l.PasswordModify(passwordModifyRequest)

	if err != nil {
		l.Close()
		log.Fatalf("ERROR: Cannot change password: %s\n", err)
	}

	log.Printf("Done")
	l.Close()

}
