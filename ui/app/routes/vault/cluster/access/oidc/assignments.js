/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import {
  IdentityApiEntityListByIdListEnum,
  IdentityApiGroupListByIdListEnum,
} from '@hashicorp/vault-client-typescript';

export default class OidcAssignmentsRoute extends Route {
  @service api;

  async model() {
    // entities and groups are needed in the details, edit and create routes
    // fetch them here and access them with modelFor in the child routes
    const [entitiesResult, groupsResult] = await Promise.allSettled([
      this.api.identity.entityListById(IdentityApiEntityListByIdListEnum.TRUE),
      this.api.identity.groupListById(IdentityApiGroupListByIdListEnum.TRUE),
    ]);

    return {
      entities: entitiesResult.status === 'fulfilled' ? this.api.keyInfoToArray(entitiesResult.value) : [],
      groups: groupsResult.status === 'fulfilled' ? this.api.keyInfoToArray(groupsResult.value) : [],
    };
  }
}
