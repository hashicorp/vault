/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default class Alphabet extends Model {
  idPrefix = 'alphabet/';

  get idForNav() {
    const modelId = this.id || '';
    return `${this.idPrefix}${modelId}`;
  }

  @attr('string', {
    readOnly: true,
    subText: 'The alphabet name. Keep in mind that spaces are not allowed and this cannot be edited later.',
  })
  name;

  @attr('string', {
    label: 'Alphabet',
    subText:
      'Provide the set of valid UTF-8 characters contained within both the input and transformed value.',
    docLink: '/vault/api-docs/secret/transform#create-update-alphabet',
  })
  alphabet;

  get attrs() {
    const keys = ['name', 'alphabet'];
    return expandAttributeMeta(this, keys);
  }

  @attr('string', {
    readOnly: true,
  })
  backend;

  @lazyCapabilities(apiPath`${'backend'}/alphabet/${'id'}`, 'backend', 'id')
  updatePath;
}
