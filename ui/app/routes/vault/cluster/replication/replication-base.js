import Ember from 'ember';
import UnloadModelRouteMixin from 'vault/mixins/unload-model-route';

export default Ember.Route.extend(UnloadModelRouteMixin, {
  modelPath: 'model.config',
  fetchMounts() {
    return Ember.RSVP
      .hash({
        mounts: this.store.findAll('secret-engine'),
        auth: this.store.findAll('auth-method'),
      })
      .then(({ mounts, auth }) => {
        return Ember.RSVP.resolve(mounts.toArray().concat(auth.toArray()));
      });
  },

  version: Ember.inject.service(),
  rm: Ember.inject.service('replication-mode'),
  replicationMode: Ember.computed.alias('rm.mode'),
});
