/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class MessagesRoute extends Route {
  @service store;

  queryParams = {
    page: {
      refreshModel: true,
    },
    authenticated: {
      refreshModel: true,
    },
  };

  async model(params) {
    try {
      const { authenticated, page } = params;
      return await this.store.lazyPaginatedQuery('config-ui/message', {
        authenticated,
        responsePath: 'data.keys',
        page: page || 1,
      });
    } catch (e) {
      if (e.httpStatus === 404) {
        return [];
      }

      throw e;
    }
  }
}
