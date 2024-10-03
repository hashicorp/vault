/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import UnloadModel from 'vault/mixins/unload-model-route';

export default Route.extend(UnloadModel, {
  store: service(),

  queryParams: {
    page: {
      refreshModel: true,
    },
  },

  version: service(),

  beforeModel() {
    this.store.unloadAll('namespace');
    return this.version.fetchFeatures().then(() => {
      return this._super(...arguments);
    });
  },

  model(params) {
    if (this.version.hasNamespaces) {
      return this.store
        .lazyPaginatedQuery('namespace', {
          responsePath: 'data.keys',
          page: Number(params?.page) || 1,
        })
        .then((model) => {
          return model;
        })
        .catch((err) => {
          if (err.httpStatus === 404) {
            return [];
          } else {
            throw err;
          }
        });
    }
    return null;
  },

  setupController(controller, model) {
    const has404 = this.has404;
    controller.setProperties({
      model: model,
      has404,
      hasModel: true,
    });
    if (!has404) {
      controller.setProperties({
        page: Number(model?.meta?.currentPage) || 1,
      });
    }
  },

  actions: {
    error(error, transition) {
      /* eslint-disable-next-line ember/no-controller-access-in-routes */
      const hasModel = this.controllerFor(this.routeName).hasModel;
      if (hasModel && error.httpStatus === 404) {
        this.set('has404', true);
        transition.abort();
      } else {
        return true;
      }
    },
    willTransition(transition) {
      window.scrollTo(0, 0);
      if (!transition || transition.targetName !== this.routeName) {
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
