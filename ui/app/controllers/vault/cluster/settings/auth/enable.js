/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';

export default Controller.extend({
  router: service(),

  actions: {
    onMountSuccess: function (type, path) {
      const transition = this.router.transitionTo('vault.cluster.settings.auth.configure', path);
      return transition.followRedirects();
    },
  },
});
