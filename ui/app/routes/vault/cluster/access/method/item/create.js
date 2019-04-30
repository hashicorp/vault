import Route from '@ember/routing/route';
import UnloadModelRoute from 'vault/mixins/unload-model-route';
import UnsavedModelRoute from 'vault/mixins/unsaved-model-route';
import { singularize } from 'ember-inflector';

export default Route.extend(UnloadModelRoute, UnsavedModelRoute, {
  model(params) {
    const { item_type: itemType } = this.paramsFor('vault.cluster.access.method.item');
    const methodModel = this.modelFor('vault.cluster.access.method');
    const { type } = methodModel;
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    const modelType = `generated-${singularize(itemType)}-${type}`;
    return this.store.createRecord(
      modelType,
      { itemType: singularize(itemType) },
      {
        adapterOptions: { path: `${method}/${itemType}` },
      }
    );
  },
});
