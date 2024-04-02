/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { findAll } from '@ember/test-helpers';

export const GENERAL = {
  breadcrumb: '[data-test-breadcrumbs] li',
  breadcrumbAtIdx: (idx) => `[data-test-breadcrumbs] li:nth-child(${idx + 1}) a`,
  breadcrumbs: '[data-test-breadcrumbs]',
  breadcrumbLink: (label) => `[data-test-breadcrumb="${label}"] a`,
  title: '[data-test-page-title]',
  headerContainer: 'header.page-header',
  icon: (name) => `[data-test-icon="${name}"]`,
  tab: (name) => `[data-test-tab="${name}"]`,
  filter: (name) => `[data-test-filter="${name}"]`,
  filterInput: '[data-test-filter-input]',
  confirmModalInput: '[data-test-confirmation-modal-input]',
  confirmButton: '[data-test-confirm-button]',
  confirmTrigger: '[data-test-confirm-action-trigger]',
  emptyStateTitle: '[data-test-empty-state-title]',
  emptyStateMessage: '[data-test-empty-state-message]',
  emptyStateActions: '[data-test-empty-state-actions]',
  menuTrigger: '[data-test-popup-menu-trigger]',
  listItem: '[data-test-list-item-link]',
  calendarWidget: {
    trigger: '[data-test-calendar-widget-trigger]',
    currentMonth: '[data-test-current-month]',
    currentBillingPeriod: '[data-test-current-billing-period]',
    customEndMonth: '[data-test-show-calendar]',
    previousYear: '[data-test-previous-year]',
    nextYear: '[data-test-next-year]',
    calendarMonth: (month) => `[data-test-calendar-month="${month}"]`,
  },
  dateDropdown: {
    toggleMonth: '[data-test-toggle-month]',
    toggleYear: '[data-test-toggle-year]',
    selectMonth: (month) => `[data-test-dropdown-month="${month}"]`,
    selectYear: (year) => `[data-test-dropdown-year="${year}"]`,
    submit: '[data-test-date-dropdown-submit]',
  },
  // FORMS
  infoRowLabel: (label) => `[data-test-row-label="${label}"]`,
  infoRowValue: (label) => `[data-test-value-div="${label}"]`,
  inputByAttr: (attr) => `[data-test-input="${attr}"]`,
  selectByAttr: (attr) => `[data-test-select="${attr}"]`,
  fieldByAttr: (attr) => `[data-test-field="${attr}"]`,
  enableField: (attr) => `[data-test-enable-field="${attr}"] button`,
  ttl: {
    toggle: (attr) => `[data-test-toggle-label="${attr}"]`,
    input: (attr) => `[data-test-ttl-value="${attr}"]`,
  },

  validation: (attr) => `[data-test-field-validation=${attr}]`,
  validationWarning: (attr) => `[data-test-validation-warning=${attr}]`,
  messageError: '[data-test-message-error]',
  kvObjectEditor: {
    deleteRow: (idx = 0) => `[data-test-kv-delete-row="${idx}"]`,
  },
  searchSelect: {
    options: '.ember-power-select-option',
    optionIndex: (text) => findAll('.ember-power-select-options li').findIndex((e) => e.innerText === text),
    option: (index = 0) => `.ember-power-select-option:nth-child(${index + 1})`,
    selectedOption: (index = 0) => `[data-test-selected-option="${index}"]`,
    noMatch: '.ember-power-select-option--no-matches-message',
    removeSelected: '[data-test-selected-list-button="delete"]',
  },
  overviewCard: {
    title: (title) => `[data-test-overview-card-title="${title}"]`,
    description: (title) => `[data-test-overview-card-subtitle="${title}"]`,
    content: (title) => `[data-test-overview-card-content="${title}"]`,
    action: (title) => `[data-test-overview-card-container="${title}"] [data-test-action-text]`,
    actionLink: (label) => `[data-test-action-text="${label}"]`,
  },
  pagination: {
    next: '.hds-pagination-nav__arrow--direction-next',
    prev: '.hds-pagination-nav__arrow--direction-prev',
  },
  kvSuggestion: {
    input: '[data-test-kv-suggestion-input]',
    select: '[data-test-kv-suggestion-select]',
  },
  navLink: (label) => `[data-test-sidebar-nav-link="${label}"]`,
  cancelButton: '[data-test-cancel]',
  saveButton: '[data-test-save]',
};
