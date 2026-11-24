// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"errors"
	"fmt"
	"strings"
)

// EphemeralEntityIDPrefix is the reserved prefix for ephemeral entity IDs.
const EphemeralEntityIDPrefix = "eid."

// ErrEphemeralEntityID indicates that an operation cannot use an ephemeral entity ID.
var ErrEphemeralEntityID = errors.New("ephemeral entity IDs cannot be used for this operation")

// IsEphemeralEntityID reports whether id begins with the reserved ephemeral
// entity ID prefix.
func IsEphemeralEntityID(id string) bool {
	return strings.HasPrefix(id, EphemeralEntityIDPrefix)
}

// errEphemeralEntityID returns ErrEphemeralEntityID annotated with the rejected ID.
func errEphemeralEntityID(id string) error {
	return fmt.Errorf("%w: %q", ErrEphemeralEntityID, id)
}
