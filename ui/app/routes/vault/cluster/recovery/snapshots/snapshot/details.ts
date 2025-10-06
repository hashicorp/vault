/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';

import type { ModelFrom } from 'vault/vault/route';

export type SnapshotManageModel = ModelFrom<RecoverySnapshotsSnapshotDetailsRoute>;

export default class RecoverySnapshotsSnapshotDetailsRoute extends Route {
  async model() {
    const snapshot = this.modelFor('vault.cluster.recovery.snapshots.snapshot');
    return {
      snapshot,
    };
  }
}
