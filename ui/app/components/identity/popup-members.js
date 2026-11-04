/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { handlePartialUpdate, refreshIdentityShow } from 'vault/utils/identity-helpers';

export default class IdentityPopupMembers extends Component {
  @service flashMessages;
  @service api;
  @service router;
  @tracked showConfirmModal = false;

  onSuccess(memberId) {
    if (this.args.onSuccess) {
      this.args.onSuccess();
    }
    refreshIdentityShow({ router: this.router, model: this.args.model });
    this.flashMessages.success(`Successfully removed '${memberId}' from the group`);
  }
  async onError(err, memberId) {
    if (this.args.onError) {
      this.args.onError();
    }
    const { message: error } = await this.api.parseError(err);
    this.flashMessages.danger(`There was a problem removing '${memberId}' from the group - ${error}`);
  }

  transaction() {
    const { model, groupArray, memberId } = this.args;
    const members = model[groupArray].filter((id) => id !== memberId);
    return handlePartialUpdate({ api: this.api, model, fields: { [groupArray]: members } });
  }

  @action
  async removeGroup() {
    const memberId = this.args.memberId;
    try {
      await this.transaction();
      this.onSuccess(memberId);
    } catch (e) {
      await this.onError(e, memberId);
    }
  }
}
