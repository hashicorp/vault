/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { withExpandedAttributes } from 'vault/decorators/model-expanded-attributes';

@withExpandedAttributes()
export default class SshOtpCredential extends Model {
  @attr('object', {
    readOnly: true,
  })
  role;
  @attr('string', {
    label: 'IP Address',
  })
  ip;
  @attr('string') username;
  @attr('string', { masked: true }) key;
  @attr('string') keyType;
  @attr('number') port;

  get toCreds() {
    return this.key;
  }
}
