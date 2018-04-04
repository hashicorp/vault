import Ember from 'ember';
const SPLASH_DELAY = Ember.testing ? 0 : 300;

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
    return new Ember.RSVP.Promise(resolve => {
      Ember.run.later(() => {
        resolve(this.store.peekAll('cluster'));
      }, SPLASH_DELAY);
    });
  },

  redirect(model, transition) {
    if (model.get('length') === 1 && transition.targetName === 'vault.index') {
      return this.transitionTo('vault.cluster', model.get('firstObject.name'));
    }
  },
});
