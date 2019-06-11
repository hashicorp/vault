import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import { singularize } from 'ember-inflector';

export default Route.extend({
  wizard: service(),
  pathHelp: service('path-help'),

  beforeModel() {
    const { item_type: itemType } = this.paramsFor(this.routeName);
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    let methodModel = this.modelFor('vault.cluster.access.method');
    let { apiPath, type } = methodModel;
    let modelType = `generated-${singularize(itemType)}-${type}`;
    return this.pathHelp.getNewModel(modelType, method, apiPath, itemType);
  },

  setupController(controller) {
    this._super(...arguments);
    const { item_type: itemType } = this.paramsFor(this.routeName);
    const { path } = this.paramsFor('vault.cluster.access.method');
    const { apiPath } = this.modelFor('vault.cluster.access.method');
    controller.set('itemType', itemType);
    controller.set('method', path);
    this.pathHelp.getPaths(apiPath, path, itemType).then(paths => {
      controller.set('paths', Array.from(paths.list, pathInfo => pathInfo.path));
    });
  },
});
