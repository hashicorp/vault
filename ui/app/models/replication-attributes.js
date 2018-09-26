import { match, not } from '@ember/object/computed';
import { computed } from '@ember/object';
import attr from 'ember-data/attr';
import Fragment from 'ember-data-model-fragments/fragment';

export default Fragment.extend({
  clusterId: attr('string'),
  clusterIdDisplay: computed('mode', function() {
    const clusterId = this.get('clusterId');
    return clusterId ? clusterId.split('-')[0] : null;
  }),
  mode: attr('string'),
  replicationDisabled: match('mode', /disabled|unsupported/),
  replicationUnsupported: match('mode', /unsupported/),
  replicationEnabled: not('replicationDisabled'),

  // primary attrs
  isPrimary: match('mode', /primary/),

  knownSecondaries: attr('array'),

  // secondary attrs
  isSecondary: match('mode', /secondary/),

  modeForUrl: computed('mode', function() {
    const mode = this.get('mode');
    return mode === 'bootstrapping'
      ? 'bootstrapping'
      : (this.get('isSecondary') && 'secondary') || (this.get('isPrimary') && 'primary');
  }),
  secondaryId: attr('string'),
  primaryClusterAddr: attr('string'),
  knownPrimaryClusterAddrs: attr('array'),
  state: attr('string'), //stream-wal, merkle-diff, merkle-sync, idle
  lastRemoteWAL: attr('number'),

  // attrs on primary and secondary
  lastWAL: attr('number'),
  merkleRoot: attr('string'),
  merkleSyncProgress: attr('object'),
  syncProgress: computed('state', 'merkleSyncProgress', function() {
    const { state, merkleSyncProgress } = this.getProperties('state', 'merkleSyncProgress');
    if (state !== 'merkle-sync' || !merkleSyncProgress) {
      return null;
    }
    const { sync_total_keys, sync_progress } = merkleSyncProgress;
    return {
      progress: sync_progress,
      total: sync_total_keys,
    };
  }).volatile(),

  syncProgressPercent: computed('syncProgress', function() {
    const syncProgress = this.get('syncProgress');
    if (!syncProgress) {
      return null;
    }
    const { progress, total } = syncProgress;

    return Math.floor(100 * (progress / total));
  }),

  modeDisplay: computed('mode', function() {
    const displays = {
      disabled: 'Disabled',
      unknown: 'Unknown',
      bootstrapping: 'Bootstrapping',
      primary: 'Primary',
      secondary: 'Secondary',
      unsupported: 'Not supported',
    };

    return displays[this.get('mode')] || 'Disabled';
  }),
});
