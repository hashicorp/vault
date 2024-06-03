/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class MessagesMessageDetailsRoute extends Route {
  @service store;

  model() {
    const { id } = this.paramsFor('messages.message');

    return this.store.queryRecord('config-ui/message', id);
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Messages', route: 'messages', query: { authenticated: resolvedModel.authenticated } },
      { label: resolvedModel.title },
    ];
  }
}
