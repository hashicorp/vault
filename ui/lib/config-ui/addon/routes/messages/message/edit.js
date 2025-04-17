/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import CustomMessage from 'vault/forms/custom-message';
import { decodeString } from 'core/utils/b64';

export default class MessagesMessageEditRoute extends Route {
  @service api;

  async model() {
    const { id } = this.paramsFor('messages.message');
    const data = await this.api.sys.uiConfigReadCustomMessage(id);
    const { keyInfo, keys } = await this.api.sys.uiConfigListCustomMessages(
      true,
      undefined,
      data.authenticated
    );
    return {
      message: new CustomMessage({ ...data, message: decodeString(data.message) }),
      messages: keys.map((id) => ({ ...keyInfo[id], id })),
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
