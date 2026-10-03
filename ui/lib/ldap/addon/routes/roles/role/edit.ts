/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ldapBreadcrumbs, roleRoutes } from 'ldap/utils/ldap-breadcrumbs';
import LdapStaticRoleForm from 'vault/forms/secrets/ldap/roles/static';
import LdapDynamicRoleForm from 'vault/forms/secrets/ldap/roles/dynamic';

import { ModelFrom } from 'vault/route';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/vault/app-types';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type { LdapRolesRoleRouteModel } from '../role';

export type LdapRoleEditRouteModel = ModelFrom<LdapRoleEditRoute>;

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapRoleEditRouteModel;
}

export default class LdapRoleEditRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  model() {
    const { role } = this.modelFor('roles.role') as LdapRolesRoleRouteModel;
    if (role.type === 'static') {
      const staticForm = new LdapStaticRoleForm(role, { isNew: false });
      return { staticForm };
    }
    const dynamicForm = new LdapDynamicRoleForm(role, { isNew: false });
    return { dynamicForm };
  }

  setupController(controller: RouteController, resolvedModel: LdapRoleEditRouteModel) {
    super.setupController(controller, resolvedModel);

    const currentPath = this.secretMountPath.currentPath;
    const { role } = this.modelFor('roles.role') as LdapRolesRoleRouteModel;

    const routeParams = (childResource: string) => {
      return [currentPath, role.type, childResource];
    };

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: currentPath, route: 'overview' },
      { label: 'Roles', route: 'roles' },
      ...ldapBreadcrumbs(role.name, routeParams, roleRoutes),
      { label: 'Edit' },
    ];
  }
}
