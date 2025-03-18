/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import UnloadModelRoute from 'vault/mixins/unload-model-route';
import { service } from '@ember/service';
import { ROUTES } from 'vault/utils/routes';

export default Route.extend(UnloadModelRoute, {
  store: service(),
  router: service(),

  beforeModel() {
    const itemType = this.modelFor(ROUTES.VAULT_CLUSTER_ACCESS_IDENTITY);
    if (itemType !== 'entity') {
      return this.router.transitionTo(ROUTES.VAULT_CLUSTER_ACCESS_IDENTITY);
    }
  },

  model() {
    const modelType = `identity/entity-merge`;
    return this.store.createRecord(modelType);
  },
});
