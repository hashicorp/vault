/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const SECRET_ENGINE_SELECTORS = {
  configTab: '[data-test-configuration-tab]',
  configure: '[data-test-secret-backend-configure]',
  configureNote: (name: string) => `[data-test-configure-note="${name}"]`,
  configureTitle: (type: string) => `[data-test-backend-configure-title="${type}"]`,
  configurationToggle: '[data-test-mount-config-toggle]',
  createSecret: '[data-test-secret-create]',
  crumb: (path: string) => `[data-test-secret-breadcrumb="${path}"] a`,
  error: {
    title: '[data-test-backend-error-title]',
  },
  generateLink: '[data-test-backend-credentials]',
  secretHeader: '[data-test-secret-header]',
  secretLink: (name: string) => (name ? `[data-test-secret-link="${name}"]` : '[data-test-secret-link]'),
  secretLinkMenu: (name: string) => `[data-test-secret-link="${name}"] [data-test-popup-menu-trigger]`,
  secretLinkMenuDelete: (name: string) =>
    `[data-test-secret-link="${name}"] [data-test-confirm-action-trigger]`,
  secretLinkATag: (name: string) =>
    name ? `[data-test-secret-item-link="${name}"]` : '[data-test-secret-item-link]',
  viewBackend: '[data-test-backend-view-link]',
  warning: '[data-test-warning]',
  configureForm: '[data-test-configure-form]',
  additionalConfigModelTitle: '[data-test-additional-config-model-title]',
  wif: {
    accessTypeSection: '[data-test-access-type-section]',
    accessTitle: '[data-test-access-title]',
    accessType: (type: string) => `[data-test-access-type="${type}"]`,
    accessTypeSubtext: '[data-test-access-type-subtext]',
    issuerWarningCancel: '[data-test-issuer-cancel]',
    issuerWarningMessage: '[data-test-issuer-warning-message]',
    issuerWarningModal: '[data-test-issuer-warning]',
    issuerWarningSave: '[data-test-issuer-save]',
  },
  aws: {
    deleteRole: (role: string) => `[data-test-aws-role-delete="${role}"]`,
  },
  ssh: {
    editConfigSection: '[data-test-edit-config-section]',
    save: '[data-test-configure-save-button]',
    cancel: '[data-test-cancel-button]',
    delete: '[data-test-delete-public-key]',
    createRole: '[data-test-role-ssh-create]',
    deleteRole: '[data-test-ssh-role-delete]',
  },
};
