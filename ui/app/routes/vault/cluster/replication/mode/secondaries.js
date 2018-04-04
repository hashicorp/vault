import Ember from 'ember';

export default Ember.Route.extend({
  model() {
    const replicationMode = this.paramsFor('vault.cluster.replication.mode').replication_mode;

    return Ember.RSVP
      .hash({
        cluster: this.modelFor('vault.cluster.replication.mode'),
        canAddSecondary: this.store
          .findRecord('capabilities', `sys/replication/${replicationMode}/primary/secondary-token`)
          .then(c => c.get('canUpdate')),
        canRevokeSecondary: this.store
          .findRecord('capabilities', `sys/replication/${replicationMode}/primary/revoke-secondary`)
          .then(c => c.get('canUpdate')),
      })
      .then(({ cluster, canAddSecondary, canRevokeSecondary }) => {
        Ember.setProperties(cluster, {
          canRevokeSecondary,
          canAddSecondary,
        });
        return cluster;
      });
  },
  afterModel(model) {
    const replicationMode = this.paramsFor('vault.cluster.replication.mode').replication_mode;
    if (
      !model.get(`${replicationMode}.isPrimary`) ||
      model.get(`${replicationMode}.replicationDisabled`) ||
      model.get(`${replicationMode}.replicationUnsupported`)
    ) {
      return this.transitionTo('vault.cluster.replication.mode', replicationMode);
    }
  },
});
