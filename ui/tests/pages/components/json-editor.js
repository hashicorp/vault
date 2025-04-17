/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isPresent, text } from 'ember-cli-page-object';

export default {
  title: text('[data-test-component=json-editor-title]'),
  hasToolbar: isPresent('[data-test-component=json-editor-toolbar]'),
  hasJSONEditor: isPresent('[data-test-component="code-mirror-modifier"]'),
  canEdit: notHasClass('readonly-codemirror'),
  readOnlyTitle: text('.hds-code-block__title'),
  readOnlyToolbar: isPresent('.hds-code-block__header'),
  readOnlyCopyButton: isPresent('.hds-code-block__copy-button'),
  readOnlyDisplay: isPresent('.hds-code-block__code'),
};
