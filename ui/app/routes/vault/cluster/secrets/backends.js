/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default Route.extend({
  store: service(),

  model() {
    // TODO LANDING PAGE: VAULT-17008 use peekAll to avoid a network request

    return this.store.query('secret-engine', {});
  },
});
