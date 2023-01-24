import Route from '@ember/routing/route';
import UnloadModelRoute from 'vault/mixins/unload-model-route';
import { inject as service } from '@ember/service';

export default Route.extend(UnloadModelRoute, {
  store: service(),

  beforeModel() {
    const itemType = this.modelFor('vault.cluster.access.identity');
    if (itemType !== 'entity') {
      return this.transitionTo('vault.cluster.access.identity');
    }
  },

  model() {
    const modelType = `identity/entity-merge`;
    return this.store.createRecord(modelType);
  },
});
