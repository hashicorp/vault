/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: MPL-2.0
 */

export const CUSTOM_MESSAGES = {
  // General selectors that are common between custom messages
  inlineErrorMessage: `[data-test-inline-error-message]`,
  unauthCreateFormInfo: '[data-test-unauth-info]',
  navLink: '[data-test-sidebar-nav-link="Custom messages"]',
  radio: (radioName: string) => `[data-test-radio="${radioName}"]`,
  field: (fieldName: string) => `[data-test-field="${fieldName}"]`,
  input: (input: string) => `[data-test-input="${input}"]`,
  fieldValidation: (fieldName: string) => `[data-test-validation-error="${fieldName}"]`,
  modal: (name: string) => `[data-test-modal="${name}"]`,
  modalTitle: (title: string) => `[data-test-modal-title="${title}"]`,
  modalBody: (name: string) => `[data-test-modal-body="${name}"]`,
  alert: (name: string) => `data-test-custom-alert=${name}`,
  alertTitle: (name: string) => `[data-test-custom-alert-title="${name}"]`,
  alertDescription: (name: string) => `[data-test-custom-alert-description="${name}"]`,
  alertAction: (name: string) => `[data-test-custom-alert-action="${name}"]`,
  badge: (name: string) => `[data-test-badge="${name}"]`,
  tab: (name: string) => `[data-test-custom-messages-tab="${name}"]`,
};
