/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';
import { tracked } from '@glimmer/tracking';

export default class IdentityPopupPolicy extends Component {
  @service flashMessages;
  @tracked showConfirmModal = false;

  onSuccess(policyName, modelId) {
    if (this.args.onSuccess) {
      this.args.onSuccess();
    }
    this.flashMessages.success(`Successfully removed '${policyName}' policy from ${modelId}`);
  }
  onError(err, policyName) {
    if (this.args.onError) {
      this.args.onError();
    }
    const error = errorMessage(err);
    this.flashMessages.danger(`There was a problem removing '${policyName}' policy - ${error}`);
  }

  transaction() {
    const policies = this.args.model.policies;
    this.args.model.policies = policies.without(this.args.policyName);
    return this.args.model.save();
  }

  @action
  async removePolicy() {
    const {
      policyName,
      model: { id },
    } = this.args;
    try {
      await this.transaction();
      this.onSuccess(policyName, id);
    } catch (e) {
      this.onError(e, policyName);
    }
  }
}
