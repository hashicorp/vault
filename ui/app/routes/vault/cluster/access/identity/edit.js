import Ember from 'ember';

export default Ember.Route.extend({
  model(params) {
    let itemType = this.modelFor('vault.cluster.access.identity');
    let modelType = `identity/${itemType}`;
    return this.store.findRecord(modelType, params.item_id);
  },
});
