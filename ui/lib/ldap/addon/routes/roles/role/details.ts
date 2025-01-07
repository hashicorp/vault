/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { ldapBreadcrumbs, roleRoutes } from 'ldap/utils/ldap-breadcrumbs';

import type { Breadcrumb } from 'vault/vault/app-types';
import type Controller from '@ember/controller';
import type LdapRoleModel from 'vault/models/ldap/role';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type Transition from '@ember/routing/transition';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapRoleModel;
}

export default class LdapRolesRoleDetailsRoute extends Route {
  @service declare readonly secretMountPath: SecretMountPath;

  setupController(controller: RouteController, resolvedModel: LdapRoleModel, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);

    const routeParams = (childResource: string) => {
      return [this.secretMountPath.currentPath, resolvedModel.type, childResource];
    };

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'overview' },
      { label: 'Roles', route: 'roles' },
      ...ldapBreadcrumbs(resolvedModel.name, routeParams, roleRoutes, true),
    ];
  }
}
