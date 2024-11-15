// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"net/mail"
	"regexp"

	version "github.com/hashicorp/go-version"
)

// A regular expression used to validate common string ID patterns.
var reStringID = regexp.MustCompile(`^[^/\s]+$`)

// validEmail checks if the given input is a correct email
func validEmail(v string) bool {
	_, err := mail.ParseAddress(v)
	return err == nil
}

// validString checks if the given input is present and non-empty.
func validString(v *string) bool {
	return v != nil && *v != ""
}

// validStringID checks if the given string pointer is non-nil and
// contains a typical string identifier.
func validStringID(v *string) bool {
	return v != nil && reStringID.MatchString(*v)
}

// validVersion checks if the given input is a valid version.
func validVersion(v string) bool {
	_, err := version.NewVersion(v)
	return err == nil
}
