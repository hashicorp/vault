/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export const SELECTORS = {
  toggleJson: '[data-test-json-view-toggle]',
  toggleMasked: '[data-test-button="toggle-masked"]',
  jsonEditor: '[data-test-component="code-mirror-modifier"]',
  tooltipTrigger: '[data-test-tooltip-trigger]',
  pageTitle: '[data-test-header-title]',
  infoRowValue: (label) => `[data-test-value-div="${label}"]`,
};

export const parseJsonEditor = (find) => {
  return JSON.parse(find(SELECTORS.jsonEditor).innerText);
};
