/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default class IdentityToken extends Model {
  @attr('string') issuer;

  get attrs() {
    const keys = ['issuer'];
    return expandAttributeMeta(this, keys);
  }
}
