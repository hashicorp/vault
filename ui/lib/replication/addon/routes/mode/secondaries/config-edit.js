import { hash } from 'rsvp';
import Base from '../../replication-base';

export default Base.extend({
  modelPath: 'model.config',

  model(params) {
    return hash({
      cluster: this.modelFor('mode.secondaries'),
      config: this.store.findRecord('mount-filter-config', params.secondary_id),
      mounts: this.fetchMounts(),
    });
  },

  redirect(model) {
    const cluster = model.cluster;
    let replicationMode = this.paramsFor('mode').replication_mode;
    if (
      !this.get('version.hasPerfReplication') ||
      replicationMode !== 'performance' ||
      !cluster.get(`${replicationMode}.isPrimary`) ||
      !cluster.get('canAddSecondary')
    ) {
      return this.transitionTo('mode', replicationMode);
    }
  },
});
