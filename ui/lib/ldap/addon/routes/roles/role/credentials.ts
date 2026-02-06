/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ldapBreadcrumbs, roleRoutes } from 'ldap/utils/ldap-breadcrumbs';

import { ModelFrom } from 'vault/route';
import type ApiService from 'vault/services/api';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type { LdapRolesRoleRouteModel } from '../role';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/app-types';
import { LdapDynamicRoleCredentials, LdapStaticRoleCredentials } from 'vault/secrets/ldap';

export type LdapRoleCredentialsRouteModel = ModelFrom<LdapRoleCredentialsRoute>;

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapRoleCredentialsRouteModel;
}

export default class LdapRoleCredentialsRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly secretMountPath: SecretMountPath;

  async model() {
    try {
      const { role } = this.modelFor('roles.role') as LdapRolesRoleRouteModel;
      const { name, type } = role;
      const { currentPath } = this.secretMountPath;

      let credentials: LdapStaticRoleCredentials | LdapDynamicRoleCredentials;
      const response =
        type === 'static'
          ? await this.api.secrets.ldapRequestStaticRoleCredentials(name, currentPath)
          : await this.api.secrets.ldapRequestDynamicRoleCredentials(name, currentPath);

      const data = response.data as object;
      if (type === 'static') {
        credentials = { ...data, type } as LdapStaticRoleCredentials;
      } else {
        const { lease_id, lease_duration, renewable } = response;
        credentials = {
          ...data,
          lease_id,
          lease_duration,
          renewable,
          type,
        } as unknown as LdapDynamicRoleCredentials;
      }

      return { credentials };
    } catch (error) {
      return { error };
    }
  }

  setupController(controller: RouteController, resolvedModel: LdapRoleCredentialsRouteModel) {
    super.setupController(controller, resolvedModel);

    const { role } = this.modelFor('roles.role') as LdapRolesRoleRouteModel;
    const { currentPath } = this.secretMountPath;
    const routeParams = (childResource: string) => {
      return [currentPath, role.type, childResource];
    };
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault', icon: 'vault', linkExternal: true },
      { label: 'Secrets engines', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'overview' },
      { label: 'Roles', route: 'roles' },
      ...ldapBreadcrumbs(role.name, routeParams, roleRoutes),
      { label: 'Credentials' },
    ];
  }
}
