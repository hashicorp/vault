import { hash } from 'rsvp';
import Route from '@ember/routing/route';
import UnloadModelRoute from 'vault/mixins/unload-model-route';

export default Route.extend(UnloadModelRoute, {
  beforeModel() {
    const params = this.paramsFor(this.routeName);
    let policyType = this.policyType();
    if (policyType === 'acl' && params.policy_name === 'root') {
      return this.transitionTo('vault.cluster.policies', 'acl');
    }
  },

  model(params) {
    let type = this.policyType();
    return hash({
      policy: this.store.findRecord(`policy/${type}`, params.policy_name),
      capabilities: this.store.findRecord('capabilities', `sys/policies/${type}/${params.policy_name}`),
    });
  },

  setupController(controller, model) {
    controller.setProperties({
      model: model.policy,
      capabilities: model.capabilities,
      policyType: this.policyType(),
    });
  },

  policyType() {
    return this.paramsFor('vault.cluster.policy').type;
  },
});
