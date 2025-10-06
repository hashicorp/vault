/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type ApiService from 'vault/services/api';

export default class RecoverySnapshotsSnapshotRoute extends Route {
  @service declare readonly api: ApiService;

  model(params: { snapshot_id: string }) {
    // This request needs to be made within the root namespace context to grab the loaded snapshot as it is unsupported in any other namespace.
    // By default, the api service uses the current namespace context, so we'll need to specify otherwise.
    // Snapshot operations do not have this constraint.
    return this.api.sys.systemReadStorageRaftSnapshotLoadId(
      params.snapshot_id,
      this.api.buildHeaders({ namespace: '' })
    );
  }
}
