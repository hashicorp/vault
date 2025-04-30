/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';

/**
 * @module Page::MessageDetails
 * Page::MessageDetails components are used to display a message
 * @example
 * ```js
 * <Page::MessageDetails @message={{this.model.message}} @capabilities={{this.model.capabilities}}  />
 * ```
 * @param message
 * @param capabilities - capabilities for the message
 */

export default class MessageDetails extends Component {
  @service('app-router') router;
  @service flashMessages;
  @service customMessages;
  @service pagination;
  @service api;

  displayFields = ['active', 'type', 'authenticated', 'title', 'message', 'startTime', 'endTime', 'link'];

  @action
  async deleteMessage() {
    try {
      const { message } = this.args;
      await this.api.sys.uiConfigDeleteCustomMessage(message.id);
      this.router.transitionTo('vault.cluster.config-ui.messages');
      this.customMessages.fetchMessages();
      this.flashMessages.success(`Successfully deleted ${message.title}.`);
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.flashMessages.danger(message);
    }
  }
}
