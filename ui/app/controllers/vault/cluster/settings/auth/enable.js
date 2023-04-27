/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller from '@ember/controller';

export default Controller.extend({
  actions: {
    onMountSuccess: function (type, path) {
      const transition = this.transitionToRoute('vault.cluster.settings.auth.configure', path);
      return transition.followRedirects();
    },
  },
});
