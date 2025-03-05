/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
// eslint-disable-next-line ember/no-computed-properties-in-native-classes
import { alias } from '@ember/object/computed';
import { withFormFields } from 'vault/decorators/model-form-fields';
import { withModelValidations } from 'vault/decorators/model-validations';
import { withExpandedAttributes } from 'vault/decorators/model-expanded-attributes';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

const ALGORITHMS = ['SHA1', 'SHA256', 'SHA512'];
const DIGITS = [6, 8];
const SKEW = [0, 1];

const validations = {
  accountName: [
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
  name: [
    { type: 'presence', message: "Name can't be blank." },
    {
      type: 'containsWhiteSpace',
      message:
        "Name contains whitespace. If this is desired, you'll need to encode it with %20 in API requests.",
      level: 'warn',
    },
  ],
  period: [{ type: 'number', message: 'Period must be a number.' }],

  skew: [{ validator: (value) => SKEW.includes(value), message: 'Skew must be one of ' + SKEW.join(', ') }],
};

const generatedDefaultFields = ['name', 'generate', 'issuer', 'accountName'];
const nonGeneratedDefaultFields = [...generatedDefaultFields, 'url', 'key'];
const totpCodeOptions = ['algorithm', 'digits', 'period'];
const providerOptions = ['key_size', 'skew', 'exported', 'qr_size'];

const generatedFormFieldGroups = [
  {
    default: generatedDefaultFields,
  },
  {
    'TOTP Code Options': totpCodeOptions,
  },
  {
    'Provider Options': providerOptions,
  },
];

const nonGeneratedFormFieldGroups = [
  {
    default: nonGeneratedDefaultFields,
  },
  {
    'TOTP Code Options': totpCodeOptions,
  },
];

const formFieldGroupsCombined = {
  generatedFormFieldGroups,
  nonGeneratedFormFieldGroups,
};
@withModelValidations(validations)
@withExpandedAttributes()
@withFormFields(null, formFieldGroupsCombined)
export default class TotpKeyModel extends Model {
  @attr('string', {
    readOnly: true,
  })
  backend;

  @attr('string') name;
  @attr('string') accountName;

  @attr('string', {
    possibleValues: ALGORITHMS,
    defaultValue: 'SHA1',
  })
  algorithm;

  @attr('number', {
    possibleValues: DIGITS,
    defaultValue: 6,
  })
  digits;

  @attr('string') issuer;

  @attr({
    label: 'Period',
    editType: 'ttl',
    helperTextEnabled: 'How long each generated TOTP is valid.',
    defaultValue: 30, // API accepts both an integer as seconds and string with unit e.g 30 || '30s'
  })
  period;

  @attr('boolean', {
    defaultValue: false,
    label: 'Use Vault as provider for this key',
  })
  generate;

  // Used when generate is true
  @attr('number', {
    defaultValue: 20,
  })
  key_size;

  @attr('number', {
    possibleValues: SKEW,
    defaultValue: 1,
  })
  skew;

  @attr('boolean', {
    defaultValue: true,
  })
  exported;

  // Doesn't really make sense as we can generate our own QR code from the url
  @attr('number', {
    defaultValue: 0,
  })
  qr_size;

  // Used when generate is false
  @attr('string', {
    label: 'URL',
    helpText:
      'If a URL is provided the other fields can be left empty. E.g. otpauth://totp/Vault:test@test.com?secret=Y64VEVMBTSXCYIWRSHRNDZW62MPGVU2G&issuer=Vault',
  })
  url;

  @attr('string', {
    label: 'Shared master key',
  })
  key;

  // Returned when a key is created as provider
  @attr('string', {
    readOnly: true,
  })
  barcode;

  get attrs() {
    const keys = ['accountName', 'name', 'algorithm', 'digits', 'issuer', 'period'];
    return keys.map((k) => this.allByKey[k]);
  }

  get generatedAttrs() {
    const keys = ['url'];
    return keys.map((k) => this.allByKey[k]);
  }

  @lazyCapabilities(apiPath`${'backend'}/keys/${'id'}`, 'backend', 'id') keyPath;
  @alias('keyPath.canRead') canRead;
  @alias('keyPath.canDelete') canDelete;
  //TODO remove these aliases
}
