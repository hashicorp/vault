/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

import type StoreService from 'vault/services/store';

interface Params {
  type: string;
}

export default class SyncSecretsDestinationsCreateDestinationRoute extends Route {
  @service declare readonly store: StoreService;

  async model(params: Params) {
    const { type } = params;
    return this.store.createRecord(`sync/destinations/${type}`, { type });
  }
}
