/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export const SELECTORS = {
  // global urls
  issuingCertificatesLabel: '[data-test-row-label="Issuing certificates"]',
  issuingCertificatesRowVal: '[data-test-row-value="Issuing certificates"]',
  crlDistributionPointsLabel: '[data-test-row-label="CRL distribution points"]',
  crlDistributionPointsRowVal: '[data-test-row-value="CRL distribution points"]',
  // crl
  expiryLabel: '[data-test-row-label="Expiry"]',
  expiryRowVal: '[data-test-row-value="Expiry"]',
  rebuildLabel: '[data-test-row-label="Auto-rebuild"]',
  rebuildRowVal: '[data-test-row-value="Auto-rebuild"]',
  responderApiLabel: '[data-test-row-label="Responder APIs"]',
  responderApiRowVal: '[data-test-row-value="Responder APIs"]',
  intervalLabel: '[data-test-row-label="Interval"]',
  intervalRowVal: '[data-test-row-value="Interval"]',
  // mount configuration
  engineTypeLabel: '[data-test-row-label="Secret engine type"]',
  engineTypeRowVal: '[data-test-row-value="Secret engine type"]',
  pathLabel: '[data-test-row-label="Path"]',
  pathRowVal: '[data-test-row-value="Path"]',
  accessorLabel: '[data-test-row-label="Accessor"]',
  accessorRowVal: '[data-test-row-value="Accessor"]',
  localLabel: '[data-test-row-label="Local"]',
  localRowVal: '[data-test-value-div="Local"]',
  sealWrapLabel: '[data-test-row-label="Seal wrap"]',
  sealWrapRowVal: '[data-test-value-div="Seal wrap"]',
  maxLeaseTtlLabel: '[data-test-row-label="Max lease TTL"]',
  maxLeaseTtlRowVal: '[data-test-row-value="Max lease TTL"]',
  allowedManagedKeysLabel: '[data-test-row-label="Allowed managed keys"]',
  allowedManagedKeysRowVal: '[data-test-value-div="Allowed managed keys"]',
};
