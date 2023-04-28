/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { reads } from '@ember/object/computed';
import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
const CREATE_FIELDS = ['username', 'ip'];

const DISPLAY_FIELDS = ['username', 'ip', 'key', 'keyType', 'port'];
export default Model.extend({
  role: attr('object', {
    readOnly: true,
  }),
  ip: attr('string', {
    label: 'IP Address',
  }),
  username: attr('string'),
  key: attr('string'),
  keyType: attr('string'),
  port: attr('number'),
  attrs: computed('key', function () {
    const keys = this.key ? DISPLAY_FIELDS.slice(0) : CREATE_FIELDS.slice(0);
    return expandAttributeMeta(this, keys);
  }),
  toCreds: reads('key'),
});
