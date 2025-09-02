/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import { singularize } from 'ember-inflector';

export default Route.extend({
  pathHelp: service('path-help'),

  beforeModel() {
    const { apiPath, type, authMethodPath, itemType } = this.getMethodAndModelInfo();
    const modelType = `generated-${singularize(itemType)}-${type}`;
    return this.pathHelp.getNewModel(modelType, authMethodPath, apiPath, itemType);
  },

  getMethodAndModelInfo() {
    const { item_type: itemType } = this.paramsFor(this.routeName);
    const { path: authMethodPath } = this.paramsFor('vault.cluster.access.method');
    const methodModel = this.modelFor('vault.cluster.access.method');
    const { apiPath, type } = methodModel;
    return { apiPath, type, authMethodPath, itemType };
  },

  setupController(controller) {
    this._super(...arguments);
    const { apiPath, authMethodPath, itemType } = this.getMethodAndModelInfo();
    controller.set('itemType', itemType);
    this.pathHelp.getPaths(apiPath, authMethodPath, itemType).then((paths) => {
      const navigationPaths = paths.paths.filter((path) => path.navigation);
      controller.set(
        'paths',
        navigationPaths.filter((path) => path.itemType.includes(itemType)).map((path) => path.path)
      );
    });
  },
});
