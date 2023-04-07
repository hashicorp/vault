/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { clickable, fillable } from 'ember-cli-page-object';

import fields from '../form-field';
export default {
  ...fields,
  ttlValue: fillable('[data-test-ttl-value]'),
  ttlUnit: fillable('[data-test-ttl-value]'),
  save: clickable('[data-test-save-config]'),
};
