import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import { singularize } from 'ember-inflector';
import ListRoute from 'vault/mixins/list-route';

export default Route.extend(ListRoute, {
  wizard: service(),
  pathHelp: service('path-help'),

  getMethodAndModelInfo() {
    const { item_type: itemType } = this.paramsFor('vault.cluster.access.method.item');
    let subItemType, parentID, parentType;
    //we have a nested item
    if (itemType.includes('~*')) {
      let types = itemType.split('~*');
      parentType = types[0];
      subItemType = itemType;
      //we have an ID (e.g. role~*my-role~*secret-id)
      if (types.length === 3) {
        subItemType = `${types[0]}~*${types[2]}`;
        parentID = types[1];
      }
    }
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    const methodModel = this.modelFor('vault.cluster.access.method');
    const itemModel = this.modelFor('vault.cluster.access.method.item');
    const { type } = methodModel;
    return { type, method, itemType, itemModel, methodModel, subItemType, parentID, parentType };
  },

  model() {
    const { type, method, itemType, parentID, subItemType } = this.getMethodAndModelInfo();
    const { page, pageFilter } = this.paramsFor(this.routeName);
    let modelType = `generated-${singularize(subItemType || itemType)}-${type}`;
    debugger;
    console.log(modelType);

    return this.store
      .lazyPaginatedQuery(modelType, {
        responsePath: 'data.keys',
        page,
        pageFilter,
        type: itemType,
        method,
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
    const {
      method,
      itemType,
      itemModel,
      methodModel,
      parentType,
      parentID,
      subItemType,
    } = this.getMethodAndModelInfo();
    controller.set('itemType', itemType);
    controller.set('subItemType', subItemType);
    controller.set('method', method);
    controller.set('parentType', parentType);
    controller.set('parentID', parentID);
    controller.set('methodModel', methodModel);
    controller.set('model.paths', itemModel.paths);
  },
});
