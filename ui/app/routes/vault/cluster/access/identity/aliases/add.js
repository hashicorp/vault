/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import AliasIdentityForm from 'vault/forms/identity/alias';

export default class VaultClusterAccessIdentityAliasesAddRoute extends Route {
  model(params) {
    const identityType = this.modelFor('vault.cluster.access.identity');

    return {
      canonicalId: params.item_id,
      form: new AliasIdentityForm({ canonical_id: params.item_id }, { isNew: true }),
      identityType,
    };
  }
}
