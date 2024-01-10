/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

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

  model(params) {
    const { authenticated } = params;
    return hash({
      message: this.store.createRecord('config-ui/message', {
        authenticated,
      }),
      messages: this.getMessages(authenticated),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Messages', route: 'messages', query: { authenticated: controller.authenticated } },
      { label: 'Create Message' },
    ];
  }
}
