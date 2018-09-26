import Route from '@ember/routing/route';
import UnloadModelRoute from 'vault/mixins/unload-model-route';
import UnsavedModelRoute from 'vault/mixins/unsaved-model-route';

export default Route.extend(UnloadModelRoute, UnsavedModelRoute, {
  model(params) {
    let itemType = this.modelFor('vault.cluster.access.identity');
    let modelType = `identity/${itemType}-alias`;
    return this.store.findRecord(modelType, params.item_alias_id);
  },
});
