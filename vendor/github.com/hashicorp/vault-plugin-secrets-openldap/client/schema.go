// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"golang.org/x/text/encoding/unicode"
)

const (
	SchemaOpenLDAP = "openldap"
	SchemaAD       = "ad"
	SchemaRACF     = "racf"
)

// SupportedSchemas returns a slice of different LDAP schemas supported
// by the plugin. This is used to change the FieldRegistry when modifying
// user passwords and to set the default user attribute (userattr).
func SupportedSchemas() []string {
	return []string{SchemaOpenLDAP, SchemaRACF, SchemaAD}
}

// ValidSchema checks if the configured schema is supported by the plugin.
func ValidSchema(schema string) bool {
	return strutil.StrListContains(SupportedSchemas(), schema)
}

// GetSchemaFieldRegistry type switches field registries depending on the configured schema.
// For example, IBM RACF has a custom LDAP schema so the password is stored in a different
// attribute.
func GetSchemaFieldRegistry(schema string, newPassword string) (map[*Field][]string, error) {
	switch schema {
	case SchemaOpenLDAP:
		fields := map[*Field][]string{FieldRegistry.UserPassword: {newPassword}}
		return fields, nil
	case SchemaRACF:
		fields := map[*Field][]string{
			FieldRegistry.RACFPassword:   {newPassword},
			FieldRegistry.RACFAttributes: {"noexpired"},
		}
		return fields, nil
	case SchemaAD:
		pwdEncoded, err := formatPassword(newPassword)
		if err != nil {
			return nil, err
		}
		fields := map[*Field][]string{FieldRegistry.UnicodePassword: {pwdEncoded}}
		return fields, nil
	default:
		return nil, fmt.Errorf("configured schema %s not valid", schema)
	}
}

// According to the MS docs, the password needs to be utf16 and enclosed in quotes.
func formatPassword(original string) (string, error) {
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	return utf16.NewEncoder().String("\"" + original + "\"")
}
