/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import { singularize } from 'ember-inflector';
import ListRoute from 'vault/mixins/list-route';

export default Route.extend(ListRoute, {
  store: service(),
  pathHelp: service('path-help'),

  getMethodAndModelInfo() {
    const { item_type: itemType } = this.paramsFor('vault.cluster.access.method.item');
    const { path: authMethodPath } = this.paramsFor('vault.cluster.access.method');
    const methodModel = this.modelFor('vault.cluster.access.method');
    const { apiPath, type } = methodModel;
    return { apiPath, type, authMethodPath, itemType, methodModel };
  },

  model() {
    const { type, authMethodPath, itemType } = this.getMethodAndModelInfo();
    const { page, pageFilter } = this.paramsFor(this.routeName);
    const modelType = `generated-${singularize(itemType)}-${type}`;

    return this.store
      .lazyPaginatedQuery(modelType, {
        responsePath: 'data.keys',
        page: page,
        pageFilter: pageFilter,
        type: itemType,
        id: authMethodPath,
      })
      .catch((err) => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      });
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

  setupController(controller) {
    this._super(...arguments);
    const { apiPath, authMethodPath, itemType, methodModel } = this.getMethodAndModelInfo();
    controller.set('itemType', itemType);
    controller.set('methodModel', methodModel);
    this.pathHelp.getPaths(apiPath, authMethodPath, itemType).then((paths) => {
      controller.set(
        'paths',
        paths.paths.filter((path) => path.navigation && path.itemType.includes(itemType))
      );
    });
  },
});
