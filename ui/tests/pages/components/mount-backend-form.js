/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { clickable, collection, fillable, text, value, attribute } from 'ember-cli-page-object';
import fields from './form-field';

export default {
  ...fields,
  header: text('[data-test-mount-form-header]'),
  submit: clickable('[data-test-mount-submit]'),
  back: clickable('[data-test-mount-back]'),
  path: fillable('[data-test-input="path"]'),
  toggleOptions: clickable('[data-test-toggle-group="Method Options"]'),
  pathValue: value('[data-test-input="path"]'),
  types: collection('[data-test-mount-type]', {
    select: clickable(),
    id: attribute('id'),
  }),
  type: fillable('[name="mount-type"]'),
  async selectType(type) {
    return this.types.filterBy('id', type)[0].select();
  },
  async mount(type, path) {
    await this.selectType(type);
    if (path) {
      await this.path(path).submit();
    } else {
      await this.submit();
    }
  },
};
