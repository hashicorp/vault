/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Base from '../../replication-base';

export default Base.extend({
  model() {
    return this.modelFor('mode.secondaries');
  },

  redirect(model) {
    const replicationMode = this.replicationMode;
    if (!model[replicationMode].isPrimary || !model.canRevokeSecondary) {
      return this.router.transitionTo('vault.cluster.replication.index');
    }
  },

  resetController(controller) {
    controller.reset();
  },
});
