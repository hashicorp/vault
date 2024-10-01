/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import UnsavedModelRoute from 'vault/mixins/unsaved-model-route';
import { service } from '@ember/service';

export default Route.extend(UnsavedModelRoute, {
  store: service(),

  model(params) {
    const itemType = this.modelFor('vault.cluster.access.identity');
    const modelType = `identity/${itemType}-alias`;
    return this.store.createRecord(modelType, {
      canonicalId: params.item_id,
    });
  },
});
