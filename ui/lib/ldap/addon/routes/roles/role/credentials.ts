/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ldapBreadcrumbs, roleRoutes } from 'ldap/utils/ldap-breadcrumbs';

import type Store from '@ember-data/store';
import type LdapRoleModel from 'vault/models/ldap/role';
import type Controller from '@ember/controller';
import type Transition from '@ember/routing/transition';
import type { Breadcrumb } from 'vault/vault/app-types';
import type AdapterError from 'ember-data/adapter'; // eslint-disable-line ember/use-ember-data-rfc-395-imports

export interface StaticCredentials {
  dn: string;
  last_vault_rotation: string;
  password: string;
  last_password: string;
  rotation_period: number;
  ttl: number;
  username: string;
  type: string;
}
export interface DynamicCredentials {
  distinguished_names: Array<string>;
  password: string;
  username: string;
  lease_id: string;
  lease_duration: string;
  renewable: boolean;
  type: string;
}
interface RouteModel {
  credentials: undefined | StaticCredentials | DynamicCredentials;
  error: undefined | AdapterError;
}
interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: RouteModel;
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
  setupController(controller: RouteController, resolvedModel: RouteModel, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);

    const role = this.modelFor('roles.role') as LdapRoleModel;
    const routeParams = (childResource: string) => {
      return [role.backend, role.type, childResource];
    };
    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: role.backend, route: 'overview' },
      { label: 'Roles', route: 'roles' },
      ...ldapBreadcrumbs(role.name, routeParams, roleRoutes),
      { label: 'Credentials' },
    ];
  }
}
