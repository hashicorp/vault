/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import GroupIdentityForm from 'vault/forms/identity/group';

export default class IdentityEditRoute extends Route {
  @service api;
  @service capabilities;

  async model(params) {
    const canCreatePolicy = await this.capabilities.for('policies').canCreate;
    const identityCapabilities = await this.capabilities.for('identityCapabilities', {
      identityType: 'group',
      id: params.item_id,
    });

    const { data } = await this.api.identity.groupReadById(params.item_id);
    const form = new GroupIdentityForm(data, { isNew: false });

    return {
      canCreatePolicy,
      canDelete: identityCapabilities?.canDelete || false,
      form,
      identityType: 'group',
      itemId: params.item_id,
    };
  }
}
