/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import type { ModelFrom } from 'vault/vault/route';

import type NamespaceService from 'vault/services/namespace';

export type SnapshotManageModel = ModelFrom<RecoverySnapshotsSnapshotManageRoute>;

export default class RecoverySnapshotsSnapshotManageRoute extends Route {
  @service declare readonly namespace: NamespaceService;

  async model() {
    const snapshot = this.modelFor('vault.cluster.recovery.snapshots.snapshot');
    const namespaces = this.namespace.inRootNamespace ? await this.fetchNamespaces() : [];
    return {
      snapshot,
      namespaces,
    };
  }

  async fetchNamespaces() {
    try {
      // TODO update with api service to request sys/internal/ui/namespaces
      await this.namespace?.findNamespacesForUser?.perform();
      return this.namespace.accessibleNamespaces;
    } catch {
      return [];
    }
  }
}
