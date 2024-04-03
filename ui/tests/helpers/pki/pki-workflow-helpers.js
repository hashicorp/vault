/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const PKI_WORKFLOW = {
  // ROLES
  deleteRoleButton: '[data-test-pki-role-delete]',
  generateCertLink: '[data-test-pki-role-generate-cert]',
  signCertLink: '[data-test-pki-role-sign-cert]',
  editRoleLink: '[data-test-pki-role-edit-link]',
  createRoleLink: '[data-test-pki-role-create-link]',
  // ISSUERS
  issuerListItem: (id) => `[data-test-issuer-list="${id}"]`,
  importIssuerLink: '[data-test-generate-issuer="import"]',
  generateIssuerDropdown: '[data-test-issuer-generate-dropdown]',
  generateIssuerRoot: '[data-test-generate-issuer="root"]',
  generateIssuerIntermediate: '[data-test-generate-issuer="intermediate"]',
  issuerPopupMenu: '[data-test-popup-menu-trigger]',
  issuerPopupDetails: '[data-test-popup-menu-details]',
};

// Clears pki-related data and capabilities so that admin
// capabilities from setup don't rollover
export function clearRecords(store) {
  store.unloadAll('pki/action');
  store.unloadAll('pki/issuer');
  store.unloadAll('pki/key');
  store.unloadAll('pki/role');
  store.unloadAll('pki/sign-intermediate');
  store.unloadAll('pki/tidy');
  store.unloadAll('pki/config/urls');
  store.unloadAll('pki/config/crl');
  store.unloadAll('pki/config/cluster');
  store.unloadAll('pki/config/acme');
  store.unloadAll('pki/certificate/generate');
  store.unloadAll('pki/certificate/sign');
  store.unloadAll('capabilities');
}
