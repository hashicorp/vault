/**
 * Copyright (c) HashiCorp, Inc.
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
      const { keyInfo } = await this.api.sys.uiConfigListCustomMessages(true, undefined, authenticated);
      return Object.values(keyInfo);
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
        startTime: addDays(startOfDay(timestamp.now()), 1).toISOString(),
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
      { label: 'Messages', route: 'messages', query: { authenticated: !!resolvedModel.authenticated } },
      { label: 'Create Message' },
    ];
  }
}
