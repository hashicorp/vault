import Ember from 'ember';
import UnloadModelRoute from 'vault/mixins/unload-model-route';
import UnsavedModelRoute from 'vault/mixins/unsaved-model-route';

export default Ember.Route.extend(UnloadModelRoute, UnsavedModelRoute, {
  model(params) {
    let itemType = this.modelFor('vault.cluster.access.identity');
    let modelType = `identity/${itemType}-alias`;
    return this.store.createRecord(modelType, {
      canonicalId: params.item_id,
    });
  },
});
