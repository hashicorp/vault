package client

import (
	"fmt"
	"github.com/hashicorp/vault/sdk/helper/strutil"
)

// SupportedSchemas returns a slice of different OpenLDAP schemas supported
// by the plugin.  This is used to change the FieldRegistry when modifying
// user passwords.
func SupportedSchemas() []string {
	return []string{"openldap", "racf"}
}

// ValidSchema checks if the configured schema is supported by the plugin.
func ValidSchema(schema string) bool {
	return strutil.StrListContains(SupportedSchemas(), schema)
}

// GetSchemaFieldRegistry type switches field registries depending on the configured schema.
// For example, IBM RACF has a custom OpenLDAP schema so the password is stored in a different
// attribute.
func GetSchemaFieldRegistry(schema string, newPassword string) (map[*Field][]string, error) {
	switch schema {
	case "openldap":
		fields := map[*Field][]string{FieldRegistry.UserPassword: {newPassword}}
		return fields, nil
	case "racf":
		fields := map[*Field][]string{
			FieldRegistry.RACFPassword:   {newPassword},
			FieldRegistry.RACFAttributes: {"noexpire"},
		}
		return fields, nil
	default:
		return nil, fmt.Errorf("configured schema %s not valid", schema)
	}
}
