/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
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
