/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

export default class IdentityPopupMetadata extends Component {
  @service flashMessages;
  @tracked showConfirmModal = false;

  onSuccess() {
    if (this.args.onSuccess) {
      this.args.onSuccess();
    }
    this.flashMessages.success(`Successfully removed '${this.args.key}' from metadata`);
  }
  onError(err) {
    if (this.args.onError) {
      this.args.onError(this.args.model, this.args.key);
    }
    const error = this.errorMessage(err);
    this.flashMessages.error(error);
  }

  errorMessage(e) {
    const error = e.errors ? e.errors.join(' ') : e.message;
    return `There was a problem removing '${this.args.key}' from the metadata - ${error}`;
  }

  transaction() {
    const metadata = this.args.model.metadata;
    delete metadata[this.args.key];
    this.args.model.metadata = { ...metadata };
    return this.args.model.save();
  }

  @action
  async removeMetadata() {
    try {
      await this.transaction();
      this.onSuccess();
    } catch (e) {
      this.onError(e);
    }
  }
}
