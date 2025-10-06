/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Factory } from 'miragejs';

export default Factory.extend({
  algorithm: 'SHA1',
  digits: 6,
  issuer: 'Vault',
  key_size: 20,
  max_validation_attempts: 5,
  name: '', // returned but cannot be set at this time
  namespace_path: 'admin/',
  period: 30,
  qr_size: 200,
  skew: 1,
  self_enrollment_enabled: false,
  type: 'totp',

  afterCreate(record) {
    if (record.name) {
      console.warn('Endpoint ignored these unrecognized parameters: [name]'); // eslint-disable-line
      record.name = '';
    }
  },
});
