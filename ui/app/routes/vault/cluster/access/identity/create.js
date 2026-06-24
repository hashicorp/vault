/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import EntityIdentityForm from 'vault/forms/identity/entity';
import GroupIdentityForm from 'vault/forms/identity/group';
import { fetchIdentityItems } from 'vault/utils/identity-helpers';

export default class IdentityCreateRoute extends Route {
  @service api;
  @service capabilities;

  async model() {
    const identityType = this.modelFor('vault.cluster.access.identity');
    const canCreatePolicy = (await this.capabilities.for('policies')).canCreate;

    const form =
      identityType === 'group'
        ? new GroupIdentityForm({ type: 'internal' }, { isNew: true })
        : new EntityIdentityForm({}, { isNew: true });

    // Fetch groups and entities in parallel, handling 404s gracefully
    const [groupsResult, entitiesResult] = await Promise.allSettled([
      fetchIdentityItems({ identityType: 'group', api: this.api }),
      fetchIdentityItems({ identityType: 'entity', api: this.api }),
    ]);

    const groups = groupsResult.status === 'fulfilled' ? groupsResult.value : [];
    const entities = entitiesResult.status === 'fulfilled' ? entitiesResult.value : [];

    return {
      canCreatePolicy,
      entities,
      form,
      groups,
      identityType,
    };
  }
}
