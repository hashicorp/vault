/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { findAll } from '@ember/test-helpers';

export const GENERAL = {
  /* ────── Header / Breadcrumbs ────── */
  breadcrumb: '[data-test-breadcrumbs] li',
  breadcrumbAtIdx: (idx: string) => `[data-test-breadcrumbs] li:nth-child(${idx + 1}) a`,
  breadcrumbLink: (label: string) => `[data-test-breadcrumb="${label}"] a`,
  breadcrumbs: '[data-test-breadcrumbs]',
  headerContainer: 'header.page-header',
  title: '[data-test-page-title]',

  /* ────── Tabs & Navigation ────── */
  tab: (name: string) => `[data-test-tab="${name}"]`,
  hdsTab: (name: string) => `[data-test-tab="${name}"] button`, // HDS tab buttons
  secretTab: (name: string) => `[data-test-secret-list-tab="${name}"]`,
  navLink: (label: string) => `[data-test-sidebar-nav-link="${label}"]`,
  linkTo: (label: string) => `[data-test-link-to="${label}"]`,

  /* ────── Buttons ────── */
  backButton: '[data-test-back-button]',
  cancelButton: '[data-test-cancel]',
  confirmButton: '[data-test-confirm-button]', // used most often on modal or confirm popups
  confirmTrigger: '[data-test-confirm-action-trigger]',
  copyButton: '[data-test-copy-button]',
  // there should only be one save button per view (e.g. one per form) so this does not need to be dynamic
  // this button should be used for any kind of "submit" on a form or "save" action.
  submitButton: '[data-test-submit]',
  button: (label: string) => (label ? `[data-test-button="${label}"]` : '[data-test-button]'),

  /* ────── Menus & Lists ────── */
  menuTrigger: '[data-test-popup-menu-trigger]',
  menuItem: (name: string) => `[data-test-popup-menu="${name}"]`,
  listItem: (label: string) => `[data-test-list-item="${label}"]`,
  listItemLink: '[data-test-list-item-link]',
  linkedBlock: (item: string) => `[data-test-linked-block="${item}"]`,

  /* ────── Inputs / Form Fields ────── */
  checkboxByAttr: (attr: string) => `[data-test-checkbox="${attr}"]`,
  confirmModalInput: '[data-test-confirmation-modal-input]',
  confirmMessage: '[data-test-confirm-action-message]',
  docLinkByAttr: (attr: string) => `[data-test-doc-link="${attr}"]`,
  enableField: (attr: string) => `[data-test-enable-field="${attr}"] button`,
  fieldByAttr: (attr: string) => `[data-test-field="${attr}"]`,
  fieldLabel: () => `[data-test-form-field-label]`,
  fieldLabelbyAttr: (attr: string) => `[data-test-form-field-label="${attr}"]`,
  groupControlByIndex: (index: number) => `.hds-form-group__control-field:nth-of-type(${index})`,
  helpText: () => `[data-test-help-text]`,
  helpTextByAttr: (attr: string) => `[data-test-help-text="${attr}"]`,
  helpTextByGroupControlIndex: (index: number) =>
    `.hds-form-group__control-field:nth-of-type(${index}) [data-test-help-text]`,
  inputByAttr: (attr: string) => `[data-test-input="${attr}"]`,
  inputGroupByAttr: (attr: string) => `[data-test-input-group="${attr}"]`,
  inputSearch: (attr: string) => `[data-test-input-search="${attr}"]`,
  filterInput: '[data-test-filter-input]',
  filterInputExplicit: '[data-test-filter-input-explicit]',
  labelById: (id: string) => `label[id="${id}"]`,
  labelByGroupControlIndex: (index: number) => `.hds-form-group__control-field:nth-of-type(${index}) label`,
  radioByAttr: (attr: string) => `[data-test-radio="${attr}"]`,
  selectByAttr: (attr: string) => `[data-test-select="${attr}"]`,
  toggleInput: (attr: string) => `[data-test-toggle-input="${attr}"]`,
  textToggle: '[data-test-text-toggle]',
  textToggleTextarea: '[data-test-text-file-textarea]',
  filter: (name: string) => `[data-test-filter="${name}"]`,

  /* ────── Code Blocks / Editor ────── */
  codemirror: `[data-test-component="code-mirror-modifier"]`,
  codemirrorTextarea: `[data-test-component="code-mirror-modifier"] textarea`,
  codeBlock: (label: string) => `[data-test-code-block="${label}"]`,

  /* ────── Key/Value Editors ────── */
  kvObjectEditor: {
    key: (idx = 0) => `[data-test-kv-key="${idx}"]`,
    value: (idx = 0) => `[data-test-kv-value="${idx}"]`,
    addRow: '[data-test-kv-add-row]',
    deleteRow: (idx = 0) => `[data-test-kv-delete-row="${idx}"]`,
  },
  kvSuggestion: {
    input: '[data-test-kv-suggestion-input]',
    select: '[data-test-kv-suggestion-select]',
  },

  /* ────── Search Select ────── */
  searchSelect: {
    trigger: (id: string) => `[data-test-component="search-select"]#${id} .ember-basic-dropdown-trigger`,
    options: '.ember-power-select-option',
    option: (index = 0) => `.ember-power-select-option:nth-child(${index + 1})`,
    optionIndex: (text: string) =>
      findAll('.ember-power-select-options li').findIndex((e) => e.textContent?.trim() === text),
    selectedOption: (index = 0) => `[data-test-selected-option="${index}"]`,
    noMatch: '.ember-power-select-option--no-matches-message',
    removeSelected: '[data-test-selected-list-button="delete"]',
    searchInput: '.ember-power-select-search-input',
  },

  /* ────── TTL Fields ────── */
  ttl: {
    toggle: (attr: string) => `[data-test-toggle-label="${attr}"]`,
    input: (attr: string) => `[data-test-ttl-value="${attr}"]`,
  },

  /* ────── Info Table / Rows ────── */
  infoRowLabel: (label: string) => `[data-test-row-label="${label}"]`,
  infoRowValue: (label: string) => `[data-test-row-value="${label}"]`,

  /* ────── Empty / Error / Alert States ────── */
  emptyStateTitle: '[data-test-empty-state-title]',
  emptyStateSubtitle: '[data-test-empty-state-subtitle]',
  emptyStateMessage: '[data-test-empty-state-message]',
  emptyStateActions: '[data-test-empty-state-actions]',
  flashMessage: '[data-test-flash-message]',
  latestFlashContent: '[data-test-flash-message]:last-of-type [data-test-flash-message-body]',
  inlineAlert: '[data-test-inline-alert]',
  inlineError: '[data-test-inline-error-message]',
  messageError: '[data-test-message-error]',
  notFound: '[data-test-not-found]',
  validationErrorByAttr: (attr: string) => `[data-test-validation-error=${attr}]`,
  validationWarningByAttr: (attr: string) => `[data-test-validation-warning=${attr}]`,

  pageError: {
    error: '[data-test-page-error]',
    errorTitle: (httpStatus: number) => `[data-test-page-error-title="${httpStatus}"]`,
    errorSubtitle: '[data-test-page-error-subtitle]',
    errorMessage: '[data-test-page-error-message]',
    errorDetails: '[data-test-page-error-details]',
  },

  /* ────── Pagination ────── */
  pagination: {
    next: '.hds-pagination-nav__arrow--direction-next',
    prev: '.hds-pagination-nav__arrow--direction-prev',
  },

  /* ────── Overview Cards ────── */
  overviewCard: {
    container: (title: string) => `[data-test-overview-card-container="${title}"]`,
    title: (title: string) => `[data-test-overview-card-title="${title}"]`,
    description: (title: string) => `[data-test-overview-card-subtitle="${title}"]`,
    content: (title: string) => `[data-test-overview-card-content="${title}"]`,
    actionText: (text: string) => `[data-test-action-text="${text}"]`,
    actionLink: (label: string) => `[data-test-overview-card="${label}"] a`,
  },

  /* ────── Misc ────── */
  icon: (name: string) => `[data-test-icon="${name}"]`,
};
