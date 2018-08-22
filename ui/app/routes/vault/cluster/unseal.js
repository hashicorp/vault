import Ember from 'ember';
import ClusterRoute from './cluster-route-base';

const { inject } = Ember;

export default ClusterRoute.extend({
  wizard: inject.service(),

  activate() {
    let wizard = this.get('wizard');
    this.get('wizard').set('initEvent', 'UNSEAL');
    this.get('wizard').transitionTutorialMachine(this.get('wizard.currentState'), 'TOUNSEAL');
  },
});
