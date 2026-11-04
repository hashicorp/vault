/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { handlePartialUpdate, refreshIdentityShow } from 'vault/utils/identity-helpers';

export default class IdentityItemDetails extends Component {
  @service api;
  @service flashMessages;
  @service router;

  @action
  async enable() {
    const { model } = this.args;
    try {
      await handlePartialUpdate({ api: this.api, model, fields: { disabled: false } });
      refreshIdentityShow({ router: this.router, model });
      this.flashMessages.success(`Successfully enabled entity: ${model.id}`);
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.flashMessages.danger(`There was a problem enabling the entity: ${model.id} - ${message}`);
    }
  }
}
