/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { PKI_ROLE_FORM } from '../components/pki/pki-role-form';
import { PKI_ROLE_GENERATE } from '../components/pki/pki-role-generate';
import { PKI_KEY_FORM } from '../components/pki/pki-key-form';
import { PKI_KEYS } from '../components/pki/page/pki-keys';
import { PKI_ISSUER_DETAILS } from '../components/pki/pki-issuer-details';
import { PKI_CONFIGURE_CREATE } from '../components/pki/pki-configure-create';
import { PKI_DELETE_ALL_ISSUERS } from '../components/pki/pki-delete-all-issuers';
import { PKI_TIDY_FORM } from '../components/pki/page/pki-tidy-form';
import { PKI_CONFIG_EDIT } from '../components/pki/page/pki-configuration-edit';
import { PKI_GENERATE_ROOT } from '../components/pki/pki-generate-root';

export const SELECTORS = {
  breadcrumbContainer: '[data-test-breadcrumbs]',
  breadcrumbs: '[data-test-breadcrumbs] li',
  overviewBreadcrumb: '[data-test-breadcrumbs] li:nth-of-type(2) > a',
  pageTitle: '[data-test-pki-role-page-title]',
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
  tidyTab: '[data-test-secret-list-tab="Tidy"]',
  configTab: '[data-test-secret-list-tab="Configuration"]',
  // ROLES
  deleteRoleButton: '[data-test-pki-role-delete]',
  generateCertLink: '[data-test-pki-role-generate-cert]',
  signCertLink: '[data-test-pki-role-sign-cert]',
  editRoleLink: '[data-test-pki-role-edit-link]',
  createRoleLink: '[data-test-pki-role-create-link]',
  roleForm: {
    ...PKI_ROLE_FORM,
  },
  generateCertForm: {
    ...PKI_ROLE_GENERATE,
  },
  // KEYS
  keyForm: {
    ...PKI_KEY_FORM,
  },
  keyPages: {
    ...PKI_KEYS,
  },
  // ISSUERS
  issuerListItem: (id) => `[data-test-issuer-list="${id}"]`,
  importIssuerLink: '[data-test-generate-issuer="import"]',
  generateIssuerDropdown: '[data-test-issuer-generate-dropdown]',
  generateIssuerRoot: '[data-test-generate-issuer="root"]',
  generateIssuerIntermediate: '[data-test-generate-issuer="intermediate"]',
  issuerPopupMenu: '[data-test-popup-menu-trigger]',
  issuerPopupDetails: '[data-test-popup-menu-details]',
  issuerDetails: {
    title: '[data-test-pki-issuer-page-title]',
    ...PKI_ISSUER_DETAILS,
  },
  // CONFIGURATION
  configuration: {
    title: '[data-test-pki-configuration-page-title]',
    emptyState: '[data-test-configuration-empty-state]',
    nextStepsBanner: '[data-test-config-next-steps]',
    importError: '[data-test-message-error]',
    pkiBetaBanner: '[data-test-pki-configuration-banner]',
    pkiBetaBannerLink: '[data-test-pki-configuration-banner] a',
    ...PKI_CONFIGURE_CREATE,
    ...PKI_DELETE_ALL_ISSUERS,
    ...PKI_TIDY_FORM,
    ...PKI_GENERATE_ROOT,
  },
  // EDIT CONFIGURATION
  configEdit: {
    ...PKI_CONFIG_EDIT,
  },
};
