/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export const PAGE = {
  // General selectors that are common between pages
  title: '[data-test-header-title]',
  breadcrumbs: '[data-test-breadcrumbs]',
  breadcrumb: '[data-test-breadcrumbs] li',
  breadcrumbAtIdx: (idx) => `[data-test-crumb="${idx}"] a`,
  infoRowValue: (label) => `[data-test-value-div="${label}"]`,
  infoRowToggleMasked: (label) => `[data-test-value-div="${label}"] [data-test-button="toggle-masked"]`,
  secretTab: (tab) => `[data-test-secrets-tab="${tab}"]`,
  emptyStateTitle: '[data-test-empty-state-title]',
  emptyStateMessage: '[data-test-empty-state-message]',
  emptyStateActions: '[data-test-empty-state-actions]',
  toolbar: 'nav.toolbar',
  toolbarAction: 'nav.toolbar-actions .toolbar-link',
  secretRow: '[data-test-component="info-table-row"]',
  // specific page selectors
  metadata: {
    editBtn: '[data-test-edit-metadata]',
    addCustomMetadataBtn: '[data-test-add-custom-metadata]',
    customMetadataSection: '[data-test-kv-custom-metadata-section]',
    secretMetadataSection: '[data-test-kv-metadata-section]',
  },
  detail: {
    versionTooltip: '[data-test-kv-version-tooltip-trigger]',
    versionDropdown: '[data-test-version-dropdown]',
    version: (number) => `[data-test-version="${number}"]`,
    createNewVersion: '[data-test-create-new-version]',
  },
  list: {
    createSecret: '[data-test-toolbar-create-secret]',
    item: (secret) => `[data-test-list-item="${secret}"]`,
    filter: `[data-test-component="kv-list-filter"]`,
    overviewCard: '[data-test-overview-card-container="View secret"]',
    overviewInput: '[data-test-search-roles] input',
    overviewButton: '[data-test-get-secret-detail]',
  },
  versions: {
    popup: '[data-test-popup-menu-trigger]',
    icon: (version) => `[data-test-icon-holder="${version}"]`,
    linkedBlock: (version) => `[data-test-version-linked-block="${version}"]`,
    button: (version) => `[data-test-version-button="${version}"]`,
  },
  create: {
    metadataSection: '[data-test-metadata-section]',
  },
};

// Form/Interactive selectors that are common between pages and forms
export const FORM = {
  inputByAttr: (attr) => `[data-test-input="${attr}"]`,
  fieldByAttr: (attr) => `[data=test=field="${attr}"]`, // formfield
  toggleJson: '[data-test-toggle-input="json"]',
  toggleMasked: '[data-test-button="toggle-masked"]',
  toggleMetadata: '[data-test-metadata-toggle]',
  jsonEditor: '[data-test-component="code-mirror-modifier"]',
  ttlValue: (name) => `[data-test-ttl-value="${name}"]`,
  toggleByLabel: (label) => `[data-test-ttl-toggle="${label}"]`,
  dataInputLabel: ({ isJson = false }) =>
    isJson ? '[data-test-component="json-editor-title"]' : '[data-test-kv-label]',
  // <KvObjectEditor>
  kvLabel: '[data-test-kv-label]',
  kvRow: '[data-test-kv-row]',
  keyInput: (idx = 0) => `[data-test-kv-key="${idx}"]`,
  valueInput: (idx = 0) => `[data-test-kv-value="${idx}"]`,
  maskedValueInput: (idx = 0) => `[data-test-kv-value="${idx}"] [data-test-textarea]`,
  addRow: (idx = 0) => `[data-test-kv-add-row="${idx}"]`,
  deleteRow: (idx = 0) => `[data-test-kv-delete-row="${idx}"]`,
  // Alerts & validation
  inlineAlert: '[data-test-inline-alert]',
  validation: (attr) => `[data-test-field="${attr}"] [data-test-inline-alert]`,
  messageError: '[data-test-message-error]',
  validationWarning: '[data-test-validation-warning]',
  invalidFormAlert: '[data-test-invalid-form-alert]',
  versionAlert: '[data-test-secret-version-alert]',
  noReadAlert: '[data-test-warning-no-read-permissions]',
  // Form btns
  saveBtn: '[data-test-kv-save]',
  cancelBtn: '[data-test-kv-cancel]',
};

export const parseJsonEditor = (find) => {
  return JSON.parse(find(FORM.jsonEditor).innerText);
};
