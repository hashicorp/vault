/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class MessagesCreateRoute extends Route {
  @service store;

  queryParams = {
    authenticated: {
      refreshModel: true,
    },
  };

  model(params) {
    const { authenticated } = params;

    return this.store.createRecord('config-ui/message', { authenticated });
  }
}
