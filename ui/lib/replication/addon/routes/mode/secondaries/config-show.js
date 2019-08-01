import { hash, resolve } from 'rsvp';
import Base from '../../replication-base';

export default Base.extend({
  modelPath: 'model.config',

  model(params) {
    const id = params.secondary_id;
    return hash({
      cluster: this.modelFor('application'),
      config: this.store.findRecord('mount-filter-config', id).catch(e => {
        if (e.httpStatus === 404) {
          // return an empty obj to let them nav to create
          return resolve({ id });
        } else {
          throw e;
        }
      }),
    });
  },
  redirect(model) {
    const cluster = model.cluster;
    const replicationMode = this.paramsFor('mode').replication_mode;
    if (
      !this.get('version.hasPerfReplication') ||
      replicationMode !== 'performance' ||
      !cluster.get(`${replicationMode}.isPrimary`)
    ) {
      return this.transitionTo('mode', cluster.get('name'), replicationMode);
    }
  },
});
