/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { paginate } from 'core/utils/paginate-list';

export default class MessagesRoute extends Route {
  @service api;

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
    const active =
      {
        active: true,
        inactive: false,
      }[status] || undefined;

    try {
      const { keyInfo, keys } = await this.api.sys.uiConfigListCustomMessages(
        true,
        active,
        authenticated,
        type
      );
      // ids are in the keys array and can be mapped to the object in keyInfo
      // map and set id property on keyInfo object
      const data = keys.map((id) => ({ id, ...keyInfo[id] }));
      const messages = paginate(data, {
        page,
        pageSize: 2,
        filter: pageFilter,
        filterKey: 'title',
      });
      return { params, messages };
    } catch (e) {
      if (e.response?.status === 404) {
        return { params, messages: [] };
      }
      throw e;
    }

    // const messages = this.pagination
    //   .lazyPaginatedQuery('config-ui/message', {
    //     authenticated,
    //     pageFilter: filter,
    //     active,
    //     type,
    //     responsePath: 'data.keys',
    //     page: page || 1,
    //     size: 10,
    //   })
    //   .catch((e) => {
    //     if (e.httpStatus === 404) {
    //       return [];
    //     }
    //     throw e;
    //   });
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
