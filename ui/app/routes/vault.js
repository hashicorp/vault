import Ember from 'ember';

export default Ember.Route.extend({
  version: Ember.inject.service(),
  beforeModel() {
    return this.get('version').fetchVersion();
  },
  model() {
    // hardcode single cluster
    const fixture = {
      data: {
        id: '1',
        type: 'cluster',
        attributes: {
          name: 'vault',
        },
      },
    };
    this.store.push(fixture);
    return this.store.peekAll('cluster');
  },

  redirect(model, transition) {
    if (model.get('length') === 1 && transition.targetName === 'vault.index') {
      return this.transitionTo('vault.cluster', model.get('firstObject.name'));
    }
  },
});
