/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';

export default Route.extend({
  replicationMode: service(),
  beforeModel() {
    this.replicationMode.setMode(null);
  },
  model() {
    return this.modelFor('application');
  },
});
