/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';

export default Route.extend({
  store: service(),
  async model() {
    const backend = this.modelFor('vault.cluster.secrets.backend');
    if (backend.isV2KV) {
      const canRead = await this.store
        .findRecord('capabilities', `${backend.id}/config`)
        .then((response) => response.canRead);
      // only set these config params if they can read the config endpoint.
      if (canRead) {
        // design wants specific default to show that can't be set in the model
        backend.set('casRequired', backend.casRequired ? backend.casRequired : 'False');
        backend.set(
          'deleteVersionAfter',
          backend.deleteVersionAfter !== '0s' ? backend.deleteVersionAfter : 'Never delete'
        );
      } else {
        // remove the default values from the model if they don't have read access otherwise it will display the defaults even if they've been set (because they error on returning config data)
        backend.set('casRequired', null);
        backend.set('deleteVersionAfter', null);
        backend.set('maxVersions', null);
      }
    }
    return backend;
  },
});
