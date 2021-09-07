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
    // if KV2 then we pull in specific attrs from the config endpoint saved on the secret-engine record and display them
    // and only set if they haven't already been set
    if (backend.isV2KV && backend.attrs[backend.attrs.length - 1].name !== 'maxVersions') {
      let secretEngineRecord = this.store.peekRecord('secret-engine', backend.id);
      // create objects like you would normally pull from the model
      let casRequired = {
        name: 'casRequired',
        options: {
          label: 'Check-and-Set required',
        },
      };
      let deleteVersionAfter = {
        name: 'deleteVersionAfter',
        options: {
          label: 'Delete version after',
        },
      };
      let maxVersions = {
        name: 'maxVersions',
        options: {
          label: 'Maximum versions',
        },
      };

      backend.attrs.pushObject(casRequired);
      backend.attrs.pushObject(deleteVersionAfter);
      backend.attrs.pushObject(maxVersions);

      backend.set('casRequired', secretEngineRecord.casRequired ? secretEngineRecord.casRequired : 'False');
      backend.set(
        'deleteVersionAfter',
        secretEngineRecord.deleteVersionAfter ? secretEngineRecord.deleteVersionAfter : 'Never delete'
      );
      backend.set('maxVersions', secretEngineRecord.maxVersions ? secretEngineRecord.maxVersions : 'Not set');
    }
    return backend;
  },
});
