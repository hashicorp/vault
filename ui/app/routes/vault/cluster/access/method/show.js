import { next } from '@ember/runloop';
import { singularize } from 'ember-inflector';
import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import { getOwner } from '@ember/application';

export default Route.extend({
  pathHelp: service('path-help'),
  beforeModel() {
    const { item_type: itemType, item_id: itemName } = this.paramsFor(this.routeName);
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    let methodModel = this.modelFor('vault.cluster.access.method');
    let { apiPath, type } = methodModel;
    let modelType = `generated-${singularize(itemType)}-${type}`;
    let path = `${apiPath}${itemType}/${itemName}`;
    return this.pathHelp.getNewModel(modelType, getOwner(this), method, path, itemType);
  },
  model() {
    let { item_type: itemType, item_id: itemName } = this.paramsFor(this.routeName);
    let { path: method } = this.paramsFor('vault.cluster.access.method');
    let methodModel = this.modelFor('vault.cluster.access.method');
    let { type } = methodModel;
    let modelType = `generated-${singularize(itemType)}-${type}`;
    return this.store.findRecord(modelType, itemName, {
      adapterOptions: { path: `${method}/${itemType}` },
    });
  },

  setupController(controller, model) {
    debugger; // eslint-disable-line
    this._super(...arguments);
    const { item_type: itemType, item_id: itemId } = this.paramsFor(this.routeName);
    let { path } = this.paramsFor('vault.cluster.access.method');
    controller.set('props', model.toJSON());
    controller.set('itemType', singularize(itemType));
    controller.set('method', path);
    controller.set('id', itemId);
  },
});
