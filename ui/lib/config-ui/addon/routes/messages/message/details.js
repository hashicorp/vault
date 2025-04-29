/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { decodeString } from 'core/utils/b64';
import { PATH_MAP } from 'core/utils/capabilities';

export default class MessagesMessageDetailsRoute extends Route {
  @service api;
  @service capabilities;

  async model() {
    const { id } = this.paramsFor('messages.message');

    const requests = [
      this.api.sys.uiConfigReadCustomMessage(id),
      this.capabilities.fetchPathCapabilities(`${PATH_MAP.customMessages}/${id}`),
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
      { label: 'Messages', route: 'messages', query: { authenticated: message.authenticated } },
      { label: message.title },
    ];
  }
}
