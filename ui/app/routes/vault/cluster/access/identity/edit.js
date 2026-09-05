/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import GroupIdentityForm from 'vault/forms/identity/group';
import EntityIdentityForm from 'vault/forms/identity/entity';

export default class IdentityEditRoute extends Route {
  @service api;
  @service capabilities;

  async model(params) {
    const identityType = this.modelFor('vault.cluster.access.identity');
    const canCreatePolicy = await this.capabilities.for('policies').canCreate;
    const identityCapabilities = await this.capabilities.for('identityCapabilities', {
      identityType,
      id: params.item_id,
    });

    const methodType = identityType === 'entity' ? 'entityReadById' : 'groupReadById';
    const { data } = await this.api.identity[methodType](params.item_id);
    const form =
      identityType === 'group'
        ? new GroupIdentityForm(data, { isNew: false })
        : new EntityIdentityForm(data, { isNew: false });

    return {
      canCreatePolicy,
      canDelete: identityCapabilities?.canDelete || false,
      form,
      identityType,
      itemId: params.item_id,
    };
  }
}
