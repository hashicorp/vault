import ClusterRouteBase from './cluster-route-base';
import Ember from 'ember';

const { inject } = Ember;

export default ClusterRouteBase.extend({
  wizard: inject.service(),

  beforeModel() {
    this._super(...arguments);
    debugger;
    if (this.get('wizard.currentState') === 'active.init.save') {
      this.get('wizard').transitionTutorialMachine(this.get('wizard.currentState'), 'CONTINUE');
    }
  },
  afterModel() {
    this._super(...arguments);
    debugger;
  },
});
