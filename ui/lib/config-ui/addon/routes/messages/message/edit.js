/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'ember-concurrency';

export default class MessagesMessageEditRoute extends Route {
  @service store;

  async getMessages(authenticated) {
    try {
      return await this.store.query('config-ui/message', {
        authenticated,
      });
    } catch {
      return [];
    }
  }

  async model() {
    const { id } = this.paramsFor('messages.message');
    const message = await this.store.queryRecord('config-ui/message', id);

    return hash({
      message,
      messages: this.getMessages(message.authenticated),
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
