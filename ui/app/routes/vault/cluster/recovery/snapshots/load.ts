/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { SystemListStorageRaftSnapshotAutoConfigListEnum } from '@hashicorp/vault-client-typescript';

import type ApiService from 'vault/services/api';
import type { Breadcrumb } from 'vault/vault/app-types';
import type Controller from '@ember/controller';
import type { ModelFrom } from 'vault/vault/route';

export type SnapshotsLoadModel = ModelFrom<RecoverySnapshotsLoadRoute>;

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
}

export default class RecoverySnapshotsLoadRoute extends Route {
  @service declare readonly api: ApiService;

  async model() {
    const { configs, configError } = await this.fetchConfigs();
    return {
      configs,
      configError,
    };
  }

  async fetchConfigs() {
    let configs: string[], configError;

    try {
      const { keys } = await this.api.sys.systemListStorageRaftSnapshotAutoConfig(
        SystemListStorageRaftSnapshotAutoConfigListEnum.TRUE
      );
      configs = keys ?? [];
    } catch (e) {
      const error = await this.api.parseError(e);

      configError = error;

      if (!configError.message) {
        configError.message = 'Something went wrong';
      }

      configs = [];
    }

    return {
      configs,
      configError,
    };
  }

  setupController(controller: RouteController, resolvedModel: SnapshotsLoadModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Secrets Recovery', route: 'vault.cluster.recovery.snapshots' },
      { label: 'Upload', route: 'vault.cluster.recovery.snapshots.load' },
    ];
  }
}
