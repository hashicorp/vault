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

export default class IdentityCreateRoute extends Route {
  @service api;
  @service capabilities;

  async model() {
    const canCreatePolicy = (await this.capabilities.for('policies')).canCreate;
    const form = new GroupIdentityForm({ type: 'internal' }, { isNew: true });

    // Fetch groups and entities in parallel, handling 404s gracefully
    const [groupsResult, entitiesResult] = await Promise.allSettled([
      this.api.identity.groupListById(IdentityApiGroupListByIdListEnum.TRUE),
      this.api.identity.entityListById(IdentityApiEntityListByIdListEnum.TRUE),
    ]);

    const groups = groupsResult.status === 'fulfilled' ? this.api.keyInfoToArray(groupsResult.value) : [];
    const entities =
      entitiesResult.status === 'fulfilled' ? this.api.keyInfoToArray(entitiesResult.value) : [];

    return {
      canCreatePolicy,
      entities,
      form,
      groups,
      identityType: 'group',
    };
  }
}
