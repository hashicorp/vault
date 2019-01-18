import { resolve } from 'rsvp';
import Route from '@ember/routing/route';

const SUPPORTED_DYNAMIC_BACKENDS = ['ssh', 'aws', 'pki'];

export default Route.extend({
  templateName: 'vault/cluster/secrets/backend/credentials',

  backendModel() {
    return this.modelFor('vault.cluster.secrets.backend');
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
