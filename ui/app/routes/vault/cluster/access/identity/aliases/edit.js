/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import UnloadModelRoute from 'vault/mixins/unload-model-route';
import UnsavedModelRoute from 'vault/mixins/unsaved-model-route';
import { inject as service } from '@ember/service';

export default Route.extend(UnloadModelRoute, UnsavedModelRoute, {
  store: service(),

  model(params) {
    const itemType = this.modelFor('vault.cluster.access.identity');
    const modelType = `identity/${itemType}-alias`;
    return this.store.findRecord(modelType, params.item_alias_id);
  },
});
