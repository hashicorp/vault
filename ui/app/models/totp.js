/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import { withModelValidations } from 'vault/decorators/model-validations';
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

const ALGORITHMS = ['SHA1', 'SHA256', 'SHA512'];
const DIGITS = [6, 8];

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
};

@withModelValidations(validations)
export default class TotpModel extends Model {
  @attr('string') backend;

  @alias('account_name') name;

  @attr('string', {
    fieldValue: 'name',
    readOnly: true,
  })
  account_name;
  @attr('string', {
    readOnly: true,
    possibleValues: ALGORITHMS,
  })
  algorithm;
  @attr('number', {
    readOnly: true,
    possibleValues: DIGITS,
  })
  digits;
  @attr('string', {
    readOnly: true,
  })
  issuer;
  @attr('number', {
    readOnly: true,
  })
  period;

  @attr('string') url;

  @computed('account_name', function () {
    const keys = ['account_name', 'algorithm', 'digits', 'issuer', 'period'];
    return expandAttributeMeta(this, keys);
  })
  attrs;

  @computed('url', function () {
    const keys = ['url'];
    return expandAttributeMeta(this, keys);
  })
  createAttrs;

  @computed('account_name', function () {
    const defaultFields = ['account_name', 'algorithm', 'digits', 'issuer', 'period'];
    const groups = [
      { default: defaultFields },
      //{
      //  Options: [...fields],
      //},
    ];
    return fieldToAttrs(this, groups);
  })
  fieldGroups;

  @lazyCapabilities(apiPath`${'backend'}/keys/${'id'}`, 'backend', 'id') keyPath;
  @alias('keyPath.canRead') canRead;
  @alias('keyPath.canUpdate') canUpdate;
  @alias('keyPath.canDelete') canDelete;
}
