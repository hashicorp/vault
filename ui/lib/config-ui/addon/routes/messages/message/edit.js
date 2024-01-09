/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'ember-concurrency';

export default class MessagesMessageEditRoute extends Route {
  @service store;

  model() {
    const { id, authenticated } = this.paramsFor('messages.message');
    // console.log('authenticatedd', authenticated);
    return hash({
      message: this.store.queryRecord('config-ui/message', id),
      messages: this.store.query('config-ui/message', {
        authenticated,
      }),
    });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Messages', route: 'messages', query: { authenticated: resolvedModel.authenticated } },
      { label: 'Edit Message' },
    ];
  }
}
