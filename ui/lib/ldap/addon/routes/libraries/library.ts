/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

import type Store from '@ember-data/store';
import type SecretMountPath from 'vault/services/secret-mount-path';

interface LdapLibraryRouteParams {
  name: string;
}

export default class LdapLibraryRoute extends Route {
  @service declare readonly store: Store;
  @service declare readonly secretMountPath: SecretMountPath;

  model(params: LdapLibraryRouteParams) {
    const backend = this.secretMountPath.currentPath;
    const { name } = params;
    return this.store.queryRecord('ldap/library', { backend, name });
  }
}
