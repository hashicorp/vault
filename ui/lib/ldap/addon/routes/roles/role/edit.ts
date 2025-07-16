/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { ldapBreadcrumbs, roleRoutes } from 'ldap/utils/ldap-breadcrumbs';

import type LdapRoleModel from 'vault/models/ldap/role';
import type Controller from '@ember/controller';
import type Transition from '@ember/routing/transition';
import type { Breadcrumb } from 'vault/vault/app-types';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapRoleModel;
}

export default class LdapRoleEditRoute extends Route {
  setupController(controller: RouteController, resolvedModel: LdapRoleModel, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);

    const routeParams = (childResource: string) => {
      return [resolvedModel.backend, resolvedModel.type, childResource];
    };

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backend, route: 'overview' },
      { label: 'Roles', route: 'roles' },
      ...ldapBreadcrumbs(resolvedModel.name, routeParams, roleRoutes),
      { label: 'Edit' },
    ];
  }
}
