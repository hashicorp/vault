/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import ListRoute from 'core/mixins/list-route';
import { service } from '@ember/service';

export default Route.extend(ListRoute, {
  store: service(),

  model(params) {
    const itemType = this.modelFor('vault.cluster.access.identity');
    const modelType = `identity/${itemType}`;
    return this.store
      .lazyPaginatedQuery(modelType, {
        responsePath: 'data.keys',
        page: params.page,
        pageFilter: params.pageFilter,
        sortBy: 'name',
      })
      .catch((err) => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      });
  },

  setupController(controller) {
    this._super(...arguments);
    controller.set('identityType', this.modelFor('vault.cluster.access.identity'));
  },

  actions: {
    willTransition(transition) {
      window.scrollTo(0, 0);
      if (transition.targetName !== this.routeName) {
        this.store.clearAllDatasets();
      }
      return true;
    },
    reload() {
      this.store.clearAllDatasets();
      this.refresh();
    },
  },
});
