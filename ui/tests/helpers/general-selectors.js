/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const SELECTORS = {
  breadcrumb: '[data-test-breadcrumbs] li',
  breadcrumbAtIdx: (idx) => `[data-test-crumb="${idx}"] a`,
  breadcrumbs: '[data-test-breadcrumbs]',
  title: '[data-test-page-title]',
  headerContainer: 'header.page-header',
  icon: (name) => `[data-test-icon="${name}"]`,
  tab: (name) => `[data-test-tab="${name}"]`,
  confirmModalInput: '[data-test-confirmation-modal-input]',
  confirmButton: '[data-test-confirm-button]',
  emptyStateTitle: '[data-test-empty-state-title]',
  emptyStateMessage: '[data-test-empty-state-message]',
  emptyStateActions: '[data-test-empty-state-actions]',
  menuTrigger: '[data-test-popup-menu-trigger]',
  // FORMS
  infoRowValue: (label) => `[data-test-value-div="${label}"]`,
  inputByAttr: (attr) => `[data-test-input="${attr}"]`,
  validation: (attr) => `[data-test-field="${attr}"] [data-test-inline-alert]`,
  messageError: '[data-test-message-error]',
};
