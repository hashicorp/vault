/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isPresent, notHasClass, text } from 'ember-cli-page-object';

export const SELECTORS = {
  codeBlock: '.hds-code-block__code',
  copy: '.hds-code-block__copy-button',
  title: '[data-test-component="json-editor-title"]',
  toolbar: '.hds-code-block__header',
};
