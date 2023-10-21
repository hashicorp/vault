/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

import type StoreService from 'vault/services/store';

interface SyncSecretsDestinationsIndexRouteParams {
  page: string;
  pageFilter: string;
}

export default class SyncSecretsDestinationsIndexRoute extends Route {
  @service declare readonly store: StoreService;

  async model(params: SyncSecretsDestinationsIndexRouteParams) {
    return this.store.lazyPaginatedQuery('sync/destination', {
      page: Number(params.page) || 1,
      pageFilter: params.pageFilter,
      responsePath: 'data.keys',
    });
  }
}
