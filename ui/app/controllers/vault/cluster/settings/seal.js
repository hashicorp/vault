/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { inject as service } from '@ember/service';
import Controller from '@ember/controller';

export default Controller.extend({
  auth: service(),
  router: service(),
  version: service(),

  actions: {
    seal() {
      return this.model.cluster.store
        .adapterFor('cluster')
        .seal()
        .then(() => {
          this.model.cluster.get('leaderNode').set('sealed', true);
          this.auth.deleteCurrentToken();
          // Reset version so it doesn't show on footer
          this.version.version = null;
          return this.router.transitionTo('vault.cluster.unseal');
        });
    },
  },
});
