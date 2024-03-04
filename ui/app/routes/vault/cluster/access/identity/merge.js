/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import UnloadModelRoute from 'vault/mixins/unload-model-route';
import { service } from '@ember/service';

export default Route.extend(UnloadModelRoute, {
  store: service(),
  router: service(),

  beforeModel() {
    const itemType = this.modelFor('vault.cluster.access.identity');
    if (itemType !== 'entity') {
      return this.router.transitionTo('vault.cluster.access.identity');
    }
  },

  model() {
    const modelType = `identity/entity-merge`;
    return this.store.createRecord(modelType);
  },
});
