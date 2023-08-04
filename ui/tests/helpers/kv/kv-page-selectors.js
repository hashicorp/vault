/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export const PAGE = {
  list: {
    createSecret: '[data-test-toolbar-create-secret]',
    item: (secret) => `[data-test-list-item="${secret}"]`,
  },
  details: {
    versionDropdown: '[data-test-version-dropdown]',
    createNewVersion: '[data-test-create-new-version]',
    version: (number) => `[data-test-version="${number}"]`,
    editMetadataBtn: '[data-test-edit-metadata]',
  },
  form: {
    kvRow: '[data-test-kv-row]',
    inputByAttr: (attr) => `[data-test-input="${attr}"]`,
    automateSecretDeletion: '[data-test-ttl-value="Automate secret deletion"]',
    metadataCancel: '[data-test-kv-metadata-cancel]',
    metadataUpdate: '[data-test-kv-metadata-update]',
    secretSave: '[data-test-kv-secret-save]',
    secretCancel: '[data-test-kv-secret-cancel]',
    versionAlert: '[data-test-secret-version-alert]',
    // <KvObjectEditor>
    keyInput: (idx = 0) => `[data-test-kv-key="${idx}"]`,
    valueInput: (idx = 0) => `[data-test-kv-value="${idx}"]`,
    maskedValueInput: (idx = 0) => `[data-test-kv-value="${idx}"] [data-test-textarea]`,
    deleteRow: (idx = 0) => `[data-test-kv-delete-row="${idx}"]`,
    dataInputLabel: ({ isJson = false }) =>
      isJson ? '[data-test-component="json-editor-title"]' : '[data-test-kv-label]',
  },
  versions: {
    popup: '[data-test-popup-menu-trigger]',
    icon: (version) => `[data-test-icon-holder="${version}"]`,
    linkedBlock: (version) => `[data-test-version-linked-block="${version}"]`,
  },
};
