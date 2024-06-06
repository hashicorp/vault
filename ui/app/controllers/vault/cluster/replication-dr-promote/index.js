/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';

export default Controller.extend({
  router: service(),
  queryParams: ['action'],
  action: '',
  actions: {
    onPromote() {
      this.router.transitionTo('vault.cluster.replication.mode.index', 'dr');
    },
  },
});
