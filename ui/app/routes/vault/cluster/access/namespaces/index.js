/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
// eslint-disable-next-line ember/no-mixins
import UnloadModel from 'vault/mixins/unload-model-route';

export default Route.extend(UnloadModel, {
  store: service(),
  version: service(),

  beforeModel() {
    this.store.unloadAll('namespace');
    return this.version.fetchFeatures().then(() => {
      return this._super(...arguments);
    });
  },

  model() {
    if (this.version.hasNamespaces) {
      return this.store.findAll('namespace').catch((err) => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      });
    }
    return null;
  },

  actions: {
    willTransition(transition) {
      window.scrollTo(0, 0);
      if (!transition || transition.targetName !== this.routeName) {
        this.store.clearAllDatasets();
      }
      return true;
    },
    reload() {
      this.store.clearAllDatasets();
      this.refresh();
    },
  },
});
