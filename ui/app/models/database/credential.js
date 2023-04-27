/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';

export default Model.extend({
  username: attr('string'),
  password: attr('string'),
  leaseId: attr('string'),
  leaseDuration: attr('string'),
  lastVaultRotation: attr('string'),
  rotationPeriod: attr('number'),
  ttl: attr('number'),
  roleType: attr('string'),
});
