import { singularize } from 'ember-inflector';
import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';

export default Route.extend({
  pathHelp: service('path-help'),
  model() {
    const { item_id: itemName } = this.paramsFor(this.routeName);
    const { item_type: itemType } = this.paramsFor('vault.cluster.access.method.item');
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    const methodModel = this.modelFor('vault.cluster.access.method');
    const { type } = methodModel;
    const modelType = `generated-${singularize(itemType)}-${type}`;
    // this is necessary to prevent name clashing in the case where there are 2 generated items
    // with the same name under different auth backends. for example, userpass/bob and radius/bob
    // TODO: this could cause problems if users named their user or group mygroup/foo/bar
    const id = `${method}/${itemName}`;
    return this.store.findRecord(modelType, id);
  },

  setupController(controller) {
    this._super(...arguments);
    const { item_type: itemType } = this.paramsFor('vault.cluster.access.method.item');
    controller.set('itemType', singularize(itemType));
  },
});
