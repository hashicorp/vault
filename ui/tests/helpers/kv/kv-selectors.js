/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export const PAGE = {
  // General selectors that are common between pages
  title: '[data-test-header-title]',
  breadcrumbs: '[data-test-breadcrumbs]',
  breadcrumb: '[data-test-breadcrumbs] li',
  infoRowValue: (label) => `[data-test-value-div="${label}"]`,
  secretTab: (tab) => `[data-test-secrets-tab="${tab}"]`,
  emptyStateTitle: '[data-test-empty-state-title]',
  emptyStateMessage: '[data-test-empty-state-message]',
  emptyStateActions: '[data-test-empty-state-actions]',
  // specific page selectors
  metadata: {
    editBtn: '[data-test-edit-metadata]',
    addCustomMetadataBtn: '[data-test-add-custom-metadata]',
    customMetadataSection: '[data-test-kv-custom-metadata-section]',
    secretMetadataSection: '[data-test-kv-metadata-section]',
  },
  detail: {
    versionCreated: '[data-test-kv-version-tooltip-trigger]',
    versionDropdown: '[data-test-version-dropdown]',
    version: (number) => `[data-test-version="${number}"]`,
    createNewVersion: '[data-test-create-new-version]',
  },
  list: {
    createSecret: '[data-test-toolbar-create-secret]',
    item: (secret) => `[data-test-list-item="${secret}"]`,
  },
  versions: {
    popup: '[data-test-popup-menu-trigger]',
    icon: (version) => `[data-test-icon-holder="${version}"]`,
    linkedBlock: (version) => `[data-test-version-linked-block="${version}"]`,
    button: (version) => `[data-test-version-button="${version}"]`,
  },
};

// Form/Interactive selectors that are common between pages and forms
export const FORM = {
  inputByAttr: (attr) => `[data-test-input="${attr}"]`,
  toggleJson: '[data-test-toggle-input="json"]',
  toggleMasked: '[data-test-button="toggle-masked"]',
  jsonEditor: '[data-test-component="code-mirror-modifier"]',
  ttlValue: (name) => `[data-test-ttl-value="${name}"]`,
  dataInputLabel: ({ isJson = false }) =>
    isJson ? '[data-test-component="json-editor-title"]' : '[data-test-kv-label]',
  // <KvObjectEditor>
  kvRow: '[data-test-kv-row]',
  keyInput: (idx = 0) => `[data-test-kv-key="${idx}"]`,
  valueInput: (idx = 0) => `[data-test-kv-value="${idx}"]`,
  maskedValueInput: (idx = 0) => `[data-test-kv-value="${idx}"] [data-test-textarea]`,
  deleteRow: (idx = 0) => `[data-test-kv-delete-row="${idx}"]`,
  // Alerts & validation
  inlineAlert: '[data-test-inline-alert]',
  validation: (attr) => `[data-test-field="${attr}"] [data-test-inline-alert]`,
  messageError: '[data-test-message-error]',
  versionAlert: '[data-test-secret-version-alert]',
  // Form btns
  saveBtn: '[data-test-kv-save]',
  cancelBtn: '[data-test-kv-cancel]',
};

export const parseJsonEditor = (find) => {
  return JSON.parse(find(FORM.jsonEditor).innerText);
};
