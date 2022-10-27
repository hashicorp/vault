import UnloadModelRecord from 'vault/routes/vault/cluster/unload-model-record';
import { inject as service } from '@ember/service';

export default class PolicyRoute extends UnloadModelRecord {
  @service store;

  beforeModel() {
    const params = this.paramsFor(this.routeName);
    let policyType = this.policyType();
    if (policyType === 'acl' && params.policy_name === 'root') {
      return this.transitionTo('vault.cluster.policies', 'acl');
    }
  }

  model(params) {
    let type = this.policyType();
    return this.store.findRecord(`policy/${type}`, params.policy_name);
  }

  policyType() {
    return this.paramsFor('vault.cluster.policies').type;
  }
}
