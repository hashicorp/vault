// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package blackbox

import (
	"fmt"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
)

// TestPKI_IssueCertificate verifies PKI engine functionality by setting up a root CA,
// issuing a certificate with specific parameters, and validating the certificate
// response contains the expected fields and values.
func TestPKI_IssueCertificate(t *testing.T) {
	v := blackbox.New(t)

	roleName := v.MustSetupPKIRoot("pki")

	// issue a cert
	issuePath := fmt.Sprintf("pki/issue/%s", roleName)
	payload := map[string]any{
		"common_name": "api.example.com",
		"ttl":         "24h",
	}

	var secret *api.Secret
	v.Eventually(func() error {
		var err error
		secret, err = v.Client.Logical().Write(issuePath, payload)
		return err
	})

	if secret == nil {
		t.Fatal("Expected certificate secret, got nil")
	}

	assertions := v.AssertSecret(secret)
	assertions.Data().
		HasKeyExists("certificate").
		HasKeyExists("issuing_ca").
		HasKeyExists("private_key").
		HasKeyCustom("serial_number", func(val any) bool {
			s, ok := val.(string)
			return ok && len(s) > 0
		})
}
