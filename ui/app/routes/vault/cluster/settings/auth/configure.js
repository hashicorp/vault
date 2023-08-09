/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default Route.extend({
  store: service(),

  model() {
    const { method } = this.paramsFor(this.routeName);
    return this.store.findAll('auth-method').then(() => {
      return this.store.peekRecord('auth-method', method);
    });
  },
});
