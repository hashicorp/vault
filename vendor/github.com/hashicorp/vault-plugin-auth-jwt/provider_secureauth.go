// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jwtauth

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/oauth2"
)

// SecureAuthProvider is used for SecureAuth-specific configuration
type SecureAuthProvider struct{}

// Initialize anything in the SecureAuthProvider struct - satisfying the CustomProvider interface
func (a *SecureAuthProvider) Initialize(_ context.Context, _ *jwtConfig) error {
	return nil
}

// SensitiveKeys - satisfying the CustomProvider interface
func (a *SecureAuthProvider) SensitiveKeys() []string {
	return []string{}
}

// FetchGroups - custom groups fetching for secureauth - satisfying GroupsFetcher interface
// SecureAuth by default will return groups not as a json list but as a list of comma seperated strings
// We need to convert this to a json list
func (a *SecureAuthProvider) FetchGroups(_ context.Context, b *jwtAuthBackend, allClaims map[string]interface{}, role *jwtRole, _ oauth2.TokenSource) (interface{}, error) {
	groupsClaimRaw := getClaim(b.Logger(), allClaims, role.GroupsClaim)

	if groupsClaimRaw != nil {
		// Try to convert the comma seperated list of strings into a list
		if groupsstr, ok := groupsClaimRaw.(string); ok {
			rawsecureauthGroups := strings.Split(groupsstr, ",")

			secureauthGroups := make([]interface{}, 0, len(rawsecureauthGroups))
			for group := range rawsecureauthGroups {
				secureauthGroups = append(secureauthGroups, rawsecureauthGroups[group])
			}
			groupsClaimRaw = secureauthGroups
		}
	}
	b.Logger().Debug(fmt.Sprintf("post: groups claim raw is %v", groupsClaimRaw))
	return groupsClaimRaw, nil
}
