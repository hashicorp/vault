/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { capitalize } from '@ember/string';
import { findAuthMethod } from 'vault/utils/mountable-auth-methods';

export default Route.extend({
  store: service(),

  model() {
    const { method } = this.paramsFor(this.routeName);
    return this.store.findAll('auth-method').then(() => {
      return this.store.peekRecord('auth-method', method);
    });
  },

  setupController(controller) {
    this._super(...arguments);
    const methodData = findAuthMethod(controller.model.type);
    // right now token is the only method that's not mountable
    const displayName = methodData ? methodData.displayName : capitalize(controller.model.type);
    controller.set('displayName', displayName);
  },
});
