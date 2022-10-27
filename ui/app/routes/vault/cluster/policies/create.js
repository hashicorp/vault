import UnloadModelRecord from 'vault/routes/vault/cluster/unload-model-record';
import { inject as service } from '@ember/service';

export default class PoliciesCreateRoute extends UnloadModelRecord {
  @service store;
  @service version;
  @service wizard;

  model() {
    let policyType = this.policyType();
    if (
      policyType === 'acl' &&
      this.wizard.currentMachine === 'policies' &&
      this.wizard.featureState === 'idle'
    ) {
      this.wizard.transitionFeatureMachine(this.wizard.featureState, 'CONTINUE');
    }
    if (!this.version.hasSentinel && policyType !== 'acl') {
      return this.transitionTo('vault.cluster.policies', policyType);
    }
    return this.store.createRecord(`policy/${policyType}`, {});
  }

  setupController(controller) {
    super.setupController(...arguments);
    controller.set('policyType', this.policyType());
  }

  policyType() {
    return this.paramsFor('vault.cluster.policies').type;
  }
}
