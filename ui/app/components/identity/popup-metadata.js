/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';

export default class IdentityPopupMetadata extends Component {
  @service flashMessages;
  @tracked showConfirmModal = false;

  onSuccess(key) {
    if (this.args.onSuccess) {
      this.args.onSuccess();
    }
    this.flashMessages.success(`Successfully removed '${key}' from metadata`);
  }
  onError(err, key) {
    if (this.args.onError) {
      this.args.onError();
    }
    const error = errorMessage(err);
    this.flashMessages.danger(`There was a problem removing '${key}' from the metadata - ${error}`);
  }

  transaction() {
    const metadata = this.args.model.metadata;
    delete metadata[this.args.key];
    this.args.model.metadata = { ...metadata };
    return this.args.model.save();
  }

  @action
  async removeMetadata() {
    const key = this.args.key;
    try {
      await this.transaction();
      this.onSuccess(key);
    } catch (e) {
      this.onError(e, key);
    }
  }
}
