/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default Model.extend({
  idPrefix: 'alphabet/',
  idForNav: computed('id', 'idPrefix', function () {
    const modelId = this.id || '';
    return `${this.idPrefix}${modelId}`;
  }),

  name: attr('string', {
    readOnly: true,
    subText: 'The alphabet name. Keep in mind that spaces are not allowed and this cannot be edited later.',
  }),
  alphabet: attr('string', {
    label: 'Alphabet',
    subText:
      'Provide the set of valid UTF-8 characters contained within both the input and transformed value. Read more.',
  }),

  attrs: computed(function () {
    const keys = ['name', 'alphabet'];
    return expandAttributeMeta(this, keys);
  }),

  backend: attr('string', { readOnly: true }),
  updatePath: lazyCapabilities(apiPath`${'backend'}/alphabet/${'id'}`, 'backend', 'id'),
});
