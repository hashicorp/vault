import Route from '@ember/routing/route';
import UnloadModelRoute from 'vault/mixins/unload-model-route';
import UnsavedModelRoute from 'vault/mixins/unsaved-model-route';

export default Route.extend(UnloadModelRoute, UnsavedModelRoute, {
  model(params) {
    debugger; // eslint-disable-line
    let methodType = this.modelFor('vault.cluster.access.method');
    let { itemType } = this.paramsFor(params);
    let modelType = `access/${methodType}`;
    return this.store.createRecord(modelType);
  },
});
