/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { inject as service } from '@ember/service';
import Controller from '@ember/controller';

export default Controller.extend({
  wizard: service(),
  actions: {
    onMountSuccess: function (type, path) {
      this.wizard.transitionFeatureMachine(this.wizard.featureState, 'CONTINUE', type);
      const transition = this.transitionToRoute('vault.cluster.settings.auth.configure', path);
      return transition.followRedirects();
    },
  },
});
