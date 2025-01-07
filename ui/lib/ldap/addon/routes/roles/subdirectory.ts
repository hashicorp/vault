/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import LdapRolesRoute from '../roles';
import { hash } from 'rsvp';
import { ldapBreadcrumbs, roleRoutes } from 'ldap/utils/ldap-breadcrumbs';

import type { Breadcrumb } from 'vault/vault/app-types';
import type Controller from '@ember/controller';
import type LdapRoleModel from 'vault/models/ldap/role';
import type SecretEngineModel from 'vault/models/secret-engine';
import type Transition from '@ember/routing/transition';

interface RouteModel {
  backendModel: SecretEngineModel;
  roleAncestry: { path_to_role: string; type: string };
  roles: Array<LdapRoleModel>;
}

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: RouteModel;
}

interface RouteParams {
  page?: string;
  pageFilter: string;
  path_to_role: string;
  type: string;
}

export default class LdapRolesSubdirectoryRoute extends LdapRolesRoute {
  queryParams = {
    pageFilter: {
      refreshModel: true,
    },
    page: {
      refreshModel: true,
    },
  };

  model(params: RouteParams) {
    const backendModel = this.modelFor('application') as SecretEngineModel;
    const { path_to_role, type } = params;
    const roleAncestry = { path_to_role, type };
    return hash({
      backendModel,
      roleAncestry,
      roles: this.lazyQuery(backendModel.id, params, { roleAncestry }),
    });
  }

  setupController(controller: RouteController, resolvedModel: RouteModel, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);
    const { backendModel, roleAncestry } = resolvedModel;

    const routeParams = (childResource: string) => {
      return [backendModel.id, roleAncestry.type, childResource];
    };

    const crumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: backendModel.id, route: 'overview' },
      { label: 'Roles', route: 'roles' },
      ...ldapBreadcrumbs(roleAncestry.path_to_role, routeParams, roleRoutes, true),
    ];

    // must call 'set' so breadcrumbs update as we navigate through directories
    controller.set('breadcrumbs', crumbs);
  }

  resetController(controller: RouteController, isExiting: boolean) {
    if (isExiting) {
      controller.set('pageFilter', undefined);
      controller.set('page', undefined);
    }
  }
}
