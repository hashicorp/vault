/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class MessagesCreateRoute extends Route {
  @service store;

  queryParams = {
    authenticated: {
      refreshModel: true,
    },
  };

  async getMessages(authenticated) {
    try {
      return await this.store.query('config-ui/message', {
        authenticated,
      });
    } catch {
      return [];
    }
  }

  async model(params) {
    const { authenticated } = params;
    const message = this.store.createRecord('config-ui/message', {
      authenticated,
    });
    const messages = await this.getMessages(authenticated);

    return {
      message,
      messages,
      authenticated,
      hasSomeActiveModals:
        messages.length && messages?.some((message) => message.type === 'modal' && message.active),
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Messages', route: 'messages', query: { authenticated: !!resolvedModel.authenticated } },
      { label: 'Create Message' },
    ];
  }
}
