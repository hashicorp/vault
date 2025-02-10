// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestManualChainValidation(t *testing.T) {
	// Set Up a Cluster
	cluster, client := setupTestPkiCluster(t)
	defer cluster.Cleanup()

	// Set Up Root-A
	mount := "pki"

	resp, err := client.Logical().WriteWithContext(context.Background(), mount+"/issuers/generate/root/internal",
		map[string]interface{}{
			"issuer_name": "rootCa",
			"key_name":    "root-key",
			"key_type":    "ec",
			"common_name": "Test Root",
			"ttl":         "7200h",
		})
	require.NoError(t, err, "failed creating root CA")

	// Set Up Int-A
	resp, err = client.Logical().WriteWithContext(context.Background(), mount+"/issuers/generate/intermediate/internal",
		map[string]interface{}{
			"key_name":    "int-a-key",
			"key_type":    "ec",
			"common_name": "Test Int A",
		})
	require.NoError(t, err, "failed creating intermediary CSR")
	intermediateCSR := resp.Data["csr"].(string)

	// Sign the intermediate CSR
	resp, err = client.Logical().Write(mount+"/issuer/rootCa/sign-intermediate", map[string]interface{}{
		"csr": intermediateCSR,
		"ttl": "7100h",
	})
	require.NoError(t, err, "failed signing intermediary CSR")
	intermediateCertPEM := resp.Data["certificate"].(string)

	// Import the intermediate cert
	resp, err = client.Logical().Write(mount+"/issuers/import/cert", map[string]interface{}{
		"pem_bundle": intermediateCertPEM,
	})
	require.NoError(t, err, "failed importing intermediary cert")
	importedIssuersRaw := resp.Data["imported_issuers"].([]interface{})
	require.Len(t, importedIssuersRaw, 1)
	intCaUuid := importedIssuersRaw[0].(string)

	_, err = client.Logical().Write(mount+"/issuer/"+intCaUuid, map[string]interface{}{
		"issuer_name": "intA",
	})
	require.NoError(t, err, "failed updating issuer name")

	// Set Up Int-B (Child of Int-A)
	resp, err = client.Logical().WriteWithContext(context.Background(), mount+"/issuers/generate/intermediate/internal",
		map[string]interface{}{
			"key_name":    "int-b-key",
			"key_type":    "ec",
			"common_name": "Test Int B",
		})
	require.NoError(t, err, "failed creating intermediary CSR")
	subIntermediateCSR := resp.Data["csr"].(string)

	// Sign the intermediate CSR
	resp, err = client.Logical().Write(mount+"/issuer/intA/sign-intermediate", map[string]interface{}{
		"csr": subIntermediateCSR,
		"ttl": "7100h",
	})
	require.NoError(t, err, "failed signing intermediary CSR")
	subIntermediateCertPEM := resp.Data["certificate"].(string)

	// Import the intermediate cert
	resp, err = client.Logical().Write(mount+"/issuers/import/cert", map[string]interface{}{
		"pem_bundle": subIntermediateCertPEM,
	})
	require.NoError(t, err, "failed importing intermediary cert")
	subImportedIssuersRaw := resp.Data["imported_issuers"].([]interface{})
	require.Len(t, subImportedIssuersRaw, 1)
	subIntCaUuid := subImportedIssuersRaw[0].(string)

	resp, err = client.Logical().Write(mount+"/issuer/"+subIntCaUuid, map[string]interface{}{
		"issuer_name": "intB",
	})
	require.NoError(t, err, "failed updating issuer name")

	// Try to Set Int-B Manual Chain to Just Be Root-A; Expect An Error
	resp, err = client.Logical().Write(mount+"/issuer/intB", map[string]interface{}{
		"manual_chain": []string{"intB", "rootCa"}, // Misses "intA" which issued "intB"
	})
	require.Error(t, err, "failed updating intermediary cert")

	resp, err = client.Logical().Read(mount + "/issuer/intB")
	require.NoError(t, err, "failed reading intermediary cert")
	require.Equal(t, resp.Data["manual_chain"], "")
}
