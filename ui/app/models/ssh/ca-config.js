/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { withModelValidations } from 'vault/decorators/model-validations';

const validations = {
  publicKey: [
    {
      validator(model) {
        const { publicKey, privateKey } = model;
        return (publicKey && privateKey) || (!publicKey && !privateKey) ? true : false;
      },
      message: 'Public Key and Private Key are both required if one of them is set.',
    },
  ],
  generateSigningKey: [
    {
      validator(model) {
        const { publicKey, privateKey, generateSigningKey } = model;
        if (!generateSigningKey && (!publicKey || !privateKey)) {
          return false;
        }
        return true;
      },
      message: 'Public Key and Private Key are both required if Generate Signing Key is false.',
    },
  ],
};
// there are more options available on the API, but the UI does not support them yet.
@withModelValidations(validations)
export default class SshCaConfig extends Model {
  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord
  @attr('string', { sensitive: true }) privateKey; // obfuscated, never returned by API
  @attr('string', { sensitive: true }) publicKey;
  @attr('boolean', { defaultValue: true })
  generateSigningKey;

  // do not return private key for configuration.index view
  get attrs() {
    const keys = ['publicKey', 'generateSigningKey'];
    return expandAttributeMeta(this, keys);
  }
  // return private key for edit/create view
  get formFields() {
    const keys = ['privateKey', 'publicKey', 'generateSigningKey'];
    return expandAttributeMeta(this, keys);
  }
}
