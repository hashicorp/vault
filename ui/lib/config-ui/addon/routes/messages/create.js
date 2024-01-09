/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'ember-concurrency';

export default class MessagesCreateRoute extends Route {
  @service store;

  queryParams = {
    authenticated: {
      refreshModel: true,
    },
  };

  model(params) {
    const { authenticated } = params;
    return hash({
      message: this.store.createRecord('config-ui/message', {
        authenticated,
      }),
      messages: this.store.query('config-ui/message', {
        authenticated,
      }),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Messages', route: 'messages', query: { authenticated: !!resolvedModel.authenticated } },
      { label: 'Create Message' },
    ];
  }
}
