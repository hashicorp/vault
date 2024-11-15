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
// there are more fields in LDAP because schema can be defined by the user.
// Here are some of the more common fields.
type fieldRegistry struct {
	CommonName         *Field `ldap:"cn"`
	DisplayName        *Field `ldap:"displayName"`
	DistinguishedName  *Field `ldap:"distinguishedName"`
	DomainComponent    *Field `ldap:"dc"`
	DomainName         *Field `ldap:"dn"`
	Name               *Field `ldap:"name"`
	ObjectCategory     *Field `ldap:"objectCategory"`
	ObjectClass        *Field `ldap:"objectClass"`
	ObjectGUID         *Field `ldap:"objectGUID"`
	ObjectSID          *Field `ldap:"objectSid"`
	OrganizationalUnit *Field `ldap:"ou"`
	PasswordLastSet    *Field `ldap:"passwordLastSet"`
	RACFID             *Field `ldap:"racfid"`
	RACFPassword       *Field `ldap:"racfPassword"`
	RACFAttributes     *Field `ldap:"racfAttributes"`
	SAMAccountName     *Field `ldap:"sAMAccountName"`
	UnicodePassword    *Field `ldap:"unicodePwd"`
	UID                *Field `ldap:"uid"`
	UserPassword       *Field `ldap:"userPassword"`
	UserPrincipalName  *Field `ldap:"userPrincipalName"`

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
