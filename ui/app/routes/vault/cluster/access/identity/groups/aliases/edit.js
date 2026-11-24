/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import AliasIdentityForm from 'vault/forms/identity/alias';

export default class IdentityAliasesEditRoute extends Route {
  @service api;
  @service capabilities;

  async model(params) {
    const { data } = await this.api.identity.groupReadAliasById(params.item_alias_id);

    // Check canDelete capability for this alias
    const identityCapabilities = await this.capabilities.for('identityCapabilities', {
      identityType: 'group',
      id: params.item_alias_id,
    });

    return {
      canDelete: identityCapabilities?.canDelete || false,
      canonicalId: params.item_alias_id,
      form: new AliasIdentityForm(data, { isNew: false }),
      identityType: 'group',
      name: data.name,
    };
  }
}
