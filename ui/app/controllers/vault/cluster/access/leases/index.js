import Ember from 'ember';

export default Ember.Controller.extend({
  actions: {
    lookupLease(id) {
      this.transitionToRoute('vault.cluster.access.leases.show', id);
    },
  },
});
