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
