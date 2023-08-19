/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';

export default Controller.extend({
  queryParams: ['action'],
  action: '',
  actions: {
    onPromote() {
      this.transitionToRoute('vault.cluster.replication.mode.index', 'dr');
    },
  },
});
