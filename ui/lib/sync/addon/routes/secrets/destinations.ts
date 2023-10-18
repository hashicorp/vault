/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

import type StoreService from 'vault/services/store';

interface SyncSecretsDestinationsRouteParams {
  page: string;
  pageFilter: string;
}

export default class SyncSecretsDestinationsRoute extends Route {
  @service declare readonly store: StoreService;

  async model(params: SyncSecretsDestinationsRouteParams) {
    return this.store.lazyPaginatedQuery(
      'sync/destination',
      {
        page: Number(params.page) || 1,
        pageFilter: params.pageFilter,
        responsePath: 'data.keys',
      },
      { adapterOptions: { showPartialError: true } }
    );
  }
}
