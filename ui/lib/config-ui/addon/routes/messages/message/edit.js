/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class MessagesMessageEditRoute extends Route {
  @service store;

  model() {
    const { id } = this.paramsFor('messages.message');

    return this.store.queryRecord('config-ui/message', id);
  }
}
