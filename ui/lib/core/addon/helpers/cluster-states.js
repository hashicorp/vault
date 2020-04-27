import { helper as buildHelper } from '@ember/component/helper';

// A hash of cluster states to ensure that the status menu and replication dashboards
// display states and glyphs consistently
// this includes states for the primary vault cluster and the connection_state

export const CLUSTER_STATES = {
  running: {
    glyph: 'check-circle-outline',
    isOk: true,
    isSyncing: false,
  },
  ready: {
    glyph: 'check-circle-outline',
    isOk: true,
    isSyncing: false,
  },
  'stream-wals': {
    glyph: 'check-circle-outline',
    display: 'Streaming',
    isOk: true,
    isSyncing: false,
  },
  'merkle-diff': {
    glyph: 'android-sync',
    display: 'Determining sync status',
    isOk: true,
    isSyncing: true,
  },
  connecting: {
    glyph: 'android-sync',
    display: 'Streaming',
    isOk: true,
    isSyncing: true,
  },
  'merkle-sync': {
    glyph: 'android-sync',
    display: 'Syncing',
    isOk: true,
    isSyncing: true,
  },
  idle: {
    glyph: 'cancel-square-outline',
    isOk: false,
    isSyncing: false,
  },
  transient_failure: {
    glyph: 'cancel-circle-outline',
    isOk: false,
    isSyncing: false,
  },
  shutdown: {
    glyph: 'cancel-circle-outline',
    isOk: false,
    isSyncing: false,
  },
};

export function clusterStates([state]) {
  const defaultDisplay = {
    glyph: '',
    display: '',
    isOk: null,
    isSyncing: null,
  };
  return CLUSTER_STATES[state] || defaultDisplay;
}

export default buildHelper(clusterStates);
