/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export const SELECTORS = {
  details: {
    versionDropdown: '[data-test-version-dropdown]',
    version: (number) => `[data-test-version="${number}"]`,
    editMetadataBtn: '[data-test-edit-metadata]',
  },
  form: {
    kvRow: '[data-test-kv-row]',
    inputByAttr: (attr) => `[data-test-input="${attr}"]`,
    automateSecretDeletion: '[data-test-ttl-value="Automate secret deletion"]',
    metadataCancel: '[data-test-kv-metadata-cancel]',
    metadataUpdate: '[data-test-kv-metadata-update]',
    keyInput: '[data-test-kv-key]',
    valueInput: '[data-test-kv-value] textarea',
    secretSave: '[data-test-kv-secret-save]',
    secretCancel: '[data-test-kv-secret-cancel]',
  },
};
