import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import { singularize } from 'ember-inflector';

export default Route.extend({
  wizard: service(),
  pathHelp: service('path-help'),

  beforeModel() {
    let { apiPath, type, method, itemType, itemID, subItemType, parentType } = this.getMethodAndModelInfo();
    let newModelFetches = [];
    return this.pathHelp.getPaths(apiPath, method, itemType).then(paths => {
      paths.itemTypes.forEach(item => {
        let modelType = `generated-${item}-${type}`;
        newModelFetches.push(this.pathHelp.getNewModel(modelType, method, apiPath, item, itemID));
      });
      return Promise.all(newModelFetches);
    });
  },

  model() {
    const { type, method, itemType, subItemType } = this.getMethodAndModelInfo();
    const modelType = `generated-${singularize(subItemType)}-${type}`;
    return this.store.createRecord(modelType, {
      itemType,
      method,
      adapterOptions: { path: `${method}/${itemType}` },
    });
  },

  getMethodAndModelInfo() {
    const { item_type: itemType } = this.paramsFor(this.routeName);
    const { item_id: itemID } = this.paramsFor('vault.cluster.access.method.item.show');
    debugger;
    let subItemType, parentID, parentType;
    //we have a nested item
    if (itemType.includes('~*')) {
      let types = itemType.split('~*');
      parentType = types[0];
      subItemType = itemType;
      //we have an ID
      if (types.length === 3) {
        subItemType = `${types[0]}~*${types[2]}`;
        parentID = types[1];
      }
    }
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    const methodModel = this.modelFor('vault.cluster.access.method');
    const itemModel = this.modelFor('vault.cluster.access.method.item');
    const { type, apiPath } = methodModel;
    return {
      type,
      apiPath,
      method,
      itemType,
      itemModel,
      methodModel,
      subItemType,
      itemID,
      parentID,
      parentType,
    };
  },

  setupController(controller) {
    this._super(...arguments);
    const { apiPath, method, itemType, itemID, subItemType } = this.getMethodAndModelInfo();
    controller.set('itemType', itemType);
    controller.set('method', method);
    this.pathHelp.getPaths(apiPath, method, subItemType || itemType, itemID).then(paths => {
      let navigationPaths = paths.paths.filter(path => path.navigation);
      controller.set(
        'paths',
        navigationPaths.filter(path => path.itemType.includes(itemType)).map(path => path.path)
      );
    });
  },
});
