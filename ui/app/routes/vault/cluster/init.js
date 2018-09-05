import Ember from 'ember';
import ClusterRoute from './cluster-route-base';

const { inject } = Ember;

export default ClusterRoute.extend({
  wizard: inject.service(),

  activate() {
    // always start from idle instead of using the current state
    this.get('wizard').transitionTutorialMachine('idle', 'INIT');
    this.get('wizard').set('initEvent', 'START');
  },
});
