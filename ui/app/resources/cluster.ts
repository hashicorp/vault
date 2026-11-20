/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { tracked } from '@glimmer/tracking';
import { dependentKeyCompat } from '@ember/object/compat';
import { service } from '@ember/service';
import { baseResourceFactory } from 'vault/resources/base-factory';
import ReplicationAttributesResource from 'vault/resources/replication-attributes';

import type ApiService from 'vault/services/api';
import type VersionService from 'vault/services/version';
import type NamespaceService from 'vault/services/namespace';
import type ReplicationModeService from 'vault/services/replication-mode';
import type { ReplicationData } from 'vault/resources/replication-attributes';

interface ClusterBase {
  id: string;
  name: string;
}

// sys/health plus the flags set while loading the cluster
interface ClusterData {
  standby?: boolean;
  status?: string;
  type?: string;
  mode?: string; // only ever 'unsupported', otherwise the mode is on dr and performance
  license?: { state?: string; expiry_time?: string };
  has_chroot_namespace?: boolean;
  replication_redacted?: boolean;
}

// sys/seal-status
interface SealStatus {
  initialized?: boolean;
  sealed?: boolean;
  t?: number; // seal threshold
  n?: number; // number of shares
  progress?: number;
  type?: string;
  storage_type?: string;
}

interface ReplicationStatus extends Partial<ClusterData> {
  dr?: ReplicationData;
  performance?: ReplicationData;
}

type ReplicationMode = 'dr' | 'performance';

// Values are coerced like the Ember Data transforms did, and only when the API returned the key.
const asBoolean = (value: unknown) => {
  if (value === undefined) return undefined;
  if (typeof value === 'string') return /^(true|t|1)$/i.test(value);
  if (typeof value === 'number') return value === 1;
  return Boolean(value);
};
const asString = (value: unknown) => (value === undefined || value === null ? value : String(value));

// The state of the single Vault cluster, built from sys/health, sys/seal-status and sys/replication/status.
export default class ClusterResource extends baseResourceFactory<ClusterBase>() {
  @service declare readonly api: ApiService;
  @service declare readonly version: VersionService;
  @service('namespace') declare readonly namespaceService: NamespaceService;
  // this service exposes what mode the UI is currently viewing
  @service('replication-mode') declare readonly rm: ReplicationModeService;

  @tracked private data: ClusterData = {};
  @tracked private sealStatus?: SealStatus;
  @tracked dr: ReplicationAttributesResource | null = null;
  @tracked performance: ReplicationAttributesResource | null = null;

  // set by the replication routes
  @tracked canEnablePrimaryDr?: boolean;
  @tracked canEnableSecondaryDr?: boolean;
  @tracked canEnablePrimaryPerformance?: boolean;
  @tracked canEnableSecondaryPerformance?: boolean;
  @tracked canAddSecondary?: boolean;
  @tracked canRevokeSecondary?: boolean;
  @tracked drModeInit?: string;

  /* Reading the cluster */

  private reloads = 0;
  private appliedReload = 0;

  // fetches the cluster state again and updates this cluster in place
  async reload() {
    const reload = ++this.reloads;
    const [health, sealStatus, replication] = await Promise.all([
      this.readHealth(),
      this.readSealStatus(),
      this.readReplicationStatus(),
    ]);
    // a response that was requested before one that is already applied would bring back stale state
    if (reload < this.appliedReload) return this;
    this.appliedReload = reload;

    // like an Ember Data push, values that are missing from a response keep their previous value
    let data = { ...this.data, ...health };
    if (sealStatus) {
      this.sealStatus = { ...this.sealStatus, ...sealStatus };
    }
    if (replication?.data) {
      const { dr, performance, ...rest } = replication.data;
      data = { ...data, ...rest };
      if (dr) this.dr = this.updateReplication(this.dr, dr);
      if (performance) this.performance = this.updateReplication(this.performance, performance);
    } else if (replication?.redacted) {
      data = { ...data, replication_redacted: true };
    }
    this.data = data;
    return this;
  }

  private async readHealth(): Promise<ClusterData> {
    try {
      const health = await this.api.sys.readHealthStatus(async (requestContext) => {
        await this.api.addQueryParams(requestContext, {
          standbycode: 200,
          sealedcode: 200,
          uninitcode: 200,
          drsecondarycode: 200,
          performancestandbycode: 200,
        });
        return {
          ...requestContext.init,
          headers: this.api.buildHeaders({ token: '', namespace: '' }).headers,
        };
      });
      return health as unknown as ClusterData;
    } catch {
      // sys/health will only fail when chroot set
      // because it's allowed in root namespace only and
      // configured to return a 200 response in other fail scenarios
      return { has_chroot_namespace: true };
    }
  }

  private async readSealStatus(): Promise<SealStatus | null> {
    try {
      // seal-status returns its fields at the top level (no `data` wrapper), so the typed sealStatus() resolves to undefined
      const { raw } = await this.api.sys.sealStatusRaw({
        headers: this.api.buildHeaders({ token: '', namespace: '' }).headers,
      });
      return await raw.json();
    } catch {
      return null;
    }
  }

  private async readReplicationStatus(): Promise<{ data?: ReplicationStatus; redacted?: boolean } | null> {
    if (!(this.version.isEnterprise && this.namespaceService.inRootNamespace)) return null;
    try {
      const { data } = await this.api.sys.systemReadReplicationStatus(this.api.buildHeaders({ token: '' }));
      return { data: data as ReplicationStatus | undefined };
    } catch (error) {
      // set redacted if result is an error which only happens when redacted
      const { response } = await this.api.parseError(error);
      return { redacted: !!response?.errors?.includes('disabled path') };
    }
  }

  // an enabled mode is updated in place, any other becomes a new record
  private updateReplication(current: ReplicationAttributesResource | null, data: ReplicationData) {
    if (current && data.cluster_id && current.id === data.cluster_id) {
      current.update(data);
      return current;
    }
    return new ReplicationAttributesResource(data);
  }

  /* Cluster */

  get standby() {
    return asBoolean(this.data.standby);
  }
  get status() {
    return asString(this.data.status);
  }
  get type() {
    return asString(this.data.type);
  }
  get license() {
    const { license } = this.data;
    return license === undefined ? undefined : license && typeof license === 'object' ? license : {};
  }
  // set when sys/health cannot be read
  get hasChrootNamespace() {
    return asBoolean(this.data.has_chroot_namespace);
  }
  // set when the replication endpoints are redacted
  get replicationRedacted() {
    return asBoolean(this.data.replication_redacted);
  }

  /* Licensing concerns */
  get licenseExpiry() {
    return this.license?.expiry_time;
  }
  get licenseState() {
    return this.license?.state;
  }

  /* Seal status */

  // needsInit and sealed are watched by the auth service's classic computed properties
  // a cluster without a readable seal status needs init
  @dependentKeyCompat
  get needsInit() {
    return !this.sealStatus || asBoolean(this.sealStatus.initialized) === false;
  }
  get unsealed() {
    return !!this.sealStatus && asBoolean(this.sealStatus.sealed) === false;
  }
  @dependentKeyCompat
  get sealed() {
    return !this.unsealed;
  }
  get sealThreshold() {
    return this.sealStatus?.t;
  }
  get sealProgress() {
    return this.sealStatus?.progress;
  }
  get sealType() {
    return this.sealStatus?.type;
  }
  get storageType() {
    return this.sealStatus?.storage_type;
  }
  get hasProgress() {
    return (this.sealProgress as number) >= 1;
  }
  get usingRaft() {
    return this.storageType === 'raft';
  }

  /* Replication */

  get mode() {
    return asString(this.data.mode);
  }
  get allReplicationDisabled() {
    return this.dr?.replicationDisabled && this.performance?.replicationDisabled;
  }
  get anyReplicationEnabled() {
    return this.dr?.replicationEnabled || this.performance?.replicationEnabled;
  }
  get drMode() {
    return this.dr?.mode;
  }
  get replicationMode() {
    return this.rm.mode as ReplicationMode | null;
  }
  get replicationModeForDisplay() {
    return this.replicationMode === 'dr' ? 'Disaster recovery' : 'Performance';
  }
  get replicationIsInitializing() {
    // a mode of null only happens when a cluster is being initialized
    // otherwise the mode will be 'disabled', 'primary', 'secondary'
    return !this.dr?.mode || !this.performance?.mode;
  }
  // returns the replication attributes of the mode that the UI is viewing
  get replicationAttrs() {
    return this.replicationMode ? this[this.replicationMode] : null;
  }
}
