/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const KV_SECRET = {
  versionTimestamp: '[data-test-kv-version-tooltip-trigger]',
  versionDropdown: '[data-test-version-dropdown]',
  version: (number: number) => `[data-test-version="${number}"]`,
  createNewVersion: '[data-test-create-new-version]',
  delete: '[data-test-kv-delete="delete"]',
  destroy: '[data-test-kv-delete="destroy"]',
  undelete: '[data-test-kv-delete="undelete"]',
  copy: '[data-test-copy-menu-trigger]',
  deleteModal: '[data-test-delete-modal]',
  deleteModalTitle: '[data-test-delete-modal] [data-test-modal-title]',
  deleteOption: 'input#delete-version',
  deleteOptionLatest: 'input#delete-latest-version',
  deleteConfirm: '[data-test-delete-modal-confirm]',
  syncAlert: (name: string) => (name ? `[data-test-sync-alert="${name}"]` : '[data-test-sync-alert]'),
};
