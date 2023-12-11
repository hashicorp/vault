/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const SELECTORS = {
  breadcrumb: '[data-test-breadcrumbs] li',
  breadcrumbAtIdx: (idx) => `[data-test-breadcrumbs] li:nth-child(${idx + 1}) a`,
  breadcrumbs: '[data-test-breadcrumbs]',
  title: '[data-test-page-title]',
  headerContainer: 'header.page-header',
  icon: (name) => `[data-test-icon="${name}"]`,
  tab: (name) => `[data-test-tab="${name}"]`,
  filter: (name) => `[data-test-filter="${name}"]`,
  confirmModalInput: '[data-test-confirmation-modal-input]',
  confirmButton: '[data-test-confirm-button]',
  emptyStateTitle: '[data-test-empty-state-title]',
  emptyStateMessage: '[data-test-empty-state-message]',
  emptyStateActions: '[data-test-empty-state-actions]',
  menuTrigger: '[data-test-popup-menu-trigger]',
  // FORMS
  infoRowValue: (label) => `[data-test-value-div="${label}"]`,
  inputByAttr: (attr) => `[data-test-input="${attr}"]`,
  validation: (attr) => `[data-test-field-validation=${attr}]`,
  messageError: '[data-test-message-error]',
  searchSelect: {
    option: (index = 0) => `.ember-power-select-option:nth-child(${index + 1})`,
    selectedOption: (index = 0) => `[data-test-selected-option="${index}"]`,
    noMatch: '.ember-power-select-option--no-matches-message',
    removeSelected: '[data-test-selected-list-button="delete"]',
  },
  overviewCard: {
    title: (title) => `[data-test-overview-card="${title}"] h3`,
    description: (title) => `[data-test-overview-card-container="${title}"] p`,
    content: (title) => `[data-test-overview-card-content="${title}"]`,
    action: (title) => `[data-test-overview-card="${title}"] a`,
  },
  pagination: {
    next: '.hds-pagination-nav__arrow--direction-next',
    prev: '.hds-pagination-nav__arrow--direction-prev',
  },
};
