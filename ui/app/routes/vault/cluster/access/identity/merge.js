/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import MergeEntitiesForm from 'vault/forms/identity/merge-entities';

export default class VaultClusterAccessIdentityMergeRoute extends Route {
  @service router;

  beforeModel() {
    const itemType = this.modelFor('vault.cluster.access.identity');
    if (itemType !== 'entity') {
      return this.router.transitionTo('vault.cluster.access.identity');
    }
  }

  model() {
    return { form: new MergeEntitiesForm({}, { isNew: true }) };
  }
}
