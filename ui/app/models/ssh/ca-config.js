/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default class SshCaConfig extends Model {
  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord
  @attr('string', { sensitive: true }) privateKey; // obfuscated, never returned by API
  @attr('string', { sensitive: true }) publicKey;
  @attr('boolean', { defaultValue: true })
  generateSigningKey;
  // there are more options available on the API, but the UI does not support them yet.
  get attrs() {
    const keys = ['publicKey', 'generateSigningKey'];
    return expandAttributeMeta(this, keys);
  }
}
