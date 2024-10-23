/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type Store from '@ember-data/store';
import type SecretMountPath from 'vault/services/secret-mount-path';
import { ModelFrom } from 'vault/vault/route';

interface LdapRoleRouteParams {
  name: string;
  type: string;
}

export type LdapRoleRouteModel = ModelFrom<LdapRoleRoute>;
export default class LdapRoleRoute extends Route {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;

  async fetchHierarchalRoles(backend: string, roleName: string, type: string) {
    const backendModel = this.modelFor('roles');
    try {
      const roles = await this.store.query(
        'ldap/role',
        { backend, roleName },
        { adapterOptions: { roleType: type } }
      );
      if (roles.length === 1) {
        const [role] = roles;
        const name = `${roleName}${role.name}`;
        return this.store.queryRecord('ldap/role', { backend, name, type });
      }
      return { roles, backendModel, roleName };
    } catch (error) {
      return error;
    }
  }

  model(params: LdapRoleRouteParams) {
    const backend = this.secretMountPath.currentPath;
    const { name: roleName, type } = params;
    if (roleName.endsWith('/')) {
      return this.fetchHierarchalRoles(backend, roleName, type);
    }
    return this.store.queryRecord('ldap/role', { backend, name: roleName, type });
  }
}
