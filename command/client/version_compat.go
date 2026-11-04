// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import "strings"

// IsCompatibleVersion reports whether clientVersion and serverVersion match.
// A vault-client (ent.client) build also matches any enterprise server edition
// of the same version.
func IsCompatibleVersion(clientVersion, serverVersion string) bool {
	clientV, clientEdition, _ := strings.Cut(clientVersion, "+")
	serverV, serverEdition, _ := strings.Cut(serverVersion, "+")
	if clientV != serverV {
		// The version or prerelease doesn't match
		return false
	}

	if clientEdition == "ent.client" {
		// The client is compatible with all ent editions
		return strings.HasPrefix(serverEdition, "ent")
	}

	return clientEdition == serverEdition
}
