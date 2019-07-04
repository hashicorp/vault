import { hash } from 'rsvp';
import Base from '../../replication-base';

export default Base.extend({
  model() {
    return hash({
      cluster: this.modelFor('mode.secondaries'),
      mounts: this.fetchMounts(),
    });
  },

  redirect(model) {
    const replicationMode = this.paramsFor('mode').replication_mode;
    if (!model.cluster.get(`${replicationMode}.isPrimary`) || !model.cluster.get('canAddSecondary')) {
      return this.transitionTo('mode', replicationMode);
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
