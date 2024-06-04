// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package uicustommessages

import "errors"

// FindFilter is a struct to capture the search criteria applied when searching
// for messages.
type FindFilter struct {
	// IncludeAncestors is used to look for messages in the current namespace
	// and every ancestor namespace all the way up to the root namespace.
	IncludeAncestors bool

	// MessageType indicates whether to filter messages whose type property
	// doesn't match the value set. If it's an empty string, messages are not
	// filtered by the type property.
	messageType string

	// authenticated indicates whether to filter messages whose authenticated
	// property doesn't match the value set. If it's nil, messages are not
	// filtered by the authenticated property.
	authenticated *bool

	// active indicates whether to filter messages whose active property doesn't
	// match the value set. If it's nil, messages are not filtered by the active
	// property.
	active *bool
}

// Authenticated sets the authenticated field of the receiver FindFilter struct
// to the address of the provided bool value.
func (f *FindFilter) Authenticated(value bool) {
	f.authenticated = &value
}

// Active sets the active field of the receiver FindFilter struct to the address
// of the provided bool value.
func (f *FindFilter) Active(value bool) {
	f.active = &value
}

// Type sets the messageType field of the receiver FindFilter struct to the
// provided value if it matches one of the allowed message types.
func (f *FindFilter) Type(value string) error {
	for _, el := range allowedMessageTypes {
		if value == el {
			f.messageType = value
			return nil
		}
	}

	return errors.New("unrecognized type value for filter")
}
