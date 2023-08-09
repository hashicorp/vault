/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Base } from '../create';
import { isPresent, clickable, visitable, create, fillable } from 'ember-cli-page-object';

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
  toggleMetadata: clickable('[data-test-show-metadata-toggle]'),
  metadataTab: clickable('[data-test-secret-metadata-tab]'),
  hasMetadataFields: isPresent('[data-test-metadata-fields]'),
  maxVersion: fillable('[data-test-input="maxVersions"]'),
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
  createSecretWithMetadata: async function (path, key, value, maxVersion) {
    return this.path(path).secretKey(key).secretValue(value).toggleMetadata().maxVersion(maxVersion).save();
  },
  editSecret: async function (key, value) {
    return this.secretKey(key).secretValue(value).save();
  },
});
