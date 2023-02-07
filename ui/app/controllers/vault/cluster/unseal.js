/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { inject as service } from '@ember/service';
import Controller from '@ember/controller';

export default Controller.extend({
  wizard: service(),
  showLicenseError: false,

  actions: {
    transitionToCluster(resp) {
      return this.model.reload().then(() => {
        this.wizard.transitionTutorialMachine(this.wizard.currentState, 'CONTINUE', resp);
        return this.transitionToRoute('vault.cluster', this.model.name);
      });
    },

    setUnsealState(resp) {
      this.wizard.set('componentState', resp);
    },

    isUnsealed(data) {
      return data.sealed === false;
    },

    handleLicenseError() {
      this.set('showLicenseError', true);
    },
  },
});
