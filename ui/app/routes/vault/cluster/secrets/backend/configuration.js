/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';

export default Route.extend({
  store: service(),
  async model() {
    const backend = this.modelFor('vault.cluster.secrets.backend');
    // TODO kv engine cleanup - this can be removed when KV has fully moved to separate ember engine and list view config details menu is refactored
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
        backend.set('maxVersions', backend.maxVersions ? backend.maxVersions : 'Not set');
      } else {
        // remove the default values from the model if they don't have read access otherwise it will display the defaults even if they've been set (because they error on returning config data)
        // normally would catch the config error in the secret-v2 adapter, but I need the functions to proceed, not stop. So we remove the values here.
        backend.set('casRequired', null);
        backend.set('deleteVersionAfter', null);
        backend.set('maxVersions', null);
      }
    }
    return backend;
  },
});
