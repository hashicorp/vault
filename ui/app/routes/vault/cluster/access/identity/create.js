import Ember from 'ember';

export default Ember.Route.extend({
  model() {
    let itemType = this.modelFor('vault.cluster.access.identity');
    let modelType = `identity/${itemType}`;
    return this.store.createRecord(modelType);
  },
});
