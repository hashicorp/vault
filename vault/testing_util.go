// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/cert_count"
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

func (c *Core) StopPkiCertificateCountConsumerJob() {
	mgr := c.certCountManager.(cert_count.CertificateCountManager)
	mgr.StopConsumerJob()
}

func (c *Core) ResetPkiCertificateCounts() logical.CertCount {
	mgr := c.certCountManager.(cert_count.CertificateCountManager)
	return mgr.GetCounts()
}

func (c *Core) RequirePkiCertificateCounts(t testing.TB, baseline *logical.CertCount, expectedIssuedCount, expectedStoredCount int) {
	t.Helper()
	mgr := c.certCountManager.(cert_count.CertificateCountManager)
	actualCount := mgr.GetCounts()

	actualCount.IssuedCerts -= baseline.IssuedCerts
	actualCount.StoredCerts -= baseline.StoredCerts

	baseline.IssuedCerts += uint64(expectedIssuedCount)
	baseline.StoredCerts += uint64(expectedStoredCount)

	require.Equal(t, expectedIssuedCount, int(actualCount.IssuedCerts), "PKI certificate issued count mismatch")
	require.Equal(t, expectedStoredCount, int(actualCount.StoredCerts), "PKI certificate stored count mismatch")
}
