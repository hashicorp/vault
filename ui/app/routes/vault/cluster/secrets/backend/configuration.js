import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';

export default Route.extend({
  wizard: service(),
  store: service(),
  async model() {
    let backend = this.modelFor('vault.cluster.secrets.backend');
    if (this.wizard.featureState === 'list') {
      this.wizard.transitionFeatureMachine(this.wizard.featureState, 'CONTINUE', backend.get('type'));
    }
    // ARG TODO confirm backend.id is kv name not just kv.
    if (backend.isV2KV) {
      let canRead = await this.store
        .findRecord('capabilities', `${backend.id}/config`)
        .then(response => response.canRead);
      // only set these config params if they can read the config endpoint.
      if (canRead) {
        // design wants specific default to show that can't be set in the model
        backend.set('casRequired', backend.casRequired ? backend.casRequired : 'False');
        backend.set(
          'deleteVersionAfter',
          backend.deleteVersionAfter !== '0s' ? backend.deleteVersionAfter : 'Never delete'
        );
        backend.set('maxVersions', backend.maxVersions ? backend.maxVersions : 'Not set');
      } else {
        backend.set('casRequired', null);
        backend.set('deleteVersionAfter', null);
        backend.set('maxVersions', null);
      }
    }
    return backend;
  },
});
