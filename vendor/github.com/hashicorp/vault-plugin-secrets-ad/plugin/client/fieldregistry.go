// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"reflect"
)

// FieldRegistry is designed to look and feel
// like an enum from another language like Python.
//
// Example: Accessing constants
//
//	FieldRegistry.AccountExpires
//	FieldRegistry.BadPasswordCount
//
// Example: Utility methods
//
//	FieldRegistry.List()
//	FieldRegistry.Parse("givenName")
var FieldRegistry = newFieldRegistry()

// newFieldRegistry iterates through all the fields in the registry,
// pulls their ldap strings, and sets up each field to contain its ldap string
func newFieldRegistry() *fieldRegistry {
	reg := &fieldRegistry{}
	vOfReg := reflect.ValueOf(reg)

	registryFields := vOfReg.Elem()
	for i := 0; i < registryFields.NumField(); i++ {

		if registryFields.Field(i).Kind() == reflect.Ptr {

			field := registryFields.Type().Field(i)
			ldapString := field.Tag.Get("ldap")
			ldapField := &Field{ldapString}
			vOfLDAPField := reflect.ValueOf(ldapField)

			registryFields.FieldByName(field.Name).Set(vOfLDAPField)

			reg.fieldList = append(reg.fieldList, ldapField)
		}
	}
	return reg
}

// fieldRegistry isn't currently intended to be an exhaustive list -
// there are more fields in ActiveDirectory. However, these are the ones
// that may be useful to Vault. Feel free to add to this list!
type fieldRegistry struct {
	AccountExpires              *Field `ldap:"accountExpires"`
	AdminCount                  *Field `ldap:"adminCount"`
	BadPasswordCount            *Field `ldap:"badPwdCount"`
	BadPasswordTime             *Field `ldap:"badPasswordTime"`
	CodePage                    *Field `ldap:"codePage"`
	CommonName                  *Field `ldap:"cn"`
	CountryCode                 *Field `ldap:"countryCode"`
	DisplayName                 *Field `ldap:"displayName"`
	DistinguishedName           *Field `ldap:"distinguishedName"`
	DomainComponent             *Field `ldap:"dc"`
	DomainName                  *Field `ldap:"dn"`
	DSCorePropogationData       *Field `ldap:"dSCorePropagationData"`
	GivenName                   *Field `ldap:"givenName"`
	GroupType                   *Field `ldap:"groupType"`
	Initials                    *Field `ldap:"initials"`
	InstanceType                *Field `ldap:"instanceType"`
	LastLogoff                  *Field `ldap:"lastLogoff"`
	LastLogon                   *Field `ldap:"lastLogon"`
	LastLogonTimestamp          *Field `ldap:"lastLogonTimestamp"`
	LockoutTime                 *Field `ldap:"lockoutTime"`
	LogonCount                  *Field `ldap:"logonCount"`
	MemberOf                    *Field `ldap:"memberOf"`
	Name                        *Field `ldap:"name"`
	ObjectCategory              *Field `ldap:"objectCategory"`
	ObjectClass                 *Field `ldap:"objectClass"`
	ObjectGUID                  *Field `ldap:"objectGUID"`
	ObjectSID                   *Field `ldap:"objectSid"`
	OrganizationalUnit          *Field `ldap:"ou"`
	PasswordLastSet             *Field `ldap:"pwdLastSet"`
	PrimaryGroupID              *Field `ldap:"primaryGroupID"`
	SAMAccountName              *Field `ldap:"sAMAccountName"`
	SAMAccountType              *Field `ldap:"sAMAccountType"`
	Surname                     *Field `ldap:"sn"`
	UnicodePassword             *Field `ldap:"unicodePwd"`
	UpdateSequenceNumberChanged *Field `ldap:"uSNChanged"`
	UpdateSequenceNumberCreated *Field `ldap:"uSNCreated"`
	UserAccountControl          *Field `ldap:"userAccountControl"`
	UserPrincipalName           *Field `ldap:"userPrincipalName"`
	WhenCreated                 *Field `ldap:"whenCreated"`
	WhenChanged                 *Field `ldap:"whenChanged"`

	fieldList []*Field
}

func (r *fieldRegistry) List() []*Field {
	return r.fieldList
}

func (r *fieldRegistry) Parse(s string) *Field {
	for _, f := range r.List() {
		if f.String() == s {
			return f
		}
	}
	return nil
}

type Field struct {
	str string
}

func (f *Field) String() string {
	return f.str
}
