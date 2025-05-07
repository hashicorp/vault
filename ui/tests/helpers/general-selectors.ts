/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { findAll } from '@ember/test-helpers';
/* 
Selector patterns:
1. keep them generic. Think element wide selectors. If you need to be specific, make them dynamic. ex: do not do data-test-my-component-button, instead do data-test-button="specific action"
2. Alphabetize within the sections.
3. If you need component or feature specific selectors, put them in the component or feature specific selector files. ex: do not put kv specific selectors here.
**/
export const GENERAL = {
  // Breadcrumbs
  breadcrumb: '[data-test-breadcrumbs] li',
  breadcrumbAtIdx: (idx: string) => `[data-test-breadcrumbs] li:nth-child(${idx + 1}) a`,
  breadcrumbLink: (label: string) => `[data-test-breadcrumb="${label}"] a`,
  breadcrumbs: '[data-test-breadcrumbs]',
  component: (name: string) => `[data-test-component="${name}"]`,
  // Page elements
  hdsTab: (name: string) => `[data-test-tab="${name}"] button`, // hds tabs are li elements and QUnit needs a clickable element so add button to the selector
  headerContainer: 'header.page-header',
  icon: (name: string) => `[data-test-icon="${name}"]`,
  secretTab: (name: string) => `[data-test-secret-list-tab="${name}"]`,
  tab: (name: string) => `[data-test-tab="${name}"]`,
  title: '[data-test-page-title]',

  // Form inputs
  checkboxByAttr: (attr: string) => `[data-test-checkbox="${attr}"]`,
  codeBlock: (label: string) => `[data-test-code-block="${label}"]`,
  codemirror: `[data-test-component="code-mirror-modifier"]`,
  codemirrorTextarea: `[data-test-component="code-mirror-modifier"] textarea`,
  confirmButton: '[data-test-confirm-button]',
  confirmMessage: '[data-test-confirm-action-message]',
  confirmModalInput: '[data-test-confirmation-modal-input]',
  confirmTrigger: '[data-test-confirm-action-trigger]',
  enableField: (attr: string) => `[data-test-enable-field="${attr}"] button`,
  fieldByAttr: (attr: string) => `[data-test-field="${attr}"]`,
  filter: (name: string) => `[data-test-filter="${name}"]`,
  filterInput: '[data-test-filter-input]',
  filterInputExplicit: '[data-test-filter-input-explicit]',
  filterInputExplicitSearch: '[data-test-filter-input-explicit-search]',
  infoRowLabel: (label: string) => `[data-test-row-label="${label}"]`,
  infoRowValue: (label: string) => `[data-test-value-div="${label}"]`,
  inputByAttr: (attr: string) => `[data-test-input="${attr}"]`,
  inputSearch: (attr: string) => `[data-test-input-search="${attr}"]`,
  listItem: '[data-test-list-item-link]',
  kvObjectEditor: {
    deleteRow: (idx = 0) => `[data-test-kv-delete-row="${idx}"]`,
  },
  menuItem: (name: string) => `[data-test-popup-menu="${name}"]`,
  menuTrigger: '[data-test-popup-menu-trigger]',
  selectByAttr: (attr: string) => `[data-test-select="${attr}"]`,
  searchSelect: {
    trigger: (id: string) => `[data-test-component="search-select"]#${id} .ember-basic-dropdown-trigger`,
    options: '.ember-power-select-option',
    optionIndex: (text: string) =>
      findAll('.ember-power-select-options li').findIndex((e) => e.textContent?.trim() === text),
    option: (index = 0) => `.ember-power-select-option:nth-child(${index + 1})`,
    selectedOption: (index = 0) => `[data-test-selected-option="${index}"]`,
    noMatch: '.ember-power-select-option--no-matches-message',
    removeSelected: '[data-test-selected-list-button="delete"]',
    searchInput: '.ember-power-select-search-input',
  },
  textToggle: '[data-test-text-toggle]',
  textToggleTextarea: '[data-test-text-file-textarea]',
  toggleGroup: (attr: string) => `[data-test-toggle-group="${attr}"]`,
  toggleInput: (attr: string) => `[data-test-toggle-input="${attr}"]`,
  ttl: {
    toggle: (attr: string) => `[data-test-toggle-label="${attr}"]`,
    input: (attr: string) => `[data-test-ttl-value="${attr}"]`,
  },
  // Links and Buttons
  backButton: '[data-test-back-button]',
  cancelButton: '[data-test-cancel]',
  navLink: (label: string) => `[data-test-sidebar-nav-link="${label}"]`,
  saveButton: '[data-test-save]',
  testButton: (label: string) => `[data-test-button="${label}"]`,
  // Validation messages
  validation: (attr: string) => `[data-test-field-validation=${attr}]`,
  validationWarning: (attr: string) => `[data-test-validation-warning=${attr}]`,
  // Error messages
  messageError: '[data-test-message-error]',
  notFound: '[data-test-not-found]',
  pageError: {
    error: '[data-test-page-error]',
    errorDetails: '[data-test-page-error-details]',
    errorMessage: '[data-test-page-error-message]',
    errorTitle: (httpStatus: number) => `[data-test-page-error-title="${httpStatus}"]`,
  },
  inlineError: '[data-test-inline-error-message]',

  // Flash messages
  flashMessage: '[data-test-flash-message]',
  latestFlashContent: '[data-test-flash-message]:last-of-type [data-test-flash-message-body]',
  inlineAlert: '[data-test-inline-alert]',
  // Empty states
  emptyStateTitle: '[data-test-empty-state-title]',
  emptyStateSubtitle: '[data-test-empty-state-subtitle]',
  emptyStateMessage: '[data-test-empty-state-message]',
  emptyStateActions: '[data-test-empty-state-actions]',
  // TODO probably shouldn't be shiloh to the overview component
  overviewCard: {
    container: (title: string) => `[data-test-overview-card-container="${title}"]`,
    title: (title: string) => `[data-test-overview-card-title="${title}"]`,
    description: (title: string) => `[data-test-overview-card-subtitle="${title}"]`,
    content: (title: string) => `[data-test-overview-card-content="${title}"]`,
    actionText: (text: string) => `[data-test-action-text="${text}"]`,
    actionLink: (label: string) => `[data-test-overview-card="${label}"] a`,
  },
  pagination: {
    next: '.hds-pagination-nav__arrow--direction-next',
    prev: '.hds-pagination-nav__arrow--direction-prev',
  },
  // TODO move to KV selectors
  kvSuggestion: {
    input: '[data-test-kv-suggestion-input]',
    select: '[data-test-kv-suggestion-select]',
  },
};
