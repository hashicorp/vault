/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';

export default Controller.extend({
  showLicenseError: false,

  actions: {
    transitionToCluster() {
      return this.model.reload().then(() => {
        return this.transitionToRoute('vault.cluster', this.model.name);
      });
    },

    reloadCluster() {
      return this.model.reload();
    },

    isUnsealed(data) {
      return data.sealed === false;
    },

    handleLicenseError() {
      this.set('showLicenseError', true);
    },
  },
});
