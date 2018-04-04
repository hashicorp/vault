import Ember from 'ember';
import ClusterRoute from 'vault/mixins/cluster-route';
const { inject } = Ember;

export default Ember.Route.extend(ClusterRoute, {
  version: inject.service(),

  beforeModel() {
    return this.get('version').fetchFeatures().then(() => {
      return this._super(...arguments);
    });
  },

  model() {
    return this.modelFor('vault.cluster');
  },

  afterModel(model) {
    return Ember.RSVP
      .hash({
        canEnablePrimary: this.store
          .findRecord('capabilities', 'sys/replication/primary/enable')
          .then(c => c.get('canUpdate')),
        canEnableSecondary: this.store
          .findRecord('capabilities', 'sys/replication/secondary/enable')
          .then(c => c.get('canUpdate')),
      })
      .then(({ canEnablePrimary, canEnableSecondary }) => {
        Ember.setProperties(model, {
          canEnablePrimary,
          canEnableSecondary,
        });
        return model;
      });
  },
});
