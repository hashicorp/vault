import { alias, and, equal } from '@ember/object/computed';
import DS from 'ember-data';
const { attr } = DS;

export default DS.Model.extend({
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

  //https://www.vaultproject.io/docs/http/sys-leader.html
  haEnabled: attr('boolean'),
  isSelf: attr('boolean'),
  leaderAddress: attr('string'),
});
