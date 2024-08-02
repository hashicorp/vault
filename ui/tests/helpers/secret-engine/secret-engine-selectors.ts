/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const SECRET_ENGINE_SELECTORS = {
  backButton: '[data-test-back-button]',
  configTab: '[data-test-configuration-tab]',
  configure: '[data-test-secret-backend-configure]',
  configureTitle: (type: string) => `[data-test-backend-configure-title="${type}"]`,
  configurationToggle: '[data-test-mount-config-toggle]',
  createSecret: '[data-test-secret-create]',
  crumb: (path: string) => `[data-test-secret-breadcrumb="${path}"] a`,
  generateLink: '[data-test-backend-credentials]',
  mountType: (name: string) => `[data-test-mount-type="${name}"]`,
  mountSubmit: '[data-test-mount-submit]',
  secretHeader: '[data-test-secret-header]',
  secretLink: (name: string) => (name ? `[data-test-secret-link="${name}"]` : '[data-test-secret-link]'),
  viewBackend: '[data-test-backend-view-link]',
  warning: '[data-test-warning]',
  aws: {
    rootForm: '[data-test-aws-root-creds-form]',
    delete: (role: string) => `[data-test-aws-role-delete="${role}"]`,
  },
  ssh: {
    configureForm: '[data-test-ssh-configure-form]',
    sshInput: (name: string) => `[data-test-ssh-input="${name}"]`,
  },
};
