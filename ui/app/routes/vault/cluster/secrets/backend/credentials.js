import { resolve } from 'rsvp';
import Route from '@ember/routing/route';
import { getOwner } from '@ember/application';
import { inject as service } from '@ember/service';

const SUPPORTED_DYNAMIC_BACKENDS = ['ssh', 'aws', 'pki'];

export default Route.extend({
  templateName: 'vault/cluster/secrets/backend/credentials',
  pathHelp: service('path-help'),

  backendModel() {
    return this.modelFor('vault.cluster.secrets.backend');
  },

  beforeModel() {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    if (backend != 'ssh') {
      return;
    }
    let modelType = 'ssh-otp-credential';
    let owner = getOwner(this);
    return this.pathHelp.getNewModel(modelType, backend, owner);
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
