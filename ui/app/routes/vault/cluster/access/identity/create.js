import Route from '@ember/routing/route';
import UnloadModelRoute from 'vault/mixins/unload-model-route';
import UnsavedModelRoute from 'vault/mixins/unsaved-model-route';

export default Route.extend(UnloadModelRoute, UnsavedModelRoute, {
  model() {
    let itemType = this.modelFor('vault.cluster.access.identity');
    let modelType = `identity/${itemType}`;
    return this.store.createRecord(modelType);
  },
});
