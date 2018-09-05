import ClusterRouteBase from './cluster-route-base';
import Ember from 'ember';
import config from 'vault/config/environment';

const { inject } = Ember;

export default ClusterRouteBase.extend({
  flashMessages: inject.service(),
  version: inject.service(),
  wizard: inject.service(),
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
  activate() {
    this.get('wizard').set('initEvent', 'LOGIN');
    this.get('wizard').transitionTutorialMachine(this.get('wizard.currentState'), 'TOLOGIN');
  },
  actions: {
    willTransition(transition) {
      if (transition.targetName !== this.routeName) {
        this.get('wizard').transitionTutorialMachine(this.get('wizard.currentState'), 'INITDONE');
      }
    },
  },
});
