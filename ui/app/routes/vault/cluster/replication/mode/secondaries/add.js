import { hash } from 'rsvp';
import Base from '../../replication-base';

export default Base.extend({
  model() {
    return hash({
      cluster: this.modelFor('vault.cluster.replication.mode.secondaries'),
      mounts: this.fetchMounts(),
    });
  },

  redirect(model) {
    const replicationMode = this.get('replicationMode');
    if (!model.cluster.get(`${replicationMode}.isPrimary`) || !model.cluster.get('canAddSecondary')) {
      return this.transitionTo('vault.cluster.replication.mode', model.cluster.get('name'), replicationMode);
    }
  },

  setupController(controller, model) {
    controller.set('model', model.cluster);
    controller.set('mounts', model.mounts);
  },

  resetController(controller) {
    controller.reset();
  },
});
