/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { inject as service } from '@ember/service';
import ClusterRoute from './cluster-route-base';

export default ClusterRoute.extend({
  wizard: service(),

  activate() {
    // always start from idle instead of using the current state
    this.wizard.transitionTutorialMachine('idle', 'INIT');
    this.wizard.set('initEvent', 'START');
  },
});
