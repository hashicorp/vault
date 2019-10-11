import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import { singularize } from 'ember-inflector';
import ListRoute from 'vault/mixins/list-route';

export default Route.extend(ListRoute, {
  wizard: service(),
  pathHelp: service('path-help'),

  getMethodAndModelInfo() {
    const { item_type: itemType } = this.paramsFor('vault.cluster.access.method.item');
    const { path: id } = this.paramsFor('vault.cluster.access.method');
    const methodModel = this.modelFor('vault.cluster.access.method');
    const { apiPath, type } = methodModel;
    return { apiPath, type, id, itemType, methodModel };
  },

  model() {
    const { type, id, itemType } = this.getMethodAndModelInfo();
    const { page, pageFilter } = this.paramsFor(this.routeName);
    let modelType = `generated-${singularize(itemType)}-${type}`;

    return this.store
      .lazyPaginatedQuery(modelType, {
        responsePath: 'data.keys',
        page: page,
        pageFilter: pageFilter,
        type: itemType,
        id: id,
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
    reload() {
      this.store.clearAllDatasets();
      this.refresh();
    },
  },
  setupController(controller) {
    this._super(...arguments);
    const { apiPath, method, itemType, methodModel } = this.getMethodAndModelInfo();
    controller.set('itemType', itemType);
    controller.set('method', method);
    controller.set('methodModel', methodModel);
    this.pathHelp.getPaths(apiPath, method, itemType).then(paths => {
      controller.set(
        'paths',
        paths.paths.filter(path => path.navigation && path.itemType.includes(itemType))
      );
    });
  },
});
