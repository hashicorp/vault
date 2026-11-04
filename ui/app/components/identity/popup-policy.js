/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { handlePartialUpdate, refreshIdentityShow } from 'vault/utils/identity-helpers';
import { tracked } from '@glimmer/tracking';

export default class IdentityPopupPolicy extends Component {
  @service flashMessages;
  @service api;
  @service router;
  @tracked showConfirmModal = false;

  onSuccess(policyName, modelId) {
    if (this.args.onSuccess) {
      this.args.onSuccess();
    }
    refreshIdentityShow({ router: this.router, model: this.args.model });
    this.flashMessages.success(`Successfully removed '${policyName}' policy from ${modelId}`);
  }
  async onError(err, policyName) {
    if (this.args.onError) {
      this.args.onError();
    }
    const { message: error } = await this.api.parseError(err);
    this.flashMessages.danger(`There was a problem removing '${policyName}' policy - ${error}`);
  }

  transaction() {
    const { model, policyName } = this.args;
    const policies = model.policies.filter((name) => name !== policyName);
    return handlePartialUpdate({ api: this.api, model, fields: { policies } });
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
      await this.onError(e, policyName);
    }
  }
}
