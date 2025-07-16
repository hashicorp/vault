// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package parsing

import (
	"encoding/pem"
	"errors"
	"strings"
)

func DecodePem(certBytes []byte) (*pem.Block, error) {
	block, extra := pem.Decode(certBytes)
	if block == nil {
		return nil, errors.New("invalid PEM")
	}
	if len(strings.TrimSpace(string(extra))) > 0 {
		return nil, errors.New("trailing PEM data")
	}
	return block, nil
}
