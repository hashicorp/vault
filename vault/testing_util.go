// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/vault/pki_cert_count"
	"github.com/hashicorp/vault/version"
	"github.com/stretchr/testify/require"
)

func init() {
	// The BuildDate is set as part of the build process in CI so we need to
	// initialize it for testing. By setting it to now minus one year we
	// provide some headroom to ensure that test license expiration (for enterprise)
	// does not exceed the BuildDate as that is invalid.
	if version.BuildDate == "" {
		version.BuildDate = time.Now().UTC().AddDate(-1, 0, 0).Format(time.RFC3339)
	}
}

func (c *TestClusterCore) StopPkiCertificateCountConsumerJob() {
	mgr := c.Core.pkiCertCountManager.(pki_cert_count.PkiCertificateCountManager)
	mgr.StopConsumerJob()
}

func (c *TestClusterCore) ResetPkiCertificateCounts() {
	mgr := c.Core.pkiCertCountManager.(pki_cert_count.PkiCertificateCountManager)

	c.pkiCertificateCountData = mgr.GetCounts()
}

func (c *TestClusterCore) RequirePkiCertificateCounts(t testing.TB, expectedIssuedCount, expectedStoredCount int) {
	t.Helper()
	mgr := c.Core.pkiCertCountManager.(pki_cert_count.PkiCertificateCountManager)
	actualCount := mgr.GetCounts()

	actualCount.IssuedCerts -= c.pkiCertificateCountData.IssuedCerts
	actualCount.StoredCerts -= c.pkiCertificateCountData.StoredCerts

	c.pkiCertificateCountData.IssuedCerts += uint64(expectedIssuedCount)
	c.pkiCertificateCountData.StoredCerts += uint64(expectedStoredCount)

	require.Equal(t, expectedIssuedCount, int(actualCount.IssuedCerts), "PKI certificate issued count mismatch")
	require.Equal(t, expectedStoredCount, int(actualCount.StoredCerts), "PKI certificate stored count mismatch")
}
