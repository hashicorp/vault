import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';

export default Route.extend({
  wizard: service(),
  model() {
    let backend = this.modelFor('vault.cluster.secrets.backend');
    if (this.get('wizard.featureState') === 'list') {
      this.get('wizard').transitionFeatureMachine(
        this.get('wizard.featureState'),
        'CONTINUE',
        backend.get('type')
      );
    }
    return backend;
  },
});
