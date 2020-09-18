import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';

export default Route.extend({
  wizard: service(),
  model() {
    let backend = this.modelFor('vault.cluster.secrets.backend');
    if (this.wizard.featureState === 'list') {
      this.wizard.transitionFeatureMachine(this.wizard.featureState, 'CONTINUE', backend.get('type'));
    }
    return backend;
  },
});
