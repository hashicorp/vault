/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

import type Store from '@ember-data/store';
import type LdapRoleModel from 'vault/models/ldap/role';
import type Controller from '@ember/controller';
import type Transition from '@ember/routing/transition';
import type { Breadcrumb } from 'vault/vault/app-types';

interface LdapRoleCredentialsController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapRoleModel;
}
export interface LdapStaticRoleCredentials {
  dn: string;
  last_vault_rotation: string;
  password: string;
  last_password: string;
  rotation_period: number;
  ttl: number;
  username: string;
  type: string;
}
export interface LdapDynamicRoleCredentials {
  distinguished_names: Array<string>;
  password: string;
  username: string;
  lease_id: string;
  lease_duration: string;
  renewable: boolean;
  type: string;
}

export default class LdapRoleCredentialsRoute extends Route {
  @service declare readonly store: Store;

  model() {
    const role = this.modelFor('roles.role') as LdapRoleModel;
    return role.fetchCredentials();
  }
  setupController(
    controller: LdapRoleCredentialsController,
    resolvedModel: LdapStaticRoleCredentials | LdapDynamicRoleCredentials,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);

    const role = this.modelFor('roles.role') as LdapRoleModel;
    controller.breadcrumbs = [
      { label: role.backend, route: 'overview' },
      { label: 'roles', route: 'roles' },
      { label: role.name, route: 'roles.role' },
      { label: 'credentials' },
    ];
  }
}
