import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import { singularize } from 'ember-inflector';

export default Route.extend({
  wizard: service(),
  pathHelp: service('path-help'),

  beforeModel() {
    const { apiPath, type, method, itemType } = this.getMethodAndModelInfo();
    let modelType = `generated-${singularize(itemType)}-${type}`;
    return this.pathHelp.getNewModel(modelType, method, apiPath, itemType);
  },

  getMethodAndModelInfo() {
    const { item_type: itemType } = this.paramsFor(this.routeName);
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    const methodModel = this.modelFor('vault.cluster.access.method');
    const { apiPath, type } = methodModel;
    return { apiPath, type, method, itemType };
  },

  setupController(controller) {
    this._super(...arguments);
    const { apiPath, method, itemType } = this.getMethodAndModelInfo();
    controller.set('itemType', itemType);
    controller.set('method', method);
    this.pathHelp.getPaths(apiPath, method, itemType).then(paths => {
      controller.set('paths', Array.from(paths.list, pathInfo => pathInfo.path));
    });
  },
});
