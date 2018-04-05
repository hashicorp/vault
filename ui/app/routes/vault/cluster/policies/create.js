import Ember from 'ember';
import UnloadModelRoute from 'vault/mixins/unload-model-route';
import UnsavedModelRoute from 'vault/mixins/unsaved-model-route';

const { inject } = Ember;
export default Ember.Route.extend(UnloadModelRoute, UnsavedModelRoute, {
  version: inject.service(),
  model() {
    let policyType = this.policyType();

    if (!this.get('version.hasSentinel') && policyType !== 'acl') {
      return this.transitionTo('vault.cluster.policies', policyType);
    }
    return this.store.createRecord(`policy/${policyType}`, {});
  },

  setupController(controller) {
    this._super(...arguments);
    controller.set('policyType', this.policyType());
  },

  policyType() {
    return this.paramsFor('vault.cluster.policies').type;
  },
});
