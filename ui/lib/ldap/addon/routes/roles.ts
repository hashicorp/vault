/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type SecretMountPath from 'vault/services/secret-mount-path';
import type StoreService from 'vault/services/store';

// Base class for roles/index and roles/subdirectory routes
export default class LdapRolesRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;
  @service declare readonly store: StoreService;

  lazyQuery(backendId: string, params: { page?: string; pageFilter: string }, adapterOptions: object) {
    const page = Number(params.page) || 1;
    return this.store.lazyPaginatedQuery(
      'ldap/role',
      {
        backend: backendId,
        page,
        pageFilter: params.pageFilter,
        responsePath: 'data.keys',
        skipCache: page === 1,
      },
      { adapterOptions }
    );
  }
}
