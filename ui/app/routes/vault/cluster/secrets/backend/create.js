import { hash } from 'rsvp';
import { inject as service } from '@ember/service';
import EmberObject from '@ember/object';
import EditBase from './secret-edit';
import KeyMixin from 'vault/models/key-mixin';

var SecretProxy = EmberObject.extend(KeyMixin, {
  store: null,

  toModel() {
    return this.getProperties('id', 'secretData', 'backend');
  },

  createRecord(backend) {
    let backendModel = this.store.peekRecord('secret-engine', backend);
    return this.store.createRecord(backendModel.get('modelTypeForKV'), this.toModel());
  },

  willDestroy() {
    this.store = null;
  },
});

export default EditBase.extend({
  wizard: service(),
  createModel(transition, parentKey) {
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
    const key = transition.queryParams.initialKey || '';
    const model = SecretProxy.create({
      initialParentKey: parentKey,
      store: this.store,
    });

    if (key) {
      // have to set this after so that it will be
      // computed properly in the template (it's dependent on `initialParentKey`)
      model.set('keyWithoutParent', key);
    }
    return model;
  },

  model(params, transition) {
    const parentKey = params.secret ? params.secret : '';
    return hash({
      secret: this.createModel(transition, parentKey),
      capabilities: {},
    });
  },
});
