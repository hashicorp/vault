/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import UnloadModel from 'vault/mixins/unload-model-route';

export default Route.extend(UnloadModel, {
  store: service(),
  version: service(),

  beforeModel() {
    return this.version.fetchFeatures().then(() => {
      return this._super(...arguments);
    });
  },

  model() {
    const type = 'control-group-config';
    return this.version.hasControlGroups
      ? this.store.findRecord(type, 'config').catch((e) => {
          // if you haven't saved a config, the API 404s, so create one here to edit and return it
          if (e.httpStatus === 404) {
            return this.store.createRecord(type, {
              id: 'config',
            });
          }
          throw e;
        })
      : null;
  },

  actions: {
    reload() {
      this.refresh();
    },
  },
});
