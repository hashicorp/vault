import { inject as service } from '@ember/service';
import ClusterRoute from './cluster-route-base';

export default ClusterRoute.extend({
  wizard: service(),

  activate() {
    this.get('wizard').set('initEvent', 'UNSEAL');
    this.get('wizard').transitionTutorialMachine(this.get('wizard.currentState'), 'TOUNSEAL');
  },
});
