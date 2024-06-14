/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

export default Model.extend({
  failedServerRead: attr('boolean'),
  auth: attr('string'),
  lease_duration: attr('number'),
  lease_id: attr('string'),
  renewable: attr('boolean'),

  backend: attr('string'),

  name: alias('account_name'),

  account_name: attr('string', {
    readOnly: true,
  }),
  algorithm: attr('string', {
    readOnly: true,
    possibleValues: ['SHA1', 'SHA256', 'SHA512'],
  }),
  digits: attr('number', {
    readOnly: true,
    possibleValues: [6, 8],
  }),
  issuer: attr('string', {
    readOnly: true,
  }),
  period: attr('number', {
    readOnly: true,
  }),

  url: attr('string'),

  attrs: computed('account_name', function () {
    const keys = ['account_name', 'algorithm', 'digits', 'issuer', 'period'];
    return expandAttributeMeta(this, keys);
  }),

  createAttrs: computed('url', function () {
    const keys = ['url'];
    return expandAttributeMeta(this, keys);
  }),

  fieldGroups: computed('account_name', function () {
    const defaultFields = ['account_name', 'algorithm', 'digits', 'issuer', 'period'];
    const groups = [
      { default: defaultFields },
      //{
      //  Options: [...fields],
      //},
    ];
    return fieldToAttrs(this, groups);
  }),
  keyPath: lazyCapabilities(apiPath`${'backend'}/keys/${'id'}`, 'backend', 'id'),
  canRead: alias('keyPath.canRead'),
  canUpdate: alias('keyPath.canUpdate'),
  canDelete: alias('keyPath.canDelete'),
});
