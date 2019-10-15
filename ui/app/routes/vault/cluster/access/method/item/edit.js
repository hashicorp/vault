import Route from '@ember/routing/route';
import UnloadModelRoute from 'vault/mixins/unload-model-route';
import UnsavedModelRoute from 'vault/mixins/unsaved-model-route';
import { singularize } from 'ember-inflector';

export default Route.extend(UnloadModelRoute, UnsavedModelRoute, {
  model(params) {
    const id = params.item_id;
    const { item_type: itemType } = this.paramsFor('vault.cluster.access.method.item');
    const methodModel = this.modelFor('vault.cluster.access.method');
    const modelType = `generated-${singularize(itemType)}-${methodModel.type}`;
    return this.store.queryRecord(modelType, { id, authMethodPath: methodModel.id });
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
