/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export const SELECTORS = {
  inputByAttr: (attr) => `[data-test-input="${attr}"]`,
  toggleInput: (attr) => `[data-test-input="${attr}"] input`,
  intervalDuration: '[data-test-ttl-value="Automatic tidy enabled"]',
};
