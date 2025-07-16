/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import LdapRolesRoute from '../roles';
import { service } from '@ember/service';
import { withConfig } from 'core/decorators/fetch-secrets-engine-config';
import { hash } from 'rsvp';

import type Store from '@ember-data/store';
import type Transition from '@ember/routing/transition';
import type LdapRoleModel from 'vault/models/ldap/role';
import type SecretEngineModel from 'vault/models/secret-engine';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/vault/app-types';

interface RouteModel {
  backendModel: SecretEngineModel;
  promptConfig: boolean;
  roles: Array<LdapRoleModel>;
}

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: RouteModel;
}

@withConfig('ldap/config')
export default class LdapRolesIndexRoute extends LdapRolesRoute {
  @service declare readonly store: Store; // necessary for @withConfig decorator

  declare promptConfig: boolean;

  queryParams = {
    pageFilter: {
      refreshModel: true,
    },
    page: {
      refreshModel: true,
    },
  };

  model(params: { page?: string; pageFilter: string }) {
    const backendModel = this.modelFor('application') as SecretEngineModel;
    return hash({
      backendModel,
      promptConfig: this.promptConfig,
      roles: this.lazyQuery(backendModel.id, params, { showPartialError: true }),
    });
  }

  setupController(controller: RouteController, resolvedModel: RouteModel, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.backendModel.id, route: 'overview' },
      { label: 'Roles' },
    ];
  }

  resetController(controller: RouteController, isExiting: boolean) {
    if (isExiting) {
      controller.set('pageFilter', undefined);
      controller.set('page', undefined);
    }
  }
}
