/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { alias } from '@ember/object/computed';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

export default Model.extend({
  failedServerRead: attr('boolean'),
  auth: attr('string'),
  lease_duration: attr('number'),
  lease_id: attr('string'),
  renewable: attr('boolean'),

  backend: attr('string'),
  account_name: attr('string'),
  algorithm: attr('string'),
  digits: attr('number'),
  issuer: attr('string'),
  period: attr('number'),
  url: attr('string'),

  secretPath: lazyCapabilities(apiPath`${'backend'}/keys/${'id'}`, 'backend', 'id'),
  codePath: lazyCapabilities(apiPath`${'backend'}/codes/${'id'}`, 'backend', 'id'),
  canDelete: alias('secretPath.canDelete'),
  get canRead() {
    return this.secretPath.get('canRead') !== false && this.codePath.get('canRead') !== false;
  },
});
