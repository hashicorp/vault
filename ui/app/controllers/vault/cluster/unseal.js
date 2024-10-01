/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';

export default Controller.extend({
  router: service(),
  showLicenseError: false,

  actions: {
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
