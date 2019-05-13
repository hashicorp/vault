import { hash } from 'rsvp';
import { inject as service } from '@ember/service';
import EditBase from './secret-edit';

let secretModel = (store, backend, key) => {
  let backendModel = store.peekRecord('secret-engine', backend);
  let modelType = backendModel.get('modelTypeForKV');
  if (modelType !== 'secret-v2') {
    let model = store.createRecord(modelType);
    model.set('id', key);
    return model;
  }
  let secret = store.createRecord(modelType);
  secret.set('engine', backendModel);
  let version = store.createRecord('secret-v2-version', {
    path: key,
  });
  secret.set('selectedVersion', version);
  return secret;
};

export default EditBase.extend({
  wizard: service(),
  createModel(transition) {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    const modelType = this.modelType(backend);
    if (modelType === 'role-ssh') {
      return this.store.createRecord(modelType, { keyType: 'ca' });
    }
    if (modelType !== 'secret' && modelType !== 'secret-v2') {
      if (this.get('wizard.featureState') === 'details' && this.get('wizard.componentState') === 'transit') {
        this.get('wizard').transitionFeatureMachine('details', 'CONTINUE', 'transit');
      }
      return this.store.createRecord(modelType);
    }
    return secretModel(this.store, backend, transition.queryParams.initialKey);
  },

  model(params, transition) {
    return hash({
      secret: this.createModel(transition),
      capabilities: {},
    });
  },
});
