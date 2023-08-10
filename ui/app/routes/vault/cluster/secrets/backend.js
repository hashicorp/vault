/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
export default Route.extend({
  store: service(),
  flashMessages: service(),
  secretMountPath: service(),
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
          return model.get('firstObject');
        }
      });
  },

  afterModel(model, transition) {
    const path = model && model.get('path');
    if (transition.targetName === this.routeName) {
      return this.replaceWith('vault.cluster.secrets.backend.list-root', path);
    }
  },
});
