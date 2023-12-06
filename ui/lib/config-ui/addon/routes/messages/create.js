/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class MessagesCreateRoute extends Route {
  @service store;

  model() {
    return this.store.createRecord('config-ui/message');
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'Messages', route: 'messages.index', query: { authenticated: false } },
      { label: 'Create Message' },
    ];
  }
}
