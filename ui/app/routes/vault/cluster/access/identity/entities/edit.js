/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import EntityIdentityForm from 'vault/forms/identity/entity';

export default class IdentityEditRoute extends Route {
  @service api;
  @service capabilities;

  async model(params) {
    const canCreatePolicy = await this.capabilities.for('policies').canCreate;
    const identityCapabilities = await this.capabilities.for('identityCapabilities', {
      identityType: 'entity',
      id: params.item_id,
    });

    const { data } = await this.api.identity.entityReadById(params.item_id);
    const form = new EntityIdentityForm(data, { isNew: false });

    return {
      canCreatePolicy,
      canDelete: identityCapabilities?.canDelete || false,
      form,
      identityType: 'entity',
      itemId: params.item_id,
    };
  }
}
