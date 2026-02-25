/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { SystemListStorageRaftSnapshotLoadListEnum } from '@hashicorp/vault-client-typescript';

import type ApiService from 'vault/services/api';
import type Capabilities from 'vault/services/capabilities';
import type { ModelFrom } from 'vault/vault/route';
import type RouterService from '@ember/routing/router-service';
import type VersionService from 'vault/services/version';

export type SnapshotsRouteModel = ModelFrom<RecoverySnapshotsRoute>;

export default class RecoverySnapshotsRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly capabilities: Capabilities;
  @service declare readonly router: RouterService;
  @service declare readonly version: VersionService;

  async model() {
    if (this.version.isEnterprise) {
      const { canUpdate } = await this.capabilities.fetchPathCapabilities(
        'sys/storage/raft/snapshot/snapshot-load'
      );

      const snapshots = await this.fetchSnapshots();

      return {
        snapshots,
        canLoadSnapshot: canUpdate,
      };
    }
    return { snapshots: [], showCommunityMessage: true };
  }

  redirect(model: SnapshotsRouteModel) {
    if (Array.isArray(model.snapshots) && model.snapshots.length === 1) {
      const snapshot_id = model.snapshots[0];
      this.router.transitionTo('vault.cluster.recovery.snapshots.snapshot.manage', snapshot_id);
    }
  }

  async fetchSnapshots() {
    try {
      // This request needs to be made within the root namespace context to grab loaded snapshot keys as it is unsupported in any other namespace.
      // By default, the api service uses the current namespace context, so we'll need to specify otherwise.
      // Snapshot operations do not have this constraint.
      const { keys } = await this.api.sys.systemListStorageRaftSnapshotLoad(
        SystemListStorageRaftSnapshotLoadListEnum.TRUE,
        this.api.buildHeaders({ namespace: '' })
      );
      return keys as string[];
    } catch (e) {
      const error = await this.api.parseError(e);
      if (error.status === 404) {
        return [];
      }

      if (error.message === 'raft storage is not in use') {
        return {
          showRaftStorageMessage: true,
        };
      }

      throw error;
    }
  }
}
