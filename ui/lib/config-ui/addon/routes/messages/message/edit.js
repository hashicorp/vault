/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';

export default class MessagesMessageEditRoute extends Route {
  @service store;

  getMessages(authenticated = true) {
    return this.store.query('config-ui/message', { authenticated }).catch(() => []);
  }

  async model() {
    const { id } = this.paramsFor('messages.message');
    const message = await this.store.queryRecord('config-ui/message', id);
    const messages = await this.getMessages(message.authenticated);
    return hash({
      message,
      messages,
      hasSomeActiveModals:
        messages.length && messages?.some((message) => message.type === 'modal' && message.active),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Messages', route: 'messages', query: { authenticated: resolvedModel.message.authenticated } },
      { label: 'Edit Message' },
    ];
  }
}
