/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import { generatedItemHelpUrl, generatedItemUrls } from 'vault/utils/generated-item-urls';

export default Route.extend({
  pathHelp: service('path-help'),

  // Field metadata and URLs come from the mount's OpenAPI help, shared by the list, create, edit and show routes.
  async model() {
    const { apiPath, authMethodPath, itemType } = this.getMethodAndModelInfo();
    const pathInfo = await this.pathHelp.getPaths(apiPath, authMethodPath, itemType);
    const formMetadata = await this.pathHelp.getFormMetadata(generatedItemHelpUrl(pathInfo, itemType));
    return { itemType, pathInfo, formMetadata, urls: generatedItemUrls(pathInfo, itemType) };
  },

  getMethodAndModelInfo() {
    const { item_type: itemType } = this.paramsFor(this.routeName);
    const { path: authMethodPath } = this.paramsFor('vault.cluster.access.method');
    const methodModel = this.modelFor('vault.cluster.access.method');
    const { apiPath, type } = methodModel;
    return { apiPath, type, authMethodPath, itemType };
  },

  setupController(controller, model) {
    this._super(...arguments);
    const { itemType, pathInfo } = model;
    controller.set('itemType', itemType);
    const navigationPaths = pathInfo.paths.filter((path) => path.navigation);
    controller.set(
      'paths',
      navigationPaths.filter((path) => path.itemType.includes(itemType)).map((path) => path.path)
    );
  },
});
