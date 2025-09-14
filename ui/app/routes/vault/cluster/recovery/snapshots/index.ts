/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type RouterService from '@ember/routing/router-service';
import type { SnapshotsRouteModel } from '../snapshots';

export default class RecoverySnapshotsIndexRoute extends Route {
  @service declare readonly router: RouterService;

  beforeModel() {
    const parentModel = this.modelFor('vault.cluster.recovery.snapshots') as SnapshotsRouteModel;

    if (parentModel.snapshots.length === 1) {
      const snapshot_id = parentModel.snapshots[0];
      this.router.transitionTo('vault.cluster.recovery.snapshots.snapshot.manage', snapshot_id);
    }
  }
}
