import { resolve } from 'rsvp';
import Route from '@ember/routing/route';
import { getOwner } from '@ember/application';
import { inject as service } from '@ember/service';
import { combineAttributes } from 'vault/utils/openapi-to-attrs';

const SUPPORTED_DYNAMIC_BACKENDS = ['ssh', 'aws', 'pki'];

export default Route.extend({
  templateName: 'vault/cluster/secrets/backend/credentials',
  pathHelp: service('path-help'),

  backendModel() {
    return this.modelFor('vault.cluster.secrets.backend');
  },

  modelType(action) {
    let types = {
      sign: 'pki-certificate-sign',
      issue: 'pki-certificate',
      signVerbatim: 'pki-certificate',
    };
    return types[action];
  },

  beforeModel() {
    const { action } = this.paramsFor(this.routeName);
    return this.buildModel(action);
  },

  buildModel(action) {
    debugger; //eslint-disable-line
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    let modelType = this.modelType(action);
    let name = `model:${modelType}`;
    let owner = getOwner(this);
    return this.pathHelp.getProps(modelType, backend).then(props => {
      let newModel = owner.factoryFor(name).class;
      if (owner.hasRegistration(name) && !newModel.merged) {
        //combine them
        let attrs = combineAttributes(newModel.attributes, props);
        debugger; //eslint-disable-line
        newModel = newModel.extend(attrs);
      } else {
        //generate a whole new model
      }
      newModel.reopenClass({ merged: true });
      owner.unregister(name);
      owner.register(name, newModel);
    });
  },

  model(params) {
    let role = params.secret;
    let backendModel = this.backendModel();
    let backendPath = backendModel.get('id');
    let backendType = backendModel.get('type');

    if (!SUPPORTED_DYNAMIC_BACKENDS.includes(backendModel.get('type'))) {
      return this.transitionTo('vault.cluster.secrets.backend.list-root', backendPath);
    }
    return resolve({
      backendPath,
      backendType,
      roleName: role,
    });
  },

  resetController(controller) {
    controller.reset();
  },
});
