/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ldapBreadcrumbs, roleRoutes } from 'ldap/utils/ldap-breadcrumbs';

import type { Breadcrumb } from 'vault/vault/app-types';
import type Controller from '@ember/controller';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type Transition from '@ember/routing/transition';
import type { LdapRolesRoleRouteModel } from '../role';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapRolesRoleRouteModel;
}

export default class LdapRolesRoleDetailsRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  setupController(
    controller: RouteController,
    resolvedModel: LdapRolesRoleRouteModel,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);

    const routeParams = (childResource: string) => {
      return [this.secretMountPath.currentPath, resolvedModel.role.type, childResource];
    };

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'Roles', route: 'roles' },
      ...ldapBreadcrumbs(resolvedModel.role.name, routeParams, roleRoutes, true),
    ];
  }
}
