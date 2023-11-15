/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { hash } from 'rsvp';
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
    let { authenticated, page } = params;
    page = page || 1;

    return hash({
      messages: this.store.lazyPaginatedQuery('config-ui/message', {
        authenticated,
        responsePath: 'data.keys',
        page,
      }),
      page,
    });
  }
}
