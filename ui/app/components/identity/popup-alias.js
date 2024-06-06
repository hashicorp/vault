/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import errorMessage from 'vault/utils/error-message';

export default class IdentityPopupAlias extends Component {
  @service flashMessages;
  @tracked showConfirmModal = false;

  onSuccess(type, id) {
    if (this.args.onSuccess) {
      this.args.onSuccess();
    }
    this.flashMessages.success(`Successfully deleted ${type}: ${id}`);
  }
  onError(err, type, id) {
    if (this.args.onError) {
      this.args.onError();
    }
    const error = errorMessage(err);
    this.flashMessages.danger(`There was a problem deleting ${type}: ${id} - ${error}`);
  }

  @action
  async deleteAlias() {
    const { identityType, id } = this.args.item;
    try {
      await this.args.item.destroyRecord();
      this.onSuccess(identityType, id);
    } catch (e) {
      this.onError(e, identityType, id);
    }
  }
}
