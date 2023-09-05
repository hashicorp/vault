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
  infoRow: '[data-test-component="info-table-row"]',
  infoRowValue: (label) => `[data-test-value-div="${label}"]`,
  infoRowToggleMasked: (label) => `[data-test-value-div="${label}"] [data-test-button="toggle-masked"]`,
  secretTab: (tab) => (tab ? `[data-test-secrets-tab="${tab}"]` : '[data-test-secrets-tab]'),
  emptyStateTitle: '[data-test-empty-state-title]',
  emptyStateMessage: '[data-test-empty-state-message]',
  emptyStateActions: '[data-test-empty-state-actions]',
  popup: '[data-test-popup-menu-trigger]',
  error: {
    title: '[data-test-page-error] h1',
    message: '[data-test-page-error] p',
  },
  toolbar: 'nav.toolbar',
  toolbarAction: 'nav.toolbar-actions .toolbar-link',
  secretRow: '[data-test-component="info-table-row"]', // replace with infoRow
  // specific page selectors
  backends: {
    link: (backend) => `[data-test-secrets-backend-link="${backend}"]`,
  },
  metadata: {
    editBtn: '[data-test-edit-metadata]',
    addCustomMetadataBtn: '[data-test-add-custom-metadata]',
    customMetadataSection: '[data-test-kv-custom-metadata-section]',
    secretMetadataSection: '[data-test-kv-metadata-section]',
    deleteMetadata: '[data-test-kv-delete="delete-metadata"]',
  },
  detail: {
    versionTimestamp: '[data-test-kv-version-tooltip-trigger]',
    versionDropdown: '[data-test-version-dropdown]',
    version: (number) => `[data-test-version="${number}"]`,
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
  },
  edit: {
    toggleDiff: '[data-test-toggle-input="Show diff"',
    toggleDiffDescription: '[data-test-diff-description]',
    visualDiff: '[data-test-visual-diff]',
    added: `.jsondiffpatch-added`,
    deleted: `.jsondiffpatch-deleted`,
  },
  list: {
    createSecret: '[data-test-toolbar-create-secret]',
    item: (secret) => (!secret ? '[data-test-list-item]' : `[data-test-list-item="${secret}"]`),
    filter: `[data-test-kv-list-filter]`,
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
    button: (version) => `[data-test-version-button="${version}"]`,
    versionMenu: (version) => `[data-test-version-linked-block="${version}"] [data-test-popup-menu-trigger]`,
    createFromVersion: (version) => `[data-test-create-new-version-from="${version}"]`,
  },
  create: {
    metadataSection: '[data-test-metadata-section]',
  },
  paths: {
    copyButton: (label) => `${PAGE.infoRowValue(label)} button`,
    codeSnippet: (section) => `[data-test-code-snippet][data-test-commands="${section}"] code`,
    snippetCopy: (section) => `[data-test-code-snippet][data-test-commands="${section}"] button`,
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
  noReadAlert: '[data-test-secret-no-read-alert]',
  // Form btns
  saveBtn: '[data-test-kv-save]',
  cancelBtn: '[data-test-kv-cancel]',
};

export const parseJsonEditor = (find) => {
  return JSON.parse(find(FORM.jsonEditor).innerText);
};
