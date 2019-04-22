import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import { getOwner } from '@ember/application';
import { singularize } from 'ember-inflector';
import ListRoute from 'vault/mixins/list-route';

export default Route.extend(ListRoute, {
  wizard: service(),
  pathHelp: service('path-help'),

  beforeModel(params) {
    const { item_type: itemType } = this.paramsFor(this.routeName);
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    let methodModel = this.modelFor('vault.cluster.access.method');
    let { apiPath, type } = methodModel;
    let modelType = `generated-${singularize(itemType)}-${type}`;
    let path = `${apiPath}${itemType}/example`;
    return this.pathHelp.getNewModel(modelType, getOwner(this), method, path, itemType);
  },

  model(params) {
    let { item_type: itemType, page, pageFilter } = this.paramsFor(this.routeName);
    const methodModel = this.modelFor('vault.cluster.access.method');
    const { type } = methodModel;
    let modelType = `generated-${singularize(itemType)}-${type}`;
    const { path: method } = this.paramsFor('vault.cluster.access.method');

    return this.store
      .lazyPaginatedQuery(modelType, {
        responsePath: 'data.keys',
        page: page,
        pageFilter: pageFilter,
        type: itemType,
        method: method,
      })
      .catch(err => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      });
  },

  setupController(controller) {
    this._super(...arguments);
    const { item_type: itemType } = this.paramsFor(this.routeName);
    const { apiPath } = this.modelFor('vault.cluster.access.method');
    let { path } = this.paramsFor('vault.cluster.access.method');
    controller.set('itemType', singularize(itemType));
    controller.set('method', path);
    this.pathHelp.getPaths(apiPath, path).then(paths => {
      controller.set('paths', paths);
    });
  },
});
