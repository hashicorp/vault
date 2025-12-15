/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import trimRight from 'vault/utils/trim-right';
import { tracked } from '@glimmer/tracking';
import { formatStanzas, PolicyStanza } from 'core/utils/code-generators/policy';

/**
 * @module PolicyForm
 * PolicyForm components are the forms to create and edit all types of policies. This is only the form, not the outlying layout, and expects that the form model is passed from the parent.
 *
 * @example
 *  <PolicyForm
 *    @model={{this.model}}
 *    @onSave={{transition-to "vault.cluster.policy.show" this.model.policyType this.model.name}}
 *    @onCancel={{transition-to "vault.cluster.policies.index"}}
 *    @isCompact={{false}}
 *  />
 * ```
 * @callback onCancel - callback triggered when cancel button is clicked
 * @callback onSave - callback triggered when save button is clicked. Passes saved model
 * @param {object} model - ember data model from createRecord
 * @param {boolean} isCompact - renders a compact version of the form component, such as when rendering in a modal (see policy-template.hbs)
 */

export default class PolicyFormComponent extends Component {
  @service flashMessages;

  editTypes = { visual: 'Visual editor', code: 'Code editor' };

  @tracked editType = 'visual';
  @tracked errorBanner = '';
  @tracked showFileUpload = false;
  @tracked showSwitchEditorsModal = false;
  @tracked showTemplateModal = false;
  @tracked stanzas = [new PolicyStanza()];

  constructor() {
    super(...arguments);
    // Only ACL policies support the visual editor
    this.editType = this.args.model.policyType === 'acl' ? 'visual' : 'code';
  }

  get hasPolicyDiff() {
    const { policy } = this.args.model;
    // Make sure policy has a value (if it's undefined, neither editor has been used)
    // Return true if there is a difference between stanzas and policy arg
    // which means the user has made changes using the code editor
    return policy && formatStanzas(this.stanzas) !== policy;
  }

  get visualEditorSupported() {
    const { model, isCompact } = this.args;
    return model.isNew && model.policyType === 'acl' && !isCompact;
  }

  @action
  confirmEditorSwitch() {
    // User has confirmed discarding changes so switch to "visual" editor
    this.editType = 'visual';
    this.showSwitchEditorsModal = false;
    // Reset this.args.model.policy to match visual editor stanzas
    this.setPolicy(formatStanzas(this.stanzas));
  }

  @action
  handleNameInput(event) {
    const { value } = event.target;
    this.setName(value);
  }

  @action
  handlePolicyChange({ policy, stanzas }) {
    this.setPolicy(policy);
    this.stanzas = stanzas;
  }

  @action
  handleRadioChange(event) {
    const { value } = event.target;
    // Users cannot make changes using the code editor and have those parsed BACK to the visual editor
    if (value === 'visual' && this.hasPolicyDiff) {
      // Open modal to confirm user wants to switch back to "visual" editor and lose changes
      this.showSwitchEditorsModal = true;
    } else {
      this.editType = value;
    }
  }

  @task
  *save(event) {
    event.preventDefault();
    try {
      const { name, policyType, isNew } = this.args.model;
      yield this.args.model.save();
      this.flashMessages.success(
        `${policyType.toUpperCase()} policy "${name}" was successfully ${isNew ? 'created' : 'updated'}.`
      );
      this.args.onSave(this.args.model);
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.errorBanner = message;
    }
  }

  @action
  setName(name) {
    this.args.model.name = name.toLowerCase();
  }

  @action
  setPolicyFromFile(fileInfo) {
    const { value, filename } = fileInfo;
    this.setPolicy(value);
    if (!this.args.model.name) {
      const trimmedFileName = trimRight(filename, ['.json', '.txt', '.hcl', '.policy']);
      this.setName(trimmedFileName);
    }
    this.showFileUpload = false;
    // Switch to the code editor if they've uploaded a policy
    this.editType = 'code';
  }

  @action
  setPolicy(policy) {
    this.args.model.policy = policy;
  }

  @action
  cancel() {
    const method = this.args.model.isNew ? 'unloadRecord' : 'rollbackAttributes';
    this.args.model[method]();
    this.args.onCancel();
  }
}
