import Ember from 'ember';
import ModelBoundaryRoute from 'vault/mixins/model-boundary-route';

const { inject } = Ember;
export default Ember.Route.extend(ModelBoundaryRoute, {
  auth: inject.service(),
  flashMessages: inject.service(),
  console: inject.service(),

  modelTypes: ['secret', 'secret-engine'],

  beforeModel() {
    this.get('auth').deleteCurrentToken();
    this.get('console').set('isOpen', false);
    this.get('console').clearLog(true);
    this.clearModelCache();
    this.replaceWith('vault.cluster');
    this.get('flashMessages').clearMessages();
  },
});
