/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

export const TOOLS_SELECTORS = {
  submit: '[data-test-tools-submit]',
  toolsInput: (attr: string) => `[data-test-tools-input="${attr}"]`,
  button: (action: string) => `[data-test-button="${action}"]`,
};
