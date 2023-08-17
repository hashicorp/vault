/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { clickable, fillable, attribute } from 'ember-cli-page-object';
import { waitFor } from '@ember/test-helpers';
import fields from '../form-field';

export default {
  ...fields,
  cancelLinkHref: attribute('href', '[data-test-cancel-link]'),
  cancelLink: clickable('[data-test-cancel-link]'),
  name: fillable('[data-test-input="name"]'),
  disabled: clickable('[data-test-input="disabled"]'),
  metadataKey: fillable('[data-test-kv-key]'),
  metadataValue: fillable('[data-test-kv-value]'),
  type: fillable('[data-test-input="type"]'),
  submit: clickable('[data-test-identity-submit]'),
  delete: clickable('[data-test-confirm-action-trigger]'),
  confirmDelete: clickable('[data-test-confirm-button]'),
  waitForConfirm() {
    return waitFor('[data-test-confirm-button]');
  },
  waitForFlash() {
    return waitFor('[data-test-flash-message-body]');
  },
};
