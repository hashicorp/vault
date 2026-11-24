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
    const [aliases, entity] = await Promise.all([
      fetchAliases({ identityType, api: this.api }),
      this.api.identity.entityReadById(params.item_id),
    ]);

    return {
      aliases,
      canonicalId: params.item_id,
      canonicalName: entity.data.name,
      form: new AliasIdentityForm({ canonical_id: params.item_id }, { isNew: true }),
      identityType,
    };
  }
}
