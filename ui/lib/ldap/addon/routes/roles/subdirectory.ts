/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import LdapRolesRoute from '../roles';
import { ldapBreadcrumbs, roleRoutes } from 'ldap/utils/ldap-breadcrumbs';

import type { Breadcrumb } from 'vault/vault/app-types';
import type Controller from '@ember/controller';
import type LdapRoleModel from 'vault/models/ldap/role';
import type Transition from '@ember/routing/transition';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
import type { LdapApplicationModel } from '../application';

interface RouteModel {
  secretsEngine: SecretsEngineResource;
  roleAncestry: { path_to_role: string; type: string };
  roles: Array<LdapRoleModel>;
}

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: RouteModel;
}

interface RouteParams {
  page?: string;
  pageFilter?: string;
  path_to_role?: string;
  type?: string;
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

  async model(params: RouteParams) {
    const { page, pageFilter: filter } = params;
    const { secretsEngine } = this.modelFor('application') as LdapApplicationModel;
    const { path_to_role, type } = params as { path_to_role: string; type: string };
    const roleAncestry = { path_to_role, type };
    const { roles, capabilities } = await this.fetchRolesAndCapabilities(
      { page: Number(page) || 1, filter },
      roleAncestry
    );
    return {
      secretsEngine,
      roleAncestry,
      roles,
      capabilities,
    };
  }

  setupController(controller: RouteController, resolvedModel: RouteModel, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);
    const { secretsEngine, roleAncestry } = resolvedModel;

    const routeParams = (childResource: string) => {
      return [secretsEngine.id, roleAncestry.type, childResource];
    };

    const crumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: secretsEngine.id, route: 'overview' },
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
