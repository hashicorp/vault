import Ember from 'ember';
import ClusterRoute from './cluster-route-base';

const { inject } = Ember;

export default ClusterRoute.extend({
  wizard: inject.service(),

  activate() {
    this.get('wizard').set('initEvent', 'UNSEAL');
    this.get('wizard').transitionTutorialMachine(this.get('wizard.currentState'), 'TOUNSEAL');
  },
});
