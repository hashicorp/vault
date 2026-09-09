/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

// A hash of cluster states to ensure that the status menu and replication dashboards
// display states and glyphs consistently
// this includes states for the primary vault cluster and the connection_state

export const CLUSTER_STATES = {
  running: {
    glyph: 'check-circle',
    isOk: true,
    isSyncing: false,
    color: 'success',
  },
  ready: {
    glyph: 'check-circle',
    isOk: true,
    isSyncing: false,
    color: 'success',
  },
  'stream-wals': {
    glyph: 'check-circle',
    isOk: true,
    isSyncing: false,
    color: 'success',
  },
  'merkle-diff': {
    glyph: 'sync-reverse',
    isOk: true,
    isSyncing: true,
    color: 'warning',
  },
  connecting: {
    glyph: 'sync-reverse',
    isOk: true,
    isSyncing: true,
    color: 'warning',
  },
  'merkle-sync': {
    glyph: 'sync-reverse',
    isOk: true,
    isSyncing: true,
    color: 'warning',
  },
  idle: {
    glyph: 'x-square',
    isOk: false,
    isSyncing: false,
    color: 'critical',
  },
  transient_failure: {
    glyph: 'x-circle',
    isOk: false,
    isSyncing: false,
    color: 'critical',
  },
  shutdown: {
    glyph: 'x-circle',
    isOk: false,
    isSyncing: false,
    color: 'critical',
  },
};

export function clusterStates([state]) {
  const defaultDisplay = {
    glyph: '',
    isOk: null,
    isSyncing: null,
  };
  return CLUSTER_STATES[state] || defaultDisplay;
}

export default buildHelper(clusterStates);
