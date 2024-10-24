/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type PaginationService from 'vault/services/pagination';
import type Store from '@ember-data/store';
import type SecretMountPath from 'vault/services/secret-mount-path';
import { ModelFrom } from 'vault/vault/route';
import SecretEngineModel from 'vault/models/secret-engine';

interface LdapRoleRouteParams {
  name: string;
  page?: string;
  pageFilter: string;
  type: string;
}

export type LdapRoleRouteModel = ModelFrom<LdapRoleRoute>;
export default class LdapRoleRoute extends Route {
  @service declare readonly pagination: PaginationService;
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;

  async fetchHierarchalRoles(params: LdapRoleRouteParams) {
    const backendModel = this.modelFor('roles') as SecretEngineModel;
    const parentRole = { name: params.name, type: params.type };
    const page = Number(params.page) || 1;
    try {
      const roles = await this.pagination.lazyPaginatedQuery(
        'ldap/role',
        {
          backend: backendModel.id,
          page,
          pageFilter: params.pageFilter,
          responsePath: 'data.keys',
          skipCache: page === 1,
        },
        { adapterOptions: { parentRole } }
      );

      return { roles, backendModel, parentRole };
    } catch (error) {
      return error;
    }
  }

  model(params: LdapRoleRouteParams) {
    const backend = this.secretMountPath.currentPath;
    const { name: roleName, type } = params;
    if (roleName.endsWith('/')) {
      return this.fetchHierarchalRoles(params);
    }
    return this.store.queryRecord('ldap/role', { backend, name: roleName, type });
  }
}
