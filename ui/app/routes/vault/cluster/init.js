import { inject as service } from '@ember/service';
import ClusterRoute from './cluster-route-base';

export default ClusterRoute.extend({
  wizard: service(),

  activate() {
    // always start from idle instead of using the current state
    this.get('wizard').transitionTutorialMachine('idle', 'INIT');
    this.get('wizard').set('initEvent', 'START');
  },
});
