// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package agent

import (
	"bytes"
	"encoding/json"
	"os"

	"github.com/hashicorp/vault/sdk/logical"
)

const envVarRunAccTests = "VAULT_ACC"

var runAcceptanceTests = os.Getenv(envVarRunAccTests) == "1"

func readToken(fileName string) (*logical.HTTPWrapInfo, error) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	wrapper := &logical.HTTPWrapInfo{}
	if err := json.NewDecoder(bytes.NewReader(b)).Decode(wrapper); err != nil {
		return nil, err
	}
	return wrapper, nil
}
