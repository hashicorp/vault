import Ember from 'ember';
import UnloadModel from 'vault/mixins/unload-model-route';

const SUPPORTED_DYNAMIC_BACKENDS = ['ssh', 'aws', 'pki'];

export default Ember.Route.extend(UnloadModel, {
  templateName: 'vault/cluster/secrets/backend/credentials',

  backendModel() {
    return this.modelFor('vault.cluster.secrets.backend');
  },

  model(params) {
    const role = params.secret;
    const backendModel = this.backendModel();
    const backend = backendModel.get('id');

    if (!SUPPORTED_DYNAMIC_BACKENDS.includes(backendModel.get('type'))) {
      return this.transitionTo('vault.cluster.secrets.backend.list-root', backend);
    }
    return Ember.RSVP.resolve({
      backend,
      id: role,
      name: role,
    });
  },

  setupController(controller) {
    this._super(...arguments);
    controller.set('backend', this.backendModel());
  },

  resetController(controller) {
    controller.reset();
  },
});
