/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { decodeString } from 'core/utils/b64';

export default class MessagesMessageDetailsRoute extends Route {
  @service api;
  @service capabilities;

  async model() {
    const { id } = this.paramsFor('messages.message');

    const requests = [
      this.api.sys.uiConfigReadCustomMessage(id),
      this.capabilities.for('customMessages', { id }),
    ];
    const [customMessage, capabilities] = await Promise.all(requests);
    customMessage.message = decodeString(customMessage.message);

    return {
      message: customMessage,
      capabilities,
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const { message } = resolvedModel;

    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      {
        label: 'Custom messages',
        route: 'messages',
        query: { authenticated: message.authenticated },
      },
      { label: message.title },
    ];
  }
}
