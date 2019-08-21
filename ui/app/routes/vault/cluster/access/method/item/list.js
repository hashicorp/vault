import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import { singularize } from 'ember-inflector';
import ListRoute from 'vault/mixins/list-route';

export default Route.extend(ListRoute, {
  wizard: service(),
  pathHelp: service('path-help'),

  getMethodAndModelInfo() {
    const { item_type: itemType } = this.paramsFor('vault.cluster.access.method.item');
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    const methodModel = this.modelFor('vault.cluster.access.method');
    const itemModel = this.modelFor('vault.cluster.access.method.item');
    const { apiPath, type } = methodModel;
    return { apiPath, type, method, itemType, itemModel, methodModel };
  },

  model() {
    const { type, method, itemType } = this.getMethodAndModelInfo();
    const { page, pageFilter } = this.paramsFor(this.routeName);
    let modelType = `generated-${singularize(itemType)}-${type}`;
    console.log(modelType);

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
    reload() {
      this.store.clearAllDatasets();
      this.refresh();
    },
  },
  setupController(controller) {
    this._super(...arguments);
    const { method, itemType, itemModel, methodModel } = this.getMethodAndModelInfo();
    controller.set('itemType', itemType);
    controller.set('method', method);
    if (itemType.includes('_')) {
      controller.set('parentType', itemType.split('_')[0]);
    }
    controller.set('methodModel', methodModel);
    controller.set('model.paths', itemModel.paths);
  },
});
