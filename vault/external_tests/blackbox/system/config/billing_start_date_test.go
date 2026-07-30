//go:build system
// +build system

// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
	"github.com/stretchr/testify/require"
)

// customBuildRe matches VAULT_VERSION values that contain a git short SHA
// suffix (e.g. "v1.21.0-beta1+ent-2cf0b2f"), which identifies a custom build
// produced by the test-hcp-image workflow. Release versions never have this
// suffix (e.g. "v1.21.0+ent" or "v1.21.0-beta1+ent").
var customBuildRe = regexp.MustCompile(`-[0-9a-f]{7,}$`)

func isCustomBuildVersion() bool {
	return customBuildRe.MatchString(os.Getenv("VAULT_VERSION"))
}

// readBillingStartTimestamp reads sys/internal/counters/config from the root
// namespace and returns the billing_start_timestamp value.
func readBillingStartTimestamp(v *blackbox.Session) (time.Time, error) {
	secret, err := v.WithRootNamespace(func() (*api.Secret, error) {
		return v.Client.Logical().Read("sys/internal/counters/config")
	})
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to read sys/internal/counters/config: %w", err)
	}
	if secret == nil {
		return time.Time{}, fmt.Errorf("nil response from sys/internal/counters/config")
	}

	billingStartStr, ok := secret.Data["billing_start_timestamp"].(string)
	if !ok {
		return time.Time{}, fmt.Errorf("billing_start_timestamp missing or not a string in response")
	}

	billingStart, err := time.Parse(time.RFC3339, billingStartStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse billing_start_timestamp %q as RFC3339: %w", billingStartStr, err)
	}

	return billingStart, nil
}

// checkBillingStartInCurrentYear verifies that the billing start date is
// within the current billing year (not older than 1 year and not in the future).
// This mirrors the verify_date_is_in_current_year logic in verify-billing-start.sh.
func checkBillingStartInCurrentYear(billingStart time.Time) error {
	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	if billingStart.Before(oneYearAgo) {
		return fmt.Errorf("billing start date %s is not in the current billing year (more than 1 year old, cutoff: %s)",
			billingStart.Format(time.RFC3339),
			oneYearAgo.Format(time.RFC3339))
	}
	if billingStart.After(time.Now()) {
		return fmt.Errorf("billing start date %s is in the future",
			billingStart.Format(time.RFC3339))
	}
	return nil
}

// TestBillingStartDate verifies that the billing start date is within the
// current billing year.
//
// WHERE THIS RUNS:
//   - Scenario: cloud-ent (enos/enos-scenario-cloud-ent.hcl) — regular CI runs
//     against a published HCP Vault version (dev/int/prod environments).
//   - The cloud-ent scenario includes "system/config" in blackbox_test_packages_default.
//
// SKIPPED WHEN:
//   - VAULT_VERSION contains a hex SHA suffix (e.g. "v1.21.0-beta1+ent-2cf0b2f"),
//     which indicates a custom build image from the test-hcp-image workflow.
//     Billing start date validation is not meaningful for ephemeral custom builds.
func TestBillingStartDate(t *testing.T) {
	if isCustomBuildVersion() {
		t.Skipf("skipping: VAULT_VERSION=%q is a custom build (test-hcp-image workflow); billing date validation is only meaningful for released versions", os.Getenv("VAULT_VERSION"))
	}
	// Check version and edition before blackbox.New to avoid namespace creation on unsupported builds.
	blackbox.SkipIfEdition(t, "ce")
	blackbox.SkipIfVersionBelow(t, "2.0.0")
	t.Parallel()
	v := blackbox.New(t)

	v.AssertUnsealedAny()

	// sys/internal/counters/config must be read from the root namespace.
	// The enos verify-billing-start.sh script sets no VAULT_NAMESPACE (root).
	billingStart, err := readBillingStartTimestamp(v)
	require.NoError(t, err)

	require.NoError(t, checkBillingStartInCurrentYear(billingStart),
		"billing start date validation failed")

	t.Logf("✓ Billing start date %s is within the current billing year", billingStart.Format(time.RFC3339))
}

// TestBillingStartDateRollover verifies that the billing start date has rolled
// over to the current billing year after a cluster upgrade. This test mirrors
// the enos vault_verify_billing_start_date module which retries the check up
// to 10 times (with 30s between retries) to allow time for automatic rollover
// to complete after an upgrade.
//
// WHERE THIS RUNS:
//   - Scenario: cloud-ent (enos/enos-scenario-cloud-ent.hcl) — regular CI runs only.
//   - Same as TestBillingStartDate — included via "system/config" default package.
//
// SKIPPED WHEN:
//   - VAULT_VERSION contains a hex SHA suffix — same custom build skip as TestBillingStartDate.
func TestBillingStartDateRollover(t *testing.T) {
	if isCustomBuildVersion() {
		t.Skipf("skipping: VAULT_VERSION=%q is a custom build (test-hcp-image workflow); billing date validation is only meaningful for released versions", os.Getenv("VAULT_VERSION"))
	}
	// Check version and edition before blackbox.New to avoid namespace creation on unsupported builds.
	blackbox.SkipIfEdition(t, "ce")
	blackbox.SkipIfVersionBelow(t, "2.0.0")
	t.Parallel()
	v := blackbox.New(t)

	v.AssertUnsealedAny()

	// Retry for up to 5 minutes (matching the enos module's retry 10 × 30s logic)
	// to allow billing start date rollover to complete after a cluster upgrade.
	v.EventuallyWithTimeout(func() error {
		billingStart, err := readBillingStartTimestamp(v)
		if err != nil {
			return err
		}
		return checkBillingStartInCurrentYear(billingStart)
	}, 5*time.Minute)

	// Log final value after rollover confirmed.
	billingStart, err := readBillingStartTimestamp(v)
	require.NoError(t, err)
	t.Logf("✓ Billing start date %s is within the current billing year", billingStart.Format(time.RFC3339))
}
