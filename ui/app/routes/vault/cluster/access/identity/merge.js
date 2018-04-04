import Ember from 'ember';
import UnloadModelRoute from 'vault/mixins/unload-model-route';

export default Ember.Route.extend(UnloadModelRoute, {
  beforeModel() {
    let itemType = this.modelFor('vault.cluster.access.identity');
    if (itemType !== 'entity') {
      return this.transitionTo('vault.cluster.access.identity');
    }
  },
  model() {
    let modelType = `identity/entity-merge`;
    return this.store.createRecord(modelType);
  },
});
