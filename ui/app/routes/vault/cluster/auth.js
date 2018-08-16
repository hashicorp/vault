import ClusterRouteBase from './cluster-route-base';
import Ember from 'ember';
import config from 'vault/config/environment';

const { inject } = Ember;

export default ClusterRouteBase.extend({
  flashMessages: inject.service(),
  version: inject.service(),
  beforeModel() {
    return this._super().then(() => {
      return this.get('version').fetchFeatures();
    });
  },
  model() {
    return this._super(...arguments);
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
