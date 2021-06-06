import { inject as service } from '@ember/service';
import ClusterRoute from './cluster-route-base';

export default ClusterRoute.extend({
  wizard: service(),

  activate() {
    this.wizard.set('initEvent', 'UNSEAL');
    this.wizard.transitionTutorialMachine(this.wizard.currentState, 'TOUNSEAL');
  },
});
