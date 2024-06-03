// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package postgresql

import "fmt"

// passwordAuthentication determines whether to send passwords in plaintext (password) or hashed (scram-sha-256).
type passwordAuthentication string

var (
	// passwordAuthenticationPassword is the default. If set, passwords will be sent to PostgreSQL in plain text.
	passwordAuthenticationPassword    passwordAuthentication = "password"
	passwordAuthenticationSCRAMSHA256 passwordAuthentication = "scram-sha-256"
)

var passwordAuthentications = map[passwordAuthentication]struct{}{
	passwordAuthenticationSCRAMSHA256: {},
	passwordAuthenticationPassword:    {},
}

func parsePasswordAuthentication(s string) (passwordAuthentication, error) {
	if _, ok := passwordAuthentications[passwordAuthentication(s)]; !ok {
		return "", fmt.Errorf("'%s' is not a valid password authentication type", s)
	}

	return passwordAuthentication(s), nil
}
