/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import trimRight from 'vault/utils/trim-right';
import { tracked } from '@glimmer/tracking';

/**
 * @module PolicyForm
 * PolicyForm components are the forms to create and edit all types of policies. This is only the form, not the outlying layout, and expects that the form model is passed from the parent.
 *
 * @example
 *  <PolicyForm
 *    @model={{this.model}}
 *    @onSave={{transition-to "vault.cluster.policy.show" this.model.policyType this.model.name}}
 *    @onCancel={{transition-to "vault.cluster.policies.index"}}
 *    @renderPolicyExampleModal={{true}}
 *  />
 * ```
 * @callback onCancel - callback triggered when cancel button is clicked
 * @callback onSave - callback triggered when save button is clicked. Passes saved model
 * @param {object} model - ember data model from createRecord
 * @param {boolean} renderPolicyExampleModal - whether or not the policy form should render the modal containing the policy example
 */

export default class PolicyFormComponent extends Component {
  @service flashMessages;

  @tracked errorBanner = '';
  @tracked showFileUpload = false;
  @tracked showTemplateModal = false;

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
  setModelName({ target }) {
    this.args.model.name = target.value.toLowerCase();
  }

  @action
  setPolicyFromFile(fileInfo) {
    const { value, filename } = fileInfo;
    this.args.model.policy = value;
    if (!this.args.model.name) {
      const trimmedFileName = trimRight(filename, ['.json', '.txt', '.hcl', '.policy']);
      this.args.model.name = trimmedFileName.toLowerCase();
    }
    this.showFileUpload = false;
  }

  @action
  cancel() {
    const method = this.args.model.isNew ? 'unloadRecord' : 'rollbackAttributes';
    this.args.model[method]();
    this.args.onCancel();
  }
}
