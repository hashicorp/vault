import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';

/**
 * @module PolicyForm
 * PolicyForm components are used to display the create and edit interface for all types of forms. This component is specifically for the policy pages for edit and create.
 *
 * @example
 *  <PolicyForm
 *    @model={{this.model}}
 *    @onSave={{transition-to "vault.cluster.policy.show" this.model.policyType this.model.name}}
 *    @onCancel={{transition-to "vault.cluster.policies.index"}}
 *  />
 * ```
 * @callback onCancel - callback triggered when cancel button is clicked
 * @callback onSave - callback triggered when save button is clicked
 * @param {object} model - ember data model from createRecord
 */

export default class PolicyFormComponent extends Component {
  @service wizard;
  @service flashMessages;

  @action onSave(policyModel) {
    const { name, policyType, isNew } = policyModel;
    this.flashMessages.success(
      `${policyType.toUpperCase()} policy "${name}" was successfully ${isNew ? 'created' : 'updated'}.`
    );
    if (this.wizard.featureState === 'create') {
      this.wizard.transitionFeatureMachine('create', 'CONTINUE', policyType);
    }
    this.args.onSave(policyModel);
  }
}
