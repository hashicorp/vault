// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package issuing

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/hashicorp/vault/sdk/helper/certutil"
)

func AddCertificateFieldsToAuditData(auditData map[string]any, certificateType string, certificate *x509.Certificate) error {
	if auditData == nil {
		return fmt.Errorf("no audit data on which to add certificate fields")
	}
	if certificate == nil {
		return fmt.Errorf("no certificate whose fields should be added to audit data")
	}
	certificateFieldPrefix := ""
	if certificateType != "" {
		certificateFieldPrefix = certificateType + "_"
	}
	auditData[certificateFieldPrefix+"certificate"] = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificate.Raw}))
	certificateFields, err := certutil.ParseCertificateToFields(*certificate)
	if err != nil {
		return fmt.Errorf("error parsing certificate fields for audit data: %v", err)
	}
	for field, value := range certificateFields {
		switch field {
		case "common_name", "serial_number", "ttl", "skid", "key_type":
			auditData[certificateFieldPrefix+field] = value
		default:
			continue
		}
	}
	return nil
}
