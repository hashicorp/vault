/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { hash } from 'rsvp';
import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default Route.extend({
  store: service(),

  model() {
    return hash({
      cluster: this.modelFor('vault.cluster'),
      seal: this.store.findRecord('capabilities', 'sys/seal'),
    });
  },
});
