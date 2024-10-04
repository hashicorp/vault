/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Base } from '../create';
import { clickable, create, fillable } from 'ember-cli-page-object';
import { settled } from '@ember/test-helpers';

export default create({
  ...Base,
  path: fillable('[data-test-secret-path="true"]'),
  secretKey: fillable('[data-test-secret-key]'),
  secretValue: fillable('[data-test-secret-value] textarea'),
  save: clickable('[data-test-secret-save]'),
  toggleJSON: clickable('[data-test-toggle-input="json"]'),
  createSecret: async function (path, key, value) {
    await this.path(path).secretKey(key).secretValue(value).save();
    await settled();
    return;
  },
});
