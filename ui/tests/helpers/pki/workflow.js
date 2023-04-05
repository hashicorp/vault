/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { SELECTORS as ROLEFORM } from './pki-role-form';
import { SELECTORS as GENERATECERT } from './pki-role-generate';
import { SELECTORS as KEYFORM } from './pki-key-form';
import { SELECTORS as KEYPAGES } from './page/pki-keys';
import { SELECTORS as ISSUERDETAILS } from './pki-issuer-details';
import { SELECTORS as CONFIGURATION } from './pki-configure-create';
import { SELECTORS as DELETE } from './pki-delete-all-issuers';

export const SELECTORS = {
  breadcrumbContainer: '[data-test-breadcrumbs]',
  breadcrumbs: '[data-test-breadcrumbs] li',
  overviewBreadcrumb: '[data-test-breadcrumbs] li:nth-of-type(2) > a',
  pageTitle: '[data-test-pki-role-page-title]',
  alertBanner: '[data-test-alert-banner="alert"]',
  emptyState: '[data-test-component="empty-state"]',
  emptyStateTitle: '[data-test-empty-state-title]',
  emptyStateLink: '.empty-state-actions a',
  emptyStateMessage: '[data-test-empty-state-message]',
  // TABS
  overviewTab: '[data-test-secret-list-tab="Overview"]',
  rolesTab: '[data-test-secret-list-tab="Roles"]',
  issuersTab: '[data-test-secret-list-tab="Issuers"]',
  certsTab: '[data-test-secret-list-tab="Certificates"]',
  keysTab: '[data-test-secret-list-tab="Keys"]',
  configTab: '[data-test-secret-list-tab="Configuration"]',
  // ROLES
  deleteRoleButton: '[data-test-pki-role-delete]',
  generateCertLink: '[data-test-pki-role-generate-cert]',
  signCertLink: '[data-test-pki-role-sign-cert]',
  editRoleLink: '[data-test-pki-role-edit-link]',
  createRoleLink: '[data-test-pki-role-create-link]',
  roleForm: {
    ...ROLEFORM,
  },
  generateCertForm: {
    ...GENERATECERT,
  },
  // KEYS
  keyForm: {
    ...KEYFORM,
  },
  keyPages: {
    ...KEYPAGES,
  },
  // ISSUERS
  importIssuerLink: '[data-test-generate-issuer="import"]',
  generateIssuerDropdown: '[data-test-issuer-generate-dropdown]',
  generateIssuerRoot: '[data-test-generate-issuer="root"]',
  generateIssuerIntermediate: '[data-test-generate-issuer="intermediate"]',
  issuerPopupMenu: '[data-test-popup-menu-trigger]',
  issuerPopupDetails: '[data-test-popup-menu-details] a',
  issuerDetails: {
    title: '[data-test-pki-issuer-page-title]',
    ...ISSUERDETAILS,
  },
  // CONFIGURATION
  configuration: {
    title: '[data-test-pki-configuration-page-title]',
    emptyState: '[data-test-configuration-empty-state]',
    pkiBetaBanner: '[data-test-pki-configuration-banner]',
    pkiBetaBannerLink: '[data-test-pki-configuration-banner] a',
    ...CONFIGURATION,
    ...DELETE,
  },
};
