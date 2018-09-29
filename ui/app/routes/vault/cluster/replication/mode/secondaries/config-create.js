import { hash } from 'rsvp';
import { inject as service } from '@ember/service';
import Base from '../../replication-base';

export default Base.extend({
  flashMessages: service(),

  modelPath: 'model.config',

  findOrCreate(id) {
    const flash = this.get('flashMessages');
    return this.store
      .findRecord('mount-filter-config', id)
      .then(() => {
        // if we find a record, transition to the edit view
        return this.transitionTo('vault.cluster.replication.mode.secondaries.config-edit', id)
          .followRedirects()
          .then(() => {
            flash.info(
              `${id} already had a mount filter config, so we loaded the config edit screen for you.`
            );
          });
      })
      .catch(e => {
        if (e.httpStatus === 404) {
          return this.store.createRecord('mount-filter-config', {
            id,
          });
        } else {
          throw e;
        }
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

  model(params) {
    return hash({
      cluster: this.modelFor('vault.cluster.replication.mode'),
      config: this.findOrCreate(params.secondary_id),
      mounts: this.fetchMounts(),
    });
  },
});
