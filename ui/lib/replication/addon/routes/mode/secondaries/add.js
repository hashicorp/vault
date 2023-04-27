/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Base from '../../replication-base';

export default Base.extend({
  model() {
    return this.modelFor('mode.secondaries');
  },

  redirect(model) {
    const replicationMode = this.paramsFor('mode').replication_mode;
    if (!model.get(`${replicationMode}.isPrimary`) || !model.get('canAddSecondary')) {
      return this.transitionTo('mode', replicationMode);
    }
  },

  resetController(controller) {
    controller.reset();
  },
});
