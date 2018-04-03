import Ember from 'ember';

export default Ember.Route.extend({
  model() {
    return Ember.RSVP.hash({
      cluster: this.modelFor('vault.cluster'),
      seal: this.store.findRecord('capabilities', 'sys/seal'),
    });
  },
});
