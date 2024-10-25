/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { withConfig } from 'core/decorators/fetch-secrets-engine-config';
import { hash } from 'rsvp';

import type StoreService from 'vault/services/store';
import type PaginationService from 'vault/services/pagination';
import type SecretMountPath from 'vault/services/secret-mount-path';
import type Transition from '@ember/routing/transition';
import type LdapRoleModel from 'vault/models/ldap/role';
import type SecretEngineModel from 'vault/models/secret-engine';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/vault/app-types';

interface LdapRolesIndexRouteModel {
  backendModel: SecretEngineModel;
  promptConfig: boolean;
  roles: Array<LdapRoleModel>;
}
interface LdapRolesController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: LdapRolesIndexRouteModel;
  pageFilter: string | undefined;
  page: number | undefined;
}

interface LdapRolesRouteParams {
  page?: string;
  pageFilter: string;
}

@withConfig('ldap/config')
export default class LdapRolesIndexRoute extends Route {
  @service declare readonly store: StoreService;
  @service declare readonly pagination: PaginationService;
  @service declare readonly secretMountPath: SecretMountPath;

  declare promptConfig: boolean;

  queryParams = {
    pageFilter: {
      refreshModel: true,
    },
    page: {
      refreshModel: true,
    },
  };

  // consolidate logic into lazyQuery function?
  model(params: LdapRolesRouteParams) {
    const backendModel = this.modelFor('application') as SecretEngineModel;
    const page = Number(params.page) || 1;
    return hash({
      backendModel,
      promptConfig: this.promptConfig,
      roles: this.pagination.lazyPaginatedQuery(
        'ldap/role',
        {
          backend: backendModel.id,
          page,
          pageFilter: params.pageFilter,
          responsePath: 'data.keys',
          skipCache: page === 1,
        },
        { adapterOptions: { showPartialError: true } }
      ),
    });
  }

  setupController(
    controller: LdapRolesController,
    resolvedModel: LdapRolesIndexRouteModel,
    transition: Transition
  ) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backendModel.id, route: 'overview' },
      { label: 'Roles' },
    ];
  }

  resetController(controller: LdapRolesController, isExiting: boolean) {
    if (isExiting) {
      controller.set('pageFilter', undefined);
      controller.set('page', undefined);
    }
  }
}
