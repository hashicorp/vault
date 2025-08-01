/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { paginate } from 'core/utils/paginate-list';

export default class MessagesRoute extends Route {
  @service api;
  @service capabilities;

  queryParams = {
    page: {
      refreshModel: true,
    },
    authenticated: {
      refreshModel: true,
    },
    pageFilter: {
      refreshModel: true,
    },
    status: {
      refreshModel: true,
    },
    type: {
      refreshModel: true,
    },
  };

  async model(params) {
    const { authenticated, page, pageFilter, status, type } = params;
    const active = {
      active: true,
      inactive: false,
    }[status];

    try {
      const { keyInfo, keys } = await this.api.sys.uiConfigListCustomMessages(
        true,
        active,
        authenticated,
        type
      );
      // ids are in the keys array and can be mapped to the object in keyInfo
      // map and set id property on keyInfo object
      const data = keys.map((id) => {
        const { startTime, endTime, ...message } = keyInfo[id];
        // dates returned from list endpoint are strings -- convert to date
        return {
          id,
          ...message,
          startTime: startTime ? new Date(startTime) : startTime,
          endTime: endTime ? new Date(endTime) : endTime,
        };
      });
      const messages = paginate(data, {
        page,
        pageSize: 2,
        filter: pageFilter,
        filterKey: 'title',
      });
      // fetch capabilities for each message path
      const paths = messages.map((message) => this.capabilities.pathFor('customMessages', message));
      const capabilities = await this.capabilities.fetch(paths);

      return { params, messages, capabilities };
    } catch (e) {
      if (e.response?.status === 404) {
        return { params, messages: [] };
      }
      throw e;
    }
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const label = controller.authenticated ? 'After User Logs In' : 'On Login Page';
    controller.breadcrumbs = [{ label: 'Messages' }, { label }];
  }

  resetController(controller, isExiting) {
    if (isExiting) {
      controller.set('pageFilter', null);
      controller.set('page', 1);
      controller.set('status', null);
      controller.set('type', null);
    }
  }
}
