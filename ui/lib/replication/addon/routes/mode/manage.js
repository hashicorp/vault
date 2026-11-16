/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';

export default Route.extend({
  router: service('app-router'),
  model() {
    return this.modelFor('mode');
  },

  beforeModel() {
    const model = this.modelFor('mode');
    const replicationMode = this.paramsFor('mode').replication_mode;
    const modeModel = model[replicationMode];
    if (modeModel.replicationDisabled || modeModel.replicationUnsupported) {
      this.router.transitionTo('vault.cluster.replication.mode', replicationMode);
    }
  },
});
