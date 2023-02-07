/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { inject as service } from '@ember/service';
import ClusterRoute from './cluster-route-base';

export default ClusterRoute.extend({
  wizard: service(),

  activate() {
    this.wizard.set('initEvent', 'UNSEAL');
    this.wizard.transitionTutorialMachine(this.wizard.currentState, 'TOUNSEAL');
  },
});
