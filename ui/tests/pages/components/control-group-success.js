/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isPresent, clickable } from 'ember-cli-page-object';

export default {
  showsJsonViewer: isPresent('[data-test-json-viewer]'),
  showsNavigateMessage: isPresent('[data-test-navigate-message]'),
  showsUnwrapForm: isPresent('[data-test-unwrap-form]'),
  navigate: clickable('[data-test-navigate-button]'),
  unwrap: clickable('[data-test-unwrap-button]'),
};
