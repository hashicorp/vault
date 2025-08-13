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
    return this.api.sys.systemReadStorageRaftSnapshotLoadId(params.snapshot_id);
  }
}
