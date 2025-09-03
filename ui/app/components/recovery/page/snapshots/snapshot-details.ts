/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';

import type ApiService from 'vault/services/api';
import type RouterService from '@ember/routing/router-service';
import type { SnapshotManageModel } from 'vault/routes/vault/cluster/recovery/snapshots/snapshot/manage';

interface Args {
  model: SnapshotManageModel;
}

export default class SnapshotDetails extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly router: RouterService;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
  }

  @action
  async unloadSnapshot() {
    try {
      const { snapshot_id } = this.args.model.snapshot as { snapshot_id: string };

      await this.api.sys.systemDeleteStorageRaftSnapshotLoadId(snapshot_id);
      this.router.transitionTo('vault.cluster.recovery.snapshots');
    } catch (e) {
      // TODO error handling
      // const error = this.api.parseError(e);
    }
  }
}
