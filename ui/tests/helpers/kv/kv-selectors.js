/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

export const PAGE = {
  // General selectors that are common between pages
  title: '.hds-page-header__title',
  breadcrumbs: '[data-test-breadcrumbs]',
  breadcrumb: '[data-test-breadcrumbs] li',
  breadcrumbAtIdx: (idx) => `[data-test-breadcrumbs] li:nth-child(${idx + 1}) a`,
  breadcrumbCurrentAtIdx: (idx) =>
    `[data-test-breadcrumbs] li:nth-child(${idx + 1}) .hds-breadcrumb__current`,
  infoRow: '[data-test-component="info-table-row"]',
  infoRowValue: (label) => `[data-test-row-value="${label}"]`, // TODO replace with GENERAL.infoRowValue
  infoRowToggleMasked: (label) => `[data-test-row-value="${label}"] [data-test-button="toggle-masked"]`,
  secretTab: (tab) => (tab ? `[data-test-secrets-tab="${tab}"]` : '[data-test-secrets-tab]'),
  emptyStateTitle: '[data-test-empty-state-title]',
  emptyStateMessage: '[data-test-empty-state-message]',
  emptyStateActions: '[data-test-empty-state-actions]',
  popup: '[data-test-popup-menu-trigger]',
  toolbar: 'nav.toolbar',
  toolbarAction: 'nav.toolbar-actions .toolbar-link, nav.toolbar-actions .toolbar-button',
  secretRow: '[data-test-component="info-table-row"]', // replace with infoRow
  // specific page selectors
  metadata: {
    requestData: '[data-test-request-data]',
    editBtn: '[data-test-edit-metadata]',
    addCustomMetadataBtn: '[data-test-add-custom-metadata]',
    customMetadataSection: '[data-test-kv-custom-metadata-section]',
    secretMetadataSection: '[data-test-kv-metadata-section]',
    deleteMetadata: '[data-test-kv-delete="delete-metadata"]',
  },
  detail: {
    versionTimestamp: '[data-test-tooltip="kv-version"]',
    versionDropdown: '[data-test-version-dropdown]',
    version: (number) => `[data-test-version="${number}"]`,
    createNewVersion: '[data-test-create-new-version]',
    patchLatest: '[data-test-patch-latest-version]',
    delete: '[data-test-kv-delete="delete"]',
    destroy: '[data-test-kv-delete="destroy"]',
    undelete: '[data-test-kv-delete="undelete"]',
    copy: '[data-test-copy-menu-trigger]',
    deleteModal: '[data-test-delete-modal]',
    deleteModalTitle: '[data-test-delete-modal] [data-test-modal-title]',
    deleteOption: 'input#delete-version',
    deleteOptionLatest: 'input#delete-latest-version',
    deleteConfirm: '[data-test-delete-modal-confirm]',
    syncAlert: (name) => (name ? `[data-test-sync-alert="${name}"]` : '[data-test-sync-alert]'),
  },
  edit: {
    toggleDiffDescription: '[data-test-diff-description]',
  },
  list: {
    createSecret: '[data-test-button="create secret"]',
    item: (secret) => (!secret ? '[data-test-list-item]' : `[data-test-list-item="${secret}"]`),
    menuItem: (label) => `[data-test-list-menu-item="${label}"]`,
    filter: `[data-test-kv-list-filter]`,
    overviewCard: '[data-test-overview-card-container="View secret"]',
    overviewInput: '[data-test-view-secret] input',
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
    codeSnippet: (section) => `[data-test-code-block="${section}"] code`,
    snippetCopy: (section) => `[data-test-code-block="${section}"] button`,
  },
};

// Form/Interactive selectors that are common between pages and forms
export const FORM = {
  inputByAttr: (attr) => `[data-test-input="${attr}"]`,
  fieldByAttr: (attr) => `[data=test=field="${attr}"]`, // formfield
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
  maskedValueInput: (idx = 0) => `[data-test-kv-value="${idx}"] [data-test-input]`,
  addRow: (idx = 0) => `[data-test-kv-add-row="${idx}"]`,
  deleteRow: (idx = 0) => `[data-test-kv-delete-row="${idx}"]`,
  // <KvPatchEditor>
  patchEditorForm: '[data-test-kv-patch-editor]',
  patchEdit: (idx = 0) => `[data-test-edit-button="${idx}"]`,
  patchDelete: (idx = 0) => `[data-test-delete-button="${idx}"]`,
  patchUndo: (idx = 0) => `[data-test-undo-button="${idx}"]`,
  patchAdd: '[data-test-add-button]',
  patchAlert: (type, idx) => `[data-test-alert-${type}="${idx}"]`,
  // Alerts & validation
  inlineAlert: '[data-test-inline-alert]',
  validation: (attr) => `[data-test-field="${attr}"] [data-test-inline-alert]`,
  messageError: '[data-test-message-error]',
  validationError: (attr) => `[data-test-validation-error="${attr}"]`,
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

export const parseObject = (cm) => JSON.parse(cm().getValue());
