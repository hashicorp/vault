/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

import type Store from '@ember-data/store';
import type LdapRoleModel from 'vault/models/ldap/role';
import type Controller from '@ember/controller';
import type Transition from '@ember/routing/transition';
import type { Breadcrumb } from 'vault/vault/app-types';
import type AdapterError from 'ember-data/adapter'; // eslint-disable-line ember/use-ember-data-rfc-395-imports

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
interface LdapRoleCredentialsRouteModel {
  credentials: undefined | LdapStaticRoleCredentials | LdapDynamicRoleCredentials;
  error: undefined | AdapterError;
}
interface LdapRoleCredentialsController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapRoleCredentialsRouteModel;
}

export default class LdapRoleCredentialsRoute extends Route {
  @service declare readonly store: Store;

  async model() {
    try {
      const role = this.modelFor('roles.role') as LdapRoleModel;
      const credentials = await role.fetchCredentials();
      return { credentials };
    } catch (error) {
      return { error };
    }
  }
  setupController(
    controller: LdapRoleCredentialsController,
    resolvedModel: LdapRoleCredentialsRouteModel,
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
