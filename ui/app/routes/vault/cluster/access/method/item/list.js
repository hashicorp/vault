import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import { singularize } from 'ember-inflector';
import ListRoute from 'vault/mixins/list-route';

export default Route.extend(ListRoute, {
  wizard: service(),
  pathHelp: service('path-help'),
  model() {
    const { item_type: itemType } = this.paramsFor('vault.cluster.access.method.item');
    const { page, pageFilter } = this.paramsFor(this.routeName);
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
  actions: {
    willTransition(transition) {
      window.scrollTo(0, 0);
      if (transition.targetName !== this.routeName) {
        this.store.clearAllDatasets();
      }
      return true;
    },
  },
  setupController(controller) {
    this._super(...arguments);
    const { item_type: itemType } = this.paramsFor('vault.cluster.access.method.item');
    let { path } = this.paramsFor('vault.cluster.access.method');
    controller.set('itemType', singularize(itemType));
    controller.set('method', path);
    const { apiPath } = this.modelFor('vault.cluster.access.method');
    this.pathHelp.getPaths(apiPath, path, itemType).then(paths => {
      controller.set('paths', Array.from(paths.list, pathInfo => pathInfo.path));
    });
  },
});
