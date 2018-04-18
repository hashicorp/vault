import Ember from 'ember';
import ModelBoundaryRoute from 'vault/mixins/model-boundary-route';

export default Ember.Route.extend(ModelBoundaryRoute, {
  auth: Ember.inject.service(),
  flashMessages: Ember.inject.service(),

  modelTypes: ['secret', 'secret-engine'],

  beforeModel() {
    this.get('auth').deleteCurrentToken();
    this.clearModelCache();
    this.replaceWith('vault.cluster');
    this.get('flashMessages').clearMessages();
  },
});
