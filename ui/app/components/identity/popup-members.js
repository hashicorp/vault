/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';

export default class IdentityPopupMembers extends Component {
  @service flashMessages;
  @tracked showConfirmModal = false;

  onSuccess(memberId) {
    if (this.args.onSuccess) {
      this.args.onSuccess();
    }
    this.flashMessages.success(`Successfully removed '${memberId}' from the group`);
  }
  onError(err, memberId) {
    if (this.args.onError) {
      this.args.onError();
    }
    const error = errorMessage(err);
    this.flashMessages.danger(`There was a problem removing '${memberId}' from the group - ${error}`);
  }

  transaction() {
    const members = this.args.model[this.args.groupArray];
    this.args.model[this.args.groupArray] = members.without(this.args.memberId);
    return this.args.model.save();
  }

  @action
  async removeGroup() {
    const memberId = this.args.memberId;
    try {
      await this.transaction();
      this.onSuccess(memberId);
    } catch (e) {
      this.onError(e, memberId);
    }
  }
}
