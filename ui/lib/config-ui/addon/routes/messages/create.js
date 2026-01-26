/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import CustomMessage from 'vault/forms/custom-message';
import { addDays, startOfDay } from 'date-fns';
import timestamp from 'core/utils/timestamp';

export default class MessagesCreateRoute extends Route {
  @service api;

  queryParams = {
    authenticated: {
      refreshModel: true,
    },
  };

  async getMessages(authenticated) {
    try {
      const { key_info } = await this.api.sys.uiConfigListCustomMessages(true, undefined, authenticated);
      return Object.values(key_info);
    } catch {
      return [];
    }
  }

  async model(params) {
    const { authenticated } = params;
    const message = new CustomMessage(
      {
        authenticated,
        type: 'banner',
        start_time: addDays(startOfDay(timestamp.now()), 1).toISOString(),
      },
      { isNew: true }
    );

    const messages = await this.getMessages(authenticated);

    return {
      message,
      messages,
      authenticated,
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      {
        label: 'Custom messages',
        route: 'messages',
        query: { authenticated: !!resolvedModel.authenticated },
      },
      { label: 'Create message' },
    ];
  }
}
