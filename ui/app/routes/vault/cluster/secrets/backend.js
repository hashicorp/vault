/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
export default Route.extend({
  flashMessages: service(),
  router: service(),
  secretMountPath: service(),
  store: service(),
  oldModel: null,

  model(params) {
    const { backend } = params;
    this.secretMountPath.update(backend);
    return this.store
      .query('secret-engine', {
        path: backend,
      })
      .then((model) => {
        if (model) {
          return model[0];
        }
      });
  },

  afterModel(model, transition) {
    const path = model && model.path;
    if (transition.targetName === this.routeName) {
      return this.router.replaceWith('vault.cluster.secrets.backend.list-root', path);
    }
  },
});
