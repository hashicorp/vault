/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
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
  },
  ready: {
    glyph: 'check-circle',
    isOk: true,
    isSyncing: false,
  },
  'stream-wals': {
    glyph: 'check-circle',
    isOk: true,
    isSyncing: false,
  },
  'merkle-diff': {
    glyph: 'sync-reverse',
    isOk: true,
    isSyncing: true,
  },
  connecting: {
    glyph: 'sync-reverse',
    isOk: true,
    isSyncing: true,
  },
  'merkle-sync': {
    glyph: 'sync-reverse',
    isOk: true,
    isSyncing: true,
  },
  idle: {
    glyph: 'x-square',
    isOk: false,
    isSyncing: false,
  },
  transient_failure: {
    glyph: 'x-circle',
    isOk: false,
    isSyncing: false,
  },
  shutdown: {
    glyph: 'x-circle',
    isOk: false,
    isSyncing: false,
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
