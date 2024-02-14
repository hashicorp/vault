/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';

export default class MessagesRoute extends Route {
  @service store;

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
  };

  model(params) {
    const { authenticated, page, pageFilter } = params;
    const filter = pageFilter
      ? (dataset) => dataset.filter((item) => item?.title.toLowerCase().includes(pageFilter.toLowerCase()))
      : null;
    const messages = this.store
      .lazyPaginatedQuery('config-ui/message', {
        authenticated,
        pageFilter: filter,
        responsePath: 'data.keys',
        page: page || 1,
        size: 10,
      })
      .catch((e) => {
        if (e.httpStatus === 404) {
          return [];
        }
        throw e;
      });
    return hash({
      pageFilter,
      messages,
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    const label = controller.authenticated ? 'After User Logs In' : 'On Login Page';
    controller.breadcrumbs = [{ label: 'Messages' }, { label }];
  }

  resetController(controller, isExiting) {
    if (isExiting) {
      controller.set('pageFilter', null);
    }
  }
}
