import ClusterRouteBase from './cluster-route-base';
import Ember from 'ember';
import config from 'vault/config/environment';

const { RSVP, inject } = Ember;

export default ClusterRouteBase.extend({
  flashMessages: inject.service(),
  beforeModel() {
    this.store.unloadAll('auth-method');
    return this._super();
  },
  model() {
    let cluster = this._super(...arguments);
    return this.store
      .findAll('auth-method', {
        adapterOptions: {
          unauthenticated: true,
        },
      })
      .then(result => {
        return RSVP.hash({
          cluster,
          methods: result,
        });
      });
  },
  resetController(controller) {
    controller.set('wrappedToken', '');
    controller.set('authMethod', '');
  },

  afterModel() {
    if (config.welcomeMessage) {
      this.get('flashMessages').stickyInfo(config.welcomeMessage);
    }
  },
});
