/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { clickable, fillable, text, isPresent } from 'ember-cli-page-object';
import fields from '../form-field';

export default {
  ...fields,
  scope: '.config-pki',
  text: text('[data-test-text]'),
  title: text('[data-test-title]'),
  hasTitle: isPresent('[data-test-title]'),
  hasError: isPresent('[data-test-error]'),
  submit: clickable('[data-test-submit]'),
  enableTtl: clickable('[data-test-toggle-input]'),
  fillInValue: fillable('[data-test-ttl-value]'),
  fillInUnit: fillable('[data-test-select="ttl-unit"]'),
};
