/**
 * Copyright IBM Corp. 2016, 2025
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
      const { key_info, keys } = await this.api.sys.uiConfigListCustomMessages(
        true,
        active,
        authenticated,
        type
      );
      // ids are in the keys array and can be mapped to the object in key_info
      // map and set id property on key_info object
      const data = keys.map((id) => {
        const { start_time, end_time, ...message } = key_info[id];
        // dates returned from list endpoint are strings -- convert to date
        return {
          id,
          ...message,
          start_time: start_time ? new Date(start_time) : start_time,
          end_time: end_time ? new Date(end_time) : end_time,
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

    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Custom messages' },
    ];
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
