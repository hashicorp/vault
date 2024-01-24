/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

export default class IdentityPopupMembers extends Component {
  @service flashMessages;
  @tracked showConfirmModal = false;

  onSuccess() {
    if (this.args.onSuccess) {
      this.args.onSuccess();
    }
    this.flashMessages.success(`Successfully removed '${this.args.memberId}' from the group`);
  }
  onError(err) {
    if (this.args.onError) {
      this.args.onError(this.args.model, this.args.memberId);
    }
    const error = this.errorMessage(err);
    this.flashMessages.error(error);
  }

  errorMessage(e) {
    const error = e.errors ? e.errors.join(' ') : e.message;
    return `There was a problem removing '${this.args.memberId}' from the group - ${error}`;
  }

  transaction() {
    const members = this.args.model[this.args.groupArray];
    this.args.model[this.args.groupArray] = members.without(this.args.memberId);
    return this.args.model.save();
  }

  @action
  async removeGroup() {
    try {
      await this.transaction();
      this.onSuccess();
    } catch (e) {
      this.onError(e);
    }
  }
}
