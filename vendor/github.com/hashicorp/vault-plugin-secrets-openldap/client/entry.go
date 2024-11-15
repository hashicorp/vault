// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"strings"

	"github.com/go-ldap/ldap/v3"
)

// Entry is an LDAP-specific construct
// to make knowing and grabbing fields more convenient,
// while retaining all original information.
func NewEntry(ldapEntry *ldap.Entry) *Entry {
	fieldMap := make(map[string][]string)
	for _, attribute := range ldapEntry.Attributes {
		field := FieldRegistry.Parse(attribute.Name)
		if field == nil {
			// This field simply isn't in the registry, no big deal.
			continue
		}
		fieldMap[field.String()] = attribute.Values
	}
	return &Entry{fieldMap: fieldMap, Entry: ldapEntry}
}

type Entry struct {
	*ldap.Entry
	fieldMap map[string][]string
}

func (e *Entry) Get(field *Field) ([]string, bool) {
	values, found := e.fieldMap[field.String()]
	return values, found
}

func (e *Entry) GetJoined(field *Field) (string, bool) {
	values, found := e.Get(field)
	if !found {
		return "", false
	}
	return strings.Join(values, ","), true
}
