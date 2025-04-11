/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class MessagesMessageEditRoute extends Route {
  @service api;

  async model() {
    const { id } = this.paramsFor('messages.message');
    const message = await this.api.sys.uiConfigReadCustomMessage(id);
    const { keyInfo = {} } = await this.api.sys.uiConfigListCustomMessages(
      true,
      undefined,
      message.authenticated
    );
    const messages = Object.values(keyInfo);
    return {
      message,
      messages: Object.values(messages),
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Messages', route: 'messages', query: { authenticated: resolvedModel.message.authenticated } },
      { label: 'Edit Message' },
    ];
  }
}
