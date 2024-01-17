/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { clickable, collection, fillable, isPresent } from 'ember-cli-page-object';

export default {
  showsDuplicateError: isPresent('[data-test-duplicate-keys-warning]'),
  addRow: clickable('[data-test-kv-add-row]'),
  rows: collection('[data-test-kv-row]', {
    kvKey: fillable('[data-test-kv-key]'),
    kvVal: fillable('[data-test-kv-value]'),
    deleteRow: clickable('[data-test-kv-delete-row]'),
  }),
};
