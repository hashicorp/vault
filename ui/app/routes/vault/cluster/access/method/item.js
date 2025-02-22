/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import { singularize } from 'ember-inflector';

export default Route.extend({
  pathHelp: service('path-help'),
  store: service(),

  async beforeModel() {
    const { apiPath, type, authMethodPath, itemType } = this.getMethodAndModelInfo();
    const modelType = `generated-${singularize(itemType)}-${type}`;
    await this.pathHelp.getNewModel(modelType, authMethodPath, apiPath, itemType);
    // getNewModel also creates an adapter if one does not exist and sets the apiPath value initially
    // this value will not change when routing between auth methods of the same type
    // in the generated-item-list adapter there is a short circuit to update the apiPath value on query({ list: true })
    // since we have removed that request to test the generated client it breaks the workflow
    // example -> navigate to userpass1, then userpass2, create a user and they will be created in userpass1
    // the apiPath value should be kept in sync at all times but since this will all be removed eventually -- hack it!
    const adapter = this.store.adapterFor(modelType);
    adapter.apiPath = apiPath;
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
