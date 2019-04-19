import { hash } from 'rsvp';
import Base from '../../replication-base';

export default Base.extend({
  modelPath: 'model.config',

  model(params) {
    return hash({
      cluster: this.modelFor('vault.cluster.replication.mode.secondaries'),
      config: this.store.findRecord('mount-filter-config', params.secondary_id),
      mounts: this.fetchMounts(),
    });
  },

  redirect(model) {
    const cluster = model.cluster;
    const replicationMode = this.get('replicationMode');
    if (
      !this.get('version.hasPerfReplication') ||
      replicationMode !== 'performance' ||
      !cluster.get(`${replicationMode}.isPrimary`) ||
      !cluster.get('canAddSecondary')
    ) {
      return this.transitionTo('vault.cluster.replication.mode', cluster.get('name'), replicationMode);
    }
  },
});
