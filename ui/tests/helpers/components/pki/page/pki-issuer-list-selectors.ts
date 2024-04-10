/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const PKI_ISSUER_LIST = {
  // Page::PkiIssuerList
  issuerListItem: (id: string) => `[data-test-issuer-list="${id}"]`,
  importIssuerLink: '[data-test-generate-issuer="import"]',
  generateIssuerDropdown: '[data-test-issuer-generate-dropdown]',
  generateIssuerRoot: '[data-test-generate-issuer="root"]',
  generateIssuerIntermediate: '[data-test-generate-issuer="intermediate"]',
  issuerPopupDetails: '[data-test-pki-issuer-details]',
};
