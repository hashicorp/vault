/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { SystemListStorageRaftSnapshotLoadListEnum } from '@hashicorp/vault-client-typescript';

import type ApiService from 'vault/services/api';
import type Capabilities from 'vault/services/capabilities';
import type RouterService from '@ember/routing/router-service';
import type { ModelFrom } from 'vault/vault/route';

type SnapshotRouteModel = ModelFrom<RecoverySnapshotsRoute>;

export default class RecoverySnapshotsRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly capabilities: Capabilities;
  @service declare readonly router: RouterService;

  async model() {
    const path = 'sys/storage/raft/snapshot/snapshot-load';
    const capabilities = await this.capabilities.fetch([path]);
    const canLoadSnapshot = capabilities[path]?.canUpdate;
    const snapshots = await this.fetchSnapshots();

    return {
      snapshots,
      canLoadSnapshot,
    };
  }

  afterModel(model: SnapshotRouteModel) {
    if (model.snapshots.length === 1) {
      const snapshot_id = model.snapshots[0];
      this.router.transitionTo('vault.cluster.recovery.snapshots.snapshot.manage', snapshot_id);
    }
  }

  // todo: If this req fails bc user cannot list snapshots, they cannot use the UI (confirm with product)
  async fetchSnapshots() {
    try {
      const { keys } = await this.api.sys.systemListStorageRaftSnapshotLoad(
        SystemListStorageRaftSnapshotLoadListEnum.TRUE
      );
      return keys;
    } catch (e) {
      const { message, status } = await this.api.parseError(e);
      if (status === 404) {
        return [];
      }
      throw message;
    }
  }
}
