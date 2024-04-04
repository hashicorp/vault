/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const PKI_BASE_URL = `/vault/cluster/secrets/backend/pki/roles`;

export const PKI_ROLE_FORM = {
  domainHandling: '[data-test-toggle-group="Domain handling"]',
  keyParams: '[data-test-toggle-group="Key parameters"]',
  keyUsage: '[data-test-toggle-group="Key usage"]',
  digitalSignature: '[data-test-checkbox="DigitalSignature"]',
  keyAgreement: '[data-test-checkbox="KeyAgreement"]',
  keyEncipherment: '[data-test-checkbox="KeyEncipherment"]',
  any: '[data-test-checkbox="Any"]',
  serverAuth: '[data-test-checkbox="ServerAuth"]',
  policyIdentifiers: '[data-test-toggle-group="Policy identifiers"]',
  san: '[data-test-toggle-group="Subject Alternative Name (SAN) Options"]',
  additionalSubjectFields: '[data-test-toggle-group="Additional subject fields"]',
};
