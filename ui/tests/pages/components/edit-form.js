/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { isPresent, text } from 'ember-cli-page-object';
import fields from './form-field';
export default {
  ...fields,
  deleteText: text('[data-test-edit-delete-text]'),
  showsDelete: isPresent('[data-test-edit-delete-text]'),
  errorText: text('[data-test-edit-form-error]'),
};
