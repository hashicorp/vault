/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';
import { alias, and, equal } from '@ember/object/computed';

export default Model.extend({
  name: attr('string'),
  //https://www.vaultproject.io/docs/http/sys-health.html
  initialized: attr('boolean'),
  sealed: attr('boolean'),
  isSealed: alias('sealed'),
  standby: attr('boolean'),
  isActive: equal('standby', false),
  clusterName: attr('string'),
  clusterId: attr('string'),

  isLeader: and('initialized', 'isActive'),

  //https://www.vaultproject.io/docs/http/sys-seal-status.html
  //The "t" parameter is the threshold, and "n" is the number of shares.
  t: attr('number'),
  n: attr('number'),
  progress: attr('number'),
  sealThreshold: alias('t'),
  sealNumShares: alias('n'),
  version: attr('string'),
  type: attr('string'),
  storageType: attr('string'),
  hcpLinkStatus: attr('string'),

  //https://www.vaultproject.io/docs/http/sys-leader.html
  haEnabled: attr('boolean'),
  isSelf: attr('boolean'),
  leaderAddress: attr('string'),
});
