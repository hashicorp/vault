/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/vault/app-types';
import type { ModelFrom } from 'vault/vault/route';

export type SnapshotManageModel = ModelFrom<RecoverySnapshotsSnapshotDetailsRoute>;

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
}

export default class RecoverySnapshotsSnapshotDetailsRoute extends Route {
  async model() {
    const snapshot = this.modelFor('vault.cluster.recovery.snapshots.snapshot');
    return {
      snapshot,
    };
  }

  setupController(controller: RouteController, resolvedModel: SnapshotManageModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard' },
      { label: 'Secrets recovery', route: 'vault.cluster.recovery.snapshots' },
      { label: 'Details' },
    ];
  }
}
