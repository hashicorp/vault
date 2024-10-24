/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

// TODO add validations
// there are more options available on the API, but the UI does not support them yet.
export default class GcpConfig extends Model {
  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord
  @attr('string') credentials; // JSON string
  @attr('string') serviceAccountEmail;
  @attr('string', {
    subText:
      'The audience claim value for plugin identity tokens. Must match an allowed audience configured for the target IAM OIDC identity provider.',
  })
  identityTokenAudience;
  @attr({
    label: 'Identity token TTL',
    helperTextDisabled:
      'The TTL of generated tokens. Defaults to 1 hour, turn on the toggle to specify a different value.',
    helperTextEnabled: 'The TTL of generated tokens.',
    subText: '',
    editType: 'ttl',
  })
  identityTokenTtl;
  @attr({
    label: 'Default Config TTL',
    editType: 'ttl',
  })
  ttl;
  @attr({
    label: 'Max TTL',
    editType: 'ttl',
  })
  maxTtl;

  get attrs() {
    const keys = [
      'credentials',
      'serviceAccountEmail',
      'identityTokenAudience',
      'identityTokenTtl',
      'ttl',
      'maxTtl',
    ];
    return expandAttributeMeta(this, keys);
  }
  // return private key for edit/create view
  // get formFields() {
  //   const keys = ['privateKey', 'publicKey', 'generateSigningKey'];
  //   return expandAttributeMeta(this, keys);
  // }
}
