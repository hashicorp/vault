/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import { withModelValidations } from 'vault/decorators/model-validations';
import { withExpandedAttributes } from 'vault/decorators/model-expanded-attributes';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

const ALGORITHMS = ['SHA1', 'SHA256', 'SHA512'];
const DIGITS = [6, 8];
const SKEW = [0, 1];

const validations = {
  account_name: [
    { type: 'presence', message: "Account name can't be blank." },
    {
      type: 'containsWhiteSpace',
      message:
        "Account name contains whitespace. If this is desired, you'll need to encode it with %20 in API requests.",
      level: 'warn',
    },
  ],
  algorithm: [
    {
      validator: (value) => ALGORITHMS.includes(value),
      message: 'Algorithm must be one of ' + ALGORITHMS.join(', '),
    },
  ],
  digits: [
    { validator: (value) => DIGITS.includes(value), message: 'Digits must be one of ' + DIGITS.join(', ') },
  ],
  period: [{ type: 'number', message: 'Period must be a number.' }],

  skew: [{ validator: (value) => SKEW.includes(value), message: 'Skew must be one of ' + SKEW.join(', ') }],
};

@withModelValidations(validations)
@withExpandedAttributes()
export default class TotpModel extends Model {
  @attr('string', {
    readOnly: true,
  })
  backend;

  @alias('account_name') name;

  @attr('string', {
    fieldValue: 'name',
    editDisabled: true,
  })
  account_name;
  @attr('string', {
    editDisabled: true,
    possibleValues: ALGORITHMS,
    defaultValue: 'SHA1',
  })
  algorithm;
  @attr('number', {
    editDisabled: true,
    possibleValues: DIGITS,
    defaultValue: 6,
  })
  digits;
  @attr('string', {
    editDisabled: true,
  })
  issuer;
  @attr('number', {
    editDisabled: true,
    defaultValue: 30,
  })
  period;

  @attr('boolean', {
    defaultValue: false,
    label: 'Use Vault as provider for this key',
    editDisabled: true,
  })
  generate;

  // Used when generate is true
  @attr('number', {
    //defaultValue: 20,
    editDisabled: true,
  })
  key_size;
  @attr('number', {
    possibleValues: SKEW,
    //defaultValue: 1,
    editDisabled: true,
  })
  skew;
  @attr('boolean', {
    //defaultValue: true,
    editDisabled: true,
  })
  exported;

  // Doesn't really make sense as we can generate our own QR code from the url
  @attr('number', {
    defaultValue: 0,
    editDisabled: true,
  })
  qr_size;

  // Used when generate is false
  @attr('string', {
    editDisabled: true,
  })
  url;
  @attr('string', {
    editDisabled: true,
  })
  key;

  // Returned when a key is created as provider
  @attr('string', {
    readOnly: true,
  })
  barcode;

  get attrs() {
    const keys = ['account_name', 'algorithm', 'digits', 'issuer', 'period'];
    return keys.map((k) => this.allByKey[k]);
  }

  get generatedAttrs() {
    const keys = ['url'];
    return keys.map((k) => this.allByKey[k]);
  }

  @computed('generate', function () {
    const defaultFields = ['generate'];
    const options = ['algorithm', 'digits', 'period'];
    const providerOptions = [];

    if (this.generate) {
      providerOptions.push('key_size', 'skew', 'exported', 'qr_size');
    } else {
      defaultFields.push('url', 'key');
    }

    defaultFields.push('account_name', 'issuer');

    const groups = [
      { default: defaultFields },
      {
        Options: [...options],
      },
    ];

    if (this.generate) {
      groups.push({
        'Provider options': [...providerOptions],
      });
    }

    return this._expandGroups(groups);
  })
  fieldGroups;

  @lazyCapabilities(apiPath`${'backend'}/keys/${'id'}`, 'backend', 'id') keyPath;
  @alias('keyPath.canRead') canRead;
  @alias('keyPath.canUpdate') canUpdate;
  @alias('keyPath.canDelete') canDelete;
}
