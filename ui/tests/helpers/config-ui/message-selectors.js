/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { GENERAL } from 'vault/tests/helpers/general-selectors';

export const PAGE = {
  // General selectors that are common between pages
  ...GENERAL,
  inlineErrorMessage: `[data-test-inline-error-message]`,
  unauthCreateFormInfo: '[data-test-unauth-info]',
  navLink: '[data-test-sidebar-nav-link="Custom Messages"]',
  radio: (radioName) => `[data-test-radio="${radioName}"]`,
  field: (fieldName) => `[data-test-field="${fieldName}"]`,
  input: (input) => `[data-test-input="${input}"]`,
  button: (buttonName) => `[data-test-button="${buttonName}"]`,
  fieldValidation: (fieldName) => `[data-test-field-validation="${fieldName}"]`,
  modal: (name) => `[data-test-modal="${name}"]`,
  modalTitle: (title) => `[data-test-modal-title="${title}"]`,
  modalBody: (name) => `[data-test-modal-body="${name}"]`,
  modalButton: (name) => `[data-test-modal-button="${name}"]`,
  alert: (name) => `data-test-custom-alert=${name}`,
  alertTitle: (name) => `[data-test-custom-alert-title="${name}"]`,
  alertDescription: (name) => `[data-test-custom-alert-description="${name}"]`,
  alertAction: (name) => `[data-test-custom-alert-action="${name}"]`,
  badge: (name) => `[data-test-badge="${name}"]`,
  tab: (name) => `[data-test-custom-messages-tab="${name}"]`,
  confirmActionButton: (name) => `[data-test-confirm-action="${name}"]`,
  listItem: (name) => `[data-test-list-item="${name}"]`,
};
