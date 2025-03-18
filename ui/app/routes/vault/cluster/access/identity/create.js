/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import UnsavedModelRoute from 'vault/mixins/unsaved-model-route';
import { service } from '@ember/service';
import { ROUTES } from 'vault/utils/routes';

export default Route.extend(UnsavedModelRoute, {
  store: service(),

  model() {
    const itemType = this.modelFor(ROUTES.VAULT_CLUSTER_ACCESS_IDENTITY);
    const modelType = `identity/${itemType}`;
    return this.store.createRecord(modelType);
  },
});
