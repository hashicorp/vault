// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package jwtauth

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/oauth2"
)

// IBMISAMProvider is used for IBMISAM-specific configuration
type IBMISAMProvider struct{}

// Initialize anything in the IBMISAMProvider struct - satisfying the CustomProvider interface
func (a *IBMISAMProvider) Initialize(_ context.Context, _ *jwtConfig) error {
	return nil
}

// SensitiveKeys - satisfying the CustomProvider interface
func (a *IBMISAMProvider) SensitiveKeys() []string {
	return []string{}
}

// FetchGroups - custom groups fetching for ibmisam - satisfying GroupsFetcher interface
// IBMISAM by default will return groups not as a json list but as a list of space seperated strings
// We need to convert this to a json list
func (a *IBMISAMProvider) FetchGroups(_ context.Context, b *jwtAuthBackend, allClaims map[string]interface{}, role *jwtRole, _ oauth2.TokenSource) (interface{}, error) {
	groupsClaimRaw := getClaim(b.Logger(), allClaims, role.GroupsClaim)

	if groupsClaimRaw != nil {
		// Try to convert the comma seperated list of strings into a list
		if groupsstr, ok := groupsClaimRaw.(string); ok {
			rawibmisamGroups := strings.Split(groupsstr, " ")

			ibmisamGroups := make([]interface{}, 0, len(rawibmisamGroups))
			for group := range rawibmisamGroups {
				ibmisamGroups = append(ibmisamGroups, rawibmisamGroups[group])
			}
			groupsClaimRaw = ibmisamGroups
		}
	}
	b.Logger().Debug(fmt.Sprintf("post: groups claim raw is %v", groupsClaimRaw))
	return groupsClaimRaw, nil
}
