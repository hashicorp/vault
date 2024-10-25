/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import LdapRolesRoute from '../roles';
import { hash } from 'rsvp';
import { ldapBreadcrumbs } from 'ldap/utils/ldap-breadcrumbs';

import type { Breadcrumb } from 'vault/vault/app-types';
import type Controller from '@ember/controller';
import type LdapRoleModel from 'vault/models/ldap/role';
import type SecretEngineModel from 'vault/models/secret-engine';
import type Transition from '@ember/routing/transition';

interface LdapRolesSubdirectoryRouteModel {
  backendModel: SecretEngineModel;
  parentRole: { path_to_role: string; type: string };
  roles: Array<LdapRoleModel>;
}

interface LdapRolesSubdirectoryController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapRolesSubdirectoryRouteModel;
}
interface LdapRolesSubdirectoryParams {
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

  model(params: LdapRolesSubdirectoryParams) {
    const backendModel = this.modelFor('application') as SecretEngineModel;
    const { path_to_role, type } = params;
    const parentRole = { path_to_role, type };
    return hash({
      backendModel,
      parentRole,
      roles: this.lazyQuery({ parentRole }, backendModel.id, params),
    });
  }

  setupController(
    controller: LdapRolesSubdirectoryController,
    resolvedModel: LdapRolesSubdirectoryRouteModel,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);
    const { backendModel, parentRole } = resolvedModel;
    const crumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: backendModel.id, route: 'overview' },
      { label: 'Roles', route: 'roles' },
      ...ldapBreadcrumbs(parentRole.path_to_role, parentRole.type, backendModel.id),
    ];

    // must call 'set' so breadcrumbs update as we navigate through directories
    controller.set('breadcrumbs', crumbs);
  }

  resetController(controller: LdapRolesSubdirectoryController, isExiting: boolean) {
    if (isExiting) {
      controller.set('pageFilter', undefined);
      controller.set('page', undefined);
    }
  }
}
