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

  onSuccess(type, name) {
    const itemType = type === 'group' ? 'groups' : 'entities';
    this.router.transitionTo(`vault.cluster.access.identity.${itemType}.aliases.index`);
    this.flashMessages.success(`Successfully deleted alias: ${name}`);
  }
  onError(err, name) {
    if (this.args.onError) {
      this.args.onError();
    }
    const error = errorMessage(err);
    this.flashMessages.danger(`There was a problem deleting alias: ${name} - ${error}`);
  }

  @action
  async deleteAlias() {
    const { identityType } = this.args;
    const { id, name } = this.args.item;
    try {
      const methodType = identityType === 'group' ? 'groupDeleteAliasById' : 'entityDeleteAliasById';
      await this.api.identity[methodType](id);
      this.onSuccess(identityType, name);
    } catch (e) {
      this.onError(e, name);
    }
  }
}
