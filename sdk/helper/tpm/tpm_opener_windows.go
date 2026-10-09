// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

//go:build windows

package tpm

import (
	"fmt"

	"github.com/google/go-attestation/attest"
)

func openTPM(_ string) (*attest.TPM, error) {
	return nil, fmt.Errorf("tpm is not supported on windows")
}
