/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const KV_WORKFLOW = {
  // General selectors that are common between pages
  breadcrumb: '[data-test-breadcrumbs] li',
  infoRow: '[data-test-component="info-table-row"]',
  infoRowToggleMasked: (label) => `[data-test-value-div="${label}"] [data-test-button="toggle-masked"]`,
  error: {
    title: '[data-test-page-error] h1',
    message: '[data-test-page-error] p',
  },
  toolbar: 'nav.toolbar',
  toolbarAction: 'nav.toolbar-actions .toolbar-link, nav.toolbar-actions .toolbar-button',
  secretRow: '[data-test-component="info-table-row"]', // replace with infoRow
  // specific page selectors
  backends: {
    link: (backend) => `[data-test-secrets-backend-link="${backend}"]`,
  },
  edit: {
    toggleDiff: '[data-test-toggle-input="Show diff"',
    toggleDiffDescription: '[data-test-diff-description]',
  },
  list: {
    createSecret: '[data-test-toolbar-create-secret]',
    item: (secret) => (!secret ? '[data-test-list-item]' : `[data-test-list-item="${secret}"]`),
    filter: `[data-test-kv-list-filter]`,
    listMenuDelete: `[data-test-popup-metadata-delete]`,
    listMenuCreate: `[data-test-popup-create-new-version]`,
    overviewCard: '[data-test-overview-card-container="View secret"]',
    overviewInput: '[data-test-view-secret] input',
    overviewButton: '[data-test-get-secret-detail]',
    pagination: '[data-test-pagination]',
    paginationInfo: '.hds-pagination-info',
    paginationNext: '.hds-pagination-nav__arrow--direction-next',
    paginationSelected: '.hds-pagination-nav__number--is-selected',
  },
  versions: {
    icon: (version) => `[data-test-icon-holder="${version}"]`,
    linkedBlock: (version) =>
      version ? `[data-test-version-linked-block="${version}"]` : '[data-test-version-linked-block]',
    versionMenu: (version) => `[data-test-version-linked-block="${version}"] [data-test-popup-menu-trigger]`,
    createFromVersion: (version) => `[data-test-create-new-version-from="${version}"]`,
  },
  diff: {
    visualDiff: '[data-test-visual-diff]',
    added: `.jsondiffpatch-added`,
    deleted: `.jsondiffpatch-deleted`,
  },
  create: {
    metadataSection: '[data-test-metadata-section]',
  },
  paths: {
    codeSnippet: (section) => `[data-test-commands="${section}"] code`,
    snippetCopy: (section) => `[data-test-commands="${section}"] button`,
  },
};

// Form/Interactive selectors that are common between pages and forms
export const KV_FORM = {
  inputByAttr: (attr) => `[data-test-input="${attr}"]`,
  fieldByAttr: (attr) => `[data=test=field="${attr}"]`, // formfield
  toggleJson: '[data-test-toggle-input="json"]',
  toggleJsonValues: '[data-test-toggle-input="revealValues"]',
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
  noReadAlert: '[data-test-secret-no-read-alert]',
  // Form btns
  saveBtn: '[data-test-save]',
  cancelBtn: '[data-test-cancel]',
};

export const parseJsonEditor = (find) => {
  return JSON.parse(find(KV_FORM.jsonEditor).innerText);
};
