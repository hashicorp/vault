/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';

import type ApiService from 'vault/services/api';
import type Capabilities from 'vault/services/capabilities';

export default class RecoverySnapshotsRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly capabilities: Capabilities;

  async model() {
    const path = 'sys/storage/raft/snapshot/snapshot-load';
    const capabilities = await this.capabilities.fetch([path]);

    const canLoadSnapshot = capabilities[path]?.canUpdate;

    try {
      const data = await this.api.sys.systemListStorageRaftSnapshotLoad(true);
      const snapshots = this.api.keyInfoToArray(data);

      return hash({
        snapshots,
        canLoadSnapshot,
      });
    } catch (e) {
      // return empty list of snapshots
      return hash({
        snapshots: [],
        canLoadSnapshot,
      });
    }
  }
}
