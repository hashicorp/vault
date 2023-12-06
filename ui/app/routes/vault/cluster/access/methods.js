/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default Route.extend({
  store: service(),

  queryParams: {
    page: {
      refreshModel: true,
    },
    pageFilter: {
      refreshModel: true,
    },
  },

  model() {
    return this.store.findAll('auth-method');
  },
});
