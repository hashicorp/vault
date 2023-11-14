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
    authenticated: {
      refreshModel: true,
    },
  };

  async model(params) {
    const { authenticated } = params;
    return hash({
      messages: this.store.query('config-ui/message', { authenticated }),
    });
  }
}
