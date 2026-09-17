// Copyright IBM Corp. 2026
// SPDX-License-Identifier: MPL-2.0

package blackbox

import (
	"errors"
)

func (s *Session) sealed() (bool, error) {
	sealStatus, err := s.Client.Logical().Read("sys/seal-status")
	if err != nil {
		return true, err
	}

	if sealStatus == nil {
		return true, errors.New("seal status response was nil")
	}

	// Verify cluster is unsealed
	sealed, ok := sealStatus.Data["sealed"].(bool)
	if !ok {
		return true, errors.New("could not determine seal status")
	}

	return sealed, nil
}
