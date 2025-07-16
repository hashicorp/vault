/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default class TransformTemplate extends Model {
  idPrefix = 'template/';

  @attr('string', {
    readOnly: true,
    subText:
      'Templates allow Vault to determine what and how to capture the value to be transformed. This cannot be edited later.',
  })
  name;

  @attr('string', { defaultValue: 'regex' }) type;

  @attr('string', {
    editType: 'regex',
    subText: 'The template’s pattern defines the data format. Expressed in regex.',
  })
  pattern;

  @attr('array', {
    subText:
      'Alphabet defines a set of characters (UTF-8) that is used for FPE to determine the validity of plaintext and ciphertext values. You can choose a built-in one, or create your own.',
    editType: 'searchSelect',
    isSectionHeader: true,
    fallbackComponent: 'string-list',
    label: 'Alphabet',
    models: ['transform/alphabet'],
    selectLimit: 1,
  })
  alphabet;

  @attr('string') encodeFormat;
  @attr('') decodeFormats;

  @attr('string', { readOnly: true }) backend;

  get readAttrs() {
    const keys = ['name', 'pattern', 'encodeFormat', 'decodeFormats', 'alphabet'];
    return expandAttributeMeta(this, keys);
  }

  get writeAttrs() {
    return expandAttributeMeta(this, ['name', 'pattern', 'alphabet']);
  }

  @lazyCapabilities(apiPath`${'backend'}/template/${'id'}`, 'backend') updatePath;
}
