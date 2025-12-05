/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import LdapRolesRoute from '../roles';
import { hash } from 'rsvp';

import type Transition from '@ember/routing/transition';
import type LdapRoleModel from 'vault/models/ldap/role';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/vault/app-types';
import type SecretsEngineResource from 'vault/resources/secrets/engine';
import { LdapApplicationModel } from '../application';

interface RouteModel {
  secretsEngine: SecretsEngineResource;
  promptConfig: boolean;
  roles: Array<LdapRoleModel>;
}

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  model: RouteModel;
}

export default class LdapRolesIndexRoute extends LdapRolesRoute {
  queryParams = {
    pageFilter: {
      refreshModel: true,
    },
    page: {
      refreshModel: true,
    },
  };

  model(params: { page?: string; pageFilter: string }) {
    const { secretsEngine, promptConfig } = this.modelFor('application') as LdapApplicationModel;
    return hash({
      secretsEngine,
      promptConfig,
      roles: this.lazyQuery(secretsEngine.id, params, { showPartialError: true }),
    });
  }

  setupController(controller: RouteController, resolvedModel: RouteModel, transition: Transition) {
    super.setupController(controller, resolvedModel, transition);

    controller.breadcrumbs = [
      { label: 'Secrets', route: 'secrets', linkExternal: true },
      { label: resolvedModel.secretsEngine.id, route: 'overview' },
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
