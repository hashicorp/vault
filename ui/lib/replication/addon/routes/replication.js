import { inject as service } from '@ember/service';
import { setProperties } from '@ember/object';
import { hash } from 'rsvp';
import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';

export default Route.extend(ClusterRoute, {
  version: service(),

  beforeModel() {
    return this.get('version')
      .fetchFeatures()
      .then(() => {
        return this._super(...arguments);
      });
  },

  model() {
    return this.modelFor('vault.cluster');
  },

  afterModel(model) {
    return hash({
      canEnablePrimary: this.store
        .findRecord('capabilities', 'sys/replication/primary/enable')
        .then(c => c.get('canUpdate')),
      canEnableSecondary: this.store
        .findRecord('capabilities', 'sys/replication/secondary/enable')
        .then(c => c.get('canUpdate')),
    }).then(({ canEnablePrimary, canEnableSecondary }) => {
      setProperties(model, {
        canEnablePrimary,
        canEnableSecondary,
      });
      return model;
    });
  },
});
