/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { withModelValidations } from 'vault/decorators/model-validations';

const validations = {
  generateSigningKey: [
    {
      validator(model) {
        const { publicKey, privateKey, generateSigningKey } = model;
        // if generateSigningKey is false, both public and private keys are required
        if (!generateSigningKey && (!publicKey || !privateKey)) {
          return false;
        }
        return true;
      },
      message: 'Provide a Public and Private key or set "Generate Signing Key" to true.',
    },
  ],
  publicKey: [
    {
      validator(model) {
        const { publicKey, privateKey } = model;
        // regardless of generateSigningKey, if one key is set they both need to be set.
        return publicKey || privateKey ? publicKey && privateKey : true;
      },
      message: 'You must provide a Public and Private keys or leave both unset.',
    },
  ],
};
// there are more options available on the API, but the UI does not support them yet.
@withModelValidations(validations)
export default class SshCaConfig extends Model {
  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord
  @attr('string', { sensitive: true }) privateKey; // obfuscated, never returned by API
  @attr('string') publicKey;
  @attr('boolean', { defaultValue: true }) generateSigningKey;

  configurableParams = ['privateKey', 'publicKey', 'generateSigningKey'];

  get displayAttrs() {
    return this.formFields.filter((attr) => attr.name !== 'privateKey');
  }

  get formFields() {
    return expandAttributeMeta(this, this.configurableParams);
  }
}
