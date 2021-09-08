import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';

export default Route.extend({
  wizard: service(),
  store: service(),
  model() {
    let backend = this.modelFor('vault.cluster.secrets.backend');
    if (this.wizard.featureState === 'list') {
      this.wizard.transitionFeatureMachine(this.wizard.featureState, 'CONTINUE', backend.get('type'));
    }
    if (backend.isV2KV) {
      // design wants specific default to show that can't be set in the model
      backend.set('casRequired', backend.casRequired ? backend.casRequired : 'False');
      backend.set(
        'deleteVersionAfter',
        backend.deleteVersionAfter ? backend.deleteVersionAfter : 'Never delete'
      );
      backend.set('maxVersions', backend.maxVersions ? backend.maxVersions : 'Not set');
    }
    return backend;
  },
});
