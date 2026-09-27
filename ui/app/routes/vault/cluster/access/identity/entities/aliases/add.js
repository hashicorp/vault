/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import AliasIdentityForm from 'vault/forms/identity/alias';
import { fetchAliases } from 'vault/utils/identity-helpers';

export default class VaultClusterAccessIdentityAliasesAddRoute extends Route {
  @service api;

  async model(params) {
    const identityType = 'entity';
    const aliases = await fetchAliases({ identityType, api: this.api });

    return {
      aliases,
      canonicalId: params.item_id,
      form: new AliasIdentityForm({ canonical_id: params.item_id }, { isNew: true }),
      identityType,
    };
  }
}
