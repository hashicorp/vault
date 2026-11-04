/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { handlePartialUpdate, refreshIdentityShow } from 'vault/utils/identity-helpers';

export default class IdentityPopupMetadata extends Component {
  @service flashMessages;
  @service api;
  @service router;
  @tracked showConfirmModal = false;

  onSuccess(key) {
    if (this.args.onSuccess) {
      this.args.onSuccess();
    }
    refreshIdentityShow({ router: this.router, model: this.args.model });
    this.flashMessages.success(`Successfully removed '${key}' from metadata`);
  }
  async onError(err, key) {
    if (this.args.onError) {
      this.args.onError();
    }
    const { message: error } = await this.api.parseError(err);
    this.flashMessages.danger(`There was a problem removing '${key}' from the metadata - ${error}`);
  }

  transaction() {
    const { model, key } = this.args;
    const metadata = Object.fromEntries(Object.entries(model.metadata).filter(([k]) => k !== key));
    return handlePartialUpdate({ api: this.api, model, fields: { metadata } });
  }

  @action
  async removeMetadata() {
    const key = this.args.key;
    try {
      await this.transaction();
      this.onSuccess(key);
    } catch (e) {
      await this.onError(e, key);
    }
  }
}
