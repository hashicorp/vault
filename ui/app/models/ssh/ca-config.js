/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default class AwsRootConfig extends Model {
  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord
  @attr('string', { sensitive: true }) privateKey;
  @attr('string', { sensitive: true }) publicKey;
  @attr('boolean', {
    defaultValue: true,
  })
  generateSigningKey;
  // TODO: there are more options available on the API, but the UI does not support them yet.
  get attrs() {
    // do not show private key, not returned by the API
    const keys = ['publicKey', 'generateSigningKey'];
    return expandAttributeMeta(this, keys);
  }
}
