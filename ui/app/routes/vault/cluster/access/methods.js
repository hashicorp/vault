/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

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
    return this.store.findAll('auth-method');
  }
}
