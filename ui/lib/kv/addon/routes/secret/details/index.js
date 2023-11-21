/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { breadcrumbsForSecret } from 'kv/utils/kv-breadcrumbs';
import { inject as service } from '@ember/service';

export default class KvSecretDetailsIndexRoute extends Route {
  @service store;

  // polled by controller and called once when route initializes
  async fetchSyncStatus(model) {
    const { backend: mount, path: secretName } = model;
    const syncAdapter = this.store.adapterFor('sync/association');
    return syncAdapter.fetchDestinations({ mount, secretName });
  }

  async setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    const breadcrumbsArray = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'list' },
      ...breadcrumbsForSecret(resolvedModel.path, true),
    ];
    controller.breadcrumbs = breadcrumbsArray;
    controller.syncDestinations = await this.fetchSyncStatus(resolvedModel);
    controller.fetchSyncStatus = this.fetchSyncStatus;
    controller.pollSyncStatus.perform();
  }

  resetController(controller, isExiting) {
    if (isExiting) {
      controller.pollSyncStatus.cancelAll();
    }
  }
}
