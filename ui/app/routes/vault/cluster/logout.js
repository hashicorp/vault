import { computed } from '@ember/object';
import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import ModelBoundaryRoute from 'vault/mixins/model-boundary-route';

export default Route.extend(ModelBoundaryRoute, {
  auth: service(),
  controlGroup: service(),
  flashMessages: service(),
  console: service(),
  permissions: service(),

  modelTypes: computed(function() {
    return ['secret', 'secret-engine'];
  }),

  beforeModel() {
    this.get('auth').deleteCurrentToken();
    this.get('controlGroup').deleteTokens();
    this.get('console').set('isOpen', false);
    this.get('console').clearLog(true);
    this.clearModelCache();
    this.replaceWith('vault.cluster');
    this.get('flashMessages').clearMessages();
    this.get('permissions').reset();
  },
});
