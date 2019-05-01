import Route from '@ember/routing/route';
import UnloadModelRoute from 'vault/mixins/unload-model-route';
import UnsavedModelRoute from 'vault/mixins/unsaved-model-route';
import { singularize } from 'ember-inflector';

export default Route.extend(UnloadModelRoute, UnsavedModelRoute, {
  model(params) {
    const methodModel = this.modelFor('vault.cluster.access.method');
    const { type } = methodModel;
    const { item_type: itemType } = this.paramsFor('vault.cluster.access.method.item');
    let modelType = `generated-${singularize(itemType)}-${type}`;
    return this.store.findRecord(modelType, params.item_id);
  },

  setupController(controller) {
    this._super(...arguments);
    const { item_type: itemType } = this.paramsFor('vault.cluster.access.method.item');
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    const { item_id: itemName } = this.paramsFor(this.routeName);
    controller.set('itemType', singularize(itemType));
    controller.set('mode', 'edit');
    controller.set('method', method);
    controller.set('itemName', itemName);
  },
});
