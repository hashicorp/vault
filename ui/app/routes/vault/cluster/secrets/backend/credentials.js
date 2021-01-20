import { resolve } from 'rsvp';
import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

const SUPPORTED_DYNAMIC_BACKENDS = ['database', 'ssh', 'aws', 'pki'];

export default Route.extend({
  templateName: 'vault/cluster/secrets/backend/credentials',
  pathHelp: service('path-help'),
  store: service(),

  backendModel() {
    return this.modelFor('vault.cluster.secrets.backend');
  },

  beforeModel() {
    const { backend } = this.paramsFor('vault.cluster.secrets.backend');
    if (backend != 'ssh') {
      return;
    }
    let modelType = 'ssh-otp-credential';
    return this.pathHelp.getNewModel(modelType, backend);
  },

  model(params) {
    let role = params.secret;
    let backendModel = this.backendModel();
    let backendPath = backendModel.get('id');
    let backendType = backendModel.get('type');
    let roleType = params.roleType;

    if (!SUPPORTED_DYNAMIC_BACKENDS.includes(backendModel.get('type'))) {
      return this.transitionTo('vault.cluster.secrets.backend.list-root', backendPath);
    }
    return resolve({
      backendPath,
      backendType,
      roleName: role,
      roleType,
    });
  },

  resetController(controller) {
    controller.reset();
  },

  actions: {
    willTransition() {
      // we do not want to save any of the credential information in the store.
      // once the user navigates away from this page, remove all credential info.
      this.store.unloadAll('database/credential');
    },
  },
});
