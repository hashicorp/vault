/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class VaultClusterAccessMethodsRoute extends Route {
  @service store;

  queryParams = {
    page: {
      refreshModel: true,
    },
    pageFilter: {
      refreshModel: true,
    },
  };

  model() {
    // when initially mounting and configuring secret engine, ember-data creates a ghost auth-method model
    // we don't want this record appearing in the access list so unload records
    this.store.unloadAll('auth-method');
    return this.store.findAll('auth-method');
  }
}
