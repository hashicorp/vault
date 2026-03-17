/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { findAll } from '@ember/test-helpers';

export const GENERAL = {
  /* ────── Header / Breadcrumbs ────── */
  breadcrumb: '[data-test-breadcrumbs] li',
  breadcrumbAtIdx: (idx: string) => `[data-test-breadcrumbs] li:nth-child(${idx + 1}) a`,
  breadcrumbLink: (label: string) => `[data-test-breadcrumb="${label}"] a`,
  currentBreadcrumb: (label: string) => `[data-test-breadcrumb="${label}"]`,
  breadcrumbs: '[data-test-breadcrumbs]',
  headerContainer: 'header.page-header',
  title: '[data-test-page-title]',
  hdsPageHeaderTitle: '.hds-page-header__title',
  hdsPageHeaderSubtitle: '.hds-page-header__subtitle',
  hdsPageHeaderDescription: '.hds-page-header__description',

  /* ────── Tabs & Navigation ────── */
  tab: (name: string) => `[data-test-tab="${name}"]`,
  tabLink: (name: string) => `[data-test-tab="${name}"] a`,
  hdsTab: (name: string) => (name ? `[data-test-tab="${name}"] button` : '[data-test-tab] button'), // HDS tab buttons
  hdsTabPanel: (name: string) => (name ? `[data-test-panel="${name}"]` : '[data-test-panel]'),
  secretTab: (name: string) => `[data-test-secret-list-tab="${name}"]`,
  navLink: (label: string) =>
    label ? `[data-test-sidebar-nav-link="${label}"]` : '[data-test-sidebar-nav-link]',
  navHeading: (label: string) =>
    label ? `[data-test-sidebar-nav-heading="${label}"]` : '[data-test-sidebar-nav-heading]',
  linkTo: (label: string) => `[data-test-link-to="${label}"]`,

  /* ────── Buttons ────── */
  backButton: '[data-test-back-button]',
  cancelButton: '[data-test-cancel]',
  confirmButton: '[data-test-confirm-button]', // used most often on modal or confirm popups
  confirmTrigger: '[data-test-confirm-action-trigger]',
  copyButton: '[data-test-copy-button]',
  revealButton: (label: string) => `[data-test-reveal="${label}"] button`, // intended for Hds::Reveal components
  accordionButton: (label: string) => `[data-test-accordion="${label}"] button`, // intended for Hds::Accordion components
  // there should only be one submit button per view (e.g. one per form) so this does not need to be dynamic
  // this button should be used for any kind of "submit" on a form or "save" action.
  submitButton: '[data-test-submit]',
  button: (label: string) => (label ? `[data-test-button="${label}"]` : '[data-test-button]'),
  copySnippet: (name: string) => `[data-test-copy-snippet=${name}]`,

  /* ────── Menus & Lists ────── */
  dropdownToggle: (text: string) => `[data-test-dropdown="${text}"]`, // Use when dropdown toggle has text
  menuTrigger: '[data-test-popup-menu-trigger]', // Use when dropdown toggle is just an icon
  menuItem: (name: string) => `[data-test-popup-menu="${name}"]`,
  listItem: (label: string) => (label ? `[data-test-list-item="${label}"]` : '[data-test-list-item]'),
  listItemLink: '[data-test-list-item-link]',
  linkedBlock: (item: string) => `[data-test-linked-block="${item}"]`,

  /* ────── Tables ────── */
  table: (title: string) => `[data-test-table="${title}"]`,
  tableRow: (idx?: number) => (idx ? `[data-test-table-row="${idx}"]` : '[data-test-table-row]'),
  tableData: (idx?: number, key?: string) => `[data-test-table-row="${idx}"] [data-test-table-data="${key}"]`,
  tableColumnHeader: (col: number, { isAdvanced = false } = {}) =>
    `${isAdvanced ? '.hds-advanced-table__th' : 'hds-table__th'}:nth-child(${col})`, // number is not 0-indexed, first column header is 1
  tableColumnHeaderSortButton: (col: number, { isAdvanced = false } = {}) =>
    `${
      isAdvanced ? '.hds-advanced-table__th' : 'hds-table__th'
    }:nth-child(${col}) .hds-advanced-table__th-button--sort`, // number is not 0-indexed, first column header is 1

  /* ────── Inputs / Form Fields ────── */
  checkboxByAttr: (attr: string) => `[data-test-checkbox="${attr}"]`,
  confirmModalInput: '[data-test-confirmation-modal-input]',
  confirmMessage: '[data-test-confirm-action-message]',
  docLinkByAttr: (attr: string) => `[data-test-doc-link="${attr}"]`,
  enableField: (attr: string) => `[data-test-enable-field="${attr}"] button`,
  fieldByAttr: (attr: string) => `[data-test-field="${attr}"]`,
  fieldLabel: (attr: string) =>
    attr ? `[data-test-form-field-label="${attr}"]` : `[data-test-form-field-label]`,
  fileInput: '[data-test-file-input]',
  filter: (name: string) => `[data-test-filter="${name}"]`,
  filterInput: '[data-test-filter-input]',
  filterInputExplicit: '[data-test-filter-input-explicit]',
  groupControlByIndex: (index: number) => `.hds-form-group__control-field:nth-of-type(${index})`,
  helpText: '[data-test-help-text]',
  helpTextByAttr: (attr: string) => `[data-test-help-text="${attr}"]`,
  helpTextByGroupControlIndex: (index: number) =>
    `.hds-form-group__control-field:nth-of-type(${index}) [data-test-help-text]`,
  inputByAttr: (attr: string) => `[data-test-input="${attr}"]`,
  inputGroupByAttr: (attr: string) => `[data-test-input-group="${attr}"]`,
  inputSearch: (attr: string) => `[data-test-input-search="${attr}"]`,
  labelById: (id: string) => `label[id="${id}"]`,
  labelByGroupControlIndex: (index: number) => `.hds-form-group__control-field:nth-of-type(${index}) label`,
  maskedInput: '[data-test-masked-input]',
  radioByAttr: (attr: string) => (attr ? `[data-test-radio="${attr}"]` : '[data-test-radio]'),
  radioCardByAttr: (attr: string) => (attr ? `[data-test-radio-card="${attr}"]` : '[data-test-radio-card]'),
  selectByAttr: (attr: string) => `[data-test-select="${attr}"]`,
  stringListByIdx: (index: number) => `[data-test-string-list-input="${index}"]`,
  textToggle: '[data-test-text-toggle]',
  textareaByAttr: (attr: string) => `textarea[name="${attr}"]`,
  toggleInput: (attr: string) => `[data-test-toggle-input="${attr}"]`,
  superSelect: (name: string) => `[data-test-super-select="${name}"]`,

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
  inlineAlertByAttr: (attr: string) => `[data-test-inline-alert="${attr}"]`,
  inlineError: '[data-test-inline-error-message]',
  messageError: '[data-test-message-error]',
  messageDescription: '[data-test-message-error-description]',
  validationErrorByAttr: (attr: string) => `[data-test-validation-error=${attr}]`,
  validationWarningByAttr: (attr: string) => `[data-test-validation-warning=${attr}]`,

  pageError: {
    error: '[data-test-page-error]',
    title: (httpStatus?: number) =>
      httpStatus ? `[data-test-page-error-title="${httpStatus}"]` : '[data-test-page-error-title]',
    message: '[data-test-page-error-message]',
    details: '[data-test-page-error-details]',
  },

  /* ────── Pagination ────── */
  pagination: '[data-test-pagination]',
  paginationInfo: '.hds-pagination-info',
  paginationSizeSelector: '.hds-pagination-size-selector select',
  nextPage: '.hds-pagination-nav__arrow--direction-next',
  prevPage: '.hds-pagination-nav__arrow--direction-prev',

  /* ────── Overview Cards ────── */
  overviewCard: {
    container: (title: string) => `[data-test-overview-card-container="${title}"]`,
    title: (title: string) => `[data-test-overview-card-title="${title}"]`,
    description: (title: string) => `[data-test-overview-card-subtitle="${title}"]`,
    content: (title: string) => `[data-test-overview-card-content="${title}"]`,
    actionText: (text: string) => `[data-test-action-text="${text}"]`,
    actionLink: (label: string) => `[data-test-overview-card="${label}"] a`,
  },

  /* ────── Cards ────── */
  cardContainer: (title: string) =>
    title ? `[data-test-card-container="${title}"]` : '[data-test-card-container]',

  /* ────── Modals & Flyouts ────── */
  flyout: '[data-test-flyout]',
  modal: {
    container: (title: string) => `[data-test-modal="${title}"]`,
    header: (title: string) => `[data-test-modal-header="${title}"]`,
    body: (title: string) => `[data-test-modal-body="${title}"]`,
  },

  /* ────── Misc ────── */
  icon: (name: string) => (name ? `[data-test-icon="${name}"]` : '[data-test-icon]'),
  badge: (name: string) => (name ? `[data-test-badge="${name}"]` : '[data-test-badge]'),
  licenseBanner: (name: string) => `[data-test-license-banner="${name}"]`,
  tooltip: (label: string) => `[data-test-tooltip="${label}"]`,
  tooltipText: '.hds-tooltip-container',
  textDisplay: (attr: string) => `[data-test-text-display="${attr}"]`,
};
