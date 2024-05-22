/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Controller from '@ember/controller';

export default Controller.extend({
  auth: service(),
  router: service(),
  version: service(),
  store: service(),

  actions: {
    seal() {
      return this.model.cluster.store
        .adapterFor('cluster')
        .seal()
        .then(() => {
          this.store.peekAll('cluster')[0].reload();
          this.auth.deleteCurrentToken();
          // Reset version so it doesn't show on footer
          this.version.version = null;
          return this.router.transitionTo('vault.cluster.unseal');
        });
    },
  },
});
