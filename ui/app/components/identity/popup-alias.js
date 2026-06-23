/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import errorMessage from 'vault/utils/error-message';

export default class IdentityPopupAlias extends Component {
  @service flashMessages;
  @service api;
  @service router;
  @tracked showConfirmModal = false;

  onSuccess(type, id) {
    this.router.transitionTo('vault.cluster.access.identity.aliases');
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
      const methodType = identityType === 'group' ? 'groupDeleteAliasById' : 'entityDeleteAliasById';
      await this.api.identity[methodType](id);
      this.onSuccess(identityType, id);
    } catch (e) {
      this.onError(e, identityType, id);
    }
  }
}
