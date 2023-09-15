/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Base } from '../create';
import { clickable, visitable, create, fillable } from 'ember-cli-page-object';

export default create({
  ...Base,
  path: fillable('[data-test-secret-path="true"]'),
  secretKey: fillable('[data-test-secret-key]'),
  secretValue: fillable('[data-test-secret-value] textarea'),
  save: clickable('[data-test-secret-save]'),
  deleteBtn: clickable('[data-test-secret-delete] button'),
  confirmBtn: clickable('[data-test-confirm-button]'),
  visitEdit: visitable('/vault/secrets/:backend/edit/:id'),
  visitEditRoot: visitable('/vault/secrets/:backend/edit'),
  toggleJSON: clickable('[data-test-toggle-input="json"]'),
  startCreateSecret: clickable('[data-test-secret-create]'),
  deleteSecret() {
    return this.deleteBtn().confirmBtn();
  },
  createSecret: async function (path, key, value) {
    return this.path(path).secretKey(key).secretValue(value).save();
  },
  createSecretDontSave: async function (path, key, value) {
    return this.path(path).secretKey(key).secretValue(value);
  },
  editSecret: async function (key, value) {
    return this.secretKey(key).secretValue(value).save();
  },
});
