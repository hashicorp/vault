/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

import type Store from '@ember-data/store';

interface RouteParams {
  name: string;
  type: string;
}

export default class SyncSecretsDestinationsDestinationRoute extends Route {
  @service declare readonly store: Store;

  model(params: RouteParams) {
    const { name, type } = params;
    return this.store.findRecord(`sync/destinations/${type}`, name);
  }
}
