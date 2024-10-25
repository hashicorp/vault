/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { hash } from 'rsvp';
import { ldapBreadcrumbs } from 'ldap/utils/ldap-breadcrumbs';
import { service } from '@ember/service';

import type { Breadcrumb } from 'vault/vault/app-types';
import type Controller from '@ember/controller';
import type LdapRoleModel from 'vault/models/ldap/role';
import type PaginationService from 'vault/services/pagination';
import type SecretEngineModel from 'vault/models/secret-engine';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type Transition from '@ember/routing/transition';

interface LdapRolesSubdirectoryRouteModel {
  backendModel: SecretEngineModel;
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

export default class LdapRolesSubdirectoryRoute extends Route {
  @service declare readonly pagination: PaginationService;
  @service declare readonly secretMountPath: SecretMountPath;

  model(params: LdapRolesSubdirectoryParams) {
    const backendModel = this.modelFor('application') as SecretEngineModel;
    const { path_to_role, type } = params;
    const page = Number(params.page) || 1;
    return hash({
      backendModel,
      roles: this.pagination.lazyPaginatedQuery(
        'ldap/role',
        {
          backend: backendModel.id,
          page,
          pageFilter: params.pageFilter,
          responsePath: 'data.keys',
          skipCache: page === 1,
        },
        { adapterOptions: { parentRole: { path_to_role, type } } }
      ),
    });
  }

  setupController(
    controller: LdapRolesSubdirectoryController,
    resolvedModel: LdapRolesSubdirectoryRouteModel,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);
    const { backendModel } = resolvedModel;
    const { path_to_role, type } = this.paramsFor('roles.subdirectory');
    const crumbs = [
      { label: backendModel.id, route: 'overview' },
      { label: 'Roles', route: 'roles' },
      ...ldapBreadcrumbs(path_to_role, type, backendModel.id),
    ];

    // must call 'set' so breadcrumbs update as we navigate through directories
    controller.set('breadcrumbs', crumbs);
  }
}
