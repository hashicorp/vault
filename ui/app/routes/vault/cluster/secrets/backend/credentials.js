import Ember from 'ember';
import UnloadModel from 'vault/mixins/unload-model-route';

const SUPPORTED_DYNAMIC_BACKENDS = ['ssh', 'aws', 'pki'];

export default Ember.Route.extend(UnloadModel, {
  templateName: 'vault/cluster/secrets/backend/credentials',

  backendModel() {
    const backend = this.paramsFor('vault.cluster.secrets.backend').backend;
    return this.store.peekRecord('secret-engine', backend);
  },

  pathQuery(role, backend) {
    const type = this.backendModel().get('type');
    if (type === 'pki') {
      return `${backend}/issue/${role}`;
    }
    return `${backend}/creds/${role}`;
  },

  model(params) {
    const role = params.secret;
    const backendModel = this.backendModel();
    const backend = backendModel.get('id');

    if (!SUPPORTED_DYNAMIC_BACKENDS.includes(backendModel.get('type'))) {
      return this.transitionTo('vault.cluster.secrets.backend.list-root', backend);
    }
    return this.store
      .queryRecord('capabilities', { id: this.pathQuery(role, backend) })
      .then(capabilities => {
        if (!capabilities.get('canUpdate')) {
          return this.transitionTo('vault.cluster.secrets.backend.list-root', backend);
        }
        return Ember.RSVP.resolve({
          backend,
          id: role,
          name: role,
        });
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
