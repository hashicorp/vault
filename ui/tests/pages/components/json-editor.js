/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isPresent, notHasClass, text } from 'ember-cli-page-object';

export default {
  title: text('[data-test-component=json-editor-title]'),
  hasToolbar: isPresent('[data-test-component=json-editor-toolbar]'),
  hasJSONEditor: isPresent('[data-test-component="code-mirror-modifier"]'),
  canEdit: notHasClass('readonly-codemirror'),
};
