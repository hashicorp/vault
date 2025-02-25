/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import { ModelFrom } from 'vault/vault/route';
import type Store from '@ember-data/store';
import type SecretMountPath from 'vault/services/secret-mount-path';

export type LdapRolesRoleRouteModel = ModelFrom<LdapRolesRoleRoute>;

export default class LdapRolesRoleRoute extends Route {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;

  model(params: { name: string; type: string }) {
    const backend = this.secretMountPath.currentPath;
    const { name, type } = params;
    return this.store.queryRecord('ldap/role', { backend, name, type });
  }
}
