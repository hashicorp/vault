import Model, { attr } from '@ember-data/model';
import { match, not } from '@ember/object/computed';
import { computed } from '@ember/object';

export default Model.extend({
  clusterId: attr('string'),
  clusterIdDisplay: computed('clusterId', 'mode', function () {
    const clusterId = this.clusterId;
    return clusterId ? clusterId.split('-')[0] : null;
  }),
  mode: attr('string'),
  replicationDisabled: match('mode', /disabled|unsupported/),
  replicationUnsupported: match('mode', /unsupported/),
  replicationEnabled: not('replicationDisabled'),

  // primary attrs
  isPrimary: match('mode', /primary/),

  knownSecondaries: attr('array'),
  secondaries: attr('array'),

  // secondary attrs
  isSecondary: match('mode', /secondary/),
  connection_state: attr('string'),
  modeForUrl: computed('isPrimary', 'isSecondary', 'mode', function () {
    const mode = this.mode;
    return mode === 'bootstrapping'
      ? 'bootstrapping'
      : (this.isSecondary && 'secondary') || (this.isPrimary && 'primary');
  }),
  modeForHeader: computed('mode', function () {
    const mode = this.mode;
    if (!mode) {
      // mode will be false or undefined if it calls the status endpoint while still setting up the cluster
      return 'loading';
    }
    return mode;
  }),
  secondaryId: attr('string'),
  primaryClusterAddr: attr('string'),
  knownPrimaryClusterAddrs: attr('array'),
  primaries: attr('array'),
  state: attr('string'), //stream-wal, merkle-diff, merkle-sync, idle
  lastRemoteWAL: attr('number'),

  // attrs on primary and secondary
  lastWAL: attr('number'),
  merkleRoot: attr('string'),
  merkleSyncProgress: attr('object'),
  get syncProgress() {
    const { state, merkleSyncProgress } = this;
    if (state !== 'merkle-sync' || !merkleSyncProgress) {
      return null;
    }
    const { sync_total_keys, sync_progress } = merkleSyncProgress;
    return {
      progress: sync_progress,
      total: sync_total_keys,
    };
  },

  syncProgressPercent: computed('syncProgress', function () {
    const syncProgress = this.syncProgress;
    if (!syncProgress) {
      return null;
    }
    const { progress, total } = syncProgress;

    return Math.floor(100 * (progress / total));
  }),

  modeDisplay: computed('mode', function () {
    const displays = {
      disabled: 'Disabled',
      unknown: 'Unknown',
      bootstrapping: 'Bootstrapping',
      primary: 'Primary',
      secondary: 'Secondary',
      unsupported: 'Not supported',
    };

    return displays[this.mode] || 'Disabled';
  }),
});
