/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

import type Store from '@ember-data/store';
import type SecretMountPath from 'vault/services/secret-mount-path';

interface LdapRoleRouteParams {
  name: string;
  type: string;
}

export default class LdapRoleRoute extends Route {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;

  model(params: LdapRoleRouteParams) {
    const backend = this.secretMountPath.currentPath;
    const { name, type } = params;
    return this.store.queryRecord('ldap/role', { backend, name, type });
  }
}
