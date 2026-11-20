/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { tracked } from '@glimmer/tracking';
import { dependentKeyCompat } from '@ember/object/compat';
import uuid from 'core/utils/uuid';

export interface ReplicationData {
  cluster_id?: string;
  mode?: string;
  known_secondaries?: string[];
  secondaries?: unknown[];
  connection_state?: string;
  secondary_id?: string;
  primary_cluster_addr?: string;
  known_primary_cluster_addrs?: string[];
  primaries?: unknown[];
  state?: string; // stream-wal, merkle-diff, merkle-sync, idle
  last_remote_wal?: number;
  last_wal?: number;
  merkle_root?: string;
  merkle_sync_progress?: { sync_total_keys?: number; sync_progress?: number };
}

// Values are coerced like the Ember Data transforms did, and only when the API returned the key.
const asString = (value: unknown) => (value === undefined || value === null ? value : String(value));
const asNumber = (value: unknown) => {
  if (value === undefined || value === null || value === '') return value === undefined ? undefined : null;
  const number = Number(value);
  return Number.isNaN(number) ? null : number;
};
const asArray = (value: unknown) =>
  value === undefined ? undefined : Array.isArray(value) ? [...value] : [];
const asObject = (value: unknown) =>
  value === undefined ? undefined : value && typeof value === 'object' && !Array.isArray(value) ? value : {};

const MODE_DISPLAY: Record<string, string> = {
  disabled: 'Disabled',
  unknown: 'Unknown',
  bootstrapping: 'Bootstrapping',
  primary: 'Primary',
  secondary: 'Secondary',
  unsupported: 'Not supported',
};

/**
 * The replication state of one mode (dr or performance) from `sys/replication/status`.
 * An enabled mode keeps its id (the cluster id) between refreshes so that it can be updated in place.
 * Getters marked with @dependentKeyCompat are watched by classic computed properties (auth, replication cards).
 */
export default class ReplicationAttributesResource {
  readonly id: string;
  @tracked private data: ReplicationData;

  constructor(data: ReplicationData = {}) {
    this.data = data;
    this.id = data.cluster_id || uuid();
  }

  update(data: ReplicationData) {
    this.data = { ...this.data, ...data };
  }

  @dependentKeyCompat
  get clusterId() {
    return asString(this.data.cluster_id) as string | undefined;
  }
  get clusterIdDisplay() {
    return this.clusterId ? this.clusterId.split('-')[0] : null;
  }
  get mode() {
    return asString(this.data.mode) as string | undefined;
  }
  get replicationDisabled() {
    return typeof this.mode === 'string' && /disabled|unsupported/.test(this.mode);
  }
  get replicationUnsupported() {
    return typeof this.mode === 'string' && /unsupported/.test(this.mode);
  }
  get replicationEnabled() {
    return !this.replicationDisabled;
  }

  // primary attrs
  get isPrimary() {
    return typeof this.mode === 'string' && /primary/.test(this.mode);
  }
  get knownSecondaries() {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    return asArray(this.data.known_secondaries) as any[] | undefined;
  }
  get secondaries() {
    return asArray(this.data.secondaries) as unknown[] | undefined;
  }

  // secondary attrs
  @dependentKeyCompat
  get isSecondary() {
    return typeof this.mode === 'string' && /secondary/.test(this.mode);
  }
  @dependentKeyCompat
  get connection_state() {
    return asString(this.data.connection_state) as string | undefined;
  }
  get modeForUrl() {
    return this.mode === 'bootstrapping'
      ? 'bootstrapping'
      : (this.isSecondary && 'secondary') || (this.isPrimary && 'primary');
  }
  get modeForHeader() {
    // mode will be false or undefined if it calls the status endpoint while still setting up the cluster
    return this.mode || 'loading';
  }
  @dependentKeyCompat
  get secondaryId() {
    return asString(this.data.secondary_id) as string | undefined;
  }
  @dependentKeyCompat
  get primaryClusterAddr() {
    return asString(this.data.primary_cluster_addr) as string | undefined;
  }
  @dependentKeyCompat
  get knownPrimaryClusterAddrs() {
    return asArray(this.data.known_primary_cluster_addrs) as string[] | undefined;
  }
  @dependentKeyCompat
  get primaries() {
    return asArray(this.data.primaries) as unknown[] | undefined;
  }
  @dependentKeyCompat
  get state() {
    return asString(this.data.state) as string | undefined;
  }
  @dependentKeyCompat
  get lastRemoteWAL() {
    return asNumber(this.data.last_remote_wal) as number | null | undefined;
  }

  // attrs on primary and secondary
  get lastWAL() {
    return asNumber(this.data.last_wal) as number | null | undefined;
  }
  @dependentKeyCompat
  get merkleRoot() {
    return asString(this.data.merkle_root) as string | undefined;
  }
  get merkleSyncProgress() {
    return asObject(this.data.merkle_sync_progress) as ReplicationData['merkle_sync_progress'];
  }
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
  }
  get syncProgressPercent() {
    const syncProgress = this.syncProgress;
    if (!syncProgress) {
      return null;
    }
    const { progress, total } = syncProgress;
    return Math.floor(100 * ((progress as number) / (total as number)));
  }
  get modeDisplay() {
    return MODE_DISPLAY[this.mode as string] || 'Disabled';
  }
}
