/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import GroupIdentityForm from 'vault/forms/identity/group';
import {
  IdentityApiGroupListByIdListEnum,
  IdentityApiEntityListByIdListEnum,
} from '@hashicorp/vault-client-typescript';

export default class IdentityEditRoute extends Route {
  @service api;
  @service capabilities;

  async model(params) {
    const canCreatePolicy = await this.capabilities.for('policies').canCreate;
    const identityCapabilities = await this.capabilities.for('identityCapabilities', {
      identityType: 'group',
      id: params.item_id,
    });

    // Fetch the group data alongside available groups and entities in parallel.
    // groups and entities are needed to resolve UUIDs to display names in the
    // Member group IDs / Member entity IDs SearchSelect fields.
    const [groupReadResult, groupsResult, entitiesResult] = await Promise.allSettled([
      this.api.identity.groupReadById(params.item_id),
      this.api.identity.groupListById(IdentityApiGroupListByIdListEnum.TRUE),
      this.api.identity.entityListById(IdentityApiEntityListByIdListEnum.TRUE),
    ]);

    if (groupReadResult.status === 'rejected') {
      throw groupReadResult.reason;
    }

    const { data } = groupReadResult.value;
    const form = new GroupIdentityForm(data, { isNew: false });

    const groups = groupsResult.status === 'fulfilled' ? this.api.keyInfoToArray(groupsResult.value) : [];
    const entities =
      entitiesResult.status === 'fulfilled' ? this.api.keyInfoToArray(entitiesResult.value) : [];

    return {
      canCreatePolicy,
      canDelete: identityCapabilities?.canDelete || false,
      entities,
      form,
      groups,
      identityType: 'group',
      itemId: params.item_id,
    };
  }
}
