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
    const replicationMode = this.paramsFor('mode').replication_mode;
    if (!model[replicationMode].isPrimary || !model.canAddSecondary) {
      return this.router.transitionTo('vault.cluster.replication.mode', replicationMode);
    }
  },

  resetController(controller) {
    controller.reset();
  },
});
