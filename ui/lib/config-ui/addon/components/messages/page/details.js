/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import errorMessage from 'vault/utils/error-message';

/**
 * @module Page::MessageDetails
 * Page::MessageDetails components are used to display a message
 * @example
 * ```js
 * <Page::MessageDetails @message={{this.message}}  />
 * ```
 * @param {model} message - message model
 */

export default class MessageDetails extends Component {
  @service('app-router') router;
  @service flashMessages;
  @service customMessages;
  @service namespace;
  @service pagination;

  @action
  async deleteMessage() {
    try {
      await this.args.message.destroyRecord(this.args.message.id);
      this.pagination.clearDataset('config-ui/message');
      this.router.transitionTo('vault.cluster.config-ui.messages');
      this.customMessages.fetchMessages(this.namespace.path);
      this.flashMessages.success(`Successfully deleted ${this.args.message.title}.`);
    } catch (e) {
      const message = errorMessage(e);
      this.flashMessages.danger(message);
    }
  }
}
